package workspace_service

import (
	"abstractions/clients/storage"
	"abstractions/db/models"
	"abstractions/logger"
	"abstractions/repositories"
	"context"
	"mime/multipart"

	"github.com/google/uuid"
)

func (s *Service) UpdateWorkspaceLogo(ctx context.Context, workspaceID uuid.UUID, file multipart.File) (*models.Workspace, error) {
	storageClient, err := storage.NewMinioStorageClient()
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	filePath := "workspace-logos/" + workspaceID.String()
	fileLink, err := storageClient.UploadPublic(ctx, filePath, file)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	workspaceRepo := repositories.NewWorkspaceRepository(s.DB)

	workspace, err := workspaceRepo.GetByID(ctx, workspaceID)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	workspace.Logo = &fileLink
	err = workspaceRepo.Update(ctx, workspace)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	return workspace, nil
}
