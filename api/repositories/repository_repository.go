package repositories

import (
	"abstractions/db/models"
	"context"
	"errors"

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
		return nil, err
	}
	return &repo, nil
}

func (r *RepositoryRepository) GetByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) ([]*models.Repository, error) {
	var repos []*models.Repository
	if err := r.db.WithContext(ctx).Where("workspace_id = ? AND deleted_at IS NULL", workspaceID).Find(&repos).Error; err != nil {
		return nil, err
	}
	return repos, nil
}

func (r *RepositoryRepository) GetByWorkspaceAndServiceID(ctx context.Context, workspaceID uuid.UUID, serviceType, serviceID string) (*models.Repository, error) {
	var repo models.Repository
	if err := r.db.WithContext(ctx).Where("workspace_id = ? AND service_type = ? AND service_id = ? AND deleted_at IS NULL", workspaceID, serviceType, serviceID).First(&repo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &repo, nil
}

func (r *RepositoryRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.RepositoryStatus) error {
	return r.db.WithContext(ctx).Model(&models.Repository{}).Where("id = ?", id).Update("status", status).Error
}

func (r *RepositoryRepository) Update(ctx context.Context, repo *models.Repository) error {
	return r.db.WithContext(ctx).Save(repo).Error
}

func (r *RepositoryRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Repository{}).Where("id = ?", id).Update("deleted_at", "NOW()").Error
}

func (r *RepositoryRepository) GetDueForNewsletter(ctx context.Context, sendDay, sendHour int) ([]*models.Repository, error) {
	var repos []*models.Repository
	err := r.db.WithContext(ctx).
		Where("send_day = ? AND send_hour = ? AND status = ? AND deleted_at IS NULL", sendDay, sendHour, models.RepositoryStatusIndexed).
		Find(&repos).Error
	return repos, err
}
