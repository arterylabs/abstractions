package models

import (
	"time"

	"github.com/google/uuid"
)

// TeamMember represents a subscriber to a repository's newsletter.
type TeamMember struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	RepoID     uuid.UUID `json:"repoId" gorm:"type:uuid;not null"`
	Email      string    `json:"email" gorm:"not null"`
	Subscribed bool      `json:"subscribed" gorm:"not null;default:true"`
	AddedAt    time.Time `json:"addedAt" gorm:"column:added_at;not null;default:NOW()"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
