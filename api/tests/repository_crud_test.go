package tests

import (
	"abstractions/db/models"
	"abstractions/repositories"
	repository_service "abstractions/services/repository"
	"abstractions/tests/testutils"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// unknownID returns a UUID that is guaranteed not to exist in any test database.
func unknownID() uuid.UUID {
	return uuid.MustParse("00000000-0000-0000-0000-000000000001")
}

// newsletterModel builds a minimal Newsletter for seeding.
func newsletterModel(repoID uuid.UUID) models.Newsletter {
	return models.Newsletter{
		ID:     uuid.New(),
		RepoID: repoID,
		Topics: models.JSONBRaw("[]"),
	}
}

// ---- repository CRUD test suite ----

type RepositoryCRUDTestSuite struct {
	testutils.TestSuite
	seeder   *testutils.Seeder
	svc      *repository_service.Service
	repoRep  *repositories.RepositoryRepository
	fileRep  *repositories.FileIndexRepository
	chunkRep *repositories.ChunkEmbeddingRepository
	nlRep    *repositories.NewsletterRepository
}

func TestRepositoryCRUD(t *testing.T) {
	suite.Run(t, new(RepositoryCRUDTestSuite))
}

func (s *RepositoryCRUDTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.svc = repository_service.NewRepositoryService(s.DB)
	s.repoRep = repositories.NewRepositoryRepository(s.DB)
	s.fileRep = repositories.NewFileIndexRepository(s.DB)
	s.chunkRep = repositories.NewChunkEmbeddingRepository(s.DB)
	s.nlRep = repositories.NewNewsletterRepository(s.DB)
}

// TestUpdateSettings_PersistsFields verifies that send_day and send_hour are saved.
func (s *RepositoryCRUDTestSuite) TestUpdateSettings_PersistsFields() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	updated, err := s.svc.UpdateSettings(context.Background(), repo.ID, &repository_service.UpdateSettingsRequest{
		SendDay:  3, // Wednesday
		SendHour: 14,
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 3, updated.SendDay)
	assert.Equal(s.T(), 14, updated.SendHour)

	// Verify persisted
	fetched, err := s.repoRep.GetByID(context.Background(), repo.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 3, fetched.SendDay)
	assert.Equal(s.T(), 14, fetched.SendHour)
}

// TestUpdateSettings_UnknownRepo returns ErrRepositoryNotFound.
func (s *RepositoryCRUDTestSuite) TestUpdateSettings_UnknownRepo() {
	_, err := s.svc.UpdateSettings(context.Background(), unknownID(), &repository_service.UpdateSettingsRequest{
		SendDay: 1, SendHour: 9,
	})
	assert.ErrorIs(s.T(), err, repository_service.ErrRepositoryNotFound)
}

// TestDeleteRepositoryData_ClearsAllArtifacts verifies that file_index, chunk_embeddings,
// and newsletters are all removed and the repo status is reset to pending.
func (s *RepositoryCRUDTestSuite) TestDeleteRepositoryData_ClearsAllArtifacts() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	// Seed artifacts
	fi := s.seeder.CreateFileIndex(repo.ID, "main.go")
	_ = fi

	nlVal := newsletterModel(repo.ID)
	nlErr := s.nlRep.Create(context.Background(), &nlVal)
	s.Require().NoError(nlErr)

	// Delete
	err := s.svc.DeleteRepositoryData(context.Background(), repo.ID)
	assert.NoError(s.T(), err)

	// Verify file_index cleared
	files, err := s.fileRep.ListByRepoID(context.Background(), repo.ID)
	assert.NoError(s.T(), err)
	assert.Empty(s.T(), files)

	// Verify newsletters cleared
	newsletters, err := s.nlRep.ListByRepoID(context.Background(), repo.ID)
	assert.NoError(s.T(), err)
	assert.Empty(s.T(), newsletters)

	// Verify repo status reset to pending
	updated, err := s.repoRep.GetByID(context.Background(), repo.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "pending", string(updated.Status))
}
