package repositories

import (
	"abstractions/db/models"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthenticatorRepositoryQ interface {
	GetAuthenticatorsByUserID(userID uuid.UUID) (*[]models.Authenticator, error)
	HasAuthenticatorByType(userID uuid.UUID, authType models.AuthenticatorType) bool
	GetAuthenticatorByType(userID uuid.UUID, authType models.AuthenticatorType) (*models.Authenticator, error)
	GetAuthenticatorTypesForUser(userID uuid.UUID) ([]models.AuthenticatorType, error)
	Update(authenticator *models.Authenticator) error
}

type AuthenticatorRepository struct {
	DB *gorm.DB
}

func NewAuthenticatorRepository(db *gorm.DB) *AuthenticatorRepository {
	return &AuthenticatorRepository{DB: db}
}

func (r *AuthenticatorRepository) GetAuthenticatorsByUserID(userID uuid.UUID) (*[]models.Authenticator, error) {
	authenticators := []models.Authenticator{}
	res := r.DB.Where("user_id = ?", userID).Find(&authenticators)
	return &authenticators, res.Error
}

func (r *AuthenticatorRepository) HasAuthenticatorByType(userID uuid.UUID, authType models.AuthenticatorType) bool {
	var count int64
	r.DB.Model(&models.Authenticator{}).Where("user_id = ? AND type = ?", userID, authType).Count(&count)
	return count > 0
}

func (r *AuthenticatorRepository) GetAuthenticatorByType(userID uuid.UUID, authType models.AuthenticatorType) (*models.Authenticator, error) {
	authenticator := &models.Authenticator{}
	res := r.DB.Where("user_id = ? AND type = ?", userID, authType).First(authenticator)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return authenticator, res.Error
}

func (r *AuthenticatorRepository) GetAuthenticatorTypesForUser(userID uuid.UUID) ([]models.AuthenticatorType, error) {
	authenticators, err := r.GetAuthenticatorsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("authenticator_repository: get_authenticator_types_for_user: %w", err)
	}

	authenticatorTypes := []models.AuthenticatorType{}
	for _, authenticator := range *authenticators {
		authenticatorTypes = append(authenticatorTypes, authenticator.Type)
	}

	return authenticatorTypes, nil
}

func (r *AuthenticatorRepository) Update(authenticator *models.Authenticator) error {
	return r.DB.Save(authenticator).Error
}
