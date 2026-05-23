package workspace_service_test

import (
	authorization_service "abstractions/services/authorization"
	workspace_service "abstractions/services/workspace"
	"abstractions/tests/testutils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WorkspaceServiceTestSuite struct {
	testutils.TestSuite
	seeder  *testutils.Seeder
	service *workspace_service.Service
}

func TestWorkspaceService(t *testing.T) {
	suite.Run(t, new(WorkspaceServiceTestSuite))
}

func (s *WorkspaceServiceTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.service = workspace_service.NewWorkspaceService(s.DB)
}

func (s *WorkspaceServiceTestSuite) TestCreateWorkspace_Success() {
	owner := s.seeder.CreateUser()

	ws, err := s.service.CreateWorkspace(s.Ctx, "Test Workspace", owner.ID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), ws)
	assert.Equal(s.T(), "Test Workspace", ws.Name)
	assert.Equal(s.T(), owner.ID, ws.CreatedBy)
}

func (s *WorkspaceServiceTestSuite) TestCreateWorkspace_CreatorIsOwner() {
	owner := s.seeder.CreateUser()

	ws, err := s.service.CreateWorkspace(s.Ctx, "Owned WS", owner.ID)
	s.Require().NoError(err)

	users, err := s.service.GetWorkspaceUsers(s.Ctx, ws.ID, owner.ID)
	assert.NoError(s.T(), err)

	var found bool
	for _, u := range users {
		if u.ID == owner.ID {
			found = true
			break
		}
	}
	assert.True(s.T(), found, "creator should be a member of the workspace")
}

func (s *WorkspaceServiceTestSuite) TestGetWorkspace_ReturnsCorrectWorkspace() {
	owner := s.seeder.CreateUser()
	ws, err := s.service.CreateWorkspace(s.Ctx, "Fetch WS", owner.ID)
	s.Require().NoError(err)

	fetched, err := s.service.GetWorkspace(s.Ctx, ws.ID)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), ws.ID, fetched.ID)
	assert.Equal(s.T(), "Fetch WS", fetched.Name)
}

func (s *WorkspaceServiceTestSuite) TestGetWorkspace_NotFound() {
	_, err := s.service.GetWorkspace(s.Ctx, uuid.New())

	assert.ErrorIs(s.T(), err, workspace_service.ErrWorkspaceNotFound)
}

func (s *WorkspaceServiceTestSuite) TestAddUserToWorkspace_Success() {
	owner := s.seeder.CreateUser()
	ws, err := s.service.CreateWorkspace(s.Ctx, "Add User WS", owner.ID)
	s.Require().NoError(err)

	newMember := s.seeder.CreateUser()

	err = s.service.AddUserToWorkspace(s.Ctx, ws.ID, newMember.ID, authorization_service.RoleMember.String(), owner.ID)
	assert.NoError(s.T(), err)

	users, err := s.service.GetWorkspaceUsers(s.Ctx, ws.ID, owner.ID)
	assert.NoError(s.T(), err)

	var found bool
	for _, u := range users {
		if u.ID == newMember.ID {
			found = true
			break
		}
	}
	assert.True(s.T(), found, "new member should appear in workspace users")
}

func (s *WorkspaceServiceTestSuite) TestAddUserToWorkspace_MemberCannotGrantRoles() {
	owner := s.seeder.CreateUser()
	ws, err := s.service.CreateWorkspace(s.Ctx, "Role Grant WS", owner.ID)
	s.Require().NoError(err)

	member := s.seeder.CreateUser()
	s.Require().NoError(
		s.service.AddUserToWorkspace(s.Ctx, ws.ID, member.ID, authorization_service.RoleMember.String(), owner.ID),
	)

	outsider := s.seeder.CreateUser()
	err = s.service.AddUserToWorkspace(s.Ctx, ws.ID, outsider.ID, authorization_service.RoleMember.String(), member.ID)
	assert.ErrorIs(s.T(), err, workspace_service.ErrInsufficientPermissions)
}

func (s *WorkspaceServiceTestSuite) TestUpdateWorkspace_PersistsName() {
	owner := s.seeder.CreateUser()
	ws, err := s.service.CreateWorkspace(s.Ctx, "Original Name", owner.ID)
	s.Require().NoError(err)

	err = s.service.UpdateWorkspace(s.Ctx, ws.ID, workspace_service.UpdateWorkspaceParams{
		Name: "Updated Name",
	}, owner.ID)
	assert.NoError(s.T(), err)

	fetched, err := s.service.GetWorkspace(s.Ctx, ws.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Updated Name", fetched.Name)
}

func (s *WorkspaceServiceTestSuite) TestRemoveUserFromWorkspace_Success() {
	owner := s.seeder.CreateUser()
	ws, err := s.service.CreateWorkspace(s.Ctx, "Remove User WS", owner.ID)
	s.Require().NoError(err)

	member := s.seeder.CreateUser()
	s.Require().NoError(
		s.service.AddUserToWorkspace(s.Ctx, ws.ID, member.ID, authorization_service.RoleMember.String(), owner.ID),
	)

	err = s.service.RemoveUserFromWorkspace(s.Ctx, ws.ID, member.ID, owner.ID)
	assert.NoError(s.T(), err)

	users, err := s.service.GetWorkspaceUsers(s.Ctx, ws.ID, owner.ID)
	assert.NoError(s.T(), err)

	for _, u := range users {
		assert.NotEqual(s.T(), member.ID, u.ID, "removed user should not appear in workspace users")
	}
}

func (s *WorkspaceServiceTestSuite) TestRemoveUserFromWorkspace_CannotRemoveSelf() {
	owner := s.seeder.CreateUser()
	ws, err := s.service.CreateWorkspace(s.Ctx, "Self Remove WS", owner.ID)
	s.Require().NoError(err)

	err = s.service.RemoveUserFromWorkspace(s.Ctx, ws.ID, owner.ID, owner.ID)
	assert.ErrorIs(s.T(), err, workspace_service.ErrCannotModifyOwnRole)
}
