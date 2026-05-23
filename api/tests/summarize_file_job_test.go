package tests

import (
	ai_client "abstractions/clients/ai"
	github_client "abstractions/clients/github"
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

func (m *mockAIProvider) SummarizeCode(_ context.Context, _, _, _ string, _, _ int, _ string) (*ai_client.SummarizeResult, error) {
	return &ai_client.SummarizeResult{
		Summary:    m.summaryText,
		DomainTags: m.domainTags,
	}, nil
}
func (m *mockAIProvider) SummarizeFile(_ context.Context, _, _ string, _ []string) (*ai_client.SummarizeResult, error) {
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

func (m *mockAIProvider) GenerateEmbeddings(_ context.Context, texts []string) ([][]float32, error) {
	out := make([][]float32, len(texts))
	for i := range out {
		out[i] = m.embedding
	}
	return out, nil
}

// ---- mock file fetcher ----

type mockFileFetcher struct {
	content string
}

func (m *mockFileFetcher) GetFileContent(_ context.Context, _, _, _, _, _ string) ([]byte, *github_client.GitHubFileMetadata, error) {
	meta := &github_client.GitHubFileMetadata{
		Type: "file",
		Name: "mock_file.go",
		Path: "mock/mock_file.go",
	}
	return []byte(m.content), meta, nil
}

// ---- test suite ----

type SummarizeFileJobTestSuite struct {
	testutils.TestSuite
	seeder        *testutils.Seeder
	intSvc        *integration_service.Service
	fileIndexRepo *repositories.RepositoryFileRepository
	chunkRepo     *repositories.ChunkEmbeddingRepository
	repoRepo      *repositories.RepositoryRepository
}

func TestSummarizeFileJob(t *testing.T) {
	suite.Run(t, new(SummarizeFileJobTestSuite))
}

func (s *SummarizeFileJobTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.intSvc = integration_service.NewIntegrationService(s.DB)
	s.fileIndexRepo = repositories.NewRepositoryFileRepository(s.DB)
	s.chunkRepo = repositories.NewChunkEmbeddingRepository(s.DB)
	s.repoRepo = repositories.NewRepositoryRepository(s.DB)
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
	s.seeder.CreateFileIndex(repo.ID, "auth/handler.go", func(fi *models.RepositoryFile) {
		fi.GitSHA = "abc123"
		fi.ProcessedAt = &now
	})

	goSrc := `func ValidateToken(ctx context.Context, token string) error {
	return nil
}

func RefreshToken(ctx context.Context, old string) (string, error) {
	return "", nil
}
`
	job := jobs.NewSummarizeFileJobWithDeps(s.DB, &config.Config{}, newMockAIProvider(), &mockFileFetcher{content: goSrc}, nil)

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

// TestSummarizeFileJob_UpdatesDomainTagsAndProcessedAt verifies that after processing
// the file_index row has its domain_tags and processed_at updated.
func (s *SummarizeFileJobTestSuite) TestSummarizeFileJob_UpdatesDomainTagsAndProcessedAt() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	s.createGitHubIntegration(workspace.ID, user.ID)
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)
	s.seeder.CreateFileIndex(repo.ID, "main.go")

	job := jobs.NewSummarizeFileJobWithDeps(s.DB, &config.Config{}, newMockAIProvider(), &mockFileFetcher{content: "func main() {}\n"}, nil)

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
	assert.NotNil(s.T(), fi.ProcessedAt)
	assert.Contains(s.T(), []string(fi.Tags), "auth")
	assert.Contains(s.T(), []string(fi.Tags), "security")
}

// TestSummarizeFileJob_Idempotency verifies that running the handler twice with the same
// payload updates the existing chunk and file-level summary rows in place.
func (s *SummarizeFileJobTestSuite) TestSummarizeFileJob_Idempotency() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	s.createGitHubIntegration(workspace.ID, user.ID)
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)
	s.seeder.CreateFileIndex(repo.ID, "util.go")

	goSrc := "func Helper() {}\n"
	job := jobs.NewSummarizeFileJobWithDeps(s.DB, &config.Config{}, newMockAIProvider(), &mockFileFetcher{content: goSrc}, nil)
	payload := jobs.SummarizeFilePayload{
		RepositoryID: repo.ID.String(),
		FilePath:     "util.go",
		GitSHA:       "sha_x",
	}

	assert.NoError(s.T(), job.Handler(s.Ctx, buildTask(payload)))

	var firstRunCount int64
	s.DB.Model(&models.ChunkEmbedding{}).Where("repo_id = ? AND repo_file_path = ?", repo.ID, "util.go").Count(&firstRunCount)
	assert.Greater(s.T(), firstRunCount, int64(0))

	assert.NoError(s.T(), job.Handler(s.Ctx, buildTask(payload)))

	var secondRunCount int64
	s.DB.Model(&models.ChunkEmbedding{}).Where("repo_id = ? AND repo_file_path = ?", repo.ID, "util.go").Count(&secondRunCount)
	assert.Equal(s.T(), firstRunCount, secondRunCount, "re-processing should upsert, not duplicate")
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
	job := jobs.NewSummarizeFileJobWithDeps(s.DB, &config.Config{}, newMockAIProvider(), &mockFileFetcher{content: secretCode}, nil)
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

// TestSummarizeFileJob_MarksRepositoryIndexedWhenAllFilesProcessed verifies that repository
// status transitions indexing -> indexed only when the final pending file completes.
func (s *SummarizeFileJobTestSuite) TestSummarizeFileJob_MarksRepositoryIndexedWhenAllFilesProcessed() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	s.createGitHubIntegration(workspace.ID, user.ID)
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	now := time.Now()
	s.seeder.CreateFileIndex(repo.ID, "done.go", func(fi *models.RepositoryFile) {
		fi.GitSHA = "sha_done"
		fi.ProcessedAt = &now
	})
	s.seeder.CreateFileIndex(repo.ID, "pending.go", func(fi *models.RepositoryFile) {
		fi.GitSHA = "sha_pending"
		fi.ProcessedAt = nil
	})

	assert.NoError(s.T(), s.repoRepo.UpdateStatus(s.Ctx, repo.ID, models.RepositoryStatusIndexing))

	enqueuer := &mockTaskEnqueuer{}
	job := jobs.NewSummarizeFileJobWithDeps(
		s.DB,
		&config.Config{},
		newMockAIProvider(),
		&mockFileFetcher{content: "func Pending() {}\n"},
		enqueuer,
	)

	payload := jobs.SummarizeFilePayload{
		RepositoryID: repo.ID.String(),
		FilePath:     "pending.go",
		GitSHA:       "sha_pending",
	}
	assert.NoError(s.T(), job.Handler(s.Ctx, buildTask(payload)))

	updated, err := s.repoRepo.GetByID(s.Ctx, repo.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), models.RepositoryStatusIndexed, updated.Status)
	assert.Len(s.T(), enqueuer.generateTopicsCalls, 1)
}

// TestSummarizeFileJob_NoChunksStillMarksProcessed verifies that files yielding no chunks
// are still marked processed so repository indexing can complete.
func (s *SummarizeFileJobTestSuite) TestSummarizeFileJob_NoChunksStillMarksProcessed() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	s.createGitHubIntegration(workspace.ID, user.ID)
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)
	s.seeder.CreateFileIndex(repo.ID, "empty.go")

	assert.NoError(s.T(), s.repoRepo.UpdateStatus(s.Ctx, repo.ID, models.RepositoryStatusIndexing))

	enqueuer := &mockTaskEnqueuer{}
	job := jobs.NewSummarizeFileJobWithDeps(
		s.DB,
		&config.Config{},
		newMockAIProvider(),
		&mockFileFetcher{content: ""},
		enqueuer,
	)
	payload := jobs.SummarizeFilePayload{
		RepositoryID: repo.ID.String(),
		FilePath:     "empty.go",
		GitSHA:       "sha_empty",
	}
	assert.NoError(s.T(), job.Handler(s.Ctx, buildTask(payload)))

	fi, err := s.fileIndexRepo.GetByRepoAndPath(s.Ctx, repo.ID, "empty.go")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), fi.ProcessedAt)

	updated, err := s.repoRepo.GetByID(s.Ctx, repo.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), models.RepositoryStatusIndexed, updated.Status)
	assert.Len(s.T(), enqueuer.generateTopicsCalls, 1)
}
