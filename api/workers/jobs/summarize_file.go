package jobs

import (
	ai_client "abstractions/clients/ai"
	github_client "abstractions/clients/github"
	"abstractions/config"
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/services/chunker"
	integration_service "abstractions/services/integrations"
	github_integration "abstractions/services/integrations/core/github"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"abstractions/logger"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	pgvector "github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

const SummarizeFileJobID = "summarize_file"

// FileContentFetcher is the interface for fetching raw file content from GitHub.
// It is satisfied by github_client.GetFileContent; the interface allows mocking in tests.
type FileContentFetcher interface {
	GetFileContent(ctx context.Context, accessToken, owner, repo, path, ref string) ([]byte, error)
}

// defaultFileContentFetcher delegates to the real github_client.GetFileContent package function.
type defaultFileContentFetcher struct{}

func (d *defaultFileContentFetcher) GetFileContent(ctx context.Context, accessToken, owner, repo, path, ref string) ([]byte, error) {
	return github_client.GetFileContent(ctx, accessToken, owner, repo, path, ref)
}

// SummarizeFilePayload is the JSON payload for a SummarizeFileJob task.
type SummarizeFilePayload struct {
	RepositoryID string `json:"repositoryId"`
	FilePath     string `json:"filePath"`
	GitSHA       string `json:"gitSha"`
}

// SummarizeFileJob processes a single file: chunks it, summarises each chunk via AI,
// generates embeddings, upserts chunk_embeddings, and updates repository_files metadata.
// Raw file content is never persisted in the database or written to logs.
type SummarizeFileJob struct {
	db                 *gorm.DB
	config             *config.Config
	aiProvider         ai_client.AIProvider
	fileFetcher        FileContentFetcher
	integrationService *integration_service.Service
	fileIndexRepo      *repositories.RepositoryFileRepository
	chunkEmbeddingRepo *repositories.ChunkEmbeddingRepository
	repositoryRepo     *repositories.RepositoryRepository
}

func NewSummarizeFileJob(db *gorm.DB, cfg *config.Config) *SummarizeFileJob {
	return &SummarizeFileJob{
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

// NewSummarizeFileJobWithDeps creates a job with injected dependencies for testing.
func NewSummarizeFileJobWithDeps(
	db *gorm.DB,
	cfg *config.Config,
	aiProvider ai_client.AIProvider,
	fileFetcher FileContentFetcher,
) *SummarizeFileJob {
	return &SummarizeFileJob{
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

func (j *SummarizeFileJob) ID() string {
	return SummarizeFileJobID
}

// NewSummarizeFileTask builds an asynq.Task for the given repository, file path, and git SHA.
func NewSummarizeFileTask(repositoryID uuid.UUID, filePath, gitSHA string) (*asynq.Task, error) {
	payload := SummarizeFilePayload{
		RepositoryID: repositoryID.String(),
		FilePath:     filePath,
		GitSHA:       gitSHA,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(SummarizeFileJobID, b), nil
}

// Handler processes a SummarizeFileJob task.
func (j *SummarizeFileJob) Handler(ctx context.Context, t *asynq.Task) error {
	var payload SummarizeFilePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("summarize_file: decode payload: %w", err)
	}

	repoID, err := uuid.Parse(payload.RepositoryID)
	if err != nil {
		return fmt.Errorf("summarize_file: invalid repo id: %w", err)
	}

	logger.Info(ctx, "summarize_file: starting filePath="+payload.FilePath)

	// Load the repository row to get workspace/owner details
	repo, err := j.repositoryRepo.GetByID(ctx, repoID)
	if err != nil || repo == nil {
		return fmt.Errorf("summarize_file: repo not found %s: %w", repoID, err)
	}

	// Resolve the GitHub access token for the workspace owner
	integration, err := j.integrationService.GetIntegrationByType(ctx, repo.WorkspaceID, github_integration.IntegrationTypeGitHub)
	if err != nil {
		return fmt.Errorf("summarize_file: github integration not found: %w", err)
	}
	accessToken, err := j.integrationService.GetDecryptedToken(ctx, integration.ID)
	if err != nil {
		return fmt.Errorf("summarize_file: decrypt token: %w", err)
	}

	// Fetch file content — held in memory only; never written to DB or logs
	owner, repoName, err := splitOwnerRepo(repo.Name)
	if err != nil {
		return fmt.Errorf("summarize_file: %w", err)
	}
	rawContent, err := j.fileFetcher.GetFileContent(ctx, accessToken, owner, repoName, payload.FilePath, payload.GitSHA)
	if err != nil {
		return fmt.Errorf("summarize_file: fetch content: %w", err)
	}

	language := chunker.DetectLanguage(payload.FilePath)
	chunks := chunker.ChunkFile(payload.FilePath, string(rawContent), language)

	// Zero out rawContent immediately after chunking — no code content persisted beyond this point
	rawContent = nil

	if len(chunks) == 0 {
		logger.Info(ctx, "summarize_file: no chunks produced, skipping filePath="+payload.FilePath)
		return nil
	}

	// Summarise and embed each chunk; accumulate domain tags
	embeddings := make([]models.ChunkEmbedding, 0, len(chunks))
	allTags := map[string]struct{}{}
	var chunkSummaries []string

	for i, chunk := range chunks {
		summary, err := j.aiProvider.SummarizeCode(ctx, payload.FilePath, language, chunk.Content, chunk.StartLine, chunk.EndLine, chunk.Type)
		if err != nil {
			logger.Info(ctx, fmt.Sprintf("summarize_file: SummarizeCode chunk %d error: %v", i, err))
			continue
		}
		if summary == nil {
			continue
		}

		embedding, err := j.aiProvider.GenerateEmbedding(ctx, summary.Summary)
		if err != nil {
			logger.Info(ctx, fmt.Sprintf("summarize_file: GenerateEmbedding chunk %d error: %v", i, err))
			continue
		}

		for _, tag := range summary.DomainTags {
			allTags[tag] = struct{}{}
		}

		chunkSummaries = append(chunkSummaries, summary.Summary)

		embeddings = append(embeddings, models.ChunkEmbedding{
			ID:           uuid.New(),
			RepoID:       repoID,
			RepoFilePath: payload.FilePath,
			ChunkIndex:   i,
			Summary:      summary.Summary,
			SymbolName:   summary.SymbolName,
			SymbolType:   summary.SymbolType,
			Tags:         summary.DomainTags,
			Embedding:    pgvector.NewVector(embedding),
			CreatedAt:    time.Now(),
		})
	}

	// Generate and upsert the file-level summary row (chunk_index = -1)
	if len(chunkSummaries) > 0 {
		fileSummary, err := j.aiProvider.SummarizeFile(ctx, payload.FilePath, language, strings.Join(chunkSummaries, "\n"))
		if err != nil {
			logger.Info(ctx, fmt.Sprintf("summarize_file: SummarizeFile error: %v", err))
		} else if fileSummary != nil {
			fileEmbedding, err := j.aiProvider.GenerateEmbedding(ctx, fileSummary.Summary)
			if err != nil {
				logger.Info(ctx, fmt.Sprintf("summarize_file: GenerateEmbedding file-level error: %v", err))
			} else {
				embeddings = append(embeddings, models.ChunkEmbedding{
					ID:           uuid.New(),
					RepoID:       repoID,
					RepoFilePath: payload.FilePath,
					ChunkIndex:   -1,
					Summary:      fileSummary.Summary,
					Tags:         fileSummary.DomainTags,
					Embedding:    pgvector.NewVector(fileEmbedding),
					CreatedAt:    time.Now(),
				})
			}
		}
	}

	// Upsert chunk embeddings in bulk
	if len(embeddings) > 0 {
		if err := j.chunkEmbeddingRepo.Upsert(ctx, embeddings); err != nil {
			return fmt.Errorf("summarize_file: upsert embeddings: %w", err)
		}
	}

	// Merge collected domain tags and update repository_files
	mergedTags := sortedKeys(allTags)
	now := time.Now()
	if err := j.fileIndexRepo.UpdateMetadata(ctx, repoID, payload.FilePath, language, mergedTags, now); err != nil {
		return fmt.Errorf("summarize_file: update repository_files: %w", err)
	}

	logger.Info(ctx, fmt.Sprintf("summarize_file: done filePath=%s chunks=%d tags=%d", payload.FilePath, len(embeddings), len(mergedTags)))
	return nil
}

// sortedKeys returns the keys of a string set as a sorted StringArray.
func sortedKeys(m map[string]struct{}) models.StringArray {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return models.StringArray(keys)
}

// splitOwnerRepo splits "owner/repo" into its two parts.
func splitOwnerRepo(fullName string) (string, string, error) {
	for i, c := range fullName {
		if c == '/' {
			return fullName[:i], fullName[i+1:], nil
		}
	}
	return "", "", fmt.Errorf("invalid full_name %q", fullName)
}
