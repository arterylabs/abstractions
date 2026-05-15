package subscription_service

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrSubscriptionNotFound      = errors.New("subscription not found")
	ErrSubscriptionExpired       = errors.New("subscription has expired")
	ErrInvalidSubscriptionStatus = errors.New("invalid subscription status")
	ErrNoActiveSubscription      = errors.New("no active subscription found")
)

type Service struct {
	DB   *gorm.DB
	Repo *repositories.SubscriptionRepository
}

func NewSubscriptionService(db *gorm.DB) *Service {
	return &Service{Repo: repositories.NewSubscriptionRepository(db), DB: db}
}

type CreateSubscriptionParams struct {
	WorkspaceID uuid.UUID
	PlanID      models.PlanID
	Status      models.SubscriptionStatus
	ExpiresAt   time.Time
	Limit       int
}

type UpdateSubscriptionParams struct {
	PlanID    *models.PlanID
	Status    *models.SubscriptionStatus
	ExpiresAt *time.Time
	Limit     *int
}

// Basic CRUD operations
func (s *Service) GetSubscriptionsByWorkspaceID(workspaceID uuid.UUID) ([]models.Subscription, error) {
	return s.Repo.GetSubscriptionsByWorkspaceID(workspaceID)
}

func (s *Service) GetActiveSubscriptionByWorkspaceID(workspaceID uuid.UUID) (*models.Subscription, error) {
	subscriptions, err := s.Repo.GetSubscriptionsByWorkspaceID(workspaceID)
	if err != nil {
		return nil, fmt.Errorf("subscription_service: get_active_subscription_by_workspace_id: %w", err)
	}

	for _, subscription := range subscriptions {
		if subscription.Status == models.SubscriptionStatusActive && subscription.ExpiresAt.After(time.Now()) {
			return &subscription, nil
		}
	}

	return nil, nil
}

func (s *Service) UpdateSubscription(subscription *models.Subscription) error {
	return s.Repo.UpdateSubscription(subscription)
}

func (s *Service) CountActiveSubscriptions() (int, error) {
	return s.Repo.CountActiveSubscriptions()
}

func (s *Service) Count() (int, error) {
	return s.Repo.Count()
}
