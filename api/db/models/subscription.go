package models

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive   SubscriptionStatus = "active"
	SubscriptionStatusInactive SubscriptionStatus = "inactive"
	SubscriptionStatusTrialing SubscriptionStatus = "trialing"
)

func (s SubscriptionStatus) String() string {
	return string(s)
}

type PlanID string

const (
	PlanIDPro      PlanID = "pro"
	PlanIDLifetime PlanID = "lifetime"
)

func NewPlanID(s string) PlanID {
	return PlanID(s)
}

func (p PlanID) String() string {
	return string(p)
}

type Subscription struct {
	ID             uuid.UUID          `json:"id"`
	WorkspaceID    uuid.UUID          `json:"workspace_id"`
	Status         SubscriptionStatus `json:"status"`
	PlanID         PlanID             `json:"plan_id"`
	ExpiresAt      time.Time          `json:"expires_at"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	DisabledAt     *time.Time         `json:"disabled_at"`
	DisabledReason *string            `json:"disabled_reason"`
	Limit          int                `json:"limit"`
}
