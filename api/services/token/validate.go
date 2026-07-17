package token_service

import (
	"abstractions/db/models"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func (s *Service) ValidateToken(ctx context.Context, token string) (bool, error) {
	tokenModel, err := s.repo.GetByValue(ctx, token)
	if err != nil {
		return false, fmt.Errorf("token_service: validate_token: %w", err)
	}
	if tokenModel == nil {
		return false, gorm.ErrRecordNotFound
	}
	// Token is invalid if it's expired or has been used
	return tokenModel.ExpiresAt.After(time.Now()) && tokenModel.UsedAt == nil, nil
}

func (s *Service) Validate(ctx context.Context, token *models.Token) (bool, error) {
	_ = ctx
	return token.ExpiresAt.After(time.Now()) && token.UsedAt == nil, nil
}

func (s *Service) ValidateTokenByType(ctx context.Context, token string, tokenType models.TokenType) (*models.Token, error) {
	t, err := s.GetByValue(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("token_service: validate_token_by_type: get_token: %w", err)
	}
	if t.Type != tokenType || !t.ExpiresAt.After(time.Now()) || t.UsedAt != nil {
		return nil, ErrTokenInvalid
	}
	return t, nil
}
