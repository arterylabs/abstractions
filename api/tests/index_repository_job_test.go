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
	files []github_client.FileTreeItem
	err   error
}

func (m *mockFileTreeFetcher) GetFilteredFileTree(_ context.Context, _, _, _ string) ([]github_client.FileTreeItem, error) {
	return m.files, m.err
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
	fileIndexRepo  *repositories.FileIndexRepository
}

func TestIndexRepositoryJob(t *testing.T) {
	suite.Run(t, new(IndexRepositoryJobTestSuite))
}

func (s *IndexRepositoryJobTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.intSvc = integration_service.NewIntegrationService(s.DB)
	s.repositoryRepo = repositories.NewRepositoryRepository(s.DB)
	s.fileIndexRepo = repositories.NewFileIndexRepository(s.DB)
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

	// Should have one task enqueued per file (all three fit within MaxConcurrentSummarisations)
	assert.Len(s.T(), enqueuer.enqueuedTasks, len(files))
	for _, t := range enqueuer.enqueuedTasks {
		assert.Equal(s.T(), jobs.SummarizeFileJobID, t.Type())
	}
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
	assert.LessOrEqual(s.T(), len(enqueuer.enqueuedTasks), jobs.MaxConcurrentSummarisations)
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
