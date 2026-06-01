package user_service

import (
	workspace_service "abstractions/services/workspace"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *Service) UpdateDefaultWorkspace(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID) error {
	workspaceRepository := workspace_service.NewService(s.db)

	// First check if the workspace exists and the user has access to it
	_, err := workspaceRepository.GetWorkspace(ctx, workspaceID)
	if err != nil {
		return fmt.Errorf("user_service: update_default_workspace: get_workspace: %w", err)
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
