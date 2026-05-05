package handlers

import (
	"abstractions/server"
	"abstractions/server/requests"
	"abstractions/server/responses"
	repository_service "abstractions/services/repository"
	"abstractions/utils/context_utils"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RepositoryHandler struct {
	*BaseHandler
}

func NewRepositoryHandler(server *server.Server) *RepositoryHandler {
	return &RepositoryHandler{
		BaseHandler: NewBaseHandler(server),
	}
}

// ListGitHubRepos returns all repositories from the user's connected GitHub account.
func (h *RepositoryHandler) ListGitHubRepos(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("workspace_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}

	repos, err := h.Services(c).Repository.ListGitHubRepos(c.Request().Context(), workspaceID)
	if err != nil {
		if errors.Is(err, repository_service.ErrGitHubNotConnected) {
			return responses.ErrorResponse(c, http.StatusUnprocessableEntity, "GitHub integration not connected")
		}
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, responses.NewGitHubRepoResponses(repos))
}

// ConnectRepository connects a GitHub repository to the workspace.
func (h *RepositoryHandler) ConnectRepository(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("workspace_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "user not found")
	}

	req := new(requests.ConnectRepositoryRequest)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	defaultBranch := req.DefaultBranch
	if defaultBranch == "" {
		defaultBranch = "main"
	}

	svcReq := &repository_service.ConnectRepositoryRequest{
		WorkspaceID:   workspaceID,
		UserID:        user.ID,
		GitHubRepoID:  req.GitHubRepoID,
		FullName:      req.FullName,
		DefaultBranch: defaultBranch,
	}

	repo, err := h.Services(c).Repository.ConnectRepository(c.Request().Context(), svcReq)
	if err != nil {
		switch {
		case errors.Is(err, repository_service.ErrAlreadyConnected):
			return responses.ErrorResponse(c, http.StatusConflict, "repository already connected to this workspace")
		case errors.Is(err, repository_service.ErrGitHubNotConnected):
			return responses.ErrorResponse(c, http.StatusUnprocessableEntity, "GitHub integration not connected")
		case errors.Is(err, repository_service.ErrWebhookCreationFailed):
			return responses.ErrorResponse(c, http.StatusBadGateway, "failed to create webhook on GitHub")
		default:
			return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
	}

	return responses.Response(c, http.StatusCreated, responses.NewRepositoryResponse(repo))
}

// ListRepositories returns all repositories connected to the workspace.
func (h *RepositoryHandler) ListRepositories(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("workspace_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}

	repos, err := h.Services(c).Repository.ListRepositories(c.Request().Context(), workspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, responses.NewRepositoryResponses(repos))
}
