package tests

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/tests/testutils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FileIndexRepositoryTestSuite struct {
	testutils.TestSuite
	seeder *testutils.Seeder
	repo   *repositories.FileIndexRepository
}

func TestFileIndexRepository(t *testing.T) {
	suite.Run(t, new(FileIndexRepositoryTestSuite))
}

func (s *FileIndexRepositoryTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.repo = repositories.NewFileIndexRepository(s.DB)
}

// TestUpsert_Create inserts a new FileIndex entry and confirms it can be retrieved.
func (s *FileIndexRepositoryTestSuite) TestUpsert_Create() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repository := s.seeder.CreateRepository(workspace.ID, user.ID)

	now := time.Now()
	fi := &models.FileIndex{
		RepoID:          repository.ID,
		FilePath:        "cmd/main.go",
		GitSHA:          "sha_abc",
		DomainTags:      models.StringArray{},
		LastProcessedAt: &now,
		PineconeIDs:     models.StringArray{},
	}

	err := s.repo.Upsert(s.Ctx, fi)
	assert.NoError(s.T(), err)

	got, err := s.repo.GetByRepoAndPath(s.Ctx, repository.ID, "cmd/main.go")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), got)
	assert.Equal(s.T(), "sha_abc", got.GitSHA)
}

// TestListByRepoID_ReturnsOnlyTargetRepo ensures workspace/repo isolation.
func (s *FileIndexRepositoryTestSuite) TestListByRepoID_ReturnsOnlyTargetRepo() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo1 := s.seeder.CreateRepository(workspace.ID, user.ID)
	repo2 := s.seeder.CreateRepository(workspace.ID, user.ID)

	s.seeder.CreateFileIndex(repo1.ID, "main.go")
	s.seeder.CreateFileIndex(repo1.ID, "utils.go")
	s.seeder.CreateFileIndex(repo2.ID, "other.go")

	files, err := s.repo.ListByRepoID(s.Ctx, repo1.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), files, 2)
}

// TestDeleteByRepoID removes all entries for a repository.
func (s *FileIndexRepositoryTestSuite) TestDeleteByRepoID() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	s.seeder.CreateFileIndex(repo.ID, "a.go")
	s.seeder.CreateFileIndex(repo.ID, "b.go")

	err := s.repo.DeleteByRepoID(s.Ctx, repo.ID)
	assert.NoError(s.T(), err)

	files, err := s.repo.ListByRepoID(s.Ctx, repo.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), files, 0)
}

// TestDeleteByPath removes only the specified file, leaving others intact.
func (s *FileIndexRepositoryTestSuite) TestDeleteByPath() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	s.seeder.CreateFileIndex(repo.ID, "keep.go")
	s.seeder.CreateFileIndex(repo.ID, "delete.go")

	err := s.repo.DeleteByPath(s.Ctx, repo.ID, "delete.go")
	assert.NoError(s.T(), err)

	files, err := s.repo.ListByRepoID(s.Ctx, repo.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), files, 1)
	assert.Equal(s.T(), "keep.go", files[0].FilePath)
}
