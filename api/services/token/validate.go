package token_service

import (
	"abstractions/db/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *Service) ValidateToken(ctx context.Context, token string) (bool, error) {
	tokenModel := &models.Token{}
	res := s.DB.WithContext(ctx).Where("token = ?", token).Find(tokenModel)
	if res.Error != nil || tokenModel.ID == uuid.Nil {
		// If no token is found, return nil
		if res.Error == nil {
			return false, gorm.ErrRecordNotFound
		}
		return false, fmt.Errorf("token_service: validate_token: %w", res.Error)
	}
	// Token is invalid if it's expired or has been used
	return tokenModel.ExpiresAt.After(time.Now()) && tokenModel.UsedAt == nil, nil
}

func (s *Service) Validate(token *models.Token) (bool, error) {
	return token.ExpiresAt.After(time.Now()) && token.UsedAt == nil, nil
}

func (s *Service) ValidateTokenByType(ctx context.Context, token string, tokenType models.TokenType) (*models.Token, error) {
	t, err := s.GetTokenByValue(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("token_service: validate_token_by_type: get_token: %w", err)
	}
	if t.Type != tokenType || !t.ExpiresAt.After(time.Now()) || t.UsedAt != nil {
		return nil, ErrTokenInvalid
	}
	return t, nil
}
