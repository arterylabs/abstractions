package subscription_service

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Complex business operations for subscriptions

func (s *Service) CreateBulkSubscription(ctx context.Context, workspaceIDs []uuid.UUID, params CreateSubscriptionParams) (*[]models.Subscription, error) {
	subscriptionsToCreate := []models.Subscription{}
	repo := repositories.NewSubscriptionRepository(s.DB.WithContext(ctx))

	for _, workspaceID := range workspaceIDs {
		subscription := models.Subscription{
			ID:          uuid.New(),
			WorkspaceID: workspaceID,
			PlanID:      params.PlanID,
			Status:      params.Status,
			ExpiresAt:   params.ExpiresAt,
			Limit:       params.Limit,
		}
		subscriptionsToCreate = append(subscriptionsToCreate, subscription)
	}
	err := repo.CreateBatch(ctx, &subscriptionsToCreate)
	if err != nil {
		return nil, fmt.Errorf("subscription_service: create_bulk_subscription: %w", err)
	}
	return &subscriptionsToCreate, nil
}

func (s *Service) DisableSubscription(ctx context.Context, id uuid.UUID, reason string) error {
	repo := repositories.NewSubscriptionRepository(s.DB.WithContext(ctx))
	subscription, err := repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("subscription_service: disable_subscription: get_subscription: %w", err)
	}

	currentTime := time.Now()
	subscription.DisabledAt = &currentTime
	subscription.DisabledReason = &reason
	subscription.Status = models.SubscriptionStatusInactive
	return repo.Update(ctx, subscription)
}

func (s *Service) EnableSubscription(ctx context.Context, id uuid.UUID) error {
	repo := repositories.NewSubscriptionRepository(s.DB.WithContext(ctx))
	subscription, err := repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("subscription_service: enable_subscription: get_subscription: %w", err)
	}

	subscription.DisabledAt = nil
	subscription.DisabledReason = nil
	subscription.Status = models.SubscriptionStatusActive
	return repo.Update(ctx, subscription)
}

func (s *Service) UpdateSubscriptionsForWorkspaces(ctx context.Context, workspaceIDs []uuid.UUID, params UpdateSubscriptionParams) error {
	repo := repositories.NewSubscriptionRepository(s.DB.WithContext(ctx))
	subscriptions, err := repo.ListActiveByWorkspaceIDs(ctx, workspaceIDs)
	if err != nil {
		return fmt.Errorf("subscription_service: update_subscriptions_for_workspaces: get_subscriptions: %w", err)
	}

	for idx := range subscriptions {
		if params.PlanID != nil {
			subscriptions[idx].PlanID = *params.PlanID
		}
		if params.Status != nil {
			subscriptions[idx].Status = *params.Status
		}
		if params.ExpiresAt != nil {
			subscriptions[idx].ExpiresAt = *params.ExpiresAt
		}
		if params.Limit != nil {
			subscriptions[idx].Limit = *params.Limit
		}
	}

	return repo.BulkUpdate(ctx, &subscriptions)
}

func (s *Service) DisableSubscriptions(ctx context.Context, ids []uuid.UUID, reason string) error {
	repo := repositories.NewSubscriptionRepository(s.DB.WithContext(ctx))
	subscriptions, err := repo.GetByIDs(ctx, ids)
	if err != nil {
		return fmt.Errorf("subscription_service: disable_subscriptions: get_subscriptions: %w", err)
	}

	currentTime := time.Now()
	for idx := range subscriptions {
		subscriptions[idx].DisabledAt = &currentTime
		subscriptions[idx].DisabledReason = &reason
		subscriptions[idx].Status = models.SubscriptionStatusInactive
	}

	return repo.BulkUpdate(ctx, &subscriptions)
}
