package tests

import (
	invitation_service "abstractions/services/invitation"
	"abstractions/tests/testutils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type InvitationTestSuite struct {
	testutils.TestSuite
	seeder            *testutils.Seeder
	invitationService *invitation_service.Service
}

func TestInvitationService(t *testing.T) {
	suite.Run(t, new(InvitationTestSuite))
}

func (s *InvitationTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.invitationService = invitation_service.NewInvitationService(s.DB)
}

// Test sending invitations
func (s *InvitationTestSuite) TestSendInvitations() {
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))
	invitedByUser := s.seeder.CreateUser(testutils.WithEmail("inviter@example.com"))

	// Test successful invitation creation
	inputs := []invitation_service.InvitationInput{
		{
			Email: "invited1@example.com",
			Role:  "member",
		},
		{
			Email: "invited2@example.com",
			Role:  "admin",
		},
	}

	errors := s.invitationService.SendInvitations(s.Ctx, workspace.ID, invitedByUser.ID, inputs)
	assert.Equal(s.T(), 2, len(errors))
	for _, err := range errors {
		assert.NoError(s.T(), err)
	}
}

// Test getting pending invitations
func (s *InvitationTestSuite) TestGetPendingInvitations() {
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))
	invitedByUser := s.seeder.CreateUser(testutils.WithEmail("inviter@example.com"))
	email := "invited@example.com"

	inputs := []invitation_service.InvitationInput{
		{
			Email: email,
			Role:  "member",
		},
	}

	errors := s.invitationService.SendInvitations(s.Ctx, workspace.ID, invitedByUser.ID, inputs)
	assert.NoError(s.T(), errors[0])

	// Test retrieving pending invitations
	invitations, err := s.invitationService.GetPendingInvitations(s.Ctx, email)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, len(invitations))
	assert.Equal(s.T(), email, invitations[0].Email)
}

// Test accepting an invitation
func (s *InvitationTestSuite) TestAcceptInvitation() {
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))
	invitedByUser := s.seeder.CreateUser(testutils.WithEmail("inviter@example.com"))
	invitedUser := s.seeder.CreateUser(testutils.WithEmail("invited@example.com"))

	inputs := []invitation_service.InvitationInput{
		{
			Email: invitedUser.Email,
			Role:  "member",
		},
	}

	errors := s.invitationService.SendInvitations(s.Ctx, workspace.ID, invitedByUser.ID, inputs)
	assert.NoError(s.T(), errors[0])

	// Get the invitation
	invitations, err := s.invitationService.GetPendingInvitations(s.Ctx, invitedUser.Email)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, len(invitations))

	// Test accepting invitation
	err = s.invitationService.AcceptInvitation(s.Ctx, invitations[0].ID, invitedUser.ID)
	assert.NoError(s.T(), err)

	// Verify invitation is no longer pending
	pendingInvitations, err := s.invitationService.GetPendingInvitations(s.Ctx, invitedUser.Email)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, len(pendingInvitations))

	// Test accepting non-existent invitation
	err = s.invitationService.AcceptInvitation(s.Ctx, uuid.New(), invitedUser.ID)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), invitation_service.ErrInvalidInvitation, err)
}

// Test ignoring an invitation
func (s *InvitationTestSuite) TestIgnoreInvitation() {
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))
	invitedByUser := s.seeder.CreateUser(testutils.WithEmail("inviter@example.com"))
	email := "invited@example.com"

	inputs := []invitation_service.InvitationInput{
		{
			Email: email,
			Role:  "member",
		},
	}

	errors := s.invitationService.SendInvitations(s.Ctx, workspace.ID, invitedByUser.ID, inputs)
	assert.NoError(s.T(), errors[0])

	// Get the invitation
	invitations, err := s.invitationService.GetPendingInvitations(s.Ctx, email)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, len(invitations))

	// Test ignoring invitation
	err = s.invitationService.IgnoreInvitation(s.Ctx, invitations[0].ID)
	assert.NoError(s.T(), err)

	// Verify invitation is no longer pending
	pendingInvitations, err := s.invitationService.GetPendingInvitations(s.Ctx, email)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, len(pendingInvitations))

	// Test ignoring non-existent invitation
	err = s.invitationService.IgnoreInvitation(s.Ctx, uuid.New())
	assert.Error(s.T(), err)
	assert.Equal(s.T(), invitation_service.ErrInvalidInvitation, err)
}

// Test resending an invitation
func (s *InvitationTestSuite) TestResendInvitation() {
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))
	invitedByUser := s.seeder.CreateUser(testutils.WithEmail("inviter@example.com"))
	email := "invited@example.com"

	inputs := []invitation_service.InvitationInput{
		{
			Email: email,
			Role:  "member",
		},
	}

	errors := s.invitationService.SendInvitations(s.Ctx, workspace.ID, invitedByUser.ID, inputs)
	assert.NoError(s.T(), errors[0])

	// Get the invitation
	invitations, err := s.invitationService.GetPendingInvitations(s.Ctx, email)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, len(invitations))

	// Test resending invitation
	err = s.invitationService.ResendInvitation(s.Ctx, invitations[0].ID, invitedByUser.ID)
	assert.NoError(s.T(), err)

	// Test resending non-existent invitation
	err = s.invitationService.ResendInvitation(s.Ctx, uuid.New(), invitedByUser.ID)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), invitation_service.ErrInvalidInvitation, err)
}
