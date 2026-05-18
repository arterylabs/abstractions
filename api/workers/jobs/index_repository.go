package jobs

import (
	github_client "abstractions/clients/github"
	"abstractions/config"
	"abstractions/db/models"
	"abstractions/logger"
	"abstractions/repositories"
	integration_service "abstractions/services/integrations"
	github_integration "abstractions/services/integrations/core/github"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

const IndexRepositoryJobID = "index_repository"

// MaxConcurrentSummarisations is the maximum number of SummarizeFileJob tasks that can be
// enqueued in a single IndexRepositoryJob. asynq concurrency is controlled server-side;
// this constant documents the intended limit.
const MaxConcurrentSummarisations = 20

// FileTreeFetcher is the interface for fetching a filtered repository file tree from GitHub.
// The second return value is the commit SHA at HEAD; callers must use this as the ref when
// fetching individual file contents — individual file SHAs are blob SHAs, not commit SHAs,
// and are not accepted by the GitHub Contents API.
type FileTreeFetcher interface {
	GetFilteredFileTree(ctx context.Context, accessToken, owner, repo string) ([]github_client.FileTreeItem, string, error)
}

// defaultFileTreeFetcher delegates to the real github_client package function.
type defaultFileTreeFetcher struct{}

func (d *defaultFileTreeFetcher) GetFilteredFileTree(ctx context.Context, accessToken, owner, repo string) ([]github_client.FileTreeItem, string, error) {
	return github_client.GetFilteredFileTree(ctx, accessToken, owner, repo)
}

// JobEnqueuer is the interface used by job handlers that need to sub-enqueue tasks.
// It exposes only job-level operations, keeping asynq details inside the jobs package.
type JobEnqueuer interface {
	EnqueueSummarizeFile(repositoryID uuid.UUID, filePath, gitSHA string) error
	EnqueueGenerateTopics(repoID uuid.UUID) error
	EnqueueGenerateNewsletter(repoID uuid.UUID) error
}

// IndexRepositoryPayload is the JSON payload for an IndexRepositoryJob task.
type IndexRepositoryPayload struct {
	RepositoryID string `json:"repositoryId"`
}

// IndexRepositoryJob fetches the file tree, seeds repository_files, updates repo status, and
// enqueues one SummarizeFileJob per file (capped at MaxConcurrentSummarisations in flight).
type IndexRepositoryJob struct {
	db                 *gorm.DB
	config             *config.Config
	fileFetcher        FileTreeFetcher
	enqueuer           JobEnqueuer
	integrationService *integration_service.Service
	fileIndexRepo      *repositories.RepositoryFileRepository
	repositoryRepo     *repositories.RepositoryRepository
}

func NewIndexRepositoryJob(db *gorm.DB, cfg *config.Config, enqueuer JobEnqueuer) *IndexRepositoryJob {
	return &IndexRepositoryJob{
		db:                 db,
		config:             cfg,
		fileFetcher:        &defaultFileTreeFetcher{},
		enqueuer:           enqueuer,
		integrationService: integration_service.NewIntegrationService(db),
		fileIndexRepo:      repositories.NewRepositoryFileRepository(db),
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
		integrationService: integration_service.NewIntegrationService(db),
		fileIndexRepo:      repositories.NewRepositoryFileRepository(db),
		repositoryRepo:     repositories.NewRepositoryRepository(db),
	}
}

func (j *IndexRepositoryJob) ID() string {
	return IndexRepositoryJobID
}

// NewIndexRepositoryTask builds an asynq.Task for the given repository ID.
func NewIndexRepositoryTask(repositoryID uuid.UUID) (*asynq.Task, error) {
	payload := IndexRepositoryPayload{RepositoryID: repositoryID.String()}
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

	logger.Info(ctx, "index_repository: starting")

	repoCtx, err := LoadRepositoryContext(ctx, repoID, j.repositoryRepo)
	if err != nil {
		return j.markError(ctx, repoID, fmt.Errorf("index_repository: %w", err))
	}
	repo := repoCtx.Repository

	ctx = logger.SetMetadata(ctx, map[string]interface{}{
		"workspaceId": repo.WorkspaceID.String(),
		"repoName":    repo.Name,
	})

	// Transition: pending → indexing (also clears any previous index_error)
	if err := j.repositoryRepo.UpdateStatusAndError(ctx, repoID, models.RepositoryStatusIndexing, nil); err != nil {
		return j.markError(ctx, repoID, fmt.Errorf("index_repository: set indexing status: %w", err))
	}
	logger.Debug(ctx, "index_repository: status pending → indexing")

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

	// Fetch filtered file tree.
	// commitSHA is the HEAD commit SHA and must be used as the ref when fetching individual
	// file contents — each FileTreeItem.SHA is a blob SHA, not a commit SHA.
	logger.Debug(ctx, fmt.Sprintf("index_repository: fetching file tree owner=%s repo=%s", owner, repoName))
	files, commitSHA, err := j.fileFetcher.GetFilteredFileTree(ctx, accessToken, owner, repoName)
	if err != nil {
		return j.markError(ctx, repoID, fmt.Errorf("index_repository: fetch file tree: %w", err))
	}

	logger.Info(ctx, fmt.Sprintf("index_repository: file tree fetched fileCount=%d commitSHA=%s", len(files), commitSHA))

	// Load existing file records before upserting so we can detect SHA changes.
	// If the lookup fails we treat all files as new/changed to avoid skipping work.
	existingFiles, err := j.fileIndexRepo.ListByRepoID(ctx, repoID)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("index_repository: list existing files: %w", err))
		existingFiles = nil
	}
	existingSHAs := make(map[string]string, len(existingFiles))
	for _, ef := range existingFiles {
		existingSHAs[ef.Path] = ef.GitSHA
	}

	// Bulk upsert file paths/SHAs into repository_files.
	// Tags, Language, and ProcessedAt are intentionally omitted here — they are
	// populated exclusively by SummarizeFileJob via UpdateMetadata.
	now := time.Now()
	for _, f := range files {
		fi := &models.RepositoryFile{
			RepoID: repoID,
			Path:   f.Path,
			GitSHA: f.SHA,
			Size:   f.Size,
		}
		if err := j.fileIndexRepo.Upsert(ctx, fi); err != nil {
			logger.Error(ctx, fmt.Errorf("index_repository: upsert repository_files %s: %w", f.Path, err))
		} else {
			logger.Debug(ctx, fmt.Sprintf("index_repository: upserted file path=%s gitSha=%s size=%d", f.Path, f.SHA, f.Size))
		}
	}

	// Enqueue one SummarizeFileJob per file that is new or has a changed SHA.
	// Files whose SHA matches the previously stored value are already up-to-date and skipped.
	// Use commitSHA as the ref so the GitHub Contents API can locate the file; individual
	// file SHAs (f.SHA) are blob SHAs and are rejected by that API.
	enqueued := 0
	for _, f := range files {
		if oldSHA, exists := existingSHAs[f.Path]; exists && oldSHA == f.SHA {
			logger.Debug(ctx, fmt.Sprintf("index_repository: skipping unchanged file path=%s gitSha=%s", f.Path, f.SHA))
			continue
		}
		if err := j.enqueuer.EnqueueSummarizeFile(repoID, f.Path, commitSHA); err != nil {
			logger.Error(ctx, fmt.Errorf("index_repository: enqueue summarize task %s: %w", f.Path, err))
			continue
		}
		logger.Debug(ctx, fmt.Sprintf("index_repository: enqueued summarize_file path=%s", f.Path))
		enqueued++
		if enqueued >= MaxConcurrentSummarisations {
			logger.Warn(ctx, fmt.Sprintf("index_repository: capped summarize enqueue at %d totalFiles=%d", MaxConcurrentSummarisations, len(files)))
			break
		}
	}

	// Transition: indexing → indexed and stamp IndexedAt
	repo.Status = models.RepositoryStatusIndexed
	repo.IndexedAt = &now
	if err := j.repositoryRepo.Update(ctx, repo); err != nil {
		return fmt.Errorf("index_repository: set indexed status: %w", err)
	}
	logger.Debug(ctx, fmt.Sprintf("index_repository: status indexing → indexed indexedAt=%s", now.Format(time.RFC3339)))

	// Enqueue topic generation 30 minutes after indexing so file summaries are ready.
	if err := j.enqueuer.EnqueueGenerateTopics(repoID); err != nil {
		logger.Error(ctx, fmt.Errorf("index_repository: enqueue generate_topics: %w", err))
	} else {
		logger.Debug(ctx, "index_repository: enqueued generate_topics processIn=30m")
	}

	logger.Info(ctx, fmt.Sprintf("index_repository: done fileCount=%d enqueued=%d commitSHA=%s", len(files), enqueued, commitSHA))
	return nil
}

// markError transitions the repository to error status, persists the error message, and returns the original error.
func (j *IndexRepositoryJob) markError(ctx context.Context, repoID uuid.UUID, err error) error {
	msg := err.Error()
	_ = j.repositoryRepo.UpdateStatusAndError(ctx, repoID, models.RepositoryStatusError, &msg)
	return err
}
