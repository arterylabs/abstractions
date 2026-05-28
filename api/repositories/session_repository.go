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

func (r *SessionRepository) GetSessionByID(session *models.Session, id string) {
	r.DB.Where("id = ?", id).Find(session)
}

func (r *SessionRepository) GetValidSessionsByUserID(sessions *[]models.Session, userID string) {
	r.DB.Where("user_id = ? AND valid = ?", userID, true).Find(sessions)
}

func (r *SessionRepository) Create(session *models.Session) error {
	return r.DB.Create(session).Error
}
