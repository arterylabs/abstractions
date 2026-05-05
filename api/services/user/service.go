package user_service

import (
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
)

type Service struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

// Basic CRUD operations
func (s *Service) CreateUser(user *models.User) error {
	return s.DB.Create(user).Error
}

func (s *Service) ExistsByEmail(email string) bool {
	userRepository := repositories.NewUserRepository(s.DB)
	return userRepository.ExistsByEmail(email)
}

func (s *Service) GetUserByID(id uuid.UUID) (*models.User, error) {
	userRepository := repositories.NewUserRepository(s.DB)
	return userRepository.GetByID(id)
}

func (s *Service) VerifyUser(tokenID uuid.UUID) error {
	tokenRepository := repositories.NewTokenRepository(s.DB)

	token := &models.Token{}

	token, err := tokenRepository.GetTokenByID(tokenID)
	if err != nil {
		return err
	}

	userRespository := repositories.NewUserRepository(s.DB)

	user, err := userRespository.GetByID(token.UserID)
	if err != nil {
		return err
	}

	currentTime := time.Now()
	user.VerifiedAt = &currentTime

	err = userRespository.Update(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetByEmail(email string) (*models.User, error) {
	userRepository := repositories.NewUserRepository(s.DB)
	return userRepository.GetByEmail(email)
}

func (s *Service) UpdatePassword(userID uuid.UUID, passwordText string) error {
	// hash the user password and update the user
	userRepository := repositories.NewUserRepository(s.DB)
	user, err := userRepository.GetByID(userID)
	if err != nil {
		return err
	}

	// get password authenticator
	authenticatorRepositories := repositories.NewAuthenticatorRepository(s.DB)
	authenticator, err := authenticatorRepositories.GetAuthenticatorByType(user.ID, models.AuthenticatorTypePassword)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	hashedPassword := password_utils.GeneratePasswordHash([]byte(passwordText))

	// user doesn't have a password authenticator so we need to create one
	if authenticator == nil {
		authenticatorService := authenticator_service.NewAuthenticatorService(s.DB)
		err := authenticatorService.CreateAuthenticator(&models.Authenticator{
			UserID: user.ID,
			Type:   models.AuthenticatorTypePassword,
			Token:  hashedPassword,
		})
		return err
	}

	authenticator.Token = hashedPassword

	err = authenticatorRepositories.Update(authenticator)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Paginate(page int, perPage int, query string) ([]models.User, error) {
	userRepository := repositories.NewUserRepository(s.DB)
	return userRepository.Paginate(page, perPage, query)
}

func (s *Service) PaginateWithCount(page int, perPage int, query string) ([]models.User, int, error) {
	userRepository := repositories.NewUserRepository(s.DB)
	return userRepository.PaginateWithCount(page, perPage, query)
}

func (s *Service) PaginateWithCountForAdmin(page int, perPage int, query string) ([]repositories.AdminUserData, int, error) {
	userRepository := repositories.NewUserRepository(s.DB)
	return userRepository.PaginateWithCountForAdmin(page, perPage, query)
}

// GetAdminUserData returns a single user with admin data
func (s *Service) GetAdminUserData(id uuid.UUID) (*repositories.AdminUserData, error) {
	userRepository := repositories.NewUserRepository(s.DB)
	return userRepository.GetAdminUserDataByID(id)
}

type UpdateUserParams struct {
	Email  *string
	Name   *string
	Avatar *string
}

func (s *Service) Update(id uuid.UUID, params UpdateUserParams) (*models.User, error) {
	userRepository := repositories.NewUserRepository(s.DB)

	// Get the existing user first
	user, err := userRepository.GetByID(id)
	if err != nil {
		return nil, err
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

	err = userRepository.Update(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) BulkVerifyUsers(IDs []uuid.UUID) error {
	userRepository := repositories.NewUserRepository(s.DB)

	users, err := userRepository.GetByIDs(IDs)
	if err != nil {
		return err
	}

	for idx := range users {
		currentTime := time.Now()
		users[idx].VerifiedAt = &currentTime
	}

	err = userRepository.BulkUpdate(users)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) BulkDeleteUsers(IDs []uuid.UUID, deletedBy uuid.UUID) error {
	userRepository := repositories.NewUserRepository(s.DB)

	users, err := userRepository.GetByIDs(IDs)
	if err != nil {
		return err
	}

	for idx := range users {
		currentTime := time.Now()
		users[idx].DeletedAt = &currentTime
		users[idx].DeletedBy = &deletedBy
	}

	err = userRepository.BulkUpdate(users)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) BulkRecoverUsers(IDs []uuid.UUID) error {
	userRepository := repositories.NewUserRepository(s.DB)

	users, err := userRepository.GetByIDsIncludingDeleted(IDs)
	if err != nil {
		return err
	}

	for idx := range users {
		users[idx].DeletedAt = nil
		users[idx].DeletedBy = nil
	}

	err = userRepository.BulkUpdate(users)
	if err != nil {
		return err
	}

	return nil
}

type CreateSessionParams struct {
	UserAgent      string
	IP             string
	ImpersonatedBy *uuid.UUID
}

func (s *Service) CreateUserSession(userID uuid.UUID, params *CreateSessionParams) (*models.Session, error) {
	expires_at := time.Now().UTC().Add(time.Hour * 24 * 7 * time.Duration(6))

	sessionID := uuid.New()
	jwtToken, err := jwt_utils.GenerateJWT(userID.String(), sessionID.String(), expires_at)
	if err != nil {
		return nil, err
	}

	session := models.Session{
		ID:        sessionID,
		UserID:    userID,
		Token:     jwtToken,
		ExpiresAt: expires_at,
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

	sessionRepository := repositories.NewSessionRepository(s.DB)
	err = sessionRepository.Create(&session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

type CreateVerifiedUserParams struct {
	Email string
	Name  string
}

func (s *Service) CreateVerifiedUser(params CreateVerifiedUserParams) (*models.User, error) {
	userRepository := repositories.NewUserRepository(s.DB)

	user := models.User{
		Email: params.Email,
		Name:  params.Name,
	}

	currentTime := time.Now()
	user.VerifiedAt = &currentTime

	err := userRepository.Create(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) Count() (int, error) {
	userRepository := repositories.NewUserRepository(s.DB)
	return userRepository.Count()
}

func (s *Service) CountActiveUsers() (int, error) {
	userRepository := repositories.NewUserRepository(s.DB)
	return userRepository.CountActiveUsers()
}

// GetHistoricalCreationCounts returns daily user creation counts for the specified number of days
func (s *Service) GetHistoricalCreationCounts(days int) (map[string]int, error) {
	userRepository := repositories.NewUserRepository(s.DB)
	return userRepository.GetHistoricalCreationCounts(days)
}

func (s *Service) SendVerificationEmail(user *models.User) error {
	tokenService := token_service.NewTokenService(s.DB)
	tokenValue := uuid.New()
	tokenString := strings.ReplaceAll(tokenValue.String(), "-", "")

	token := models.Token{
		UserID:    user.ID,
		Type:      models.TokenTypeEmailVerification,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	err := tokenService.CreateToken(&token)
	if err != nil {
		return err
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

func (s *Service) SendResetPasswordEmail(userEmail string) error {
	tokenService := token_service.NewTokenService(s.DB)

	user, err := s.GetByEmail(userEmail)
	if err != nil {
		return err
	}

	tokenValue := uuid.New()
	tokenString := strings.ReplaceAll(tokenValue.String(), "-", "")

	token := models.Token{
		UserID:    user.ID,
		Type:      models.TokenTypePasswordReset,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}

	err = tokenService.CreateToken(&token)
	if err != nil {
		return err
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

func (s *Service) DeleteUser(id uuid.UUID) error {
	userRepository := repositories.NewUserRepository(s.DB)
	return userRepository.Delete(id)
}
