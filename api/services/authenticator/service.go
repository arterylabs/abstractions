package authenticator_service

import (
	"abstractions/db/models"
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
	DB *gorm.DB
}

func NewAuthenticatorService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

// Basic CRUD operations
func (s *Service) CreateAuthenticator(authenticator *models.Authenticator) error {
	if err := s.DB.Create(authenticator).Error; err != nil {
		return fmt.Errorf("authenticator_service: create_authenticator: %w", err)
	}
	return nil
}

func (s *Service) GetAuthenticatorByType(userID uuid.UUID, authType models.AuthenticatorType) (*models.Authenticator, error) {
	var authenticator models.Authenticator
	err := s.DB.Where("user_id = ? AND type = ?", userID, authType).First(&authenticator).Error
	if err != nil {
		return nil, fmt.Errorf("authenticator_service: get_by_type: %w", err)
	}
	return &authenticator, nil
}

func (s *Service) GetAuthenticatorsByUserID(userID uuid.UUID) ([]models.Authenticator, error) {
	var authenticators []models.Authenticator
	err := s.DB.Where("user_id = ?", userID).Find(&authenticators).Error
	if err != nil {
		return nil, fmt.Errorf("authenticator_service: get_by_user_id: %w", err)
	}
	return authenticators, nil
}

func (s *Service) ExistsByType(userID uuid.UUID, authenticatorType models.AuthenticatorType) (bool, error) {
	var count int64
	err := s.DB.Model(&models.Authenticator{}).Where("user_id = ? AND type = ?", userID, authenticatorType).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("authenticator_service: exists_by_type: %w", err)
	}
	return count > 0, nil
}

func (s *Service) UpdateAuthenticatorTokens(authenticatorID uuid.UUID, token string, refreshToken string) error {
	authenticator := models.Authenticator{}
	if err := s.DB.Where("id = ?", authenticatorID).First(&authenticator).Error; err != nil {
		return fmt.Errorf("authenticator_service: update_authenticator_tokens: get_authenticator: %w", err)
	}

	authenticator.Token = token
	if refreshToken != "" {
		authenticator.RefreshToken = &refreshToken
	} else {
		authenticator.RefreshToken = nil
	}

	if err := s.DB.Save(&authenticator).Error; err != nil {
		return fmt.Errorf("authenticator_service: update_authenticator_tokens: save: %w", err)
	}
	return nil
}
