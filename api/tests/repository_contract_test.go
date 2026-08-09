package tests

import (
	"abstractions/tests/testutils"
	"testing"

	"abstractions/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Soft/permanent delete semantics are security-sensitive (recovery, purge).
// Basic CRUD is exercised by service-layer suites.

type RepositoryContractTestSuite struct {
	testutils.TestSuite
	seeder *testutils.Seeder
}

func TestRepositoryContract(t *testing.T) {
	suite.Run(t, new(RepositoryContractTestSuite))
}

func (s *RepositoryContractTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
}

func (s *RepositoryContractTestSuite) TestRepositorySoftAndPermanentDelete() {
	user := s.seeder.CreateUser()
	workspace := s.seeder.CreateWorkspace()
	repoModel := s.seeder.CreateRepository(workspace.ID, user.ID)
	repo := repositories.NewRepositoryRepository(s.DB)

	assert.NoError(s.T(), repo.Delete(s.Ctx, repoModel.ID, user.ID))

	deleted, err := repo.GetByID(s.Ctx, repoModel.ID)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), deleted)

	includingDeleted, err := repo.GetByIDIncludingDeleted(s.Ctx, repoModel.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), includingDeleted)
	assert.NotNil(s.T(), includingDeleted.DeletedAt)

	assert.NoError(s.T(), repo.PermanentlyDelete(s.Ctx, repoModel.ID))
	afterHardDelete, err := repo.GetByIDIncludingDeleted(s.Ctx, repoModel.ID)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), afterHardDelete)
}
