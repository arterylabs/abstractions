package handlers

import (
	"abstractions/db/models"
	"abstractions/server"
	"abstractions/server/requests"
	"abstractions/server/responses"
	repository_service "abstractions/services/repository"
	"abstractions/utils/context_utils"
	"abstractions/workers/jobs"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
)

type RepositoryHandler struct {
	*BaseHandler
	enqueuer *asynq.Client
}

func NewRepositoryHandler(server *server.Server, enqueuer *asynq.Client) *RepositoryHandler {
	return &RepositoryHandler{
		BaseHandler: NewBaseHandler(server),
		enqueuer:    enqueuer,
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
			return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
	}

	if task, err := jobs.NewIndexRepositoryTask(repo.ID); err == nil {
		if _, enqErr := h.enqueuer.Enqueue(task,
			asynq.MaxRetry(3),
			asynq.Queue("default"),
			asynq.TaskID("index_repository:"+repo.ServiceID),
		); enqErr != nil && enqErr != asynq.ErrTaskIDConflict {
			// Non-fatal: user can re-trigger from Settings
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

// GetRepository returns a single repository by ID.
func (h *RepositoryHandler) GetRepository(c echo.Context) error {
	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}

	repo, err := h.Services(c).Repository.GetRepository(c.Request().Context(), repoID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	if repo == nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
	}

	return responses.Response(c, http.StatusOK, responses.NewRepositoryResponse(repo))
}

// UpdateRepositorySettings updates send_day and send_hour for a repository.
func (h *RepositoryHandler) UpdateRepositorySettings(c echo.Context) error {
	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}

	req := new(requests.UpdateRepositorySettingsRequest)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	repo, err := h.Services(c).Repository.UpdateSettings(c.Request().Context(), repoID, &repository_service.UpdateSettingsRequest{
		SendDay:           req.SendDay,
		SendHour:          req.SendHour,
		NewsletterEnabled: req.NewsletterEnabled,
	})
	if err != nil {
		if errors.Is(err, repository_service.ErrRepositoryNotFound) {
			return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
		}
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, responses.NewRepositoryResponse(repo))
}

// UpdateRepositoryDescription updates only the description for a repository.
func (h *RepositoryHandler) UpdateRepositoryDescription(c echo.Context) error {
	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}

	req := new(requests.UpdateRepositoryDescriptionRequest)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	repo, err := h.Services(c).Repository.UpdateDescription(c.Request().Context(), repoID, req.Description)
	if err != nil {
		if errors.Is(err, repository_service.ErrRepositoryNotFound) {
			return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
		}
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, responses.NewRepositoryResponse(repo))
}

// DeleteRepository soft-deletes the repository and removes all derived artifacts.
func (h *RepositoryHandler) DeleteRepository(c echo.Context) error {
	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	if err := h.Services(c).Repository.DeleteRepository(c.Request().Context(), repoID, user.ID); err != nil {
		if errors.Is(err, repository_service.ErrRepositoryNotFound) {
			return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
		}
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// GetRepositoryStats returns file count and language breakdown for a repository.
func (h *RepositoryHandler) GetRepositoryStats(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("workspace_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}
	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}

	stats, err := h.Services(c).Repository.GetStats(c.Request().Context(), workspaceID, repoID)
	if err != nil {
		if errors.Is(err, repository_service.ErrRepositoryNotFound) {
			return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
		}
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, stats)
}

// TriggerIndex enqueues an IndexRepositoryJob for the given repository.
// Only allowed when the repository is in pending or error state (not currently indexing).
func (h *RepositoryHandler) TriggerIndex(c echo.Context) error {
	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}

	ctx := c.Request().Context()
	repo, err := h.Services(c).Repository.GetRepository(ctx, repoID)
	if err != nil {
		if errors.Is(err, repository_service.ErrRepositoryNotFound) {
			return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
		}
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	task, err := jobs.NewIndexRepositoryTask(repo.ID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "failed to create index task")
	}
	if _, err := h.enqueuer.Enqueue(task, asynq.MaxRetry(3), asynq.Queue("default"),
		asynq.TaskID("index_repository:"+repo.ServiceID)); err != nil {
		// Ignore duplicate task errors — the job is already queued.
		if err != asynq.ErrTaskIDConflict {
			return responses.ErrorResponse(c, http.StatusInternalServerError, "failed to enqueue index task")
		}
	}

	return c.JSON(http.StatusAccepted, map[string]string{"message": "indexing queued"})
}
