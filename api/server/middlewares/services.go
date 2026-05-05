package middlewares

import (
	"abstractions/services"
	"abstractions/utils/context_utils"
	"github.com/labstack/echo/v4"
)

// InitializeServices middleware initializes all services once per request
// and stores them in the context for reuse throughout the request lifecycle.
// This middleware must be applied after the transaction middleware.
func InitializeServices() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the database transaction from context (set by transaction middleware)
			db := context_utils.GetTransactionFromContext(c)
			if db == nil {
				// This should not happen if transaction middleware is properly configured
				return echo.NewHTTPError(500, "database transaction not found in context")
			}

			// Initialize services once per request
			svcs := services.New(db)

			// Store services in context
			c = context_utils.SetServicesInContext(c, svcs)

			// Continue with the request
			return next(c)
		}
	}
}
