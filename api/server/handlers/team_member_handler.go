package handlers

import (
	"abstractions/server"
	"abstractions/server/requests"
	"abstractions/server/responses"
	team_member_service "abstractions/services/team_member"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TeamMemberHandler struct {
	*BaseHandler
}

func NewTeamMemberHandler(server *server.Server) *TeamMemberHandler {
	return &TeamMemberHandler{
		BaseHandler: NewBaseHandler(server),
	}
}

// ListMembers returns all team members for a repository.
func (h *TeamMemberHandler) ListMembers(c echo.Context) error {
	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}

	members, err := h.Services(c).TeamMember.ListMembers(c.Request().Context(), repoID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, responses.NewTeamMemberResponses(members))
}

// AddMember adds a new member to the repository's newsletter.
// On the free tier, a maximum of 10 members is enforced.
func (h *TeamMemberHandler) AddMember(c echo.Context) error {
	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}

	req := new(requests.AddTeamMemberRequest)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Determine whether the workspace has an active paid subscription.
	// We need the repository to get the workspace ID.
	repo, err := h.Services(c).Repository.GetRepository(c.Request().Context(), repoID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	if repo == nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
	}

	activeSub, _ := h.Services(c).Subscription.GetActiveSubscriptionByWorkspaceID(repo.WorkspaceID)
	hasActiveSub := activeSub != nil

	member, err := h.Services(c).TeamMember.AddMember(c.Request().Context(), repoID, req.Email, hasActiveSub)
	if err != nil {
		switch {
		case errors.Is(err, team_member_service.ErrDuplicateEmail):
			return responses.ErrorResponse(c, http.StatusConflict, err.Error())
		case errors.Is(err, team_member_service.ErrFreeTierLimitHit):
			return responses.ErrorResponse(c, http.StatusPaymentRequired, err.Error())
		default:
			return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
	}

	return responses.Response(c, http.StatusCreated, responses.NewTeamMemberResponse(member))
}

// RemoveMember deletes a team member from the repository.
func (h *TeamMemberHandler) RemoveMember(c echo.Context) error {
	memberID, err := uuid.Parse(c.Param("member_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid member_id format")
	}

	if err := h.Services(c).TeamMember.RemoveMember(c.Request().Context(), memberID); err != nil {
		if errors.Is(err, team_member_service.ErrMemberNotFound) {
			return responses.ErrorResponse(c, http.StatusNotFound, "team member not found")
		}
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
