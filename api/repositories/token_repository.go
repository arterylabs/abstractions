package repositories

import (
	"abstractions/db/models"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenRepository struct {
	DB *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{DB: db}
}

func (r *TokenRepository) Create(ctx context.Context, token *models.Token) error {
	return r.DB.WithContext(ctx).Create(token).Error
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

func (r *TokenRepository) GetByValue(ctx context.Context, token string) (*models.Token, error) {
	tokenModel := &models.Token{}
	err := r.DB.WithContext(ctx).Where("token = ?", token).First(tokenModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return tokenModel, nil
}

func (r *TokenRepository) GetByTokenAndUserID(ctx context.Context, token string, userID uuid.UUID) (*models.Token, error) {
	tokenModel := &models.Token{}
	err := r.DB.WithContext(ctx).Where("token = ? AND user_id = ?", token, userID).First(tokenModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return tokenModel, nil
}

func (r *TokenRepository) ListValidByUserID(ctx context.Context, userID uuid.UUID) ([]models.Token, error) {
	tokens := []models.Token{}
	res := r.DB.WithContext(ctx).Where("user_id = ? AND expires_at > ?", userID, time.Now()).Find(&tokens)
	return tokens, res.Error
}

func (r *TokenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&models.Token{}).Error
}

func (r *TokenRepository) MarkUsed(ctx context.Context, id uuid.UUID, usedAt time.Time) error {
	return r.DB.WithContext(ctx).
		Model(&models.Token{}).
		Where("id = ?", id).
		Update("used_at", &usedAt).Error
}
