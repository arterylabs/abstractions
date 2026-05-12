package repositories

import (
	"abstractions/db/models"
	"abstractions/services/integrations/enums"
	"abstractions/utils/pagination_utils"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IntegrationRepository struct {
	db *gorm.DB
}

func NewIntegrationRepository(db *gorm.DB) *IntegrationRepository {
	return &IntegrationRepository{
		db: db,
	}
}

func (r *IntegrationRepository) Create(ctx context.Context, integration *models.Integration) error {
	return r.db.WithContext(ctx).Create(integration).Error
}

func (r *IntegrationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Integration, error) {
	var integration models.Integration
	if err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&integration).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &integration, nil
}

func (r *IntegrationRepository) GetByIDIncludingDeleted(ctx context.Context, id uuid.UUID) (*models.Integration, error) {
	var integration models.Integration
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&integration).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &integration, nil
}

func (r *IntegrationRepository) GetByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) ([]*models.Integration, error) {
	var integrations []*models.Integration
	if err := r.db.WithContext(ctx).Where("workspace_id = ? AND deleted_at IS NULL", workspaceID).Find(&integrations).Error; err != nil {
		return nil, err
	}
	return integrations, nil
}

func (r *IntegrationRepository) GetMultipleByWorkspaceIDAndType(ctx context.Context, workspaceID uuid.UUID, integrationType enums.IntegrationTypeEnum) ([]*models.Integration, error) {
	var integrations []*models.Integration
	if err := r.db.WithContext(ctx).Where("workspace_id = ? AND type = ? AND deleted_at IS NULL", workspaceID, integrationType).Find(&integrations).Error; err != nil {
		return nil, err
	}
	return integrations, nil
}

func (r *IntegrationRepository) GetByWorkspaceIDAndType(ctx context.Context, workspaceID uuid.UUID, integrationType enums.IntegrationTypeEnum) (*models.Integration, error) {
	var integration models.Integration
	if err := r.db.WithContext(ctx).Where("workspace_id = ? AND type = ? AND deleted_at IS NULL", workspaceID, integrationType).First(&integration).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &integration, nil
}

func (r *IntegrationRepository) GetPaginated(ctx context.Context, workspaceID uuid.UUID, params pagination_utils.PaginationParams) ([]*models.Integration, int, error) {
	var integrations []*models.Integration
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.Integration{}).Where("workspace_id = ? AND deleted_at IS NULL", workspaceID)

	// Apply search if provided
	if params.Query != "" {
		query = query.Where("type ILIKE ? OR service_title ILIKE ?", "%"+params.Query+"%", "%"+params.Query+"%")
	}

	// Get total count before pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Apply ordering
	if params.Sort != "" {
		sortOrder := pagination_utils.ParseSortOrder(params.SortOrder)
		switch params.Sort {
		case "type":
			query = query.Order("type " + sortOrder.String())
		case "service_title":
			query = query.Order("service_title " + sortOrder.String())
		case "created_at":
			query = query.Order("created_at " + sortOrder.String())
		default:
			return nil, 0, errors.New("invalid sort field")
		}
	} else {
		// Default ordering by created_at if no sort field is provided
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if err := query.Offset(params.GetOffset()).Limit(params.GetLimit()).Find(&integrations).Error; err != nil {
		return nil, 0, err
	}

	return integrations, int(totalCount), nil
}

func (r *IntegrationRepository) Update(ctx context.Context, integration *models.Integration) error {
	integration.UpdatedAt = time.Now().UTC()
	return r.db.WithContext(ctx).Save(integration).Error
}

func (r *IntegrationRepository) UpdateConfiguration(ctx context.Context, id uuid.UUID, configuration models.JSONBField, updatedBy uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Integration{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{
			"configuration": configuration,
			"updated_at":    time.Now().UTC(),
			"updated_by":    updatedBy,
		}).Error
}

func (r *IntegrationRepository) UpdateTokens(ctx context.Context, id uuid.UUID, token string, refreshToken *string, expiresAt *time.Time, updatedBy uuid.UUID) error {
	updates := map[string]interface{}{
		"token":      token,
		"updated_at": time.Now().UTC(),
		"updated_by": updatedBy,
	}

	if refreshToken != nil {
		updates["refresh_token"] = *refreshToken
	}

	if expiresAt != nil {
		updates["token_expires_at"] = *expiresAt
	}

	return r.db.WithContext(ctx).Model(&models.Integration{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).Error
}

func (r *IntegrationRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).Model(&models.Integration{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{
			"deleted_at":       &now,
			"deleted_by":       &deletedBy,
			"token":            "",
			"refresh_token":    nil,
			"token_expires_at": nil,
			"configuration":    nil,
			"updated_at":       now,
			"updated_by":       deletedBy,
		}).Error
}

func (r *IntegrationRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Integration{}).Error
}

func (r *IntegrationRepository) ExistsByWorkspaceIDAndType(ctx context.Context, workspaceID uuid.UUID, integrationType enums.IntegrationTypeEnum) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Integration{}).
		Where("workspace_id = ? AND type = ? AND deleted_at IS NULL", workspaceID, integrationType).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *IntegrationRepository) CountByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Integration{}).
		Where("workspace_id = ? AND deleted_at IS NULL", workspaceID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *IntegrationRepository) GetExpiringSoon(ctx context.Context, hours int) ([]*models.Integration, error) {
	var integrations []*models.Integration
	threshold := time.Now().Add(time.Duration(hours) * time.Hour)

	if err := r.db.WithContext(ctx).
		Where("token_expires_at IS NOT NULL AND token_expires_at <= ? AND deleted_at IS NULL", threshold).
		Find(&integrations).Error; err != nil {
		return nil, err
	}
	return integrations, nil
}

func (r *IntegrationRepository) GetByConfigurationField(ctx context.Context, workspaceID uuid.UUID, jsonPath string, value interface{}) ([]*models.Integration, error) {
	var integrations []*models.Integration
	if err := r.db.WithContext(ctx).
		Where("workspace_id = ? AND deleted_at IS NULL AND configuration ->> ? = ?", workspaceID, jsonPath, value).
		Find(&integrations).Error; err != nil {
		return nil, err
	}
	return integrations, nil
}

// ExistsByWorkspaceIDTypeAndResource checks if an active integration exists for the given
// workspace, type, resource_type, and resource_id combination.
func (r *IntegrationRepository) ExistsByWorkspaceIDTypeAndResource(ctx context.Context, workspaceID uuid.UUID, integrationType enums.IntegrationTypeEnum, resourceType *string, resourceID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.Integration{}).
		Where("workspace_id = ? AND type = ? AND deleted_at IS NULL", workspaceID, integrationType)
	if resourceType != nil {
		query = query.Where("resource_type = ?", *resourceType)
	} else {
		query = query.Where("resource_type IS NULL")
	}
	if resourceID != nil {
		query = query.Where("resource_id = ?", *resourceID)
	} else {
		query = query.Where("resource_id IS NULL")
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetPaginatedByResource returns paginated integrations scoped to a specific resource.
func (r *IntegrationRepository) GetPaginatedByResource(ctx context.Context, workspaceID uuid.UUID, resourceType string, resourceID uuid.UUID, params pagination_utils.PaginationParams) ([]*models.Integration, int, error) {
	var integrations []*models.Integration
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.Integration{}).
		Where("workspace_id = ? AND resource_type = ? AND resource_id = ? AND deleted_at IS NULL", workspaceID, resourceType, resourceID)

	if params.Query != "" {
		query = query.Where("type ILIKE ? OR service_title ILIKE ?", "%"+params.Query+"%", "%"+params.Query+"%")
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	query = query.Order("created_at DESC")
	if err := query.Offset(params.GetOffset()).Limit(params.GetLimit()).Find(&integrations).Error; err != nil {
		return nil, 0, err
	}

	return integrations, int(totalCount), nil
}
