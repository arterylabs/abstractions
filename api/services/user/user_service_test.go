package user_service_test

import (
	"abstractions/db/models"
	user_service "abstractions/services/user"
	"abstractions/tests/testutils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserServiceTestSuite struct {
	testutils.TestSuite
	seeder  *testutils.Seeder
	service *user_service.Service
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (s *UserServiceTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.service = user_service.NewUserService(s.DB)
}

func (s *UserServiceTestSuite) TestCreateUser_Success() {
	user := &models.User{
		ID:    uuid.New(),
		Email: "create-test-" + uuid.New().String() + "@example.com",
		Name:  "Create Test",
	}

	err := s.service.CreateUser(user)

	assert.NoError(s.T(), err)
	assert.NotEqual(s.T(), uuid.Nil, user.ID)
}

func (s *UserServiceTestSuite) TestCreateUser_DuplicateEmail() {
	email := "dup-" + uuid.New().String() + "@example.com"

	first := &models.User{ID: uuid.New(), Email: email, Name: "First"}
	err := s.service.CreateUser(first)
	assert.NoError(s.T(), err)

	second := &models.User{ID: uuid.New(), Email: email, Name: "Second"}
	err = s.service.CreateUser(second)

	assert.Error(s.T(), err, "creating a user with a duplicate email should fail")
}

func (s *UserServiceTestSuite) TestGetUserByEmail_ReturnsCorrectUser() {
	user := s.seeder.CreateUser(testutils.WithEmail("find-by-email-" + uuid.New().String() + "@example.com"))

	found, err := s.service.GetByEmail(user.Email)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), found)
	assert.Equal(s.T(), user.ID, found.ID)
	assert.Equal(s.T(), user.Email, found.Email)
}

func (s *UserServiceTestSuite) TestGetUserByEmail_NotFound() {
	_, err := s.service.GetByEmail("nonexistent-" + uuid.New().String() + "@example.com")

	assert.Error(s.T(), err)
}

func (s *UserServiceTestSuite) TestExistsByEmail_ReturnsTrueForExistingUser() {
	user := s.seeder.CreateUser()

	exists := s.service.ExistsByEmail(user.Email)

	assert.True(s.T(), exists)
}

func (s *UserServiceTestSuite) TestExistsByEmail_ReturnsFalseForMissingUser() {
	exists := s.service.ExistsByEmail("ghost-" + uuid.New().String() + "@example.com")

	assert.False(s.T(), exists)
}

func (s *UserServiceTestSuite) TestGetUserByID_ReturnsCorrectUser() {
	user := s.seeder.CreateUser()

	found, err := s.service.GetUserByID(user.ID)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), found)
	assert.Equal(s.T(), user.ID, found.ID)
}

func (s *UserServiceTestSuite) TestUpdate_ChangesFieldsOnDisk() {
	user := s.seeder.CreateUser()
	newName := "Updated Name"

	updated, err := s.service.Update(user.ID, user_service.UpdateUserParams{Name: &newName})

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), newName, updated.Name)

	// Verify persisted
	fetched, err := s.service.GetUserByID(user.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), newName, fetched.Name)
}
