package models

import (
	"time"

	"github.com/google/uuid"
)

// Newsletter status constants representing the lifecycle of a newsletter record.
const (
	NewsletterStatusQueued     = "queued"
	NewsletterStatusProcessing = "processing"
	NewsletterStatusSent       = "sent"
	NewsletterStatusError      = "error"
)

// Newsletter represents a generated and (optionally) sent newsletter for a repository.
type Newsletter struct {
	ID             uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	RepoID         uuid.UUID        `json:"repoId" gorm:"type:uuid;not null"`
	TopicID        *uuid.UUID       `json:"topicId" gorm:"type:uuid;column:topic_id"`
	Topic          *NewsletterTopic `json:"topic,omitempty" gorm:"foreignKey:TopicID;references:ID"`
	HTMLContent    string           `json:"htmlContent" gorm:"column:html_content;not null;default:''"`
	Status         string           `json:"status" gorm:"column:status;not null;default:'sent'"`
	ErrorMessage   *string          `json:"errorMessage" gorm:"column:error_message"`
	SentAt         *time.Time       `json:"sentAt"`
	RecipientCount int              `json:"recipientCount" gorm:"column:recipient_count;not null;default:0"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
	DeletedAt      *time.Time       `json:"deletedAt"`
	DeletedBy      *uuid.UUID       `json:"deletedBy" gorm:"type:uuid"`
}
