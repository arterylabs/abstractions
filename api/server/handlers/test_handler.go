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

func (h *TestHandler) VerifyTestUser(c echo.Context) error {
	ctx := c.Request().Context()

	var body struct {
		Email string `json:"email"`
	}
	if err := c.Bind(&body); err != nil || body.Email == "" {
		return responses.ErrorResponse(c, http.StatusBadRequest, "email is required")
	}

	if err := h.Services(c).Test.VerifyTestUser(ctx, body.Email); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return responses.Response(c, http.StatusOK, map[string]bool{"verified": true})
}
