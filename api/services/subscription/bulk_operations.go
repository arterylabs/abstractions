package subscription_service

import (
	"abstractions/db/models"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Complex business operations for subscriptions

func (s *Service) CreateBulkSubscription(workspaceIDs []uuid.UUID, params CreateSubscriptionParams) (*[]models.Subscription, error) {
	subscriptionsToCreate := []models.Subscription{}

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
	err := s.Repo.CreateBulkSubscriptions(&subscriptionsToCreate)
	if err != nil {
		return nil, fmt.Errorf("subscription_service: create_bulk_subscription: %w", err)
	}
	return &subscriptionsToCreate, nil
}

func (s *Service) DisableSubscription(id uuid.UUID, reason string) error {
	subscription, err := s.Repo.GetSubscriptionByID(id)
	if err != nil {
		return fmt.Errorf("subscription_service: disable_subscription: get_subscription: %w", err)
	}

	currentTime := time.Now()
	subscription.DisabledAt = &currentTime
	subscription.DisabledReason = &reason
	subscription.Status = models.SubscriptionStatusInactive
	return s.Repo.UpdateSubscription(subscription)
}

func (s *Service) EnableSubscription(id uuid.UUID) error {
	subscription, err := s.Repo.GetSubscriptionByID(id)
	if err != nil {
		return fmt.Errorf("subscription_service: enable_subscription: get_subscription: %w", err)
	}

	subscription.DisabledAt = nil
	subscription.DisabledReason = nil
	subscription.Status = models.SubscriptionStatusActive
	return s.Repo.UpdateSubscription(subscription)
}

func (s *Service) UpdateSubscriptionsForWorkspaces(workspaceIDs []uuid.UUID, params UpdateSubscriptionParams) error {
	subscriptions, err := s.Repo.GetActiveSubscriptionsByWorkspaceIDs(workspaceIDs)
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

	return s.Repo.BulkUpdate(&subscriptions)
}

func (s *Service) DisableSubscriptions(ids []uuid.UUID, reason string) error {
	subscriptions, err := s.Repo.GetSubscriptionsByIDs(ids)
	if err != nil {
		return fmt.Errorf("subscription_service: disable_subscriptions: get_subscriptions: %w", err)
	}

	currentTime := time.Now()
	for idx := range subscriptions {
		subscriptions[idx].DisabledAt = &currentTime
		subscriptions[idx].DisabledReason = &reason
		subscriptions[idx].Status = models.SubscriptionStatusInactive
	}

	return s.Repo.BulkUpdate(&subscriptions)
}
