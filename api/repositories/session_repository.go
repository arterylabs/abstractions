package repositories

import (
	"abstractions/db/models"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionRepositoryQ interface {
	Create(ctx context.Context, session *models.Session) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Session, error)
	ListValidByUserID(ctx context.Context, userID uuid.UUID) ([]models.Session, error)
}

type SessionRepository struct {
	DB *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{DB: db}
}

func (r *SessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	var session models.Session
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) ListValidByUserID(ctx context.Context, userID uuid.UUID) ([]models.Session, error) {
	var sessions []models.Session
	err := r.DB.WithContext(ctx).
		Where("user_id = ? AND expires_at > now() AND revoked_at IS NULL", userID).
		Find(&sessions).Error
	return sessions, err
}

func (r *SessionRepository) Create(ctx context.Context, session *models.Session) error {
	return r.DB.WithContext(ctx).Create(session).Error
}
