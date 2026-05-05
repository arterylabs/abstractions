package tests

import (
	"abstractions/db/models"
	session_service "abstractions/services/session"
	"abstractions/tests/testutils"
	"abstractions/utils/jwt_utils"
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
	s.sessionService = session_service.NewSessionService(s.DB)
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

	err := s.sessionService.CreateSession(session)
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
	s.sessionService.CreateSession(session)

	// Test valid session with IsSessionValid
	valid, err := s.sessionService.IsSessionValid(session.ID)
	assert.NoError(s.T(), err)
	assert.True(s.T(), valid)

	// Also test token validation for completeness
	result, err := s.sessionService.ValidateToken(token)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), user.ID, result.UserID)

	// Test invalid token
	result, err = s.sessionService.ValidateToken("invalid-token")
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *SessionTestSuite) TestSessionInvalidation() {
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))

	// Create multiple sessions
	session1 := s.createTestSession(user.ID)
	session2 := s.createTestSession(user.ID)

	// Verify sessions are initially valid
	valid, err := s.sessionService.IsSessionValid(session1.ID)
	assert.NoError(s.T(), err)
	assert.True(s.T(), valid)

	valid, err = s.sessionService.IsSessionValid(session2.ID)
	assert.NoError(s.T(), err)
	assert.True(s.T(), valid)

	// Test invalidating single session
	err = s.sessionService.InvalidateSession(session1.ID)
	assert.NoError(s.T(), err)

	// Verify first session is invalidated
	valid, err = s.sessionService.IsSessionValid(session1.ID)
	assert.NoError(s.T(), err)
	assert.False(s.T(), valid)

	// Second session should still be valid
	valid, err = s.sessionService.IsSessionValid(session2.ID)
	assert.NoError(s.T(), err)
	assert.True(s.T(), valid)

	// Test invalidating all sessions
	err = s.sessionService.InvalidateAllSessionsForUser(user.ID)
	assert.NoError(s.T(), err)

	// Verify all sessions are invalidated
	valid, err = s.sessionService.IsSessionValid(session1.ID)
	assert.NoError(s.T(), err)
	assert.False(s.T(), valid)

	valid, err = s.sessionService.IsSessionValid(session2.ID)
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
	s.sessionService.CreateSession(session)

	// Test expired session validation
	valid, err := s.sessionService.IsSessionValid(session.ID)
	assert.NoError(s.T(), err)
	assert.False(s.T(), valid)

	// Create an active but invalidated session
	sessionID2 := uuid.New()
	expiresAt2 := time.Now().Add(24 * time.Hour) // Not expired
	token2, _ := jwt_utils.GenerateJWT(user.ID.String(), sessionID2.String(), expiresAt2)
	now := time.Now()

	session2 := &models.Session{
		ID:            sessionID2,
		UserID:        user.ID,
		Token:         token2,
		ExpiresAt:     expiresAt2,
		InvalidatedAt: &now,
	}
	s.sessionService.CreateSession(session2)

	// Test invalidated session validation
	valid, err = s.sessionService.IsSessionValid(session2.ID)
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
	s.sessionService.CreateSession(session)
	return session
}
