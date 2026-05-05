package tests

import (
	"abstractions/db/models"
	token_service "abstractions/services/token"
	"abstractions/tests/testutils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TokenTestSuite struct {
	testutils.TestSuite
	seeder       *testutils.Seeder
	tokenService *token_service.Service
}

func TestTokenService(t *testing.T) {
	suite.Run(t, new(TokenTestSuite))
}

func (s *TokenTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.tokenService = token_service.NewTokenService(s.DB)
}

func (s *TokenTestSuite) TestCreateToken() {
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))

	// Test creating verification token
	verificationToken := &models.Token{
		UserID:    user.ID,
		Type:      models.TokenTypeEmailVerification,
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err := s.tokenService.CreateToken(verificationToken)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), verificationToken.ID)

	// Test creating password reset token
	resetToken := &models.Token{
		UserID:    user.ID,
		Type:      models.TokenTypePasswordReset,
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = s.tokenService.CreateToken(resetToken)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resetToken.ID)
}

func (s *TokenTestSuite) TestTokenValidation() {
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))

	// Create a valid token
	token := &models.Token{
		UserID:    user.ID,
		Type:      models.TokenTypeEmailVerification,
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	s.tokenService.CreateToken(token)

	// Test token validation
	valid, err := s.tokenService.ValidateToken(token.Token)
	assert.NoError(s.T(), err)
	assert.True(s.T(), valid)

	// Test invalid token
	valid, err = s.tokenService.ValidateToken("invalid-token")
	assert.Error(s.T(), err)
	assert.False(s.T(), valid)
}

func (s *TokenTestSuite) TestTokenExpiration() {
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))

	// Create an expired token
	expiredToken := &models.Token{
		UserID:    user.ID,
		Type:      models.TokenTypeEmailVerification,
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(-24 * time.Hour), // Expired
	}
	s.tokenService.CreateToken(expiredToken)

	// Test expired token validation
	valid, err := s.tokenService.Validate(expiredToken)
	assert.NoError(s.T(), err)
	assert.False(s.T(), valid)
}

func (s *TokenTestSuite) TestTokenTypeValidation() {
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))

	// Create a token
	token := &models.Token{
		UserID:    user.ID,
		Type:      models.TokenTypePasswordReset,
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	s.tokenService.CreateToken(token)

	// Test correct token type
	validToken, err := s.tokenService.ValidateTokenByType(token.Token, models.TokenTypePasswordReset)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), validToken)

	// Test incorrect token type
	invalidToken, err := s.tokenService.ValidateTokenByType(token.Token, models.TokenTypeEmailVerification)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), invalidToken)
}

func (s *TokenTestSuite) TestTokenUsage() {
	user := s.seeder.CreateUser(testutils.WithEmail("test@example.com"))

	// Create a token
	token := &models.Token{
		UserID:    user.ID,
		Type:      models.TokenTypePasswordReset,
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	s.tokenService.CreateToken(token)

	// Test marking token as used
	err := s.tokenService.MarkTokenAsUsed(token.ID)
	assert.NoError(s.T(), err)

	// Verify token cannot be used again
	validToken, err := s.tokenService.ValidateTokenByType(token.Token, models.TokenTypePasswordReset)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), validToken)
}
