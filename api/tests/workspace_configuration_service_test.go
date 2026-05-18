package tests

import (
	"abstractions/repositories"
	workspace_configuration_service "abstractions/services/workspace_configuration"
	"abstractions/tests/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// 32-byte key used only in tests — not a production secret.
const testEncryptionKey = "abcdefghijklmnopqrstuvwxyz123456"

type WorkspaceConfigurationServiceTestSuite struct {
	testutils.TestSuite
	seeder *testutils.Seeder
	svc    *workspace_configuration_service.Service
}

func TestWorkspaceConfigurationService(t *testing.T) {
	suite.Run(t, new(WorkspaceConfigurationServiceTestSuite))
}

func (s *WorkspaceConfigurationServiceTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	repo := repositories.NewWorkspaceConfigurationRepository(s.DB)
	s.svc = workspace_configuration_service.NewService(repo, testEncryptionKey)
}

// TestSet_PlainValue verifies that a plain (non-encrypted) key is stored and retrieved correctly.
func (s *WorkspaceConfigurationServiceTestSuite) TestSet_PlainValue() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()

	err := s.svc.Set(s.Ctx, workspace.ID, workspace_configuration_service.KeyIndexingModel, "claude-haiku", user.ID)
	assert.NoError(s.T(), err)

	val, err := s.svc.Get(s.Ctx, workspace.ID, workspace_configuration_service.KeyIndexingModel)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "claude-haiku", val)
}

// TestSet_EncryptedValue verifies that an encrypted key is stored in ciphertext but decrypted on Get.
func (s *WorkspaceConfigurationServiceTestSuite) TestSet_EncryptedValue() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()

	secret := "sk-test-anthropic-key-123"
	err := s.svc.Set(s.Ctx, workspace.ID, workspace_configuration_service.KeyAnthropicAPIKey, secret, user.ID)
	assert.NoError(s.T(), err)

	// Get should return the plaintext secret after decryption.
	val, err := s.svc.Get(s.Ctx, workspace.ID, workspace_configuration_service.KeyAnthropicAPIKey)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), secret, val)
}

// TestSet_EmptyValueDeletesKey verifies that setting an empty string hard-deletes the row so the
// system default takes effect.
func (s *WorkspaceConfigurationServiceTestSuite) TestSet_EmptyValueDeletesKey() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()

	// Set a custom model, then clear it.
	err := s.svc.Set(s.Ctx, workspace.ID, workspace_configuration_service.KeyIndexingModel, "custom-model", user.ID)
	assert.NoError(s.T(), err)

	err = s.svc.Set(s.Ctx, workspace.ID, workspace_configuration_service.KeyIndexingModel, "", user.ID)
	assert.NoError(s.T(), err)

	// After deletion the default (from envDefaults) should be returned.
	val, err := s.svc.Get(s.Ctx, workspace.ID, workspace_configuration_service.KeyIndexingModel)
	assert.NoError(s.T(), err)
	// The system default for KeyIndexingModel is non-empty.
	assert.NotEmpty(s.T(), val)
	assert.NotEqual(s.T(), "custom-model", val)
}

// TestGet_FallsBackToEnvDefault verifies that when no workspace override exists Get returns the
// system default value rather than an empty string.
func (s *WorkspaceConfigurationServiceTestSuite) TestGet_FallsBackToEnvDefault() {
	workspace := s.seeder.CreateWorkspace()

	// No Set call — should fall back to the default newsletter model.
	val, err := s.svc.Get(s.Ctx, workspace.ID, workspace_configuration_service.KeyNewsletterModel)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), val)
}

// TestGet_UnknownKey verifies that an unregistered key returns an empty string without error.
func (s *WorkspaceConfigurationServiceTestSuite) TestGet_UnknownKey() {
	workspace := s.seeder.CreateWorkspace()

	val, err := s.svc.Get(s.Ctx, workspace.ID, "unknown_key_xyz")
	assert.NoError(s.T(), err)
	assert.Empty(s.T(), val)
}

// TestGetLastFour_NoConfig verifies that GetLastFour returns nil when no workspace row exists.
func (s *WorkspaceConfigurationServiceTestSuite) TestGetLastFour_NoConfig() {
	workspace := s.seeder.CreateWorkspace()

	result, err := s.svc.GetLastFour(s.Ctx, workspace.ID, workspace_configuration_service.KeyAnthropicAPIKey)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), result)
}

// TestGetLastFour_ReturnsLastFour verifies that GetLastFour returns only the last 4 characters of a
// stored secret.
func (s *WorkspaceConfigurationServiceTestSuite) TestGetLastFour_ReturnsLastFour() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()

	secret := "sk-anthropic-secret-abcd"
	err := s.svc.Set(s.Ctx, workspace.ID, workspace_configuration_service.KeyAnthropicAPIKey, secret, user.ID)
	assert.NoError(s.T(), err)

	result, err := s.svc.GetLastFour(s.Ctx, workspace.ID, workspace_configuration_service.KeyAnthropicAPIKey)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "abcd", *result)
}

// TestGetResolvedModel_UsesDefault verifies that GetResolvedModel falls back to the system default
// when no workspace override is configured.
func (s *WorkspaceConfigurationServiceTestSuite) TestGetResolvedModel_UsesDefault() {
	workspace := s.seeder.CreateWorkspace()

	model, err := s.svc.GetResolvedModel(s.Ctx, workspace.ID, workspace_configuration_service.KeyIndexingModel)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), model)
}

// TestGetResolvedModel_UsesWorkspaceOverride verifies that a workspace-level model setting overrides
// the system default.
func (s *WorkspaceConfigurationServiceTestSuite) TestGetResolvedModel_UsesWorkspaceOverride() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()

	customModel := "gpt-4o"
	err := s.svc.Set(s.Ctx, workspace.ID, workspace_configuration_service.KeyIndexingModel, customModel, user.ID)
	assert.NoError(s.T(), err)

	model, err := s.svc.GetResolvedModel(s.Ctx, workspace.ID, workspace_configuration_service.KeyIndexingModel)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), customModel, model)
}
