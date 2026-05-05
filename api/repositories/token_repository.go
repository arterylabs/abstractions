package repositories

import (
	"abstractions/db/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenRepositoryQ interface {
	GetTokenByID(id uuid.UUID) (*models.Token, error)
}

type TokenRepository struct {
	DB *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{DB: db}
}

func (tokenRepository *TokenRepository) GetTokenByID(id uuid.UUID) (*models.Token, error) {
	token := &models.Token{}
	res := tokenRepository.DB.Where("id = ?", id).Find(token)
	return token, res.Error
}

func (tokenRepository *TokenRepository) GetValidTokensByUserID(userID uuid.UUID) (*[]models.Token, error) {
	tokens := []models.Token{}
	res := tokenRepository.DB.Where("user_id = ? AND expires_at > ?", userID, time.Now()).Find(tokens)
	return &tokens, res.Error
}
