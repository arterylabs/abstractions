package tests

import (
	user_service "abstractions/services/user"
	"abstractions/tests/testutils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RegisterTestSuite struct {
	testutils.TestSuite
	seeder      *testutils.Seeder
	userService *user_service.Service
}

func TestRegisterService(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}

func (s *RegisterTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.userService = user_service.NewUserService(s.DB)
}

func (s *RegisterTestSuite) TestUserRegistration() {
	// Test successful user registration
	params := user_service.CreateVerifiedUserParams{
		Email: "test@example.com",
		Name:  "Test User",
	}

	user, err := s.userService.CreateVerifiedUser(params)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), user)
	assert.Equal(s.T(), params.Email, user.Email)
	assert.Equal(s.T(), params.Name, user.Name)
	assert.NotNil(s.T(), user.VerifiedAt)

	// Test duplicate email registration
	duplicateParams := user_service.CreateVerifiedUserParams{
		Email: "test@example.com",
		Name:  "Another User",
	}

	duplicateUser, err := s.userService.CreateVerifiedUser(duplicateParams)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), duplicateUser)

	// Test password update for registered user
	password := "testPassword123"
	err = s.userService.UpdatePassword(user.ID, password)
	assert.NoError(s.T(), err)

	// Test email verification
	err = s.userService.SendVerificationEmail(user)
	assert.NoError(s.T(), err)

	// Test user retrieval
	foundUser, err := s.userService.GetByEmail(user.Email)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user.ID, foundUser.ID)
	assert.Equal(s.T(), user.Email, foundUser.Email)
	assert.Equal(s.T(), user.Name, foundUser.Name)
}

func (s *RegisterTestSuite) TestInvalidRegistration() {
	// Test empty email
	params := user_service.CreateVerifiedUserParams{
		Email: "",
		Name:  "Test User",
	}

	user, err := s.userService.CreateVerifiedUser(params)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), user)

	// Test empty name
	params = user_service.CreateVerifiedUserParams{
		Email: "test@example.com",
		Name:  "",
	}

	user, err = s.userService.CreateVerifiedUser(params)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), user)
}

func (s *RegisterTestSuite) TestPasswordHandling() {
	// Create a test user
	params := user_service.CreateVerifiedUserParams{
		Email: "password@example.com",
		Name:  "Password Test User",
	}

	user, err := s.userService.CreateVerifiedUser(params)
	assert.NoError(s.T(), err)

	// Test password reset flow
	err = s.userService.SendResetPasswordEmail(user.Email)
	assert.NoError(s.T(), err)

	// Test invalid password reset email
	err = s.userService.SendResetPasswordEmail("nonexistent@example.com")
	assert.Error(s.T(), err)
}

func (s *RegisterTestSuite) TestUserSession() {
	// Create a test user
	params := user_service.CreateVerifiedUserParams{
		Email: "session@example.com",
		Name:  "Session User",
	}

	user, err := s.userService.CreateVerifiedUser(params)
	assert.NoError(s.T(), err)

	// Test session creation
	sessionParams := &user_service.CreateSessionParams{
		UserAgent: "test-agent",
		IP:        "127.0.0.1",
	}

	session, err := s.userService.CreateUserSession(user.ID, sessionParams)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), session)
	assert.Equal(s.T(), user.ID, session.UserID)
	assert.NotEmpty(s.T(), session.Token)
	assert.Equal(s.T(), sessionParams.UserAgent, session.UserAgent)
	assert.Equal(s.T(), sessionParams.IP, session.IP)
	assert.Nil(s.T(), session.ImpersonatedAt)
	assert.Nil(s.T(), session.ImpersonatedBy)
}

func (s *RegisterTestSuite) TestUserSessionWithImpersonation() {
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

	// Test impersonation session creation
	sessionParams := &user_service.CreateSessionParams{
		UserAgent:      "test-agent",
		IP:             "127.0.0.1",
		ImpersonatedBy: &admin.ID,
	}

	session, err := s.userService.CreateUserSession(user.ID, sessionParams)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), session)
	assert.Equal(s.T(), user.ID, session.UserID)
	assert.NotEmpty(s.T(), session.Token)
	assert.Equal(s.T(), sessionParams.UserAgent, session.UserAgent)
	assert.Equal(s.T(), sessionParams.IP, session.IP)
	assert.NotNil(s.T(), session.ImpersonatedAt)
	assert.NotNil(s.T(), session.ImpersonatedBy)
	assert.Equal(s.T(), admin.ID, *session.ImpersonatedBy)
}

func (s *RegisterTestSuite) TestGetAdminUserData() {
	// Create a test user first
	params := user_service.CreateVerifiedUserParams{
		Email: "admin-test@example.com",
		Name:  "Admin Test User",
	}

	user, err := s.userService.CreateVerifiedUser(params)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), user)

	// Test getting admin user data
	adminUserData, err := s.userService.GetAdminUserData(user.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), adminUserData)
	assert.Equal(s.T(), user.ID, adminUserData.User.ID)
	assert.Equal(s.T(), user.Email, adminUserData.User.Email)
	assert.Equal(s.T(), user.Name, adminUserData.User.Name)
	assert.Equal(s.T(), 0, adminUserData.WorkspacesCount) // No workspaces for new user

	// Test with non-existent user
	nonExistentUser, err := s.userService.GetAdminUserData(uuid.New())
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), nonExistentUser)
}
