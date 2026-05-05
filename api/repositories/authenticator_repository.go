package repositories

import (
	"abstractions/db/models"

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

func (authenticatorRepository *AuthenticatorRepository) GetAuthenticatorsByUserID(userID uuid.UUID) (*[]models.Authenticator, error) {
	authenticators := []models.Authenticator{}
	res := authenticatorRepository.DB.Where("user_id = ?", userID).Find(&authenticators)
	return &authenticators, res.Error
}

func (authenticatorRepository *AuthenticatorRepository) HasAuthenticatorByType(userID uuid.UUID, authType models.AuthenticatorType) bool {
	var count int64
	authenticatorRepository.DB.Model(&models.Authenticator{}).Where("user_id = ? AND type = ?", userID, authType).Count(&count)
	return count > 0
}

func (authenticatorRepository *AuthenticatorRepository) GetAuthenticatorByType(userID uuid.UUID, authType models.AuthenticatorType) (*models.Authenticator, error) {
	authenticator := &models.Authenticator{}
	res := authenticatorRepository.DB.Where("user_id = ? AND type = ?", userID, authType).First(authenticator)
	if res.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return authenticator, res.Error
}

func (authenticatorRepository *AuthenticatorRepository) GetAuthenticatorTypesForUser(userID uuid.UUID) ([]models.AuthenticatorType, error) {
	authenticators, err := authenticatorRepository.GetAuthenticatorsByUserID(userID)
	if err != nil {
		return nil, err
	}

	authenticatorTypes := []models.AuthenticatorType{}
	for _, authenticator := range *authenticators {
		authenticatorTypes = append(authenticatorTypes, authenticator.Type)
	}

	return authenticatorTypes, nil
}

func (authenticatorRepository *AuthenticatorRepository) Update(authenticator *models.Authenticator) error {
	return authenticatorRepository.DB.Save(authenticator).Error
}
