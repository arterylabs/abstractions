package invitation_service

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/utils/pagination_utils"
	"context"

	"github.com/google/uuid"
)

// GetPendingInvitations gets all pending invitations for a given email
func (s *Service) GetPendingInvitations(ctx context.Context, email string) ([]models.Invitation, error) {
	invitationRepo := repositories.NewInvitationRepository(s.DB)
	invitations, err := invitationRepo.GetPendingByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Get user information for all invitations in one query
	userRepo := repositories.NewUserRepository(s.DB)
	userIDs := make([]uuid.UUID, len(invitations))
	workspaceIDs := make([]uuid.UUID, len(invitations))
	for i := range invitations {
		userIDs[i] = invitations[i].InvitedBy
		workspaceIDs[i] = invitations[i].WorkspaceID
	}

	users, err := userRepo.GetByIDs(userIDs)
	if err != nil {
		return nil, err
	}

	// Get workspace information
	workspaceRepo := repositories.NewWorkspaceRepository(s.DB)
	workspaces, err := workspaceRepo.GetByIDs(ctx, workspaceIDs)
	if err != nil {
		return nil, err
	}

	// Create maps for quick lookups
	userMap := make(map[uuid.UUID]*models.User)
	for i := range users {
		userMap[users[i].ID] = &users[i]
	}

	workspaceMap := make(map[uuid.UUID]*models.Workspace)
	for i := range workspaces {
		workspaceMap[workspaces[i].ID] = &workspaces[i]
	}

	// Assign users and workspaces to invitations
	for i := range invitations {
		if user, ok := userMap[invitations[i].InvitedBy]; ok {
			invitations[i].InvitedByUser = user
		}
		if workspace, ok := workspaceMap[invitations[i].WorkspaceID]; ok {
			invitations[i].Workspace = workspace
		}
	}

	return invitations, nil
}

// GetPendingInvitationsPaginated gets paginated pending invitations for a given email
func (s *Service) GetPendingInvitationsPaginated(ctx context.Context, email string, params pagination_utils.PaginationParams) ([]models.Invitation, int, error) {
	invitationRepo := repositories.NewInvitationRepository(s.DB)
	params.Validate()

	invitations, total, err := invitationRepo.GetPendingByEmailPaginated(ctx, email, params)
	if err != nil {
		return nil, 0, err
	}

	// Get user information for all invitations in one query
	userRepo := repositories.NewUserRepository(s.DB)
	userIDs := make([]uuid.UUID, len(invitations))
	workspaceIDs := make([]uuid.UUID, len(invitations))
	for i := range invitations {
		userIDs[i] = invitations[i].InvitedBy
		workspaceIDs[i] = invitations[i].WorkspaceID
	}

	users, err := userRepo.GetByIDs(userIDs)
	if err != nil {
		return nil, 0, err
	}

	// Get workspace information
	workspaceRepo := repositories.NewWorkspaceRepository(s.DB)
	workspaces, err := workspaceRepo.GetByIDs(ctx, workspaceIDs)
	if err != nil {
		return nil, 0, err
	}

	// Create maps for quick lookups
	userMap := make(map[uuid.UUID]*models.User)
	for i := range users {
		userMap[users[i].ID] = &users[i]
	}

	workspaceMap := make(map[uuid.UUID]*models.Workspace)
	for i := range workspaces {
		workspaceMap[workspaces[i].ID] = &workspaces[i]
	}

	// Assign users and workspaces to invitations
	for i := range invitations {
		if user, ok := userMap[invitations[i].InvitedBy]; ok {
			invitations[i].InvitedByUser = user
		}
		if workspace, ok := workspaceMap[invitations[i].WorkspaceID]; ok {
			invitations[i].Workspace = workspace
		}
	}

	return invitations, total, nil
}

// GetWorkspaceInvitationsPaginated gets paginated invitations for a workspace
func (s *Service) GetWorkspaceInvitationsPaginated(ctx context.Context, workspaceID uuid.UUID, params pagination_utils.PaginationParams) ([]models.Invitation, int, error) {
	invitationRepo := repositories.NewInvitationRepository(s.DB)
	params.Validate()

	invitations, total, err := invitationRepo.GetByWorkspacePaginated(ctx, workspaceID, params)
	if err != nil {
		return nil, 0, err
	}

	// Get user information for all invitations in one query
	userRepo := repositories.NewUserRepository(s.DB)
	userIDs := make([]uuid.UUID, len(invitations))
	for i := range invitations {
		userIDs[i] = invitations[i].InvitedBy
	}

	users, err := userRepo.GetByIDs(userIDs)
	if err != nil {
		return nil, 0, err
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
		return nil, err
	}

	return invitations, nil
}
