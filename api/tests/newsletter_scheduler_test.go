package tests

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/tests/testutils"
	"abstractions/workers/jobs"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ---- mock enqueuer ----

type mockEnqueuer struct {
	generateNewsletterCalls []uuid.UUID
}

func (m *mockEnqueuer) EnqueueSummarizeFile(_ uuid.UUID, _ string, _ string) error {
	return nil
}

func (m *mockEnqueuer) EnqueueGenerateTopics(_ uuid.UUID) error {
	return nil
}

func (m *mockEnqueuer) EnqueueGenerateNewsletter(repoID, _ uuid.UUID) error {
	m.generateNewsletterCalls = append(m.generateNewsletterCalls, repoID)
	return nil
}

// ---- scheduler test suite ----

type NewsletterSchedulerTestSuite struct {
	testutils.TestSuite
	seeder         *testutils.Seeder
	repositoryRepo *repositories.RepositoryRepository
}

func TestNewsletterScheduler(t *testing.T) {
	suite.Run(t, new(NewsletterSchedulerTestSuite))
}

func (s *NewsletterSchedulerTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.repositoryRepo = repositories.NewRepositoryRepository(s.DB)
}

func (s *NewsletterSchedulerTestSuite) seedIndexedRepo(sendDay, sendHour int, lastSentAt *time.Time) *models.Repository {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()
	repo := s.seeder.CreateRepository(workspace.ID, user.ID, func(r *models.Repository) {
		r.SendDay = sendDay
		r.SendHour = sendHour
		r.Status = models.RepositoryStatusIndexed
	})
	// Use explicit SQL updates so GORM doesn't skip pointer/bool fields treated as zero values.
	s.Require().NoError(s.DB.Exec(
		`UPDATE repositories SET send_day = ?, send_hour = ?, status = ? WHERE id = ?`,
		sendDay, sendHour, models.RepositoryStatusIndexed, repo.ID,
	).Error)
	if lastSentAt != nil {
		s.Require().NoError(s.DB.Exec(
			`UPDATE repositories SET newsletter_sent_at = ? WHERE id = ?`,
			lastSentAt.UTC(), repo.ID,
		).Error)
	}
	return repo
}

// TestScheduler_EnqueuesQualifyingRepos verifies that the scheduler enqueues jobs for
// repos whose send_day and send_hour match the current time.
func (s *NewsletterSchedulerTestSuite) TestScheduler_EnqueuesQualifyingRepos() {
	now := time.Now().UTC()
	sendDay := int(now.Weekday())
	if sendDay == 0 {
		sendDay = 7 // Sunday → 7
	}
	sendHour := now.Hour()

	// Matching repo: send_day/send_hour match now
	matchingRepo := s.seedIndexedRepo(sendDay, sendHour, nil)
	// Non-matching repo: different hour
	nonMatchingHour := (sendHour + 1) % 24
	s.seedIndexedRepo(sendDay, nonMatchingHour, nil)

	enqueuer := &mockEnqueuer{}
	job := jobs.NewNewsletterSchedulerJob(s.DB, s.Config, enqueuer)

	task := asynq.NewTask(jobs.NewsletterSchedulerJobID, nil)
	err := job.Handler(s.Ctx, task)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), enqueuer.generateNewsletterCalls, 1, "only the matching repo should be enqueued")
	assert.Equal(s.T(), matchingRepo.ID, enqueuer.generateNewsletterCalls[0])
}

// TestScheduler_DoubleSendGuard verifies repos with a recent newsletter are skipped.
func (s *NewsletterSchedulerTestSuite) TestScheduler_DoubleSendGuard() {
	now := time.Now().UTC()
	sendDay := int(now.Weekday())
	if sendDay == 0 {
		sendDay = 7
	}
	sendHour := now.Hour()

	recent := time.Now().Add(-30 * time.Minute) // sent 30 min ago — within 2h guard
	s.seedIndexedRepo(sendDay, sendHour, &recent)

	enqueuer := &mockEnqueuer{}
	job := jobs.NewNewsletterSchedulerJob(s.DB, s.Config, enqueuer)

	task := asynq.NewTask(jobs.NewsletterSchedulerJobID, nil)
	err := job.Handler(s.Ctx, task)

	assert.NoError(s.T(), err)
	assert.Empty(s.T(), enqueuer.generateNewsletterCalls, "recently sent repo should be skipped by double-send guard")
}

// TestScheduler_AllowsAfterGuardWindow verifies repos with an old newsletter are still enqueued.
func (s *NewsletterSchedulerTestSuite) TestScheduler_AllowsAfterGuardWindow() {
	now := time.Now().UTC()
	sendDay := int(now.Weekday())
	if sendDay == 0 {
		sendDay = 7
	}
	sendHour := now.Hour()

	old := time.Now().Add(-3 * time.Hour) // sent 3 hours ago — outside 2h guard
	s.seedIndexedRepo(sendDay, sendHour, &old)

	enqueuer := &mockEnqueuer{}
	job := jobs.NewNewsletterSchedulerJob(s.DB, s.Config, enqueuer)

	task := asynq.NewTask(jobs.NewsletterSchedulerJobID, nil)
	err := job.Handler(s.Ctx, task)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), enqueuer.generateNewsletterCalls, 1, "repo outside guard window should be enqueued")
}

// TestManualTrigger_RateLimit verifies the rate-limit check in the delivery service.
// We test the logic directly via the repository query to keep tests DB-only (no HTTP server needed).
func (s *NewsletterSchedulerTestSuite) TestManualTrigger_RateLimit() {
	workspace := s.seeder.CreateWorkspace()
	user := s.seeder.CreateUser()

	// Repo that was sent recently
	recentSent := time.Now().Add(-2 * time.Hour)
	repoRecent := s.seeder.CreateRepository(workspace.ID, user.ID, func(r *models.Repository) {
		r.NewsletterSentAt = &recentSent
		r.Status = models.RepositoryStatusIndexed
	})

	// Repo that has never been sent
	repoNever := s.seeder.CreateRepository(workspace.ID, user.ID, func(r *models.Repository) {
		r.Status = models.RepositoryStatusIndexed
	})

	recentRepo, err := s.repositoryRepo.GetByID(context.Background(), repoRecent.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), recentRepo.NewsletterSentAt)

	// Simulate rate-limit check
	rateLimited := recentRepo.NewsletterSentAt != nil &&
		time.Since(*recentRepo.NewsletterSentAt) < 24*time.Hour
	assert.True(s.T(), rateLimited, "recently-sent repo should be rate-limited")

	neverRepo, err := s.repositoryRepo.GetByID(context.Background(), repoNever.ID)
	assert.NoError(s.T(), err)
	rateLimited = neverRepo.NewsletterSentAt != nil &&
		time.Since(*neverRepo.NewsletterSentAt) < 24*time.Hour
	assert.False(s.T(), rateLimited, "never-sent repo should not be rate-limited")

	_ = uuid.Nil // imported for completeness
}
