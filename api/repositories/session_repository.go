package repositories

import (
	"abstractions/db/models"

	"gorm.io/gorm"
)

type SessionRepositoryQ interface {
	Create(session *models.Session) error
	GetSessionByID(session *models.Session, id string)
	GetValidSessionsByUserID(sessions *[]models.Session, userID string)
}

type SessionRepository struct {
	DB *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{DB: db}
}

func (sessionRepository *SessionRepository) GetSessionByID(session *models.Session, id string) {
	sessionRepository.DB.Where("id = ?", id).Find(session)
}

func (sessionRepository *SessionRepository) GetValidSessionsByUserID(sessions *[]models.Session, userID string) {
	sessionRepository.DB.Where("user_id = ? AND valid = ?", userID, true).Find(sessions)
}

func (sessionRepository *SessionRepository) Create(session *models.Session) error {
	return sessionRepository.DB.Create(session).Error
}
