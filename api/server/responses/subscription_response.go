package responses

import (
	"abstractions/db/models"
	"time"

	"github.com/google/uuid"
)

type SubscriptionResponse struct {
	ID          uuid.UUID `json:"id"`
	WorkspaceID uuid.UUID `json:"workspaceId"`
	Status      string    `json:"status"`
	PlanID      string    `json:"planId"`
	ExpiresAt   time.Time `json:"expiresAt"`
	CreatedAt   time.Time `json:"createdAt"`
	Limit       int       `json:"limit"`
}

func NewSubscriptionResponse(subscription *models.Subscription) *SubscriptionResponse {
	if subscription == nil {
		return nil
	}
	return &SubscriptionResponse{
		ID:          subscription.ID,
		WorkspaceID: subscription.WorkspaceID,
		Status:      subscription.Status.String(),
		PlanID:      subscription.PlanID.String(),
		ExpiresAt:   subscription.ExpiresAt,
		CreatedAt:   subscription.CreatedAt,
		Limit:       subscription.Limit,
	}
}
