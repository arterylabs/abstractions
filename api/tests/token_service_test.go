package tests

import (
	"abstractions/db/models"
	token_service "abstractions/services/token"
	"abstractions/tests/testutils"
	"context"
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

	err := s.tokenService.CreateToken(context.Background(), verificationToken)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), verificationToken.ID)

	// Test creating password reset token
	resetToken := &models.Token{
		UserID:    user.ID,
		Type:      models.TokenTypePasswordReset,
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = s.tokenService.CreateToken(context.Background(), resetToken)
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
	s.tokenService.CreateToken(context.Background(), token)

	// Test token validation
	valid, err := s.tokenService.ValidateToken(context.Background(), token.Token)
	assert.NoError(s.T(), err)
	assert.True(s.T(), valid)

	// Test invalid token
	valid, err = s.tokenService.ValidateToken(context.Background(), "invalid-token")
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
	s.tokenService.CreateToken(context.Background(), expiredToken)

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
	s.tokenService.CreateToken(context.Background(), token)

	// Test correct token type
	validToken, err := s.tokenService.ValidateTokenByType(context.Background(), token.Token, models.TokenTypePasswordReset)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), validToken)

	// Test incorrect token type
	invalidToken, err := s.tokenService.ValidateTokenByType(context.Background(), token.Token, models.TokenTypeEmailVerification)
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
	s.tokenService.CreateToken(context.Background(), token)

	// Test marking token as used
	err := s.tokenService.MarkTokenAsUsed(context.Background(), token.ID)
	assert.NoError(s.T(), err)

	// Verify token cannot be used again
	validToken, err := s.tokenService.ValidateTokenByType(context.Background(), token.Token, models.TokenTypePasswordReset)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), validToken)
}

// TestCreateAndValidateOAuthExchangeCode verifies the happy-path OAuth exchange flow:
// CreateOAuthExchangeCode persists a short-lived token and ValidateOAuthExchangeCode
// returns it by its UUID code.
func (s *TokenTestSuite) TestCreateAndValidateOAuthExchangeCode() {
	user := s.seeder.CreateUser(testutils.WithEmail("oauth-happy@example.com"))
	sessionJWT := "fake.jwt.token"

	exchangeToken, err := s.tokenService.CreateOAuthExchangeCode(context.Background(), user.ID, sessionJWT)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), exchangeToken)
	assert.Equal(s.T(), models.TokenTypeOAuthExchange, exchangeToken.Type)
	assert.Equal(s.T(), sessionJWT, exchangeToken.Token)

	// ValidateOAuthExchangeCode uses the token UUID as the opaque code.
	validated, err := s.tokenService.ValidateOAuthExchangeCode(context.Background(), exchangeToken.ID.String())
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), validated)
	assert.Equal(s.T(), exchangeToken.ID, validated.ID)
	assert.Equal(s.T(), sessionJWT, validated.Token)
}

// TestValidateOAuthExchangeCode_AlreadyConsumed verifies that a used exchange code cannot be
// validated a second time.
func (s *TokenTestSuite) TestValidateOAuthExchangeCode_AlreadyConsumed() {
	user := s.seeder.CreateUser(testutils.WithEmail("oauth-consumed@example.com"))

	exchangeToken, err := s.tokenService.CreateOAuthExchangeCode(context.Background(), user.ID, "jwt-value")
	assert.NoError(s.T(), err)

	// Mark it as used, simulating the first successful consumption.
	err = s.tokenService.MarkTokenAsUsed(context.Background(), exchangeToken.ID)
	assert.NoError(s.T(), err)

	// Second validation must fail.
	_, err = s.tokenService.ValidateOAuthExchangeCode(context.Background(), exchangeToken.ID.String())
	assert.Error(s.T(), err)
}

// TestValidateOAuthExchangeCode_WrongType verifies that a token of a different type is rejected
// even if its UUID happens to be passed as an exchange code.
func (s *TokenTestSuite) TestValidateOAuthExchangeCode_WrongType() {
	user := s.seeder.CreateUser(testutils.WithEmail("oauth-wrongtype@example.com"))

	// Create a password-reset token (not an OAuth exchange token).
	wrongType := &models.Token{
		UserID:    user.ID,
		Type:      models.TokenTypePasswordReset,
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(time.Hour),
	}
	err := s.tokenService.CreateToken(context.Background(), wrongType)
	assert.NoError(s.T(), err)

	_, err = s.tokenService.ValidateOAuthExchangeCode(context.Background(), wrongType.ID.String())
	assert.Error(s.T(), err)
}
