package handlers

import (
	"abstractions/server"
	"abstractions/server/responses"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SubscriptionHandler struct {
	*BaseHandler
}

func NewSubscriptionHandler(server *server.Server) *SubscriptionHandler {
	return &SubscriptionHandler{
		BaseHandler: NewBaseHandler(server),
	}
}

func (h *SubscriptionHandler) GetActiveWorkspaceSubscription(c echo.Context) error {

	workspaceID := c.Param("workspace_id")
	if workspaceID == "" {
		return responses.ErrorResponse(c, http.StatusBadRequest, "workspace_id is required")
	}

	wsUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}

	activeSubscription, err := h.Services(c).Subscription.GetActiveSubscriptionByWorkspaceID(wsUUID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, responses.NewSubscriptionResponse(activeSubscription))
}
