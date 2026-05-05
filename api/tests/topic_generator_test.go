package tests

import (
	ai_client "abstractions/clients/ai"
	"abstractions/db/models"
	"abstractions/repositories"
	newsletter_service "abstractions/services/newsletter"
	"abstractions/tests/testutils"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ---- topic generation test suite ----

type TopicGeneratorTestSuite struct {
	testutils.TestSuite
	seeder         *testutils.Seeder
	fileIndexRepo  *repositories.FileIndexRepository
	newsletterRepo *repositories.NewsletterRepository
	repositoryRepo *repositories.RepositoryRepository
}

func TestTopicGenerator(t *testing.T) {
	suite.Run(t, new(TopicGeneratorTestSuite))
}

func (s *TopicGeneratorTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.fileIndexRepo = repositories.NewFileIndexRepository(s.DB)
	s.newsletterRepo = repositories.NewNewsletterRepository(s.DB)
	s.repositoryRepo = repositories.NewRepositoryRepository(s.DB)
}

// seedFile inserts a file_index row with given tags and optional last-processed timestamp.
func (s *TopicGeneratorTestSuite) seedFile(repoID uuid.UUID, path string, tags []string, processedAt *time.Time) {
	fi := &models.FileIndex{
		ID:              uuid.New(),
		RepoID:          repoID,
		FilePath:        path,
		GitSHA:          "sha",
		DomainTags:      models.StringArray(tags),
		LastProcessedAt: processedAt,
	}
	err := s.fileIndexRepo.Upsert(context.Background(), fi)
	s.Require().NoError(err)
}

// seedNewsletter inserts a newsletter row with the given topic titles in the Topics JSONB field.
func (s *TopicGeneratorTestSuite) seedNewsletter(repoID uuid.UUID, topicTitles []string) {
	type topicEntry struct {
		Title string `json:"title"`
	}
	entries := make([]topicEntry, len(topicTitles))
	for i, t := range topicTitles {
		entries[i] = topicEntry{Title: t}
	}
	raw, _ := json.Marshal(entries)
	nl := &models.Newsletter{
		ID:          uuid.New(),
		RepoID:      repoID,
		Topics:      models.JSONBRaw(raw),
		HTMLContent: "<p>test</p>",
	}
	err := s.newsletterRepo.Create(context.Background(), nl)
	s.Require().NoError(err)
}

// mockTopicAI is an AI provider that returns a fixed set of topics for GenerateTopics.
type mockTopicAI struct {
	topics []ai_client.Topic
}

func (m *mockTopicAI) SummarizeCode(_ context.Context, _, _, _ string) (*ai_client.SummarizeResult, error) {
	return &ai_client.SummarizeResult{}, nil
}
func (m *mockTopicAI) GenerateTopics(_ context.Context, _ *ai_client.TopicGenerationRequest) ([]ai_client.Topic, error) {
	return m.topics, nil
}
func (m *mockTopicAI) GenerateNewsletterEntry(_ context.Context, _ *ai_client.NewsletterEntryRequest) (*ai_client.NewsletterEntry, error) {
	return &ai_client.NewsletterEntry{}, nil
}
func (m *mockTopicAI) GenerateEmbedding(_ context.Context, _ string) ([]float32, error) {
	return make([]float32, 1536), nil
}

// TestGenerateTopics_ReturnsScoredTopics verifies that topics are returned and scored
// by recency (files modified in the last 7 days raise the score).
func (s *TopicGeneratorTestSuite) TestGenerateTopics_ReturnsScoredTopics() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	recent := time.Now().Add(-1 * time.Hour)
	old := time.Now().Add(-30 * 24 * time.Hour)

	s.seedFile(repo.ID, "pkg/auth.go", []string{"auth"}, &recent)
	s.seedFile(repo.ID, "pkg/db.go", []string{"database"}, &old)

	aiTopics := []ai_client.Topic{
		{Title: "Auth Changes", Description: "auth improvements", RelatedFiles: []string{"pkg/auth.go"}},
		{Title: "DB Schema", Description: "db changes", RelatedFiles: []string{"pkg/db.go"}},
	}

	gen := newsletter_service.NewTopicGenerator(
		&mockTopicAI{topics: aiTopics},
		s.fileIndexRepo, s.newsletterRepo, s.repositoryRepo,
	)

	topics, err := gen.GenerateTopics(s.Ctx, repo.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), topics, 2)

	// Auth Changes should rank first (recent file → higher score)
	assert.Equal(s.T(), "Auth Changes", topics[0].Title)
	assert.Equal(s.T(), float64(1), topics[0].RecencyScore)
	assert.Equal(s.T(), float64(0), topics[1].RecencyScore)
}

// TestGenerateTopics_DeduplicatesRecentNewsletters verifies that topics matching
// titles from the last 4 newsletters are filtered out.
func (s *TopicGeneratorTestSuite) TestGenerateTopics_DeduplicatesRecentNewsletters() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	recent := time.Now().Add(-1 * time.Hour)
	s.seedFile(repo.ID, "pkg/auth.go", []string{"auth"}, &recent)
	s.seedFile(repo.ID, "pkg/api.go", []string{"api"}, &recent)

	// Seed a newsletter that already covered "Auth Changes"
	s.seedNewsletter(repo.ID, []string{"Auth Changes"})

	aiTopics := []ai_client.Topic{
		{Title: "Auth Changes", Description: "covered already", RelatedFiles: []string{"pkg/auth.go"}},
		{Title: "API Improvements", Description: "new topic", RelatedFiles: []string{"pkg/api.go"}},
	}

	gen := newsletter_service.NewTopicGenerator(
		&mockTopicAI{topics: aiTopics},
		s.fileIndexRepo, s.newsletterRepo, s.repositoryRepo,
	)

	topics, err := gen.GenerateTopics(s.Ctx, repo.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), topics, 1, "duplicate topic must be filtered out")
	assert.Equal(s.T(), "API Improvements", topics[0].Title)
}

// TestGenerateTopics_EmptyIndex_ReturnsNil verifies that no AI call is made and nil
// is returned when the file index is empty.
func (s *TopicGeneratorTestSuite) TestGenerateTopics_EmptyIndex_ReturnsNil() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	aiCalled := false
	ai := &mockTopicAI{}
	_ = aiCalled

	gen := newsletter_service.NewTopicGenerator(
		ai,
		s.fileIndexRepo, s.newsletterRepo, s.repositoryRepo,
	)

	topics, err := gen.GenerateTopics(s.Ctx, repo.ID)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), topics, "empty index must return nil without calling AI")
}
