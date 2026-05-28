package responses

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/utils/context_utils"
	"abstractions/utils/impersonation_utils"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// ImpersonationInfo represents impersonation status and details
type ImpersonationInfo struct {
	IsImpersonated   bool              `json:"isImpersonated"`
	ImpersonatedBy   *uuid.UUID        `json:"impersonatedBy,omitempty"`
	ImpersonatedAt   *string           `json:"impersonatedAt,omitempty"`
	ImpersonatorUser *ImpersonatorUser `json:"impersonatorUser,omitempty"`
}

// ImpersonatorUser represents basic info about the user doing the impersonation
type ImpersonatorUser struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	Avatar *string   `json:"avatar"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expiresAt"`
}

func NewLoginResponse(token string, expiresAt string) *LoginResponse {
	return &LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}
}

type GenericSuccessResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func NewGenericSuccessResponse(err error) *GenericSuccessResponse {
	resp := GenericSuccessResponse{}
	if err != nil {
		resp.Error = err.Error()
	} else {
		resp.Success = true
	}
	return &resp
}

type CurrentUserResponse struct {
	ID                 uuid.UUID                  `json:"id"`
	Email              string                     `json:"email"`
	Name               string                     `json:"name"`
	Avatar             *string                    `json:"avatar"`
	CreatedAt          time.Time                  `json:"createdAt"`
	UpdatedAt          time.Time                  `json:"updatedAt"`
	DefaultWorkspaceID *uuid.UUID                 `json:"defaultWorkspaceId,omitempty"`
	VerifiedAt         *time.Time                 `json:"verifiedAt,omitempty"`
	Authenticators     []models.AuthenticatorType `json:"authenticators"`
	GodMode            bool                       `json:"godMode"`
	Impersonation      *ImpersonationInfo         `json:"impersonation,omitempty"`
}

func NewCurrentUserResponse(db *gorm.DB, user *models.User, c echo.Context) *CurrentUserResponse {
	authenticatorRepository := repositories.NewAuthenticatorRepository(db)
	configuredAuthenticators, _ := authenticatorRepository.GetAuthenticatorTypesForUser(user.ID)

	response := &CurrentUserResponse{
		ID:                 user.ID,
		Email:              user.Email,
		Name:               user.Name,
		Avatar:             user.Avatar,
		CreatedAt:          user.CreatedAt,
		UpdatedAt:          user.UpdatedAt,
		DefaultWorkspaceID: user.DefaultWorkspaceID,
		VerifiedAt:         user.VerifiedAt,
		Authenticators:     configuredAuthenticators,
		GodMode:            user.GodMode,
	}

	// Add impersonation info only if session is impersonated
	if context_utils.IsSessionImpersonated(c) {
		impersonationInfo, err := impersonation_utils.GetImpersonationInfo(c)
		if err == nil && impersonationInfo.IsImpersonated {
			response.Impersonation = &ImpersonationInfo{
				IsImpersonated: impersonationInfo.IsImpersonated,
				ImpersonatedBy: impersonationInfo.ImpersonatedBy,
				ImpersonatedAt: impersonationInfo.ImpersonatedAt,
			}

			// Add impersonator user info if available
			if impersonationInfo.ImpersonatorUser != nil {
				response.Impersonation.ImpersonatorUser = &ImpersonatorUser{
					ID:     impersonationInfo.ImpersonatorUser.ID,
					Name:   impersonationInfo.ImpersonatorUser.Name,
					Email:  impersonationInfo.ImpersonatorUser.Email,
					Avatar: impersonationInfo.ImpersonatorUser.Avatar,
				}
			}
		}
	}

	return response
}
