package session_service_test

import (
	"abstractions/db/models"
	session_service "abstractions/services/session"
	user_service "abstractions/services/user"
	"abstractions/tests/testutils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SessionServiceTestSuite struct {
	testutils.TestSuite
	seeder  *testutils.Seeder
	service *session_service.Service
	userSvc *user_service.Service
}

func TestSessionService(t *testing.T) {
	suite.Run(t, new(SessionServiceTestSuite))
}

func (s *SessionServiceTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.service = session_service.NewSessionService(s.DB)
	s.userSvc = user_service.NewUserService(s.DB)
}

func (s *SessionServiceTestSuite) TestCreateSession_PersistsAndReturns() {
	user := s.seeder.CreateUser()

	sess, err := s.userSvc.CreateUserSession(user.ID, nil)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), sess)
	assert.NotEmpty(s.T(), sess.Token)
	assert.Equal(s.T(), user.ID, sess.UserID)
}

func (s *SessionServiceTestSuite) TestValidateToken_ValidJWT_ReturnsUser() {
	user := s.seeder.CreateUser()

	sess, err := s.userSvc.CreateUserSession(user.ID, nil)
	s.Require().NoError(err)

	result, err := s.service.ValidateToken(sess.Token)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), user.ID, result.UserID)
	assert.Equal(s.T(), sess.ID, result.SessionID)
}

func (s *SessionServiceTestSuite) TestValidateToken_InvalidJWT_ReturnsError() {
	_, err := s.service.ValidateToken("this.is.not.a.jwt")

	assert.Error(s.T(), err)
}

func (s *SessionServiceTestSuite) TestValidateToken_RevokedSession_ReturnsError() {
	user := s.seeder.CreateUser()

	// Create a session then immediately revoke it
	sess, err := s.userSvc.CreateUserSession(user.ID, nil)
	s.Require().NoError(err)

	err = s.service.InvalidateSession(sess.ID)
	s.Require().NoError(err)

	// GetSessionIfActive filters by revoked_at IS NULL, so a revoked session
	// yields "record not found" which ValidateToken wraps and returns as an error.
	result, err := s.service.ValidateToken(sess.Token)

	assert.Error(s.T(), err, "ValidateToken should error for a revoked session")
	assert.Nil(s.T(), result)
}

func (s *SessionServiceTestSuite) TestInvalidateSession_SetsRevokedAt() {
	user := s.seeder.CreateUser()

	sess, err := s.userSvc.CreateUserSession(user.ID, nil)
	s.Require().NoError(err)

	err = s.service.InvalidateSession(sess.ID)
	assert.NoError(s.T(), err)

	fetched, err := s.service.GetSessionByID(sess.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), fetched.RevokedAt)
}

func (s *SessionServiceTestSuite) TestIsSessionValid_ActiveSession_ReturnsTrue() {
	user := s.seeder.CreateUser()

	sess, err := s.userSvc.CreateUserSession(user.ID, nil)
	s.Require().NoError(err)

	valid, err := s.service.IsSessionValid(sess.ID)

	assert.NoError(s.T(), err)
	assert.True(s.T(), valid)
}

func (s *SessionServiceTestSuite) TestIsSessionValid_RevokedSession_ReturnsFalse() {
	user := s.seeder.CreateUser()

	sess, err := s.userSvc.CreateUserSession(user.ID, nil)
	s.Require().NoError(err)

	s.Require().NoError(s.service.InvalidateSession(sess.ID))

	valid, err := s.service.IsSessionValid(sess.ID)
	assert.NoError(s.T(), err)
	assert.False(s.T(), valid)
}

func (s *SessionServiceTestSuite) TestCreateSession_WithMetadata() {
	user := s.seeder.CreateUser()

	sess, err := s.userSvc.CreateUserSession(user.ID, &user_service.CreateSessionParams{
		UserAgent: "TestAgent/1.0",
		IP:        "127.0.0.1",
	})

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "TestAgent/1.0", sess.UserAgent)
	assert.Equal(s.T(), "127.0.0.1", sess.IP)
}

func (s *SessionServiceTestSuite) TestCreateSession_ExpiryIsInFuture() {
	user := s.seeder.CreateUser()

	sess, err := s.userSvc.CreateUserSession(user.ID, nil)
	s.Require().NoError(err)

	assert.True(s.T(), sess.ExpiresAt.After(time.Now()))
}

func (s *SessionServiceTestSuite) TestGetSessionIfActive_RevokedSession_ReturnsError() {
	user := s.seeder.CreateUser()

	sess, err := s.userSvc.CreateUserSession(user.ID, nil)
	s.Require().NoError(err)

	s.Require().NoError(s.service.InvalidateSession(sess.ID))

	_, err = s.service.GetSessionIfActive(sess.ID)
	assert.Error(s.T(), err, "revoked session should not be returned as active")
}

func (s *SessionServiceTestSuite) TestCreateSession_Direct_PersistsSession() {
	user := s.seeder.CreateUser()

	sess := &models.Session{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     "test-token-" + uuid.New().String(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err := s.service.CreateSession(sess)
	assert.NoError(s.T(), err)

	fetched, err := s.service.GetSessionByID(sess.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), sess.ID, fetched.ID)
}
