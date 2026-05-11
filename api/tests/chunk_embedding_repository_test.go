package tests

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/tests/testutils"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	pgvector "github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ChunkEmbeddingRepositoryTestSuite struct {
	testutils.TestSuite
	seeder *testutils.Seeder
	repo   *repositories.ChunkEmbeddingRepository
}

func TestChunkEmbeddingRepository(t *testing.T) {
	suite.Run(t, new(ChunkEmbeddingRepositoryTestSuite))
}

func (s *ChunkEmbeddingRepositoryTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.repo = repositories.NewChunkEmbeddingRepository(s.DB)
}

// makeEmbedding builds a 1536-dim vector with the first element set to val.
func makeEmbedding(val float32) []float32 {
	vec := make([]float32, 1536)
	vec[0] = val
	return vec
}

// TestUpsert_Idempotency verifies that upserting the same (repo_id, file_path, chunk_index)
// updates the existing row rather than inserting a duplicate.
func (s *ChunkEmbeddingRepositoryTestSuite) TestUpsert_Idempotency() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	original := models.ChunkEmbedding{
		ID:           uuid.New(),
		RepoID:       repo.ID,
		RepoFilePath: "main.go",
		ChunkIndex:   0,
		Summary:      "original summary",
		Tags:         models.StringArray{},
		Embedding:    pgvector.NewVector(makeEmbedding(0.1)),
		CreatedAt:    time.Now(),
	}

	err := s.repo.Upsert(context.Background(), []models.ChunkEmbedding{original})
	assert.NoError(s.T(), err)

	updated := original
	updated.Summary = "updated summary"
	updated.Embedding = pgvector.NewVector(makeEmbedding(0.9))

	err = s.repo.Upsert(context.Background(), []models.ChunkEmbedding{updated})
	assert.NoError(s.T(), err)

	// Should still have exactly one row
	var count int64
	s.DB.Model(&models.ChunkEmbedding{}).
		Where("repo_id = ? AND repo_file_path = ? AND chunk_index = ?", repo.ID, "main.go", 0).
		Count(&count)
	assert.Equal(s.T(), int64(1), count)

	// Summary should reflect the update
	var got models.ChunkEmbedding
	s.DB.Where("repo_id = ? AND repo_file_path = ? AND chunk_index = ?", repo.ID, "main.go", 0).
		First(&got)
	assert.Equal(s.T(), "updated summary", got.Summary)
}

// TestDeleteByFile_RemovesAllChunks verifies that all chunks for a given file are deleted,
// while chunks for other files in the same repo are preserved.
func (s *ChunkEmbeddingRepositoryTestSuite) TestDeleteByFile_RemovesAllChunks() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	s.seeder.CreateChunkEmbedding(repo.ID, "delete.go", 0, makeEmbedding(0.1))
	s.seeder.CreateChunkEmbedding(repo.ID, "delete.go", 1, makeEmbedding(0.2))
	s.seeder.CreateChunkEmbedding(repo.ID, "keep.go", 0, makeEmbedding(0.3))

	err := s.repo.DeleteByFile(context.Background(), repo.ID, "delete.go")
	assert.NoError(s.T(), err)

	var deleted []models.ChunkEmbedding
	s.DB.Where("repo_id = ? AND repo_file_path = ?", repo.ID, "delete.go").Find(&deleted)
	assert.Len(s.T(), deleted, 0)

	var kept []models.ChunkEmbedding
	s.DB.Where("repo_id = ? AND repo_file_path = ?", repo.ID, "keep.go").Find(&kept)
	assert.Len(s.T(), kept, 1)
}

// TestDeleteByRepo_RemovesAllChunksForRepo verifies that all embeddings for a repo are deleted.
func (s *ChunkEmbeddingRepositoryTestSuite) TestDeleteByRepo_RemovesAllChunksForRepo() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo1 := s.seeder.CreateRepository(workspace.ID, user.ID)
	repo2 := s.seeder.CreateRepository(workspace.ID, user.ID)

	s.seeder.CreateChunkEmbedding(repo1.ID, "a.go", 0, makeEmbedding(0.1))
	s.seeder.CreateChunkEmbedding(repo1.ID, "b.go", 0, makeEmbedding(0.2))
	s.seeder.CreateChunkEmbedding(repo2.ID, "c.go", 0, makeEmbedding(0.3))

	err := s.repo.DeleteByRepo(context.Background(), repo1.ID)
	assert.NoError(s.T(), err)

	var repo1Chunks []models.ChunkEmbedding
	s.DB.Where("repo_id = ?", repo1.ID).Find(&repo1Chunks)
	assert.Len(s.T(), repo1Chunks, 0)

	var repo2Chunks []models.ChunkEmbedding
	s.DB.Where("repo_id = ?", repo2.ID).Find(&repo2Chunks)
	assert.Len(s.T(), repo2Chunks, 1)
}

// TestQuerySimilar_ReturnsTopK verifies that QuerySimilar returns at most topK results
// and that results belong to the requested repo.
func (s *ChunkEmbeddingRepositoryTestSuite) TestQuerySimilar_ReturnsTopK() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)
	otherRepo := s.seeder.CreateRepository(workspace.ID, user.ID)

	for i := 0; i < 5; i++ {
		vec := makeEmbedding(float32(i) * 0.1)
		s.seeder.CreateChunkEmbedding(repo.ID, "file.go", i, vec)
	}
	s.seeder.CreateChunkEmbedding(otherRepo.ID, "other.go", 0, makeEmbedding(0.5))

	query := makeEmbedding(0.5)
	results, err := s.repo.QuerySimilar(context.Background(), repo.ID, query, 3)
	assert.NoError(s.T(), err)
	assert.LessOrEqual(s.T(), len(results), 3)

	for _, r := range results {
		assert.Equal(s.T(), repo.ID, r.RepoID)
	}
}
