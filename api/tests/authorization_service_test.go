package tests

import (
	authorization_service "abstractions/services/authorization"
	"abstractions/tests/testutils"
	"context"
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
	s.authorizationService = authorization_service.NewService(s.DB)
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

// TestOwnerPermissions verifies that the Owner role has the full expected permission set.
func (s *AuthorizationTestSuite) TestOwnerPermissions() {
	user := s.seeder.CreateUser(testutils.WithEmail("owner-perms@example.com"))
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Owner Workspace"))
	s.seeder.CreateUserWorkspaceRelationship(user, workspace, authorization_service.RoleOwner.String())

	// Owner should hold all permissions.
	for _, perm := range []authorization_service.Permission{
		authorization_service.WorkspaceUpdate,
		authorization_service.WorkspaceDelete,
		authorization_service.WorkspaceManageUsers,
		authorization_service.WorkspaceManageSubscription,
		authorization_service.WorkspaceManageIntegrations,
	} {
		err := s.authorizationService.Validate(s.Ctx, user, workspace.ID, []authorization_service.Permission{perm})
		assert.NoError(s.T(), err, "owner should have permission %s", perm)
	}
}

// TestMemberPermissions verifies that the Member role has no elevated permissions and cannot
// perform workspace management actions.
func (s *AuthorizationTestSuite) TestMemberPermissions() {
	user := s.seeder.CreateUser(testutils.WithEmail("member-perms@example.com"))
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Member Workspace"))
	s.seeder.CreateUserWorkspaceRelationship(user, workspace, authorization_service.RoleMember.String())

	// Member should be denied all management permissions.
	for _, perm := range []authorization_service.Permission{
		authorization_service.WorkspaceUpdate,
		authorization_service.WorkspaceDelete,
		authorization_service.WorkspaceManageUsers,
		authorization_service.WorkspaceManageSubscription,
		authorization_service.WorkspaceManageIntegrations,
	} {
		err := s.authorizationService.Validate(s.Ctx, user, workspace.ID, []authorization_service.Permission{perm})
		assert.ErrorIs(s.T(), err, authorization_service.ErrUnauthorized, "member should NOT have permission %s", perm)
	}
}

// TestAdminCannotDeleteWorkspace verifies that the Admin role does not hold WorkspaceDelete
// (only Owners may delete a workspace).
func (s *AuthorizationTestSuite) TestAdminCannotDeleteWorkspace() {
	user := s.seeder.CreateUser(testutils.WithEmail("admin-nodelete@example.com"))
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Admin No-Delete Workspace"))
	s.seeder.CreateUserWorkspaceRelationship(user, workspace, authorization_service.RoleAdmin.String())

	err := s.authorizationService.Validate(s.Ctx, user, workspace.ID, []authorization_service.Permission{
		authorization_service.WorkspaceDelete,
	})
	assert.ErrorIs(s.T(), err, authorization_service.ErrUnauthorized)
}

// TestCanGrantRole_Hierarchy verifies the role-granting hierarchy:
//   - Members cannot grant any role
//   - Admins can grant Admin and Member but not Owner
//   - Owners can grant any role including Owner
func (s *AuthorizationTestSuite) TestCanGrantRole_Hierarchy() {
	// Members cannot grant any role
	assert.False(s.T(), authorization_service.CanGrantRole(authorization_service.RoleMember.String(), authorization_service.RoleMember.String()))
	assert.False(s.T(), authorization_service.CanGrantRole(authorization_service.RoleMember.String(), authorization_service.RoleAdmin.String()))
	assert.False(s.T(), authorization_service.CanGrantRole(authorization_service.RoleMember.String(), authorization_service.RoleOwner.String()))

	// Admins can grant Member and Admin but not Owner (Owner is a higher level)
	assert.True(s.T(), authorization_service.CanGrantRole(authorization_service.RoleAdmin.String(), authorization_service.RoleMember.String()))
	assert.True(s.T(), authorization_service.CanGrantRole(authorization_service.RoleAdmin.String(), authorization_service.RoleAdmin.String()))
	assert.False(s.T(), authorization_service.CanGrantRole(authorization_service.RoleAdmin.String(), authorization_service.RoleOwner.String()))

	// Owners can grant any role including Owner
	assert.True(s.T(), authorization_service.CanGrantRole(authorization_service.RoleOwner.String(), authorization_service.RoleMember.String()))
	assert.True(s.T(), authorization_service.CanGrantRole(authorization_service.RoleOwner.String(), authorization_service.RoleAdmin.String()))
	assert.True(s.T(), authorization_service.CanGrantRole(authorization_service.RoleOwner.String(), authorization_service.RoleOwner.String()))
}

// TestCanModifyUserRole_Hierarchy verifies that a user can only modify the role of a user with a
// strictly lower role level.
func (s *AuthorizationTestSuite) TestCanModifyUserRole_Hierarchy() {
	// Members cannot modify anyone
	assert.False(s.T(), authorization_service.CanModifyUserRole(authorization_service.RoleMember.String(), authorization_service.RoleMember.String()))
	assert.False(s.T(), authorization_service.CanModifyUserRole(authorization_service.RoleMember.String(), authorization_service.RoleAdmin.String()))
	assert.False(s.T(), authorization_service.CanModifyUserRole(authorization_service.RoleMember.String(), authorization_service.RoleOwner.String()))

	// Admins can modify Members but not other Admins or Owners
	assert.True(s.T(), authorization_service.CanModifyUserRole(authorization_service.RoleAdmin.String(), authorization_service.RoleMember.String()))
	assert.False(s.T(), authorization_service.CanModifyUserRole(authorization_service.RoleAdmin.String(), authorization_service.RoleAdmin.String()))
	assert.False(s.T(), authorization_service.CanModifyUserRole(authorization_service.RoleAdmin.String(), authorization_service.RoleOwner.String()))

	// Owners can modify Admins and Members but not other Owners
	assert.True(s.T(), authorization_service.CanModifyUserRole(authorization_service.RoleOwner.String(), authorization_service.RoleAdmin.String()))
	assert.True(s.T(), authorization_service.CanModifyUserRole(authorization_service.RoleOwner.String(), authorization_service.RoleMember.String()))
	assert.False(s.T(), authorization_service.CanModifyUserRole(authorization_service.RoleOwner.String(), authorization_service.RoleOwner.String()))
}

// TestGetRoleByID_And_GetRoleByName verifies the role lookup helpers.
func (s *AuthorizationTestSuite) TestGetRoleByID_And_GetRoleByName() {
	// GetRoleByID: known IDs return the correct role
	ownerRole := authorization_service.GetRoleByID(authorization_service.RoleOwner.String())
	assert.NotNil(s.T(), ownerRole)
	assert.Equal(s.T(), authorization_service.RoleOwner, ownerRole.ID)

	adminRole := authorization_service.GetRoleByID(authorization_service.RoleAdmin.String())
	assert.NotNil(s.T(), adminRole)
	assert.Equal(s.T(), authorization_service.RoleAdmin, adminRole.ID)

	memberRole := authorization_service.GetRoleByID(authorization_service.RoleMember.String())
	assert.NotNil(s.T(), memberRole)
	assert.Equal(s.T(), authorization_service.RoleMember, memberRole.ID)

	// GetRoleByID: unknown ID returns nil
	assert.Nil(s.T(), authorization_service.GetRoleByID("unknown_role"))

	// GetRoleByName: known names return the correct role
	byName := authorization_service.GetRoleByName("Owner")
	assert.NotNil(s.T(), byName)
	assert.Equal(s.T(), authorization_service.RoleOwner, byName.ID)

	// GetRoleByName: unknown name returns nil
	assert.Nil(s.T(), authorization_service.GetRoleByName("SuperAdmin"))
}

// TestRolePermissions_MapIsComplete verifies that RolePermissions() returns an entry for every
// role defined in ALL_ROLES.
func (s *AuthorizationTestSuite) TestRolePermissions_MapIsComplete() {
	all := s.authorizationService.GetAllRoles(context.Background())
	permMap := authorization_service.RolePermissions()

	for _, role := range all {
		_, ok := permMap[role.ID]
		assert.True(s.T(), ok, "RolePermissions() should contain an entry for role %s", role.ID)
	}
}

// TestAllRoleIDsInterface_ContainsAllRoles verifies that AllRoleIDsInterface() returns a slice
// whose length matches the number of defined roles and that each element is a string.
func (s *AuthorizationTestSuite) TestAllRoleIDsInterface_ContainsAllRoles() {
	all := s.authorizationService.GetAllRoles(context.Background())
	ids := authorization_service.AllRoleIDsInterface()

	assert.Equal(s.T(), len(all), len(ids), "AllRoleIDsInterface length should match ALL_ROLES length")

	for _, id := range ids {
		_, ok := id.(string)
		assert.True(s.T(), ok, "each element returned by AllRoleIDsInterface should be a string")
	}
}

// TestGetAllRoles_ReturnsCopy verifies that the service's GetAllRoles() returns all three
// expected roles with non-empty IDs and names.
func (s *AuthorizationTestSuite) TestGetAllRoles_ReturnsCopy() {
	roles := s.authorizationService.GetAllRoles(context.Background())
	assert.Len(s.T(), roles, 3, "there should be exactly 3 roles: owner, admin, member")

	for _, role := range roles {
		assert.NotEmpty(s.T(), role.ID, "role ID should not be empty")
		assert.NotEmpty(s.T(), role.Name, "role Name should not be empty")
	}
}
