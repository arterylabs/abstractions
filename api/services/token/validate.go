package token_service

import (
	"abstractions/db/models"
	"context"
	"fmt"
	"time"
)

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
