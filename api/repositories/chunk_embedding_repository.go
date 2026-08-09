package repositories

import (
	"abstractions/db/models"
	"context"
	"fmt"
	"strings"

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
// On conflict (repo_id, repo_file_path, chunk_index) the existing row is updated in place.
func (r *ChunkEmbeddingRepository) Upsert(ctx context.Context, chunks []models.ChunkEmbedding) error {
	if len(chunks) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Exec(`INSERT INTO file_chunk_embeddings (
                id, repo_id, repo_file_id, repo_file_path, chunk_index, summary, tags,
                symbol_name, symbol_type, start_line, end_line, chunk_type, language,
                search_document, indexer_version, embedding, created_at
              ) VALUES `+buildUpsertPlaceholders(len(chunks))+`
              ON CONFLICT (repo_id, repo_file_path, chunk_index)
              DO UPDATE SET
                repo_file_id     = EXCLUDED.repo_file_id,
                summary          = EXCLUDED.summary,
                tags             = EXCLUDED.tags,
                symbol_name      = EXCLUDED.symbol_name,
                symbol_type      = EXCLUDED.symbol_type,
                start_line       = EXCLUDED.start_line,
                end_line         = EXCLUDED.end_line,
                chunk_type       = EXCLUDED.chunk_type,
                language         = EXCLUDED.language,
                search_document  = EXCLUDED.search_document,
                indexer_version  = EXCLUDED.indexer_version,
                embedding        = EXCLUDED.embedding`,
			buildUpsertArgs(chunks)...).Error
}

const colsPerRow = 17

func buildUpsertPlaceholders(n int) string {
	var result strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			result.WriteString(",")
		}
		result.WriteString("(")
		for j := 0; j < colsPerRow; j++ {
			if j > 0 {
				result.WriteString(",")
			}
			result.WriteString(fmt.Sprintf("$%d", i*colsPerRow+j+1))
		}
		result.WriteString(")")
	}
	return result.String()
}

func buildUpsertArgs(chunks []models.ChunkEmbedding) []interface{} {
	args := make([]interface{}, 0, len(chunks)*colsPerRow)
	for i := range chunks {
		c := &chunks[i]
		if c.ID == uuid.Nil {
			c.ID = uuid.New()
		}
		args = append(args,
			c.ID,
			c.RepoID,
			c.RepoFileID,
			c.RepoFilePath,
			c.ChunkIndex,
			c.Summary,
			c.Tags,
			c.SymbolName,
			c.SymbolType,
			c.StartLine,
			c.EndLine,
			c.ChunkType,
			c.Language,
			c.SearchDocument,
			c.IndexerVersion,
			c.Embedding,
			c.CreatedAt,
		)
	}
	return args
}

// QueryOptions controls hybrid retrieval behaviour.
type QueryOptions struct {
	// TopK is the maximum number of results (default 12).
	TopK int
	// ExcludeFileLevel drops chunk_index = -1 rows.
	ExcludeFileLevel bool
	// Language filters to a single language when non-empty.
	Language string
	// QueryText enables keyword ranking fused with vector distance via RRF.
	QueryText string
	// Diversify caps results per file path (0 = no cap; default 3 when TopK > 0).
	DiversifyPerFile int
}

// QuerySimilar returns the topK chunks most similar to the given embedding using cosine distance.
// When opts.QueryText is set, results are fused with full-text rank (Reciprocal Rank Fusion).
func (r *ChunkEmbeddingRepository) QuerySimilar(ctx context.Context, repoID uuid.UUID, embedding []float32, topK int) ([]models.ChunkEmbedding, error) {
	return r.QuerySimilarWithOptions(ctx, repoID, embedding, QueryOptions{TopK: topK})
}

// QuerySimilarWithOptions performs vector and optional hybrid retrieval.
func (r *ChunkEmbeddingRepository) QuerySimilarWithOptions(ctx context.Context, repoID uuid.UUID, embedding []float32, opts QueryOptions) ([]models.ChunkEmbedding, error) {
	if opts.TopK <= 0 {
		opts.TopK = 12
	}
	if opts.DiversifyPerFile <= 0 {
		opts.DiversifyPerFile = 3
	}

	fetchLimit := opts.TopK * 4
	if fetchLimit < 40 {
		fetchLimit = 40
	}

	where := []string{"repo_id = ?"}
	args := []interface{}{repoID}
	if opts.ExcludeFileLevel {
		where = append(where, "chunk_index >= 0")
	}
	if opts.Language != "" {
		where = append(where, "language = ?")
		args = append(args, opts.Language)
	}
	whereSQL := strings.Join(where, " AND ")

	// Vector candidates
	vecArgs := append(append([]interface{}{}, args...), pgvector.NewVector(embedding), fetchLimit)
	var vectorHits []models.ChunkEmbedding
	err := r.db.WithContext(ctx).
		Raw(`SELECT id, repo_id, repo_file_id, repo_file_path, chunk_index, summary, tags,
		            symbol_name, symbol_type, start_line, end_line, chunk_type, language,
		            search_document, indexer_version, embedding, created_at
		     FROM file_chunk_embeddings
		     WHERE `+whereSQL+`
		     ORDER BY embedding <=> ?
		     LIMIT ?`, vecArgs...).
		Scan(&vectorHits).Error
	if err != nil {
		return nil, fmt.Errorf("chunk_embedding: query similar vector: %w", err)
	}

	if strings.TrimSpace(opts.QueryText) == "" {
		return diversifyChunks(vectorHits, opts.TopK, opts.DiversifyPerFile), nil
	}

	// Keyword candidates
	textArgs := append(append([]interface{}{}, args...), opts.QueryText, opts.QueryText, fetchLimit)
	var textHits []models.ChunkEmbedding
	err = r.db.WithContext(ctx).
		Raw(`SELECT id, repo_id, repo_file_id, repo_file_path, chunk_index, summary, tags,
		            symbol_name, symbol_type, start_line, end_line, chunk_type, language,
		            search_document, indexer_version, embedding, created_at
		     FROM file_chunk_embeddings
		     WHERE `+whereSQL+` AND search_tsv @@ plainto_tsquery('english', ?)
		     ORDER BY ts_rank(search_tsv, plainto_tsquery('english', ?)) DESC
		     LIMIT ?`, textArgs...).
		Scan(&textHits).Error
	if err != nil {
		// Fall back to vector-only if FTS fails (e.g. empty tsv column before backfill).
		return diversifyChunks(vectorHits, opts.TopK, opts.DiversifyPerFile), nil
	}

	fused := reciprocalRankFusion(vectorHits, textHits, opts.TopK*2)
	return diversifyChunks(fused, opts.TopK, opts.DiversifyPerFile), nil
}

func reciprocalRankFusion(vectorHits, textHits []models.ChunkEmbedding, limit int) []models.ChunkEmbedding {
	const k = 60.0
	scores := map[string]float64{}
	byID := map[string]models.ChunkEmbedding{}

	for rank, c := range vectorHits {
		id := c.ID.String()
		scores[id] += 1.0 / (k + float64(rank+1))
		byID[id] = c
	}
	for rank, c := range textHits {
		id := c.ID.String()
		scores[id] += 1.0 / (k + float64(rank+1))
		byID[id] = c
	}

	type scored struct {
		id    string
		score float64
	}
	ordered := make([]scored, 0, len(scores))
	for id, s := range scores {
		ordered = append(ordered, scored{id: id, score: s})
	}
	// Simple insertion sort by score desc (N is small).
	for i := 1; i < len(ordered); i++ {
		j := i
		for j > 0 && ordered[j].score > ordered[j-1].score {
			ordered[j], ordered[j-1] = ordered[j-1], ordered[j]
			j--
		}
	}
	if limit > len(ordered) {
		limit = len(ordered)
	}
	out := make([]models.ChunkEmbedding, 0, limit)
	for i := 0; i < limit; i++ {
		out = append(out, byID[ordered[i].id])
	}
	return out
}

func diversifyChunks(chunks []models.ChunkEmbedding, topK, perFile int) []models.ChunkEmbedding {
	if topK <= 0 {
		return nil
	}
	counts := map[string]int{}
	out := make([]models.ChunkEmbedding, 0, topK)
	for _, c := range chunks {
		if counts[c.RepoFilePath] >= perFile {
			continue
		}
		counts[c.RepoFilePath]++
		out = append(out, c)
		if len(out) >= topK {
			break
		}
	}
	// If diversity left us short, fill from remaining without cap.
	if len(out) < topK {
		seen := map[string]struct{}{}
		for _, c := range out {
			seen[c.ID.String()] = struct{}{}
		}
		for _, c := range chunks {
			if _, ok := seen[c.ID.String()]; ok {
				continue
			}
			out = append(out, c)
			if len(out) >= topK {
				break
			}
		}
	}
	return out
}

// ListFileSummaries returns all file-level summary rows (chunk_index = -1) for a repository.
func (r *ChunkEmbeddingRepository) ListFileSummaries(ctx context.Context, repoID uuid.UUID) ([]models.ChunkEmbedding, error) {
	var results []models.ChunkEmbedding
	err := r.db.WithContext(ctx).
		Where("repo_id = ? AND chunk_index = -1", repoID).
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("chunk_embedding: list file summaries: %w", err)
	}
	return results, nil
}

// DeleteByFile removes all chunk embeddings for a specific file within a repository.
func (r *ChunkEmbeddingRepository) DeleteByFile(ctx context.Context, repoID uuid.UUID, filePath string) error {
	return r.db.WithContext(ctx).
		Where("repo_id = ? AND repo_file_path = ?", repoID, filePath).
		Delete(&models.ChunkEmbedding{}).Error
}

// DeleteByRepo removes all chunk embeddings for a repository.
func (r *ChunkEmbeddingRepository) DeleteByRepo(ctx context.Context, repoID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("repo_id = ?", repoID).
		Delete(&models.ChunkEmbedding{}).Error
}
