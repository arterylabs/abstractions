package authenticator_service

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrAuthenticatorNotFound    = errors.New("authenticator not found")
	ErrAuthenticatorExists      = errors.New("authenticator already exists for this user and type")
	ErrInvalidAuthenticatorType = errors.New("invalid authenticator type")
)

type Service struct {
	db   *gorm.DB
	repo *repositories.AuthenticatorRepository
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db, repo: repositories.NewAuthenticatorRepository(db)}
}

// Basic CRUD operations
func (s *Service) Create(ctx context.Context, authenticator *models.Authenticator) error {
	if err := s.repo.Create(ctx, authenticator); err != nil {
		return fmt.Errorf("authenticator_service: create: %w", err)
	}
	return nil
}

func (s *Service) GetAuthenticatorByType(ctx context.Context, userID uuid.UUID, authType models.AuthenticatorType) (*models.Authenticator, error) {
	authenticator, err := s.repo.GetByType(ctx, userID, authType)
	if err != nil {
		return nil, fmt.Errorf("authenticator_service: get_by_type: %w", err)
	}
	if authenticator == nil {
		return nil, ErrAuthenticatorNotFound
	}
	return authenticator, nil
}

func (s *Service) GetAuthenticatorsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Authenticator, error) {
	authenticators, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("authenticator_service: get_by_user_id: %w", err)
	}
	return authenticators, nil
}

func (s *Service) ExistsByType(ctx context.Context, userID uuid.UUID, authenticatorType models.AuthenticatorType) (bool, error) {
	return s.repo.ExistsByType(ctx, userID, authenticatorType), nil
}

func (s *Service) UpdateAuthenticatorTokens(ctx context.Context, authenticatorID uuid.UUID, token string, refreshToken string) error {
	authenticator, err := s.repo.GetByID(ctx, authenticatorID)
	if err != nil {
		return fmt.Errorf("authenticator_service: update_authenticator_tokens: get_authenticator: %w", err)
	}
	if authenticator == nil {
		return ErrAuthenticatorNotFound
	}

	authenticator.Token = token
	if refreshToken != "" {
		authenticator.RefreshToken = &refreshToken
	} else {
		authenticator.RefreshToken = nil
	}

	if err := s.repo.Update(ctx, authenticator); err != nil {
		return fmt.Errorf("authenticator_service: update_authenticator_tokens: save: %w", err)
	}
	return nil
}
