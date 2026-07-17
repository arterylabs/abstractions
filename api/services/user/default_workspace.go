package user_service

import (
	"abstractions/repositories"
	workspace_service "abstractions/services/workspace"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *Service) UpdateDefaultWorkspace(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID) error {
	workspaceRepository := workspace_service.NewService(s.db)

	// The workspace must exist and the user must be a member of it. GetWorkspace
	// only checks existence, so verify membership explicitly via the user's role.
	if _, err := workspaceRepository.GetWorkspace(ctx, workspaceID); err != nil {
		return fmt.Errorf("user_service: update_default_workspace: get_workspace: %w", err)
	}

	role, err := repositories.NewWorkspaceRepository(s.db).GetUserRole(ctx, userID, workspaceID)
	if err != nil {
		return fmt.Errorf("user_service: update_default_workspace: get_user_role: %w", err)
	}
	if role == "" {
		return ErrWorkspaceAccessDenied
	}

	// Update the user's default workspace
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user_service: update_default_workspace: get_user: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	user.DefaultWorkspaceID = &workspaceID
	return s.userRepo.Update(ctx, user)
}
