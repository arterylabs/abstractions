package token_service

import (
	"abstractions/db/models"
	"context"
	"errors"
	"fmt"
	"time"

	"abstractions/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrTokenNotFound = errors.New("token not found")
	ErrTokenExpired  = errors.New("token has expired")
	ErrTokenInvalid  = errors.New("token is invalid")
)

type Service struct {
	db   *gorm.DB
	repo *repositories.TokenRepository
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db, repo: repositories.NewTokenRepository(db)}
}

// Basic CRUD operations
func (s *Service) Create(ctx context.Context, token *models.Token) error {
	if token.ID == uuid.Nil {
		token.ID = uuid.New()
	}
	if err := s.repo.Create(ctx, token); err != nil {
		return fmt.Errorf("token_service: create: %w", err)
	}
	return nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Token, error) {
	token, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("token_service: get_by_id: %w", err)
	}
	if token == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return token, nil
}

func (s *Service) FindTokenByTokenAndID(ctx context.Context, token string, id uuid.UUID) (*models.Token, error) {
	tokenModel, err := s.repo.GetByTokenAndUserID(ctx, token, id)
	if err != nil {
		return nil, fmt.Errorf("token_service: find_token_by_token_and_id: %w", err)
	}
	return tokenModel, nil
}

func (s *Service) GetByValue(ctx context.Context, token string) (*models.Token, error) {
	tokenModel, err := s.repo.GetByValue(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("token_service: get_by_value: %w", err)
	}
	if tokenModel == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return tokenModel, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("token_service: delete: %w", err)
	}
	return nil
}

func (s *Service) MarkTokenAsUsed(ctx context.Context, tokenID uuid.UUID) error {
	now := time.Now()
	if err := s.repo.MarkUsed(ctx, tokenID, now); err != nil {
		return fmt.Errorf("token_service: mark_token_as_used: %w", err)
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
	if err := s.Create(ctx, t); err != nil {
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
	token, err := s.GetByID(ctx, id)
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
