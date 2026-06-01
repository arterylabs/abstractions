package repositories

import (
	"abstractions/db/models"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthenticatorRepositoryQ interface {
	Create(ctx context.Context, authenticator *models.Authenticator) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Authenticator, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]models.Authenticator, error)
	ExistsByType(ctx context.Context, userID uuid.UUID, authType models.AuthenticatorType) bool
	GetByType(ctx context.Context, userID uuid.UUID, authType models.AuthenticatorType) (*models.Authenticator, error)
	ListTypesForUser(ctx context.Context, userID uuid.UUID) ([]models.AuthenticatorType, error)
	Update(ctx context.Context, authenticator *models.Authenticator) error
}

type AuthenticatorRepository struct {
	DB *gorm.DB
}

func NewAuthenticatorRepository(db *gorm.DB) *AuthenticatorRepository {
	return &AuthenticatorRepository{DB: db}
}

func (r *AuthenticatorRepository) Create(ctx context.Context, authenticator *models.Authenticator) error {
	return r.DB.WithContext(ctx).Create(authenticator).Error
}

func (r *AuthenticatorRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Authenticator, error) {
	authenticator := &models.Authenticator{}
	res := r.DB.WithContext(ctx).Where("id = ?", id).First(authenticator)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return authenticator, res.Error
}

func (r *AuthenticatorRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]models.Authenticator, error) {
	authenticators := []models.Authenticator{}
	res := r.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&authenticators)
	return authenticators, res.Error
}

func (r *AuthenticatorRepository) ExistsByType(ctx context.Context, userID uuid.UUID, authType models.AuthenticatorType) bool {
	var count int64
	r.DB.WithContext(ctx).Model(&models.Authenticator{}).Where("user_id = ? AND type = ?", userID, authType).Count(&count)
	return count > 0
}

func (r *AuthenticatorRepository) GetByType(ctx context.Context, userID uuid.UUID, authType models.AuthenticatorType) (*models.Authenticator, error) {
	authenticator := &models.Authenticator{}
	res := r.DB.WithContext(ctx).Where("user_id = ? AND type = ?", userID, authType).First(authenticator)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return authenticator, res.Error
}

func (r *AuthenticatorRepository) ListTypesForUser(ctx context.Context, userID uuid.UUID) ([]models.AuthenticatorType, error) {
	authenticators, err := r.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("authenticator_repository: get_authenticator_types_for_user: %w", err)
	}

	authenticatorTypes := []models.AuthenticatorType{}
	for _, authenticator := range authenticators {
		authenticatorTypes = append(authenticatorTypes, authenticator.Type)
	}

	return authenticatorTypes, nil
}

func (r *AuthenticatorRepository) Update(ctx context.Context, authenticator *models.Authenticator) error {
	return r.DB.WithContext(ctx).Save(authenticator).Error
}
