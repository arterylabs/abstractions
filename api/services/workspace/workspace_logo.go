package workspace_service

import (
	"abstractions/clients/storage"
	"abstractions/db/models"
	"abstractions/logger"
	"context"
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
)

func (s *Service) UpdateWorkspaceLogo(ctx context.Context, workspaceID uuid.UUID, file multipart.File) (*models.Workspace, error) {
	storageClient, err := storage.NewMinioStorageClient()
	if err != nil {
		logger.Error(ctx, err)
		return nil, fmt.Errorf("workspace_service: update_workspace_logo: create_storage_client: %w", err)
	}

	filePath := "workspace-logos/" + workspaceID.String()
	fileLink, err := storageClient.UploadPublic(ctx, filePath, file)
	if err != nil {
		logger.Error(ctx, err)
		return nil, fmt.Errorf("workspace_service: update_workspace_logo: upload: %w", err)
	}

	workspace, err := s.workspaceRepo.GetByID(ctx, workspaceID)
	if err != nil {
		logger.Error(ctx, err)
		return nil, fmt.Errorf("workspace_service: update_workspace_logo: get_workspace: %w", err)
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	workspace.Logo = &fileLink
	err = s.workspaceRepo.Update(ctx, workspace)
	if err != nil {
		logger.Error(ctx, err)
		return nil, fmt.Errorf("workspace_service: update_workspace_logo: update: %w", err)
	}

	return workspace, nil
}
