package handlers

import (
	"abstractions/clients/storage"
	"abstractions/db/models"
	"abstractions/logger"
	"abstractions/server"
	"abstractions/server/requests"
	"abstractions/server/responses"
	user_service "abstractions/services/user"
	"abstractions/utils/context_utils"
	"abstractions/utils/impersonation_utils"
	"abstractions/utils/password_utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

type UsersHandler struct {
	*BaseHandler
}

func NewUsersHandler(server *server.Server) *UsersHandler {
	return &UsersHandler{
		BaseHandler: NewBaseHandler(server),
	}
}

func (h *UsersHandler) GetCurrentUser(c echo.Context) error {
	db := context_utils.GetTransactionFromContext(c)

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, responses.NewCurrentUserResponse(db, user, c))
}

func (h *UsersHandler) UpdateCurrentUser(c echo.Context) error {
	updateUserRequest := new(requests.UpdateUserRequest)
	if err := c.Bind(updateUserRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	db := context_utils.GetTransactionFromContext(c)

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	user, err = h.Services(c).User.Update(user.ID, user_service.UpdateUserParams{
		Name: &updateUserRequest.Name,
	})
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, responses.NewCurrentUserResponse(db, user, c))
}

func (h *UsersHandler) DeleteCurrentUser(c echo.Context) error {
	// Block account deletion during impersonation for security
	if impersonation_utils.IsImpersonatedSession(c) {
		return responses.ErrorResponse(c, http.StatusForbidden, "Account deletion is not allowed during impersonation")
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	err = h.Services(c).User.DeleteUser(user.ID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, nil)
}

func (h *UsersHandler) UpdateUserAvatar(c echo.Context) error {
	db := context_utils.GetTransactionFromContext(c)

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	file, err := c.FormFile("file")
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "file is required")
	}
	src, err := file.Open()
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "file is required")
	}
	defer src.Close()

	storageClient, err := storage.NewMinioStorageClient()
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	filePath := "avatars/" + user.ID.String()

	fileLink, err := storageClient.UploadPublic(c.Request().Context(), filePath, src)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	user, err = h.Services(c).User.Update(user.ID, user_service.UpdateUserParams{
		Avatar: &fileLink,
	})
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, responses.NewCurrentUserResponse(db, user, c))
}

func (h *UsersHandler) UpdateUserPassword(c echo.Context) error {
	// Block password updates during impersonation for security
	if impersonation_utils.IsImpersonatedSession(c) {
		return responses.ErrorResponse(c, http.StatusForbidden, "Password updates are not allowed during impersonation")
	}

	updateUserAuthenticatorRequest := new(requests.UpdateUserPasswordRequest)
	if err := c.Bind(updateUserAuthenticatorRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	// Get the password authenticator if it exists
	authenticator, err := h.Services(c).Authenticator.GetAuthenticatorByType(user.ID, models.AuthenticatorTypePassword)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error(c.Request().Context(), "error getting authenticator", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	// If the user has a password authenticator, verify the current password
	if authenticator != nil {
		if updateUserAuthenticatorRequest.CurrentPassword == "" {
			return responses.ErrorResponse(c, http.StatusBadRequest, "Current password is required")
		}
		if !password_utils.ComparePasswords(authenticator.Token, []byte(updateUserAuthenticatorRequest.CurrentPassword)) {
			return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid password")
		}
	}

	err = h.Services(c).User.UpdatePassword(user.ID, updateUserAuthenticatorRequest.NewPassword)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	err = h.Services(c).Session.InvalidateAllSessionsForUser(user.ID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, nil)
}

func (h *UsersHandler) SwitchWorkspace(c echo.Context) error {
	switchWorkspaceRequest := new(requests.SwitchWorkspaceRequest)
	if err := c.Bind(switchWorkspaceRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := switchWorkspaceRequest.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	db := context_utils.GetTransactionFromContext(c)

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	err = h.Services(c).User.UpdateDefaultWorkspace(c.Request().Context(), user.ID, switchWorkspaceRequest.WorkspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	user.DefaultWorkspaceID = &switchWorkspaceRequest.WorkspaceID

	return responses.Response(c, http.StatusOK, responses.NewCurrentUserResponse(db, user, c))
}
