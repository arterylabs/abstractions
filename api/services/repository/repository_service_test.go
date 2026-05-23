package repository_service_test

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
	service *repository_service.Service
}

func TestRepositoryService(t *testing.T) {
	suite.Run(t, new(RepositoryServiceTestSuite))
}

func (s *RepositoryServiceTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.service = repository_service.NewRepositoryService(s.DB)
}

func (s *RepositoryServiceTestSuite) TestGetRepository_ReturnsCorrectRecord() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	fetched, err := s.service.GetRepository(s.Ctx, repo.ID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), fetched)
	assert.Equal(s.T(), repo.ID, fetched.ID)
	assert.Equal(s.T(), repo.Name, fetched.Name)
}

func (s *RepositoryServiceTestSuite) TestGetRepository_NotFound() {
	_, err := s.service.GetRepository(s.Ctx, uuid.New())

	assert.ErrorIs(s.T(), err, repository_service.ErrRepositoryNotFound)
}

func (s *RepositoryServiceTestSuite) TestConnectRepository_CreatesWithPendingStatus() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()

	req := &repository_service.ConnectRepositoryRequest{
		WorkspaceID:   workspace.ID,
		CreatedBy:     user.ID,
		Name:          "test-org/new-repo",
		Description:   "A test repo",
		ServiceType:   "GITHUB",
		ServiceID:     "gh-" + uuid.New().String(),
		DefaultBranch: "main",
	}

	repo, err := s.service.ConnectRepository(s.Ctx, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), repo)
	assert.Equal(s.T(), models.RepositoryStatusPending, repo.Status)
	assert.Equal(s.T(), workspace.ID, repo.WorkspaceID)
}

func (s *RepositoryServiceTestSuite) TestConnectRepository_DuplicateServiceID_ReturnsError() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	serviceID := "gh-" + uuid.New().String()

	req := &repository_service.ConnectRepositoryRequest{
		WorkspaceID:   workspace.ID,
		CreatedBy:     user.ID,
		Name:          "test-org/repo-a",
		ServiceType:   "GITHUB",
		ServiceID:     serviceID,
		DefaultBranch: "main",
	}

	_, err := s.service.ConnectRepository(s.Ctx, req)
	s.Require().NoError(err)

	req2 := &repository_service.ConnectRepositoryRequest{
		WorkspaceID:   workspace.ID,
		CreatedBy:     user.ID,
		Name:          "test-org/repo-a-copy",
		ServiceType:   "GITHUB",
		ServiceID:     serviceID, // same service ID
		DefaultBranch: "main",
	}

	_, err = s.service.ConnectRepository(s.Ctx, req2)
	assert.ErrorIs(s.T(), err, repository_service.ErrAlreadyConnected)
}

func (s *RepositoryServiceTestSuite) TestListRepositories_ReturnsWorkspaceScopedRepos() {
	ws1 := s.seeder.CreateWorkspace()
	ws2 := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()

	s.seeder.CreateRepository(ws1.ID, user.ID)
	s.seeder.CreateRepository(ws1.ID, user.ID)
	s.seeder.CreateRepository(ws2.ID, user.ID)

	repos, err := s.service.ListRepositories(s.Ctx, ws1.ID)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), repos, 2, "should only return repos for ws1")
}

func (s *RepositoryServiceTestSuite) TestConnectRepository_AppearsInList() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()

	req := &repository_service.ConnectRepositoryRequest{
		WorkspaceID:   workspace.ID,
		CreatedBy:     user.ID,
		Name:          "test-org/list-repo",
		ServiceType:   "GITHUB",
		ServiceID:     "gh-list-" + uuid.New().String(),
		DefaultBranch: "main",
	}

	created, err := s.service.ConnectRepository(s.Ctx, req)
	s.Require().NoError(err)

	list, err := s.service.ListRepositories(s.Ctx, workspace.ID)
	assert.NoError(s.T(), err)

	var found bool
	for _, r := range list {
		if r.ID == created.ID {
			found = true
			break
		}
	}
	assert.True(s.T(), found, "newly connected repo should appear in workspace list")
}
