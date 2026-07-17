package handlers

import (
	"abstractions/logger"
	"abstractions/server"
	"abstractions/server/responses"
	"abstractions/services"
	"abstractions/utils/context_utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// BaseHandler provides common functionality for all handlers
type BaseHandler struct {
	server *server.Server
}

// NewBaseHandler creates a new base handler
func NewBaseHandler(server *server.Server) *BaseHandler {
	return &BaseHandler{server: server}
}

// GetServices retrieves the pre-initialized Services instance from the request context
func (h *BaseHandler) Services(c echo.Context) *services.Services {
	svcs := context_utils.GetServicesFromContext(c)
	if svcs == nil {
		// This should not happen if services middleware is properly configured
		panic("services not found in context - ensure services middleware is applied")
	}

	services, ok := svcs.(*services.Services)
	if !ok {
		panic("invalid services type in context")
	}

	return services
}

// GetServer returns the server instance
func (h *BaseHandler) Server() *server.Server {
	return h.server
}

// internalError logs the original error and returns a generic 500 response,
// ensuring internal details (DB schema, SQL state codes, etc.) are never exposed to clients.
func (h *BaseHandler) internalError(c echo.Context, err error) error {
	logger.Error(c.Request().Context(), err)
	return responses.ErrorResponse(c, http.StatusInternalServerError, "an internal error occurred")
}

func (h *BaseHandler) uuidParam(c echo.Context, name string) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param(name))
	if err != nil {
		return uuid.Nil, responses.ErrorResponse(c, http.StatusBadRequest, "invalid "+name+" format")
	}
	return id, nil
}
