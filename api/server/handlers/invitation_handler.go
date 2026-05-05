package handlers

import (
	s "abstractions/server"
	"abstractions/server/requests"
	"abstractions/server/responses"
	invitation_service "abstractions/services/invitation"
	"abstractions/utils/context_utils"
	"abstractions/utils/pagination_utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type InvitationHandler struct {
	*BaseHandler
}

func NewInvitationHandler(server *s.Server) *InvitationHandler {
	return &InvitationHandler{
		BaseHandler: NewBaseHandler(server),
	}
}

func (h *InvitationHandler) SendInvitations(c echo.Context) error {
	var req requests.SendInvitationsRequest
	if err := c.Bind(&req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid request format")
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	// Validate each invitation input
	for _, invite := range req.Invitations {
		if err := invite.Validate(); err != nil {
			return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
		}
	}

	workspaceID, err := uuid.Parse(c.Param("workspace_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid workspace ID format")
	}

	// Convert request invitations to service inputs
	invitationInputs := make([]invitation_service.InvitationInput, len(req.Invitations))
	for i, inv := range req.Invitations {
		invitationInputs[i] = invitation_service.InvitationInput{
			Email: inv.Email,
			Role:  inv.Role,
		}
	}

	errors := h.Services(c).Invitation.SendInvitations(c.Request().Context(), workspaceID, user.ID, invitationInputs)

	response := responses.NewSendInvitationsResponse(req.Invitations, errors)
	return responses.Response(c, http.StatusOK, response)
}

func (h *InvitationHandler) GetPendingInvitations(c echo.Context) error {
	email := c.Get("user_email").(string)
	invitations, err := h.Services(c).Invitation.GetPendingInvitations(c.Request().Context(), email)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to get pending invitations")
	}

	return responses.Response(c, http.StatusOK, responses.NewInvitationListResponse(invitations))
}

func (h *InvitationHandler) AcceptInvitation(c echo.Context) error {
	invitationID, err := uuid.Parse(c.Param("invitation_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid invitation ID format")
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	err = h.Services(c).Invitation.AcceptInvitation(c.Request().Context(), invitationID, user.ID)
	if err != nil {
		if err == invitation_service.ErrInvalidInvitation {
			return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid invitation")
		}
		if err == invitation_service.ErrInvitationExpired {
			return responses.ErrorResponse(c, http.StatusBadRequest, "Invitation has expired")
		}
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to accept invitation")
	}

	return responses.MessageResponse(c, http.StatusOK, "Invitation accepted successfully")
}

func (h *InvitationHandler) IgnoreInvitation(c echo.Context) error {
	invitationID, err := uuid.Parse(c.Param("invitation_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid invitation ID format")
	}

	err = h.Services(c).Invitation.IgnoreInvitation(c.Request().Context(), invitationID)
	if err != nil {
		if err == invitation_service.ErrInvalidInvitation {
			return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid invitation")
		}
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to ignore invitation")
	}

	return responses.MessageResponse(c, http.StatusOK, "Invitation ignored successfully")
}

// GetWorkspaceInvitations returns all invitations for a workspace with pagination
func (h *InvitationHandler) GetWorkspaceInvitations(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid workspace ID format")
	}

	paginationParams := pagination_utils.PaginationParams{}
	if err := c.Bind(&paginationParams); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	paginationParams.Validate()

	invitations, total, err := h.Services(c).Invitation.GetWorkspaceInvitationsPaginated(c.Request().Context(), workspaceUUID, paginationParams)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to get workspace invitations")
	}

	return responses.Response(c, http.StatusOK, responses.NewPaginatedInvitationsResponse(invitations, paginationParams.Page, paginationParams.PerPage, total))
}

// GetUserPendingInvitations returns all pending invitations for the current user with pagination
func (h *InvitationHandler) GetUserPendingInvitations(c echo.Context) error {
	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	paginationParams := pagination_utils.PaginationParams{}
	if err := c.Bind(&paginationParams); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	paginationParams.Validate()

	invitations, total, err := h.Services(c).Invitation.GetPendingInvitationsPaginated(c.Request().Context(), user.Email, paginationParams)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to get pending invitations")
	}

	return responses.Response(c, http.StatusOK, responses.NewPaginatedInvitationsResponse(invitations, paginationParams.Page, paginationParams.PerPage, total))
}

// ResendInvitation resends an existing invitation with a new expiry date
func (h *InvitationHandler) ResendInvitation(c echo.Context) error {
	invitationID, err := uuid.Parse(c.Param("invitation_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid invitation ID format")
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	err = h.Services(c).Invitation.ResendInvitation(c.Request().Context(), invitationID, user.ID)
	if err != nil {
		if err == invitation_service.ErrInvalidInvitation {
			return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid invitation")
		}
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to resend invitation")
	}

	return responses.MessageResponse(c, http.StatusOK, "Invitation resent successfully")
}

// RevokeInvitation revokes an existing invitation by marking it as ignored
func (h *InvitationHandler) RevokeInvitation(c echo.Context) error {
	invitationID, err := uuid.Parse(c.Param("invitation_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid invitation ID format")
	}

	err = h.Services(c).Invitation.IgnoreInvitation(c.Request().Context(), invitationID)
	if err != nil {
		if err == invitation_service.ErrInvalidInvitation {
			return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid invitation")
		}
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to revoke invitation")
	}

	return responses.MessageResponse(c, http.StatusOK, "Invitation revoked successfully")
}
