package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"abstractions/db/models"
	"abstractions/server"
	"abstractions/server/routes"
	"abstractions/tests/testutils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// SecurityEmailVerificationTestSuite verifies that email-verification tokens are
// single-use and honour expiry, exercised through the real POST /users/verify route.
type SecurityEmailVerificationTestSuite struct {
	testutils.TestSuite
	seeder *testutils.Seeder
}

func TestSecurityEmailVerificationSuite(t *testing.T) {
	suite.Run(t, new(SecurityEmailVerificationTestSuite))
}

func (s *SecurityEmailVerificationTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
}

func (s *SecurityEmailVerificationTestSuite) router() *echo.Echo {
	s.T().Helper()
	srv := &server.Server{
		Echo:   echo.New(),
		DB:     s.DB,
		Config: s.Config,
	}
	routes.ConfigureRoutes(srv)
	return srv.Echo
}

func (s *SecurityEmailVerificationTestSuite) postVerify(e *echo.Echo, token string) *httptest.ResponseRecorder {
	s.T().Helper()
	body, _ := json.Marshal(map[string]string{"token": token})
	req := httptest.NewRequest(http.MethodPost, "/users/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

// seedVerificationToken creates an unverified user plus an email-verification token.
func (s *SecurityEmailVerificationTestSuite) seedVerificationToken(expiresAt time.Time, usedAt *time.Time) (*models.User, *models.Token) {
	s.T().Helper()
	user := s.seeder.CreateUser() // unverified (no VerifiedAt)

	token := &models.Token{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     "verify-" + uuid.NewString(),
		Type:      models.TokenTypeEmailVerification,
		ExpiresAt: expiresAt,
		UsedAt:    usedAt,
	}
	s.Require().NoError(s.DB.Create(token).Error)
	return user, token
}

func (s *SecurityEmailVerificationTestSuite) isVerified(userID uuid.UUID) bool {
	s.T().Helper()
	var user models.User
	s.Require().NoError(s.DB.First(&user, "id = ?", userID).Error)
	return user.VerifiedAt != nil
}

// TestExpiredVerificationTokenRejected: an expired token must not verify the user.
func (s *SecurityEmailVerificationTestSuite) TestExpiredVerificationTokenRejected() {
	// Use a comfortably-past expiry; the DB column is a naive timestamp, so a
	// small offset can round-trip ambiguously across timezones (see token_service_test).
	user, token := s.seedVerificationToken(time.Now().Add(-24*time.Hour), nil)

	rec := s.postVerify(s.router(), token.Token)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code,
		"expired verification token must be rejected, got %d: %s", rec.Code, rec.Body.String())
	assert.False(s.T(), s.isVerified(user.ID), "user must remain unverified after an expired token is used")
}

// TestUsedVerificationTokenRejected: a token already marked used must not verify the user.
func (s *SecurityEmailVerificationTestSuite) TestUsedVerificationTokenRejected() {
	usedAt := time.Now().Add(-1 * time.Minute)
	user, token := s.seedVerificationToken(time.Now().Add(24*time.Hour), &usedAt)

	rec := s.postVerify(s.router(), token.Token)

	assert.Equal(s.T(), http.StatusBadRequest, rec.Code,
		"already-used verification token must be rejected, got %d: %s", rec.Code, rec.Body.String())
	assert.False(s.T(), s.isVerified(user.ID), "user must remain unverified when replaying a used token")
}

// TestVerificationTokenIsSingleUse: a valid token verifies once, then cannot be replayed.
func (s *SecurityEmailVerificationTestSuite) TestVerificationTokenIsSingleUse() {
	_, token := s.seedVerificationToken(time.Now().Add(24*time.Hour), nil)
	e := s.router()

	first := s.postVerify(e, token.Token)
	assert.Equal(s.T(), http.StatusOK, first.Code,
		"a valid, unexpired token should verify successfully the first time, got %d: %s", first.Code, first.Body.String())

	second := s.postVerify(e, token.Token)
	assert.Equal(s.T(), http.StatusBadRequest, second.Code,
		"replaying a consumed verification token must be rejected, got %d: %s", second.Code, second.Body.String())
}
