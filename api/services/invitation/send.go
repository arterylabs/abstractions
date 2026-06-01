package invitation_service

import (
	"abstractions/clients/email"
	"abstractions/config"
	"abstractions/db/models"
	"abstractions/repositories"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SendInvitations sends workspace invitations to multiple users
// Returns an array of errors corresponding to each invitation attempt
func (s *Service) SendInvitations(ctx context.Context, workspaceID uuid.UUID, invitedByUserID uuid.UUID, inputs []InvitationInput) []error {
	errors := make([]error, len(inputs))

	workspaceRepo := repositories.NewWorkspaceRepository(s.DB)
	// Fetch workspace details
	workspace, err := workspaceRepo.GetByID(ctx, workspaceID)
	if err != nil {
		for i := range errors {
			errors[i] = err
		}
		return errors
	}

	// Fetch inviter details
	userRepo := repositories.NewUserRepository(s.DB)
	invitedByUser, err := userRepo.GetByID(ctx, invitedByUserID)
	if err != nil {
		for i := range errors {
			errors[i] = err
		}
		return errors
	}

	invitationRepo := repositories.NewInvitationRepository(s.DB)

	for i, input := range inputs {
		invitation := &models.Invitation{
			WorkspaceID: workspaceID,
			InvitedBy:   invitedByUserID,
			Email:       input.Email,
			Role:        input.Role,
			ExpiresAt:   time.Now().UTC().Add(7 * 24 * time.Hour), // 7 days expiry
		}
		errors[i] = invitationRepo.Create(ctx, invitation)
		if errors[i] == nil {
			// Send invitation email
			go s.sendInvitationEmail(ctx, invitation, workspace, invitedByUser)
		}
	}
	return errors
}

// sendInvitationEmail sends an invitation email to the user
func (s *Service) sendInvitationEmail(ctx context.Context, invitation *models.Invitation, workspace *models.Workspace, invitedByUser *models.User) {
	bodyHTML := fmt.Sprintf("<p>%s (%s) has invited you to join the workspace <strong>%s</strong> on Abstractions as a %s.</p>", invitedByUser.Name, invitedByUser.Email, workspace.Name, invitation.Role)
	bodyPlainText := fmt.Sprintf("%s (%s) has invited you to join the workspace %s on Abstractions as a %s.", invitedByUser.Name, invitedByUser.Email, workspace.Name, invitation.Role)

	email.SendGenericEmail(ctx, email.SendGenericEmailParams{
		To:            invitation.Email,
		Subject:       fmt.Sprintf("Invitation to join %s on Abstractions", workspace.Name),
		BodyHTML:      bodyHTML,
		BodyPlainText: bodyPlainText,
		CTAText:       "Accept Invitation",
		CTALink:       fmt.Sprintf("%s/me/invitations?id=%s", config.Environment.ClientURL, invitation.ID),
		Reason:        fmt.Sprintf("You are receiving this email because %s has invited you to join their workspace on Abstractions", invitedByUser.Name),
	})
}

// ResendInvitation resends an invitation email
func (s *Service) ResendInvitation(ctx context.Context, id uuid.UUID, invitedByUserID uuid.UUID) error {
	invitationRepo := repositories.NewInvitationRepository(s.DB)
	invitation, err := invitationRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("invitation_service: resend_invitation: get_invitation: %w", err)
	}

	// Check if invitation is valid
	if invitation == nil || invitation.AcceptedAt != nil || invitation.IgnoredAt != nil {
		return ErrInvalidInvitation
	}

	if invitation.ExpiresAt.Before(time.Now()) {
		return ErrInvitationExpired
	}

	// Fetch workspace and user details for email
	workspaceRepo := repositories.NewWorkspaceRepository(s.DB)
	workspace, err := workspaceRepo.GetByID(ctx, invitation.WorkspaceID)
	if err != nil {
		return fmt.Errorf("invitation_service: resend_invitation: get_workspace: %w", err)
	}

	userRepo := repositories.NewUserRepository(s.DB)
	invitedByUser, err := userRepo.GetByID(ctx, invitedByUserID)
	if err != nil {
		return fmt.Errorf("invitation_service: resend_invitation: get_user: %w", err)
	}

	// Send invitation email
	go s.sendInvitationEmail(ctx, invitation, workspace, invitedByUser)

	return nil
}
