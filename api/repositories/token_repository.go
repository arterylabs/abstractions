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

func (r *TokenRepository) GetTokenByID(id uuid.UUID) (*models.Token, error) {
	token := &models.Token{}
	res := r.DB.Where("id = ?", id).Find(token)
	return token, res.Error
}

func (r *TokenRepository) GetValidTokensByUserID(userID uuid.UUID) (*[]models.Token, error) {
	tokens := []models.Token{}
	res := r.DB.Where("user_id = ? AND expires_at > ?", userID, time.Now()).Find(&tokens)
	return &tokens, res.Error
}
