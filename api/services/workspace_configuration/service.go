package workspace_configuration_service

import (
	"abstractions/config"
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/utils/encryption_utils"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	defaultIndexingModel   = "claude-haiku-4-5-20251001"
	defaultNewsletterModel = "claude-sonnet-4-6"
	defaultEmbeddingsModel = "text-embedding-3-small"
)

var envDefaults = map[string]func() string{
	KeyAnthropicAPIKey: func() string {
		if config.Environment == nil || config.Environment.DisableDefaultAIUsage {
			return ""
		}
		return config.Environment.AnthropicAPIKey
	},
	KeyOpenAIAPIKey: func() string {
		if config.Environment == nil || config.Environment.DisableDefaultAIUsage {
			return ""
		}
		return config.Environment.OpenAIAPIKey
	},
	KeyIndexingModel:   func() string { return defaultIndexingModel },
	KeyNewsletterModel: func() string { return defaultNewsletterModel },
	KeyEmbeddingsModel: func() string { return defaultEmbeddingsModel },
}

type Service struct {
	repo          *repositories.WorkspaceConfigurationRepository
	encryptionKey string
}

func NewService(repo *repositories.WorkspaceConfigurationRepository, encryptionKey string) *Service {
	return &Service{repo: repo, encryptionKey: encryptionKey}
}

// Set stores a configuration value for a workspace key. If the key is in EncryptedKeys the value
// is AES-256-GCM encrypted before storage. An empty value hard-deletes the row so the system
// default takes effect. For encrypted keys the unencrypted last 4 characters are stored in the
// last_four column so they can be surfaced in responses without decryption.
func (s *Service) Set(ctx context.Context, workspaceID uuid.UUID, key, value string, userID uuid.UUID) error {
	if value == "" {
		return s.repo.Delete(ctx, workspaceID, key)
	}

	storedValue := value
	encrypted := EncryptedKeys[key]
	if encrypted {
		enc, err := encryption_utils.EncryptString(s.encryptionKey, value)
		if err != nil {
			return fmt.Errorf("workspace_configuration_service: set: encrypt: %w", err)
		}
		storedValue = enc
	}

	last4 := value
	if len(value) > 4 {
		last4 = value[len(value)-4:]
	}

	cfg := &models.WorkspaceConfiguration{
		WorkspaceID: workspaceID,
		Key:         key,
		Value:       storedValue,
		Encrypted:   encrypted,
		LastFour:    &last4,
		CreatedBy:   userID,
		UpdatedBy:   userID,
		UpdatedAt:   time.Now(),
	}
	return s.repo.Upsert(ctx, cfg)
}

// Get returns the resolved value for a configuration key. If a workspace override exists it is
// returned (decrypted if needed). Otherwise the system environment default is returned.
func (s *Service) Get(ctx context.Context, workspaceID uuid.UUID, key string) (string, error) {
	cfg, err := s.repo.GetByWorkspaceAndKey(ctx, workspaceID, key)
	if err != nil {
		return "", fmt.Errorf("workspace_configuration_service: get: %w", err)
	}
	if cfg == nil {
		if fn, ok := envDefaults[key]; ok {
			return fn(), nil
		}
		return "", nil
	}
	if cfg.Encrypted {
		plain, err := encryption_utils.DecryptString(s.encryptionKey, cfg.Value)
		if err != nil {
			return "", fmt.Errorf("workspace_configuration_service: get: decrypt: %w", err)
		}
		return plain, nil
	}
	return cfg.Value, nil
}

// GetLastFour returns the last 4 characters stored in the last_four column at write time, or nil
// if no workspace override exists. Reading from the dedicated column avoids decrypting the secret.
func (s *Service) GetLastFour(ctx context.Context, workspaceID uuid.UUID, key string) (*string, error) {
	cfg, err := s.repo.GetByWorkspaceAndKey(ctx, workspaceID, key)
	if err != nil {
		return nil, fmt.Errorf("workspace_configuration_service: get_last_four: %w", err)
	}
	if cfg == nil {
		return nil, nil
	}
	return cfg.LastFour, nil
}

// GetResolvedModel returns the resolved model value for a model key (indexing or newsletter).
// Falls back to the system default if not set in the workspace.
func (s *Service) GetResolvedModel(ctx context.Context, workspaceID uuid.UUID, key string) (string, error) {
	return s.Get(ctx, workspaceID, key)
}
