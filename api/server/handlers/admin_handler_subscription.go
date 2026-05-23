package handlers

import (
	"abstractions/db/models"
	"abstractions/server/requests"
	"abstractions/server/responses"
	subscription_service "abstractions/services/subscription"
	"abstractions/utils/context_utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *AdminHandler) GetWorkspaceSubscription(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	if workspaceID == "" {
		return responses.ErrorResponse(c, http.StatusBadRequest, "workspace_id is required")
	}

	wsUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}

	db := context_utils.GetTransactionFromContext(c)
	subscriptionService := subscription_service.NewSubscriptionService(db)

	subscription, err := subscriptionService.GetActiveSubscriptionByWorkspaceID(wsUUID)
	if err != nil {
		return h.internalError(c, err)
	}

	if subscription == nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "no active subscription found for this workspace")
	}

	return responses.Response(c, http.StatusOK, responses.NewAdminSubscriptionsResponse(&[]models.Subscription{*subscription})[0])
}

func (h *AdminHandler) CreateWorkspaceSubscription(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	if workspaceID == "" {
		return responses.ErrorResponse(c, http.StatusBadRequest, "workspace_id is required")
	}

	wsUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}

	req := new(requests.AdminCreateSubscription)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	db := context_utils.GetTransactionFromContext(c)
	subscriptionService := subscription_service.NewSubscriptionService(db)
	subs, err := subscriptionService.CreateBulkSubscription([]uuid.UUID{wsUUID}, subscription_service.CreateSubscriptionParams{
		PlanID:    req.PlanID,
		Status:    req.Status,
		ExpiresAt: req.ExpiresAt,
		Limit:     req.Limit,
	})
	if err != nil {
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusOK, responses.NewAdminSubscriptionsResponse(subs)[0])
}

func (h *AdminHandler) UpdateWorkspaceSubscription(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	if workspaceID == "" {
		return responses.ErrorResponse(c, http.StatusBadRequest, "workspace_id is required")
	}

	wsUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}

	req := new(requests.AdminUpdateSubscriptions)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	db := context_utils.GetTransactionFromContext(c)
	subscriptionService := subscription_service.NewSubscriptionService(db)
	err = subscriptionService.UpdateSubscriptionsForWorkspaces([]uuid.UUID{wsUUID}, subscription_service.UpdateSubscriptionParams{
		PlanID:    req.PlanID,
		Status:    req.Status,
		ExpiresAt: req.ExpiresAt,
		Limit:     req.Limit,
	})
	if err != nil {
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusOK, "Subscription updated successfully")
}

func (h *AdminHandler) DisableWorkspaceSubscription(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	if workspaceID == "" {
		return responses.ErrorResponse(c, http.StatusBadRequest, "workspace_id is required")
	}

	wsUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}

	req := new(requests.AdminDisableSubscriptions)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	db := context_utils.GetTransactionFromContext(c)
	subscriptionService := subscription_service.NewSubscriptionService(db)
	err = subscriptionService.DisableSubscriptions([]uuid.UUID{wsUUID}, req.DisabledReason)
	if err != nil {
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusOK, "Subscription disabled successfully")
}
