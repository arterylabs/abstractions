package handlers

import (
	"abstractions/clients/storage"
	"abstractions/db/models"
	"abstractions/logger"
	"abstractions/server"
	"abstractions/server/requests"
	"abstractions/server/responses"
	authenticator_service "abstractions/services/authenticator"
	user_service "abstractions/services/user"
	"abstractions/utils/context_utils"
	"abstractions/utils/impersonation_utils"
	"abstractions/utils/password_utils"
	"errors"
	"github.com/labstack/echo/v4"
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
	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return h.internalError(c, err)
	}

	authenticatorTypes, err := h.Services(c).Authenticator.ListTypesForUser(c.Request().Context(), user.ID)
	if err != nil {
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusOK, responses.NewCurrentUserResponse(user, authenticatorTypes, c))
}

func (h *UsersHandler) UpdateCurrentUser(c echo.Context) error {
	updateUserRequest := new(requests.UpdateUserRequest)
	if err := c.Bind(updateUserRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return h.internalError(c, err)
	}

	user, err = h.Services(c).User.Update(c.Request().Context(), user.ID, user_service.UpdateUserParams{
		Name: &updateUserRequest.Name,
	})
	if err != nil {
		return h.internalError(c, err)
	}

	authenticatorTypes, err := h.Services(c).Authenticator.ListTypesForUser(c.Request().Context(), user.ID)
	if err != nil {
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusOK, responses.NewCurrentUserResponse(user, authenticatorTypes, c))
}

func (h *UsersHandler) DeleteCurrentUser(c echo.Context) error {
	// Block account deletion during impersonation for security
	if impersonation_utils.IsImpersonatedSession(c) {
		return responses.ErrorResponse(c, http.StatusForbidden, "Account deletion is not allowed during impersonation")
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return h.internalError(c, err)
	}

	err = h.Services(c).User.Delete(c.Request().Context(), user.ID, &user.ID)
	if err != nil {
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusOK, nil)
}

func (h *UsersHandler) UpdateUserAvatar(c echo.Context) error {
	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return h.internalError(c, err)
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
		return h.internalError(c, err)
	}

	filePath := "avatars/" + user.ID.String()

	fileLink, err := storageClient.UploadPublic(c.Request().Context(), filePath, src)
	if err != nil {
		if errors.Is(err, storage.ErrFileTooLarge) || errors.Is(err, storage.ErrMIMENotAllowed) {
			return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		return responses.ErrorResponse(c, http.StatusInternalServerError, "upload failed")
	}

	user, err = h.Services(c).User.Update(c.Request().Context(), user.ID, user_service.UpdateUserParams{
		Avatar: &fileLink,
	})
	if err != nil {
		return h.internalError(c, err)
	}

	authenticatorTypes, err := h.Services(c).Authenticator.ListTypesForUser(c.Request().Context(), user.ID)
	if err != nil {
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusOK, responses.NewCurrentUserResponse(user, authenticatorTypes, c))
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
		return h.internalError(c, err)
	}

	// Get the password authenticator if it exists
	authenticator, err := h.Services(c).Authenticator.GetAuthenticatorByType(c.Request().Context(), user.ID, models.AuthenticatorTypePassword)
	if err != nil && !errors.Is(err, authenticator_service.ErrAuthenticatorNotFound) {
		logger.Error(c.Request().Context(), "error getting authenticator", err)
		return h.internalError(c, err)
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

	err = h.Services(c).User.UpdatePassword(c.Request().Context(), user.ID, updateUserAuthenticatorRequest.NewPassword)
	if err != nil {
		return h.internalError(c, err)
	}

	err = h.Services(c).Session.InvalidateAllSessionsForUser(c.Request().Context(), user.ID)
	if err != nil {
		return h.internalError(c, err)
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

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return h.internalError(c, err)
	}

	err = h.Services(c).User.UpdateDefaultWorkspace(c.Request().Context(), user.ID, switchWorkspaceRequest.WorkspaceID)
	if err != nil {
		if errors.Is(err, user_service.ErrWorkspaceAccessDenied) {
			return responses.ErrorResponse(c, http.StatusForbidden, "user does not have access to this workspace")
		}
		return h.internalError(c, err)
	}

	user.DefaultWorkspaceID = &switchWorkspaceRequest.WorkspaceID

	authenticatorTypes, err := h.Services(c).Authenticator.ListTypesForUser(c.Request().Context(), user.ID)
	if err != nil {
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusOK, responses.NewCurrentUserResponse(user, authenticatorTypes, c))
}
