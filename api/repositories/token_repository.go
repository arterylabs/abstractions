package repositories

import (
	"abstractions/db/models"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenRepositoryQ interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.Token, error)
}

type TokenRepository struct {
	DB *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{DB: db}
}

func (r *TokenRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Token, error) {
	token := &models.Token{}
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return token, nil
}

func (r *TokenRepository) ListValidByUserID(ctx context.Context, userID uuid.UUID) ([]models.Token, error) {
	tokens := []models.Token{}
	res := r.DB.WithContext(ctx).Where("user_id = ? AND expires_at > ?", userID, time.Now()).Find(&tokens)
	return tokens, res.Error
}
