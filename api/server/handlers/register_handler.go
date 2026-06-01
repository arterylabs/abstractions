package handlers

import (
	s "abstractions/server"

	"github.com/labstack/echo/v4"
)

type RegisterHandler struct {
	*BaseHandler
}

func NewRegisterHandler(server *s.Server) *RegisterHandler {
	return &RegisterHandler{
		BaseHandler: NewBaseHandler(server),
	}
}

// Register godoc
//
//	@Summary		Register
//	@Description	New user registration
//	@ID				user-register
//	@Tags			User Actions
//	@Accept			json
//	@Produce		json
//	@Param			params	body		requests.RegisterRequest	true	"User's email, user's password"
//	@Success		201		{object}	responses.Data
//	@Failure		400		{object}	responses.Error
//	@Router			/register [post]
func (h *RegisterHandler) Register(c echo.Context) error {
	return nil
}
