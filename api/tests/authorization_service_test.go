package tests

import (
	authorization_service "abstractions/services/authorization"
	"abstractions/tests/testutils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthorizationTestSuite struct {
	testutils.TestSuite
	seeder               *testutils.Seeder
	authorizationService *authorization_service.Service
}

func TestAuthorizationService(t *testing.T) {
	suite.Run(t, new(AuthorizationTestSuite))
}

func (s *AuthorizationTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.authorizationService = authorization_service.NewAuthorizationService(s.DB)
}

// Test god mode bypass
func (s *AuthorizationTestSuite) TestGodModeBypass() {
	user := s.seeder.CreateUser(testutils.WithEmail("god@example.com"))
	user.GodMode = true // Set god mode directly
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))

	// Even without any role, god mode should bypass all permission checks
	err := s.authorizationService.Validate(s.Ctx, user, workspace.ID, []authorization_service.Permission{
		authorization_service.WorkspaceUpdate,
		authorization_service.WorkspaceManageUsers,
		authorization_service.WorkspaceManageSubscription,
	})
	assert.NoError(s.T(), err)

	// Test HasPermission with god mode
	hasPermission := s.authorizationService.HasPermission(s.Ctx, user, workspace.ID, authorization_service.WorkspaceManageUsers.String())
	assert.True(s.T(), hasPermission)
}

// Test role-based access control
func (s *AuthorizationTestSuite) TestRoleBasedAccess() {
	user := s.seeder.CreateUser(testutils.WithEmail("member@example.com"))
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))
	s.seeder.CreateUserWorkspaceRelationship(user, workspace, authorization_service.RoleMember.String())

	// Test member role permissions
	err := s.authorizationService.Validate(s.Ctx, user, workspace.ID, []authorization_service.Permission{})
	assert.NoError(s.T(), err)

	// Test member role should not have admin permissions
	err = s.authorizationService.Validate(s.Ctx, user, workspace.ID, []authorization_service.Permission{
		authorization_service.WorkspaceUpdate,
	})
	assert.Error(s.T(), err)
	assert.Equal(s.T(), authorization_service.ErrUnauthorized, err)
}

// Test admin permissions
func (s *AuthorizationTestSuite) TestAdminPermissions() {
	user := s.seeder.CreateUser(testutils.WithEmail("admin@example.com"))
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))
	s.seeder.CreateUserWorkspaceRelationship(user, workspace, authorization_service.RoleAdmin.String())

	// Test admin role permissions
	err := s.authorizationService.Validate(s.Ctx, user, workspace.ID, []authorization_service.Permission{
		authorization_service.WorkspaceUpdate,
		authorization_service.WorkspaceManageUsers,
		authorization_service.WorkspaceManageSubscription,
	})
	assert.NoError(s.T(), err)

	// Test HasPermission for admin
	hasPermission := s.authorizationService.HasPermission(s.Ctx, user, workspace.ID, authorization_service.WorkspaceManageUsers.String())
	assert.True(s.T(), hasPermission)
}

// Test no role in workspace
func (s *AuthorizationTestSuite) TestNoRole() {
	user := s.seeder.CreateUser(testutils.WithEmail("norole@example.com"))
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))

	// User has no role in workspace
	err := s.authorizationService.Validate(s.Ctx, user, workspace.ID, []authorization_service.Permission{})
	assert.Error(s.T(), err)
	assert.Equal(s.T(), authorization_service.ErrNoRole, err)

	// Test HasPermission with no role
	hasPermission := s.authorizationService.HasPermission(s.Ctx, user, workspace.ID, authorization_service.WorkspaceDelete.String())
	assert.False(s.T(), hasPermission)
}

// Test invalid role
func (s *AuthorizationTestSuite) TestInvalidRole() {
	user := s.seeder.CreateUser(testutils.WithEmail("invalid@example.com"))
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))
	s.seeder.CreateUserWorkspaceRelationship(user, workspace, "invalid_role")

	// User has invalid role
	err := s.authorizationService.Validate(s.Ctx, user, workspace.ID, []authorization_service.Permission{})
	assert.Error(s.T(), err)
	assert.Equal(s.T(), authorization_service.ErrUnauthorized, err)
}

// Test non-existent workspace
func (s *AuthorizationTestSuite) TestNonExistentWorkspace() {
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))
	nonExistentWorkspaceID := uuid.New()

	// Test with non-existent workspace
	err := s.authorizationService.Validate(s.Ctx, user, nonExistentWorkspaceID, []authorization_service.Permission{})
	assert.Error(s.T(), err)
	assert.Equal(s.T(), authorization_service.ErrNoRole, err)
}

// Test multiple permission checks
func (s *AuthorizationTestSuite) TestMultiplePermissions() {
	user := s.seeder.CreateUser(testutils.WithEmail("member@example.com"))
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))
	s.seeder.CreateUserWorkspaceRelationship(user, workspace, authorization_service.RoleOwner.String())

	// Test multiple valid permissions
	err := s.authorizationService.Validate(s.Ctx, user, workspace.ID, []authorization_service.Permission{
		authorization_service.WorkspaceUpdate,
		authorization_service.WorkspaceManageUsers,
	})
	assert.NoError(s.T(), err)

	// Test mixture of valid and invalid permissions
	err = s.authorizationService.Validate(s.Ctx, user, workspace.ID, []authorization_service.Permission{
		authorization_service.WorkspaceDelete,
		"workspaceManageUsers",
	})
	assert.Error(s.T(), err)
	assert.Equal(s.T(), authorization_service.ErrUnauthorized, err)
}

// Test empty permissions list
func (s *AuthorizationTestSuite) TestEmptyPermissions() {
	user := s.seeder.CreateUser(testutils.WithEmail("member@example.com"))
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))
	s.seeder.CreateUserWorkspaceRelationship(user, workspace, authorization_service.RoleMember.String())

	// Test with empty permissions list
	err := s.authorizationService.Validate(s.Ctx, user, workspace.ID, []authorization_service.Permission{})
	assert.NoError(s.T(), err)
}
