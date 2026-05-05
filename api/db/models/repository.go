package models

import (
	"time"

	"github.com/google/uuid"
)

// RepositoryStatus represents the indexing state of a connected repository.
type RepositoryStatus string

const (
	RepositoryStatusPending  RepositoryStatus = "pending"
	RepositoryStatusIndexing RepositoryStatus = "indexing"
	RepositoryStatusIndexed  RepositoryStatus = "indexed"
	RepositoryStatusError    RepositoryStatus = "error"
)

// Repository represents a GitHub repository connected to a workspace.
type Repository struct {
	ID                   uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	WorkspaceID          uuid.UUID        `json:"workspaceId" gorm:"type:uuid;not null"`
	UserID               uuid.UUID        `json:"userId" gorm:"type:uuid;not null"`
	GitHubRepoID         int64            `json:"githubRepoId" gorm:"column:github_repo_id;not null"`
	FullName             string           `json:"fullName" gorm:"not null"`
	DefaultBranch        string           `json:"defaultBranch" gorm:"not null;default:'main'"`
	WebhookID            *int64           `json:"webhookId" gorm:"column:webhook_id"`
	WebhookSecret        *string          `json:"webhookSecret" gorm:"column:webhook_secret"`
	Status               RepositoryStatus `json:"status" gorm:"not null;default:'pending'"`
	IndexedAt            *time.Time       `json:"indexedAt"`
	SendDay              int              `json:"sendDay" gorm:"not null;default:1"`
	SendHour             int              `json:"sendHour" gorm:"not null;default:9"`
	LastNewsletterSentAt *time.Time       `json:"lastNewsletterSentAt"`
	CreatedAt            time.Time        `json:"createdAt"`
	UpdatedAt            time.Time        `json:"updatedAt"`
	DeletedAt            *time.Time       `json:"deletedAt"`
}
