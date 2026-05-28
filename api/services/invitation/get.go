package invitation_service

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/utils/pagination_utils"
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetInvitationByID fetches an invitation by its UUID
func (s *Service) GetInvitationByID(ctx context.Context, id uuid.UUID) (*models.Invitation, error) {
	invitationRepo := repositories.NewInvitationRepository(s.DB)
	invitation, err := invitationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("invitation_service: get_invitation_by_id: %w", err)
	}
	return invitation, nil
}

// GetPendingInvitations gets all pending invitations for a given email
func (s *Service) GetPendingInvitations(ctx context.Context, email string) ([]models.Invitation, error) {
	invitationRepo := repositories.NewInvitationRepository(s.DB)
	invitations, err := invitationRepo.GetPendingByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("invitation_service: get_pending_invitations: get_by_email: %w", err)
	}

	if err := s.enrichInvitations(ctx, invitations, true, "get_pending_invitations"); err != nil {
		return nil, err
	}

	return invitations, nil
}

// GetPendingInvitationsPaginated gets paginated pending invitations for a given email
func (s *Service) GetPendingInvitationsPaginated(ctx context.Context, email string, params pagination_utils.PaginationParams) ([]models.Invitation, int, error) {
	invitationRepo := repositories.NewInvitationRepository(s.DB)
	params.Validate()

	invitations, total, err := invitationRepo.GetPendingByEmailPaginated(ctx, email, params)
	if err != nil {
		return nil, 0, fmt.Errorf("invitation_service: get_pending_invitations_paginated: get_by_email: %w", err)
	}

	if err := s.enrichInvitations(ctx, invitations, true, "get_pending_invitations_paginated"); err != nil {
		return nil, 0, err
	}

	return invitations, total, nil
}

// GetWorkspaceInvitationsPaginated gets paginated invitations for a workspace
func (s *Service) GetWorkspaceInvitationsPaginated(ctx context.Context, workspaceID uuid.UUID, params pagination_utils.PaginationParams) ([]models.Invitation, int, error) {
	invitationRepo := repositories.NewInvitationRepository(s.DB)
	params.Validate()

	invitations, total, err := invitationRepo.GetByWorkspacePaginated(ctx, workspaceID, params)
	if err != nil {
		return nil, 0, fmt.Errorf("invitation_service: get_workspace_invitations_paginated: get_by_workspace: %w", err)
	}

	// Get user information for all invitations in one query
	userRepo := repositories.NewUserRepository(s.DB)
	userIDs := make([]uuid.UUID, len(invitations))
	for i := range invitations {
		userIDs[i] = invitations[i].InvitedBy
	}

	users, err := userRepo.GetByIDs(userIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("invitation_service: get_workspace_invitations_paginated: get_users: %w", err)
	}

	// Create a map of user IDs to users for quick lookup
	userMap := make(map[uuid.UUID]*models.User)
	for i := range users {
		userMap[users[i].ID] = &users[i]
	}

	// Assign users to invitations
	for i := range invitations {
		if user, ok := userMap[invitations[i].InvitedBy]; ok {
			invitations[i].InvitedByUser = user
		}
	}

	return invitations, total, nil
}

// GetWorkspaceInvitations gets all invitations for a workspace (deprecated, use paginated version)
func (s *Service) GetWorkspaceInvitations(ctx context.Context, workspaceID uuid.UUID) ([]models.Invitation, error) {
	invitations, _, err := s.GetWorkspaceInvitationsPaginated(ctx, workspaceID, pagination_utils.PaginationParams{
		Page:    1,
		PerPage: 100, // Use a reasonable default
	})
	if err != nil {
		return nil, fmt.Errorf("invitation_service: get_workspace_invitations: get_paginated: %w", err)
	}

	return invitations, nil
}

func (s *Service) enrichInvitations(ctx context.Context, invitations []models.Invitation, includeWorkspaces bool, operation string) error {
	userIDs := make([]uuid.UUID, len(invitations))
	workspaceIDs := make([]uuid.UUID, len(invitations))
	for i := range invitations {
		userIDs[i] = invitations[i].InvitedBy
		workspaceIDs[i] = invitations[i].WorkspaceID
	}

	userRepo := repositories.NewUserRepository(s.DB)
	users, err := userRepo.GetByIDs(userIDs)
	if err != nil {
		return fmt.Errorf("invitation_service: %s: get_users: %w", operation, err)
	}

	userMap := make(map[uuid.UUID]*models.User)
	for i := range users {
		userMap[users[i].ID] = &users[i]
	}

	var workspaceMap map[uuid.UUID]*models.Workspace
	if includeWorkspaces {
		workspaceRepo := repositories.NewWorkspaceRepository(s.DB)
		workspaces, err := workspaceRepo.GetByIDs(ctx, workspaceIDs)
		if err != nil {
			return fmt.Errorf("invitation_service: %s: get_workspaces: %w", operation, err)
		}

		workspaceMap = make(map[uuid.UUID]*models.Workspace)
		for i := range workspaces {
			workspaceMap[workspaces[i].ID] = &workspaces[i]
		}
	}

	for i := range invitations {
		if user, ok := userMap[invitations[i].InvitedBy]; ok {
			invitations[i].InvitedByUser = user
		}
		if workspace, ok := workspaceMap[invitations[i].WorkspaceID]; ok {
			invitations[i].Workspace = workspace
		}
	}

	return nil
}
