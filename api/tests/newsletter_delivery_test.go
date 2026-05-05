package tests

import (
	"abstractions/clients/email"
	"abstractions/config"
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/server/handlers"
	newsletter_service "abstractions/services/newsletter"
	"abstractions/tests/testutils"
	jwt_utils "abstractions/utils/jwt_utils"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ---- delivery test suite ----

type NewsletterDeliveryTestSuite struct {
	testutils.TestSuite
	seeder         *testutils.Seeder
	newsletterRepo *repositories.NewsletterRepository
	teamMemberRepo *repositories.TeamMemberRepository
	repositoryRepo *repositories.RepositoryRepository
}

func TestNewsletterDelivery(t *testing.T) {
	suite.Run(t, new(NewsletterDeliveryTestSuite))
}

func (s *NewsletterDeliveryTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.newsletterRepo = repositories.NewNewsletterRepository(s.DB)
	s.teamMemberRepo = repositories.NewTeamMemberRepository(s.DB)
	s.repositoryRepo = repositories.NewRepositoryRepository(s.DB)
}

func (s *NewsletterDeliveryTestSuite) seedNewsletter(repoID uuid.UUID, html string) *models.Newsletter {
	nl := &models.Newsletter{
		ID:          uuid.New(),
		RepoID:      repoID,
		Topics:      models.JSONBRaw("[]"),
		HTMLContent: html,
	}
	err := s.newsletterRepo.Create(context.Background(), nl)
	s.Require().NoError(err)
	return nl
}

func (s *NewsletterDeliveryTestSuite) seedTeamMember(repoID uuid.UUID, emailAddr string, subscribed bool) *models.TeamMember {
	tm := &models.TeamMember{
		ID:         uuid.New(),
		RepoID:     repoID,
		Email:      emailAddr,
		Subscribed: subscribed,
		AddedAt:    time.Now(),
	}
	err := s.DB.Exec(
		`INSERT INTO team_members (id, repo_id, email, subscribed, added_at) VALUES (?, ?, ?, ?, ?)`,
		tm.ID, tm.RepoID, tm.Email, tm.Subscribed, tm.AddedAt,
	).Error
	s.Require().NoError(err)
	return tm
}

// TestSendNewsletter_OnlySubscribedMembersReceiveEmail verifies that unsubscribed
// members are excluded from the delivery.
func (s *NewsletterDeliveryTestSuite) TestSendNewsletter_OnlySubscribedMembersReceiveEmail() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)
	nl := s.seedNewsletter(repo.ID, "<html><body><p>Hello</p></body></html>")

	s.seedTeamMember(repo.ID, "alice@example.com", true)
	s.seedTeamMember(repo.ID, "bob@example.com", false) // unsubscribed
	s.seedTeamMember(repo.ID, "carol@example.com", true)

	var mu sync.Mutex
	sentTo := []string{}
	mockSender := func(ctx context.Context, params email.AWSEmailParams) {
		mu.Lock()
		sentTo = append(sentTo, params.Email)
		mu.Unlock()
	}

	svc := newsletter_service.NewDeliveryServiceWithSender(
		&config.Config{Environment: &config.ENV{ClientURL: "http://localhost:3000"}},
		s.newsletterRepo, s.teamMemberRepo, s.repositoryRepo,
		mockSender,
	)

	err := svc.SendNewsletter(s.Ctx, repo.ID, nl.ID)
	assert.NoError(s.T(), err)

	assert.Len(s.T(), sentTo, 2, "only subscribed members should receive email")
	assert.Contains(s.T(), sentTo, "alice@example.com")
	assert.NotContains(s.T(), sentTo, "bob@example.com")
}

// TestSendNewsletter_UpdatesNewsletterRow verifies that sent_at and recipient_count are persisted.
func (s *NewsletterDeliveryTestSuite) TestSendNewsletter_UpdatesNewsletterRow() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)
	nl := s.seedNewsletter(repo.ID, "<html><body><p>update test</p></body></html>")

	s.seedTeamMember(repo.ID, "dev@example.com", true)

	svc := newsletter_service.NewDeliveryServiceWithSender(
		&config.Config{Environment: &config.ENV{ClientURL: "http://localhost:3000"}},
		s.newsletterRepo, s.teamMemberRepo, s.repositoryRepo,
		func(_ context.Context, _ email.AWSEmailParams) {}, // no-op sender
	)

	err := svc.SendNewsletter(s.Ctx, repo.ID, nl.ID)
	assert.NoError(s.T(), err)

	updated, err := s.newsletterRepo.GetByID(s.Ctx, nl.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), updated.SentAt)
	assert.Equal(s.T(), 1, updated.RecipientCount)
}

// TestSendBulkEmail_RespectsBatchSize verifies that bulk sender processes all recipients
// even when there are more than one batch.
func (s *NewsletterDeliveryTestSuite) TestSendBulkEmail_RespectsBatchSize() {
	const totalRecipients = 16 // more than one batch of 14

	var mu sync.Mutex
	sentCount := 0
	mockSender := func(_ context.Context, _ email.AWSEmailParams) {
		mu.Lock()
		sentCount++
		mu.Unlock()
	}

	params := make([]email.AWSEmailParams, totalRecipients)
	for i := range params {
		params[i] = email.AWSEmailParams{Email: "user@example.com", Subject: "test", BodyHTML: "<p>hi</p>", BodyText: "hi"}
	}

	err := email.SendBulkEmailUsingSES(context.Background(), params, mockSender)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), totalRecipients, sentCount)
}

// TestUnsubscribeToken_RoundTrip verifies that GenerateUnsubscribeToken and
// ParseUnsubscribeToken correctly serialise/deserialise the team member ID.
func (s *NewsletterDeliveryTestSuite) TestUnsubscribeToken_RoundTrip() {
	// requires a SIGNING_KEY to be set in the environment config
	original := uuid.New()

	// Bootstrap minimal config so the jwt_utils package uses a test key
	config.Environment = &config.ENV{SigningKey: "test-signing-key-for-unit-tests"}

	token, err := jwt_utils.GenerateUnsubscribeToken(original)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), token)

	parsed, err := jwt_utils.ParseUnsubscribeToken(token)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), original, parsed)
}

// TestUnsubscribeEndpoint_SetsSubscribedFalse verifies the HTTP handler sets subscribed=false.
func (s *NewsletterDeliveryTestSuite) TestUnsubscribeEndpoint_SetsSubscribedFalse() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)
	tm := s.seedTeamMember(repo.ID, "sub@example.com", true)

	config.Environment = &config.ENV{SigningKey: "test-signing-key-for-unit-tests"}

	token, err := jwt_utils.GenerateUnsubscribeToken(tm.ID)
	s.Require().NoError(err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/unsubscribe?token="+token, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/unsubscribe")

	h := handlers.NewUnsubscribeHandler(s.DB)
	err = h.Unsubscribe(c)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), http.StatusOK, rec.Code)
	assert.Contains(s.T(), rec.Body.String(), "unsubscribed")

	updated, err := s.teamMemberRepo.GetByID(s.Ctx, tm.ID)
	assert.NoError(s.T(), err)
	assert.False(s.T(), updated.Subscribed)
}

// ---- helpers ----

// newsletterTopicJSON marshals topic titles into JSONBRaw for test newsletters.
func newsletterTopicJSON(titles []string) models.JSONBRaw {
	type entry struct {
		Title string `json:"title"`
	}
	entries := make([]entry, len(titles))
	for i, t := range titles {
		entries[i] = entry{Title: t}
	}
	b, _ := json.Marshal(entries)
	return models.JSONBRaw(b)
}
