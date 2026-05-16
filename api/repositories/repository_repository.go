package repositories

import (
	"abstractions/db/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RepositoryRepository struct {
	db *gorm.DB
}

func NewRepositoryRepository(db *gorm.DB) *RepositoryRepository {
	return &RepositoryRepository{db: db}
}

func (r *RepositoryRepository) Create(ctx context.Context, repo *models.Repository) error {
	return r.db.WithContext(ctx).Create(repo).Error
}

func (r *RepositoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Repository, error) {
	var repo models.Repository
	if err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&repo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("repository_repository: get_by_id: %w", err)
	}
	return &repo, nil
}

func (r *RepositoryRepository) GetByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) ([]*models.Repository, error) {
	var repos []*models.Repository
	if err := r.db.WithContext(ctx).Where("workspace_id = ? AND deleted_at IS NULL", workspaceID).Find(&repos).Error; err != nil {
		return nil, fmt.Errorf("repository_repository: get_by_workspace_id: %w", err)
	}
	return repos, nil
}

func (r *RepositoryRepository) GetByWorkspaceAndServiceID(ctx context.Context, workspaceID uuid.UUID, serviceType, serviceID string) (*models.Repository, error) {
	var repo models.Repository
	if err := r.db.WithContext(ctx).
		Where("workspace_id = ? AND service_type = ? AND service_id = ? AND deleted_at IS NULL", workspaceID, serviceType, serviceID).
		First(&repo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("repository_repository: get_by_workspace_and_service_id: %w", err)
	}
	return &repo, nil
}

func (r *RepositoryRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.RepositoryStatus) error {
	return r.db.WithContext(ctx).Model(&models.Repository{}).Where("id = ?", id).Update("status", status).Error
}

// GetAllIndexed returns all non-deleted repositories that have been successfully indexed.
func (r *RepositoryRepository) GetAllIndexed(ctx context.Context) ([]*models.Repository, error) {
	var repos []*models.Repository
	err := r.db.WithContext(ctx).
		Where("status = ? AND deleted_at IS NULL", models.RepositoryStatusIndexed).
		Find(&repos).Error
	return repos, err
}

func (r *RepositoryRepository) Update(ctx context.Context, repo *models.Repository) error {
	return r.db.WithContext(ctx).Save(repo).Error
}

func (r *RepositoryRepository) SoftDelete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).Model(&models.Repository{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"deleted_at": now, "deleted_by": deletedBy}).Error
}

func (r *RepositoryRepository) GetDueForNewsletter(ctx context.Context, sendDay, sendHour int) ([]*models.Repository, error) {
	var repos []*models.Repository
	err := r.db.WithContext(ctx).
		Where("send_day = ? AND send_hour = ? AND status = ? AND newsletter_enabled = TRUE AND deleted_at IS NULL", sendDay, sendHour, models.RepositoryStatusIndexed).
		Find(&repos).Error
	return repos, err
}
