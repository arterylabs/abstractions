package models

import (
	"time"

	"github.com/google/uuid"
)

// RepositoryFile represents an indexed file within a connected repository.
type RepositoryFile struct {
	ID              uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	RepoID          uuid.UUID   `json:"repoId" gorm:"type:uuid;not null"`
	Path            string      `json:"path" gorm:"column:path;not null"`
	GitSHA          string      `json:"gitSha" gorm:"column:git_sha;not null;default:''"`
	Tags            StringArray `json:"tags" gorm:"column:tags;type:text[];not null;default:'{}'"`
	ProcessedAt     *time.Time  `json:"processedAt" gorm:"column:processed_at"`
	Size            int         `json:"size" gorm:"column:size;not null;default:0"`
	Language        string      `json:"language" gorm:"column:language;not null;default:''"`
	ServiceMetadata JSONB       `json:"serviceMetadata" gorm:"column:service_metadata;type:jsonb"`
	DeletedAt       *time.Time  `json:"deletedAt" gorm:"type:timestamptz;column:deleted_at"`
	DeletedBy       *uuid.UUID  `json:"deletedBy" gorm:"type:uuid;column:deleted_by"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt"`
}

func (RepositoryFile) TableName() string {
	return "repository_files"
}
