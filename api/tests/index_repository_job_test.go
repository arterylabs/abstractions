package tests

import (
	github_client "abstractions/clients/github"
	"abstractions/config"
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/tests/testutils"
	"abstractions/workers/jobs"
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	integration_service "abstractions/services/integrations"
	github_integration "abstractions/services/integrations/core/github"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ---- mock file tree fetcher ----

type mockFileTreeFetcher struct {
	files     []github_client.FileTreeItem
	commitSHA string
	err       error
}

func (m *mockFileTreeFetcher) GetFilteredFileTree(_ context.Context, _, _, _ string) ([]github_client.FileTreeItem, string, error) {
	return m.files, m.commitSHA, m.err
}

// ---- mock job enqueuer ----

type mockTaskEnqueuer struct {
	summarizeFileCalls  []string // tracked by filePath
	generateTopicsCalls []uuid.UUID
	returnErr           error
}

func (m *mockTaskEnqueuer) EnqueueSummarizeFile(_ uuid.UUID, filePath, _ string) error {
	if m.returnErr != nil {
		return m.returnErr
	}
	m.summarizeFileCalls = append(m.summarizeFileCalls, filePath)
	return nil
}

func (m *mockTaskEnqueuer) EnqueueGenerateTopics(repoID uuid.UUID) error {
	m.generateTopicsCalls = append(m.generateTopicsCalls, repoID)
	return nil
}

func (m *mockTaskEnqueuer) EnqueueGenerateNewsletter(_, _ uuid.UUID) error {
	return nil
}

// ---- test suite ----

type IndexRepositoryJobTestSuite struct {
	testutils.TestSuite
	seeder         *testutils.Seeder
	intSvc         *integration_service.Service
	repositoryRepo *repositories.RepositoryRepository
	fileIndexRepo  *repositories.RepositoryFileRepository
}

func TestIndexRepositoryJob(t *testing.T) {
	suite.Run(t, new(IndexRepositoryJobTestSuite))
}

func (s *IndexRepositoryJobTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.intSvc = integration_service.NewIntegrationService(s.DB)
	s.repositoryRepo = repositories.NewRepositoryRepository(s.DB)
	s.fileIndexRepo = repositories.NewRepositoryFileRepository(s.DB)
}

func (s *IndexRepositoryJobTestSuite) createGitHubIntegration(workspaceID, userID uuid.UUID) {
	_, err := s.intSvc.CreateIntegration(s.Ctx, &integration_service.CreateIntegrationRequest{
		Type:        github_integration.IntegrationTypeGitHub,
		WorkspaceID: workspaceID,
		Token:       "ghp_test_token",
		CreatedBy:   userID,
	})
	s.Require().NoError(err)
}

func buildIndexTask(repoID uuid.UUID) *asynq.Task {
	b, _ := json.Marshal(jobs.IndexRepositoryPayload{RepositoryID: repoID.String()})
	return asynq.NewTask(jobs.IndexRepositoryJobID, b)
}

// TestIndexRepositoryJob_StatusTransitions verifies pending → indexing when file processing is pending.
func (s *IndexRepositoryJobTestSuite) TestIndexRepositoryJob_StatusTransitions() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	s.createGitHubIntegration(workspace.ID, user.ID)
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	assert.Equal(s.T(), models.RepositoryStatusPending, repo.Status)

	fetcher := &mockFileTreeFetcher{
		files: []github_client.FileTreeItem{
			{Path: "main.go", SHA: "sha1"},
			{Path: "util.go", SHA: "sha2"},
		},
	}
	enqueuer := &mockTaskEnqueuer{}
	job := jobs.NewIndexRepositoryJobWithDeps(s.DB, &config.Config{}, fetcher, enqueuer)

	err := job.Handler(s.Ctx, buildIndexTask(repo.ID))
	assert.NoError(s.T(), err)

	updated, err := s.repositoryRepo.GetByID(s.Ctx, repo.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), models.RepositoryStatusIndexing, updated.Status)
}

// TestIndexRepositoryJob_EnqueuesOneSummarizeJobPerFile verifies that a SummarizeFileJob
// task is enqueued for every file in the filtered tree.
func (s *IndexRepositoryJobTestSuite) TestIndexRepositoryJob_EnqueuesOneSummarizeJobPerFile() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	s.createGitHubIntegration(workspace.ID, user.ID)
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	files := []github_client.FileTreeItem{
		{Path: "a.go", SHA: "sha_a"},
		{Path: "b.go", SHA: "sha_b"},
		{Path: "c.go", SHA: "sha_c"},
	}
	fetcher := &mockFileTreeFetcher{files: files}
	enqueuer := &mockTaskEnqueuer{}
	job := jobs.NewIndexRepositoryJobWithDeps(s.DB, &config.Config{}, fetcher, enqueuer)

	err := job.Handler(s.Ctx, buildIndexTask(repo.ID))
	assert.NoError(s.T(), err)

	// One summarize task per file. Topic generation waits for full processing completion.
	assert.Len(s.T(), enqueuer.summarizeFileCalls, len(files))
	assert.Empty(s.T(), enqueuer.generateTopicsCalls)
}

// TestIndexRepositoryJob_EnqueuesAllChangedFiles verifies that every changed/new file
// is enqueued for summarization in a single run.
func (s *IndexRepositoryJobTestSuite) TestIndexRepositoryJob_EnqueuesAllChangedFiles() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	s.createGitHubIntegration(workspace.ID, user.ID)
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	files := make([]github_client.FileTreeItem, 40)
	for i := range files {
		files[i] = github_client.FileTreeItem{Path: fmt.Sprintf("file_%d.go", i), SHA: fmt.Sprintf("sha_%d", i)}
	}
	fetcher := &mockFileTreeFetcher{files: files}
	enqueuer := &mockTaskEnqueuer{}
	job := jobs.NewIndexRepositoryJobWithDeps(s.DB, &config.Config{}, fetcher, enqueuer)

	err := job.Handler(s.Ctx, buildIndexTask(repo.ID))
	assert.NoError(s.T(), err)

	assert.Len(s.T(), enqueuer.summarizeFileCalls, len(files))
	assert.Empty(s.T(), enqueuer.generateTopicsCalls)
}

// TestIndexRepositoryJob_ErrorOnFetchTreeSetsErrorStatus verifies that when the file tree
// fetch fails the repository status is set to "error".
func (s *IndexRepositoryJobTestSuite) TestIndexRepositoryJob_ErrorOnFetchTreeSetsErrorStatus() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	s.createGitHubIntegration(workspace.ID, user.ID)
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	fetcher := &mockFileTreeFetcher{err: assert.AnError}
	enqueuer := &mockTaskEnqueuer{}
	job := jobs.NewIndexRepositoryJobWithDeps(s.DB, &config.Config{}, fetcher, enqueuer)

	err := job.Handler(s.Ctx, buildIndexTask(repo.ID))
	assert.Error(s.T(), err)

	updated, _ := s.repositoryRepo.GetByID(s.Ctx, repo.ID)
	assert.Equal(s.T(), models.RepositoryStatusError, updated.Status)
}

// TestIndexRepositoryJob_UpsertFileIndex verifies that all files appear in the file_index table.
func (s *IndexRepositoryJobTestSuite) TestIndexRepositoryJob_UpsertFileIndex() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	s.createGitHubIntegration(workspace.ID, user.ID)
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	files := []github_client.FileTreeItem{
		{Path: "x.go", SHA: "sha_x"},
		{Path: "y.go", SHA: "sha_y"},
	}
	fetcher := &mockFileTreeFetcher{files: files}
	job := jobs.NewIndexRepositoryJobWithDeps(s.DB, &config.Config{}, fetcher, &mockTaskEnqueuer{})

	err := job.Handler(s.Ctx, buildIndexTask(repo.ID))
	assert.NoError(s.T(), err)

	indexed, err := s.fileIndexRepo.ListByRepoID(s.Ctx, repo.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), indexed, len(files))
	for _, fi := range indexed {
		assert.Equal(s.T(), "go", fi.Language)
	}
}

// TestIndexRepositoryJob_SkipsUnchangedFiles verifies that files whose GitSHA matches the stored
// value are not enqueued for re-processing on subsequent index runs.
func (s *IndexRepositoryJobTestSuite) TestIndexRepositoryJob_SkipsUnchangedFiles() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	s.createGitHubIntegration(workspace.ID, user.ID)
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	// Seed existing repository_files with the same SHAs the fetcher will return.
	s.seeder.CreateFileIndex(repo.ID, "main.go", func(fi *models.RepositoryFile) {
		fi.GitSHA = "sha_unchanged"
	})
	s.seeder.CreateFileIndex(repo.ID, "util.go", func(fi *models.RepositoryFile) {
		fi.GitSHA = "sha_unchanged_2"
	})

	fetcher := &mockFileTreeFetcher{
		files: []github_client.FileTreeItem{
			{Path: "main.go", SHA: "sha_unchanged"},
			{Path: "util.go", SHA: "sha_unchanged_2"},
		},
	}
	enqueuer := &mockTaskEnqueuer{}
	job := jobs.NewIndexRepositoryJobWithDeps(s.DB, &config.Config{}, fetcher, enqueuer)

	err := job.Handler(s.Ctx, buildIndexTask(repo.ID))
	assert.NoError(s.T(), err)

	assert.Empty(s.T(), enqueuer.summarizeFileCalls, "unchanged files must not be re-enqueued")
}

// TestIndexRepositoryJob_ReprocessesChangedFiles verifies that only files with a new SHA
// (or new files not yet in repository_files) are enqueued for summarization.
func (s *IndexRepositoryJobTestSuite) TestIndexRepositoryJob_ReprocessesChangedFiles() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	s.createGitHubIntegration(workspace.ID, user.ID)
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	// Seed one file that hasn't changed and one that has changed.
	s.seeder.CreateFileIndex(repo.ID, "unchanged.go", func(fi *models.RepositoryFile) {
		fi.GitSHA = "sha_old"
	})
	s.seeder.CreateFileIndex(repo.ID, "changed.go", func(fi *models.RepositoryFile) {
		fi.GitSHA = "sha_old_2"
	})

	fetcher := &mockFileTreeFetcher{
		files: []github_client.FileTreeItem{
			{Path: "unchanged.go", SHA: "sha_old"},   // same SHA — skip
			{Path: "changed.go", SHA: "sha_new_2"},   // different SHA — enqueue
			{Path: "brand_new.go", SHA: "sha_brand"}, // new file — enqueue
		},
	}
	enqueuer := &mockTaskEnqueuer{}
	job := jobs.NewIndexRepositoryJobWithDeps(s.DB, &config.Config{}, fetcher, enqueuer)

	err := job.Handler(s.Ctx, buildIndexTask(repo.ID))
	assert.NoError(s.T(), err)

	// Only "changed.go" and "brand_new.go" should be re-processed.
	assert.Len(s.T(), enqueuer.summarizeFileCalls, 2)
}

// TestIndexRepositoryJob_ChangedFilesResetProcessingState verifies that changed files
// are marked pending and stale chunk rows are removed before re-enqueueing.
func (s *IndexRepositoryJobTestSuite) TestIndexRepositoryJob_ChangedFilesResetProcessingState() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	s.createGitHubIntegration(workspace.ID, user.ID)
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	now := time.Now()
	s.seeder.CreateFileIndex(repo.ID, "changed.go", func(fi *models.RepositoryFile) {
		fi.GitSHA = "old_sha"
		fi.ProcessedAt = &now
		fi.Tags = models.StringArray{"stale"}
	})
	s.seeder.CreateChunkEmbedding(repo.ID, "changed.go", 0, make([]float32, 1536))

	fetcher := &mockFileTreeFetcher{
		files: []github_client.FileTreeItem{
			{Path: "changed.go", SHA: "new_sha"},
		},
	}
	enqueuer := &mockTaskEnqueuer{}
	job := jobs.NewIndexRepositoryJobWithDeps(s.DB, &config.Config{}, fetcher, enqueuer)

	err := job.Handler(s.Ctx, buildIndexTask(repo.ID))
	assert.NoError(s.T(), err)

	fi, err := s.fileIndexRepo.GetByRepoAndPath(s.Ctx, repo.ID, "changed.go")
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), fi.ProcessedAt)
	assert.Empty(s.T(), fi.Tags)

	var chunkCount int64
	s.DB.Model(&models.ChunkEmbedding{}).
		Where("repo_id = ? AND repo_file_path = ?", repo.ID, "changed.go").
		Count(&chunkCount)
	assert.Equal(s.T(), int64(0), chunkCount)
}

// TestIndexRepositoryJob_MarksIndexedWhenNothingPending verifies that if no file
// requires re-processing, the repository is immediately marked indexed.
func (s *IndexRepositoryJobTestSuite) TestIndexRepositoryJob_MarksIndexedWhenNothingPending() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	s.createGitHubIntegration(workspace.ID, user.ID)
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	now := time.Now()
	s.seeder.CreateFileIndex(repo.ID, "main.go", func(fi *models.RepositoryFile) {
		fi.GitSHA = "sha_1"
		fi.ProcessedAt = &now
	})
	s.seeder.CreateFileIndex(repo.ID, "util.go", func(fi *models.RepositoryFile) {
		fi.GitSHA = "sha_2"
		fi.ProcessedAt = &now
	})

	fetcher := &mockFileTreeFetcher{
		files: []github_client.FileTreeItem{
			{Path: "main.go", SHA: "sha_1"},
			{Path: "util.go", SHA: "sha_2"},
		},
	}
	enqueuer := &mockTaskEnqueuer{}
	job := jobs.NewIndexRepositoryJobWithDeps(s.DB, &config.Config{}, fetcher, enqueuer)

	err := job.Handler(s.Ctx, buildIndexTask(repo.ID))
	assert.NoError(s.T(), err)

	updated, err := s.repositoryRepo.GetByID(s.Ctx, repo.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), models.RepositoryStatusIndexed, updated.Status)
	assert.Len(s.T(), enqueuer.generateTopicsCalls, 1)
}
