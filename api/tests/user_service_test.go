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

// TestExistsByEmail verifies that ExistsByEmail returns true for a known email and false for an
// unknown one.
func (s *RegisterTestSuite) TestExistsByEmail() {
	params := user_service.CreateVerifiedUserParams{
		Name:  "Exists Test",
		Email: "exists-test@example.com",
	}
	_, err := s.userService.CreateVerifiedUser(params)
	assert.NoError(s.T(), err)

	assert.True(s.T(), s.userService.ExistsByEmail("exists-test@example.com"))
	assert.False(s.T(), s.userService.ExistsByEmail("nobody@example.com"))
}

// TestGetUserByID verifies that GetUserByID returns the correct user for a known ID and an empty
// struct for an unknown one (repository uses Find which doesn't return error on not-found).
func (s *RegisterTestSuite) TestGetUserByID() {
	params := user_service.CreateVerifiedUserParams{
		Name:  "GetByID User",
		Email: "getbyid@example.com",
	}
	created, err := s.userService.CreateVerifiedUser(params)
	assert.NoError(s.T(), err)

	found, err := s.userService.GetUserByID(created.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), found)
	assert.Equal(s.T(), created.ID, found.ID)

	// Unknown ID returns an empty User (repository uses Find, not First — no ErrRecordNotFound).
	unknown, err := s.userService.GetUserByID(uuid.New())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), uuid.Nil, unknown.ID)
}

// TestUpdate_ProfileFields verifies that Update persists name and email changes.
func (s *RegisterTestSuite) TestUpdate_ProfileFields() {
	params := user_service.CreateVerifiedUserParams{
		Name:  "Original Name",
		Email: "update-test@example.com",
	}
	created, err := s.userService.CreateVerifiedUser(params)
	assert.NoError(s.T(), err)

	newName := "Updated Name"
	updated, err := s.userService.Update(created.ID, user_service.UpdateUserParams{Name: &newName})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Updated Name", updated.Name)

	// Confirm persisted via GetUserByID.
	fetched, err := s.userService.GetUserByID(created.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Updated Name", fetched.Name)
}

// TestCount_And_CountActiveUsers verifies that Count and CountActiveUsers reflect the number of
// users in the database.
func (s *RegisterTestSuite) TestCount_And_CountActiveUsers() {
	initialCount, err := s.userService.Count()
	assert.NoError(s.T(), err)

	// Create two users.
	for i, email := range []string{"count1@example.com", "count2@example.com"} {
		_, err = s.userService.CreateVerifiedUser(user_service.CreateVerifiedUserParams{
			Name:  "Count User",
			Email: email,
		})
		assert.NoError(s.T(), err, "creating user %d", i)
	}

	newCount, err := s.userService.Count()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), initialCount+2, newCount)

	activeCount, err := s.userService.CountActiveUsers()
	assert.NoError(s.T(), err)
	assert.GreaterOrEqual(s.T(), activeCount, 2)
}

// TestBulkVerifyUsers verifies that BulkVerifyUsers sets VerifiedAt on all supplied user IDs.
func (s *RegisterTestSuite) TestBulkVerifyUsers() {
	u1, err := s.userService.CreateVerifiedUser(user_service.CreateVerifiedUserParams{
		Name: "Bulk Verify 1", Email: "bulk-verify1@example.com",
	})
	assert.NoError(s.T(), err)
	u2, err := s.userService.CreateVerifiedUser(user_service.CreateVerifiedUserParams{
		Name: "Bulk Verify 2", Email: "bulk-verify2@example.com",
	})
	assert.NoError(s.T(), err)

	// CreateVerifiedUser already marks them as verified; reset to unverified for the test.
	s.DB.Model(u1).Update("verified_at", nil)
	s.DB.Model(u2).Update("verified_at", nil)

	err = s.userService.BulkVerifyUsers([]uuid.UUID{u1.ID, u2.ID})
	assert.NoError(s.T(), err)

	for _, uid := range []uuid.UUID{u1.ID, u2.ID} {
		fetched, ferr := s.userService.GetUserByID(uid)
		assert.NoError(s.T(), ferr)
		assert.NotNil(s.T(), fetched.VerifiedAt, "expected VerifiedAt to be set for user %s", uid)
	}
}

// TestBulkDeleteUsers verifies that BulkDeleteUsers soft-deletes all supplied user IDs.
func (s *RegisterTestSuite) TestBulkDeleteUsers() {
	deletedBy := uuid.New()
	u1, err := s.userService.CreateVerifiedUser(user_service.CreateVerifiedUserParams{
		Name: "Bulk Delete 1", Email: "bulk-delete1@example.com",
	})
	assert.NoError(s.T(), err)
	u2, err := s.userService.CreateVerifiedUser(user_service.CreateVerifiedUserParams{
		Name: "Bulk Delete 2", Email: "bulk-delete2@example.com",
	})
	assert.NoError(s.T(), err)

	err = s.userService.BulkDeleteUsers([]uuid.UUID{u1.ID, u2.ID}, deletedBy)
	assert.NoError(s.T(), err)

	// Soft-deleted users should have DeletedAt set.
	for _, uid := range []uuid.UUID{u1.ID, u2.ID} {
		var raw struct {
			DeletedAt *string
		}
		s.DB.Raw("SELECT deleted_at FROM users WHERE id = ?", uid).Scan(&raw)
		assert.NotNil(s.T(), raw.DeletedAt, "expected DeletedAt to be set for user %s", uid)
	}
}
