package tests

import (
	"abstractions/db/models"
	repository_service "abstractions/services/repository"
	"abstractions/tests/testutils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RepositoryServiceTestSuite struct {
	testutils.TestSuite
	seeder  *testutils.Seeder
	repoSvc *repository_service.Service
}

func TestRepositoryService(t *testing.T) {
	suite.Run(t, new(RepositoryServiceTestSuite))
}

func (s *RepositoryServiceTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.repoSvc = repository_service.NewRepositoryService(s.DB)
}

// TestConnectRepository_CreatesRowWithCorrectFields verifies that ConnectRepository inserts
// a row with the expected name, service_type, service_id, and status=pending.
func (s *RepositoryServiceTestSuite) TestConnectRepository_CreatesRowWithCorrectFields() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()

	req := &repository_service.ConnectRepositoryRequest{
		WorkspaceID:   workspace.ID,
		CreatedBy:     user.ID,
		Name:          "octocat/hello-world",
		ServiceType:   "GITHUB",
		ServiceID:     "123456",
		DefaultBranch: "main",
	}

	repo, err := s.repoSvc.ConnectRepository(s.Ctx, req)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), repo)
	assert.Equal(s.T(), "octocat/hello-world", repo.Name)
	assert.Equal(s.T(), "GITHUB", repo.ServiceType)
	assert.Equal(s.T(), "123456", repo.ServiceID)
	assert.Equal(s.T(), "pending", string(repo.Status))
	assert.Equal(s.T(), workspace.ID, repo.WorkspaceID)
	assert.Equal(s.T(), user.ID, repo.CreatedBy)
}

// TestConnectRepository_DuplicateReturnsError verifies that connecting the same repo twice
// returns ErrAlreadyConnected.
func (s *RepositoryServiceTestSuite) TestConnectRepository_DuplicateReturnsError() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()

	req := &repository_service.ConnectRepositoryRequest{
		WorkspaceID:   workspace.ID,
		CreatedBy:     user.ID,
		Name:          "octocat/dup-repo",
		ServiceType:   "GITHUB",
		ServiceID:     "789",
		DefaultBranch: "main",
	}

	_, err := s.repoSvc.ConnectRepository(s.Ctx, req)
	assert.NoError(s.T(), err)

	_, err = s.repoSvc.ConnectRepository(s.Ctx, req)
	assert.ErrorIs(s.T(), err, repository_service.ErrAlreadyConnected)
}

// TestConnectRepository_AfterDeleteAllowsReconnect verifies that a soft-deleted repository can be
// connected again using the same service identifier.
func (s *RepositoryServiceTestSuite) TestConnectRepository_AfterDeleteAllowsReconnect() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()

	req := &repository_service.ConnectRepositoryRequest{
		WorkspaceID:   workspace.ID,
		CreatedBy:     user.ID,
		Name:          "octocat/reconnect-repo",
		ServiceType:   "GITHUB",
		ServiceID:     "reconnect-1",
		DefaultBranch: "main",
	}

	original, err := s.repoSvc.ConnectRepository(s.Ctx, req)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), original)

	err = s.repoSvc.DeleteRepository(s.Ctx, original.ID, user.ID)
	assert.NoError(s.T(), err)

	reconnected, err := s.repoSvc.ConnectRepository(s.Ctx, req)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), reconnected)
	assert.NotEqual(s.T(), original.ID, reconnected.ID)
}

// TestConnectRepository_DefaultsServiceTypeToGitHub verifies that omitting ServiceType
// defaults to "GITHUB".
func (s *RepositoryServiceTestSuite) TestConnectRepository_DefaultsServiceTypeToGitHub() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()

	repo, err := s.repoSvc.ConnectRepository(s.Ctx, &repository_service.ConnectRepositoryRequest{
		WorkspaceID:   workspace.ID,
		CreatedBy:     user.ID,
		Name:          "octocat/default-type",
		ServiceID:     "42",
		DefaultBranch: "main",
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "GITHUB", repo.ServiceType)
}

// TestGetRepository_Found verifies that GetRepository returns the correct repo when it exists.
func (s *RepositoryServiceTestSuite) TestGetRepository_Found() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	seeded := s.seeder.CreateRepository(workspace.ID, user.ID)

	repo, err := s.repoSvc.GetRepository(s.Ctx, seeded.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), repo)
	assert.Equal(s.T(), seeded.ID, repo.ID)
	assert.Equal(s.T(), seeded.Name, repo.Name)
}

// TestGetRepository_NotFound verifies that GetRepository returns ErrRepositoryNotFound for an
// unknown ID.
func (s *RepositoryServiceTestSuite) TestGetRepository_NotFound() {
	_, err := s.repoSvc.GetRepository(s.Ctx, uuid.New())
	assert.ErrorIs(s.T(), err, repository_service.ErrRepositoryNotFound)
}

// TestGetStatus_ReturnsLiveProgress verifies that GetStatus returns the repository state
// and file stats derived from repository_files.
func (s *RepositoryServiceTestSuite) TestGetStatus_ReturnsLiveProgress() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	s.seeder.CreateFileIndex(repo.ID, "cmd/main.go", func(file *models.RepositoryFile) {
		file.Language = "Go"
		file.ProcessedAt = nil
	})
	s.seeder.CreateFileIndex(repo.ID, "web/app.ts", func(file *models.RepositoryFile) {
		file.Language = "TypeScript"
	})

	status, err := s.repoSvc.GetStatus(s.Ctx, workspace.ID, repo.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), status)
	assert.Equal(s.T(), "pending", status.Status)
	assert.Equal(s.T(), 2, status.FileCount)
	assert.Equal(s.T(), 2, status.TotalFiles)
	assert.Equal(s.T(), 1, status.SummarizedFiles)
	assert.Equal(s.T(), 50.0, status.ProgressPercentage)
	assert.Len(s.T(), status.LanguageBreakdown, 2)
}

// TestUpdateDescription verifies that UpdateDescription persists the new description.
func (s *RepositoryServiceTestSuite) TestUpdateDescription() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	seeded := s.seeder.CreateRepository(workspace.ID, user.ID)

	updated, err := s.repoSvc.UpdateDescription(s.Ctx, seeded.ID, "new description")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "new description", updated.Description)

	// Confirm via GetRepository that the change is persisted.
	fetched, err := s.repoSvc.GetRepository(s.Ctx, seeded.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "new description", fetched.Description)
}

// TestUpdateDescription_NotFound verifies that UpdateDescription returns ErrRepositoryNotFound for
// an unknown ID.
func (s *RepositoryServiceTestSuite) TestUpdateDescription_NotFound() {
	_, err := s.repoSvc.UpdateDescription(s.Ctx, uuid.New(), "desc")
	assert.ErrorIs(s.T(), err, repository_service.ErrRepositoryNotFound)
}

// TestUpdateSettings verifies that UpdateSettings persists send_day, send_hour, and
// newsletter_enabled.
func (s *RepositoryServiceTestSuite) TestUpdateSettings() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	seeded := s.seeder.CreateRepository(workspace.ID, user.ID)

	enabled := true
	updated, err := s.repoSvc.UpdateSettings(s.Ctx, seeded.ID, &repository_service.UpdateSettingsRequest{
		SendDay:           3,
		SendHour:          9,
		NewsletterEnabled: &enabled,
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 3, updated.SendDay)
	assert.Equal(s.T(), 9, updated.SendHour)
	assert.True(s.T(), updated.NewsletterEnabled)
}

// TestDeleteRepository_RemovesArtifacts verifies that DeleteRepository soft-deletes the repository
// and cascades the deletion to chunk_embeddings, file_index, and newsletter rows.
func (s *RepositoryServiceTestSuite) TestDeleteRepository_RemovesArtifacts() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	seeded := s.seeder.CreateRepository(workspace.ID, user.ID)

	// Seed derived artifacts that should be removed on delete.
	s.seeder.CreateFileIndex(seeded.ID, "src/main.go")
	s.seeder.CreateNewsletter(seeded.ID)
	s.seeder.CreateChunkEmbedding(seeded.ID, "src/main.go", 0, make([]float32, 1536))

	err := s.repoSvc.DeleteRepository(s.Ctx, seeded.ID, user.ID)
	assert.NoError(s.T(), err)

	// Repository should no longer be accessible.
	_, err = s.repoSvc.GetRepository(s.Ctx, seeded.ID)
	assert.ErrorIs(s.T(), err, repository_service.ErrRepositoryNotFound)
}

// TestDeleteRepository_NotFound verifies that DeleteRepository returns ErrRepositoryNotFound for an
// unknown ID.
func (s *RepositoryServiceTestSuite) TestDeleteRepository_NotFound() {
	err := s.repoSvc.DeleteRepository(s.Ctx, uuid.New(), uuid.New())
	assert.ErrorIs(s.T(), err, repository_service.ErrRepositoryNotFound)
}

// TestListRepositories_ReturnsOnlyWorkspaceRepos verifies workspace isolation.
func (s *RepositoryServiceTestSuite) TestListRepositories_ReturnsOnlyWorkspaceRepos() {
	ws1 := s.seeder.CreateWorkspace()
	ws2 := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()

	_, err := s.repoSvc.ConnectRepository(s.Ctx, &repository_service.ConnectRepositoryRequest{
		WorkspaceID: ws1.ID, CreatedBy: user.ID, Name: "a/repo1", ServiceType: "GITHUB", ServiceID: "1", DefaultBranch: "main",
	})
	assert.NoError(s.T(), err)

	_, err = s.repoSvc.ConnectRepository(s.Ctx, &repository_service.ConnectRepositoryRequest{
		WorkspaceID: ws2.ID, CreatedBy: user.ID, Name: "a/repo2", ServiceType: "GITHUB", ServiceID: "2", DefaultBranch: "main",
	})
	assert.NoError(s.T(), err)

	repos1, err := s.repoSvc.ListRepositories(s.Ctx, ws1.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), repos1, 1)
	assert.Equal(s.T(), "a/repo1", repos1[0].Name)

	repos2, err := s.repoSvc.ListRepositories(s.Ctx, ws2.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), repos2, 1)
	assert.Equal(s.T(), "a/repo2", repos2[0].Name)
}
