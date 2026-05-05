package repository_service

import (
	github_client "abstractions/clients/github"
	"abstractions/config"
	"abstractions/db/models"
	"abstractions/repositories"
	integration_service "abstractions/services/integrations"
	github_integration "abstractions/services/integrations/core/github"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrRepositoryNotFound    = errors.New("repository not found")
	ErrAlreadyConnected      = errors.New("repository already connected to this workspace")
	ErrGitHubNotConnected    = errors.New("github integration not connected for this workspace")
	ErrWebhookCreationFailed = errors.New("failed to create webhook on github")
)

// WebhookCreator is an interface for creating GitHub webhooks, allowing the real client
// to be swapped out for a mock in tests.
type WebhookCreator interface {
	CreateWebhook(ctx context.Context, accessToken, owner, repo, payloadURL, secret string) (int64, error)
}

// defaultWebhookCreator wraps the real GitHub client.
type defaultWebhookCreator struct{}

func (d *defaultWebhookCreator) CreateWebhook(ctx context.Context, accessToken, owner, repo, payloadURL, secret string) (int64, error) {
	return github_client.CreateWebhook(ctx, accessToken, owner, repo, payloadURL, secret)
}

// Service handles repository connection business logic.
type Service struct {
	db                 *gorm.DB
	repo               *repositories.RepositoryRepository
	integrationService *integration_service.Service
	webhookCreator     WebhookCreator
}

func NewRepositoryService(db *gorm.DB) *Service {
	return &Service{
		db:                 db,
		repo:               repositories.NewRepositoryRepository(db),
		integrationService: integration_service.NewIntegrationService(db),
		webhookCreator:     &defaultWebhookCreator{},
	}
}

// NewRepositoryServiceWithWebhookCreator creates a service with a custom webhook creator (for tests).
func NewRepositoryServiceWithWebhookCreator(db *gorm.DB, wc WebhookCreator) *Service {
	svc := NewRepositoryService(db)
	svc.webhookCreator = wc
	return svc
}

// ConnectRepositoryRequest carries the parameters needed to connect a GitHub repository.
type ConnectRepositoryRequest struct {
	WorkspaceID   uuid.UUID
	UserID        uuid.UUID
	GitHubRepoID  int64
	FullName      string
	DefaultBranch string
}

// ConnectRepository creates a repository row, creates a GitHub webhook, and sets status=pending.
func (s *Service) ConnectRepository(ctx context.Context, req *ConnectRepositoryRequest) (*models.Repository, error) {
	// Check for duplicate
	existing, err := s.repo.GetByWorkspaceAndGitHubID(ctx, req.WorkspaceID, req.GitHubRepoID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrAlreadyConnected
	}

	// Fetch the GitHub access token for this workspace
	integration, err := s.integrationService.GetIntegrationByType(ctx, req.WorkspaceID, github_integration.IntegrationTypeGitHub)
	if err != nil {
		return nil, ErrGitHubNotConnected
	}
	accessToken, err := s.integrationService.GetDecryptedToken(ctx, integration.ID)
	if err != nil {
		return nil, fmt.Errorf("repository_service: decrypt token: %w", err)
	}

	// Generate a cryptographically random per-repo webhook secret
	webhookSecret, err := generateWebhookSecret()
	if err != nil {
		return nil, fmt.Errorf("repository_service: generate webhook secret: %w", err)
	}

	// Split full_name into owner/repo for the webhook API call
	owner, repoName, err := splitFullName(req.FullName)
	if err != nil {
		return nil, err
	}

	payloadURL := config.Environment.GitHubWebhookBaseURL + "/webhooks/github"

	webhookID, err := s.webhookCreator.CreateWebhook(ctx, accessToken, owner, repoName, payloadURL, webhookSecret)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrWebhookCreationFailed, err)
	}

	repository := &models.Repository{
		WorkspaceID:   req.WorkspaceID,
		UserID:        req.UserID,
		GitHubRepoID:  req.GitHubRepoID,
		FullName:      req.FullName,
		DefaultBranch: req.DefaultBranch,
		WebhookID:     &webhookID,
		WebhookSecret: &webhookSecret,
		Status:        models.RepositoryStatusPending,
	}

	if err := s.repo.Create(ctx, repository); err != nil {
		return nil, fmt.Errorf("repository_service: create repository: %w", err)
	}

	return repository, nil
}

// GetRepository returns a repository by ID.
func (s *Service) GetRepository(ctx context.Context, id uuid.UUID) (*models.Repository, error) {
	repo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
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

// generateWebhookSecret creates a 32-byte random hex string.
func generateWebhookSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// splitFullName splits "owner/repo" into its parts.
func splitFullName(fullName string) (string, string, error) {
	for i, c := range fullName {
		if c == '/' {
			return fullName[:i], fullName[i+1:], nil
		}
	}
	return "", "", fmt.Errorf("repository_service: invalid full_name %q (expected owner/repo)", fullName)
}
