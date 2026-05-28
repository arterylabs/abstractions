package middlewares

import (
	"abstractions/logger"
	"abstractions/server/responses"
	authorization_service "abstractions/services/authorization"
	"abstractions/utils/context_utils"
	"abstractions/utils/string_utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

// RequirePermissions creates a middleware that checks if the authenticated user has all the required permissions
func RequirePermissions(permissions ...authorization_service.Permission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			// skip public routes
			if string_utils.StringArrayContainsString(publicRoutes, c.Path()) {
				return next(c)
			}

			// Get the user from context
			user, err := context_utils.GetUserFromContext(c)
			if err != nil {
				return responses.ErrorResponse(c, http.StatusUnauthorized, "user not authenticated")
			}

			// Get workspace ID from path parameter
			workspaceID := c.Param("workspace_id")
			if workspaceID == "" {
				return next(c) // If no workspace ID, skip permission check
			}

			workspaceUUID, err := uuid.Parse(workspaceID)
			if err != nil {
				return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
			}

			// Create authorization service
			db := context_utils.GetTransactionFromContext(c)
			authService := authorization_service.NewAuthorizationService(db)

			// Validate permissions
			if err := authService.Validate(c.Request().Context(), user, workspaceUUID, permissions); err != nil {
				switch err {
				case authorization_service.ErrNoRole:
					return responses.ErrorResponse(c, http.StatusForbidden, "user has no access to this workspace")
				case authorization_service.ErrUnauthorized:
					return responses.ErrorResponse(c, http.StatusForbidden, "user does not have required permissions")
				default:
					logger.Error(c.Request().Context(), err)
					return responses.ErrorResponse(c, http.StatusInternalServerError, "an internal error occurred")
				}
			}

			return next(c)
		}
	}
}
