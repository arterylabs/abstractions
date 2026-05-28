package repository_service

import (
	"abstractions/clients/app_notifications"
	github_client "abstractions/clients/github"
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/server/responses"
	integration_service "abstractions/services/integrations"
	github_integration "abstractions/services/integrations/core/github"
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrRepositoryNotFound = errors.New("repository not found")
	ErrAlreadyConnected   = errors.New("repository already connected to this workspace")
	ErrGitHubNotConnected = errors.New("github integration not connected for this workspace")
)

// Service handles repository connection business logic.
type Service struct {
	db                 *gorm.DB
	repo               *repositories.RepositoryRepository
	integrationService *integration_service.Service
}

func NewRepositoryService(db *gorm.DB) *Service {
	return &Service{
		db:                 db,
		repo:               repositories.NewRepositoryRepository(db),
		integrationService: integration_service.NewIntegrationService(db),
	}
}

func (s *Service) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

// ConnectRepositoryRequest carries the parameters needed to connect a repository.
type ConnectRepositoryRequest struct {
	WorkspaceID   uuid.UUID
	CreatedBy     uuid.UUID
	Name          string
	Description   string
	ServiceType   string
	ServiceID     string
	DefaultBranch string
}

// ConnectRepository creates a repository row and sets status=pending.
func (s *Service) ConnectRepository(ctx context.Context, req *ConnectRepositoryRequest) (*models.Repository, error) {
	serviceType := req.ServiceType
	if serviceType == "" {
		serviceType = "GITHUB"
	}

	existing, err := s.repo.GetByWorkspaceAndServiceID(ctx, req.WorkspaceID, serviceType, req.ServiceID)
	if err != nil {
		return nil, fmt.Errorf("repository_service: connect_repository: check_existing: %w", err)
	}
	if existing != nil {
		return nil, ErrAlreadyConnected
	}

	repository := &models.Repository{
		WorkspaceID:   req.WorkspaceID,
		CreatedBy:     req.CreatedBy,
		Name:          req.Name,
		Description:   req.Description,
		ServiceType:   serviceType,
		ServiceID:     req.ServiceID,
		DefaultBranch: req.DefaultBranch,
		Status:        models.RepositoryStatusPending,
	}

	if err := s.repo.Create(ctx, repository); err != nil {
		return nil, fmt.Errorf("repository_service: create repository: %w", err)
	}

	go func() {
		workspaceName := req.WorkspaceID.String()
		if ws, err := repositories.NewWorkspaceRepository(s.db).GetByID(context.Background(), req.WorkspaceID); err == nil && ws != nil {
			workspaceName = ws.Name
		}
		app_notifications.SendRepositoryConnectedNotification(context.Background(), app_notifications.RepositoryConnectedParams{
			RepoName:      repository.Name,
			ServiceType:   repository.ServiceType,
			WorkspaceName: workspaceName,
		})
	}()

	return repository, nil
}

// GetRepositoryInWorkspace returns a repository by ID scoped to a workspace.
func (s *Service) GetRepositoryInWorkspace(ctx context.Context, workspaceID, id uuid.UUID) (*models.Repository, error) {
	repo, err := s.repo.GetByIDInWorkspace(ctx, id, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("repository_service: get_repository_in_workspace: %w", err)
	}
	if repo == nil {
		return nil, ErrRepositoryNotFound
	}
	return repo, nil
}

// ListRepositories returns all repositories for a workspace.
func (s *Service) ListRepositories(ctx context.Context, workspaceID uuid.UUID) ([]*models.Repository, error) {
	return s.repo.GetByWorkspaceID(ctx, workspaceID)
}

// ListGitHubRepos fetches the authenticated user's GitHub repos using the stored OAuth token.
func (s *Service) ListGitHubRepos(ctx context.Context, workspaceID uuid.UUID) ([]github_client.GitHubRepo, error) {
	integration, err := s.integrationService.GetIntegrationByType(ctx, workspaceID, github_integration.IntegrationTypeGitHub)
	if err != nil {
		return nil, ErrGitHubNotConnected
	}
	accessToken, err := s.integrationService.GetDecryptedToken(ctx, integration.ID)
	if err != nil {
		return nil, fmt.Errorf("repository_service: decrypt token: %w", err)
	}
	return github_client.ListUserRepos(ctx, accessToken)
}

// UpdateSettingsRequest holds the fields callers may update on a repository.
type UpdateSettingsRequest struct {
	SendDay           int
	SendHour          int
	NewsletterEnabled *bool
}

func (s *Service) buildRepositoryFileStats(ctx context.Context, repoID uuid.UUID) ([]responses.RepositoryLanguageStat, *repositories.ProcessingCounts, error) {
	fileIndexRepo := repositories.NewRepositoryFileRepository(s.db)
	langStats, err := fileIndexRepo.GetLanguageStats(ctx, repoID)
	if err != nil {
		return nil, nil, err
	}

	processingCounts, err := fileIndexRepo.GetProcessingCounts(ctx, repoID)
	if err != nil {
		return nil, nil, err
	}

	breakdown := make([]responses.RepositoryLanguageStat, len(langStats))
	for i, ls := range langStats {
		pct := 0.0
		if processingCounts.TotalFiles > 0 {
			pct = float64(ls.FileCount) / float64(processingCounts.TotalFiles) * 100
		}
		breakdown[i] = responses.RepositoryLanguageStat{
			Language:   ls.Language,
			FileCount:  ls.FileCount,
			Percentage: pct,
		}
	}

	return breakdown, &processingCounts, nil
}

// UpdateDescription saves the description for a repository in a workspace.
func (s *Service) UpdateDescription(ctx context.Context, workspaceID, repoID uuid.UUID, description string) (*models.Repository, error) {
	repo, err := s.repo.GetByIDInWorkspace(ctx, repoID, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("repository_service: update_description: get_repository: %w", err)
	}
	if repo == nil {
		return nil, ErrRepositoryNotFound
	}
	repo.Description = description
	if err := s.repo.Update(ctx, repo); err != nil {
		return nil, fmt.Errorf("repository_service: update_description: update: %w", err)
	}
	return repo, nil
}

// UpdateSettings saves send_day, send_hour, and newsletter_enabled for a repository in a workspace.
func (s *Service) UpdateSettings(ctx context.Context, workspaceID, repoID uuid.UUID, req *UpdateSettingsRequest) (*models.Repository, error) {
	repo, err := s.repo.GetByIDInWorkspace(ctx, repoID, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("repository_service: update_settings: get_repository: %w", err)
	}
	if repo == nil {
		return nil, ErrRepositoryNotFound
	}
	repo.SendDay = req.SendDay
	repo.SendHour = req.SendHour
	if req.NewsletterEnabled != nil {
		repo.NewsletterEnabled = *req.NewsletterEnabled
	}
	if err := s.repo.Update(ctx, repo); err != nil {
		return nil, fmt.Errorf("repository_service: update_settings: update: %w", err)
	}
	return repo, nil
}

// GetStats returns file count and per-language breakdown for a repository.
func (s *Service) GetStats(ctx context.Context, workspaceID uuid.UUID, repoID uuid.UUID) (*responses.RepositoryStatsResponse, error) {
	repo, err := s.repo.GetByIDInWorkspace(ctx, repoID, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("repository_service: get_stats: get_repository: %w", err)
	}
	if repo == nil {
		return nil, ErrRepositoryNotFound
	}

	breakdown, counts, err := s.buildRepositoryFileStats(ctx, repoID)
	if err != nil {
		return nil, fmt.Errorf("repository_service: get_stats: %w", err)
	}

	return &responses.RepositoryStatsResponse{
		FileCount:         counts.TotalFiles,
		LanguageBreakdown: breakdown,
	}, nil
}

// GetStatus returns the repository indexing status together with live file stats.
func (s *Service) GetStatus(ctx context.Context, workspaceID uuid.UUID, repoID uuid.UUID) (*responses.RepositoryStatusResponse, error) {
	repo, err := s.repo.GetByIDInWorkspace(ctx, repoID, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("repository_service: get_status: get_repository: %w", err)
	}
	if repo == nil {
		return nil, ErrRepositoryNotFound
	}

	breakdown, counts, err := s.buildRepositoryFileStats(ctx, repoID)
	if err != nil {
		return nil, fmt.Errorf("repository_service: get_status: %w", err)
	}

	progressPercentage := 0.0
	if counts.TotalFiles > 0 {
		progressPercentage = float64(counts.SummarizedFiles) / float64(counts.TotalFiles) * 100
	}
	if repo.Status == models.RepositoryStatusIndexed {
		progressPercentage = 100
	}
	progressPercentage = math.Max(0, math.Min(100, progressPercentage))

	var indexedAt *string
	if repo.IndexedAt != nil {
		formatted := repo.IndexedAt.UTC().Format(time.RFC3339)
		indexedAt = &formatted
	}

	return &responses.RepositoryStatusResponse{
		Status:             string(repo.Status),
		IndexedAt:          indexedAt,
		IndexError:         repo.IndexError,
		TotalFiles:         counts.TotalFiles,
		SummarizedFiles:    counts.SummarizedFiles,
		ProgressPercentage: progressPercentage,
		LanguageBreakdown:  breakdown,
	}, nil
}

// DeleteRepository soft-deletes the repository row and removes all derived artifacts:
// chunk_embeddings, file_index, newsletter_topics, and newsletter rows.
func (s *Service) DeleteRepository(ctx context.Context, workspaceID, repoID uuid.UUID, userID uuid.UUID) error {
	repo, err := s.repo.GetByIDInWorkspace(ctx, repoID, workspaceID)
	if err != nil {
		return fmt.Errorf("repository_service: get repo: %w", err)
	}
	if repo == nil {
		return ErrRepositoryNotFound
	}

	if err := s.repo.SoftDelete(ctx, repoID, userID); err != nil {
		return fmt.Errorf("repository_service: soft delete: %w", err)
	}

	chunkRepo := repositories.NewChunkEmbeddingRepository(s.db)
	if err := chunkRepo.DeleteByRepo(ctx, repoID); err != nil {
		return fmt.Errorf("repository_service: delete chunks: %w", err)
	}

	fileIndexRepo := repositories.NewRepositoryFileRepository(s.db)
	if err := fileIndexRepo.DeleteByRepoID(ctx, repoID); err != nil {
		return fmt.Errorf("repository_service: delete repository_files: %w", err)
	}

	newsletterRepo := repositories.NewNewsletterRepository(s.db)
	if err := newsletterRepo.DeleteByRepoID(ctx, repoID); err != nil {
		return fmt.Errorf("repository_service: delete newsletters: %w", err)
	}

	return nil
}
