package jobs

import (
	ai_client "abstractions/clients/ai"
	github_client "abstractions/clients/github"
	"abstractions/config"
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/services/chunker"
	"abstractions/services/indexer"
	integration_service "abstractions/services/integrations"
	github_integration "abstractions/services/integrations/core/github"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
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
func (noopJobEnqueuer) EnqueueDeleteFile(_ uuid.UUID, _ string) error       { return nil }

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
	chunk      chunker.Chunk
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
			return j.handleProcessError(ctx, repoID, payload.FilePath, fmt.Errorf("summarize_file: %w", err))
		}
	}

	// Resolve the GitHub access token for the workspace owner
	integration, err := j.integrationService.GetIntegrationByType(ctx, repo.WorkspaceID, github_integration.IntegrationTypeGitHub)
	if err != nil {
		return j.handleProcessError(ctx, repoID, payload.FilePath, fmt.Errorf("summarize_file: github integration not found: %w", err))
	}
	accessToken, err := j.integrationService.GetDecryptedToken(ctx, integration.ID)
	if err != nil {
		return j.handleProcessError(ctx, repoID, payload.FilePath, fmt.Errorf("summarize_file: decrypt token: %w", err))
	}

	// Fetch file content — held in memory only; never written to DB or logs
	owner, repoName, err := splitOwnerRepo(repo.Name)
	if err != nil {
		return j.handleProcessError(ctx, repoID, payload.FilePath, fmt.Errorf("summarize_file: %w", err))
	}
	logger.Debug(ctx, fmt.Sprintf("summarize_file: fetching file content owner=%s repo=%s ref=%s", owner, repoName, payload.GitSHA))
	rawContent, fileMeta, err := j.fileFetcher.GetFileContent(ctx, accessToken, owner, repoName, payload.FilePath, payload.GitSHA)
	if err != nil {
		return j.handleProcessError(ctx, repoID, payload.FilePath, fmt.Errorf("summarize_file: fetch content: %w", err))
	}

	language := chunker.DetectLanguage(payload.FilePath)
	fileSize := 0
	if fileMeta != nil {
		fileSize = fileMeta.Size
	}
	logger.Debug(ctx, fmt.Sprintf("summarize_file: file fetched language=%s fileSize=%d", language, fileSize))

	// Skip binary / generated / minified content without burning AI tokens.
	if indexer.IsBinaryContent(rawContent) {
		logger.Info(ctx, "summarize_file: skipping binary content")
		return j.markProcessedEmpty(ctx, repoID, payload.FilePath, language, "binary content")
	}
	contentStr := string(rawContent)
	rawContent = nil
	if indexer.IsGeneratedContent(contentStr) {
		logger.Info(ctx, "summarize_file: skipping generated content")
		return j.markProcessedEmpty(ctx, repoID, payload.FilePath, language, "generated content")
	}
	if indexer.IsMinifiedOrDense(contentStr) {
		logger.Info(ctx, "summarize_file: skipping minified/dense content")
		return j.markProcessedEmpty(ctx, repoID, payload.FilePath, language, "minified content")
	}

	chunks := chunker.ChunkFile(payload.FilePath, contentStr, language)
	contentStr = "" // drop source as soon as chunked

	logger.Debug(ctx, fmt.Sprintf("summarize_file: chunked chunkCount=%d language=%s", len(chunks), language))

	if len(chunks) == 0 {
		return j.markProcessedEmpty(ctx, repoID, payload.FilePath, language, "")
	}

	// Fetch the repository_files row to populate repo_file_id on each embedding.
	repoFile, err := j.fileIndexRepo.GetByRepoAndPath(ctx, repoID, payload.FilePath)
	if err != nil {
		return j.handleProcessError(ctx, repoID, payload.FilePath, fmt.Errorf("summarize_file: lookup repository_files: %w", err))
	}
	var repoFileID *uuid.UUID
	if repoFile != nil {
		repoFileID = &repoFile.ID
	}

	// Always clear prior chunks before insert so shrunk files leave no stale high indices.
	if err := j.chunkEmbeddingRepo.DeleteByFile(ctx, repoID, payload.FilePath); err != nil {
		return j.handleProcessError(ctx, repoID, payload.FilePath, fmt.Errorf("summarize_file: clear chunks: %w", err))
	}

	summaryResults, err := summarizeFileChunks(ctx, aiProvider, payload.FilePath, language, chunks)
	if err != nil {
		// Rate limits propagate for asynq retry; other errors already skipped per-chunk.
		if ai_client.IsRateLimitError(err) {
			logger.Error(ctx, err)
			return fmt.Errorf("summarize_file: %w", err)
		}
		return j.handleProcessError(ctx, repoID, payload.FilePath, fmt.Errorf("summarize_file: %w", err))
	}
	textsToEmbed := embeddingInputs(payload.FilePath, language, summaryResults)
	embeddingVecs, err := generateEmbeddings(ctx, aiProvider, textsToEmbed)
	if err != nil {
		if ai_client.IsRateLimitError(err) {
			logger.Error(ctx, err)
			return fmt.Errorf("summarize_file: %w", err)
		}
		return j.handleProcessError(ctx, repoID, payload.FilePath, fmt.Errorf("summarize_file: %w", err))
	}
	embeddings := buildChunkEmbeddings(ctx, repoID, repoFileID, payload.FilePath, language, summaryResults, embeddingVecs)

	if len(embeddings) > 0 {
		if err := j.chunkEmbeddingRepo.Upsert(ctx, embeddings); err != nil {
			return j.handleProcessError(ctx, repoID, payload.FilePath, fmt.Errorf("summarize_file: upsert embeddings: %w", err))
		}
	}

	mergedTags := sortedKeys(summaryResults.tags)
	now := time.Now()
	logger.Debug(ctx, fmt.Sprintf("summarize_file: updating metadata language=%s tagCount=%d tags=%v", language, len(mergedTags), []string(mergedTags)))
	if err := j.fileIndexRepo.UpdateMetadataFull(ctx, repoID, payload.FilePath, language, mergedTags, now, "", 0, indexer.Version); err != nil {
		return j.handleProcessError(ctx, repoID, payload.FilePath, fmt.Errorf("summarize_file: update file_index: %w", err))
	}

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

func (j *SummarizeFileJob) markProcessedEmpty(ctx context.Context, repoID uuid.UUID, filePath, language, note string) error {
	now := time.Now()
	if err := j.fileIndexRepo.UpdateMetadataFull(ctx, repoID, filePath, language, models.StringArray{}, now, note, 0, indexer.Version); err != nil {
		logger.Error(ctx, err)
		return fmt.Errorf("summarize_file: update file_index (empty): %w", err)
	}
	if err := j.finalizeRepositoryIndexing(ctx, repoID); err != nil {
		return err
	}
	logger.Info(ctx, "summarize_file: marked processed with no chunks")
	return nil
}

// handleProcessError increments attempts and either re-queues (return err) or permanently fails the file.
func (j *SummarizeFileJob) handleProcessError(ctx context.Context, repoID uuid.UUID, filePath string, err error) error {
	logger.Error(ctx, err)
	if ai_client.IsRateLimitError(err) {
		return err
	}
	attempts, incErr := j.fileIndexRepo.IncrementProcessAttempts(ctx, repoID, filePath)
	if incErr != nil {
		logger.Error(ctx, incErr)
	}
	if attempts >= indexer.MaxProcessAttempts {
		msg := err.Error()
		if markErr := j.fileIndexRepo.MarkFailed(ctx, repoID, filePath, msg, attempts); markErr != nil {
			logger.Error(ctx, markErr)
		}
		logger.Info(ctx, fmt.Sprintf("summarize_file: permanently failed after attempts=%d", attempts))
		_ = j.finalizeRepositoryIndexing(ctx, repoID)
		return nil // do not retry further
	}
	return err
}

// summarizeFileChunks summarizes each chunk (with bounded parallelism) and then the whole file.
func summarizeFileChunks(ctx context.Context, aiProvider ai_client.AIProvider, filePath, language string, chunks []chunker.Chunk) (fileSummaryResults, error) {
	results := fileSummaryResults{
		tags: make(map[string]struct{}),
	}

	type workItem struct {
		index int
		chunk chunker.Chunk
	}
	var work []workItem
	for i, chunk := range chunks {
		if chunker.IsImportChunk(chunk, language) {
			logger.Debug(ctx, fmt.Sprintf("summarize_file: skipping import chunk chunkIndex=%d", i))
			continue
		}
		work = append(work, workItem{index: i, chunk: chunk})
	}

	var (
		mu       sync.Mutex
		rateErr  error
		wg       sync.WaitGroup
		sem      = make(chan struct{}, indexer.MaxConcurrentChunkSummaries)
		partials = make([]summarizedChunk, 0, len(work))
	)

	for _, w := range work {
		if rateErr != nil {
			break
		}
		wg.Add(1)
		sem <- struct{}{}
		go func(w workItem) {
			defer wg.Done()
			defer func() { <-sem }()

			logger.Debug(ctx, fmt.Sprintf("summarize_file: processing chunk chunkIndex=%d chunkType=%s startLine=%d endLine=%d", w.index, w.chunk.Type, w.chunk.StartLine, w.chunk.EndLine))
			content := chunker.ContentForSummary(w.chunk)
			summary, err := aiProvider.SummarizeCode(ctx, filePath, language, content, w.chunk.StartLine, w.chunk.EndLine, w.chunk.Type)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				if ai_client.IsRateLimitError(err) {
					rateErr = fmt.Errorf("SummarizeCode chunkIndex=%d: %w", w.index, err)
					return
				}
				logger.Info(ctx, fmt.Sprintf("summarize_file: SummarizeCode error chunkIndex=%d: %v", w.index, err))
				return
			}
			if summary == nil {
				logger.Debug(ctx, fmt.Sprintf("summarize_file: skipping chunk (no logic) chunkIndex=%d", w.index))
				return
			}
			logger.Debug(ctx, fmt.Sprintf("summarize_file: chunk summarized chunkIndex=%d symbolName=%s symbolType=%s tagCount=%d", w.index, summary.SymbolName, summary.SymbolType, len(summary.DomainTags)))
			partials = append(partials, summarizedChunk{chunkIndex: w.index, chunk: w.chunk, result: summary})
		}(w)
	}
	wg.Wait()
	if rateErr != nil {
		return results, rateErr
	}

	// Stable order by original chunk index for file-level summary.
	sort.Slice(partials, func(i, j int) bool { return partials[i].chunkIndex < partials[j].chunkIndex })
	chunkSummaryTexts := make([]string, 0, len(partials))
	for _, p := range partials {
		addTags(results.tags, p.result.DomainTags)
		chunkSummaryTexts = append(chunkSummaryTexts, p.result.Summary)
	}
	results.chunks = partials

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

func embeddingInputs(filePath, language string, results fileSummaryResults) []string {
	textsToEmbed := make([]string, 0, len(results.chunks)+1)
	for _, chunk := range results.chunks {
		r := chunk.result
		textsToEmbed = append(textsToEmbed, indexer.BuildEmbedText(filePath, language, r.SymbolName, r.SymbolType, r.Summary, r.DomainTags))
	}
	if results.fileSummary != nil {
		r := results.fileSummary
		textsToEmbed = append(textsToEmbed, indexer.BuildEmbedText(filePath, language, r.SymbolName, r.SymbolType, r.Summary, r.DomainTags))
	}
	return textsToEmbed
}

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

func buildChunkEmbeddings(ctx context.Context, repoID uuid.UUID, repoFileID *uuid.UUID, filePath, language string, results fileSummaryResults, embeddingVecs [][]float32) []models.ChunkEmbedding {
	embeddings := make([]models.ChunkEmbedding, 0, len(results.chunks)+1)
	for i, chunk := range results.chunks {
		if embeddingVecs == nil || i >= len(embeddingVecs) || embeddingVecs[i] == nil {
			logger.Debug(ctx, fmt.Sprintf("summarize_file: no embedding for chunk chunkIndex=%d", chunk.chunkIndex))
			continue
		}
		logger.Debug(ctx, fmt.Sprintf("summarize_file: embedding stored chunkIndex=%d", chunk.chunkIndex))
		embeddings = append(embeddings, newChunkEmbedding(repoID, repoFileID, filePath, language, chunk.chunkIndex, chunk.chunk, chunk.result, embeddingVecs[i]))
	}

	if results.fileSummary == nil {
		return embeddings
	}

	fileEmbedIdx := len(results.chunks)
	if embeddingVecs == nil || fileEmbedIdx >= len(embeddingVecs) || embeddingVecs[fileEmbedIdx] == nil {
		return embeddings
	}

	logger.Debug(ctx, fmt.Sprintf("summarize_file: file-level embedding stored symbolType=%s tagCount=%d", results.fileSummary.SymbolType, len(results.fileSummary.DomainTags)))
	fileChunk := chunker.Chunk{Type: "file", StartLine: 0, EndLine: 0}
	return append(embeddings, newChunkEmbedding(repoID, repoFileID, filePath, language, -1, fileChunk, results.fileSummary, embeddingVecs[fileEmbedIdx]))
}

func newChunkEmbedding(repoID uuid.UUID, repoFileID *uuid.UUID, filePath, language string, chunkIndex int, chunk chunker.Chunk, summary *ai_client.SummarizeResult, embedding []float32) models.ChunkEmbedding {
	searchDoc := indexer.BuildSearchDocument(filePath, language, summary.SymbolName, summary.SymbolType, summary.Summary, summary.DomainTags)
	return models.ChunkEmbedding{
		ID:             uuid.New(),
		RepoID:         repoID,
		RepoFileID:     repoFileID,
		RepoFilePath:   filePath,
		ChunkIndex:     chunkIndex,
		Summary:        summary.Summary,
		Tags:           summary.DomainTags,
		SymbolName:     summary.SymbolName,
		SymbolType:     summary.SymbolType,
		StartLine:      chunk.StartLine,
		EndLine:        chunk.EndLine,
		ChunkType:      chunk.Type,
		Language:       language,
		SearchDocument: searchDoc,
		IndexerVersion: indexer.Version,
		Embedding:      pgvector.NewVector(embedding),
		CreatedAt:      time.Now(),
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
