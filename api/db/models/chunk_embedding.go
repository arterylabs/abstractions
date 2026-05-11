package models

import (
	"time"

	"github.com/google/uuid"
	pgvector "github.com/pgvector/pgvector-go"
)

// ChunkEmbedding stores a vector embedding and summary for a single code chunk or file.
// The unique constraint is on (repo_id, repo_file_path, chunk_index).
// chunk_index = -1 is reserved for the file-level summary row.
type ChunkEmbedding struct {
	ID           uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	RepoID       uuid.UUID       `json:"repoId" gorm:"type:uuid;not null"`
	RepoFileID   *uuid.UUID      `json:"repoFileId" gorm:"column:repo_file_id;type:uuid"`
	RepoFilePath string          `json:"repoFilePath" gorm:"column:repo_file_path;not null"`
	ChunkIndex   int             `json:"chunkIndex" gorm:"column:chunk_index;not null"`
	Summary      string          `json:"summary" gorm:"not null;default:''"`
	Tags         StringArray     `json:"tags" gorm:"column:tags;type:text[]"`
	Embedding    pgvector.Vector `json:"embedding" gorm:"type:vector(1536)"`
	CreatedAt    time.Time       `json:"createdAt"`
	DeletedAt    *time.Time      `json:"deletedAt" gorm:"type:timestamptz"`
	DeletedBy    *uuid.UUID      `json:"deletedBy" gorm:"type:uuid"`
}

func (ChunkEmbedding) TableName() string {
	return "file_chunk_embeddings"
}
