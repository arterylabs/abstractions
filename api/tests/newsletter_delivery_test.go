package tests

import (
	"abstractions/clients/email"
	"abstractions/config"
	"abstractions/db/models"
	"abstractions/repositories"
	newsletter_service "abstractions/services/newsletter"
	"abstractions/tests/testutils"
	"context"
	"encoding/json"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ---- delivery test suite ----

type NewsletterDeliveryTestSuite struct {
	testutils.TestSuite
	seeder         *testutils.Seeder
	newsletterRepo *repositories.NewsletterRepository
	repositoryRepo *repositories.RepositoryRepository
	workspaceRepo  *repositories.WorkspaceRepository
}

func TestNewsletterDelivery(t *testing.T) {
	suite.Run(t, new(NewsletterDeliveryTestSuite))
}

func (s *NewsletterDeliveryTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.newsletterRepo = repositories.NewNewsletterRepository(s.DB)
	s.repositoryRepo = repositories.NewRepositoryRepository(s.DB)
	s.workspaceRepo = repositories.NewWorkspaceRepository(s.DB)
}

func (s *NewsletterDeliveryTestSuite) seedNewsletter(repoID uuid.UUID, html string) *models.Newsletter {
	nl := &models.Newsletter{
		ID:          uuid.New(),
		RepoID:      repoID,
		HTMLContent: html,
	}
	err := s.newsletterRepo.Create(context.Background(), nl)
	s.Require().NoError(err)
	return nl
}

// TestSendNewsletter_RecipientsListIsRespected verifies that emails in repo.Recipients receive the newsletter.
func (s *NewsletterDeliveryTestSuite) TestSendNewsletter_RecipientsListIsRespected() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	// Set explicit recipients and disable workspace member delivery
	repo.Recipients = models.StringArray{"alice@example.com", "carol@example.com"}
	repo.SendToWorkspaceMembers = false
	err := s.repositoryRepo.Update(context.Background(), repo)
	s.Require().NoError(err)

	nl := s.seedNewsletter(repo.ID, "<html><body><p>Hello</p></body></html>")

	var mu sync.Mutex
	sentTo := []string{}
	mockSender := func(ctx context.Context, params email.AWSEmailParams) {
		mu.Lock()
		sentTo = append(sentTo, params.Email)
		mu.Unlock()
	}

	svc := newsletter_service.NewDeliveryServiceWithSender(
		&config.Config{Environment: &config.ENV{ClientURL: "http://localhost:3000"}},
		s.newsletterRepo, s.workspaceRepo, s.repositoryRepo,
		mockSender,
	)

	err = svc.SendNewsletter(s.Ctx, repo.ID, nl.ID, "Test Newsletter")
	assert.NoError(s.T(), err)

	assert.Len(s.T(), sentTo, 2)
	assert.Contains(s.T(), sentTo, "alice@example.com")
	assert.Contains(s.T(), sentTo, "carol@example.com")
}

// TestSendNewsletter_UpdatesNewsletterRow verifies that sent_at and recipient_count are persisted.
func (s *NewsletterDeliveryTestSuite) TestSendNewsletter_UpdatesNewsletterRow() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	repo.Recipients = models.StringArray{"dev@example.com"}
	repo.SendToWorkspaceMembers = false
	err := s.repositoryRepo.Update(context.Background(), repo)
	s.Require().NoError(err)

	nl := s.seedNewsletter(repo.ID, "<html><body><p>update test</p></body></html>")

	svc := newsletter_service.NewDeliveryServiceWithSender(
		&config.Config{Environment: &config.ENV{ClientURL: "http://localhost:3000"}},
		s.newsletterRepo, s.workspaceRepo, s.repositoryRepo,
		func(_ context.Context, _ email.AWSEmailParams) {}, // no-op sender
	)

	err = svc.SendNewsletter(s.Ctx, repo.ID, nl.ID, "Test Newsletter")
	assert.NoError(s.T(), err)

	updated, err := s.newsletterRepo.GetByID(s.Ctx, nl.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), updated.SentAt)
	assert.Equal(s.T(), 1, updated.RecipientCount)
}

// TestSendNewsletter_EmptyRecipientsSkipsSend verifies that no emails are sent when
// Recipients is empty and SendToWorkspaceMembers is false.
func (s *NewsletterDeliveryTestSuite) TestSendNewsletter_EmptyRecipientsSkipsSend() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID)

	repo.Recipients = models.StringArray{}
	repo.SendToWorkspaceMembers = false
	err := s.repositoryRepo.Update(context.Background(), repo)
	s.Require().NoError(err)

	nl := s.seedNewsletter(repo.ID, "<html><body><p>empty test</p></body></html>")

	sentCount := 0
	svc := newsletter_service.NewDeliveryServiceWithSender(
		&config.Config{Environment: &config.ENV{ClientURL: "http://localhost:3000"}},
		s.newsletterRepo, s.workspaceRepo, s.repositoryRepo,
		func(_ context.Context, _ email.AWSEmailParams) { sentCount++ },
	)

	err = svc.SendNewsletter(s.Ctx, repo.ID, nl.ID, "Test Newsletter")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, sentCount)
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
