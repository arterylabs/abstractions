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
	GetFileContent(ctx context.Context, accessToken, owner, repo, path, ref string) ([]byte, *github_client.GitHubFileMetadata, error)
}

// defaultFileContentFetcher delegates to the real github_client.GetFileContent package function.
type defaultFileContentFetcher struct{}

func (d *defaultFileContentFetcher) GetFileContent(ctx context.Context, accessToken, owner, repo, path, ref string) ([]byte, *github_client.GitHubFileMetadata, error) {
	return github_client.GetFileContent(ctx, accessToken, owner, repo, path, ref)
}

// SummarizeFilePayload is the JSON payload for a SummarizeFileJob task.
type SummarizeFilePayload struct {
	RepositoryID string `json:"repositoryId"`
	FilePath     string `json:"filePath"`
	GitSHA       string `json:"gitSha"`
}

// SummarizeFileJob processes a single file: chunks it, summarises each chunk via AI,
// generates embeddings, upserts chunk_embeddings, and updates file_index metadata.
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
		aiProvider:         nil, // resolved per-execution from workspace configuration
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
		logger.Error(ctx, err)
		return fmt.Errorf("summarize_file: decode payload: %w", err)
	}

	repoID, err := uuid.Parse(payload.RepositoryID)
	if err != nil {
		logger.Error(ctx, err)
		return fmt.Errorf("summarize_file: invalid repo id: %w", err)
	}

	ctx = logger.SetMetadata(ctx, map[string]interface{}{
		"repoId":   repoID.String(),
		"filePath": payload.FilePath,
	})

	logger.Info(ctx, "summarize_file: starting")

	repoCtx, err := LoadRepositoryContext(ctx, repoID, j.repositoryRepo)
	if err != nil {
		logger.Error(ctx, err)
		return fmt.Errorf("summarize_file: %w", err)
	}
	repo := repoCtx.Repository

	aiProvider := j.aiProvider
	if aiProvider == nil {
		aiProvider, err = newWorkspaceAIProvider(ctx, j.db, repo.WorkspaceID)
		if err != nil {
			return fmt.Errorf("summarize_file: %w", err)
		}
	}

	// Resolve the GitHub access token for the workspace owner
	integration, err := j.integrationService.GetIntegrationByType(ctx, repo.WorkspaceID, github_integration.IntegrationTypeGitHub)
	if err != nil {
		logger.Error(ctx, err)
		return fmt.Errorf("summarize_file: github integration not found: %w", err)
	}
	accessToken, err := j.integrationService.GetDecryptedToken(ctx, integration.ID)
	if err != nil {
		logger.Error(ctx, err)
		return fmt.Errorf("summarize_file: decrypt token: %w", err)
	}

	// Fetch file content — held in memory only; never written to DB or logs
	owner, repoName, err := splitOwnerRepo(repo.Name)
	if err != nil {
		logger.Error(ctx, err)
		return fmt.Errorf("summarize_file: %w", err)
	}
	logger.Debug(ctx, fmt.Sprintf("summarize_file: fetching file content owner=%s repo=%s ref=%s", owner, repoName, payload.GitSHA))
	rawContent, fileMeta, err := j.fileFetcher.GetFileContent(ctx, accessToken, owner, repoName, payload.FilePath, payload.GitSHA)
	if err != nil {
		logger.Error(ctx, err)
		return fmt.Errorf("summarize_file: fetch content: %w", err)
	}

	language := chunker.DetectLanguage(payload.FilePath)
	fileSize := 0
	if fileMeta != nil {
		fileSize = fileMeta.Size
	}
	logger.Debug(ctx, fmt.Sprintf("summarize_file: file fetched language=%s fileSize=%d", language, fileSize))

	chunks := chunker.ChunkFile(payload.FilePath, string(rawContent), language)

	// Zero out rawContent immediately after chunking — no code content persisted beyond this point
	rawContent = nil

	logger.Debug(ctx, fmt.Sprintf("summarize_file: chunked chunkCount=%d language=%s", len(chunks), language))

	if len(chunks) == 0 {
		logger.Info(ctx, "summarize_file: no chunks produced, skipping")
		return nil
	}

	// Fetch the repository_files row to populate repo_file_id on each embedding.
	repoFile, err := j.fileIndexRepo.GetByRepoAndPath(ctx, repoID, payload.FilePath)
	if err != nil {
		logger.Error(ctx, err)
		return fmt.Errorf("summarize_file: lookup repository_files: %w", err)
	}
	var repoFileID *uuid.UUID
	if repoFile != nil {
		repoFileID = &repoFile.ID
	}

	// Summarise and embed each chunk; accumulate domain tags and summaries
	embeddings := make([]models.ChunkEmbedding, 0, len(chunks))
	allTags := map[string]struct{}{}
	chunkSummaries := make([]string, 0, len(chunks))

	for i, chunk := range chunks {
		logger.Debug(ctx, fmt.Sprintf("summarize_file: processing chunk chunkIndex=%d chunkType=%s startLine=%d endLine=%d", i, chunk.Type, chunk.StartLine, chunk.EndLine))

		// Skip chunks that are purely import/package declarations — they carry no logic worth summarising.
		if chunker.IsImportChunk(chunk, language) {
			logger.Debug(ctx, fmt.Sprintf("summarize_file: skipping import chunk chunkIndex=%d", i))
			continue
		}

		summary, err := aiProvider.SummarizeCode(ctx, payload.FilePath, language, chunk.Content, chunk.StartLine, chunk.EndLine, chunk.Type)
		if err != nil {
			logger.Info(ctx, fmt.Sprintf("summarize_file: SummarizeCode error chunkIndex=%d: %v", i, err))
			continue
		}
		// Model signalled the chunk has no meaningful logic (e.g. import block detected by AI).
		if summary == nil {
			logger.Debug(ctx, fmt.Sprintf("summarize_file: skipping chunk (no logic) chunkIndex=%d", i))
			continue
		}

		logger.Debug(ctx, fmt.Sprintf("summarize_file: chunk summarized chunkIndex=%d symbolName=%s symbolType=%s tagCount=%d", i, summary.SymbolName, summary.SymbolType, len(summary.DomainTags)))

		embedding, err := aiProvider.GenerateEmbedding(ctx, summary.Summary)
		if err != nil {
			logger.Info(ctx, fmt.Sprintf("summarize_file: GenerateEmbedding error chunkIndex=%d: %v", i, err))
			continue
		}
		logger.Debug(ctx, fmt.Sprintf("summarize_file: embedding generated chunkIndex=%d", i))

		for _, tag := range summary.DomainTags {
			allTags[tag] = struct{}{}
		}

		chunkSummaries = append(chunkSummaries, summary.Summary)

		embeddings = append(embeddings, models.ChunkEmbedding{
			ID:           uuid.New(),
			RepoID:       repoID,
			RepoFileID:   repoFileID,
			RepoFilePath: payload.FilePath,
			ChunkIndex:   i,
			Summary:      summary.Summary,
			Tags:         summary.DomainTags,
			SymbolName:   summary.SymbolName,
			SymbolType:   summary.SymbolType,
			Embedding:    pgvector.NewVector(embedding),
			CreatedAt:    time.Now(),
		})
	}

	// Generate a file-level summary from the collected chunk summaries and store it at chunk_index = -1.
	if len(chunkSummaries) > 0 {
		logger.Debug(ctx, fmt.Sprintf("summarize_file: generating file-level summary chunkSummaryCount=%d", len(chunkSummaries)))
		fileSummary, err := aiProvider.SummarizeFile(ctx, payload.FilePath, language, chunkSummaries)
		if err != nil {
			logger.Info(ctx, fmt.Sprintf("summarize_file: SummarizeFile error: %v", err))
		} else {
			fileEmbedding, err := aiProvider.GenerateEmbedding(ctx, fileSummary.Summary)
			if err != nil {
				logger.Info(ctx, fmt.Sprintf("summarize_file: GenerateEmbedding file-level error: %v", err))
			} else {
				logger.Debug(ctx, fmt.Sprintf("summarize_file: file-level summary generated symbolType=%s tagCount=%d", fileSummary.SymbolType, len(fileSummary.DomainTags)))
				for _, tag := range fileSummary.DomainTags {
					allTags[tag] = struct{}{}
				}
				embeddings = append(embeddings, models.ChunkEmbedding{
					ID:           uuid.New(),
					RepoID:       repoID,
					RepoFileID:   repoFileID,
					RepoFilePath: payload.FilePath,
					ChunkIndex:   -1,
					Summary:      fileSummary.Summary,
					Tags:         fileSummary.DomainTags,
					SymbolName:   fileSummary.SymbolName,
					SymbolType:   fileSummary.SymbolType,
					Embedding:    pgvector.NewVector(fileEmbedding),
					CreatedAt:    time.Now(),
				})
			}
		}
	}

	// Upsert chunk embeddings in bulk
	if len(embeddings) > 0 {
		if err := j.chunkEmbeddingRepo.Upsert(ctx, embeddings); err != nil {
			logger.Error(ctx, err)
			return fmt.Errorf("summarize_file: upsert embeddings: %w", err)
		}
	}

	// Merge collected domain tags and update file_index
	mergedTags := sortedKeys(allTags)
	now := time.Now()
	logger.Debug(ctx, fmt.Sprintf("summarize_file: updating metadata language=%s tagCount=%d tags=%v", language, len(mergedTags), []string(mergedTags)))
	if err := j.fileIndexRepo.UpdateMetadata(ctx, repoID, payload.FilePath, language, mergedTags, now); err != nil {
		logger.Error(ctx, err)
		return fmt.Errorf("summarize_file: update file_index: %w", err)
	}

	// Persist GitHub Contents API metadata for this file (content bytes excluded).
	if fileMeta != nil {
		if metaBytes, jsonErr := json.Marshal(fileMeta); jsonErr == nil {
			if err := j.fileIndexRepo.UpdateServiceMetadata(ctx, repoID, payload.FilePath, models.JSONB(metaBytes)); err != nil {
				logger.Info(ctx, fmt.Sprintf("summarize_file: update service_metadata error: %v", err))
			}
		}
	}

	logger.Info(ctx, fmt.Sprintf("summarize_file: done language=%s chunkCount=%d embeddingCount=%d tagCount=%d", language, len(chunks), len(embeddings), len(mergedTags)))
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
