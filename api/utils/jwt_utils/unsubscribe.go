package jwt_utils

import (
	"abstractions/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// UnsubscribeClaims holds the team member ID for an unsubscribe JWT.
type UnsubscribeClaims struct {
	TeamMemberID string `json:"team_member_id"`
	jwt.StandardClaims
}

// GenerateUnsubscribeToken creates a signed JWT for a newsletter unsubscribe link.
// The token does not expire (set to 10 years) so recipients can always unsubscribe from old emails.
func GenerateUnsubscribeToken(teamMemberID uuid.UUID) (string, error) {
	signingKey := []byte(config.Environment.SigningKey)
	claims := &UnsubscribeClaims{
		TeamMemberID: teamMemberID.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(10 * 365 * 24 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

// ParseUnsubscribeToken validates a signed unsubscribe JWT and returns the team member ID.
func ParseUnsubscribeToken(tokenString string) (uuid.UUID, error) {
	signingKey := []byte(config.Environment.SigningKey)
	claims := &UnsubscribeClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, errors.New("invalid unsubscribe token")
	}
	id, err := uuid.Parse(claims.TeamMemberID)
	if err != nil {
		return uuid.Nil, errors.New("invalid team_member_id in token")
	}
	return id, nil
}
