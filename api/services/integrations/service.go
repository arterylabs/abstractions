package integration_service

import (
	"abstractions/config"
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/services/integrations/enums"
	"abstractions/utils/encryption_utils"
	"abstractions/utils/pagination_utils"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrIntegrationNotFound    = errors.New("integration not found")
	ErrDuplicateIntegration   = errors.New("integration already exists for this type")
	ErrInvalidConfiguration   = errors.New("invalid integration configuration")
	ErrInvalidIntegrationType = errors.New("invalid integration type")
	ErrEncryptionFailed       = errors.New("failed to encrypt token")
	ErrDecryptionFailed       = errors.New("failed to decrypt token")
	ErrMissingEncryptionKey   = errors.New("encryption key not configured")
	ErrUnsupportedIntegration = errors.New("unsupported integration type")
)

type Service struct {
	DB            *gorm.DB
	Repo          *repositories.IntegrationRepository
	EncryptionKey string
}

type CreateIntegrationRequest struct {
	WorkspaceID     uuid.UUID                 `json:"workspace_id"`
	Type            enums.IntegrationTypeEnum `json:"type"`
	ResourceType    *string                   `json:"resource_type,omitempty"`
	ResourceID      *uuid.UUID                `json:"resource_id,omitempty"`
	ServiceTitle    *string                   `json:"service_title,omitempty"`
	ServiceIcon     *string                   `json:"service_icon,omitempty"`
	Token           string                    `json:"token"`
	RefreshToken    *string                   `json:"refresh_token,omitempty"`
	TokenExpiresAt  *time.Time                `json:"token_expires_at,omitempty"`
	Scopes          *string                   `json:"scopes,omitempty"`
	ServiceTenantID *string                   `json:"service_tenant_id,omitempty"`
	ServiceUserID   *string                   `json:"service_user_id,omitempty"`
	Configuration   models.JSONBField         `json:"configuration,omitempty"`
	CreatedBy       uuid.UUID                 `json:"created_by"`
}

func NewIntegrationService(db *gorm.DB) *Service {
	return &Service{
		DB:            db,
		Repo:          repositories.NewIntegrationRepository(db),
		EncryptionKey: config.Environment.EncryptionKey32Bit,
	}
}

// CreateIntegration creates a new integration with encrypted tokens
func (s *Service) CreateIntegration(ctx context.Context, req *CreateIntegrationRequest) (*models.Integration, error) {
	// Validate integration type is supported
	if !IsValidType(req.Type) {
		return nil, ErrInvalidIntegrationType
	}

	// Get integration type implementation for validation
	integrationType, err := GetIntegrationType(req.Type)
	if err != nil {
		return nil, ErrUnsupportedIntegration
	}

	// Validate configuration if provided
	if req.Configuration != nil {
		if err := integrationType.ValidateConfiguration(req.Configuration); err != nil {
			return nil, ErrInvalidConfiguration
		}
	}

	// // Check if integration already exists for this type, workspace, and resource
	// exists, err := s.Repo.ExistsByWorkspaceIDAndType(ctx, req.WorkspaceID, req.Type, req.ResourceType, req.ResourceID)
	// if err != nil {
	// 	return nil, fmt.Errorf("integration_service: create_integration: check_exists: %w", err)
	// }
	// if exists {
	// 	return nil, ErrDuplicateIntegration
	// }

	// Encrypt the access token
	encryptedToken, err := s.encryptToken(req.Token)
	if err != nil {
		return nil, fmt.Errorf("integration_service: create_integration: encrypt_token: %w", err)
	}

	// Encrypt the refresh token if provided
	var encryptedRefreshToken *string
	if req.RefreshToken != nil && *req.RefreshToken != "" {
		encrypted, err := s.encryptToken(*req.RefreshToken)
		if err != nil {
			return nil, fmt.Errorf("integration_service: create_integration: encrypt_refresh_token: %w", err)
		}
		encryptedRefreshToken = &encrypted
	}

	integration := &models.Integration{
		Type:            req.Type,
		ServiceTitle:    req.ServiceTitle,
		ServiceIcon:     req.ServiceIcon,
		WorkspaceID:     req.WorkspaceID,
		ResourceType:    req.ResourceType,
		ResourceID:      req.ResourceID,
		Token:           encryptedToken,
		RefreshToken:    encryptedRefreshToken,
		TokenExpiresAt:  req.TokenExpiresAt,
		Scopes:          req.Scopes,
		ServiceTenantID: req.ServiceTenantID,
		ServiceUserID:   req.ServiceUserID,
		Configuration:   req.Configuration,
		CreatedBy:       req.CreatedBy,
		UpdatedBy:       req.CreatedBy,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	if err := s.Repo.Create(ctx, integration); err != nil {
		return nil, fmt.Errorf("integration_service: create_integration: create: %w", err)
	}

	return integration, nil
}

// GetDecryptedToken returns the decrypted access token for an integration
func (s *Service) GetDecryptedToken(ctx context.Context, integrationID uuid.UUID) (string, error) {
	integration, err := s.GetIntegration(ctx, integrationID)
	if err != nil {
		return "", err
	}

	return s.decryptToken(integration.Token)
}

// GetDecryptedRefreshToken returns the decrypted refresh token for an integration
func (s *Service) GetDecryptedRefreshToken(ctx context.Context, integrationID uuid.UUID) (*string, error) {
	integration, err := s.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	if integration.RefreshToken == nil {
		return nil, nil
	}

	decrypted, err := s.decryptToken(*integration.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &decrypted, nil
}

// UpdateTokens updates the access and refresh tokens for an integration
func (s *Service) UpdateTokens(ctx context.Context, integrationID uuid.UUID, token string, refreshToken *string, expiresAt *time.Time, updatedBy uuid.UUID) error {
	// Encrypt the new access token
	encryptedToken, err := s.encryptToken(token)
	if err != nil {
		return fmt.Errorf("integration_service: update_tokens: encrypt_token: %w", err)
	}

	// Encrypt the refresh token if provided
	var encryptedRefreshToken *string
	if refreshToken != nil && *refreshToken != "" {
		encrypted, err := s.encryptToken(*refreshToken)
		if err != nil {
			return fmt.Errorf("integration_service: update_tokens: encrypt_refresh_token: %w", err)
		}
		encryptedRefreshToken = &encrypted
	}

	return s.Repo.UpdateTokens(ctx, integrationID, encryptedToken, encryptedRefreshToken, expiresAt, updatedBy)
}

// GetSupportedIntegrationTypes returns a list of all supported integration types
func (s *Service) GetSupportedIntegrationTypes(ctx context.Context) ([]IntegrationType, error) {
	return GetAllIntegrationTypes(), nil
}

// ValidateIntegrationConfiguration validates the configuration for a specific integration type
func (s *Service) ValidateIntegrationConfiguration(integrationType enums.IntegrationTypeEnum, config models.JSONBField) error {
	integrationImpl, err := GetIntegrationType(integrationType)
	if err != nil {
		return fmt.Errorf("integration_service: validate_integration_configuration: get_type: %w", err)
	}

	return integrationImpl.ValidateConfiguration(config)
}

// encryptToken encrypts a token using the service's encryption key
func (s *Service) encryptToken(token string) (string, error) {
	if s.EncryptionKey == "" {
		return "", ErrMissingEncryptionKey
	}

	encrypted, err := encryption_utils.EncryptString(s.EncryptionKey, token)
	if err != nil {
		return "", ErrEncryptionFailed
	}

	return encrypted, nil
}

// decryptToken decrypts a token using the service's encryption key
func (s *Service) decryptToken(encryptedToken string) (string, error) {
	if s.EncryptionKey == "" {
		return "", ErrMissingEncryptionKey
	}

	decrypted, err := encryption_utils.DecryptString(s.EncryptionKey, encryptedToken)
	if err != nil {
		return "", ErrDecryptionFailed
	}

	return decrypted, nil
}

// GetIntegration retrieves an integration by ID
func (s *Service) GetIntegration(ctx context.Context, id uuid.UUID) (*models.Integration, error) {
	integration, err := s.Repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if integration == nil {
		return nil, ErrIntegrationNotFound
	}
	return integration, nil
}

// GetIntegrationsByWorkspace retrieves all integrations for a workspace
func (s *Service) GetIntegrationsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]*models.Integration, error) {
	return s.Repo.GetByWorkspaceID(ctx, workspaceID)
}

// GetIntegrationsPaginated retrieves paginated integrations for a workspace
func (s *Service) GetIntegrationsPaginated(ctx context.Context, workspaceID uuid.UUID, params pagination_utils.PaginationParams) ([]*models.Integration, int, error) {
	return s.Repo.GetPaginated(ctx, workspaceID, params)
}

// GetIntegrationsPaginatedByResource retrieves paginated integrations for a workspace filtered by resource
func (s *Service) GetIntegrationsPaginatedByResource(ctx context.Context, workspaceID uuid.UUID, params pagination_utils.PaginationParams, resourceType *string, resourceID *uuid.UUID) ([]*models.Integration, int, error) {
	return s.Repo.GetPaginatedByResource(ctx, workspaceID, params, resourceType, resourceID)
}

// UpdateIntegration updates an existing integration
func (s *Service) UpdateIntegration(ctx context.Context, integration *models.Integration) error {
	// get integration by id
	existingIntegration, err := s.GetIntegration(ctx, integration.ID)
	if err != nil {
		return fmt.Errorf("integration_service: update_integration: get_integration: %w", err)
	}
	if existingIntegration == nil {
		return ErrIntegrationNotFound
	}
	return s.Repo.Update(ctx, integration)
}

// UpdateIntegrationConfiguration updates only the configuration of an integration
func (s *Service) UpdateIntegrationConfiguration(ctx context.Context, id uuid.UUID, configuration models.JSONBField, updatedBy uuid.UUID) error {
	// get integration by id
	existingIntegration, err := s.GetIntegration(ctx, id)
	if err != nil {
		return fmt.Errorf("integration_service: update_integration_configuration: get_integration: %w", err)
	}
	if existingIntegration == nil {
		return ErrIntegrationNotFound
	}
	return s.Repo.UpdateConfiguration(ctx, id, configuration, updatedBy)
}

// DeleteIntegration soft deletes an integration
func (s *Service) DeleteIntegration(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	// get integration by id
	existingIntegration, err := s.GetIntegration(ctx, id)
	if err != nil {
		return fmt.Errorf("integration_service: delete_integration: get_integration: %w", err)
	}
	if existingIntegration == nil {
		return ErrIntegrationNotFound
	}
	return s.Repo.Delete(ctx, id, deletedBy)
}

// GetIntegrationsByType retrieves integrations by type for a workspace
func (s *Service) GetIntegrationsByType(ctx context.Context, workspaceID uuid.UUID, integrationType enums.IntegrationTypeEnum) ([]*models.Integration, error) {
	return s.Repo.GetMultipleByWorkspaceIDAndType(ctx, workspaceID, integrationType)
}

func (s *Service) GetIntegrationByType(ctx context.Context, workspaceID uuid.UUID, integrationType enums.IntegrationTypeEnum) (*models.Integration, error) {
	integration, err := s.Repo.GetByWorkspaceIDAndType(ctx, workspaceID, integrationType)
	if err != nil {
		return nil, err
	}
	if integration == nil {
		return nil, ErrIntegrationNotFound
	}
	return integration, nil
}

// GetIntegrationByTypeAndResource retrieves an integration by type and resource for a workspace
func (s *Service) GetIntegrationByTypeAndResource(ctx context.Context, workspaceID uuid.UUID, integrationType enums.IntegrationTypeEnum, resourceType *string, resourceID *uuid.UUID) (*models.Integration, error) {
	integration, err := s.Repo.GetByWorkspaceIDTypeAndResource(ctx, workspaceID, integrationType, resourceType, resourceID)
	if err != nil {
		return nil, err
	}
	if integration == nil {
		return nil, ErrIntegrationNotFound
	}
	return integration, nil
}

// CountIntegrations returns the count of integrations for a workspace
func (s *Service) CountIntegrations(ctx context.Context, workspaceID uuid.UUID) (int64, error) {
	return s.Repo.CountByWorkspaceID(ctx, workspaceID)
}

// GetExpiringSoonIntegrations retrieves integrations that will expire soon
func (s *Service) GetExpiringSoonIntegrations(ctx context.Context, hours int) ([]*models.Integration, error) {
	return s.Repo.GetExpiringSoon(ctx, hours)
}
