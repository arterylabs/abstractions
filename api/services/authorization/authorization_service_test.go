package authorization_service_test

import (
	"abstractions/db/models"
	authorization_service "abstractions/services/authorization"
	workspace_service "abstractions/services/workspace"
	"abstractions/tests/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthorizationServiceTestSuite struct {
	testutils.TestSuite
	seeder       *testutils.Seeder
	service      *authorization_service.Service
	workspaceSvc *workspace_service.Service
}

func TestAuthorizationService(t *testing.T) {
	suite.Run(t, new(AuthorizationServiceTestSuite))
}

func (s *AuthorizationServiceTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.service = authorization_service.NewAuthorizationService(s.DB)
	s.workspaceSvc = workspace_service.NewWorkspaceService(s.DB)
}

// ---- permission matrix (pure logic, no DB needed) ----

func (s *AuthorizationServiceTestSuite) TestOwnerRole_HasAllPermissions() {
	role := authorization_service.GetRoleByID(authorization_service.RoleOwner.String())
	s.Require().NotNil(role)

	permSet := make(map[authorization_service.Permission]bool)
	for _, p := range role.Permissions {
		permSet[p] = true
	}

	assert.True(s.T(), permSet[authorization_service.WorkspaceUpdate])
	assert.True(s.T(), permSet[authorization_service.WorkspaceDelete])
	assert.True(s.T(), permSet[authorization_service.WorkspaceManageUsers])
	assert.True(s.T(), permSet[authorization_service.WorkspaceManageIntegrations])
	assert.True(s.T(), permSet[authorization_service.WorkspaceManageRepositories])
}

func (s *AuthorizationServiceTestSuite) TestMemberRole_HasNoPermissions() {
	role := authorization_service.GetRoleByID(authorization_service.RoleMember.String())
	s.Require().NotNil(role)

	assert.Empty(s.T(), role.Permissions, "member role should have no permissions")
}

func (s *AuthorizationServiceTestSuite) TestAdminRole_CannotDeleteWorkspace() {
	role := authorization_service.GetRoleByID(authorization_service.RoleAdmin.String())
	s.Require().NotNil(role)

	for _, p := range role.Permissions {
		assert.NotEqual(s.T(), authorization_service.WorkspaceDelete, p, "admin should not have WorkspaceDelete")
	}
}

func (s *AuthorizationServiceTestSuite) TestCanGrantRole_OwnerCanGrantAll() {
	assert.True(s.T(), authorization_service.CanGrantRole(
		authorization_service.RoleOwner.String(),
		authorization_service.RoleOwner.String(),
	))
	assert.True(s.T(), authorization_service.CanGrantRole(
		authorization_service.RoleOwner.String(),
		authorization_service.RoleAdmin.String(),
	))
	assert.True(s.T(), authorization_service.CanGrantRole(
		authorization_service.RoleOwner.String(),
		authorization_service.RoleMember.String(),
	))
}

func (s *AuthorizationServiceTestSuite) TestCanGrantRole_MemberCannotGrantAny() {
	assert.False(s.T(), authorization_service.CanGrantRole(
		authorization_service.RoleMember.String(),
		authorization_service.RoleMember.String(),
	))
	assert.False(s.T(), authorization_service.CanGrantRole(
		authorization_service.RoleMember.String(),
		authorization_service.RoleAdmin.String(),
	))
}

func (s *AuthorizationServiceTestSuite) TestCanModifyUserRole_OwnerCanModifyAdmin() {
	assert.True(s.T(), authorization_service.CanModifyUserRole(
		authorization_service.RoleOwner.String(),
		authorization_service.RoleAdmin.String(),
	))
}

func (s *AuthorizationServiceTestSuite) TestCanModifyUserRole_AdminCannotModifyOwner() {
	assert.False(s.T(), authorization_service.CanModifyUserRole(
		authorization_service.RoleAdmin.String(),
		authorization_service.RoleOwner.String(),
	))
}

func (s *AuthorizationServiceTestSuite) TestCanModifyUserRole_MemberCannotModifyAnyone() {
	assert.False(s.T(), authorization_service.CanModifyUserRole(
		authorization_service.RoleMember.String(),
		authorization_service.RoleMember.String(),
	))
}

// ---- DB-backed Validate tests ----

func (s *AuthorizationServiceTestSuite) TestValidate_OwnerPassesAllChecks() {
	owner := s.seeder.CreateUser()
	ws, err := s.workspaceSvc.CreateWorkspace(s.Ctx, "Auth WS Owner", owner.ID)
	s.Require().NoError(err)

	err = s.service.Validate(s.Ctx, owner, ws.ID, []authorization_service.Permission{
		authorization_service.WorkspaceUpdate,
		authorization_service.WorkspaceDelete,
		authorization_service.WorkspaceManageUsers,
	})

	assert.NoError(s.T(), err)
}

func (s *AuthorizationServiceTestSuite) TestValidate_MemberIsBlockedFromAdminActions() {
	owner := s.seeder.CreateUser()
	ws, err := s.workspaceSvc.CreateWorkspace(s.Ctx, "Auth WS Member", owner.ID)
	s.Require().NoError(err)

	member := s.seeder.CreateUser()
	s.Require().NoError(
		s.workspaceSvc.AddUserToWorkspace(s.Ctx, ws.ID, member.ID, authorization_service.RoleMember.String(), owner.ID),
	)

	err = s.service.Validate(s.Ctx, member, ws.ID, []authorization_service.Permission{
		authorization_service.WorkspaceDelete,
	})

	assert.ErrorIs(s.T(), err, authorization_service.ErrUnauthorized)
}

func (s *AuthorizationServiceTestSuite) TestValidate_GodModeBypassesAll() {
	owner := s.seeder.CreateUser()
	ws, err := s.workspaceSvc.CreateWorkspace(s.Ctx, "Auth WS GodMode", owner.ID)
	s.Require().NoError(err)

	godUser := s.seeder.CreateUser(func(u *models.User) {
		u.GodMode = true
	})

	err = s.service.Validate(s.Ctx, godUser, ws.ID, []authorization_service.Permission{
		authorization_service.WorkspaceDelete,
		authorization_service.WorkspaceManageUsers,
	})

	assert.NoError(s.T(), err)
}

func (s *AuthorizationServiceTestSuite) TestValidate_UserNotInWorkspace_ReturnsNoRole() {
	owner := s.seeder.CreateUser()
	ws, err := s.workspaceSvc.CreateWorkspace(s.Ctx, "Auth WS Outsider", owner.ID)
	s.Require().NoError(err)

	outsider := s.seeder.CreateUser()

	err = s.service.Validate(s.Ctx, outsider, ws.ID, []authorization_service.Permission{
		authorization_service.WorkspaceUpdate,
	})

	assert.Error(s.T(), err)
}
