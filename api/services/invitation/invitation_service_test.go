package invitation_service_test

import (
	"abstractions/db/models"
	"abstractions/repositories"
	invitation_service "abstractions/services/invitation"
	workspace_service "abstractions/services/workspace"
	"abstractions/tests/testutils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type InvitationServiceTestSuite struct {
	testutils.TestSuite
	seeder         *testutils.Seeder
	service        *invitation_service.Service
	workspaceSvc   *workspace_service.Service
	invitationRepo *repositories.InvitationRepository
}

func TestInvitationService(t *testing.T) {
	suite.Run(t, new(InvitationServiceTestSuite))
}

func (s *InvitationServiceTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.service = invitation_service.NewInvitationService(s.DB)
	s.workspaceSvc = workspace_service.NewWorkspaceService(s.DB)
	s.invitationRepo = repositories.NewInvitationRepository(s.DB)
}

func (s *InvitationServiceTestSuite) seedPendingInvitation(workspaceID, invitedBy uuid.UUID, role string) *models.Invitation {
	inv := &models.Invitation{
		ID:          uuid.New(),
		WorkspaceID: workspaceID,
		InvitedBy:   invitedBy,
		Email:       "invite-" + uuid.New().String() + "@example.com",
		Role:        role,
		ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
	}
	s.Require().NoError(s.invitationRepo.Create(s.Ctx, inv))
	return inv
}

func (s *InvitationServiceTestSuite) TestCreateInvitation_Success() {
	owner := s.seeder.CreateUser()
	ws, err := s.workspaceSvc.CreateWorkspace(s.Ctx, "Invite WS", owner.ID)
	s.Require().NoError(err)

	errs := s.service.SendInvitations(s.Ctx, ws.ID, owner.ID, []invitation_service.InvitationInput{
		{Email: "new-" + uuid.New().String() + "@example.com", Role: "member"},
	})

	assert.Len(s.T(), errs, 1)
	assert.NoError(s.T(), errs[0])
}

func (s *InvitationServiceTestSuite) TestAcceptInvitation_Success() {
	owner := s.seeder.CreateUser()
	ws, err := s.workspaceSvc.CreateWorkspace(s.Ctx, "Accept WS", owner.ID)
	s.Require().NoError(err)

	invitee := s.seeder.CreateUser()
	inv := s.seedPendingInvitation(ws.ID, owner.ID, "member")

	err = s.service.AcceptInvitation(s.Ctx, inv.ID, invitee.ID)
	assert.NoError(s.T(), err)

	// Verify the workspace relationship was created
	users, err := s.workspaceSvc.GetWorkspaceUsers(s.Ctx, ws.ID, owner.ID)
	assert.NoError(s.T(), err)

	var found bool
	for _, u := range users {
		if u.ID == invitee.ID {
			found = true
			break
		}
	}
	assert.True(s.T(), found, "accepted invitee should be in workspace users")
}

func (s *InvitationServiceTestSuite) TestAcceptInvitation_AlreadyAccepted() {
	owner := s.seeder.CreateUser()
	ws, err := s.workspaceSvc.CreateWorkspace(s.Ctx, "Double Accept WS", owner.ID)
	s.Require().NoError(err)

	invitee := s.seeder.CreateUser()
	inv := s.seedPendingInvitation(ws.ID, owner.ID, "member")

	s.Require().NoError(s.service.AcceptInvitation(s.Ctx, inv.ID, invitee.ID))

	// Accept again — should fail because invitation is no longer pending
	err = s.service.AcceptInvitation(s.Ctx, inv.ID, invitee.ID)
	assert.ErrorIs(s.T(), err, invitation_service.ErrInvalidInvitation)
}

func (s *InvitationServiceTestSuite) TestAcceptInvitation_ExpiredInvitation() {
	owner := s.seeder.CreateUser()
	ws, err := s.workspaceSvc.CreateWorkspace(s.Ctx, "Expired Invite WS", owner.ID)
	s.Require().NoError(err)

	// Seed an already-expired invitation
	inv := &models.Invitation{
		ID:          uuid.New(),
		WorkspaceID: ws.ID,
		InvitedBy:   owner.ID,
		Email:       "expired-" + uuid.New().String() + "@example.com",
		Role:        "member",
		ExpiresAt:   time.Now().Add(-24 * time.Hour), // already expired
	}
	s.Require().NoError(s.invitationRepo.Create(s.Ctx, inv))

	invitee := s.seeder.CreateUser()
	err = s.service.AcceptInvitation(s.Ctx, inv.ID, invitee.ID)

	// GetState() returns InvitationStateExpired (not Pending), so AcceptInvitation
	// hits the "not pending" guard first and returns ErrInvalidInvitation.
	assert.ErrorIs(s.T(), err, invitation_service.ErrInvalidInvitation)
}

func (s *InvitationServiceTestSuite) TestGetInvitationByID_ReturnsCorrectInvitation() {
	owner := s.seeder.CreateUser()
	ws, err := s.workspaceSvc.CreateWorkspace(s.Ctx, "Get Invite WS", owner.ID)
	s.Require().NoError(err)

	inv := s.seedPendingInvitation(ws.ID, owner.ID, "member")

	fetched, err := s.service.GetInvitationByID(s.Ctx, inv.ID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), fetched)
	assert.Equal(s.T(), inv.ID, fetched.ID)
}

func (s *InvitationServiceTestSuite) TestGetInvitationByID_NotFound() {
	fetched, err := s.service.GetInvitationByID(s.Ctx, uuid.New())

	assert.NoError(s.T(), err)
	assert.Nil(s.T(), fetched, "non-existent invitation should return nil")
}

func (s *InvitationServiceTestSuite) TestIgnoreInvitation_MarksAsIgnored() {
	owner := s.seeder.CreateUser()
	ws, err := s.workspaceSvc.CreateWorkspace(s.Ctx, "Ignore WS", owner.ID)
	s.Require().NoError(err)

	inv := s.seedPendingInvitation(ws.ID, owner.ID, "member")

	err = s.service.IgnoreInvitation(s.Ctx, inv.ID)
	assert.NoError(s.T(), err)

	fetched, err := s.service.GetInvitationByID(s.Ctx, inv.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), models.InvitationStateIgnored, fetched.GetState())
}
