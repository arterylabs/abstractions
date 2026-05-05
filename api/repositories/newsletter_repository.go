package repositories

import (
	"abstractions/db/models"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NewsletterRepository struct {
	db *gorm.DB
}

func NewNewsletterRepository(db *gorm.DB) *NewsletterRepository {
	return &NewsletterRepository{db: db}
}

func (r *NewsletterRepository) Create(ctx context.Context, n *models.Newsletter) error {
	return r.db.WithContext(ctx).Create(n).Error
}

func (r *NewsletterRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Newsletter, error) {
	var n models.Newsletter
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&n).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &n, nil
}

func (r *NewsletterRepository) ListByRepoID(ctx context.Context, repoID uuid.UUID) ([]*models.Newsletter, error) {
	var newsletters []*models.Newsletter
	if err := r.db.WithContext(ctx).Where("repo_id = ?", repoID).Order("created_at DESC").Find(&newsletters).Error; err != nil {
		return nil, err
	}
	return newsletters, nil
}

func (r *NewsletterRepository) Update(ctx context.Context, n *models.Newsletter) error {
	return r.db.WithContext(ctx).Save(n).Error
}

func (r *NewsletterRepository) DeleteByRepoID(ctx context.Context, repoID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("repo_id = ?", repoID).Delete(&models.Newsletter{}).Error
}
