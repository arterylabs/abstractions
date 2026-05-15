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
	"testing"

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

// ---- mock task enqueuer ----

type mockTaskEnqueuer struct {
	enqueuedTasks []*asynq.Task
	returnErr     error
}

func (m *mockTaskEnqueuer) Enqueue(task *asynq.Task, _ ...asynq.Option) (*asynq.TaskInfo, error) {
	if m.returnErr != nil {
		return nil, m.returnErr
	}
	m.enqueuedTasks = append(m.enqueuedTasks, task)
	return &asynq.TaskInfo{ID: uuid.New().String()}, nil
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

// TestIndexRepositoryJob_StatusTransitions verifies pending → indexing → indexed.
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
	assert.Equal(s.T(), models.RepositoryStatusIndexed, updated.Status)
}

// TestIndexRepositoryJob_EnqueuesOneSummarizeJobPerFile verifies that a SummarizeFileJob
// task is enqueued for every file in the filtered tree and a delayed topic generation
// task is scheduled once indexing finishes.
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

	// One summarize task per file, plus one delayed generate_topics follow-up.
	assert.Len(s.T(), enqueuer.enqueuedTasks, len(files)+1)
	summarizeTasks := 0
	for _, t := range enqueuer.enqueuedTasks {
		if t.Type() == jobs.SummarizeFileJobID {
			summarizeTasks++
			continue
		}
		assert.Equal(s.T(), jobs.GenerateTopicsJobID, t.Type())
	}
	assert.Equal(s.T(), len(files), summarizeTasks)
}

// TestIndexRepositoryJob_ConcurrencyCap verifies that no more than MaxConcurrentSummarisations
// tasks are enqueued in a single run.
func (s *IndexRepositoryJobTestSuite) TestIndexRepositoryJob_ConcurrencyCap() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	s.createGitHubIntegration(workspace.ID, user.ID)
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	// Build more files than the concurrency cap
	files := make([]github_client.FileTreeItem, jobs.MaxConcurrentSummarisations+5)
	for i := range files {
		files[i] = github_client.FileTreeItem{Path: "file.go", SHA: "sha"}
	}
	fetcher := &mockFileTreeFetcher{files: files}
	enqueuer := &mockTaskEnqueuer{}
	job := jobs.NewIndexRepositoryJobWithDeps(s.DB, &config.Config{}, fetcher, enqueuer)

	err := job.Handler(s.Ctx, buildIndexTask(repo.ID))
	assert.NoError(s.T(), err)

	summarizeTasks := 0
	generateTopicsTasks := 0
	for _, t := range enqueuer.enqueuedTasks {
		switch t.Type() {
		case jobs.SummarizeFileJobID:
			summarizeTasks++
		case jobs.GenerateTopicsJobID:
			generateTopicsTasks++
		}
	}

	assert.LessOrEqual(s.T(), summarizeTasks, jobs.MaxConcurrentSummarisations)
	assert.Equal(s.T(), 1, generateTopicsTasks)
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

	summarizeTasks := 0
	for _, t := range enqueuer.enqueuedTasks {
		if t.Type() == jobs.SummarizeFileJobID {
			summarizeTasks++
		}
	}
	assert.Equal(s.T(), 0, summarizeTasks, "unchanged files must not be re-enqueued")
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

	summarizeTasks := 0
	for _, t := range enqueuer.enqueuedTasks {
		if t.Type() == jobs.SummarizeFileJobID {
			summarizeTasks++
		}
	}
	// Only "changed.go" and "brand_new.go" should be re-processed.
	assert.Equal(s.T(), 2, summarizeTasks)
}
