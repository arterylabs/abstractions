package session_service

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/utils/jwt_utils"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session has expired")
	ErrSessionInvalid  = errors.New("session is invalid")
	ErrInvalidToken    = errors.New("invalid authentication token")
)

type Service struct {
	DB *gorm.DB
}

type TokenValidationResult struct {
	UserID    uuid.UUID
	User      *models.User
	SessionID uuid.UUID
	Session   *models.Session
}

func NewSessionService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

// Basic CRUD operations
func (s *Service) CreateSession(session *models.Session) error {
	return s.DB.Create(session).Error
}

func (s *Service) GetSessionByID(sessionID uuid.UUID) (*models.Session, error) {
	var session models.Session
	err := s.DB.Where("id = ?", sessionID).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *Service) GetSessionIfActive(sessionID uuid.UUID) (*models.Session, error) {
	var session models.Session
	err := s.DB.Where("id = ? AND expires_at > now() AND invalidated_at IS NULL", sessionID).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *Service) IsSessionValid(sessionID uuid.UUID) (bool, error) {
	session, err := s.GetSessionByID(sessionID)
	if err != nil {
		return false, err
	}
	return session.IsValid(), nil
}

func (s *Service) InvalidateSession(sessionID uuid.UUID) error {
	session, err := s.GetSessionByID(sessionID)
	if err != nil {
		return err
	}

	now := time.Now()
	session.InvalidatedAt = &now
	err = s.DB.Save(session).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) InvalidateAllSessionsForUser(userID uuid.UUID) error {
	var sessions []models.Session
	err := s.DB.Where("user_id = ?", userID).Find(&sessions).Error
	if err != nil {
		return err
	}

	now := time.Now()
	for idx := range sessions {
		sessions[idx].InvalidatedAt = &now
	}

	err = s.DB.Save(&sessions).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) InvalidateAllImpersonationSessionsForAdmin(userID uuid.UUID) error {
	var sessions []models.Session
	err := s.DB.Where("impersonated_by = ? and invalidated_at IS NULL and expires_at > now()", userID).Find(&sessions).Error
	if err != nil {
		return err
	}

	now := time.Now()
	for idx := range sessions {
		sessions[idx].InvalidatedAt = &now
	}

	err = s.DB.Save(&sessions).Error
	if err != nil {
		return err
	}

	return nil
}

// Complex business logic
func (s *Service) ValidateToken(authToken string) (*TokenValidationResult, error) {
	userID, sessionID, err := jwt_utils.DecodeJWT(authToken)
	if err != nil {
		return nil, err
	}

	session, err := s.GetSessionIfActive(uuid.MustParse(sessionID))
	if err != nil {
		return nil, err
	}

	if session.UserID.String() != userID {
		return nil, nil
	}

	if !session.IsValid() {
		return nil, nil
	}

	userRepo := repositories.NewUserRepository(s.DB)
	user, err := userRepo.GetByID(session.UserID)
	if err != nil {
		return nil, err
	}

	return &TokenValidationResult{
		UserID:    session.UserID,
		User:      user,
		SessionID: session.ID,
		Session:   session,
	}, nil
}
