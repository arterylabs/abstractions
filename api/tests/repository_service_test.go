package tests

import (
	integration_service "abstractions/services/integrations"
	github_integration "abstractions/services/integrations/core/github"
	repository_service "abstractions/services/repository"
	"abstractions/tests/testutils"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// mockWebhookCreator is a test double that avoids real GitHub API calls.
type mockWebhookCreator struct {
	webhookID    int64
	returnErr    error
	calledOwner  string
	calledRepo   string
	calledSecret string
}

func (m *mockWebhookCreator) CreateWebhook(_ context.Context, _, owner, repo, _, secret string) (int64, error) {
	m.calledOwner = owner
	m.calledRepo = repo
	m.calledSecret = secret
	return m.webhookID, m.returnErr
}

type RepositoryServiceTestSuite struct {
	testutils.TestSuite
	seeder      *testutils.Seeder
	repoSvc     *repository_service.Service
	intSvc      *integration_service.Service
	mockWebhook *mockWebhookCreator
}

func TestRepositoryService(t *testing.T) {
	suite.Run(t, new(RepositoryServiceTestSuite))
}

func (s *RepositoryServiceTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.intSvc = integration_service.NewIntegrationService(s.DB)
	s.mockWebhook = &mockWebhookCreator{webhookID: 99999}
	s.repoSvc = repository_service.NewRepositoryServiceWithWebhookCreator(s.DB, s.mockWebhook)
}

// createGitHubIntegration is a helper that seeds a GitHub integration with an encrypted token.
func (s *RepositoryServiceTestSuite) createGitHubIntegration(workspaceID, userID uuid.UUID) {
	_, err := s.intSvc.CreateIntegration(s.Ctx, &integration_service.CreateIntegrationRequest{
		Type:        github_integration.IntegrationTypeGitHub,
		WorkspaceID: workspaceID,
		Token:       "ghp_test_token",
		CreatedBy:   userID,
	})
	s.Require().NoError(err)
}

// TestConnectRepository_CreatesRowWithCorrectFields verifies that ConnectRepository inserts
// a row with the expected github_repo_id, full_name, status=pending, and webhook_id.
func (s *RepositoryServiceTestSuite) TestConnectRepository_CreatesRowWithCorrectFields() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	s.createGitHubIntegration(workspace.ID, user.ID)

	req := &repository_service.ConnectRepositoryRequest{
		WorkspaceID:   workspace.ID,
		UserID:        user.ID,
		GitHubRepoID:  123456,
		FullName:      "octocat/hello-world",
		DefaultBranch: "main",
	}

	repo, err := s.repoSvc.ConnectRepository(s.Ctx, req)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), repo)
	assert.Equal(s.T(), int64(123456), repo.GitHubRepoID)
	assert.Equal(s.T(), "octocat/hello-world", repo.FullName)
	assert.Equal(s.T(), "pending", string(repo.Status))
	assert.Equal(s.T(), int64(99999), *repo.WebhookID)
	assert.Equal(s.T(), workspace.ID, repo.WorkspaceID)
	assert.Equal(s.T(), user.ID, repo.UserID)
}

// TestConnectRepository_WebhookCalledWithCorrectOwnerAndRepo verifies the webhook client
// receives owner="octocat" and repo="hello-world" when full_name is "octocat/hello-world".
func (s *RepositoryServiceTestSuite) TestConnectRepository_WebhookCalledWithCorrectOwnerAndRepo() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	s.createGitHubIntegration(workspace.ID, user.ID)

	_, err := s.repoSvc.ConnectRepository(s.Ctx, &repository_service.ConnectRepositoryRequest{
		WorkspaceID:   workspace.ID,
		UserID:        user.ID,
		GitHubRepoID:  111,
		FullName:      "octocat/hello-world",
		DefaultBranch: "main",
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "octocat", s.mockWebhook.calledOwner)
	assert.Equal(s.T(), "hello-world", s.mockWebhook.calledRepo)
	assert.NotEmpty(s.T(), s.mockWebhook.calledSecret)
}

// TestConnectRepository_DuplicateReturnsError verifies that connecting the same repo twice
// returns ErrAlreadyConnected.
func (s *RepositoryServiceTestSuite) TestConnectRepository_DuplicateReturnsError() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	s.createGitHubIntegration(workspace.ID, user.ID)

	req := &repository_service.ConnectRepositoryRequest{
		WorkspaceID:   workspace.ID,
		UserID:        user.ID,
		GitHubRepoID:  789,
		FullName:      "octocat/dup-repo",
		DefaultBranch: "main",
	}

	_, err := s.repoSvc.ConnectRepository(s.Ctx, req)
	assert.NoError(s.T(), err)

	_, err = s.repoSvc.ConnectRepository(s.Ctx, req)
	assert.ErrorIs(s.T(), err, repository_service.ErrAlreadyConnected)
}

// TestConnectRepository_NoGitHubIntegrationReturnsError verifies that connecting a repo
// without a GitHub integration returns ErrGitHubNotConnected.
func (s *RepositoryServiceTestSuite) TestConnectRepository_NoGitHubIntegrationReturnsError() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()

	_, err := s.repoSvc.ConnectRepository(s.Ctx, &repository_service.ConnectRepositoryRequest{
		WorkspaceID:   workspace.ID,
		UserID:        user.ID,
		GitHubRepoID:  999,
		FullName:      "octocat/no-integration",
		DefaultBranch: "main",
	})
	assert.ErrorIs(s.T(), err, repository_service.ErrGitHubNotConnected)
}

// TestListRepositories_ReturnsOnlyWorkspaceRepos verifies workspace isolation.
func (s *RepositoryServiceTestSuite) TestListRepositories_ReturnsOnlyWorkspaceRepos() {
	ws1 := s.seeder.CreateWorkspace()
	ws2 := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	s.createGitHubIntegration(ws1.ID, user.ID)
	s.createGitHubIntegration(ws2.ID, user.ID)

	_, err := s.repoSvc.ConnectRepository(s.Ctx, &repository_service.ConnectRepositoryRequest{
		WorkspaceID: ws1.ID, UserID: user.ID, GitHubRepoID: 1, FullName: "a/repo1", DefaultBranch: "main",
	})
	assert.NoError(s.T(), err)

	_, err = s.repoSvc.ConnectRepository(s.Ctx, &repository_service.ConnectRepositoryRequest{
		WorkspaceID: ws2.ID, UserID: user.ID, GitHubRepoID: 2, FullName: "a/repo2", DefaultBranch: "main",
	})
	assert.NoError(s.T(), err)

	repos1, err := s.repoSvc.ListRepositories(s.Ctx, ws1.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), repos1, 1)
	assert.Equal(s.T(), int64(1), repos1[0].GitHubRepoID)

	repos2, err := s.repoSvc.ListRepositories(s.Ctx, ws2.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), repos2, 1)
	assert.Equal(s.T(), int64(2), repos2[0].GitHubRepoID)
}
