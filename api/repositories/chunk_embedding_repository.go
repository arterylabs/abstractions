package repositories

import (
	"abstractions/db/models"
	"context"
	"fmt"

	"github.com/google/uuid"
	pgvector "github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type ChunkEmbeddingRepository struct {
	db *gorm.DB
}

func NewChunkEmbeddingRepository(db *gorm.DB) *ChunkEmbeddingRepository {
	return &ChunkEmbeddingRepository{db: db}
}

// Upsert inserts or updates a batch of chunk embeddings.
// On conflict (repo_id, file_path, chunk_index) the existing row is updated in place.
func (r *ChunkEmbeddingRepository) Upsert(ctx context.Context, chunks []models.ChunkEmbedding) error {
	if len(chunks) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Exec(`INSERT INTO chunk_embeddings (id, repo_id, file_path, chunk_index, summary, domain_tags, embedding, created_at)
              VALUES `+buildUpsertPlaceholders(len(chunks))+`
              ON CONFLICT (repo_id, file_path, chunk_index)
              DO UPDATE SET
                summary     = EXCLUDED.summary,
                domain_tags = EXCLUDED.domain_tags,
                embedding   = EXCLUDED.embedding`,
			buildUpsertArgs(chunks)...).Error
}

// buildUpsertPlaceholders builds ($1,$2,...,$9),($10,...) style placeholders for bulk insert.
func buildUpsertPlaceholders(n int) string {
	const colsPerRow = 8
	result := ""
	for i := 0; i < n; i++ {
		if i > 0 {
			result += ","
		}
		result += "("
		for j := 0; j < colsPerRow; j++ {
			if j > 0 {
				result += ","
			}
			result += fmt.Sprintf("$%d", i*colsPerRow+j+1)
		}
		result += ")"
	}
	return result
}

func buildUpsertArgs(chunks []models.ChunkEmbedding) []interface{} {
	args := make([]interface{}, 0, len(chunks)*8)
	for i := range chunks {
		c := &chunks[i]
		if c.ID == uuid.Nil {
			c.ID = uuid.New()
		}
		args = append(args,
			c.ID,
			c.RepoID,
			c.FilePath,
			c.ChunkIndex,
			c.Summary,
			c.DomainTags,
			c.Embedding,
			c.CreatedAt,
		)
	}
	return args
}

// QuerySimilar returns the topK chunks most similar to the given embedding using cosine distance.
func (r *ChunkEmbeddingRepository) QuerySimilar(ctx context.Context, repoID uuid.UUID, embedding []float32, topK int) ([]models.ChunkEmbedding, error) {
	var results []models.ChunkEmbedding
	err := r.db.WithContext(ctx).
		Raw(`SELECT * FROM chunk_embeddings
             WHERE repo_id = ?
             ORDER BY embedding <=> ?
             LIMIT ?`,
			repoID, pgvector.NewVector(embedding), topK).
		Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("chunk_embedding: query similar: %w", err)
	}
	return results, nil
}

// DeleteByFile removes all chunk embeddings for a specific file within a repository.
func (r *ChunkEmbeddingRepository) DeleteByFile(ctx context.Context, repoID uuid.UUID, filePath string) error {
	return r.db.WithContext(ctx).
		Where("repo_id = ? AND file_path = ?", repoID, filePath).
		Delete(&models.ChunkEmbedding{}).Error
}

// DeleteByRepo removes all chunk embeddings for a repository.
func (r *ChunkEmbeddingRepository) DeleteByRepo(ctx context.Context, repoID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("repo_id = ?", repoID).
		Delete(&models.ChunkEmbedding{}).Error
}
