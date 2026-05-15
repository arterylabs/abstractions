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
	ID           uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	RepoID       uuid.UUID   `json:"repoId" gorm:"type:uuid;not null"`
	RepoFileID   *uuid.UUID  `json:"repoFileId" gorm:"column:repo_file_id;type:uuid"`
	RepoFilePath string      `json:"repoFilePath" gorm:"column:repo_file_path;not null"`
	ChunkIndex   int         `json:"chunkIndex" gorm:"column:chunk_index;not null"`
	Summary      string      `json:"summary" gorm:"not null;default:''"`
	Tags         StringArray `json:"tags" gorm:"column:tags;type:text[];not null;default:'{}'"`
	// SymbolName is the name of the primary symbol (function, class, etc.) in this chunk.
	// Empty for file-level summary rows.
	SymbolName string `json:"symbolName" gorm:"column:symbol_name;not null;default:''"`
	// SymbolType is the kind of symbol: "function", "class", "method", "variable", "file", etc.
	SymbolType string          `json:"symbolType" gorm:"column:symbol_type;not null;default:''"`
	Embedding  pgvector.Vector `json:"embedding" gorm:"type:vector(1536)"`
	CreatedAt  time.Time       `json:"createdAt"`
}

func (ChunkEmbedding) TableName() string {
	return "file_chunk_embeddings"
}
