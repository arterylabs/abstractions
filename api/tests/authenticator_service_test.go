package tests

import (
	"abstractions/db/models"
	authenticator_service "abstractions/services/authenticator"
	"abstractions/tests/testutils"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthenticatorTestSuite struct {
	testutils.TestSuite
	seeder               *testutils.Seeder
	authenticatorService *authenticator_service.Service
}

func TestAuthenticatorService(t *testing.T) {
	suite.Run(t, new(AuthenticatorTestSuite))
}

func (s *AuthenticatorTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.authenticatorService = authenticator_service.NewAuthenticatorService(s.DB)
}

func (s *AuthenticatorTestSuite) TestCreateAuthenticator() {
	// Test creating password authenticator
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))
	passwordAuth := &models.Authenticator{
		UserID: user.ID,
		Type:   models.AuthenticatorTypePassword,
		Token:  "hashedpassword",
	}

	err := s.authenticatorService.CreateAuthenticator(context.Background(), passwordAuth)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), passwordAuth)

	// Test creating Google authenticator
	googleAuth := &models.Authenticator{
		UserID:       user.ID,
		Type:         models.AuthenticatorTypeGoogle,
		Token:        "access_token",
		RefreshToken: stringPtr("refresh_token"),
		Scopes:       "email,profile",
	}

	err = s.authenticatorService.CreateAuthenticator(context.Background(), googleAuth)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), googleAuth)
}

func (s *AuthenticatorTestSuite) TestGetAuthenticators() {
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))

	// Create multiple authenticators
	s.seeder.CreateAuthenticator(user, models.AuthenticatorTypePassword, "hashedpass")
	s.seeder.CreateAuthenticator(user, models.AuthenticatorTypeGoogle, "googletoken")

	// Test fetching authenticators
	auths, err := s.authenticatorService.GetAuthenticatorsByUserID(context.Background(), user.ID)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), auths, 2)
}

func (s *AuthenticatorTestSuite) TestAuthenticatorExistence() {
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))
	s.seeder.CreateAuthenticator(user, models.AuthenticatorTypePassword, "hashedpass")

	// Test existing authenticator type
	exists, err := s.authenticatorService.ExistsByType(context.Background(), user.ID, models.AuthenticatorTypePassword)
	assert.NoError(s.T(), err)
	assert.True(s.T(), exists)

	// Test non-existing authenticator type
	exists, err = s.authenticatorService.ExistsByType(context.Background(), user.ID, models.AuthenticatorTypeGoogle)
	assert.NoError(s.T(), err)
	assert.False(s.T(), exists)
}

func (s *AuthenticatorTestSuite) TestUpdateAuthenticator() {
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))
	auth := s.seeder.CreateAuthenticator(user, models.AuthenticatorTypeGoogle, "oldtoken")

	// Test updating tokens
	err := s.authenticatorService.UpdateAuthenticatorTokens(context.Background(), auth.ID, "newtoken", "newrefreshtoken")
	assert.NoError(s.T(), err)

	// Verify update
	auths, err := s.authenticatorService.GetAuthenticatorsByUserID(context.Background(), user.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "newtoken", auths[0].Token)
	assert.Equal(s.T(), "newrefreshtoken", *auths[0].RefreshToken)
}

func stringPtr(s string) *string {
	return &s
}
