package repositories

import (
	"abstractions/db/models"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrTopicAlreadySent = errors.New("topic has already been sent")

type NewsletterTopicRepository struct {
	db *gorm.DB
}

func NewNewsletterTopicRepository(db *gorm.DB) *NewsletterTopicRepository {
	return &NewsletterTopicRepository{db: db}
}

func (r *NewsletterTopicRepository) Create(ctx context.Context, t *models.NewsletterTopic) error {
	return r.db.WithContext(ctx).Create(t).Error
}

// ListByRepoID returns all topics for a repo ordered by pending first, then sent, oldest first.
func (r *NewsletterTopicRepository) ListByRepoID(ctx context.Context, repoID uuid.UUID) ([]*models.NewsletterTopic, error) {
	var topics []*models.NewsletterTopic
	if err := r.db.WithContext(ctx).
		Where("repo_id = ? AND deleted_at IS NULL", repoID).
		Order("CASE WHEN sent_at IS NULL THEN 0 ELSE 1 END ASC").
		Order("created_at ASC").
		Find(&topics).Error; err != nil {
		return nil, err
	}
	return topics, nil
}

// GetNextForRepo returns the next unsent topic for a repo, oldest first.
func (r *NewsletterTopicRepository) GetNextForRepo(ctx context.Context, repoID uuid.UUID) (*models.NewsletterTopic, error) {
	var topic models.NewsletterTopic
	if err := r.db.WithContext(ctx).
		Where("repo_id = ? AND sent_at IS NULL AND deleted_at IS NULL", repoID).
		Order("created_at ASC").
		First(&topic).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &topic, nil
}

// MarkSent sets sent_at for the given topic.
func (r *NewsletterTopicRepository) MarkSent(ctx context.Context, topicID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.NewsletterTopic{}).
		Where("id = ? AND deleted_at IS NULL", topicID).
		Update("sent_at", now).Error
}

// Delete hard-deletes a topic by ID and rejects deleting sent topics.
func (r *NewsletterTopicRepository) Delete(ctx context.Context, topicID uuid.UUID) error {
	var topic models.NewsletterTopic
	if err := r.db.WithContext(ctx).
		Select("id", "sent_at").
		Where("id = ? AND deleted_at IS NULL", topicID).
		First(&topic).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if topic.SentAt != nil {
		return ErrTopicAlreadySent
	}
	return r.db.WithContext(ctx).Unscoped().Delete(&models.NewsletterTopic{}, "id = ?", topicID).Error
}

// CreateBatch inserts multiple topics and skips duplicate titles for the same repo.
func (r *NewsletterTopicRepository) CreateBatch(ctx context.Context, topics []*models.NewsletterTopic) error {
	if len(topics) == 0 {
		return nil
	}

	titlesByRepo := make(map[uuid.UUID][]string)
	seenInBatch := make(map[string]struct{})
	for _, topic := range topics {
		if topic == nil || topic.Title == "" {
			continue
		}
		key := topic.RepoID.String() + ":" + topic.Title
		if _, exists := seenInBatch[key]; exists {
			continue
		}
		seenInBatch[key] = struct{}{}
		titlesByRepo[topic.RepoID] = append(titlesByRepo[topic.RepoID], topic.Title)
	}

	existingKeys := make(map[string]struct{})
	for repoID, titles := range titlesByRepo {
		if len(titles) == 0 {
			continue
		}
		var existing []models.NewsletterTopic
		if err := r.db.WithContext(ctx).
			Select("repo_id", "title").
			Where("repo_id = ? AND title IN ? AND deleted_at IS NULL", repoID, titles).
			Find(&existing).Error; err != nil {
			return err
		}
		for _, topic := range existing {
			existingKeys[topic.RepoID.String()+":"+topic.Title] = struct{}{}
		}
	}

	toCreate := make([]*models.NewsletterTopic, 0, len(topics))
	addedKeys := make(map[string]struct{})
	for _, topic := range topics {
		if topic == nil || topic.Title == "" {
			continue
		}
		key := topic.RepoID.String() + ":" + topic.Title
		if _, exists := existingKeys[key]; exists {
			continue
		}
		if _, exists := addedKeys[key]; exists {
			continue
		}
		addedKeys[key] = struct{}{}
		toCreate = append(toCreate, topic)
	}

	if len(toCreate) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).CreateInBatches(toCreate, 100).Error
}

// GetByID returns a single topic by ID.
func (r *NewsletterTopicRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.NewsletterTopic, error) {
	var topic models.NewsletterTopic
	if err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&topic).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &topic, nil
}

// CountUnsentByRepoID returns the number of unsent topics for a repo.
func (r *NewsletterTopicRepository) CountUnsentByRepoID(ctx context.Context, repoID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.NewsletterTopic{}).
		Where("repo_id = ? AND sent_at IS NULL AND deleted_at IS NULL", repoID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
