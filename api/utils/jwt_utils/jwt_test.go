package jwt_utils_test

import (
	"testing"
	"time"

	"abstractions/config"
	"abstractions/utils/jwt_utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	config.Environment = &config.ENV{SigningKey: "test_signing_key_for_jwt_tokens_123456"}
}

func TestGenerateAndDecodeJWT_RoundTrip(t *testing.T) {
	userID := uuid.New().String()
	sessionID := uuid.New().String()
	exp := time.Now().Add(time.Hour).UTC().Truncate(time.Second)

	token, err := jwt_utils.GenerateJWT(userID, sessionID, exp)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	gotUser, gotSession, err := jwt_utils.DecodeJWT(token)
	require.NoError(t, err)
	assert.Equal(t, userID, gotUser)
	assert.Equal(t, sessionID, gotSession)

	gotExp, err := jwt_utils.GetJWTExpiry(token)
	require.NoError(t, err)
	assert.Equal(t, exp.Unix(), gotExp.Unix())
}

func TestDecodeJWT_TamperedTokenRejected(t *testing.T) {
	token, err := jwt_utils.GenerateJWT("u1", "s1", time.Now().Add(time.Hour))
	require.NoError(t, err)

	// Flip a character in the signature segment.
	tampered := token[:len(token)-4] + "xxxx"
	_, _, err = jwt_utils.DecodeJWT(tampered)
	assert.Error(t, err)
}

func TestDecodeJWT_ExpiredTokenRejected(t *testing.T) {
	token, err := jwt_utils.GenerateJWT("u1", "s1", time.Now().Add(-time.Hour))
	require.NoError(t, err)

	_, _, err = jwt_utils.DecodeJWT(token)
	assert.Error(t, err)
}

func TestDecodeJWT_GarbageRejected(t *testing.T) {
	_, _, err := jwt_utils.DecodeJWT("not.a.jwt")
	assert.Error(t, err)
}

func TestDecodeJWT_WrongSigningMethodRejected(t *testing.T) {
	// Build an unsigned "none" alg token; DecodeJWT must reject non-HMAC methods.
	claims := &jwt_utils.Claims{
		UserID:    "u1",
		SessionID: "s1",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	signed, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	_, _, err = jwt_utils.DecodeJWT(signed)
	assert.Error(t, err)
}

func TestGetJWTExpiry_InvalidToken(t *testing.T) {
	_, err := jwt_utils.GetJWTExpiry("bad-token")
	assert.Error(t, err)
}
