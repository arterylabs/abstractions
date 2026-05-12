package tests

import (
	ai_client "abstractions/clients/ai"
	"abstractions/db/models"
	"abstractions/repositories"
	newsletter_service "abstractions/services/newsletter"
	"abstractions/tests/testutils"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	pgvector "github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ---- content generator test suite ----

type ContentGeneratorTestSuite struct {
	testutils.TestSuite
	seeder    *testutils.Seeder
	chunkRepo *repositories.ChunkEmbeddingRepository
	repoRepo  *repositories.RepositoryRepository
}

func TestContentGenerator(t *testing.T) {
	suite.Run(t, new(ContentGeneratorTestSuite))
}

func (s *ContentGeneratorTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.chunkRepo = repositories.NewChunkEmbeddingRepository(s.DB)
	s.repoRepo = repositories.NewRepositoryRepository(s.DB)
}

// seedChunk inserts a chunk embedding with a deterministic embedding vector.
func (s *ContentGeneratorTestSuite) seedChunk(repoID uuid.UUID, filePath, summary string, vecVal float32) {
	vec := make([]float32, 1536)
	vec[0] = vecVal
	err := s.chunkRepo.Upsert(s.Ctx, []models.ChunkEmbedding{{
		ID:           uuid.New(),
		RepoID:       repoID,
		RepoFilePath: filePath,
		ChunkIndex:   0,
		Summary:      summary,
		Embedding:    pgvector.NewVector(vec),
		CreatedAt:    time.Now(),
	}})
	s.Require().NoError(err)
}

// mockContentAI returns a fixed short HTML content for newsletter entries.
type mockContentAI struct {
	content string
}

func (m *mockContentAI) SummarizeCode(_ context.Context, _, _, _ string, _, _ int, _ string) (*ai_client.SummarizeResult, error) {
	return &ai_client.SummarizeResult{}, nil
}
func (m *mockContentAI) SummarizeFile(_ context.Context, _, _ string, _ []string) (*ai_client.SummarizeResult, error) {
	return &ai_client.SummarizeResult{}, nil
}
func (m *mockContentAI) GenerateTopics(_ context.Context, _ *ai_client.TopicGenerationRequest) ([]ai_client.Topic, error) {
	return nil, nil
}
func (m *mockContentAI) GenerateNewsletterEntry(_ context.Context, _ *ai_client.NewsletterEntryRequest) (*ai_client.NewsletterEntry, error) {
	return &ai_client.NewsletterEntry{Content: m.content}, nil
}
func (m *mockContentAI) GenerateEmbedding(_ context.Context, _ string) ([]float32, error) {
	vec := make([]float32, 1536)
	vec[0] = 0.5
	return vec, nil
}

// TestGenerateNewsletter_ReturnsEntryPerTopic checks that one entry is produced per topic.
func (s *ContentGeneratorTestSuite) TestGenerateNewsletter_ReturnsEntryPerTopic() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	s.seedChunk(repo.ID, "pkg/auth.go", "handles JWT authentication", 0.9)
	s.seedChunk(repo.ID, "pkg/db.go", "manages database connections", 0.1)

	topics := []newsletter_service.Topic{
		{Title: "Auth Improvements", Description: "JWT changes"},
		{Title: "DB Optimisations", Description: "query speed"},
	}

	const html = "<p>Great improvements this week.</p>"
	gen := newsletter_service.NewContentGenerator(
		&mockContentAI{content: html},
		s.chunkRepo, s.repoRepo,
	)

	entries, err := gen.GenerateNewsletter(s.Ctx, repo.ID, topics)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), entries, 2)
	assert.Equal(s.T(), "Auth Improvements", entries[0].TopicTitle)
	assert.Equal(s.T(), html, entries[0].Content)
}

// TestGenerateNewsletter_FilePathsFromChunks verifies that FilePathsReferenced comes
// exclusively from the chunk_embeddings returned by QuerySimilar.
func (s *ContentGeneratorTestSuite) TestGenerateNewsletter_FilePathsFromChunks() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	s.seedChunk(repo.ID, "src/handler.go", "HTTP request handler", 0.8)
	s.seedChunk(repo.ID, "src/middleware.go", "auth middleware", 0.7)

	topics := []newsletter_service.Topic{
		{Title: "Handler Refactor", Description: "improved routing"},
	}

	gen := newsletter_service.NewContentGenerator(
		&mockContentAI{content: "<p>refactored handlers</p>"},
		s.chunkRepo, s.repoRepo,
	)

	entries, err := gen.GenerateNewsletter(s.Ctx, repo.ID, topics)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), entries, 1)

	// All referenced file paths must come from chunk_embeddings for this repo
	for _, p := range entries[0].FilePathsReferenced {
		assert.True(s.T(),
			p == "src/handler.go" || p == "src/middleware.go",
			"unexpected file path: %s", p)
	}
}

// TestAssembler_RendersHTMLEmail verifies the assembler produces valid HTML with expected content.
func (s *ContentGeneratorTestSuite) TestAssembler_RendersHTMLEmail() {
	entries := []newsletter_service.NewsletterEntry{
		{TopicTitle: "Auth Refactor", Content: "<p>JWT improvements landed.</p>", FilePathsReferenced: []string{"pkg/auth.go"}},
	}

	assembler := &newsletter_service.Assembler{}
	html, err := assembler.Assemble("owner/repo", entries)
	assert.NoError(s.T(), err)
	assert.Contains(s.T(), html, "owner/repo")
	assert.Contains(s.T(), html, "JWT improvements landed")
	assert.True(s.T(), strings.HasPrefix(strings.TrimSpace(html), "<h2"))
}
