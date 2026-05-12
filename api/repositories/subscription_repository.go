package repositories

import (
	"abstractions/db/models"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionRepositoryQ interface {
	GetSubscriptionsByWorkspaceID(workspaceID uuid.UUID) ([]models.Subscription, error)
	GetSubscriptionsByWorkspaceIDs(workspaceIDs []uuid.UUID) ([]models.Subscription, error)
	GetSubscriptionByID(id uuid.UUID) (*models.Subscription, error)
	GetSubscriptionsByIDs(ids []uuid.UUID) ([]models.Subscription, error)
	GetActiveSubscription(workspaceID uuid.UUID) (*models.Subscription, error)
	CreateSubscription(subscription *models.Subscription) error
	UpdateSubscription(subscription *models.Subscription) error
	BulkUpdate(subscriptions *[]models.Subscription) error
	CountActiveSubscriptions() (int, error)
	Count() (int, error)
}

type SubscriptionRepository struct {
	DB *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db}
}

func (subscriptionRepository *SubscriptionRepository) GetSubscriptionsByWorkspaceID(workspaceID uuid.UUID) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := subscriptionRepository.DB.Where("workspace_id = ?", workspaceID).Find(&subscriptions).Error
	if err != nil {
		return nil, fmt.Errorf("subscription_repository: get_subscriptions_by_workspace_id: %w", err)
	}
	return subscriptions, nil
}

func (subscriptionRepository *SubscriptionRepository) GetSubscriptionByID(id uuid.UUID) (*models.Subscription, error) {
	subscription := models.Subscription{}
	err := subscriptionRepository.DB.Where("id = ?", id).First(&subscription).Error
	if err != nil {
		return nil, fmt.Errorf("subscription_repository: get_subscription_by_id: %w", err)
	}
	return &subscription, nil
}

func (subscriptionRepository *SubscriptionRepository) CreateSubscription(subscription *models.Subscription) error {
	return subscriptionRepository.DB.Create(subscription).Error
}

func (subscriptionRepository *SubscriptionRepository) CreateBulkSubscriptions(subscriptions *[]models.Subscription) error {
	return subscriptionRepository.DB.Create(subscriptions).Error
}

func (subscriptionRepository *SubscriptionRepository) UpdateSubscription(subscription *models.Subscription) error {
	return subscriptionRepository.DB.Save(subscription).Error
}

func (subscriptionRepository *SubscriptionRepository) GetActiveSubscription(workspaceID uuid.UUID) (*models.Subscription, error) {
	subscription := models.Subscription{}
	err := subscriptionRepository.DB.Where("workspace_id = ? AND status = ? AND expires_at > now()", workspaceID, models.SubscriptionStatusActive).First(&subscription).Error
	if err != nil {
		return nil, fmt.Errorf("subscription_repository: get_active_subscription: %w", err)
	}
	return &subscription, nil
}

func (subscriptionRepository *SubscriptionRepository) GetActiveSubscriptionsByWorkspaceIDs(workspaceIDs []uuid.UUID) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := subscriptionRepository.DB.Where("workspace_id IN (?) AND status = ? AND expires_at > now()", workspaceIDs, models.SubscriptionStatusActive).Find(&subscriptions).Error
	if err != nil {
		return nil, fmt.Errorf("subscription_repository: get_active_subscriptions_by_workspace_ids: %w", err)
	}
	return subscriptions, nil
}

func (subscriptionRepository *SubscriptionRepository) BulkUpdate(subscriptions *[]models.Subscription) error {
	return subscriptionRepository.DB.Save(subscriptions).Error
}

func (subscriptionRepository *SubscriptionRepository) GetSubscriptionsByIDs(ids []uuid.UUID) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := subscriptionRepository.DB.Where("id IN (?)", ids).Find(&subscriptions).Error
	if err != nil {
		return nil, fmt.Errorf("subscription_repository: get_subscriptions_by_ids: %w", err)
	}
	return subscriptions, nil
}

func (subscriptionRepository *SubscriptionRepository) CountActiveSubscriptions() (int, error) {
	var count int64
	err := subscriptionRepository.DB.Model(&models.Subscription{}).Where("status = ? AND expires_at > now()", models.SubscriptionStatusActive).Count(&count).Error
	return int(count), err
}

func (subscriptionRepository *SubscriptionRepository) Count() (int, error) {
	var count int64
	err := subscriptionRepository.DB.Model(&models.Subscription{}).Count(&count).Error
	return int(count), err
}
