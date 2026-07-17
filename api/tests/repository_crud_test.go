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
	}
}

// ---- repository CRUD test suite ----

type RepositoryCRUDTestSuite struct {
	testutils.TestSuite
	seeder   *testutils.Seeder
	svc      *repository_service.Service
	repoRep  *repositories.RepositoryRepository
	fileRep  *repositories.RepositoryFileRepository
	chunkRep *repositories.ChunkEmbeddingRepository
	nlRep    *repositories.NewsletterRepository
}

func TestRepositoryCRUD(t *testing.T) {
	suite.Run(t, new(RepositoryCRUDTestSuite))
}

func (s *RepositoryCRUDTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.svc = repository_service.NewService(s.DB)
	s.repoRep = repositories.NewRepositoryRepository(s.DB)
	s.fileRep = repositories.NewRepositoryFileRepository(s.DB)
	s.chunkRep = repositories.NewChunkEmbeddingRepository(s.DB)
	s.nlRep = repositories.NewNewsletterRepository(s.DB)
}

// TestUpdateSettings_PersistsFields verifies that send_day and send_hour are saved.
func (s *RepositoryCRUDTestSuite) TestUpdateSettings_PersistsFields() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	updated, err := s.svc.UpdateSettings(context.Background(), workspace.ID, repo.ID, &repository_service.UpdateSettingsRequest{
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
	_, err := s.svc.UpdateSettings(context.Background(), unknownID(), unknownID(), &repository_service.UpdateSettingsRequest{
		SendDay: 1, SendHour: 9,
	})
	assert.ErrorIs(s.T(), err, repository_service.ErrRepositoryNotFound)
}

// TestDeleteRepository_SoftDeletesAndClearsArtifacts verifies that the repository is
// soft-deleted and all derived artifacts (file_index, newsletters) are removed.
func (s *RepositoryCRUDTestSuite) TestDeleteRepository_SoftDeletesAndClearsArtifacts() {
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
	err := s.svc.Delete(context.Background(), workspace.ID, repo.ID, user.ID)
	assert.NoError(s.T(), err)

	// Verify repo is soft-deleted (GetByID filters deleted_at IS NULL)
	deleted, err := s.repoRep.GetByID(context.Background(), repo.ID)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), deleted)

	// Verify file_index cleared
	files, err := s.fileRep.ListByRepoID(context.Background(), repo.ID)
	assert.NoError(s.T(), err)
	assert.Empty(s.T(), files)

	// Verify newsletters cleared
	newsletters, err := s.nlRep.ListByRepoID(context.Background(), repo.ID)
	assert.NoError(s.T(), err)
	assert.Empty(s.T(), newsletters)
}
