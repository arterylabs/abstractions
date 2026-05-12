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

type NewsletterTopicRepository struct {
	db *gorm.DB
}

// NewsletterTopicWithNewsletterID extends NewsletterTopic with the ID of the newsletter sent for this topic (if any).
type NewsletterTopicWithNewsletterID struct {
	models.NewsletterTopic
	NewsletterID *uuid.UUID `gorm:"column:newsletter_id"`
}

func NewNewsletterTopicRepository(db *gorm.DB) *NewsletterTopicRepository {
	return &NewsletterTopicRepository{db: db}
}

func (r *NewsletterTopicRepository) Create(ctx context.Context, t *models.NewsletterTopic) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *NewsletterTopicRepository) CreateBatch(ctx context.Context, topics []*models.NewsletterTopic) error {
	if len(topics) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&topics).Error
}

func (r *NewsletterTopicRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.NewsletterTopic, error) {
	var t models.NewsletterTopic
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("newsletter_topic_repository: get_by_id: %w", err)
	}
	return &t, nil
}

// ListByRepoID returns all topics for a repository, pending first then sent, ordered by created_at ASC within each group.
// Each result is enriched with the ID of the newsletter that was generated for the topic (if any).
func (r *NewsletterTopicRepository) ListByRepoID(ctx context.Context, repoID uuid.UUID) ([]*NewsletterTopicWithNewsletterID, error) {
	var topics []*NewsletterTopicWithNewsletterID
	if err := r.db.WithContext(ctx).
		Raw(`SELECT nt.*, n.id AS newsletter_id
		     FROM newsletter_topics nt
		     LEFT JOIN newsletters n ON n.topic_id = nt.id
		     WHERE nt.repo_id = ?
		     ORDER BY nt.sent_at IS NOT NULL, nt.created_at ASC`, repoID).
		Scan(&topics).Error; err != nil {
		return nil, fmt.Errorf("newsletter_topic_repository: list_by_repo_id: %w", err)
	}
	return topics, nil
}

// GetNextForRepo returns the oldest pending (unsent) topic for the given repository.
func (r *NewsletterTopicRepository) GetNextForRepo(ctx context.Context, repoID uuid.UUID) (*models.NewsletterTopic, error) {
	var t models.NewsletterTopic
	if err := r.db.WithContext(ctx).
		Where("repo_id = ? AND sent_at IS NULL", repoID).
		Order("created_at ASC, id ASC").
		First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("newsletter_topic_repository: get_next_for_repo: %w", err)
	}
	return &t, nil
}

// MarkSent stamps sent_at on the given topic.
func (r *NewsletterTopicRepository) MarkSent(ctx context.Context, topicID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.NewsletterTopic{}).
		Where("id = ?", topicID).
		Update("sent_at", now).Error
}

// TitleExistsForRepo returns true when a topic with the same title (case-insensitive) already exists in the queue.
func (r *NewsletterTopicRepository) TitleExistsForRepo(ctx context.Context, repoID uuid.UUID, title string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.NewsletterTopic{}).
		Where("repo_id = ? AND LOWER(title) = LOWER(?)", repoID, title).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("newsletter_topic_repository: title_exists: %w", err)
	}
	return count > 0, nil
}

func (r *NewsletterTopicRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.NewsletterTopic{}).Error
}

func (r *NewsletterTopicRepository) DeleteByRepoID(ctx context.Context, repoID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("repo_id = ?", repoID).Delete(&models.NewsletterTopic{}).Error
}
