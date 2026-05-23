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

// noopJobEnqueuer is a no-op implementation of JobEnqueuer used when no enqueuer is needed
// (e.g. in unit tests that do not verify downstream enqueueing).
type noopJobEnqueuer struct{}

func (noopJobEnqueuer) EnqueueSummarizeFile(_ uuid.UUID, _, _ string) error { return nil }
func (noopJobEnqueuer) EnqueueGenerateTopics(_ uuid.UUID) error             { return nil }
func (noopJobEnqueuer) EnqueueGenerateNewsletter(_, _ uuid.UUID) error      { return nil }

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
	enqueuer           JobEnqueuer
	integrationService *integration_service.Service
	fileIndexRepo      *repositories.RepositoryFileRepository
	chunkEmbeddingRepo *repositories.ChunkEmbeddingRepository
	repositoryRepo     *repositories.RepositoryRepository
}

func NewSummarizeFileJob(db *gorm.DB, cfg *config.Config, enqueuer JobEnqueuer) *SummarizeFileJob {
	return &SummarizeFileJob{
		db:                 db,
		config:             cfg,
		aiProvider:         nil, // resolved per-execution from workspace configuration
		fileFetcher:        &defaultFileContentFetcher{},
		enqueuer:           enqueuer,
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
	enqueuer JobEnqueuer,
) *SummarizeFileJob {
	if enqueuer == nil {
		enqueuer = noopJobEnqueuer{}
	}
	return &SummarizeFileJob{
		db:                 db,
		config:             cfg,
		aiProvider:         aiProvider,
		fileFetcher:        fileFetcher,
		enqueuer:           enqueuer,
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

	ctx = logger.SetMetadata(ctx, map[string]interface{}{
		"workspaceId": repo.WorkspaceID.String(),
		"repoName":    repo.Name,
	})

	aiProvider := j.aiProvider
	if aiProvider == nil {
		aiProvider, err = newWorkspaceAIProvider(ctx, j.db, repo.WorkspaceID)
		if err != nil {
			logger.Error(ctx, err)
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
		// Files that produce no chunks are still considered processed.
		now := time.Now()
		if err := j.fileIndexRepo.UpdateMetadata(ctx, repoID, payload.FilePath, language, models.StringArray{}, now); err != nil {
			logger.Error(ctx, err)
			return fmt.Errorf("summarize_file: update file_index (no chunks): %w", err)
		}
		if err := j.finalizeRepositoryIndexing(ctx, repoID); err != nil {
			return err
		}
		logger.Info(ctx, "summarize_file: no chunks produced, marked as processed")
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

	// Phase 1: Summarise each chunk; collect results for batch embedding later.
	type chunkRecord struct {
		chunkIndex int
		result     *ai_client.SummarizeResult
	}
	allTags := map[string]struct{}{}
	var chunkRecords []chunkRecord
	var chunkSummaryTexts []string // used as input to SummarizeFile

	for i, chunk := range chunks {
		logger.Debug(ctx, fmt.Sprintf("summarize_file: processing chunk chunkIndex=%d chunkType=%s startLine=%d endLine=%d", i, chunk.Type, chunk.StartLine, chunk.EndLine))

		if chunker.IsImportChunk(chunk, language) {
			logger.Debug(ctx, fmt.Sprintf("summarize_file: skipping import chunk chunkIndex=%d", i))
			continue
		}

		summary, err := aiProvider.SummarizeCode(ctx, payload.FilePath, language, chunk.Content, chunk.StartLine, chunk.EndLine, chunk.Type)
		if err != nil {
			logger.Info(ctx, fmt.Sprintf("summarize_file: SummarizeCode error chunkIndex=%d: %v", i, err))
			continue
		}
		if summary == nil {
			logger.Debug(ctx, fmt.Sprintf("summarize_file: skipping chunk (no logic) chunkIndex=%d", i))
			continue
		}

		logger.Debug(ctx, fmt.Sprintf("summarize_file: chunk summarized chunkIndex=%d symbolName=%s symbolType=%s tagCount=%d", i, summary.SymbolName, summary.SymbolType, len(summary.DomainTags)))
		for _, tag := range summary.DomainTags {
			allTags[tag] = struct{}{}
		}
		chunkRecords = append(chunkRecords, chunkRecord{chunkIndex: i, result: summary})
		chunkSummaryTexts = append(chunkSummaryTexts, summary.Summary)
	}

	// Phase 2: File-level summary built from per-chunk summaries (no raw code).
	var fileSummaryResult *ai_client.SummarizeResult
	if len(chunkSummaryTexts) > 0 {
		logger.Debug(ctx, fmt.Sprintf("summarize_file: generating file-level summary chunkSummaryCount=%d", len(chunkSummaryTexts)))
		fileSummaryResult, err = aiProvider.SummarizeFile(ctx, payload.FilePath, language, chunkSummaryTexts)
		if err != nil {
			logger.Info(ctx, fmt.Sprintf("summarize_file: SummarizeFile error: %v", err))
			fileSummaryResult = nil
		}
		if fileSummaryResult != nil {
			for _, tag := range fileSummaryResult.DomainTags {
				allTags[tag] = struct{}{}
			}
		}
	}

	// Phase 3: Batch-embed all summaries in one API call to minimise round-trips.
	textsToEmbed := make([]string, 0, len(chunkRecords)+1)
	for _, r := range chunkRecords {
		textsToEmbed = append(textsToEmbed, r.result.Summary)
	}
	if fileSummaryResult != nil {
		textsToEmbed = append(textsToEmbed, fileSummaryResult.Summary)
	}

	var embeddingVecs [][]float32
	if len(textsToEmbed) > 0 {
		embeddingVecs, err = aiProvider.GenerateEmbeddings(ctx, textsToEmbed)
		if err != nil {
			logger.Info(ctx, fmt.Sprintf("summarize_file: GenerateEmbeddings error: %v", err))
			embeddingVecs = nil
		}
	}

	// Phase 4: Build ChunkEmbedding rows from summarisation results + vectors.
	embeddings := make([]models.ChunkEmbedding, 0, len(textsToEmbed))
	for j, rec := range chunkRecords {
		if embeddingVecs == nil || j >= len(embeddingVecs) || embeddingVecs[j] == nil {
			logger.Debug(ctx, fmt.Sprintf("summarize_file: no embedding for chunk chunkIndex=%d", rec.chunkIndex))
			continue
		}
		logger.Debug(ctx, fmt.Sprintf("summarize_file: embedding stored chunkIndex=%d", rec.chunkIndex))
		embeddings = append(embeddings, models.ChunkEmbedding{
			ID:           uuid.New(),
			RepoID:       repoID,
			RepoFileID:   repoFileID,
			RepoFilePath: payload.FilePath,
			ChunkIndex:   rec.chunkIndex,
			Summary:      rec.result.Summary,
			Tags:         rec.result.DomainTags,
			SymbolName:   rec.result.SymbolName,
			SymbolType:   rec.result.SymbolType,
			Embedding:    pgvector.NewVector(embeddingVecs[j]),
			CreatedAt:    time.Now(),
		})
	}
	if fileSummaryResult != nil {
		fileEmbedIdx := len(chunkRecords)
		if embeddingVecs != nil && fileEmbedIdx < len(embeddingVecs) && embeddingVecs[fileEmbedIdx] != nil {
			logger.Debug(ctx, fmt.Sprintf("summarize_file: file-level embedding stored symbolType=%s tagCount=%d", fileSummaryResult.SymbolType, len(fileSummaryResult.DomainTags)))
			embeddings = append(embeddings, models.ChunkEmbedding{
				ID:           uuid.New(),
				RepoID:       repoID,
				RepoFileID:   repoFileID,
				RepoFilePath: payload.FilePath,
				ChunkIndex:   -1,
				Summary:      fileSummaryResult.Summary,
				Tags:         fileSummaryResult.DomainTags,
				SymbolName:   fileSummaryResult.SymbolName,
				SymbolType:   fileSummaryResult.SymbolType,
				Embedding:    pgvector.NewVector(embeddingVecs[fileEmbedIdx]),
				CreatedAt:    time.Now(),
			})
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

	if err := j.finalizeRepositoryIndexing(ctx, repoID); err != nil {
		return err
	}

	logger.Info(ctx, fmt.Sprintf("summarize_file: done language=%s chunkCount=%d embeddingCount=%d tagCount=%d", language, len(chunks), len(embeddings), len(mergedTags)))
	return nil
}

func (j *SummarizeFileJob) finalizeRepositoryIndexing(ctx context.Context, repoID uuid.UUID) error {
	transitioned, err := j.repositoryRepo.MarkIndexedIfReady(ctx, repoID, time.Now())
	if err != nil {
		return fmt.Errorf("summarize_file: finalize repository status: %w", err)
	}
	if !transitioned {
		return nil
	}

	logger.Debug(ctx, "summarize_file: status indexing → indexed")
	if err := j.enqueuer.EnqueueGenerateTopics(repoID); err != nil {
		logger.Error(ctx, fmt.Errorf("summarize_file: enqueue generate_topics: %w", err))
	} else {
		logger.Debug(ctx, "summarize_file: enqueued generate_topics")
	}
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
