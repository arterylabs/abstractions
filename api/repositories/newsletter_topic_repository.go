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

// NewsletterTopicWithNewsletterID extends NewsletterTopic with the newsletter ID sent for this topic and creator info.
type NewsletterTopicWithNewsletterID struct {
	models.NewsletterTopic
	NewsletterID  *uuid.UUID `gorm:"column:newsletter_id"`
	CreatorName   *string    `gorm:"column:creator_name"`
	CreatorAvatar *string    `gorm:"column:creator_avatar"`
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
	if err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("newsletter_topic_repository: get_by_id: %w", err)
	}
	return &t, nil
}

// ListByRepoID returns all topics for a repository. Pending topics come first ordered by position ASC,
// then sent topics ordered by sent_at DESC (most recently sent first). Each result is enriched with
// the newsletter ID (if sent) and the creator's name and avatar.
func (r *NewsletterTopicRepository) ListByRepoID(ctx context.Context, repoID uuid.UUID) ([]*NewsletterTopicWithNewsletterID, error) {
	var topics []*NewsletterTopicWithNewsletterID
	if err := r.db.WithContext(ctx).
		Raw(`SELECT nt.*, n.id AS newsletter_id, u.name AS creator_name, u.avatar AS creator_avatar
		     FROM newsletter_topics nt
		     LEFT JOIN newsletters n ON n.topic_id = nt.id
		     LEFT JOIN users u ON u.id = nt.created_by
		     WHERE nt.repo_id = ? AND nt.deleted_at IS NULL
		     ORDER BY
		       CASE WHEN nt.sent_at IS NULL THEN 0 ELSE 1 END ASC,
		       CASE WHEN nt.sent_at IS NULL THEN nt.position ELSE NULL END ASC,
		       nt.sent_at DESC`, repoID).
		Scan(&topics).Error; err != nil {
		return nil, fmt.Errorf("newsletter_topic_repository: list_by_repo_id: %w", err)
	}
	return topics, nil
}

// GetNextForRepo returns the oldest pending (unsent) topic for the given repository, ordered by position then created_at.
func (r *NewsletterTopicRepository) GetNextForRepo(ctx context.Context, repoID uuid.UUID) (*models.NewsletterTopic, error) {
	var t models.NewsletterTopic
	if err := r.db.WithContext(ctx).
		Where("repo_id = ? AND sent_at IS NULL AND deleted_at IS NULL", repoID).
		Order("position ASC, created_at ASC, id ASC").
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
	topic, err := r.GetByID(ctx, topicID)
	if err != nil || topic == nil {
		return err
	}
	now := time.Now()
	topic.SentAt = &now
	topic.UpdatedAt = now
	return r.db.WithContext(ctx).Save(topic).Error
}

// TitleExistsForRepo returns true when a topic with the same title (case-insensitive) already exists in the queue.
func (r *NewsletterTopicRepository) TitleExistsForRepo(ctx context.Context, repoID uuid.UUID, title string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.NewsletterTopic{}).
		Where("repo_id = ? AND LOWER(title) = LOWER(?) AND deleted_at IS NULL", repoID, title).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("newsletter_topic_repository: title_exists: %w", err)
	}
	return count > 0, nil
}

// GetMaxPendingPosition returns the highest position value among pending topics for the given repo.
// Returns -1 if no pending topics exist (so the first topic gets position 0).
func (r *NewsletterTopicRepository) GetMaxPendingPosition(ctx context.Context, repoID uuid.UUID) (int, error) {
	var maxPos *int
	if err := r.db.WithContext(ctx).
		Model(&models.NewsletterTopic{}).
		Select("MAX(position)").
		Where("repo_id = ? AND sent_at IS NULL AND deleted_at IS NULL", repoID).
		Scan(&maxPos).Error; err != nil {
		return -1, fmt.Errorf("newsletter_topic_repository: get_max_pending_position: %w", err)
	}
	if maxPos == nil {
		return -1, nil
	}
	return *maxPos, nil
}

// UpdateTopic updates the content fields of a topic.
func (r *NewsletterTopicRepository) UpdateTopic(ctx context.Context, id uuid.UUID, title, description string, relatedFiles models.JSONBRaw) error {
	topic, err := r.GetByID(ctx, id)
	if err != nil || topic == nil {
		return err
	}
	topic.Title = title
	topic.Description = description
	topic.RelatedFiles = relatedFiles
	topic.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(topic).Error
}

// ReorderTopic moves a pending topic up, down, or to the top of the queue.
// Only pending (unsent) topics can be reordered.
func (r *NewsletterTopicRepository) ReorderTopic(ctx context.Context, repoID, topicID uuid.UUID, direction string) error {
	var pendingTopics []*models.NewsletterTopic
	if err := r.db.WithContext(ctx).
		Where("repo_id = ? AND sent_at IS NULL AND deleted_at IS NULL", repoID).
		Order("position ASC, created_at ASC").
		Find(&pendingTopics).Error; err != nil {
		return fmt.Errorf("newsletter_topic_repository: reorder: list pending: %w", err)
	}

	idx := -1
	for i, t := range pendingTopics {
		if t.ID == topicID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("newsletter_topic_repository: reorder: topic not found in pending queue")
	}

	switch direction {
	case "up":
		if idx == 0 {
			return nil
		}
		pendingTopics[idx].Position, pendingTopics[idx-1].Position = pendingTopics[idx-1].Position, pendingTopics[idx].Position
		pendingTopics[idx].UpdatedAt = time.Now()
		if err := r.db.WithContext(ctx).Save(pendingTopics[idx]).Error; err != nil {
			return fmt.Errorf("newsletter_topic_repository: reorder up: %w", err)
		}
		pendingTopics[idx-1].UpdatedAt = time.Now()
		if err := r.db.WithContext(ctx).Save(pendingTopics[idx-1]).Error; err != nil {
			return fmt.Errorf("newsletter_topic_repository: reorder up swap: %w", err)
		}
	case "down":
		if idx == len(pendingTopics)-1 {
			return nil
		}
		pendingTopics[idx].Position, pendingTopics[idx+1].Position = pendingTopics[idx+1].Position, pendingTopics[idx].Position
		pendingTopics[idx].UpdatedAt = time.Now()
		if err := r.db.WithContext(ctx).Save(pendingTopics[idx]).Error; err != nil {
			return fmt.Errorf("newsletter_topic_repository: reorder down: %w", err)
		}
		pendingTopics[idx+1].UpdatedAt = time.Now()
		if err := r.db.WithContext(ctx).Save(pendingTopics[idx+1]).Error; err != nil {
			return fmt.Errorf("newsletter_topic_repository: reorder down swap: %w", err)
		}
	case "top":
		if idx == 0 {
			return nil
		}
		topPos := pendingTopics[0].Position
		for i := 0; i < idx; i++ {
			pendingTopics[i].Position++
			pendingTopics[i].UpdatedAt = time.Now()
			if err := r.db.WithContext(ctx).Save(pendingTopics[i]).Error; err != nil {
				return fmt.Errorf("newsletter_topic_repository: reorder top shift: %w", err)
			}
		}
		pendingTopics[idx].Position = topPos
		pendingTopics[idx].UpdatedAt = time.Now()
		if err := r.db.WithContext(ctx).Save(pendingTopics[idx]).Error; err != nil {
			return fmt.Errorf("newsletter_topic_repository: reorder top set: %w", err)
		}
	default:
		return fmt.Errorf("newsletter_topic_repository: reorder: unknown direction %q", direction)
	}

	return nil
}

func (r *NewsletterTopicRepository) Delete(ctx context.Context, id uuid.UUID) error {
	topic, err := r.GetByID(ctx, id)
	if err != nil || topic == nil {
		return err
	}
	now := time.Now()
	topic.DeletedAt = &now
	return r.db.WithContext(ctx).Save(topic).Error
}

func (r *NewsletterTopicRepository) DeleteByRepoID(ctx context.Context, repoID uuid.UUID) error {
	var topics []*models.NewsletterTopic
	if err := r.db.WithContext(ctx).Where("repo_id = ? AND deleted_at IS NULL", repoID).Find(&topics).Error; err != nil {
		return fmt.Errorf("newsletter_topic_repository: delete_by_repo_id: list: %w", err)
	}
	now := time.Now()
	for _, topic := range topics {
		topic.DeletedAt = &now
		if err := r.db.WithContext(ctx).Save(topic).Error; err != nil {
			return fmt.Errorf("newsletter_topic_repository: delete_by_repo_id: save: %w", err)
		}
	}
	return nil
}

func (r *NewsletterTopicRepository) PermanentlyDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.NewsletterTopic{}, "id = ?", id).Error
}
