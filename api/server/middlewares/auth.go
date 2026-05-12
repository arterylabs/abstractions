package middlewares

import (
	"abstractions/logger"
	"abstractions/server/responses"
	session_service "abstractions/services/session"
	"abstractions/utils/context_utils"
	"abstractions/utils/string_utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

var PUBLIC_ROUTES = []string{"/auth/google", "/auth/google/callback", "/auth/github", "/auth/github/callback", "/auth/signup", "/auth/login", "/users/verify", "/users/password/reset/confirm", "/email", "/users/password/reset", "/integrations/callback/:type", "/unsubscribe"}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		// skip public routes
		fmt.Println(c.Path())
		if string_utils.StringArrayContainsString(PUBLIC_ROUTES, c.Path()) {
			return next(c)
		}

		// read the Authorization header and get the bearer token
		// check if the token is valid and if the user is authenticated
		// if the user is authenticated, set the user in the context
		// if the user is not authenticated, return an error

		db := context_utils.GetTransactionFromContext(c)

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(401, "Unauthorized")
		}

		// extract the token value from the Authorization header
		_, token, _ := strings.Cut(authHeader, " ")

		// check if the token is valid
		// if the token is valid, set the user in the context
		// if the token is not valid, return an error

		session_service := session_service.NewSessionService(db)
		validationResult, err := session_service.ValidateToken(token)
		if err != nil || validationResult == nil {
			logger.Error(c.Request().Context(), err)
			return echo.NewHTTPError(401, "Unauthorized")
		}

		c = context_utils.SetUserInContext(c, *validationResult.User)
		c = context_utils.SetSessionInContext(c, *validationResult.Session)

		// Add impersonation information to request logs
		logMetadata := make(map[string]interface{})
		logMetadata["user_id"] = validationResult.User.ID.String()
		logMetadata["session_id"] = validationResult.Session.ID.String()

		if validationResult.Session.ImpersonatedBy != nil {
			logMetadata["is_impersonated"] = true
			logMetadata["impersonated_by"] = validationResult.Session.ImpersonatedBy.String()
			logMetadata["impersonated_at"] = validationResult.Session.ImpersonatedAt
		} else {
			logMetadata["is_impersonated"] = false
		}

		// Set the metadata in the request context for logging
		ctx := logger.SetMetadata(c.Request().Context(), logMetadata)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

func AdminAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		// if the path is not an admin route, skip this middleware
		if !strings.HasPrefix(c.Path(), "/admin") {
			return next(c)
		}

		// get the user from the context as we assume this middleware will be applied after the auth middleware
		// so the user object will already be in the context

		user, err := context_utils.GetUserFromContext(c)
		if err != nil {
			logger.Error(c.Request().Context(), err)
			return responses.ErrorResponse(c, http.StatusForbidden, "forbidden")
		}

		// check if the user is an admin

		if user.GodMode == true {
			return next(c)
		}

		return responses.ErrorResponse(c, http.StatusForbidden, "forbidden")
	}
}
