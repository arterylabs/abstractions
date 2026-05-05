package user_service

import (
	"abstractions/repositories"
	"context"
	"errors"

	"github.com/google/uuid"
)

func (s *Service) UpdateDefaultWorkspace(ctx context.Context, userID uuid.UUID, workspaceID uuid.UUID) error {
	userRepository := repositories.NewUserRepository(s.DB)
	workspaceRepository := repositories.NewWorkspaceRepository(s.DB)

	// First check if the workspace exists and the user has access to it
	workspace, err := workspaceRepository.GetByID(ctx, workspaceID)
	if err != nil {
		return err
	}
	if workspace == nil {
		return errors.New("workspace not found")
	}

	// Update the user's default workspace
	user, err := userRepository.GetByID(userID)
	if err != nil {
		return err
	}

	user.DefaultWorkspaceID = &workspaceID
	return userRepository.Update(user)
}
