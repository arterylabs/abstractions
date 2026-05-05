package tests

import (
	"abstractions/db/models"
	workspace_service "abstractions/services/workspace"
	"abstractions/tests/testutils"
	"abstractions/utils/pagination_utils"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WorkspaceServiceTestSuite struct {
	testutils.TestSuite
	seeder           *testutils.Seeder
	workspaceService *workspace_service.Service
}

func TestWorkspaceService(t *testing.T) {
	suite.Run(t, new(WorkspaceServiceTestSuite))
}

func (s *WorkspaceServiceTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.workspaceService = workspace_service.NewWorkspaceService(s.DB)
}

// Test workspace creation
func (s *WorkspaceServiceTestSuite) TestCreateWorkspace() {
	user := s.seeder.CreateUser(testutils.WithEmail("user@example.com"))
	workspaceName := "New Test Workspace"

	// Test successful workspace creation
	workspace, err := s.workspaceService.CreateWorkspace(s.Ctx, workspaceName, user.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), workspace)
	assert.Equal(s.T(), workspaceName, workspace.Name)

	// Verify user relationship was created
	relationships, err := s.workspaceService.GetWorkspaceUsers(s.Ctx, workspace.ID, user.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), relationships, 1)
}

// Test workspace retrieval
func (s *WorkspaceServiceTestSuite) TestGetWorkspace() {
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))

	// Test successful retrieval
	found, err := s.workspaceService.GetWorkspace(s.Ctx, workspace.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), found)
	assert.Equal(s.T(), workspace.ID, found.ID)

	// Test non-existent workspace
	found, err = s.workspaceService.GetWorkspace(s.Ctx, uuid.New())
	assert.Error(s.T(), err)
	assert.Equal(s.T(), workspace_service.ErrWorkspaceNotFound, err)
	assert.Nil(s.T(), found)
}

// Test workspace update
func (s *WorkspaceServiceTestSuite) TestUpdateWorkspace() {
	user := s.seeder.CreateUser(testutils.WithEmail("user@example.com"))
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Old Name"))
	s.seeder.CreateUserWorkspaceRelationship(user, workspace, "member")

	params := workspace_service.UpdateWorkspaceParams{
		Name: "New Name",
	}
	err := s.workspaceService.UpdateWorkspace(s.Ctx, workspace.ID, params, user.ID)
	assert.NoError(s.T(), err)

	// Test update with non-existent workspace
	err = s.workspaceService.UpdateWorkspace(s.Ctx, uuid.New(), params, user.ID)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), workspace_service.ErrWorkspaceNotFound, err)
}

// Test workspace delete
func (s *WorkspaceServiceTestSuite) TestDeleteWorkspace() {
	user := s.seeder.CreateUser(testutils.WithEmail("user@example.com"))
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))
	s.seeder.CreateUserWorkspaceRelationship(user, workspace, "member")

	err := s.workspaceService.DeleteWorkspace(s.Ctx, workspace.ID, user.ID)
	assert.NoError(s.T(), err)

	// Verify workspace is deleted
	found, err := s.workspaceService.GetWorkspace(s.Ctx, workspace.ID)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), found) // Should not be found since it's deleted
}

// Test adding user to workspace
func (s *WorkspaceServiceTestSuite) TestAddUserToWorkspace() {
	user1 := s.seeder.CreateUser(testutils.WithEmail("user1@example.com"))
	user2 := s.seeder.CreateUser(testutils.WithEmail("user2@example.com"))
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))

	// Setup initial relationship - user1 as admin so they can add other users
	s.seeder.CreateUserWorkspaceRelationship(user1, workspace, "admin")

	// Test adding another user
	err := s.workspaceService.AddUserToWorkspace(s.Ctx, workspace.ID, user2.ID, "member", user1.ID)
	assert.NoError(s.T(), err)

	// Verify users are in workspace
	relationships, err := s.workspaceService.GetWorkspaceUsers(s.Ctx, workspace.ID, user1.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), relationships, 2)
}

// Test removing user from workspace
func (s *WorkspaceServiceTestSuite) TestRemoveUserFromWorkspace() {
	user1 := s.seeder.CreateUser(testutils.WithEmail("user1@example.com"))
	user2 := s.seeder.CreateUser(testutils.WithEmail("user2@example.com"))
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))

	// Setup relationships - user1 as admin so they can remove other users
	s.seeder.CreateUserWorkspaceRelationship(user1, workspace, "admin")
	s.seeder.CreateUserWorkspaceRelationship(user2, workspace, "member")

	// Test removing a user
	err := s.workspaceService.RemoveUserFromWorkspace(s.Ctx, workspace.ID, user2.ID, user1.ID)
	assert.NoError(s.T(), err)

	// Verify user was removed
	relationships, err := s.workspaceService.GetWorkspaceUsers(s.Ctx, workspace.ID, user1.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), relationships, 1)
}

// Test workspace user listing
func (s *WorkspaceServiceTestSuite) TestGetWorkspaceUsers() {
	user1 := s.seeder.CreateUser(testutils.WithEmail("user1@example.com"))
	user2 := s.seeder.CreateUser(testutils.WithEmail("user2@example.com"))
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))

	// Setup relationships
	s.seeder.CreateUserWorkspaceRelationship(user1, workspace, "member")
	s.seeder.CreateUserWorkspaceRelationship(user2, workspace, "member")

	// Test listing workspace users
	users, err := s.workspaceService.GetWorkspaceUsers(s.Ctx, workspace.ID, user1.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), users, 2)
}

// Test getting user's workspaces
func (s *WorkspaceServiceTestSuite) TestGetUserWorkspaces() {
	user := s.seeder.CreateUser(testutils.WithEmail("user@example.com"))
	workspace1 := s.seeder.CreateWorkspace(testutils.WithName("Workspace 1"))
	workspace2 := s.seeder.CreateWorkspace(testutils.WithName("Workspace 2"))

	// Add user to workspaces
	s.seeder.CreateUserWorkspaceRelationship(user, workspace1, "member")
	s.seeder.CreateUserWorkspaceRelationship(user, workspace2, "member")

	// Test getting user's workspaces
	workspaces, err := s.workspaceService.GetUserWorkspaces(s.Ctx, user.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), workspaces, 2)
}

// Test workspace user listing with pagination
func (s *WorkspaceServiceTestSuite) TestGetWorkspaceUsersPaginated() {
	// Create test users and workspace
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))
	var users []*models.User

	// Create 15 users for testing pagination
	for i := 0; i < 15; i++ {
		user := s.seeder.CreateUser(testutils.WithEmail(fmt.Sprintf("user%d@example.com", i)))
		users = append(users, user)
		s.seeder.CreateUserWorkspaceRelationship(user, workspace, "member")
	}

	// Test first page with default page size
	params := pagination_utils.PaginationParams{
		Page:    1,
		PerPage: 10,
	}
	result, total, err := s.workspaceService.GetWorkspaceUsersPaginated(s.Ctx, workspace.ID, params)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 15, total)
	assert.Len(s.T(), result, 10)

	// Test second page
	params.Page = 2
	result, total, err = s.workspaceService.GetWorkspaceUsersPaginated(s.Ctx, workspace.ID, params)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 15, total)
	assert.Len(s.T(), result, 5)

	// Test with smaller page size
	params = pagination_utils.PaginationParams{
		Page:    1,
		PerPage: 5,
	}
	result, total, err = s.workspaceService.GetWorkspaceUsersPaginated(s.Ctx, workspace.ID, params)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 15, total)
	assert.Len(s.T(), result, 5)

	// Test page beyond available results
	params.Page = 4
	result, total, err = s.workspaceService.GetWorkspaceUsersPaginated(s.Ctx, workspace.ID, params)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 15, total)
	assert.Len(s.T(), result, 0)

	// Test with invalid page number
	params = pagination_utils.PaginationParams{
		Page:    0,
		PerPage: 10,
	}
	result, total, err = s.workspaceService.GetWorkspaceUsersPaginated(s.Ctx, workspace.ID, params)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 15, total)
	assert.Equal(s.T(), len(result), 10) // Should return first page by default

	// Test with invalid page size
	params = pagination_utils.PaginationParams{
		Page:    1,
		PerPage: 0,
	}
	result, total, err = s.workspaceService.GetWorkspaceUsersPaginated(s.Ctx, workspace.ID, params)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 15, total)
	assert.Equal(s.T(), 10, len(result)) // If per page is 0, it should default to 10 automatically
}

// Role hierarchy validation tests

// Test role hierarchy: member cannot modify other users
func (s *WorkspaceServiceTestSuite) TestMemberCannotModifyUsers() {
	memberUser := s.seeder.CreateUser(testutils.WithEmail("member@example.com"))
	targetUser := s.seeder.CreateUser(testutils.WithEmail("target@example.com"))
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))

	// Setup relationships
	s.seeder.CreateUserWorkspaceRelationship(memberUser, workspace, "member")
	s.seeder.CreateUserWorkspaceRelationship(targetUser, workspace, "member")

	// Member should not be able to update another member's role
	err := s.workspaceService.UpdateUserRole(s.Ctx, workspace.ID, targetUser.ID, "admin", memberUser.ID)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), workspace_service.ErrInsufficientPermissions, err)

	// Member should not be able to remove another member
	err = s.workspaceService.RemoveUserFromWorkspace(s.Ctx, workspace.ID, targetUser.ID, memberUser.ID)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), workspace_service.ErrInsufficientPermissions, err)

	// Member should not be able to add new users
	newUser := s.seeder.CreateUser(testutils.WithEmail("newuser@example.com"))
	err = s.workspaceService.AddUserToWorkspace(s.Ctx, workspace.ID, newUser.ID, "member", memberUser.ID)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), workspace_service.ErrInsufficientPermissions, err)
}

// Test role hierarchy: admin can modify members but not other admins or owners
func (s *WorkspaceServiceTestSuite) TestAdminRoleHierarchy() {
	adminUser := s.seeder.CreateUser(testutils.WithEmail("admin@example.com"))
	memberUser := s.seeder.CreateUser(testutils.WithEmail("member@example.com"))
	otherAdminUser := s.seeder.CreateUser(testutils.WithEmail("otheradmin@example.com"))
	ownerUser := s.seeder.CreateUser(testutils.WithEmail("owner@example.com"))
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))

	// Setup relationships
	s.seeder.CreateUserWorkspaceRelationship(adminUser, workspace, "admin")
	s.seeder.CreateUserWorkspaceRelationship(memberUser, workspace, "member")
	s.seeder.CreateUserWorkspaceRelationship(otherAdminUser, workspace, "admin")
	s.seeder.CreateUserWorkspaceRelationship(ownerUser, workspace, "owner")

	// Admin can update member's role
	err := s.workspaceService.UpdateUserRole(s.Ctx, workspace.ID, memberUser.ID, "admin", adminUser.ID)
	assert.NoError(s.T(), err)

	// Admin can remove members
	newMember := s.seeder.CreateUser(testutils.WithEmail("newmember@example.com"))
	s.seeder.CreateUserWorkspaceRelationship(newMember, workspace, "member")
	err = s.workspaceService.RemoveUserFromWorkspace(s.Ctx, workspace.ID, newMember.ID, adminUser.ID)
	assert.NoError(s.T(), err)

	// Admin cannot update another admin's role
	err = s.workspaceService.UpdateUserRole(s.Ctx, workspace.ID, otherAdminUser.ID, "member", adminUser.ID)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), workspace_service.ErrInsufficientPermissions, err)

	// Admin cannot update owner's role
	err = s.workspaceService.UpdateUserRole(s.Ctx, workspace.ID, ownerUser.ID, "admin", adminUser.ID)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), workspace_service.ErrInsufficientPermissions, err)

	// Admin cannot grant owner role
	err = s.workspaceService.UpdateUserRole(s.Ctx, workspace.ID, memberUser.ID, "owner", adminUser.ID)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), workspace_service.ErrInsufficientPermissions, err)
}

// Test role hierarchy: owner can modify all roles
func (s *WorkspaceServiceTestSuite) TestOwnerCanModifyAllRoles() {
	ownerUser := s.seeder.CreateUser(testutils.WithEmail("owner@example.com"))
	adminUser := s.seeder.CreateUser(testutils.WithEmail("admin@example.com"))
	memberUser := s.seeder.CreateUser(testutils.WithEmail("member@example.com"))
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))

	// Setup relationships
	s.seeder.CreateUserWorkspaceRelationship(ownerUser, workspace, "owner")
	s.seeder.CreateUserWorkspaceRelationship(adminUser, workspace, "admin")
	s.seeder.CreateUserWorkspaceRelationship(memberUser, workspace, "member")

	// Owner can update admin's role
	err := s.workspaceService.UpdateUserRole(s.Ctx, workspace.ID, adminUser.ID, "member", ownerUser.ID)
	assert.NoError(s.T(), err)

	// Owner can update member's role to admin
	err = s.workspaceService.UpdateUserRole(s.Ctx, workspace.ID, memberUser.ID, "admin", ownerUser.ID)
	assert.NoError(s.T(), err)

	// Owner can remove any user
	err = s.workspaceService.RemoveUserFromWorkspace(s.Ctx, workspace.ID, adminUser.ID, ownerUser.ID)
	assert.NoError(s.T(), err)

	// Owner can grant any role
	newUser := s.seeder.CreateUser(testutils.WithEmail("newuser@example.com"))
	err = s.workspaceService.AddUserToWorkspace(s.Ctx, workspace.ID, newUser.ID, "owner", ownerUser.ID)
	assert.NoError(s.T(), err)
}

// Test self-modification prevention
func (s *WorkspaceServiceTestSuite) TestCannotModifyOwnRole() {
	user := s.seeder.CreateUser(testutils.WithEmail("user@example.com"))
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))

	// Setup relationship
	s.seeder.CreateUserWorkspaceRelationship(user, workspace, "owner")

	// User cannot update their own role
	err := s.workspaceService.UpdateUserRole(s.Ctx, workspace.ID, user.ID, "admin", user.ID)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), workspace_service.ErrCannotModifyOwnRole, err)

	// User cannot remove themselves from workspace
	err = s.workspaceService.RemoveUserFromWorkspace(s.Ctx, workspace.ID, user.ID, user.ID)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), workspace_service.ErrCannotModifyOwnRole, err)
}

// Test god mode user can modify any role
func (s *WorkspaceServiceTestSuite) TestGodModeCanModifyAnyRole() {
	// Create god mode user and target users
	godModeUser := s.seeder.CreateUser(testutils.WithEmail("god@example.com"))
	godModeUser.GodMode = true
	// Update the user in the database to set god mode
	if err := s.DB.Save(godModeUser).Error; err != nil {
		panic(err)
	}

	ownerUser := s.seeder.CreateUser(testutils.WithEmail("owner@example.com"))
	adminUser := s.seeder.CreateUser(testutils.WithEmail("admin@example.com"))
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))

	// Setup relationships - god mode user doesn't need to be in the workspace
	s.seeder.CreateUserWorkspaceRelationship(ownerUser, workspace, "owner")
	s.seeder.CreateUserWorkspaceRelationship(adminUser, workspace, "admin")

	// God mode user can change owner to admin
	err := s.workspaceService.UpdateUserRole(s.Ctx, workspace.ID, ownerUser.ID, "admin", godModeUser.ID)
	assert.NoError(s.T(), err)

	// God mode user can change admin to owner
	err = s.workspaceService.UpdateUserRole(s.Ctx, workspace.ID, adminUser.ID, "owner", godModeUser.ID)
	assert.NoError(s.T(), err)

	// God mode user can add any user with any role
	newUser := s.seeder.CreateUser(testutils.WithEmail("newuser@example.com"))
	err = s.workspaceService.AddUserToWorkspace(s.Ctx, workspace.ID, newUser.ID, "owner", godModeUser.ID)
	assert.NoError(s.T(), err)

	// God mode user can remove any user
	err = s.workspaceService.RemoveUserFromWorkspace(s.Ctx, workspace.ID, ownerUser.ID, godModeUser.ID)
	assert.NoError(s.T(), err)
}
