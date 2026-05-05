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
		return signingKey, nil
	})
	if err != nil {
		fmt.Println(err)
		return "", "", errors.New("Unauthorized")
	}
	if !token.Valid {
		fmt.Println(err)
		return "", "", errors.New("Invalid Token")
	}
	return claims.UserID, claims.SessionID, nil
}
