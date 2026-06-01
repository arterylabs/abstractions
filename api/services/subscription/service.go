package subscription_service

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"context"
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
func (s *Service) GetSubscriptionsByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) ([]models.Subscription, error) {
	return repositories.NewSubscriptionRepository(s.DB.WithContext(ctx)).ListByWorkspaceID(ctx, workspaceID)
}

func (s *Service) GetActiveSubscriptionByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) (*models.Subscription, error) {
	subscriptions, err := repositories.NewSubscriptionRepository(s.DB.WithContext(ctx)).ListByWorkspaceID(ctx, workspaceID)
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

func (s *Service) UpdateSubscription(ctx context.Context, subscription *models.Subscription) error {
	return repositories.NewSubscriptionRepository(s.DB.WithContext(ctx)).Update(ctx, subscription)
}

func (s *Service) CountActiveSubscriptions(ctx context.Context) (int, error) {
	return repositories.NewSubscriptionRepository(s.DB.WithContext(ctx)).CountActive(ctx)
}

func (s *Service) Count(ctx context.Context) (int, error) {
	return repositories.NewSubscriptionRepository(s.DB.WithContext(ctx)).Count(ctx)
}
