package tests

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/server/responses"
	"abstractions/tests/testutils"
	"context"
	"encoding/json"
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
func (s *NewsletterHistoryTestSuite) seedNewsletter(repoID uuid.UUID, recipientCount, openCount, clickCount int) *models.Newsletter {
	nl := &models.Newsletter{
		ID:             uuid.New(),
		RepoID:         repoID,
		Topics:         models.JSONBRaw(`[{"title":"Feature X"}]`),
		HTMLContent:    "<p>hello</p>",
		RecipientCount: recipientCount,
		OpenCount:      openCount,
		ClickCount:     clickCount,
	}
	now := time.Now()
	nl.SentAt = &now
	if err := s.nlRepo.Create(context.Background(), nl); err != nil {
		panic(err)
	}
	return nl
}

// TestListNewsletters_OpenRateClickRate verifies open_rate and click_rate calculations.
func (s *NewsletterHistoryTestSuite) TestListNewsletters_OpenRateClickRate() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	s.seedNewsletter(repo.ID, 100, 40, 10)

	newsletters, err := s.nlRepo.ListByRepoID(context.Background(), repo.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), newsletters, 1)

	resp := responses.NewNewsletterSummaryResponse(newsletters[0])
	assert.InDelta(s.T(), 0.40, resp.OpenRate, 0.001)
	assert.InDelta(s.T(), 0.10, resp.ClickRate, 0.001)
}

// TestListNewsletters_ZeroRecipientCount returns 0.0 rates (no division by zero).
func (s *NewsletterHistoryTestSuite) TestListNewsletters_ZeroRecipientCount() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	s.seedNewsletter(repo.ID, 0, 0, 0)

	newsletters, err := s.nlRepo.ListByRepoID(context.Background(), repo.ID)
	assert.NoError(s.T(), err)

	resp := responses.NewNewsletterSummaryResponse(newsletters[0])
	assert.Equal(s.T(), 0.0, resp.OpenRate)
	assert.Equal(s.T(), 0.0, resp.ClickRate)
}

// TestListNewsletters_ScopedToRepo ensures newsletters from another repo are not returned.
func (s *NewsletterHistoryTestSuite) TestListNewsletters_ScopedToRepo() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repoA := s.seeder.CreateRepository(workspace.ID, user.ID)
	repoB := s.seeder.CreateRepository(workspace.ID, user.ID)

	s.seedNewsletter(repoA.ID, 50, 10, 5)
	s.seedNewsletter(repoA.ID, 50, 20, 5)
	s.seedNewsletter(repoB.ID, 50, 5, 1)

	newsletters, err := s.nlRepo.ListByRepoID(context.Background(), repoA.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), newsletters, 2)
	for _, n := range newsletters {
		assert.Equal(s.T(), repoA.ID, n.RepoID)
	}
}

// TestGetNewsletter_Detail verifies the detail response includes topics and HTML.
func (s *NewsletterHistoryTestSuite) TestGetNewsletter_Detail() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)
	nl := s.seedNewsletter(repo.ID, 20, 8, 2)

	fetched, err := s.nlRepo.GetByID(context.Background(), nl.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), fetched)

	resp := responses.NewNewsletterDetailResponse(fetched)
	assert.Equal(s.T(), "<p>hello</p>", resp.HTMLContent)

	var topics []map[string]interface{}
	assert.NoError(s.T(), json.Unmarshal(resp.Topics, &topics))
	assert.Len(s.T(), topics, 1)
}
