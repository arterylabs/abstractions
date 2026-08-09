package tests

import (
	"abstractions/clients/email"
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/tests/testutils"
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type EmailRepositoryTestSuite struct {
	testutils.TestSuite
	seeder          *testutils.Seeder
	logRepo         *repositories.EmailLogRepository
	suppressionRepo *repositories.EmailSuppressionRepository
}

func TestEmailRepository(t *testing.T) {
	suite.Run(t, new(EmailRepositoryTestSuite))
}

func (s *EmailRepositoryTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.logRepo = repositories.NewEmailLogRepository(s.DB)
	s.suppressionRepo = repositories.NewEmailSuppressionRepository(s.DB)
}

// TestCreateAndCountSoftFailures verifies only soft failures are counted toward the threshold.
func (s *EmailRepositoryTestSuite) TestCreateAndCountSoftFailures() {
	addr := "count@example.com"
	s.seeder.CreateEmailLog(addr, models.EmailLogStatusSent)
	s.seeder.CreateEmailLog(addr, models.EmailLogStatusFailed, func(l *models.EmailLog) {
		t := models.EmailErrorTypeSoft
		l.ErrorType = &t
	})
	s.seeder.CreateEmailLog(addr, models.EmailLogStatusFailed, func(l *models.EmailLog) {
		t := models.EmailErrorTypeSoft
		l.ErrorType = &t
	})
	// A hard failure must not be counted as a soft failure.
	s.seeder.CreateEmailLog(addr, models.EmailLogStatusFailed, func(l *models.EmailLog) {
		t := models.EmailErrorTypeHard
		l.ErrorType = &t
	})

	count, err := s.logRepo.CountSoftFailuresByEmail(s.Ctx, addr)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 2, count)
}

// TestCountSoftFailuresByEmail_ReturnsZeroForUnknownAddress ensures an address with no logs
// at all is treated the same as one with none below threshold.
func (s *EmailRepositoryTestSuite) TestCountSoftFailuresByEmail_ReturnsZeroForUnknownAddress() {
	count, err := s.logRepo.CountSoftFailuresByEmail(s.Ctx, "never-logged@example.com")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, count)
}

// TestEmailLogCreate_RequiresEmailAndStatus ensures malformed log entries are rejected
// rather than silently persisted with missing required fields.
func (s *EmailRepositoryTestSuite) TestEmailLogCreate_RequiresEmailAndStatus() {
	err := s.logRepo.Create(s.Ctx, &models.EmailLog{Status: models.EmailLogStatusSent})
	assert.Error(s.T(), err)

	err = s.logRepo.Create(s.Ctx, &models.EmailLog{Email: "missing-status@example.com"})
	assert.Error(s.T(), err)
}

// TestSuppressionUpsertIsIdempotent ensures a second Upsert for an active address is a no-op.
func (s *EmailRepositoryTestSuite) TestSuppressionUpsertIsIdempotent() {
	addr := "dupe@example.com"
	assert.NoError(s.T(), s.suppressionRepo.Upsert(s.Ctx, addr, models.SuppressionReasonHardBounce))
	assert.NoError(s.T(), s.suppressionRepo.Upsert(s.Ctx, addr, models.SuppressionReasonThresholdExceeded))

	active, err := s.suppressionRepo.GetActiveByEmail(s.Ctx, addr)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), active)
	// The first reason is retained; the second Upsert did not create a duplicate.
	assert.Equal(s.T(), models.SuppressionReasonHardBounce, active.Reason)
}

// TestSuppressionDeleteAllowsReSuppression verifies a soft-deleted suppression is inactive
// and the address can be suppressed again.
func (s *EmailRepositoryTestSuite) TestSuppressionDeleteAllowsReSuppression() {
	addr := "unblock@example.com"
	assert.NoError(s.T(), s.suppressionRepo.Upsert(s.Ctx, addr, models.SuppressionReasonManual))
	assert.NoError(s.T(), s.suppressionRepo.Delete(s.Ctx, addr))

	active, err := s.suppressionRepo.GetActiveByEmail(s.Ctx, addr)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), active)

	// Re-suppression works because the unique index only covers non-deleted rows.
	assert.NoError(s.T(), s.suppressionRepo.Upsert(s.Ctx, addr, models.SuppressionReasonComplaint))
	active, err = s.suppressionRepo.GetActiveByEmail(s.Ctx, addr)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), active)
	assert.Equal(s.T(), models.SuppressionReasonComplaint, active.Reason)
}

// TestSendSkipsBlockedDomain verifies a disposable-domain recipient is skipped in the shared
// sendEmail shell (used by SES and Cloudflare) and recorded as suppressed.
func (s *EmailRepositoryTestSuite) TestSendSkipsBlockedDomain() {
	email.SetSuppressionStore(email.NewGormSuppressionStore(s.DB))
	defer email.SetSuppressionStore(nil)

	err := email.SendEmailUsingSES(s.Ctx, email.AWSEmailParams{
		Email:   "temp@mailinator.com",
		Subject: "hello",
	})
	assert.ErrorIs(s.T(), err, email.ErrEmailBlockedDomain)

	var logs []models.EmailLog
	assert.NoError(s.T(), s.DB.Where("email = ?", "temp@mailinator.com").Find(&logs).Error)
	assert.Len(s.T(), logs, 1)
	assert.Equal(s.T(), models.EmailLogStatusSuppressed, logs[0].Status)
	assert.NotNil(s.T(), logs[0].ErrorType)
	assert.Equal(s.T(), models.EmailErrorTypeBlockedDomain, *logs[0].ErrorType)
}

// TestSendSkipsSuppressedAddress verifies an already-suppressed recipient is skipped before
// SES and recorded as suppressed. This path never contacts SES.
func (s *EmailRepositoryTestSuite) TestSendSkipsSuppressedAddress() {
	email.SetSuppressionStore(email.NewGormSuppressionStore(s.DB))
	defer email.SetSuppressionStore(nil)

	addr := "blocked@example.com"
	s.seeder.CreateEmailSuppression(addr, models.SuppressionReasonHardBounce)

	err := email.SendEmailUsingSES(s.Ctx, email.AWSEmailParams{Email: addr, Subject: "hello"})
	assert.ErrorIs(s.T(), err, email.ErrEmailSuppressed)

	var logs []models.EmailLog
	assert.NoError(s.T(), s.DB.Where("email = ?", addr).Find(&logs).Error)
	assert.Len(s.T(), logs, 1)
	assert.Equal(s.T(), models.EmailLogStatusSuppressed, logs[0].Status)
}

// TestSendBulkEmail_SkipsSuppressedButContinuesOthers verifies that a suppressed recipient
// in a batch is skipped and logged, while the rest of the batch still goes through - a
// suppressed address must not abort the whole send.
func (s *EmailRepositoryTestSuite) TestSendBulkEmail_SkipsSuppressedButContinuesOthers() {
	email.SetSuppressionStore(email.NewGormSuppressionStore(s.DB))
	defer email.SetSuppressionStore(nil)

	suppressedAddr := "bulk-blocked@example.com"
	s.seeder.CreateEmailSuppression(suppressedAddr, models.SuppressionReasonHardBounce)

	var mu sync.Mutex
	completed := 0
	tracker := func(ctx context.Context, params email.AWSEmailParams) error {
		err := email.SendEmailUsingSES(ctx, params)
		mu.Lock()
		completed++
		mu.Unlock()
		return err
	}

	params := []email.AWSEmailParams{
		{Email: suppressedAddr, Subject: "hello"},
		// The @test.arterylabs.com suffix short-circuits before SES is contacted.
		{Email: "ok-1@test.arterylabs.com", Subject: "hello"},
		{Email: "ok-2@test.arterylabs.com", Subject: "hello"},
	}

	err := email.SendBulkEmailUsingSES(s.Ctx, params, tracker)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 3, completed)

	var logs []models.EmailLog
	assert.NoError(s.T(), s.DB.Where("email = ?", suppressedAddr).Find(&logs).Error)
	assert.Len(s.T(), logs, 1)
	assert.Equal(s.T(), models.EmailLogStatusSuppressed, logs[0].Status)
}
