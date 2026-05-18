package repositories

import (
	"abstractions/db/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WorkspaceConfigurationRepository struct {
	db *gorm.DB
}

func NewWorkspaceConfigurationRepository(db *gorm.DB) *WorkspaceConfigurationRepository {
	return &WorkspaceConfigurationRepository{db: db}
}

func (r *WorkspaceConfigurationRepository) Upsert(ctx context.Context, cfg *models.WorkspaceConfiguration) error {
	cfg.UpdatedAt = time.Now()
	if cfg.ID == (uuid.UUID{}) {
		cfg.ID = uuid.New()
	}
	cfg.CreatedAt = time.Now()
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "workspace_id"}, {Name: "key"}},
			DoUpdates: clause.AssignmentColumns([]string{"value", "encrypted", "updated_at", "updated_by"}),
		}).
		Create(cfg).Error
}

func (r *WorkspaceConfigurationRepository) GetByWorkspaceAndKey(ctx context.Context, workspaceID uuid.UUID, key string) (*models.WorkspaceConfiguration, error) {
	var cfg models.WorkspaceConfiguration
	if err := r.db.WithContext(ctx).
		Where("workspace_id = ? AND key = ?", workspaceID, key).
		First(&cfg).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("workspace_configuration_repository: get_by_workspace_and_key: %w", err)
	}
	return &cfg, nil
}

func (r *WorkspaceConfigurationRepository) GetAllByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]*models.WorkspaceConfiguration, error) {
	var cfgs []*models.WorkspaceConfiguration
	if err := r.db.WithContext(ctx).
		Where("workspace_id = ?", workspaceID).
		Find(&cfgs).Error; err != nil {
		return nil, fmt.Errorf("workspace_configuration_repository: get_all_by_workspace: %w", err)
	}
	return cfgs, nil
}

func (r *WorkspaceConfigurationRepository) Delete(ctx context.Context, workspaceID uuid.UUID, key string) error {
	return r.db.WithContext(ctx).
		Where("workspace_id = ? AND key = ?", workspaceID, key).
		Delete(&models.WorkspaceConfiguration{}).Error
}
