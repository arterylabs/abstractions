package newsletter_test

import (
	"abstractions/repositories"
	"abstractions/tests/testutils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type NewsletterServiceTestSuite struct {
	testutils.TestSuite
	seeder         *testutils.Seeder
	newsletterRepo *repositories.NewsletterRepository
}

func TestNewsletterService(t *testing.T) {
	suite.Run(t, new(NewsletterServiceTestSuite))
}

func (s *NewsletterServiceTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.newsletterRepo = repositories.NewNewsletterRepository(s.DB)
}

func (s *NewsletterServiceTestSuite) TestGetNewsletter_ReturnsCorrectRecord() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)
	nl := s.seeder.CreateNewsletter(repo.ID)

	fetched, err := s.newsletterRepo.GetByID(s.Ctx, nl.ID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), fetched)
	assert.Equal(s.T(), nl.ID, fetched.ID)
	assert.Equal(s.T(), repo.ID, fetched.RepoID)
}

func (s *NewsletterServiceTestSuite) TestGetNewsletter_NotFound_ReturnsNil() {
	fetched, err := s.newsletterRepo.GetByID(s.Ctx, uuid.New())

	assert.NoError(s.T(), err)
	assert.Nil(s.T(), fetched, "non-existent newsletter should return nil, not error")
}

func (s *NewsletterServiceTestSuite) TestListByRepoID_ReturnsAllNewsletters() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	s.seeder.CreateNewsletter(repo.ID)
	s.seeder.CreateNewsletter(repo.ID)

	// Create a newsletter for a different repo to ensure scoping
	otherRepo := s.seeder.CreateRepository(workspace.ID, user.ID)
	s.seeder.CreateNewsletter(otherRepo.ID)

	results, err := s.newsletterRepo.ListByRepoID(s.Ctx, repo.ID)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), results, 2, "should only return newsletters for the specified repo")
}

func (s *NewsletterServiceTestSuite) TestCount_IncrementsOnCreate() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	before, err := s.newsletterRepo.Count(s.Ctx)
	s.Require().NoError(err)

	s.seeder.CreateNewsletter(repo.ID)

	after, err := s.newsletterRepo.Count(s.Ctx)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), before+1, after)
}
