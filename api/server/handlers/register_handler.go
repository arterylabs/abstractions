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
	// registerRequest := new(requests.RegisterRequest)

	// if err := c.Bind(registerRequest); err != nil {
	// 	return err
	// }

	// if err := registerRequest.Validate(); err != nil {
	// 	return responses.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty or not valid")
	// }

	// existUser := models.User{}
	// userRepository := repositories.NewUserRepository(registerHandler.server.DB)
	// userRepository.GetUserByEmail(&existUser, registerRequest.Email)

	// if existUser.ID != 0 {
	// 	return responses.ErrorResponse(c, http.StatusBadRequest, "User already exists")
	// }

	// userService := user_service.NewService(registerHandler.server.DB)
	// if err := userService.Register(registerRequest); err != nil {
	// 	return responses.ErrorResponse(c, http.StatusInternalServerError, "Server error")
	// }

	// return responses.MessageResponse(c, http.StatusCreated, "User successfully created")

	return nil
}
