package token_service

import (
	"abstractions/db/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *Service) ValidateToken(token string) (bool, error) {
	tokenModel := &models.Token{}
	res := s.DB.Where("token = ?", token).Find(tokenModel)
	if res.Error != nil || tokenModel.ID == uuid.Nil {
		// If no token is found, return nil
		if res.Error == nil {
			return false, gorm.ErrRecordNotFound
		}
		return false, res.Error
	}
	// Token is invalid if it's expired or has been used
	return tokenModel.ExpiresAt.After(time.Now()) && tokenModel.UsedAt == nil, nil
}

func (s *Service) Validate(token *models.Token) (bool, error) {
	return token.ExpiresAt.After(time.Now()) && token.UsedAt == nil, nil
}

func (s *Service) ValidateTokenByType(token string, tokenType models.TokenType) (*models.Token, error) {
	t, err := s.GetTokenByValue(token)
	if err != nil || t.Type != tokenType || !t.ExpiresAt.After(time.Now()) || t.UsedAt != nil {
		return nil, errors.New("invalid token")
	}
	return t, nil
}
