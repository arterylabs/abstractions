package jobs

import (
	app_notifications "abstractions/clients/app_notifications"
	github_client "abstractions/clients/github"
	"abstractions/config"
	"abstractions/db/models"
	"abstractions/logger"
	"abstractions/repositories"
	"abstractions/services/chunker"
	"abstractions/services/indexer"
	integration_service "abstractions/services/integrations"
	github_integration "abstractions/services/integrations/core/github"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

const IndexRepositoryJobID = "index_repository"

// FileTreeFetcher is the interface for fetching a filtered repository file tree from GitHub.
// The second return value is the commit SHA at HEAD; callers must use this as the ref when
// fetching individual file contents — individual file SHAs are blob SHAs, not commit SHAs,
// and are not accepted by the GitHub Contents API.
type FileTreeFetcher interface {
	GetFilteredFileTree(ctx context.Context, accessToken, owner, repo, branch string) ([]github_client.FileTreeItem, string, error)
}

// defaultFileTreeFetcher delegates to the real github_client package function.
type defaultFileTreeFetcher struct{}

func (d *defaultFileTreeFetcher) GetFilteredFileTree(ctx context.Context, accessToken, owner, repo, branch string) ([]github_client.FileTreeItem, string, error) {
	return github_client.GetFilteredFileTree(ctx, accessToken, owner, repo, branch)
}

// JobEnqueuer is the interface used by job handlers that need to sub-enqueue tasks.
// It exposes only job-level operations, keeping asynq details inside the jobs package.
type JobEnqueuer interface {
	EnqueueSummarizeFile(repositoryID uuid.UUID, filePath, gitSHA string) error
	EnqueueGenerateTopics(repoID uuid.UUID) error
	// EnqueueGenerateNewsletter enqueues a newsletter generation job for repoID.
	// newsletterID must refer to an existing newsletter record with status=queued.
	EnqueueGenerateNewsletter(repoID, newsletterID uuid.UUID) error
}

// IndexRepositoryPayload is the JSON payload for an IndexRepositoryJob task.
type IndexRepositoryPayload struct {
	RepositoryID string `json:"repositoryId"`
	// Force reprocesses every file even when the git SHA is unchanged (e.g. indexer version bump).
	Force bool `json:"force,omitempty"`
}

// IndexRepositoryJob fetches the file tree, seeds repository_files, updates repo status, and
// enqueues one SummarizeFileJob per changed/new file.
type IndexRepositoryJob struct {
	db                 *gorm.DB
	config             *config.Config
	fileFetcher        FileTreeFetcher
	enqueuer           JobEnqueuer
	integrationService *integration_service.Service
	fileIndexRepo      *repositories.RepositoryFileRepository
	chunkEmbeddingRepo *repositories.ChunkEmbeddingRepository
	repositoryRepo     *repositories.RepositoryRepository
}

func NewIndexRepositoryJob(db *gorm.DB, cfg *config.Config, enqueuer JobEnqueuer) *IndexRepositoryJob {
	return &IndexRepositoryJob{
		db:                 db,
		config:             cfg,
		fileFetcher:        &defaultFileTreeFetcher{},
		enqueuer:           enqueuer,
		integrationService: integration_service.NewService(db),
		fileIndexRepo:      repositories.NewRepositoryFileRepository(db),
		chunkEmbeddingRepo: repositories.NewChunkEmbeddingRepository(db),
		repositoryRepo:     repositories.NewRepositoryRepository(db),
	}
}

// NewIndexRepositoryJobWithDeps creates a job with injected dependencies for testing.
func NewIndexRepositoryJobWithDeps(
	db *gorm.DB,
	cfg *config.Config,
	fileFetcher FileTreeFetcher,
	enqueuer JobEnqueuer,
) *IndexRepositoryJob {
	return &IndexRepositoryJob{
		db:                 db,
		config:             cfg,
		fileFetcher:        fileFetcher,
		enqueuer:           enqueuer,
		integrationService: integration_service.NewService(db),
		fileIndexRepo:      repositories.NewRepositoryFileRepository(db),
		chunkEmbeddingRepo: repositories.NewChunkEmbeddingRepository(db),
		repositoryRepo:     repositories.NewRepositoryRepository(db),
	}
}

func (j *IndexRepositoryJob) ID() string {
	return IndexRepositoryJobID
}

// NewIndexRepositoryTask builds an asynq.Task for the given repository ID.
func NewIndexRepositoryTask(repositoryID uuid.UUID) (*asynq.Task, error) {
	return NewIndexRepositoryTaskWithForce(repositoryID, false)
}

// NewIndexRepositoryTaskWithForce builds a task that optionally forces full reprocessing.
func NewIndexRepositoryTaskWithForce(repositoryID uuid.UUID, force bool) (*asynq.Task, error) {
	payload := IndexRepositoryPayload{RepositoryID: repositoryID.String(), Force: force}
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(IndexRepositoryJobID, b), nil
}

// Handler processes an IndexRepositoryJob task.
func (j *IndexRepositoryJob) Handler(ctx context.Context, t *asynq.Task) error {
	var payload IndexRepositoryPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("index_repository: decode payload: %w", err)
	}

	repoID, err := uuid.Parse(payload.RepositoryID)
	if err != nil {
		return fmt.Errorf("index_repository: invalid repo id: %w", err)
	}

	ctx = logger.SetMetadata(ctx, map[string]interface{}{
		"repoId": repoID.String(),
	})

	logger.Info(ctx, fmt.Sprintf("index_repository: starting force=%v indexerVersion=%s", payload.Force, indexer.Version))

	repoCtx, err := LoadRepositoryContext(ctx, repoID, j.repositoryRepo)
	if err != nil {
		return j.markError(ctx, repoID, fmt.Errorf("index_repository: %w", err))
	}
	repo := repoCtx.Repository

	ctx = logger.SetMetadata(ctx, map[string]interface{}{
		"workspaceId": repo.WorkspaceID.String(),
		"repoName":    repo.Name,
	})

	// Transition → indexing (also clears any previous index_error)
	if err := j.repositoryRepo.UpdateStatusAndError(ctx, repoID, models.RepositoryStatusIndexing, nil); err != nil {
		return j.markError(ctx, repoID, fmt.Errorf("index_repository: set indexing status: %w", err))
	}
	logger.Debug(ctx, "index_repository: status → indexing")

	// Resolve GitHub access token
	integration, err := j.integrationService.GetIntegrationByType(ctx, repo.WorkspaceID, github_integration.IntegrationTypeGitHub)
	if err != nil {
		return j.markError(ctx, repoID, fmt.Errorf("index_repository: github integration not found: %w", err))
	}
	accessToken, err := j.integrationService.GetDecryptedToken(ctx, integration.ID)
	if err != nil {
		return j.markError(ctx, repoID, fmt.Errorf("index_repository: decrypt token: %w", err))
	}

	owner, repoName, err := splitOwnerRepo(repo.Name)
	if err != nil {
		return j.markError(ctx, repoID, err)
	}

	branch := repo.DefaultBranch
	if branch == "" {
		branch = "HEAD"
	}

	// Fetch filtered file tree (built-in denylist + .gitignore + .indexignore).
	logger.Debug(ctx, fmt.Sprintf("index_repository: fetching file tree owner=%s repo=%s branch=%s", owner, repoName, branch))
	files, commitSHA, err := j.fileFetcher.GetFilteredFileTree(ctx, accessToken, owner, repoName, branch)
	if err != nil {
		return j.markError(ctx, repoID, fmt.Errorf("index_repository: fetch file tree: %w", err))
	}

	logger.Info(ctx, fmt.Sprintf("index_repository: file tree fetched fileCount=%d commitSHA=%s", len(files), commitSHA))

	// Cap large monorepos: keep highest-priority files first.
	if len(files) > indexer.MaxFilesPerIndex {
		files = prioritizeAndCap(files, indexer.MaxFilesPerIndex)
		logger.Info(ctx, fmt.Sprintf("index_repository: capped file tree to maxFiles=%d", indexer.MaxFilesPerIndex))
	}

	existingMeta := j.loadExistingFileMeta(ctx, repoID)

	// Soft-delete orphan paths that are no longer in the tree.
	orphansRemoved, err := j.removeOrphanFiles(ctx, repoID, files, existingMeta)
	if err != nil {
		return j.markError(ctx, repoID, err)
	}
	if orphansRemoved > 0 {
		logger.Info(ctx, fmt.Sprintf("index_repository: removed orphan files count=%d", orphansRemoved))
	}

	j.upsertRepositoryFiles(ctx, repoID, files)

	enqueued, err := j.enqueueChangedFiles(ctx, repoID, files, existingMeta, commitSHA, payload.Force)
	if err != nil {
		return j.markError(ctx, repoID, err)
	}

	if err := j.finalizeIndexing(ctx, repoID); err != nil {
		return j.markError(ctx, repoID, err)
	}

	logger.Info(ctx, fmt.Sprintf("index_repository: done fileCount=%d enqueued=%d orphans=%d commitSHA=%s", len(files), enqueued, orphansRemoved, commitSHA))
	return nil
}

type existingFileMeta struct {
	gitSHA         string
	indexerVersion string
	processed      bool
}

func (j *IndexRepositoryJob) loadExistingFileMeta(ctx context.Context, repoID uuid.UUID) map[string]existingFileMeta {
	existingFiles, err := j.fileIndexRepo.ListByRepoID(ctx, repoID)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("index_repository: list existing files: %w", err))
		return map[string]existingFileMeta{}
	}

	meta := make(map[string]existingFileMeta, len(existingFiles))
	for _, existingFile := range existingFiles {
		meta[existingFile.Path] = existingFileMeta{
			gitSHA:         existingFile.GitSHA,
			indexerVersion: existingFile.IndexerVersion,
			processed:      existingFile.ProcessedAt != nil,
		}
	}
	return meta
}

func (j *IndexRepositoryJob) removeOrphanFiles(ctx context.Context, repoID uuid.UUID, files []github_client.FileTreeItem, existing map[string]existingFileMeta) (int, error) {
	inTree := make(map[string]struct{}, len(files))
	for _, f := range files {
		inTree[f.Path] = struct{}{}
	}
	removed := 0
	for path := range existing {
		if _, ok := inTree[path]; ok {
			continue
		}
		if err := j.chunkEmbeddingRepo.DeleteByFile(ctx, repoID, path); err != nil {
			return removed, fmt.Errorf("index_repository: delete orphan chunks %s: %w", path, err)
		}
		if err := j.fileIndexRepo.DeleteByPath(ctx, repoID, path); err != nil {
			return removed, fmt.Errorf("index_repository: delete orphan file %s: %w", path, err)
		}
		logger.Debug(ctx, fmt.Sprintf("index_repository: removed orphan path=%s", path))
		removed++
	}
	return removed, nil
}

func (j *IndexRepositoryJob) upsertRepositoryFiles(ctx context.Context, repoID uuid.UUID, files []github_client.FileTreeItem) {
	for _, file := range files {
		language := chunker.DetectLanguage(file.Path)
		repositoryFile := &models.RepositoryFile{
			RepoID:   repoID,
			Path:     file.Path,
			GitSHA:   file.SHA,
			Size:     file.Size,
			Language: language,
		}
		if err := j.fileIndexRepo.Upsert(ctx, repositoryFile); err != nil {
			logger.Error(ctx, fmt.Errorf("index_repository: upsert repository_files %s: %w", file.Path, err))
			continue
		}
		logger.Debug(ctx, fmt.Sprintf("index_repository: upserted file path=%s gitSha=%s size=%d language=%s", file.Path, file.SHA, file.Size, language))
	}
}

func (j *IndexRepositoryJob) enqueueChangedFiles(ctx context.Context, repoID uuid.UUID, files []github_client.FileTreeItem, existing map[string]existingFileMeta, commitSHA string, force bool) (int, error) {
	// Enqueue high-priority source files first for faster time-to-useful-index.
	sorted := make([]github_client.FileTreeItem, len(files))
	copy(sorted, files)
	sort.SliceStable(sorted, func(i, j int) bool {
		return indexer.SortKey(sorted[i].Path) < indexer.SortKey(sorted[j].Path)
	})

	enqueued := 0
	for _, file := range sorted {
		meta, exists := existing[file.Path]
		unchangedSHA := exists && meta.gitSHA == file.SHA
		sameVersion := exists && meta.indexerVersion == indexer.Version && meta.processed
		if !force && unchangedSHA && sameVersion {
			logger.Debug(ctx, fmt.Sprintf("index_repository: skipping unchanged file path=%s gitSha=%s", file.Path, file.SHA))
			continue
		}
		// SHA unchanged but version mismatch → reprocess without treating as "changed" orphan.
		if exists {
			if err := j.resetChangedFile(ctx, repoID, file.Path); err != nil {
				return enqueued, err
			}
		}

		if err := j.enqueuer.EnqueueSummarizeFile(repoID, file.Path, commitSHA); err != nil {
			logger.Error(ctx, fmt.Errorf("index_repository: enqueue summarize task %s: %w", file.Path, err))
			continue
		}
		logger.Debug(ctx, fmt.Sprintf("index_repository: enqueued summarize_file path=%s priority=%d", file.Path, indexer.SortKey(file.Path)))
		enqueued++
	}
	return enqueued, nil
}

func (j *IndexRepositoryJob) resetChangedFile(ctx context.Context, repoID uuid.UUID, filePath string) error {
	if err := j.fileIndexRepo.ResetProcessingState(ctx, repoID, filePath); err != nil {
		return fmt.Errorf("index_repository: reset processing state %s: %w", filePath, err)
	}
	if err := j.chunkEmbeddingRepo.DeleteByFile(ctx, repoID, filePath); err != nil {
		return fmt.Errorf("index_repository: delete stale chunks %s: %w", filePath, err)
	}
	return nil
}

func prioritizeAndCap(files []github_client.FileTreeItem, max int) []github_client.FileTreeItem {
	if len(files) <= max {
		return files
	}
	sorted := make([]github_client.FileTreeItem, len(files))
	copy(sorted, files)
	sort.SliceStable(sorted, func(i, j int) bool {
		return indexer.SortKey(sorted[i].Path) < indexer.SortKey(sorted[j].Path)
	})
	return sorted[:max]
}

func (j *IndexRepositoryJob) finalizeIndexing(ctx context.Context, repoID uuid.UUID) error {
	transitioned, err := j.repositoryRepo.MarkIndexedIfReady(ctx, repoID, time.Now())
	if err != nil {
		return fmt.Errorf("index_repository: finalize status: %w", err)
	}
	if !transitioned {
		logger.Debug(ctx, "index_repository: waiting for summarize_file jobs before indexed transition")
		return nil
	}

	logger.Debug(ctx, "index_repository: status indexing → indexed")
	if err := j.enqueuer.EnqueueGenerateTopics(repoID); err != nil {
		logger.Error(ctx, fmt.Errorf("index_repository: enqueue generate_topics: %w", err))
		return nil
	}
	logger.Debug(ctx, "index_repository: enqueued generate_topics")
	return nil
}

// markError transitions the repository to error status, persists the error message, and returns the original error.
func (j *IndexRepositoryJob) markError(ctx context.Context, repoID uuid.UUID, err error) error {
	msg := err.Error()
	_ = j.repositoryRepo.UpdateStatusAndError(ctx, repoID, models.RepositoryStatusError, &msg)
	go func() {
		repoName := repoID.String()
		workspaceName := ""
		if repo, lookupErr := j.repositoryRepo.GetByID(context.Background(), repoID); lookupErr == nil && repo != nil {
			repoName = repo.Name
			workspaceName = repo.WorkspaceID.String()
		}
		app_notifications.SendRepositoryIndexingFailedNotification(context.Background(), app_notifications.RepositoryIndexingFailedParams{
			RepoName:      repoName,
			WorkspaceName: workspaceName,
			ErrorMessage:  msg,
		})
	}()
	return err
}
