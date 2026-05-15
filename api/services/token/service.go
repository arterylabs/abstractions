package token_service

import (
	"abstractions/db/models"
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
func (s *Service) CreateToken(token *models.Token) error {
	if token.ID == uuid.Nil {
		token.ID = uuid.New()
	}
	if err := s.DB.Create(token).Error; err != nil {
		return fmt.Errorf("token_service: create_token: %w", err)
	}
	return nil
}

func (s *Service) GetTokenByID(id uuid.UUID) (*models.Token, error) {
	token := &models.Token{}
	res := s.DB.Where("id = ?", id).Find(token)
	if res.Error != nil || token.ID == uuid.Nil {
		// If no token is found, return nil
		if res.Error == nil {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("token_service: get_token_by_id: %w", res.Error)
	}
	return token, nil
}

func (s *Service) FindTokenByTokenAndID(token string, id uuid.UUID) (*models.Token, error) {
	tokenModel := &models.Token{}
	res := s.DB.Where("token = ? AND user_id = ?", token, id).Find(tokenModel)
	if res.Error != nil {
		return nil, fmt.Errorf("token_service: find_token_by_token_and_id: %w", res.Error)
	}
	return tokenModel, nil
}

func (s *Service) GetTokenByValue(token string) (*models.Token, error) {
	tokenModel := &models.Token{}
	res := s.DB.Where("token = ?", token).Find(tokenModel)
	if res.Error != nil || tokenModel.ID == uuid.Nil {
		// If no token is found, return nil
		if res.Error == nil {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("token_service: get_token_by_value: %w", res.Error)
	}
	return tokenModel, nil
}

func (s *Service) DeleteToken(id uuid.UUID) error {
	if err := s.DB.Where("id = ?", id).Delete(&models.Token{}).Error; err != nil {
		return fmt.Errorf("token_service: delete_token: %w", err)
	}
	return nil
}

func (s *Service) MarkTokenAsUsed(tokenId uuid.UUID) error {
	now := time.Now()
	res := s.DB.Model(&models.Token{}).Where("id = ?", tokenId).Update("used_at", &now)
	if res.Error != nil {
		return fmt.Errorf("token_service: mark_token_as_used: %w", res.Error)
	}
	return nil
}
