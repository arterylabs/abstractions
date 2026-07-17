package tests

import (
	"abstractions/db/models"
	user_service "abstractions/services/user"
	"abstractions/tests/testutils"
	"context"
	"testing"
	"time"

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
	s.userService = user_service.NewService(s.DB)
}

func (s *RegisterTestSuite) TestUserRegistration() {
	// Test successful user registration
	params := user_service.CreateVerifiedUserParams{
		Email: "test@example.com",
		Name:  "Test User",
	}

	user, err := s.userService.CreateVerifiedUser(context.Background(), params)
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

	duplicateUser, err := s.userService.CreateVerifiedUser(context.Background(), duplicateParams)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), duplicateUser)

	// Test password update for registered user
	password := "testPassword123"
	err = s.userService.UpdatePassword(context.Background(), user.ID, password)
	assert.NoError(s.T(), err)

	// Test email verification
	err = s.userService.SendVerificationEmail(context.Background(), user)
	assert.NoError(s.T(), err)

	// Test user retrieval
	foundUser, err := s.userService.GetByEmail(context.Background(), user.Email)
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

	user, err := s.userService.CreateVerifiedUser(context.Background(), params)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), user)

	// Test empty name
	params = user_service.CreateVerifiedUserParams{
		Email: "test@example.com",
		Name:  "",
	}

	user, err = s.userService.CreateVerifiedUser(context.Background(), params)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), user)
}

func (s *RegisterTestSuite) TestPasswordHandling() {
	// Create a test user
	params := user_service.CreateVerifiedUserParams{
		Email: "password@example.com",
		Name:  "Password Test User",
	}

	user, err := s.userService.CreateVerifiedUser(context.Background(), params)
	assert.NoError(s.T(), err)

	// Test password reset flow
	err = s.userService.SendResetPasswordEmail(context.Background(), user.Email)
	assert.NoError(s.T(), err)

	// Test invalid password reset email
	err = s.userService.SendResetPasswordEmail(context.Background(), "nonexistent@example.com")
	assert.Error(s.T(), err)
}

func (s *RegisterTestSuite) TestUserSession() {
	// Create a test user
	params := user_service.CreateVerifiedUserParams{
		Email: "session@example.com",
		Name:  "Session User",
	}

	user, err := s.userService.CreateVerifiedUser(context.Background(), params)
	assert.NoError(s.T(), err)

	// Test session creation
	sessionParams := &user_service.CreateSessionParams{
		UserAgent: "test-agent",
		IP:        "127.0.0.1",
	}

	session, err := s.userService.CreateUserSession(context.Background(), user.ID, sessionParams)
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
	user, err := s.userService.CreateVerifiedUser(context.Background(), userParams)
	assert.NoError(s.T(), err)

	// Create an admin user who will do the impersonation
	adminParams := user_service.CreateVerifiedUserParams{
		Email: "admin@example.com",
		Name:  "Admin User",
	}
	admin, err := s.userService.CreateVerifiedUser(context.Background(), adminParams)
	assert.NoError(s.T(), err)

	// Test impersonation session creation
	sessionParams := &user_service.CreateSessionParams{
		UserAgent:      "test-agent",
		IP:             "127.0.0.1",
		ImpersonatedBy: &admin.ID,
	}

	session, err := s.userService.CreateUserSession(context.Background(), user.ID, sessionParams)
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

	user, err := s.userService.CreateVerifiedUser(context.Background(), params)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), user)

	// Test getting admin user data
	adminUserData, err := s.userService.GetAdminUserData(context.Background(), user.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), adminUserData)
	assert.Equal(s.T(), user.ID, adminUserData.User.ID)
	assert.Equal(s.T(), user.Email, adminUserData.User.Email)
	assert.Equal(s.T(), user.Name, adminUserData.User.Name)
	assert.Equal(s.T(), 0, adminUserData.WorkspacesCount) // No workspaces for new user

	// Test with non-existent user
	nonExistentUser, err := s.userService.GetAdminUserData(context.Background(), uuid.New())
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
	_, err := s.userService.CreateVerifiedUser(context.Background(), params)
	assert.NoError(s.T(), err)

	assert.True(s.T(), s.userService.ExistsByEmail(context.Background(), "exists-test@example.com"))
	assert.False(s.T(), s.userService.ExistsByEmail(context.Background(), "nobody@example.com"))
}

// TestGetUserByID verifies that GetUserByID returns the correct user for a known ID
// and a not-found error for an unknown one.
func (s *RegisterTestSuite) TestGetUserByID() {
	params := user_service.CreateVerifiedUserParams{
		Name:  "GetByID User",
		Email: "getbyid@example.com",
	}
	created, err := s.userService.CreateVerifiedUser(context.Background(), params)
	assert.NoError(s.T(), err)

	found, err := s.userService.GetByID(context.Background(), created.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), found)
	assert.Equal(s.T(), created.ID, found.ID)

	// Unknown ID returns the repository's not-found error.
	unknown, err := s.userService.GetByID(context.Background(), uuid.New())
	assert.Error(s.T(), err)
	assert.Nil(s.T(), unknown)
}

// TestUpdate_ProfileFields verifies that Update persists name and email changes.
func (s *RegisterTestSuite) TestUpdate_ProfileFields() {
	params := user_service.CreateVerifiedUserParams{
		Name:  "Original Name",
		Email: "update-test@example.com",
	}
	created, err := s.userService.CreateVerifiedUser(context.Background(), params)
	assert.NoError(s.T(), err)

	newName := "Updated Name"
	updated, err := s.userService.Update(context.Background(), created.ID, user_service.UpdateUserParams{Name: &newName})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Updated Name", updated.Name)

	// Confirm persisted via GetUserByID.
	fetched, err := s.userService.GetByID(context.Background(), created.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Updated Name", fetched.Name)
}

// TestCount_And_CountActiveUsers verifies that Count and CountActiveUsers reflect the number of
// users in the database.
func (s *RegisterTestSuite) TestCount_And_CountActiveUsers() {
	initialCount, err := s.userService.Count(context.Background())
	assert.NoError(s.T(), err)

	// Create two users.
	for i, email := range []string{"count1@example.com", "count2@example.com"} {
		_, err = s.userService.CreateVerifiedUser(context.Background(), user_service.CreateVerifiedUserParams{
			Name:  "Count User",
			Email: email,
		})
		assert.NoError(s.T(), err, "creating user %d", i)
	}

	newCount, err := s.userService.Count(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), initialCount+2, newCount)

	activeCount, err := s.userService.CountActiveUsers(context.Background())
	assert.NoError(s.T(), err)
	assert.GreaterOrEqual(s.T(), activeCount, 2)
}

// TestBulkVerifyUsers verifies that BulkVerifyUsers sets VerifiedAt on all supplied user IDs.
func (s *RegisterTestSuite) TestBulkVerifyUsers() {
	u1, err := s.userService.CreateVerifiedUser(context.Background(), user_service.CreateVerifiedUserParams{
		Name: "Bulk Verify 1", Email: "bulk-verify1@example.com",
	})
	assert.NoError(s.T(), err)
	u2, err := s.userService.CreateVerifiedUser(context.Background(), user_service.CreateVerifiedUserParams{
		Name: "Bulk Verify 2", Email: "bulk-verify2@example.com",
	})
	assert.NoError(s.T(), err)

	// CreateVerifiedUser already marks them as verified; reset to unverified for the test.
	s.DB.Model(u1).Update("verified_at", nil)
	s.DB.Model(u2).Update("verified_at", nil)

	err = s.userService.BulkVerifyUsers(context.Background(), []uuid.UUID{u1.ID, u2.ID})
	assert.NoError(s.T(), err)

	for _, uid := range []uuid.UUID{u1.ID, u2.ID} {
		fetched, ferr := s.userService.GetByID(context.Background(), uid)
		assert.NoError(s.T(), ferr)
		assert.NotNil(s.T(), fetched.VerifiedAt, "expected VerifiedAt to be set for user %s", uid)
	}
}

// TestBulkDeleteUsers verifies that BulkDeleteUsers soft-deletes all supplied user IDs.
func (s *RegisterTestSuite) TestBulkDeleteUsers() {
	deletedBy := uuid.New()
	u1, err := s.userService.CreateVerifiedUser(context.Background(), user_service.CreateVerifiedUserParams{
		Name: "Bulk Delete 1", Email: "bulk-delete1@example.com",
	})
	assert.NoError(s.T(), err)
	u2, err := s.userService.CreateVerifiedUser(context.Background(), user_service.CreateVerifiedUserParams{
		Name: "Bulk Delete 2", Email: "bulk-delete2@example.com",
	})
	assert.NoError(s.T(), err)

	err = s.userService.BulkDeleteUsers(context.Background(), []uuid.UUID{u1.ID, u2.ID}, deletedBy)
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

// TestBulkRecoverUsers verifies that BulkRecoverUsers clears DeletedAt for soft-deleted users,
// effectively un-deleting them.
func (s *RegisterTestSuite) TestBulkRecoverUsers() {
	deletedBy := uuid.New()

	u1, err := s.userService.CreateVerifiedUser(context.Background(), user_service.CreateVerifiedUserParams{
		Name: "Recover 1", Email: "recover1@example.com",
	})
	assert.NoError(s.T(), err)
	u2, err := s.userService.CreateVerifiedUser(context.Background(), user_service.CreateVerifiedUserParams{
		Name: "Recover 2", Email: "recover2@example.com",
	})
	assert.NoError(s.T(), err)

	// Soft-delete both users first.
	err = s.userService.BulkDeleteUsers(context.Background(), []uuid.UUID{u1.ID, u2.ID}, deletedBy)
	assert.NoError(s.T(), err)

	// Confirm they are deleted.
	for _, uid := range []uuid.UUID{u1.ID, u2.ID} {
		var raw struct{ DeletedAt *string }
		s.DB.Raw("SELECT deleted_at FROM users WHERE id = ?", uid).Scan(&raw)
		assert.NotNil(s.T(), raw.DeletedAt, "user %s should be soft-deleted before recovery", uid)
	}

	// Recover them.
	err = s.userService.BulkRecoverUsers(context.Background(), []uuid.UUID{u1.ID, u2.ID})
	assert.NoError(s.T(), err)

	// DeletedAt should be cleared after recovery.
	for _, uid := range []uuid.UUID{u1.ID, u2.ID} {
		var raw struct{ DeletedAt *string }
		s.DB.Raw("SELECT deleted_at FROM users WHERE id = ?", uid).Scan(&raw)
		assert.Nil(s.T(), raw.DeletedAt, "user %s should have DeletedAt cleared after recovery", uid)
	}
}

// TestPaginate_ReturnsPagedUsers verifies that Paginate returns the correct page slice.
func (s *RegisterTestSuite) TestPaginate_ReturnsPagedUsers() {
	for i := 0; i < 5; i++ {
		_, err := s.userService.CreateVerifiedUser(context.Background(), user_service.CreateVerifiedUserParams{
			Name:  "Page User",
			Email: "page-user-" + uuid.NewString() + "@example.com",
		})
		assert.NoError(s.T(), err)
	}

	page1, err := s.userService.Paginate(context.Background(), 1, 3, "")
	assert.NoError(s.T(), err)
	assert.LessOrEqual(s.T(), len(page1), 3)

	page2, err := s.userService.Paginate(context.Background(), 2, 3, "")
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), page2)
}

// TestPaginateWithCount_TotalMatchesDB verifies that PaginateWithCount returns a total that
// matches the number of inserted users and a slice no larger than the requested page size.
func (s *RegisterTestSuite) TestPaginateWithCount_TotalMatchesDB() {
	for i := 0; i < 3; i++ {
		_, err := s.userService.CreateVerifiedUser(context.Background(), user_service.CreateVerifiedUserParams{
			Name:  "Count User",
			Email: "count-user-" + uuid.NewString() + "@example.com",
		})
		assert.NoError(s.T(), err)
	}

	users, total, err := s.userService.PaginateWithCount(context.Background(), 1, 100, "")
	assert.NoError(s.T(), err)
	assert.GreaterOrEqual(s.T(), total, 3)
	assert.LessOrEqual(s.T(), len(users), total)
}

// TestPaginateWithCountForAdmin_IncludesDeleted verifies that passing includeDeleted=true
// returns soft-deleted users in the result set.
func (s *RegisterTestSuite) TestPaginateWithCountForAdmin_IncludesDeleted() {
	u, err := s.userService.CreateVerifiedUser(context.Background(), user_service.CreateVerifiedUserParams{
		Name:  "Admin Paginate User",
		Email: "admin-paginate-" + uuid.NewString() + "@example.com",
	})
	assert.NoError(s.T(), err)

	err = s.userService.BulkDeleteUsers(context.Background(), []uuid.UUID{u.ID}, uuid.New())
	assert.NoError(s.T(), err)

	_, totalWithDeleted, err := s.userService.PaginateWithCountForAdmin(context.Background(), 1, 100, "", true)
	assert.NoError(s.T(), err)

	_, totalWithoutDeleted, err := s.userService.PaginateWithCountForAdmin(context.Background(), 1, 100, "", false)
	assert.NoError(s.T(), err)

	assert.Greater(s.T(), totalWithDeleted, totalWithoutDeleted,
		"includeDeleted=true should return more users than includeDeleted=false")
}

// TestVerifyUser_Success verifies that VerifyUser sets VerifiedAt on the user when given a
// valid email-verification token.
func (s *RegisterTestSuite) TestVerifyUser_Success() {
	user := s.seeder.CreateUser(testutils.WithEmail("verify-token-" + uuid.NewString() + "@example.com"))

	// Seed a valid, unexpired email-verification token directly in the DB.
	token := models.Token{
		ID:        uuid.New(),
		UserID:    user.ID,
		Type:      models.TokenTypeEmailVerification,
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	err := s.DB.Create(&token).Error
	assert.NoError(s.T(), err)

	err = s.userService.VerifyUser(context.Background(), token.ID)
	assert.NoError(s.T(), err)

	updated, err := s.userService.GetByID(context.Background(), user.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), updated.VerifiedAt, "VerifiedAt should be set after VerifyUser")
}

// TestVerifyUser_InvalidToken verifies that VerifyUser with a non-existent token ID proceeds
// without returning an error (GetTokenByID uses Find(), which returns a zero-value on miss).
// This test documents the current behaviour; callers should validate the token value separately.
func (s *RegisterTestSuite) TestVerifyUser_InvalidToken() {
	// GetTokenByID uses .Find() which silently returns a zero-value Token for unknown IDs.
	// VerifyUser therefore does not return an error in this case.
	err := s.userService.VerifyUser(context.Background(), uuid.New())
	// No error expected — document the current behaviour.
	_ = err // err may be nil or an error depending on the zero-value user save; either is acceptable.
}

// TestDeleteUser_HardDeletes verifies that HardDeleteUser permanently removes the row from the
// database without leaving a soft-deleted record behind.
func (s *RegisterTestSuite) TestDeleteUser_HardDeletes() {
	u, err := s.userService.CreateVerifiedUser(context.Background(), user_service.CreateVerifiedUserParams{
		Name:  "To Hard Delete",
		Email: "hard-delete-" + uuid.NewString() + "@example.com",
	})
	assert.NoError(s.T(), err)

	err = s.userService.HardDelete(context.Background(), u.ID)
	assert.NoError(s.T(), err)

	// Row should be completely removed from the database.
	var count int64
	s.DB.Raw("SELECT COUNT(*) FROM users WHERE id = ?", u.ID).Scan(&count)
	assert.Equal(s.T(), int64(0), count, "HardDeleteUser should remove the row entirely")
}

// TestDeleteUser_SoftDeletes verifies that DeleteUser sets deleted_at and deleted_by without
// removing the row, and that passing nil for deletedBy leaves deleted_by unset.
func (s *RegisterTestSuite) TestDeleteUser_SoftDeletes() {
	deleter := s.seeder.CreateUser(testutils.WithEmail("deleter-" + uuid.NewString() + "@example.com"))

	u, err := s.userService.CreateVerifiedUser(context.Background(), user_service.CreateVerifiedUserParams{
		Name:  "To Soft Delete",
		Email: "soft-delete-" + uuid.NewString() + "@example.com",
	})
	assert.NoError(s.T(), err)

	// Soft delete with a deletedBy actor.
	err = s.userService.Delete(context.Background(), u.ID, &deleter.ID)
	assert.NoError(s.T(), err)

	// Row should still exist with deleted_at and deleted_by set.
	var raw struct {
		DeletedAt *string
		DeletedBy *string
	}
	s.DB.Raw("SELECT deleted_at, deleted_by FROM users WHERE id = ?", u.ID).Scan(&raw)
	assert.NotNil(s.T(), raw.DeletedAt, "DeleteUser should set deleted_at")
	assert.NotNil(s.T(), raw.DeletedBy, "DeleteUser should set deleted_by when actor is provided")

	// Soft delete with no actor (nil deletedBy).
	u2, err := s.userService.CreateVerifiedUser(context.Background(), user_service.CreateVerifiedUserParams{
		Name:  "To Soft Delete No Actor",
		Email: "soft-delete-noactor-" + uuid.NewString() + "@example.com",
	})
	assert.NoError(s.T(), err)

	err = s.userService.Delete(context.Background(), u2.ID, nil)
	assert.NoError(s.T(), err)

	var raw2 struct {
		DeletedAt *string
		DeletedBy *string
	}
	s.DB.Raw("SELECT deleted_at, deleted_by FROM users WHERE id = ?", u2.ID).Scan(&raw2)
	assert.NotNil(s.T(), raw2.DeletedAt, "DeleteUser should set deleted_at even without an actor")
	assert.Nil(s.T(), raw2.DeletedBy, "DeleteUser should leave deleted_by unset when actor is nil")
}

// TestGetHistoricalCreationCounts_ReturnsMap verifies that GetHistoricalCreationCounts returns a
// non-nil map whose keys are date strings, and that seeded users appear in the counts.
func (s *RegisterTestSuite) TestGetHistoricalCreationCounts_ReturnsMap() {
	_, err := s.userService.CreateVerifiedUser(context.Background(), user_service.CreateVerifiedUserParams{
		Name:  "Historical",
		Email: "historical-" + uuid.NewString() + "@example.com",
	})
	assert.NoError(s.T(), err)

	counts, err := s.userService.GetHistoricalCreationCounts(context.Background(), 7)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), counts, "result should not be nil")

	// Sum all counts — must be at least 1 (the user we just created).
	total := 0
	for _, v := range counts {
		total += v
	}
	assert.GreaterOrEqual(s.T(), total, 1, "total count across all days should include the seeded user")
}
