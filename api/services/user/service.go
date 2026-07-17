package user_service

import (
	"abstractions/clients/app_notifications"
	"abstractions/clients/email"
	"abstractions/config"
	"abstractions/db/models"
	"abstractions/repositories"
	authenticator_service "abstractions/services/authenticator"
	token_service "abstractions/services/token"
	"abstractions/utils/jwt_utils"
	"abstractions/utils/password_utils"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrUserAlreadyVerified = errors.New("user is already verified")
	ErrInvalidEmail        = errors.New("invalid email address")
	ErrInvalidPassword     = errors.New("invalid password")

	ErrWorkspaceAccessDenied = errors.New("user is not a member of the workspace")
)

type Service struct {
	db                *gorm.DB
	userRepo          *repositories.UserRepository
	tokenRepo         *repositories.TokenRepository
	sessionRepo       *repositories.SessionRepository
	authenticatorRepo *repositories.AuthenticatorRepository
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db:                db,
		userRepo:          repositories.NewUserRepository(db),
		tokenRepo:         repositories.NewTokenRepository(db),
		sessionRepo:       repositories.NewSessionRepository(db),
		authenticatorRepo: repositories.NewAuthenticatorRepository(db),
	}
}

// Basic CRUD operations
func (s *Service) Create(ctx context.Context, user *models.User) error {
	if err := s.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("user_service: create: %w", err)
	}
	go app_notifications.SendNewSignupNotification(context.Background(), app_notifications.NewSignupParams{
		Name:  user.Name,
		Email: user.Email,
	})
	return nil
}

func (s *Service) ExistsByEmail(ctx context.Context, email string) bool {
	return s.userRepo.ExistsByEmail(ctx, email)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *Service) VerifyUser(ctx context.Context, tokenID uuid.UUID) error {
	token, err := s.tokenRepo.GetByID(ctx, tokenID)
	if err != nil {
		return fmt.Errorf("user_service: verify_user: get_token: %w", err)
	}
	if token == nil {
		return fmt.Errorf("user_service: verify_user: %w", ErrUserNotFound)
	}

	user, err := s.userRepo.GetByID(ctx, token.UserID)
	if err != nil {
		return fmt.Errorf("user_service: verify_user: get_user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user_service: verify_user: %w", ErrUserNotFound)
	}

	currentTime := time.Now()
	user.VerifiedAt = &currentTime

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("user_service: verify_user: update_user: %w", err)
	}

	return nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *Service) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordText string) error {
	// hash the user password and update the user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user_service: update_password: get_user: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	// get password authenticator
	authenticator, err := s.authenticatorRepo.GetByType(ctx, user.ID, models.AuthenticatorTypePassword)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("user_service: update_password: get_authenticator: %w", err)
	}

	hashedPassword := password_utils.GeneratePasswordHash([]byte(passwordText))

	// user doesn't have a password authenticator so we need to create one
	if authenticator == nil {
		authenticatorService := authenticator_service.NewService(s.db)
		err := authenticatorService.Create(ctx, &models.Authenticator{
			UserID: user.ID,
			Type:   models.AuthenticatorTypePassword,
			Token:  hashedPassword,
		})
		if err != nil {
			return fmt.Errorf("user_service: update_password: create_authenticator: %w", err)
		}
		return nil
	}

	authenticator.Token = hashedPassword

	err = s.authenticatorRepo.Update(ctx, authenticator)
	if err != nil {
		return fmt.Errorf("user_service: update_password: update_authenticator: %w", err)
	}

	return nil
}

func (s *Service) Paginate(ctx context.Context, page int, perPage int, query string) ([]models.User, error) {
	return s.userRepo.Paginate(ctx, page, perPage, query)
}

func (s *Service) PaginateWithCount(ctx context.Context, page int, perPage int, query string) ([]models.User, int, error) {
	return s.userRepo.PaginateWithCount(ctx, page, perPage, query)
}

func (s *Service) PaginateWithCountForAdmin(ctx context.Context, page int, perPage int, query string, includeDeleted bool) ([]repositories.AdminUserData, int, error) {
	return s.userRepo.PaginateWithCountForAdmin(ctx, page, perPage, query, includeDeleted)
}

// GetAdminUserData returns a single user with admin data
func (s *Service) GetAdminUserData(ctx context.Context, id uuid.UUID) (*repositories.AdminUserData, error) {
	return s.userRepo.GetAdminUserDataByID(ctx, id)
}

type UpdateUserParams struct {
	Email  *string
	Name   *string
	Avatar *string
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, params UpdateUserParams) (*models.User, error) {
	// Get the existing user first
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user_service: update: get_user: %w", err)
	}

	if params.Email != nil {
		user.Email = *params.Email
	}

	if params.Name != nil {
		user.Name = *params.Name
	}

	if params.Avatar != nil {
		user.Avatar = params.Avatar
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("user_service: update: update_user: %w", err)
	}
	return user, nil
}

func (s *Service) BulkVerifyUsers(ctx context.Context, ids []uuid.UUID) error {
	users, err := s.userRepo.GetByIDs(ctx, ids)
	if err != nil {
		return fmt.Errorf("user_service: bulk_verify_users: get_users: %w", err)
	}

	for idx := range users {
		currentTime := time.Now()
		users[idx].VerifiedAt = &currentTime
	}

	err = s.userRepo.BulkUpdate(ctx, users)
	if err != nil {
		return fmt.Errorf("user_service: bulk_verify_users: bulk_update: %w", err)
	}

	return nil
}

func (s *Service) BulkDeleteUsers(ctx context.Context, ids []uuid.UUID, deletedBy uuid.UUID) error {
	users, err := s.userRepo.GetByIDs(ctx, ids)
	if err != nil {
		return fmt.Errorf("user_service: bulk_delete_users: get_users: %w", err)
	}

	for idx := range users {
		currentTime := time.Now()
		users[idx].DeletedAt = &currentTime
		users[idx].DeletedBy = &deletedBy
	}

	err = s.userRepo.BulkUpdate(ctx, users)
	if err != nil {
		return fmt.Errorf("user_service: bulk_delete_users: bulk_update: %w", err)
	}

	return nil
}

func (s *Service) BulkRecoverUsers(ctx context.Context, ids []uuid.UUID) error {
	users, err := s.userRepo.GetByIDsIncludingDeleted(ctx, ids)
	if err != nil {
		return fmt.Errorf("user_service: bulk_recover_users: get_users: %w", err)
	}

	for idx := range users {
		users[idx].DeletedAt = nil
		users[idx].DeletedBy = nil
	}

	err = s.userRepo.BulkUpdate(ctx, users)
	if err != nil {
		return fmt.Errorf("user_service: bulk_recover_users: bulk_update: %w", err)
	}

	return nil
}

type CreateSessionParams struct {
	UserAgent      string
	IP             string
	ImpersonatedBy *uuid.UUID
}

func (s *Service) CreateUserSession(ctx context.Context, userID uuid.UUID, params *CreateSessionParams) (*models.Session, error) {
	expiresAt := time.Now().UTC().Add(time.Hour * 24 * 7 * time.Duration(6))

	sessionID := uuid.New()
	jwtToken, err := jwt_utils.GenerateJWT(userID.String(), sessionID.String(), expiresAt)
	if err != nil {
		return nil, fmt.Errorf("user_service: create_user_session: generate_jwt: %w", err)
	}

	session := models.Session{
		ID:        sessionID,
		UserID:    userID,
		Token:     jwtToken,
		ExpiresAt: expiresAt,
	}
	if params != nil {
		session.UserAgent = params.UserAgent
		session.IP = params.IP
		if params.ImpersonatedBy != nil {
			now := time.Now()
			session.ImpersonatedAt = &now
			session.ImpersonatedBy = params.ImpersonatedBy
		}
	}

	err = s.sessionRepo.Create(ctx, &session)
	if err != nil {
		return nil, fmt.Errorf("user_service: create_user_session: create_session: %w", err)
	}
	return &session, nil
}

type CreateVerifiedUserParams struct {
	Email string
	Name  string
}

func (s *Service) CreateVerifiedUser(ctx context.Context, params CreateVerifiedUserParams) (*models.User, error) {
	user := models.User{
		Email: params.Email,
		Name:  params.Name,
	}

	currentTime := time.Now()
	user.VerifiedAt = &currentTime

	err := s.userRepo.Create(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("user_service: create_verified_user: create_user: %w", err)
	}

	return &user, nil
}

func (s *Service) Count(ctx context.Context) (int, error) {
	return s.userRepo.Count(ctx)
}

func (s *Service) CountActiveUsers(ctx context.Context) (int, error) {
	return s.userRepo.CountActiveUsers(ctx)
}

// GetHistoricalCreationCounts returns daily user creation counts for the specified number of days
func (s *Service) GetHistoricalCreationCounts(ctx context.Context, days int) (map[string]int, error) {
	return s.userRepo.GetHistoricalCreationCounts(ctx, days)
}

func (s *Service) SendVerificationEmail(ctx context.Context, user *models.User) error {
	tokenService := token_service.NewService(s.db)
	tokenValue := uuid.New()
	tokenString := strings.ReplaceAll(tokenValue.String(), "-", "")

	token := models.Token{
		UserID:    user.ID,
		Type:      models.TokenTypeEmailVerification,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	err := tokenService.Create(ctx, &token)
	if err != nil {
		return fmt.Errorf("user_service: send_verification_email: create_token: %w", err)
	}

	go email.SendGenericEmail(context.Background(), email.SendGenericEmailParams{
		To:            user.Email,
		Subject:       "Verify Your Email",
		BodyHTML:      "<p>Click on the link below to verify your email</p>",
		BodyPlainText: "Click on the link below to verify your email",
		CTAText:       "Verify Email",
		CTALink:       config.Environment.ClientURL + "/verify?token=" + token.Token,
		Reason:        "You are receiving this email because you signed up on Abstractions",
	})
	return nil
}

func (s *Service) SendResetPasswordEmail(ctx context.Context, userEmail string) error {
	tokenService := token_service.NewService(s.db)

	user, err := s.GetByEmail(ctx, userEmail)
	if err != nil {
		return fmt.Errorf("user_service: send_reset_password_email: get_user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user_service: send_reset_password_email: %w", ErrUserNotFound)
	}

	tokenValue := uuid.New()
	tokenString := strings.ReplaceAll(tokenValue.String(), "-", "")

	token := models.Token{
		UserID:    user.ID,
		Type:      models.TokenTypePasswordReset,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}

	err = tokenService.Create(ctx, &token)
	if err != nil {
		return fmt.Errorf("user_service: send_reset_password_email: create_token: %w", err)
	}

	go email.SendGenericEmail(context.Background(), email.SendGenericEmailParams{
		To:            user.Email,
		Subject:       "Reset Your Password",
		BodyHTML:      "<p>Click on the link below to reset your password</p>",
		BodyPlainText: "Click on the link below to reset your password",
		CTAText:       "Reset Password",
		CTALink:       config.Environment.ClientURL + "/reset?token=" + token.Token + "&email=" + user.Email + "&id=" + token.ID.String(),
		Reason:        "You are receiving this email because you requested to reset your password on Abstractions",
	})

	return nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID, deletedBy *uuid.UUID) error {
	return s.userRepo.Delete(ctx, id, deletedBy)
}

func (s *Service) HardDelete(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.PermanentlyDelete(ctx, id)
}
