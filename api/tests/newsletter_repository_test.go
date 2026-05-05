package tests

import (
	"abstractions/repositories"
	"abstractions/tests/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type NewsletterRepositoryTestSuite struct {
	testutils.TestSuite
	seeder *testutils.Seeder
	repo   *repositories.NewsletterRepository
}

func TestNewsletterRepository(t *testing.T) {
	suite.Run(t, new(NewsletterRepositoryTestSuite))
}

func (s *NewsletterRepositoryTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.repo = repositories.NewNewsletterRepository(s.DB)
}

// TestCreate_GetByID creates a newsletter and retrieves it by ID.
func (s *NewsletterRepositoryTestSuite) TestCreate_GetByID() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repository := s.seeder.CreateRepository(workspace.ID, user.ID)

	n := s.seeder.CreateNewsletter(repository.ID)
	assert.NotNil(s.T(), n)

	got, err := s.repo.GetByID(s.Ctx, n.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), got)
	assert.Equal(s.T(), n.ID, got.ID)
	assert.Equal(s.T(), repository.ID, got.RepoID)
}

// TestListByRepoID_ReturnsOnlyTargetRepo ensures only newsletters for the given repo are returned.
func (s *NewsletterRepositoryTestSuite) TestListByRepoID_ReturnsOnlyTargetRepo() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo1 := s.seeder.CreateRepository(workspace.ID, user.ID)
	repo2 := s.seeder.CreateRepository(workspace.ID, user.ID)

	s.seeder.CreateNewsletter(repo1.ID)
	s.seeder.CreateNewsletter(repo1.ID)
	s.seeder.CreateNewsletter(repo2.ID)

	newsletters, err := s.repo.ListByRepoID(s.Ctx, repo1.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), newsletters, 2)
}

// TestUpdate persists changes to an existing newsletter.
func (s *NewsletterRepositoryTestSuite) TestUpdate() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repository := s.seeder.CreateRepository(workspace.ID, user.ID)

	n := s.seeder.CreateNewsletter(repository.ID)

	n.HTMLContent = "<p>Updated content</p>"
	n.RecipientCount = 42
	err := s.repo.Update(s.Ctx, n)
	assert.NoError(s.T(), err)

	got, err := s.repo.GetByID(s.Ctx, n.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "<p>Updated content</p>", got.HTMLContent)
	assert.Equal(s.T(), 42, got.RecipientCount)
}
