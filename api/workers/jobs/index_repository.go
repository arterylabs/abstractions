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
type FileTreeFetcher interface {
	GetFilteredFileTree(ctx context.Context, accessToken, owner, repo string) ([]github_client.FileTreeItem, error)
}

// defaultFileTreeFetcher delegates to the real github_client package function.
type defaultFileTreeFetcher struct{}

func (d *defaultFileTreeFetcher) GetFilteredFileTree(ctx context.Context, accessToken, owner, repo string) ([]github_client.FileTreeItem, error) {
	return github_client.GetFilteredFileTree(ctx, accessToken, owner, repo)
}

// TaskEnqueuer is the interface for enqueuing asynq tasks.
type TaskEnqueuer interface {
	Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)
}

// IndexRepositoryPayload is the JSON payload for an IndexRepositoryJob task.
type IndexRepositoryPayload struct {
	RepositoryID string `json:"repositoryId"`
}

// IndexRepositoryJob fetches the file tree, seeds file_index, updates repo status, and
// enqueues one SummarizeFileJob per file (capped at MaxConcurrentSummarisations in flight).
type IndexRepositoryJob struct {
	db                 *gorm.DB
	config             *config.Config
	fileFetcher        FileTreeFetcher
	enqueuer           TaskEnqueuer
	integrationService *integration_service.Service
	fileIndexRepo      *repositories.FileIndexRepository
	repositoryRepo     *repositories.RepositoryRepository
}

func NewIndexRepositoryJob(db *gorm.DB, cfg *config.Config, enqueuer TaskEnqueuer) *IndexRepositoryJob {
	return &IndexRepositoryJob{
		db:                 db,
		config:             cfg,
		fileFetcher:        &defaultFileTreeFetcher{},
		enqueuer:           enqueuer,
		integrationService: integration_service.NewIntegrationService(db),
		fileIndexRepo:      repositories.NewFileIndexRepository(db),
		repositoryRepo:     repositories.NewRepositoryRepository(db),
	}
}

// NewIndexRepositoryJobWithDeps creates a job with injected dependencies for testing.
func NewIndexRepositoryJobWithDeps(
	db *gorm.DB,
	cfg *config.Config,
	fileFetcher FileTreeFetcher,
	enqueuer TaskEnqueuer,
) *IndexRepositoryJob {
	return &IndexRepositoryJob{
		db:                 db,
		config:             cfg,
		fileFetcher:        fileFetcher,
		enqueuer:           enqueuer,
		integrationService: integration_service.NewIntegrationService(db),
		fileIndexRepo:      repositories.NewFileIndexRepository(db),
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

	logger.Info(ctx, "index_repository: starting repoId="+repoID.String())

	repo, err := j.repositoryRepo.GetByID(ctx, repoID)
	if err != nil || repo == nil {
		return j.markError(ctx, repoID, fmt.Errorf("index_repository: repo not found: %w", err))
	}

	// Transition: pending → indexing
	if err := j.repositoryRepo.UpdateStatus(ctx, repoID, models.RepositoryStatusIndexing); err != nil {
		return j.markError(ctx, repoID, fmt.Errorf("index_repository: set indexing status: %w", err))
	}

	// Resolve GitHub access token
	integration, err := j.integrationService.GetIntegrationByType(ctx, repo.WorkspaceID, github_integration.IntegrationTypeGitHub)
	if err != nil {
		return j.markError(ctx, repoID, fmt.Errorf("index_repository: github integration not found: %w", err))
	}
	accessToken, err := j.integrationService.GetDecryptedToken(ctx, integration.ID)
	if err != nil {
		return j.markError(ctx, repoID, fmt.Errorf("index_repository: decrypt token: %w", err))
	}

	owner, repoName, err := splitOwnerRepo(repo.FullName)
	if err != nil {
		return j.markError(ctx, repoID, err)
	}

	// Fetch filtered file tree
	files, err := j.fileFetcher.GetFilteredFileTree(ctx, accessToken, owner, repoName)
	if err != nil {
		return j.markError(ctx, repoID, fmt.Errorf("index_repository: fetch file tree: %w", err))
	}

	logger.Info(ctx, fmt.Sprintf("index_repository: found %d files repoId=%s", len(files), repoID))

	// Bulk upsert file paths/SHAs into file_index
	now := time.Now()
	for _, f := range files {
		fi := &models.FileIndex{
			RepoID:   repoID,
			FilePath: f.Path,
			GitSHA:   f.SHA,
			// domain_tags and last_processed_at are populated by SummarizeFileJob
			LastProcessedAt: &now,
		}
		if err := j.fileIndexRepo.Upsert(ctx, fi); err != nil {
			logger.Info(ctx, fmt.Sprintf("index_repository: upsert file_index error %s: %v", f.Path, err))
		}
	}

	// Enqueue one SummarizeFileJob per file (cap at MaxConcurrentSummarisations in-flight)
	enqueued := 0
	for _, f := range files {
		task, err := NewSummarizeFileTask(repoID, f.Path, f.SHA)
		if err != nil {
			logger.Info(ctx, fmt.Sprintf("index_repository: build task error %s: %v", f.Path, err))
			continue
		}
		if _, err := j.enqueuer.Enqueue(task, asynq.MaxRetry(3), asynq.Queue("default")); err != nil {
			logger.Info(ctx, fmt.Sprintf("index_repository: enqueue error %s: %v", f.Path, err))
			continue
		}
		enqueued++
		if enqueued >= MaxConcurrentSummarisations {
			break
		}
	}

	// Transition: indexing → indexed and stamp IndexedAt
	repo.Status = models.RepositoryStatusIndexed
	repo.IndexedAt = &now
	if err := j.repositoryRepo.Update(ctx, repo); err != nil {
		return fmt.Errorf("index_repository: set indexed status: %w", err)
	}

	logger.Info(ctx, fmt.Sprintf("index_repository: done repoId=%s files=%d enqueued=%d", repoID, len(files), enqueued))
	return nil
}

// markError transitions the repository to error status and returns the original error.
func (j *IndexRepositoryJob) markError(ctx context.Context, repoID uuid.UUID, err error) error {
	_ = j.repositoryRepo.UpdateStatus(ctx, repoID, models.RepositoryStatusError)
	return err
}
