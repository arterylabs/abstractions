package models

import (
	"time"

	"github.com/google/uuid"
)

// NewsletterTopic represents a queued or sent topic associated with a repository newsletter.
type NewsletterTopic struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	RepoID       uuid.UUID  `json:"repoId" gorm:"type:uuid;not null;column:repo_id"`
	CreatedBy    *uuid.UUID `json:"createdBy" gorm:"type:uuid;column:created_by"`
	UpdatedBy    *uuid.UUID `json:"updatedBy" gorm:"type:uuid;column:updated_by"`
	Title        string     `json:"title" gorm:"not null;default:''"`
	Description  string     `json:"description" gorm:"not null;default:''"`
	RelatedFiles JSONBRaw   `json:"relatedFiles" gorm:"type:jsonb;not null;default:'[]'"`
	Position     int        `json:"position" gorm:"not null;default:0"`
	SentAt       *time.Time `json:"sentAt" gorm:"type:timestamptz;column:sent_at"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `json:"deletedAt" gorm:"type:timestamptz"`
	DeletedBy    *uuid.UUID `json:"deletedBy" gorm:"type:uuid"`
}
