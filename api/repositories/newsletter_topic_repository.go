package repositories

import (
	"abstractions/db/models"
	"context"

	"gorm.io/gorm"
)

type NewsletterTopicRepository struct {
	db *gorm.DB
}

func NewNewsletterTopicRepository(db *gorm.DB) *NewsletterTopicRepository {
	return &NewsletterTopicRepository{db: db}
}

func (r *NewsletterTopicRepository) Create(ctx context.Context, t *models.NewsletterTopic) error {
	return r.db.WithContext(ctx).Create(t).Error
}
