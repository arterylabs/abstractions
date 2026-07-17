package models

import (
	"time"

	"github.com/google/uuid"
)

// Suppression reasons.
const (
	SuppressionReasonHardBounce        = "hard_bounce"
	SuppressionReasonComplaint         = "complaint"
	SuppressionReasonThresholdExceeded = "threshold_exceeded"
	SuppressionReasonBlockedDomain     = "blocked_domain"
	SuppressionReasonManual            = "manual"
)

// EmailSuppression marks an address that must not be emailed. Soft-deletable so an
// address can be manually unblocked while preserving history.
type EmailSuppression struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email     string     `json:"email" gorm:"column:email;not null"`
	Reason    string     `json:"reason" gorm:"column:reason;not null"`
	CreatedAt time.Time  `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}
