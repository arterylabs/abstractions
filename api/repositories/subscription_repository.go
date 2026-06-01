package repositories

import (
	"abstractions/db/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionRepositoryQ interface {
	ListByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) ([]models.Subscription, error)
	ListByWorkspaceIDs(ctx context.Context, workspaceIDs []uuid.UUID) ([]models.Subscription, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Subscription, error)
	GetActiveByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) (*models.Subscription, error)
	ListActiveByWorkspaceIDs(ctx context.Context, workspaceIDs []uuid.UUID) ([]models.Subscription, error)
	Create(ctx context.Context, subscription *models.Subscription) error
	CreateBatch(ctx context.Context, subscriptions *[]models.Subscription) error
	Update(ctx context.Context, subscription *models.Subscription) error
	BulkUpdate(ctx context.Context, subscriptions *[]models.Subscription) error
	CountActive(ctx context.Context) (int, error)
	Count(ctx context.Context) (int, error)
}

type SubscriptionRepository struct {
	DB *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db}
}

func (r *SubscriptionRepository) ListByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := r.DB.WithContext(ctx).Where("workspace_id = ? AND deleted_at IS NULL", workspaceID).Find(&subscriptions).Error
	if err != nil {
		return nil, fmt.Errorf("subscription_repository: get_subscriptions_by_workspace_id: %w", err)
	}
	return subscriptions, nil
}

func (r *SubscriptionRepository) ListByWorkspaceIDs(ctx context.Context, workspaceIDs []uuid.UUID) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := r.DB.WithContext(ctx).Where("workspace_id IN ? AND deleted_at IS NULL", workspaceIDs).Find(&subscriptions).Error
	if err != nil {
		return nil, fmt.Errorf("subscription_repository: list_by_workspace_ids: %w", err)
	}
	return subscriptions, nil
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	subscription := models.Subscription{}
	err := r.DB.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&subscription).Error
	if err != nil {
		return nil, fmt.Errorf("subscription_repository: get_subscription_by_id: %w", err)
	}
	return &subscription, nil
}

func (r *SubscriptionRepository) Create(ctx context.Context, subscription *models.Subscription) error {
	return r.DB.WithContext(ctx).Create(subscription).Error
}

func (r *SubscriptionRepository) CreateBatch(ctx context.Context, subscriptions *[]models.Subscription) error {
	return r.DB.WithContext(ctx).Create(subscriptions).Error
}

func (r *SubscriptionRepository) Update(ctx context.Context, subscription *models.Subscription) error {
	subscription.UpdatedAt = time.Now()
	return r.DB.WithContext(ctx).Save(subscription).Error
}

func (r *SubscriptionRepository) GetActiveByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) (*models.Subscription, error) {
	subscription := models.Subscription{}
	err := r.DB.WithContext(ctx).Where("workspace_id = ? AND status = ? AND expires_at > now() AND deleted_at IS NULL", workspaceID, models.SubscriptionStatusActive).First(&subscription).Error
	if err != nil {
		return nil, fmt.Errorf("subscription_repository: get_active_subscription: %w", err)
	}
	return &subscription, nil
}

func (r *SubscriptionRepository) ListActiveByWorkspaceIDs(ctx context.Context, workspaceIDs []uuid.UUID) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := r.DB.WithContext(ctx).Where("workspace_id IN ? AND status = ? AND expires_at > now() AND deleted_at IS NULL", workspaceIDs, models.SubscriptionStatusActive).Find(&subscriptions).Error
	if err != nil {
		return nil, fmt.Errorf("subscription_repository: get_active_subscriptions_by_workspace_ids: %w", err)
	}
	return subscriptions, nil
}

func (r *SubscriptionRepository) BulkUpdate(ctx context.Context, subscriptions *[]models.Subscription) error {
	now := time.Now()
	for i := range *subscriptions {
		(*subscriptions)[i].UpdatedAt = now
	}
	return r.DB.WithContext(ctx).Save(subscriptions).Error
}

func (r *SubscriptionRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := r.DB.WithContext(ctx).Where("id IN ? AND deleted_at IS NULL", ids).Find(&subscriptions).Error
	if err != nil {
		return nil, fmt.Errorf("subscription_repository: get_subscriptions_by_ids: %w", err)
	}
	return subscriptions, nil
}

func (r *SubscriptionRepository) CountActive(ctx context.Context) (int, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.Subscription{}).Where("status = ? AND expires_at > now() AND deleted_at IS NULL", models.SubscriptionStatusActive).Count(&count).Error
	return int(count), err
}

func (r *SubscriptionRepository) Count(ctx context.Context) (int, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.Subscription{}).Count(&count).Error
	return int(count), err
}
