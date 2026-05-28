package session_service

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/utils/jwt_utils"
	"errors"
	"fmt"
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
	if err := s.DB.Create(session).Error; err != nil {
		return fmt.Errorf("session_service: create_session: %w", err)
	}
	return nil
}

func (s *Service) GetSessionByID(sessionID uuid.UUID) (*models.Session, error) {
	var session models.Session
	err := s.DB.Where("id = ?", sessionID).First(&session).Error
	if err != nil {
		return nil, fmt.Errorf("session_service: get_session_by_id: %w", err)
	}
	return &session, nil
}

func (s *Service) GetSessionIfActive(sessionID uuid.UUID) (*models.Session, error) {
	var session models.Session
	err := s.DB.Where("id = ? AND expires_at > now() AND revoked_at IS NULL", sessionID).First(&session).Error
	if err != nil {
		return nil, fmt.Errorf("session_service: get_session_if_active: %w", err)
	}
	return &session, nil
}

func (s *Service) IsSessionValid(sessionID uuid.UUID) (bool, error) {
	session, err := s.GetSessionByID(sessionID)
	if err != nil {
		return false, fmt.Errorf("session_service: is_session_valid: %w", err)
	}
	return session.IsValid(), nil
}

func (s *Service) InvalidateSession(sessionID uuid.UUID) error {
	session, err := s.GetSessionByID(sessionID)
	if err != nil {
		return fmt.Errorf("session_service: invalidate_session: get_session: %w", err)
	}

	now := time.Now()
	session.RevokedAt = &now
	if err = s.DB.Save(session).Error; err != nil {
		return fmt.Errorf("session_service: invalidate_session: save: %w", err)
	}

	return nil
}

func (s *Service) InvalidateAllSessionsForUser(userID uuid.UUID) error {
	now := time.Now()
	if err := s.DB.Model(&models.Session{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error; err != nil {
		return fmt.Errorf("session_service: invalidate_all_sessions_for_user: %w", err)
	}
	return nil
}

func (s *Service) InvalidateAllImpersonationSessionsForAdmin(userID uuid.UUID) error {
	now := time.Now()
	if err := s.DB.Model(&models.Session{}).
		Where("impersonated_by = ? AND revoked_at IS NULL AND expires_at > now()", userID).
		Update("revoked_at", now).Error; err != nil {
		return fmt.Errorf("session_service: invalidate_all_impersonation_sessions_for_admin: %w", err)
	}
	return nil
}

// Complex business logic
func (s *Service) ValidateToken(authToken string) (*TokenValidationResult, error) {
	userID, sessionID, err := jwt_utils.DecodeJWT(authToken)
	if err != nil {
		return nil, fmt.Errorf("session_service: validate_token: decode_jwt: %w", err)
	}

	sessionUUID, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session_service: validate_token: invalid session id in token: %w", err)
	}

	session, err := s.GetSessionIfActive(sessionUUID)
	if err != nil {
		return nil, fmt.Errorf("session_service: validate_token: get_session: %w", err)
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
		return nil, fmt.Errorf("session_service: validate_token: get_user: %w", err)
	}

	return &TokenValidationResult{
		UserID:    session.UserID,
		User:      user,
		SessionID: session.ID,
		Session:   session,
	}, nil
}
