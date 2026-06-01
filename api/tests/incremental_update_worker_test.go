package tests

import (
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
	pgvector "github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ---- helpers ----

func buildDeleteFileTask(repoID uuid.UUID, filePath string) *asynq.Task {
	type p struct {
		RepositoryID string `json:"repositoryId"`
		FilePath     string `json:"filePath"`
	}
	b, _ := json.Marshal(p{RepositoryID: repoID.String(), FilePath: filePath})
	return asynq.NewTask(jobs.DeleteFileJobID, b)
}

// ---- test suite ----

type IncrementalUpdateWorkerTestSuite struct {
	testutils.TestSuite
	seeder        *testutils.Seeder
	intSvc        *integration_service.Service
	fileIndexRepo *repositories.RepositoryFileRepository
	chunkRepo     *repositories.ChunkEmbeddingRepository
}

func TestIncrementalUpdateWorker(t *testing.T) {
	suite.Run(t, new(IncrementalUpdateWorkerTestSuite))
}

func (s *IncrementalUpdateWorkerTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.intSvc = integration_service.NewService(s.DB)
	s.fileIndexRepo = repositories.NewRepositoryFileRepository(s.DB)
	s.chunkRepo = repositories.NewChunkEmbeddingRepository(s.DB)
}

func (s *IncrementalUpdateWorkerTestSuite) createIntegration(workspaceID, userID uuid.UUID) {
	_, err := s.intSvc.Create(s.Ctx, &integration_service.CreateIntegrationRequest{
		Type:        github_integration.IntegrationTypeGitHub,
		WorkspaceID: workspaceID,
		Token:       "ghp_test",
		CreatedBy:   userID,
	})
	s.Require().NoError(err)
}

func (s *IncrementalUpdateWorkerTestSuite) seedFileIndex(repoID uuid.UUID, filePath, sha string) {
	now := time.Now()
	fi := &models.RepositoryFile{
		ID:          uuid.New(),
		RepoID:      repoID,
		Path:        filePath,
		GitSHA:      sha,
		ProcessedAt: &now,
	}
	err := s.fileIndexRepo.Upsert(context.Background(), fi)
	s.Require().NoError(err)
}

// TestDeleteFileJob_RemovesFileIndexAndChunks verifies that the DeleteFileJob
// deletes both the file_index row and all associated chunk_embeddings.
func (s *IncrementalUpdateWorkerTestSuite) TestDeleteFileJob_RemovesFileIndexAndChunks() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	s.createIntegration(workspace.ID, user.ID)
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	const filePath = "internal/handler.go"
	s.seedFileIndex(repo.ID, filePath, "deadbeef")

	vec := make([]float32, 1536)
	vec[0] = 0.7
	err := s.chunkRepo.Upsert(s.Ctx, []models.ChunkEmbedding{
		{ID: uuid.New(), RepoID: repo.ID, RepoFilePath: filePath, ChunkIndex: 0, Summary: "handles requests", Embedding: pgvector.NewVector(vec)},
		{ID: uuid.New(), RepoID: repo.ID, RepoFilePath: filePath, ChunkIndex: 1, Summary: "error handling", Embedding: pgvector.NewVector(vec)},
	})
	s.Require().NoError(err)

	job := jobs.NewDeleteFileJob(s.DB, &config.Config{})
	err = job.Handler(s.Ctx, buildDeleteFileTask(repo.ID, filePath))
	assert.NoError(s.T(), err)

	// file_index must be gone
	fi, err := s.fileIndexRepo.GetByRepoAndPath(s.Ctx, repo.ID, filePath)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), fi, "file_index must be deleted")
}
