package models

import (
	"time"

	"github.com/google/uuid"
)

// Email log statuses.
const (
	EmailLogStatusSent       = "sent"
	EmailLogStatusFailed     = "failed"
	EmailLogStatusSuppressed = "suppressed"
)

// Email log error types. These classify why a send did not succeed so that the
// suppression logic can decide whether an address should be blocked.
const (
	EmailErrorTypeHard          = "hard"           // recipient rejected (invalid address) -> suppress immediately
	EmailErrorTypeSoft          = "soft"           // transient failure -> counts toward the threshold
	EmailErrorTypeSystem        = "system"         // our misconfiguration -> logged but never counts against the address
	EmailErrorTypeBlockedDomain = "blocked_domain" // disposable/blocked domain -> skipped before send
	EmailErrorTypeSuppressed    = "suppressed"     // address already on the suppression list -> skipped before send
)

// EmailLog is an append-only record of every email send attempt.
type EmailLog struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email        string    `json:"email" gorm:"column:email;not null"`
	MessageID    *string   `json:"messageId" gorm:"column:message_id"`
	Subject      *string   `json:"subject" gorm:"column:subject"`
	Status       string    `json:"status" gorm:"column:status;not null"`
	ErrorType    *string   `json:"errorType" gorm:"column:error_type"`
	ErrorMessage *string   `json:"errorMessage" gorm:"column:error_message"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
}
