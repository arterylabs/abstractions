package tests

import (
	"abstractions/db/models"
	integration_service "abstractions/services/integrations"
	github_integration "abstractions/services/integrations/core/github"
	"abstractions/tests/testutils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	testutils.TestSuite
	seeder             *testutils.Seeder
	integrationService *integration_service.Service
}

func TestIntegrationService(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.integrationService = integration_service.NewIntegrationService(s.DB)
}

// Test creating integrations
func (s *IntegrationTestSuite) TestCreateIntegration() {
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))

	// Test successful integration creation
	req := &integration_service.CreateIntegrationRequest{
		Type:        github_integration.IntegrationTypeGitHub,
		WorkspaceID: workspace.ID,
		Token:       "test-github-token",
		CreatedBy:   user.ID,
	}

	integration, err := s.integrationService.CreateIntegration(s.Ctx, req)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), integration)
	assert.Equal(s.T(), req.Type, integration.Type)
	assert.Equal(s.T(), req.WorkspaceID, integration.WorkspaceID)
	assert.Equal(s.T(), req.CreatedBy, integration.CreatedBy)
	assert.NotNil(s.T(), integration.ID)

	// Creating another integration of the same type is allowed; callers can scope
	// integrations by resource and list multiple integrations per type.
	duplicateIntegration, err := s.integrationService.CreateIntegration(s.Ctx, req)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), duplicateIntegration)
	assert.NotEqual(s.T(), integration.ID, duplicateIntegration.ID)

	integrations, err := s.integrationService.GetIntegrationsByType(s.Ctx, workspace.ID, req.Type)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), integrations, 2)
}

// Test getting integration by ID
func (s *IntegrationTestSuite) TestGetIntegrationByID() {
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))

	// Create integration
	req := &integration_service.CreateIntegrationRequest{
		Type:        github_integration.IntegrationTypeGitHub,
		WorkspaceID: workspace.ID,
		Token:       "test-token",
		CreatedBy:   user.ID,
	}

	integration, err := s.integrationService.CreateIntegration(s.Ctx, req)
	assert.NoError(s.T(), err)

	// Test successful retrieval
	found, err := s.integrationService.GetIntegration(s.Ctx, integration.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), found)
	assert.Equal(s.T(), integration.ID, found.ID)

	// Test non-existent integration
	notFound, err := s.integrationService.GetIntegration(s.Ctx, uuid.New())
	assert.Error(s.T(), err)
	assert.Nil(s.T(), notFound)
}

// Test getting integrations by workspace
func (s *IntegrationTestSuite) TestGetIntegrationsByWorkspace() {
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))

	// Create multiple integrations
	req := &integration_service.CreateIntegrationRequest{
		Type:        github_integration.IntegrationTypeGitHub,
		WorkspaceID: workspace.ID,
		Token:       "github-token",
		CreatedBy:   user.ID,
	}

	_, err := s.integrationService.CreateIntegration(s.Ctx, req)
	assert.NoError(s.T(), err)

	// Test getting all integrations for workspace
	integrations, err := s.integrationService.GetIntegrationsByWorkspace(s.Ctx, workspace.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), integrations, 1)
	assert.Equal(s.T(), github_integration.IntegrationTypeGitHub, integrations[0].Type)
}

// Test getting integrations by type
func (s *IntegrationTestSuite) TestGetIntegrationsByType() {
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))

	req := &integration_service.CreateIntegrationRequest{
		Type:        github_integration.IntegrationTypeGitHub,
		WorkspaceID: workspace.ID,
		Token:       "github-token",
		CreatedBy:   user.ID,
	}

	integration, err := s.integrationService.CreateIntegration(s.Ctx, req)
	assert.NoError(s.T(), err)

	// Test getting integrations by type
	integrations, err := s.integrationService.GetIntegrationsByType(s.Ctx, workspace.ID, github_integration.IntegrationTypeGitHub)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), integrations, 1)
	assert.Equal(s.T(), integration.ID, integrations[0].ID)

	// Test getting non-existent integration type
	nonExistentIntegrations, err := s.integrationService.GetIntegrationsByType(s.Ctx, workspace.ID, "non_existent_type")
	assert.NoError(s.T(), err)
	assert.Len(s.T(), nonExistentIntegrations, 0)
}

// Test getting single integration by type
func (s *IntegrationTestSuite) TestGetIntegrationByType() {
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))

	req := &integration_service.CreateIntegrationRequest{
		Type:        github_integration.IntegrationTypeGitHub,
		WorkspaceID: workspace.ID,
		Token:       "github-token",
		CreatedBy:   user.ID,
	}

	integration, err := s.integrationService.CreateIntegration(s.Ctx, req)
	assert.NoError(s.T(), err)

	// Test getting single integration by type
	foundIntegration, err := s.integrationService.GetIntegrationByType(s.Ctx, workspace.ID, github_integration.IntegrationTypeGitHub)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), foundIntegration)
	assert.Equal(s.T(), integration.ID, foundIntegration.ID)

	// Test getting non-existent integration type
	notFoundIntegration, err := s.integrationService.GetIntegrationByType(s.Ctx, workspace.ID, "non_existent_type")
	assert.Error(s.T(), err)
	assert.Nil(s.T(), notFoundIntegration)
}

// Test updating integration
func (s *IntegrationTestSuite) TestUpdateIntegration() {
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))

	// Create integration
	req := &integration_service.CreateIntegrationRequest{
		Type:        github_integration.IntegrationTypeGitHub,
		WorkspaceID: workspace.ID,
		Token:       "original-token",
		CreatedBy:   user.ID,
	}

	integration, err := s.integrationService.CreateIntegration(s.Ctx, req)
	assert.NoError(s.T(), err)

	// Update integration tokens using the service method
	err = s.integrationService.UpdateTokens(s.Ctx, integration.ID, "updated-token", nil, nil, user.ID)
	assert.NoError(s.T(), err)

	// Verify update by checking the decrypted token
	decryptedToken, err := s.integrationService.GetDecryptedToken(s.Ctx, integration.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "updated-token", decryptedToken)
}

// Test updating integration configuration
func (s *IntegrationTestSuite) TestUpdateIntegrationConfiguration() {
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))

	// Create integration
	req := &integration_service.CreateIntegrationRequest{
		Type:        github_integration.IntegrationTypeGitHub,
		WorkspaceID: workspace.ID,
		Token:       "github-token",
		CreatedBy:   user.ID,
	}

	integration, err := s.integrationService.CreateIntegration(s.Ctx, req)
	assert.NoError(s.T(), err)

	// Update configuration
	newConfig := models.JSONBField{
		"away_message": "Updated away message",
		"auto_reply":   true,
	}

	err = s.integrationService.UpdateIntegrationConfiguration(s.Ctx, integration.ID, newConfig, user.ID)
	assert.NoError(s.T(), err)

	// Verify configuration update
	updated, err := s.integrationService.GetIntegration(s.Ctx, integration.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), updated.Configuration)
	assert.Equal(s.T(), "Updated away message", updated.Configuration["away_message"])
	assert.Equal(s.T(), true, updated.Configuration["auto_reply"])
}

// Test deleting integration
func (s *IntegrationTestSuite) TestDeleteIntegration() {
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))

	// Create integration
	req := &integration_service.CreateIntegrationRequest{
		Type:        github_integration.IntegrationTypeGitHub,
		WorkspaceID: workspace.ID,
		Token:       "github-token",
		CreatedBy:   user.ID,
	}

	integration, err := s.integrationService.CreateIntegration(s.Ctx, req)
	assert.NoError(s.T(), err)

	// Delete integration
	err = s.integrationService.DeleteIntegration(s.Ctx, integration.ID, user.ID)
	assert.NoError(s.T(), err)

	// Verify deletion (should be soft deleted)
	deleted, err := s.integrationService.GetIntegration(s.Ctx, integration.ID)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), deleted)
}

// Test counting integrations
func (s *IntegrationTestSuite) TestCountIntegrations() {
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))

	// Initially should have 0 integrations
	count, err := s.integrationService.CountIntegrations(s.Ctx, workspace.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(0), count)

	// Create integration
	req := &integration_service.CreateIntegrationRequest{
		Type:        github_integration.IntegrationTypeGitHub,
		WorkspaceID: workspace.ID,
		Token:       "github-token",
		CreatedBy:   user.ID,
	}

	_, err = s.integrationService.CreateIntegration(s.Ctx, req)
	assert.NoError(s.T(), err)

	// Should now have 1 integration
	count, err = s.integrationService.CountIntegrations(s.Ctx, workspace.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), count)
}

// Test error handling
func (s *IntegrationTestSuite) TestErrorHandling() {
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))

	// Test creating integration with invalid workspace
	req := &integration_service.CreateIntegrationRequest{
		Type:        github_integration.IntegrationTypeGitHub,
		WorkspaceID: uuid.New(), // Non-existent workspace
		Token:       "github-token",
		CreatedBy:   user.ID,
	}

	integration, err := s.integrationService.CreateIntegration(s.Ctx, req)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), integration)

	// Test updating non-existent integration
	nonExistentIntegration := &models.Integration{
		ID:          uuid.New(),
		Type:        github_integration.IntegrationTypeGitHub,
		WorkspaceID: workspace.ID,
		Token:       "token",
		CreatedBy:   user.ID,
		UpdatedBy:   user.ID,
	}

	err = s.integrationService.UpdateIntegration(s.Ctx, nonExistentIntegration)
	assert.Error(s.T(), err)

	// Test deleting non-existent integration
	err = s.integrationService.DeleteIntegration(s.Ctx, uuid.New(), user.ID)
	assert.Error(s.T(), err)
}

// TestGitHubIntegrationEncryptsToken verifies that CreateIntegration with the GitHub type
// encrypts the access token at rest and that GetDecryptedToken recovers the original value.
func (s *IntegrationTestSuite) TestGitHubIntegrationEncryptsToken() {
	workspace := s.seeder.CreateWorkspace(testutils.WithName("GitHub Workspace"))
	user := s.seeder.CreateUser(testutils.WithEmail("github-test@example.com"))

	rawToken := "ghp_test_access_token_abc123"

	req := &integration_service.CreateIntegrationRequest{
		Type:        github_integration.IntegrationTypeGitHub,
		WorkspaceID: workspace.ID,
		Token:       rawToken,
		CreatedBy:   user.ID,
	}

	integration, err := s.integrationService.CreateIntegration(s.Ctx, req)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), integration)
	assert.Equal(s.T(), github_integration.IntegrationTypeGitHub, integration.Type)

	// The stored token must be encrypted (not equal to the raw value)
	assert.NotEqual(s.T(), rawToken, integration.Token, "token should be stored encrypted")

	// GetDecryptedToken must return the original plaintext value
	decrypted, err := s.integrationService.GetDecryptedToken(s.Ctx, integration.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), rawToken, decrypted, "decrypted token should match original")
}

// TestGitHubIntegrationGetRequiredScopes verifies the GitHub integration type returns the expected scopes.
func (s *IntegrationTestSuite) TestGitHubIntegrationGetRequiredScopes() {
	ghIntegration := github_integration.New()
	scopes := ghIntegration.GetRequiredScopes()
	assert.Contains(s.T(), scopes, "repo", "GitHub integration must request 'repo' scope")
	assert.Equal(s.T(), github_integration.IntegrationTypeGitHub, ghIntegration.GetIntegrationType())
}
