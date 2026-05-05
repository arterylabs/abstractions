package tests

import (
	ai_client "abstractions/clients/ai"
	"abstractions/config"
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/tests/testutils"
	"abstractions/workers/jobs"
	"context"
	"encoding/json"
	"testing"
	"time"

	integration_service "abstractions/services/integrations"
	github_integration "abstractions/services/integrations/core/github"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ---- mock AI provider ----

type mockAIProvider struct {
	summaryText string
	domainTags  []string
	embedding   []float32
}

func newMockAIProvider() *mockAIProvider {
	vec := make([]float32, 1536)
	vec[0] = 0.5
	return &mockAIProvider{
		summaryText: "handles user authentication",
		domainTags:  []string{"auth", "security"},
		embedding:   vec,
	}
}

func (m *mockAIProvider) SummarizeCode(_ context.Context, _, _, _ string) (*ai_client.SummarizeResult, error) {
	return &ai_client.SummarizeResult{
		Summary:    m.summaryText,
		DomainTags: m.domainTags,
	}, nil
}

func (m *mockAIProvider) GenerateTopics(_ context.Context, _ *ai_client.TopicGenerationRequest) ([]ai_client.Topic, error) {
	return nil, nil
}

func (m *mockAIProvider) GenerateNewsletterEntry(_ context.Context, _ *ai_client.NewsletterEntryRequest) (*ai_client.NewsletterEntry, error) {
	return nil, nil
}

func (m *mockAIProvider) GenerateEmbedding(_ context.Context, _ string) ([]float32, error) {
	return m.embedding, nil
}

// ---- mock file fetcher ----

type mockFileFetcher struct {
	content string
}

func (m *mockFileFetcher) GetFileContent(_ context.Context, _, _, _, _, _ string) ([]byte, error) {
	return []byte(m.content), nil
}

// ---- test suite ----

type SummarizeFileJobTestSuite struct {
	testutils.TestSuite
	seeder        *testutils.Seeder
	intSvc        *integration_service.Service
	fileIndexRepo *repositories.FileIndexRepository
	chunkRepo     *repositories.ChunkEmbeddingRepository
}

func TestSummarizeFileJob(t *testing.T) {
	suite.Run(t, new(SummarizeFileJobTestSuite))
}

func (s *SummarizeFileJobTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.intSvc = integration_service.NewIntegrationService(s.DB)
	s.fileIndexRepo = repositories.NewFileIndexRepository(s.DB)
	s.chunkRepo = repositories.NewChunkEmbeddingRepository(s.DB)
}

// createGitHubIntegration seeds a GitHub OAuth integration with a plain-text token that
// the integration service will encrypt.
func (s *SummarizeFileJobTestSuite) createGitHubIntegration(workspaceID, userID uuid.UUID) {
	_, err := s.intSvc.CreateIntegration(s.Ctx, &integration_service.CreateIntegrationRequest{
		Type:        github_integration.IntegrationTypeGitHub,
		WorkspaceID: workspaceID,
		Token:       "ghp_test_token",
		CreatedBy:   userID,
	})
	s.Require().NoError(err)
}

// buildTask serialises the given payload into an asynq.Task.
func buildTask(payload jobs.SummarizeFilePayload) *asynq.Task {
	b, _ := json.Marshal(payload)
	return asynq.NewTask(jobs.SummarizeFileJobID, b)
}

// TestSummarizeFileJob_CreatesChunkEmbeddings verifies that running the handler inserts
// chunk_embeddings rows for the processed file.
func (s *SummarizeFileJobTestSuite) TestSummarizeFileJob_CreatesChunkEmbeddings() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	s.createGitHubIntegration(workspace.ID, user.ID)
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	now := time.Now()
	s.seeder.CreateFileIndex(repo.ID, "auth/handler.go", func(fi *models.FileIndex) {
		fi.GitSHA = "abc123"
		fi.LastProcessedAt = &now
	})

	goSrc := `func ValidateToken(ctx context.Context, token string) error {
	return nil
}

func RefreshToken(ctx context.Context, old string) (string, error) {
	return "", nil
}
`
	job := jobs.NewSummarizeFileJobWithDeps(s.DB, &config.Config{}, newMockAIProvider(), &mockFileFetcher{content: goSrc})

	payload := jobs.SummarizeFilePayload{
		RepositoryID: repo.ID.String(),
		FilePath:     "auth/handler.go",
		GitSHA:       "abc123",
	}
	err := job.Handler(s.Ctx, buildTask(payload))
	assert.NoError(s.T(), err)

	chunks, err := s.chunkRepo.QuerySimilar(context.Background(), repo.ID, make([]float32, 1536), 10)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), chunks)
}

// TestSummarizeFileJob_UpdatesDomainTagsAndLastProcessedAt verifies that after processing
// the file_index row has its domain_tags and last_processed_at updated.
func (s *SummarizeFileJobTestSuite) TestSummarizeFileJob_UpdatesDomainTagsAndLastProcessedAt() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	s.createGitHubIntegration(workspace.ID, user.ID)
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)
	s.seeder.CreateFileIndex(repo.ID, "main.go")

	job := jobs.NewSummarizeFileJobWithDeps(s.DB, &config.Config{}, newMockAIProvider(), &mockFileFetcher{content: "func main() {}\n"})

	payload := jobs.SummarizeFilePayload{
		RepositoryID: repo.ID.String(),
		FilePath:     "main.go",
		GitSHA:       "sha1",
	}
	err := job.Handler(s.Ctx, buildTask(payload))
	assert.NoError(s.T(), err)

	fi, err := s.fileIndexRepo.GetByRepoAndPath(context.Background(), repo.ID, "main.go")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), fi)
	assert.NotNil(s.T(), fi.LastProcessedAt)
	assert.Contains(s.T(), []string(fi.DomainTags), "auth")
	assert.Contains(s.T(), []string(fi.DomainTags), "security")
}

// TestSummarizeFileJob_Idempotency verifies that running the handler twice with the same
// payload does not duplicate rows in chunk_embeddings.
func (s *SummarizeFileJobTestSuite) TestSummarizeFileJob_Idempotency() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	s.createGitHubIntegration(workspace.ID, user.ID)
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)
	s.seeder.CreateFileIndex(repo.ID, "util.go")

	goSrc := "func Helper() {}\n"
	job := jobs.NewSummarizeFileJobWithDeps(s.DB, &config.Config{}, newMockAIProvider(), &mockFileFetcher{content: goSrc})
	payload := jobs.SummarizeFilePayload{
		RepositoryID: repo.ID.String(),
		FilePath:     "util.go",
		GitSHA:       "sha_x",
	}

	assert.NoError(s.T(), job.Handler(s.Ctx, buildTask(payload)))
	assert.NoError(s.T(), job.Handler(s.Ctx, buildTask(payload)))

	var count int64
	s.DB.Model(&models.ChunkEmbedding{}).Where("repo_id = ? AND file_path = ?", repo.ID, "util.go").Count(&count)
	assert.Equal(s.T(), int64(1), count, "re-processing should upsert, not duplicate")
}

// TestSummarizeFileJob_RawContentNotPersisted verifies that the chunk_embeddings rows
// contain only summary text, not raw source code.
func (s *SummarizeFileJobTestSuite) TestSummarizeFileJob_RawContentNotPersisted() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	s.createGitHubIntegration(workspace.ID, user.ID)
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)
	s.seeder.CreateFileIndex(repo.ID, "secret.go")

	secretCode := "func SecretAlgorithm_DO_NOT_LOG() { /* proprietary */ }"
	job := jobs.NewSummarizeFileJobWithDeps(s.DB, &config.Config{}, newMockAIProvider(), &mockFileFetcher{content: secretCode})
	payload := jobs.SummarizeFilePayload{
		RepositoryID: repo.ID.String(),
		FilePath:     "secret.go",
		GitSHA:       "sha_s",
	}

	assert.NoError(s.T(), job.Handler(s.Ctx, buildTask(payload)))

	// Raw code should not appear in any ChunkEmbedding.Summary
	chunks, _ := s.chunkRepo.QuerySimilar(context.Background(), repo.ID, make([]float32, 1536), 10)
	for _, c := range chunks {
		assert.NotContains(s.T(), c.Summary, "SecretAlgorithm_DO_NOT_LOG",
			"raw source code must not be stored in chunk summary")
	}
}
