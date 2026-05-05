package models

import (
	"time"

	"github.com/google/uuid"
	pgvector "github.com/pgvector/pgvector-go"
)

// ChunkEmbedding stores a vector embedding and summary for a single code chunk.
// The unique constraint is on (repo_id, file_path, chunk_index).
type ChunkEmbedding struct {
	ID         uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	RepoID     uuid.UUID       `json:"repoId" gorm:"type:uuid;not null"`
	FilePath   string          `json:"filePath" gorm:"not null"`
	ChunkIndex int             `json:"chunkIndex" gorm:"column:chunk_index;not null"`
	Summary    string          `json:"summary" gorm:"not null;default:''"`
	DomainTags StringArray     `json:"domainTags" gorm:"column:domain_tags;type:text[]"`
	Embedding  pgvector.Vector `json:"embedding" gorm:"type:vector(1536)"`
	CreatedAt  time.Time       `json:"createdAt"`
}
