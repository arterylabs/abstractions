package session_service

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/utils/jwt_utils"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session has expired")
	ErrSessionInvalid  = errors.New("session is invalid")
	ErrInvalidToken    = errors.New("invalid authentication token")
)

type Service struct {
	db          *gorm.DB
	sessionRepo *repositories.SessionRepository
	userRepo    *repositories.UserRepository
}

type TokenValidationResult struct {
	UserID    uuid.UUID
	User      *models.User
	SessionID uuid.UUID
	Session   *models.Session
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db:          db,
		sessionRepo: repositories.NewSessionRepository(db),
		userRepo:    repositories.NewUserRepository(db),
	}
}

// Basic CRUD operations
func (s *Service) Create(ctx context.Context, session *models.Session) error {
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return fmt.Errorf("session_service: create: %w", err)
	}
	return nil
}

func (s *Service) GetByID(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session_service: get_by_id: %w", err)
	}
	if session == nil {
		return nil, ErrSessionNotFound
	}
	return session, nil
}

func (s *Service) GetSessionIfActive(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	session, err := s.sessionRepo.GetActiveByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session_service: get_session_if_active: %w", err)
	}
	if session == nil {
		return nil, ErrSessionNotFound
	}
	return session, nil
}

func (s *Service) IsSessionValid(ctx context.Context, sessionID uuid.UUID) (bool, error) {
	session, err := s.GetByID(ctx, sessionID)
	if err != nil {
		return false, fmt.Errorf("session_service: is_session_valid: %w", err)
	}
	return session.IsValid(), nil
}

func (s *Service) InvalidateSession(ctx context.Context, sessionID uuid.UUID) error {
	_, err := s.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session_service: invalidate_session: get_session: %w", err)
	}
	if err := s.sessionRepo.RevokeByID(ctx, sessionID); err != nil {
		return fmt.Errorf("session_service: invalidate_session: revoke: %w", err)
	}
	return nil
}

func (s *Service) InvalidateAllSessionsForUser(ctx context.Context, userID uuid.UUID) error {
	if err := s.sessionRepo.RevokeAllForUser(ctx, userID); err != nil {
		return fmt.Errorf("session_service: invalidate_all_sessions_for_user: %w", err)
	}
	return nil
}

func (s *Service) InvalidateAllImpersonationSessionsForAdmin(ctx context.Context, userID uuid.UUID) error {
	if err := s.sessionRepo.RevokeActiveImpersonationForAdmin(ctx, userID); err != nil {
		return fmt.Errorf("session_service: invalidate_all_impersonation_sessions_for_admin: %w", err)
	}
	return nil
}

// Complex business logic
func (s *Service) ValidateToken(ctx context.Context, authToken string) (*TokenValidationResult, error) {
	userID, sessionID, err := jwt_utils.DecodeJWT(authToken)
	if err != nil {
		return nil, fmt.Errorf("session_service: validate_token: decode_jwt: %w", err)
	}

	sessionUUID, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session_service: validate_token: invalid session id in token: %w", err)
	}

	session, err := s.GetSessionIfActive(ctx, sessionUUID)
	if err != nil {
		return nil, fmt.Errorf("session_service: validate_token: get_session: %w", err)
	}

	if session.UserID.String() != userID {
		return nil, nil
	}

	if !session.IsValid() {
		return nil, nil
	}

	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("session_service: validate_token: get_user: %w", err)
	}
	// GetByID returns (nil, nil) when the user row is missing (e.g. soft-deleted).
	// Treat that as an invalid token rather than returning a result with a nil User,
	// which would panic callers that dereference validationResult.User.
	if user == nil {
		return nil, nil
	}

	return &TokenValidationResult{
		UserID:    session.UserID,
		User:      user,
		SessionID: session.ID,
		Session:   session,
	}, nil
}
