package tests

import (
	"abstractions/db/models"
	session_service "abstractions/services/session"
	"abstractions/tests/testutils"
	"abstractions/utils/jwt_utils"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type SessionTestSuite struct {
	testutils.TestSuite
	seeder         *testutils.Seeder
	sessionService *session_service.Service
}

func TestSessionService(t *testing.T) {
	suite.Run(t, new(SessionTestSuite))
}

func (s *SessionTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.sessionService = session_service.NewService(s.DB)
}

func (s *SessionTestSuite) TestCreateSession() {
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))

	// Test creating a new session
	sessionID := uuid.New()
	expiresAt := time.Now().Add(24 * time.Hour)
	token, _ := jwt_utils.GenerateJWT(user.ID.String(), sessionID.String(), expiresAt)

	session := &models.Session{
		ID:        sessionID,
		UserID:    user.ID,
		Token:     token,
		UserAgent: "test-agent",
		IP:        "127.0.0.1",
		ExpiresAt: expiresAt,
	}

	err := s.sessionService.Create(context.Background(), session)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), session)
}

func (s *SessionTestSuite) TestSessionValidation() {
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))

	// Create an active session
	sessionID := uuid.New()
	expiresAt := time.Now().Add(24 * time.Hour)
	token, _ := jwt_utils.GenerateJWT(user.ID.String(), sessionID.String(), expiresAt)

	session := &models.Session{
		ID:        sessionID,
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	s.sessionService.Create(context.Background(), session)

	// Test valid session with IsSessionValid
	valid, err := s.sessionService.IsSessionValid(context.Background(), session.ID)
	assert.NoError(s.T(), err)
	assert.True(s.T(), valid)

	// Also test token validation for completeness
	result, err := s.sessionService.ValidateToken(context.Background(), token)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), user.ID, result.UserID)

	// Test invalid token
	result, err = s.sessionService.ValidateToken(context.Background(), "invalid-token")
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *SessionTestSuite) TestSessionInvalidation() {
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))

	// Create multiple sessions
	session1 := s.createTestSession(user.ID)
	session2 := s.createTestSession(user.ID)

	// Verify sessions are initially valid
	valid, err := s.sessionService.IsSessionValid(context.Background(), session1.ID)
	assert.NoError(s.T(), err)
	assert.True(s.T(), valid)

	valid, err = s.sessionService.IsSessionValid(context.Background(), session2.ID)
	assert.NoError(s.T(), err)
	assert.True(s.T(), valid)

	// Test invalidating single session
	err = s.sessionService.InvalidateSession(context.Background(), session1.ID)
	assert.NoError(s.T(), err)

	// Verify first session is invalidated
	valid, err = s.sessionService.IsSessionValid(context.Background(), session1.ID)
	assert.NoError(s.T(), err)
	assert.False(s.T(), valid)

	// Second session should still be valid
	valid, err = s.sessionService.IsSessionValid(context.Background(), session2.ID)
	assert.NoError(s.T(), err)
	assert.True(s.T(), valid)

	// Test invalidating all sessions
	err = s.sessionService.InvalidateAllSessionsForUser(context.Background(), user.ID)
	assert.NoError(s.T(), err)

	// Verify all sessions are invalidated
	valid, err = s.sessionService.IsSessionValid(context.Background(), session1.ID)
	assert.NoError(s.T(), err)
	assert.False(s.T(), valid)

	valid, err = s.sessionService.IsSessionValid(context.Background(), session2.ID)
	assert.NoError(s.T(), err)
	assert.False(s.T(), valid)
}

func (s *SessionTestSuite) TestExpiredAndInvalidatedSession() {
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))

	// Create expired session
	sessionID := uuid.New()
	expiresAt := time.Now().Add(-24 * time.Hour) // Expired
	token, _ := jwt_utils.GenerateJWT(user.ID.String(), sessionID.String(), expiresAt)

	session := &models.Session{
		ID:        sessionID,
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	s.sessionService.Create(context.Background(), session)

	// Test expired session validation
	valid, err := s.sessionService.IsSessionValid(context.Background(), session.ID)
	assert.NoError(s.T(), err)
	assert.False(s.T(), valid)

	// Create an active but invalidated session
	sessionID2 := uuid.New()
	expiresAt2 := time.Now().Add(24 * time.Hour) // Not expired
	token2, _ := jwt_utils.GenerateJWT(user.ID.String(), sessionID2.String(), expiresAt2)
	now := time.Now()

	session2 := &models.Session{
		ID:        sessionID2,
		UserID:    user.ID,
		Token:     token2,
		ExpiresAt: expiresAt2,
		RevokedAt: &now,
	}
	s.sessionService.Create(context.Background(), session2)

	// Test invalidated session validation
	valid, err = s.sessionService.IsSessionValid(context.Background(), session2.ID)
	assert.NoError(s.T(), err)
	assert.False(s.T(), valid)
}

// Helper to create a test session
func (s *SessionTestSuite) createTestSession(userID uuid.UUID) *models.Session {
	sessionID := uuid.New()
	expiresAt := time.Now().Add(24 * time.Hour)
	token, _ := jwt_utils.GenerateJWT(userID.String(), sessionID.String(), expiresAt)

	session := &models.Session{
		ID:        sessionID,
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	s.sessionService.Create(context.Background(), session)
	return session
}

// TestGetSessionByID verifies that GetSessionByID returns the correct session for a known ID and
// an error for an unknown one.
func (s *SessionTestSuite) TestGetSessionByID() {
	user := s.seeder.CreateUser(testutils.WithEmail("getsession@example.com"))
	session := s.createTestSession(user.ID)

	found, err := s.sessionService.GetByID(context.Background(), session.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), found)
	assert.Equal(s.T(), session.ID, found.ID)

	// Unknown ID returns an error (repository uses First which returns ErrRecordNotFound).
	_, err = s.sessionService.GetByID(context.Background(), uuid.New())
	assert.Error(s.T(), err)
}

// TestGetSessionIfActive verifies that GetSessionIfActive returns an active session and an error
// for an invalidated one.
func (s *SessionTestSuite) TestGetSessionIfActive() {
	user := s.seeder.CreateUser(testutils.WithEmail("activesession@example.com"))
	session := s.createTestSession(user.ID)

	active, err := s.sessionService.GetSessionIfActive(context.Background(), session.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), active)

	// Invalidate and confirm GetSessionIfActive returns an error.
	err = s.sessionService.InvalidateSession(context.Background(), session.ID)
	assert.NoError(s.T(), err)

	_, err = s.sessionService.GetSessionIfActive(context.Background(), session.ID)
	assert.Error(s.T(), err)
}

// TestInvalidateAllImpersonationSessionsForAdmin verifies that only impersonation sessions for a
// given admin are revoked while regular sessions remain untouched.
func (s *SessionTestSuite) TestInvalidateAllImpersonationSessionsForAdmin() {
	admin := s.seeder.CreateUser(testutils.WithEmail("admin-impersonation@example.com"))
	target := s.seeder.CreateUser(testutils.WithEmail("target-user@example.com"))

	// Create a regular (non-impersonation) session for admin.
	regularSession := s.createTestSession(admin.ID)

	// Create an impersonation session: the admin is browsing as target, ImpersonatedBy = admin.ID.
	sessionID := uuid.New()
	expiresAt := time.Now().Add(24 * time.Hour)
	token, _ := jwt_utils.GenerateJWT(target.ID.String(), sessionID.String(), expiresAt)
	impersonationSession := &models.Session{
		ID:             sessionID,
		UserID:         target.ID,
		Token:          token,
		ExpiresAt:      expiresAt,
		ImpersonatedBy: &admin.ID,
	}
	s.sessionService.Create(context.Background(), impersonationSession)

	// Invalidate all impersonation sessions for this admin.
	err := s.sessionService.InvalidateAllImpersonationSessionsForAdmin(context.Background(), admin.ID)
	assert.NoError(s.T(), err)

	// Impersonation session should now be revoked.
	valid, err := s.sessionService.IsSessionValid(context.Background(), impersonationSession.ID)
	assert.NoError(s.T(), err)
	assert.False(s.T(), valid)

	// The admin's own regular session should remain valid.
	valid, err = s.sessionService.IsSessionValid(context.Background(), regularSession.ID)
	assert.NoError(s.T(), err)
	assert.True(s.T(), valid)
}
