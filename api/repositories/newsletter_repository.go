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
	if err := r.db.WithContext(ctx).Preload("Topic").Where("id = ? AND deleted_at IS NULL", id).First(&n).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("newsletter_repository: get_by_id: %w", err)
	}
	return &n, nil
}

func (r *NewsletterRepository) ListByRepoID(ctx context.Context, repoID uuid.UUID) ([]*models.Newsletter, error) {
	var newsletters []*models.Newsletter
	if err := r.db.WithContext(ctx).Preload("Topic").Where("repo_id = ? AND deleted_at IS NULL", repoID).Order("created_at DESC").Find(&newsletters).Error; err != nil {
		return nil, fmt.Errorf("newsletter_repository: list_by_repo_id: %w", err)
	}
	return newsletters, nil
}

func (r *NewsletterRepository) Update(ctx context.Context, n *models.Newsletter) error {
	n.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(n).Error
}

func (r *NewsletterRepository) DeleteByRepoID(ctx context.Context, repoID uuid.UUID) error {
	newsletters, err := r.ListByRepoID(ctx, repoID)
	if err != nil {
		return err
	}
	now := time.Now()
	for _, newsletter := range newsletters {
		newsletter.DeletedAt = &now
		if err := r.Update(ctx, newsletter); err != nil {
			return err
		}
	}
	return nil
}

func (r *NewsletterRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy *uuid.UUID) error {
	n, err := r.GetByID(ctx, id)
	if err != nil || n == nil {
		return err
	}
	now := time.Now()
	n.DeletedAt = &now
	n.DeletedBy = deletedBy
	return r.Update(ctx, n)
}

func (r *NewsletterRepository) PermanentlyDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Newsletter{}, "id = ?", id).Error
}

func (r *NewsletterRepository) Count(ctx context.Context) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Newsletter{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("newsletter_repository: count: %w", err)
	}
	return int(count), nil
}
