package jwt_utils

import (
	"abstractions/config"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	jwt.StandardClaims
}

func GenerateJWT(userID string, sessionID string, expirationTime time.Time) (string, error) {
	var signingKey = []byte(config.Environment.SigningKey)
	claims := &Claims{
		UserID:    userID,
		SessionID: sessionID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func DecodeJWT(tokenString string) (string, string, error) {
	var signingKey = []byte(config.Environment.SigningKey)
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Reject any token not signed with HMAC (e.g., RSA, ECDSA, or none).
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})
	if err != nil {
		return "", "", errors.New("Unauthorized")
	}
	if !token.Valid {
		return "", "", errors.New("Invalid Token")
	}
	return claims.UserID, claims.SessionID, nil
}
