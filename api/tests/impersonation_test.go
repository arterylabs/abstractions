package tests

import (
	user_service "abstractions/services/user"
	"abstractions/tests/testutils"
	"abstractions/utils/context_utils"
	"abstractions/utils/impersonation_utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ImpersonationTestSuite struct {
	testutils.TestSuite
	userService *user_service.Service
}

func (s *ImpersonationTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.userService = user_service.NewUserService(s.DB)
}

func TestImpersonationSuite(t *testing.T) {
	suite.Run(t, new(ImpersonationTestSuite))
}

func (s *ImpersonationTestSuite) TestImpersonationUtilities() {
	// Create a test user to be impersonated
	userParams := user_service.CreateVerifiedUserParams{
		Email: "impersonated@example.com",
		Name:  "Impersonated User",
	}
	user, err := s.userService.CreateVerifiedUser(userParams)
	assert.NoError(s.T(), err)

	// Create an admin user who will do the impersonation
	adminParams := user_service.CreateVerifiedUserParams{
		Email: "admin@example.com",
		Name:  "Admin User",
	}
	admin, err := s.userService.CreateVerifiedUser(adminParams)
	assert.NoError(s.T(), err)

	// Create impersonation session
	sessionParams := &user_service.CreateSessionParams{
		UserAgent:      "test-agent",
		IP:             "127.0.0.1",
		ImpersonatedBy: &admin.ID,
	}

	session, err := s.userService.CreateUserSession(user.ID, sessionParams)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), session)
	assert.NotNil(s.T(), session.ImpersonatedBy)
	assert.Equal(s.T(), admin.ID, *session.ImpersonatedBy)

	// Test utility functions with Echo context
	e := echo.New()
	c := e.NewContext(nil, nil)

	// Set user and session in context
	c = context_utils.SetUserInContext(c, *user)
	c = context_utils.SetSessionInContext(c, *session)
	c = context_utils.SetTransactionInContext(c, s.DB)

	// Test IsImpersonatedSession
	assert.True(s.T(), impersonation_utils.IsImpersonatedSession(c))

	// Test GetImpersonatorUser
	impersonator, err := impersonation_utils.GetImpersonatorUser(c)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), impersonator)
	assert.Equal(s.T(), admin.ID, impersonator.ID)
	assert.Equal(s.T(), admin.Email, impersonator.Email)

	// Test GetImpersonationInfo
	info, err := impersonation_utils.GetImpersonationInfo(c)
	assert.NoError(s.T(), err)
	assert.True(s.T(), info.IsImpersonated)
	assert.Equal(s.T(), admin.ID, *info.ImpersonatedBy)
	assert.NotNil(s.T(), info.ImpersonatedAt)
	assert.NotNil(s.T(), info.ImpersonatorUser)
	assert.Equal(s.T(), admin.Email, info.ImpersonatorUser.Email)

	// Test BlockIfImpersonated
	assert.True(s.T(), impersonation_utils.BlockIfImpersonated(c))
}

func (s *ImpersonationTestSuite) TestNonImpersonatedSession() {
	// Create a regular user session (not impersonated)
	userParams := user_service.CreateVerifiedUserParams{
		Email: "regular@example.com",
		Name:  "Regular User",
	}
	user, err := s.userService.CreateVerifiedUser(userParams)
	assert.NoError(s.T(), err)

	// Create regular session without impersonation
	sessionParams := &user_service.CreateSessionParams{
		UserAgent: "test-agent",
		IP:        "127.0.0.1",
	}

	session, err := s.userService.CreateUserSession(user.ID, sessionParams)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), session)
	assert.Nil(s.T(), session.ImpersonatedBy)

	// Test utility functions with Echo context
	e := echo.New()
	c := e.NewContext(nil, nil)

	// Set user and session in context
	c = context_utils.SetUserInContext(c, *user)
	c = context_utils.SetSessionInContext(c, *session)
	c = context_utils.SetTransactionInContext(c, s.DB)

	// Test IsImpersonatedSession
	assert.False(s.T(), impersonation_utils.IsImpersonatedSession(c))

	// Test GetImpersonatorUser
	impersonator, err := impersonation_utils.GetImpersonatorUser(c)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), impersonator)

	// Test GetImpersonationInfo
	info, err := impersonation_utils.GetImpersonationInfo(c)
	assert.NoError(s.T(), err)
	assert.False(s.T(), info.IsImpersonated)
	assert.Nil(s.T(), info.ImpersonatedBy)
	assert.Nil(s.T(), info.ImpersonatedAt)
	assert.Nil(s.T(), info.ImpersonatorUser)

	// Test BlockIfImpersonated
	assert.False(s.T(), impersonation_utils.BlockIfImpersonated(c))
}
