package models

import (
	"time"

	"github.com/google/uuid"
)

// Newsletter represents a generated and (optionally) sent newsletter for a repository.
type Newsletter struct {
	ID             uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	RepoID         uuid.UUID  `json:"repoId" gorm:"type:uuid;not null"`
	Topics         JSONBField `json:"topics" gorm:"type:jsonb;not null;default:'[]'"`
	HTMLContent    string     `json:"htmlContent" gorm:"column:html_content;not null;default:''"`
	SentAt         *time.Time `json:"sentAt"`
	RecipientCount int        `json:"recipientCount" gorm:"column:recipient_count;not null;default:0"`
	OpenCount      int        `json:"openCount" gorm:"column:open_count;not null;default:0"`
	ClickCount     int        `json:"clickCount" gorm:"column:click_count;not null;default:0"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}
