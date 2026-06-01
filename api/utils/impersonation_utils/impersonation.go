package impersonation_utils

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/utils/context_utils"
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// IsImpersonatedSession checks if the current session is an impersonated session
func IsImpersonatedSession(c echo.Context) bool {
	session, err := context_utils.GetSessionFromContext(c)
	if err != nil {
		return false
	}
	return session.ImpersonatedBy != nil
}

// GetImpersonatorUser returns the User object of who is impersonating, if any
func GetImpersonatorUser(c echo.Context) (*models.User, error) {
	session, err := context_utils.GetSessionFromContext(c)
	if err != nil {
		return nil, err
	}

	if session.ImpersonatedBy == nil {
		return nil, nil
	}

	db := context_utils.GetTransactionFromContext(c)
	userRepo := repositories.NewUserRepository(db)
	impersonator, err := userRepo.GetByID(requestContext(c), *session.ImpersonatedBy)
	if err != nil {
		return nil, err
	}

	return impersonator, nil
}

// GetImpersonationInfo returns detailed impersonation information
type ImpersonationInfo struct {
	IsImpersonated   bool
	ImpersonatedBy   *uuid.UUID
	ImpersonatedAt   *string
	ImpersonatorUser *models.User
}

func GetImpersonationInfo(c echo.Context) (*ImpersonationInfo, error) {
	session, err := context_utils.GetSessionFromContext(c)
	if err != nil {
		return &ImpersonationInfo{IsImpersonated: false}, err
	}

	info := &ImpersonationInfo{
		IsImpersonated: session.ImpersonatedBy != nil,
		ImpersonatedBy: session.ImpersonatedBy,
	}

	if session.ImpersonatedAt != nil {
		impersonatedAtStr := session.ImpersonatedAt.Format("2006-01-02T15:04:05Z07:00")
		info.ImpersonatedAt = &impersonatedAtStr
	}

	if session.ImpersonatedBy != nil {
		db := context_utils.GetTransactionFromContext(c)
		userRepo := repositories.NewUserRepository(db)
		impersonator, err := userRepo.GetByID(requestContext(c), *session.ImpersonatedBy)
		if err != nil {
			return info, err
		}
		info.ImpersonatorUser = impersonator
	}

	return info, nil
}

func requestContext(c echo.Context) context.Context {
	if c.Request() == nil {
		return context.Background()
	}
	return c.Request().Context()
}

// BlockIfImpersonated is a helper function that can be used to block certain operations
// when the session is impersonated. Returns true if the operation should be blocked.
func BlockIfImpersonated(c echo.Context) bool {
	return IsImpersonatedSession(c)
}
