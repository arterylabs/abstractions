package handlers

import (
	"abstractions/db/models"
	"abstractions/logger"
	"abstractions/server"
	"abstractions/server/requests"
	"abstractions/server/responses"
	repository_service "abstractions/services/repository"
	"abstractions/utils/context_utils"
	"abstractions/workers/jobs"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RepositoryHandler struct {
	*BaseHandler
	jobClient *jobs.Client
}

func NewRepositoryHandler(server *server.Server, jobClient *jobs.Client) *RepositoryHandler {
	return &RepositoryHandler{
		BaseHandler: NewBaseHandler(server),
		jobClient:   jobClient,
	}
}

// ListGitHubRepos returns all repositories from the user's connected GitHub account.
func (h *RepositoryHandler) ListGitHubRepos(c echo.Context) error {
	workspaceID, err := h.uuidParam(c, "workspace_id")
	if err != nil {
		return err
	}

	repos, err := h.Services(c).Repository.ListGitHubRepos(c.Request().Context(), workspaceID)
	if err != nil {
		if errors.Is(err, repository_service.ErrGitHubNotConnected) {
			return responses.ErrorResponse(c, http.StatusUnprocessableEntity, "GitHub integration not connected")
		}
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusOK, responses.NewGitHubRepoResponses(repos))
}

// ConnectRepository connects a GitHub repository to the workspace.
func (h *RepositoryHandler) ConnectRepository(c echo.Context) error {
	workspaceID, err := h.uuidParam(c, "workspace_id")
	if err != nil {
		return err
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
		CreatedBy:     user.ID,
		ServiceType:   models.RepositoryServiceTypeGitHub,
		ServiceID:     req.ServiceID,
		Name:          req.Name,
		Description:   req.Description,
		DefaultBranch: defaultBranch,
	}

	repo, err := h.Services(c).Repository.ConnectRepository(c.Request().Context(), svcReq)
	if err != nil {
		switch {
		case errors.Is(err, repository_service.ErrAlreadyConnected):
			return responses.ErrorResponse(c, http.StatusConflict, "repository already connected to this workspace")
		case errors.Is(err, repository_service.ErrGitHubNotConnected):
			return responses.ErrorResponse(c, http.StatusUnprocessableEntity, "GitHub integration not connected")
		default:
			return h.internalError(c, err)
		}
	}

	if err := h.jobClient.EnqueueIndexRepository(repo.ID, repo.ServiceID); err != nil {
		logger.Error(c.Request().Context(), err)
	}

	return responses.Response(c, http.StatusCreated, responses.NewRepositoryResponse(repo))
}

// ListRepositories returns all repositories connected to the workspace.
func (h *RepositoryHandler) ListRepositories(c echo.Context) error {
	workspaceID, err := h.uuidParam(c, "workspace_id")
	if err != nil {
		return err
	}

	repos, err := h.Services(c).Repository.ListRepositories(c.Request().Context(), workspaceID)
	if err != nil {
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusOK, responses.NewRepositoryResponses(repos))
}

// GetRepository returns a single repository by ID in a workspace.
func (h *RepositoryHandler) GetRepository(c echo.Context) error {
	workspaceID, err := h.uuidParam(c, "workspace_id")
	if err != nil {
		return err
	}

	repoID, err := h.uuidParam(c, "repo_id")
	if err != nil {
		return err
	}

	repo, err := h.Services(c).Repository.GetRepositoryInWorkspace(c.Request().Context(), workspaceID, repoID)
	if err != nil {
		if errors.Is(err, repository_service.ErrRepositoryNotFound) {
			return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
		}
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusOK, responses.NewRepositoryResponse(repo))
}

// UpdateRepositorySettings updates send_day and send_hour for a repository.
func (h *RepositoryHandler) UpdateRepositorySettings(c echo.Context) error {
	workspaceID, err := h.uuidParam(c, "workspace_id")
	if err != nil {
		return err
	}

	repoID, err := h.uuidParam(c, "repo_id")
	if err != nil {
		return err
	}

	req := new(requests.UpdateRepositorySettingsRequest)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	repo, err := h.Services(c).Repository.UpdateSettings(c.Request().Context(), workspaceID, repoID, &repository_service.UpdateSettingsRequest{
		SendDay:           req.SendDay,
		SendHour:          req.SendHour,
		NewsletterEnabled: req.NewsletterEnabled,
	})
	if err != nil {
		if errors.Is(err, repository_service.ErrRepositoryNotFound) {
			return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
		}
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusOK, responses.NewRepositoryResponse(repo))
}

// UpdateRepositoryDescription updates only the description for a repository.
func (h *RepositoryHandler) UpdateRepositoryDescription(c echo.Context) error {
	workspaceID, err := h.uuidParam(c, "workspace_id")
	if err != nil {
		return err
	}

	repoID, err := h.uuidParam(c, "repo_id")
	if err != nil {
		return err
	}

	req := new(requests.UpdateRepositoryDescriptionRequest)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	repo, err := h.Services(c).Repository.UpdateDescription(c.Request().Context(), workspaceID, repoID, req.Description)
	if err != nil {
		if errors.Is(err, repository_service.ErrRepositoryNotFound) {
			return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
		}
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusOK, responses.NewRepositoryResponse(repo))
}

// DeleteRepository soft-deletes the repository and removes all derived artifacts.
func (h *RepositoryHandler) DeleteRepository(c echo.Context) error {
	workspaceID, err := h.uuidParam(c, "workspace_id")
	if err != nil {
		return err
	}

	repoID, err := h.uuidParam(c, "repo_id")
	if err != nil {
		return err
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return h.internalError(c, err)
	}

	if err := h.Services(c).Repository.DeleteRepository(c.Request().Context(), workspaceID, repoID, user.ID); err != nil {
		if errors.Is(err, repository_service.ErrRepositoryNotFound) {
			return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
		}
		return h.internalError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// GetRepositoryStats returns file count and language breakdown for a repository.
func (h *RepositoryHandler) GetRepositoryStats(c echo.Context) error {
	workspaceID, err := h.uuidParam(c, "workspace_id")
	if err != nil {
		return err
	}
	repoID, err := h.uuidParam(c, "repo_id")
	if err != nil {
		return err
	}

	stats, err := h.Services(c).Repository.GetStats(c.Request().Context(), workspaceID, repoID)
	if err != nil {
		if errors.Is(err, repository_service.ErrRepositoryNotFound) {
			return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
		}
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusOK, stats)
}

// GetRepositoryStatus returns live indexing progress and the repository status.
func (h *RepositoryHandler) GetRepositoryStatus(c echo.Context) error {
	workspaceID, err := h.uuidParam(c, "workspace_id")
	if err != nil {
		return err
	}
	repoID, err := h.uuidParam(c, "repo_id")
	if err != nil {
		return err
	}

	status, err := h.Services(c).Repository.GetStatus(c.Request().Context(), workspaceID, repoID)
	if err != nil {
		if errors.Is(err, repository_service.ErrRepositoryNotFound) {
			return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
		}
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusOK, status)
}

// TriggerIndex enqueues an IndexRepositoryJob for the given repository.
// Only allowed when the repository is in pending or error state (not currently indexing).
func (h *RepositoryHandler) TriggerIndex(c echo.Context) error {
	workspaceID, err := h.uuidParam(c, "workspace_id")
	if err != nil {
		return err
	}

	repoID, err := h.uuidParam(c, "repo_id")
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	repo, err := h.Services(c).Repository.GetRepositoryInWorkspace(ctx, workspaceID, repoID)
	if err != nil {
		if errors.Is(err, repository_service.ErrRepositoryNotFound) {
			return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
		}
		return h.internalError(c, err)
	}

	if err := h.jobClient.EnqueueIndexRepository(repo.ID, repo.ServiceID); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "failed to enqueue index task")
	}

	return c.JSON(http.StatusAccepted, map[string]string{"message": "indexing queued"})
}
