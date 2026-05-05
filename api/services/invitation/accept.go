package invitation_service

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"context"
	"time"

	"github.com/google/uuid"
)

// AcceptInvitation accepts a pending invitation and creates the workspace relationship
func (s *Service) AcceptInvitation(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	invitationRepo := repositories.NewInvitationRepository(s.DB)
	invitation, err := invitationRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if invitation == nil || invitation.GetState() != models.InvitationStatePending {
		return ErrInvalidInvitation
	}

	if time.Now().After(invitation.ExpiresAt) {
		return ErrInvitationExpired
	}

	// Mark invitation as accepted
	if err = invitationRepo.Accept(ctx, id); err != nil {
		return err
	}

	workspaceRepo := repositories.NewWorkspaceRepository(s.DB)
	// Create workspace relationship
	relationship := &models.UserWorkspaceRelationship{
		UserID:      userID,
		WorkspaceID: invitation.WorkspaceID,
		Role:        invitation.Role,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err = workspaceRepo.AddUserToWorkspace(ctx, relationship); err != nil {
		return err
	}

	return nil
}

// IgnoreInvitation marks an invitation as ignored
func (s *Service) IgnoreInvitation(ctx context.Context, id uuid.UUID) error {
	invitationRepo := repositories.NewInvitationRepository(s.DB)
	invitation, err := invitationRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if invitation == nil || invitation.GetState() != models.InvitationStatePending {
		return ErrInvalidInvitation
	}

	return invitationRepo.Ignore(ctx, id)
}
