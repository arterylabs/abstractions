package handlers

import (
	"abstractions/server"
	"abstractions/services"
	"abstractions/utils/context_utils"
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
