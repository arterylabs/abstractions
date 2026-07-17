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

type summarizedChunk struct {
	chunkIndex int
	result     *ai_client.SummarizeResult
}

type fileSummaryResults struct {
	chunks      []summarizedChunk
	fileSummary *ai_client.SummarizeResult
	tags        map[string]struct{}
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
		integrationService: integration_service.NewService(db),
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
		integrationService: integration_service.NewService(db),
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

	summaryResults, err := summarizeFileChunks(ctx, aiProvider, payload.FilePath, language, chunks)
	if err != nil {
		// A provider rate limit was hit and retries were exhausted in-transport. Propagate so
		// asynq re-queues the whole file rather than silently dropping the unsummarized chunks.
		logger.Error(ctx, err)
		return fmt.Errorf("summarize_file: %w", err)
	}
	textsToEmbed := embeddingInputs(summaryResults)
	embeddingVecs, err := generateEmbeddings(ctx, aiProvider, textsToEmbed)
	if err != nil {
		logger.Error(ctx, err)
		return fmt.Errorf("summarize_file: %w", err)
	}
	embeddings := buildChunkEmbeddings(ctx, repoID, repoFileID, payload.FilePath, summaryResults, embeddingVecs)

	// Upsert chunk embeddings in bulk
	if len(embeddings) > 0 {
		if err := j.chunkEmbeddingRepo.Upsert(ctx, embeddings); err != nil {
			logger.Error(ctx, err)
			return fmt.Errorf("summarize_file: upsert embeddings: %w", err)
		}
	}

	// Merge collected domain tags and update file_index
	mergedTags := sortedKeys(summaryResults.tags)
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

// summarizeFileChunks summarizes each chunk and then the whole file. A chunk that fails for an
// ordinary reason (bad JSON, transient provider error, etc.) is skipped so one bad chunk doesn't
// block the rest of the file. A rate-limit error, however, means every subsequent call would fail
// the same way — so it's returned immediately rather than skipped, letting the caller propagate it
// so asynq re-queues the whole job instead of silently under-indexing the file.
func summarizeFileChunks(ctx context.Context, aiProvider ai_client.AIProvider, filePath, language string, chunks []chunker.Chunk) (fileSummaryResults, error) {
	results := fileSummaryResults{
		tags: make(map[string]struct{}),
	}

	chunkSummaryTexts := make([]string, 0, len(chunks))
	for i, chunk := range chunks {
		logger.Debug(ctx, fmt.Sprintf("summarize_file: processing chunk chunkIndex=%d chunkType=%s startLine=%d endLine=%d", i, chunk.Type, chunk.StartLine, chunk.EndLine))

		if chunker.IsImportChunk(chunk, language) {
			logger.Debug(ctx, fmt.Sprintf("summarize_file: skipping import chunk chunkIndex=%d", i))
			continue
		}

		summary, err := aiProvider.SummarizeCode(ctx, filePath, language, chunk.Content, chunk.StartLine, chunk.EndLine, chunk.Type)
		if err != nil {
			if ai_client.IsRateLimitError(err) {
				return results, fmt.Errorf("SummarizeCode chunkIndex=%d: %w", i, err)
			}
			logger.Info(ctx, fmt.Sprintf("summarize_file: SummarizeCode error chunkIndex=%d: %v", i, err))
			continue
		}
		if summary == nil {
			logger.Debug(ctx, fmt.Sprintf("summarize_file: skipping chunk (no logic) chunkIndex=%d", i))
			continue
		}

		logger.Debug(ctx, fmt.Sprintf("summarize_file: chunk summarized chunkIndex=%d symbolName=%s symbolType=%s tagCount=%d", i, summary.SymbolName, summary.SymbolType, len(summary.DomainTags)))
		addTags(results.tags, summary.DomainTags)
		results.chunks = append(results.chunks, summarizedChunk{chunkIndex: i, result: summary})
		chunkSummaryTexts = append(chunkSummaryTexts, summary.Summary)
	}

	if len(chunkSummaryTexts) == 0 {
		return results, nil
	}

	logger.Debug(ctx, fmt.Sprintf("summarize_file: generating file-level summary chunkSummaryCount=%d", len(chunkSummaryTexts)))
	fileSummary, err := aiProvider.SummarizeFile(ctx, filePath, language, chunkSummaryTexts)
	if err != nil {
		if ai_client.IsRateLimitError(err) {
			return results, fmt.Errorf("SummarizeFile: %w", err)
		}
		logger.Info(ctx, fmt.Sprintf("summarize_file: SummarizeFile error: %v", err))
		return results, nil
	}
	if fileSummary != nil {
		results.fileSummary = fileSummary
		addTags(results.tags, fileSummary.DomainTags)
	}

	return results, nil
}

func embeddingInputs(results fileSummaryResults) []string {
	textsToEmbed := make([]string, 0, len(results.chunks)+1)
	for _, chunk := range results.chunks {
		textsToEmbed = append(textsToEmbed, chunk.result.Summary)
	}
	if results.fileSummary != nil {
		textsToEmbed = append(textsToEmbed, results.fileSummary.Summary)
	}
	return textsToEmbed
}

// generateEmbeddings returns embeddings for textsToEmbed. A rate-limit error is propagated (see
// summarizeFileChunks) so the job is re-queued instead of persisting chunk summaries with no
// embeddings; any other error is logged and swallowed, matching prior best-effort behaviour.
func generateEmbeddings(ctx context.Context, aiProvider ai_client.AIProvider, textsToEmbed []string) ([][]float32, error) {
	if len(textsToEmbed) == 0 {
		return nil, nil
	}

	embeddingVecs, err := aiProvider.GenerateEmbeddings(ctx, textsToEmbed)
	if err != nil {
		if ai_client.IsRateLimitError(err) {
			return nil, fmt.Errorf("GenerateEmbeddings: %w", err)
		}
		logger.Info(ctx, fmt.Sprintf("summarize_file: GenerateEmbeddings error: %v", err))
		return nil, nil
	}
	return embeddingVecs, nil
}

func buildChunkEmbeddings(ctx context.Context, repoID uuid.UUID, repoFileID *uuid.UUID, filePath string, results fileSummaryResults, embeddingVecs [][]float32) []models.ChunkEmbedding {
	embeddings := make([]models.ChunkEmbedding, 0, len(results.chunks)+1)
	for i, chunk := range results.chunks {
		if embeddingVecs == nil || i >= len(embeddingVecs) || embeddingVecs[i] == nil {
			logger.Debug(ctx, fmt.Sprintf("summarize_file: no embedding for chunk chunkIndex=%d", chunk.chunkIndex))
			continue
		}
		logger.Debug(ctx, fmt.Sprintf("summarize_file: embedding stored chunkIndex=%d", chunk.chunkIndex))
		embeddings = append(embeddings, newChunkEmbedding(repoID, repoFileID, filePath, chunk.chunkIndex, chunk.result, embeddingVecs[i]))
	}

	if results.fileSummary == nil {
		return embeddings
	}

	fileEmbedIdx := len(results.chunks)
	if embeddingVecs == nil || fileEmbedIdx >= len(embeddingVecs) || embeddingVecs[fileEmbedIdx] == nil {
		return embeddings
	}

	logger.Debug(ctx, fmt.Sprintf("summarize_file: file-level embedding stored symbolType=%s tagCount=%d", results.fileSummary.SymbolType, len(results.fileSummary.DomainTags)))
	return append(embeddings, newChunkEmbedding(repoID, repoFileID, filePath, -1, results.fileSummary, embeddingVecs[fileEmbedIdx]))
}

func newChunkEmbedding(repoID uuid.UUID, repoFileID *uuid.UUID, filePath string, chunkIndex int, summary *ai_client.SummarizeResult, embedding []float32) models.ChunkEmbedding {
	return models.ChunkEmbedding{
		ID:           uuid.New(),
		RepoID:       repoID,
		RepoFileID:   repoFileID,
		RepoFilePath: filePath,
		ChunkIndex:   chunkIndex,
		Summary:      summary.Summary,
		Tags:         summary.DomainTags,
		SymbolName:   summary.SymbolName,
		SymbolType:   summary.SymbolType,
		Embedding:    pgvector.NewVector(embedding),
		CreatedAt:    time.Now(),
	}
}

func addTags(tags map[string]struct{}, newTags []string) {
	for _, tag := range newTags {
		tags[tag] = struct{}{}
	}
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
