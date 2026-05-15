package user_service

import (
	"abstractions/repositories"
	workspace_service "abstractions/services/workspace"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *Service) UpdateDefaultWorkspace(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID) error {
	userRepository := repositories.NewUserRepository(s.DB)
	workspaceRepository := repositories.NewWorkspaceRepository(s.DB)

	// First check if the workspace exists and the user has access to it
	workspace, err := workspaceRepository.GetByID(ctx, workspaceID)
	if err != nil {
		return fmt.Errorf("user_service: update_default_workspace: get_workspace: %w", err)
	}
	if workspace == nil {
		return workspace_service.ErrWorkspaceNotFound
	}

	// Update the user's default workspace
	user, err := userRepository.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user_service: update_default_workspace: get_user: %w", err)
	}

	user.DefaultWorkspaceID = &workspaceID
	return userRepository.Update(user)
}
