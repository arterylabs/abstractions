package tests

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/server/responses"
	"abstractions/tests/testutils"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ---- newsletter history test suite ----

type NewsletterHistoryTestSuite struct {
	testutils.TestSuite
	nlRepo  *repositories.NewsletterRepository
	repoRep *repositories.RepositoryRepository
	seeder  *testutils.Seeder
}

func (s *NewsletterHistoryTestSuite) SetupTest() {
	s.nlRepo = repositories.NewNewsletterRepository(s.DB)
	s.repoRep = repositories.NewRepositoryRepository(s.DB)
	s.seeder = testutils.NewSeeder(s.DB)
}

func TestNewsletterHistory(t *testing.T) {
	suite.Run(t, new(NewsletterHistoryTestSuite))
}

// seedNewsletter creates a newsletter row directly for testing.
func (s *NewsletterHistoryTestSuite) seedNewsletter(repoID uuid.UUID, recipientCount int) *models.Newsletter {
	nl := &models.Newsletter{
		ID:             uuid.New(),
		RepoID:         repoID,
		HTMLContent:    "<p>hello</p>",
		RecipientCount: recipientCount,
	}
	now := time.Now()
	nl.SentAt = &now
	if err := s.nlRepo.Create(context.Background(), nl); err != nil {
		panic(err)
	}
	return nl
}

// TestListNewsletters_ScopedToRepo ensures newsletters from another repo are not returned.
func (s *NewsletterHistoryTestSuite) TestListNewsletters_ScopedToRepo() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repoA := s.seeder.CreateRepository(workspace.ID, user.ID)
	repoB := s.seeder.CreateRepository(workspace.ID, user.ID)

	s.seedNewsletter(repoA.ID, 50)
	s.seedNewsletter(repoA.ID, 50)
	s.seedNewsletter(repoB.ID, 50)

	newsletters, err := s.nlRepo.ListByRepoID(context.Background(), repoA.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), newsletters, 2)
	for _, n := range newsletters {
		assert.Equal(s.T(), repoA.ID, n.RepoID)
	}
}

// TestGetNewsletter_Detail verifies the detail response includes HTML.
func (s *NewsletterHistoryTestSuite) TestGetNewsletter_Detail() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)
	nl := s.seedNewsletter(repo.ID, 20)

	fetched, err := s.nlRepo.GetByID(context.Background(), nl.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), fetched)

	resp := responses.NewNewsletterDetailResponse(fetched)
	assert.Equal(s.T(), "<p>hello</p>", resp.HTMLContent)
	// Topic is nil when no newsletter_topics row exists for this newsletter.
	assert.Nil(s.T(), resp.Topic)
}
