package token_service

import (
	"abstractions/db/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrTokenNotFound = errors.New("token not found")
	ErrTokenExpired  = errors.New("token has expired")
	ErrTokenInvalid  = errors.New("token is invalid")
)

type Service struct {
	DB *gorm.DB
}

func NewTokenService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

// Basic CRUD operations
func (s *Service) CreateToken(ctx context.Context, token *models.Token) error {
	if token.ID == uuid.Nil {
		token.ID = uuid.New()
	}
	if err := s.DB.WithContext(ctx).Create(token).Error; err != nil {
		return fmt.Errorf("token_service: create_token: %w", err)
	}
	return nil
}

func (s *Service) GetTokenByID(ctx context.Context, id uuid.UUID) (*models.Token, error) {
	token := &models.Token{}
	res := s.DB.WithContext(ctx).Where("id = ?", id).Find(token)
	if res.Error != nil || token.ID == uuid.Nil {
		// If no token is found, return nil
		if res.Error == nil {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("token_service: get_token_by_id: %w", res.Error)
	}
	return token, nil
}

func (s *Service) FindTokenByTokenAndID(ctx context.Context, token string, id uuid.UUID) (*models.Token, error) {
	tokenModel := &models.Token{}
	res := s.DB.WithContext(ctx).Where("token = ? AND user_id = ?", token, id).Find(tokenModel)
	if res.Error != nil {
		return nil, fmt.Errorf("token_service: find_token_by_token_and_id: %w", res.Error)
	}
	return tokenModel, nil
}

func (s *Service) GetTokenByValue(ctx context.Context, token string) (*models.Token, error) {
	tokenModel := &models.Token{}
	res := s.DB.WithContext(ctx).Where("token = ?", token).Find(tokenModel)
	if res.Error != nil || tokenModel.ID == uuid.Nil {
		// If no token is found, return nil
		if res.Error == nil {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("token_service: get_token_by_value: %w", res.Error)
	}
	return tokenModel, nil
}

func (s *Service) DeleteToken(ctx context.Context, id uuid.UUID) error {
	if err := s.DB.WithContext(ctx).Where("id = ?", id).Delete(&models.Token{}).Error; err != nil {
		return fmt.Errorf("token_service: delete_token: %w", err)
	}
	return nil
}

func (s *Service) MarkTokenAsUsed(ctx context.Context, tokenID uuid.UUID) error {
	now := time.Now()
	res := s.DB.WithContext(ctx).Model(&models.Token{}).Where("id = ?", tokenID).Update("used_at", &now)
	if res.Error != nil {
		return fmt.Errorf("token_service: mark_token_as_used: %w", res.Error)
	}
	return nil
}

// CreateOAuthExchangeCode stores the session JWT under a short-lived, single-use exchange token.
// The returned token's ID (UUID) is used as the opaque code sent to the frontend in the redirect URL.
func (s *Service) CreateOAuthExchangeCode(ctx context.Context, userID uuid.UUID, sessionToken string) (*models.Token, error) {
	t := &models.Token{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      models.TokenTypeOAuthExchange,
		Token:     sessionToken,
		ExpiresAt: time.Now().UTC().Add(60 * time.Second),
	}
	if err := s.CreateToken(ctx, t); err != nil {
		return nil, fmt.Errorf("token_service: create_oauth_exchange_code: create_token: %w", err)
	}
	return t, nil
}

// ValidateOAuthExchangeCode looks up an exchange code by its UUID, checks type, expiry, and
// that it has not already been consumed. Returns the token so the caller can read the session JWT.
func (s *Service) ValidateOAuthExchangeCode(ctx context.Context, code string) (*models.Token, error) {
	id, err := uuid.Parse(code)
	if err != nil {
		return nil, ErrTokenInvalid
	}
	token, err := s.GetTokenByID(ctx, id)
	if err != nil {
		return nil, ErrTokenNotFound
	}
	if token.Type != models.TokenTypeOAuthExchange {
		return nil, ErrTokenInvalid
	}
	if time.Now().UTC().After(token.ExpiresAt) {
		return nil, ErrTokenExpired
	}
	if token.UsedAt != nil {
		return nil, ErrTokenInvalid
	}
	return token, nil
}
