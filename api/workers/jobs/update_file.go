package jobs

import (
	ai_client "abstractions/clients/ai"
	"abstractions/config"
	"abstractions/db/models"
	"abstractions/logger"
	"abstractions/repositories"
	"abstractions/services/chunker"
	integration_service "abstractions/services/integrations"
	github_integration "abstractions/services/integrations/core/github"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	pgvector "github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

const UpdateFileJobID = "update_file"

// updateFilePayload mirrors the payload published by the webhook handler.
type updateFilePayload struct {
	RepositoryID string `json:"repositoryId"`
	FilePath     string `json:"filePath"`
	SHA          string `json:"sha"`
}

// UpdateFileJob checks whether a pushed file has a new git SHA; if changed it deletes the
// old chunk_embeddings and re-runs the full summarize pipeline.
type UpdateFileJob struct {
	db                 *gorm.DB
	config             *config.Config
	aiProvider         ai_client.AIProvider
	fileFetcher        FileContentFetcher // reuse interface from summarize_file.go
	integrationService *integration_service.Service
	fileIndexRepo      *repositories.RepositoryFileRepository
	chunkEmbeddingRepo *repositories.ChunkEmbeddingRepository
	repositoryRepo     *repositories.RepositoryRepository
}

func NewUpdateFileJob(db *gorm.DB, cfg *config.Config) *UpdateFileJob {
	return &UpdateFileJob{
		db:                 db,
		config:             cfg,
		aiProvider:         ai_client.NewAIProvider(),
		fileFetcher:        &defaultFileContentFetcher{},
		integrationService: integration_service.NewIntegrationService(db),
		fileIndexRepo:      repositories.NewRepositoryFileRepository(db),
		chunkEmbeddingRepo: repositories.NewChunkEmbeddingRepository(db),
		repositoryRepo:     repositories.NewRepositoryRepository(db),
	}
}

// NewUpdateFileJobWithDeps creates a job with injected dependencies for testing.
func NewUpdateFileJobWithDeps(
	db *gorm.DB,
	cfg *config.Config,
	aiProvider ai_client.AIProvider,
	fileFetcher FileContentFetcher,
) *UpdateFileJob {
	return &UpdateFileJob{
		db:                 db,
		config:             cfg,
		aiProvider:         aiProvider,
		fileFetcher:        fileFetcher,
		integrationService: integration_service.NewIntegrationService(db),
		fileIndexRepo:      repositories.NewRepositoryFileRepository(db),
		chunkEmbeddingRepo: repositories.NewChunkEmbeddingRepository(db),
		repositoryRepo:     repositories.NewRepositoryRepository(db),
	}
}

func (j *UpdateFileJob) ID() string { return UpdateFileJobID }

// Handler processes an UpdateFileJob task:
//  1. Decode payload
//  2. Look up repository_files; if git_sha unchanged → return early (no API calls)
//  3. If changed: delete old chunk_embeddings, run summarize pipeline, update git_sha
func (j *UpdateFileJob) Handler(ctx context.Context, t *asynq.Task) error {
	var payload updateFilePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("update_file: decode payload: %w", err)
	}

	repoID, err := uuid.Parse(payload.RepositoryID)
	if err != nil {
		return fmt.Errorf("update_file: invalid repo id: %w", err)
	}

	logger.Info(ctx, "update_file: starting filePath="+payload.FilePath)

	// SHA check: skip if unchanged
	fi, err := j.fileIndexRepo.GetByRepoAndPath(ctx, repoID, payload.FilePath)
	if err != nil {
		return fmt.Errorf("update_file: lookup repository_files: %w", err)
	}
	if fi != nil && fi.GitSHA == payload.SHA {
		logger.Info(ctx, "update_file: sha unchanged, skipping filePath="+payload.FilePath)
		return nil
	}

	// Delete old chunk embeddings so stale chunks don't linger after re-processing
	if err := j.chunkEmbeddingRepo.DeleteByFile(ctx, repoID, payload.FilePath); err != nil {
		return fmt.Errorf("update_file: delete old embeddings: %w", err)
	}

	// Resolve GitHub access token
	repo, err := j.repositoryRepo.GetByID(ctx, repoID)
	if err != nil || repo == nil {
		return fmt.Errorf("update_file: repo not found %s: %w", repoID, err)
	}
	integration, err := j.integrationService.GetIntegrationByType(ctx, repo.WorkspaceID, github_integration.IntegrationTypeGitHub)
	if err != nil {
		return fmt.Errorf("update_file: github integration not found: %w", err)
	}
	accessToken, err := j.integrationService.GetDecryptedToken(ctx, integration.ID)
	if err != nil {
		return fmt.Errorf("update_file: decrypt token: %w", err)
	}

	owner, repoName, err := splitOwnerRepo(repo.Name)
	if err != nil {
		return fmt.Errorf("update_file: %w", err)
	}

	rawContent, err := j.fileFetcher.GetFileContent(ctx, accessToken, owner, repoName, payload.FilePath, payload.SHA)
	if err != nil {
		return fmt.Errorf("update_file: fetch content: %w", err)
	}

	language := chunker.DetectLanguage(payload.FilePath)
	chunks := chunker.ChunkFile(payload.FilePath, string(rawContent), language)
	rawContent = nil // clear content immediately

	if len(chunks) == 0 {
		logger.Info(ctx, "update_file: no chunks, skipping filePath="+payload.FilePath)
		return nil
	}

	embeddings := make([]models.ChunkEmbedding, 0, len(chunks))
	allTags := map[string]struct{}{}

	for i, chunk := range chunks {
		summary, err := j.aiProvider.SummarizeCode(ctx, payload.FilePath, language, chunk.Content)
		if err != nil {
			logger.Info(ctx, fmt.Sprintf("update_file: SummarizeCode chunk %d error: %v", i, err))
			continue
		}
		embedding, err := j.aiProvider.GenerateEmbedding(ctx, summary.Summary)
		if err != nil {
			logger.Info(ctx, fmt.Sprintf("update_file: GenerateEmbedding chunk %d error: %v", i, err))
			continue
		}
		for _, tag := range summary.DomainTags {
			allTags[tag] = struct{}{}
		}
		embeddings = append(embeddings, models.ChunkEmbedding{
			ID:           uuid.New(),
			RepoID:       repoID,
			RepoFilePath: payload.FilePath,
			ChunkIndex:   i,
			Summary:      summary.Summary,
			Tags:         summary.DomainTags,
			Embedding:    pgvector.NewVector(embedding),
			CreatedAt:    time.Now(),
		})
	}

	if len(embeddings) > 0 {
		if err := j.chunkEmbeddingRepo.Upsert(ctx, embeddings); err != nil {
			return fmt.Errorf("update_file: upsert embeddings: %w", err)
		}
	}

	mergedTags := sortedKeys(allTags)
	now := time.Now()
	if err := j.fileIndexRepo.UpdateMetadata(ctx, repoID, payload.FilePath, language, mergedTags, now); err != nil {
		return fmt.Errorf("update_file: update repository_files metadata: %w", err)
	}

	// Update the git_sha in repository_files so future identical pushes are skipped
	if fi != nil {
		fi.GitSHA = payload.SHA
		fi.ProcessedAt = &now
		if err := j.fileIndexRepo.Upsert(ctx, fi); err != nil {
			logger.Info(ctx, "update_file: update git_sha error: "+err.Error())
		}
	}

	logger.Info(ctx, fmt.Sprintf("update_file: done filePath=%s chunks=%d", payload.FilePath, len(embeddings)))
	return nil
}
