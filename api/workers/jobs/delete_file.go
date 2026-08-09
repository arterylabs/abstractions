package jobs

import (
	"abstractions/config"
	"abstractions/logger"
	"abstractions/repositories"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

const DeleteFileJobID = "delete_file"

// deleteFilePayload mirrors the payload published by the webhook handler.
type deleteFilePayload struct {
	RepositoryID string `json:"repositoryId"`
	FilePath     string `json:"filePath"`
}

// DeleteFileJob removes all stored data for a file that was deleted in a push event:
// the repository_files row and all associated chunk_embeddings.
type DeleteFileJob struct {
	db                 *gorm.DB
	config             *config.Config
	fileIndexRepo      *repositories.RepositoryFileRepository
	chunkEmbeddingRepo *repositories.ChunkEmbeddingRepository
}

func NewDeleteFileJob(db *gorm.DB, cfg *config.Config) *DeleteFileJob {
	return &DeleteFileJob{
		db:                 db,
		config:             cfg,
		fileIndexRepo:      repositories.NewRepositoryFileRepository(db),
		chunkEmbeddingRepo: repositories.NewChunkEmbeddingRepository(db),
	}
}

func (j *DeleteFileJob) ID() string { return DeleteFileJobID }

// NewDeleteFileTask builds an asynq.Task for deleting a file's index data.
func NewDeleteFileTask(repositoryID uuid.UUID, filePath string) (*asynq.Task, error) {
	payload := deleteFilePayload{
		RepositoryID: repositoryID.String(),
		FilePath:     filePath,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(DeleteFileJobID, b), nil
}

// Handler processes a DeleteFileJob task:
//  1. Decode payload
//  2. Delete all chunk_embeddings for the file
//  3. Delete the repository_files row
func (j *DeleteFileJob) Handler(ctx context.Context, t *asynq.Task) error {
	var payload deleteFilePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("delete_file: decode payload: %w", err)
	}

	repoID, err := uuid.Parse(payload.RepositoryID)
	if err != nil {
		return fmt.Errorf("delete_file: invalid repo id: %w", err)
	}

	logger.Info(ctx, "delete_file: starting filePath="+payload.FilePath)

	if err := j.chunkEmbeddingRepo.DeleteByFile(ctx, repoID, payload.FilePath); err != nil {
		return fmt.Errorf("delete_file: delete chunk_embeddings: %w", err)
	}

	if err := j.fileIndexRepo.DeleteByPath(ctx, repoID, payload.FilePath); err != nil {
		return fmt.Errorf("delete_file: delete repository_files: %w", err)
	}

	logger.Info(ctx, "delete_file: done filePath="+payload.FilePath)
	return nil
}
