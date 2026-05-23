package handlers

import (
	"abstractions/server"
	"abstractions/server/responses"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TestHandler struct {
	*BaseHandler
}

func NewTestHandler(server *server.Server) *TestHandler {
	return &TestHandler{
		BaseHandler: NewBaseHandler(server),
	}
}

func (h *TestHandler) CleanupTestData(c echo.Context) error {
	ctx := c.Request().Context()

	deletedCount, err := h.Services(c).Test.CleanupTestUsers(ctx)
	if err != nil {
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusOK, map[string]int{"deletedCount": deletedCount})
}
