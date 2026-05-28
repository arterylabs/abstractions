package middlewares

import (
	app_notifications "abstractions/clients/app_notifications"
	"abstractions/logger"
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo/v4"
)

// PanicRecoveryMiddleware catches any panic in a handler, logs it with a stack trace,
// sends a Slack alert, and returns a 500 to the client so the server keeps running.
func PanicRecoveryMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				msg := fmt.Sprintf("%v", r)

				logger.Error(c.Request().Context(), fmt.Errorf("panic recovered: %s\n%s", msg, stack))

				go app_notifications.SendPanicNotification(context.Background(), app_notifications.PanicParams{
					Source:  fmt.Sprintf("%s %s", c.Request().Method, c.Request().URL.Path),
					Message: msg,
					Stack:   stack,
				})

				err = echo.NewHTTPError(http.StatusInternalServerError, "an internal error occurred")
			}
		}()
		return next(c)
	}
}
