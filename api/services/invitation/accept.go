package invitation_service

import (
	"abstractions/db/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AcceptInvitation accepts a pending invitation and creates the workspace relationship
func (s *Service) AcceptInvitation(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	invitation, err := s.invitationRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("invitation_service: accept_invitation: get_invitation: %w", err)
	}

	if invitation == nil || invitation.GetState() != models.InvitationStatePending {
		return ErrInvalidInvitation
	}

	// Mark invitation as accepted
	if err = s.invitationRepo.Accept(ctx, id); err != nil {
		return fmt.Errorf("invitation_service: accept_invitation: accept: %w", err)
	}

	// Create workspace relationship
	relationship := &models.UserWorkspaceRelationship{
		UserID:      userID,
		WorkspaceID: invitation.WorkspaceID,
		Role:        invitation.Role,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err = s.workspaceRepo.AddUserToWorkspace(ctx, relationship); err != nil {
		return fmt.Errorf("invitation_service: accept_invitation: add_user_to_workspace: %w", err)
	}

	return nil
}

// IgnoreInvitation marks an invitation as ignored
func (s *Service) IgnoreInvitation(ctx context.Context, id uuid.UUID) error {
	invitation, err := s.invitationRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("invitation_service: ignore_invitation: get_invitation: %w", err)
	}

	if invitation == nil || invitation.GetState() != models.InvitationStatePending {
		return ErrInvalidInvitation
	}

	return s.invitationRepo.Ignore(ctx, id)
}
