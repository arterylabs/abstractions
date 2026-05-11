package tests

import (
	repository_service "abstractions/services/repository"
	"abstractions/tests/testutils"
	"testing"

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
