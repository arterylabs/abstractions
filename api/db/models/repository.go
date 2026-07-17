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

	RepositoryServiceTypeGitHub = "GITHUB"
)

// Repository represents a connected repository.
type Repository struct {
	ID                     uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	WorkspaceID            uuid.UUID        `json:"workspaceId" gorm:"type:uuid;not null"`
	CreatedBy              uuid.UUID        `json:"createdBy" gorm:"type:uuid;not null"`
	UpdatedBy              *uuid.UUID       `json:"updatedBy" gorm:"type:uuid"`
	Name                   string           `json:"name" gorm:"not null"`
	Description            string           `json:"description" gorm:"column:description;not null;default:''"`
	ServiceType            string           `json:"serviceType" gorm:"column:service_type;not null;default:'GITHUB'"`
	ServiceID              string           `json:"serviceId" gorm:"column:service_id;not null"`
	DefaultBranch          string           `json:"defaultBranch" gorm:"not null;default:'main'"`
	Status                 RepositoryStatus `json:"status" gorm:"not null;default:'pending'"`
	IndexedAt              *time.Time       `json:"indexedAt"`
	IndexError             *string          `json:"indexError" gorm:"column:index_error"`
	SendDay                int              `json:"sendDay" gorm:"not null;default:1"`
	SendHour               int              `json:"sendHour" gorm:"not null;default:9"`
	NewsletterSentAt       *time.Time       `json:"newsletterSentAt"`
	Recipients             StringArray      `json:"recipients" gorm:"column:recipients;type:text[];not null;default:'{}'"`
	SendToWorkspaceMembers bool             `json:"sendToWorkspaceMembers" gorm:"column:send_to_workspace_members;not null;default:true"`
	NewsletterEnabled      bool             `json:"newsletterEnabled" gorm:"column:newsletter_enabled;not null;default:true"`
	CreatedAt              time.Time        `json:"createdAt"`
	UpdatedAt              time.Time        `json:"updatedAt"`
	DeletedAt              *time.Time       `json:"deletedAt" gorm:"type:timestamptz"`
	DeletedBy              *uuid.UUID       `json:"deletedBy" gorm:"type:uuid"`
}
