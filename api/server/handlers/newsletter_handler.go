package handlers

import (
	"abstractions/repositories"
	"abstractions/server"
	"abstractions/server/responses"
	"abstractions/workers/jobs"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
)

// newsletterTriggerRateLimit is the minimum time between manual triggers for the same repo.
const newsletterTriggerRateLimit = 24 * time.Hour

// NewsletterHandler handles newsletter-related HTTP endpoints.
type NewsletterHandler struct {
	*BaseHandler
	repositoryRepo *repositories.RepositoryRepository
	enqueuer       *asynq.Client
}

func NewNewsletterHandler(srv *server.Server, enqueuer *asynq.Client) *NewsletterHandler {
	return &NewsletterHandler{
		BaseHandler:    NewBaseHandler(srv),
		repositoryRepo: repositories.NewRepositoryRepository(srv.DB),
		enqueuer:       enqueuer,
	}
}

// TriggerNewsletter handles POST /repositories/:id/trigger-newsletter.
// Rate-limited: at most once every 24 hours per repository.
func (h *NewsletterHandler) TriggerNewsletter(c echo.Context) error {
	repoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repository id")
	}

	ctx := c.Request().Context()
	repo, err := h.repositoryRepo.GetByID(ctx, repoID)
	if err != nil || repo == nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
	}

	// Rate-limit: reject if a newsletter was sent within the last 24 hours.
	if repo.LastNewsletterSentAt != nil && time.Since(*repo.LastNewsletterSentAt) < newsletterTriggerRateLimit {
		return responses.ErrorResponse(c, http.StatusTooManyRequests, "newsletter already sent within the last 24 hours")
	}

	task, err := jobs.NewGenerateNewsletterTask(repoID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "failed to create task")
	}
	if _, err := h.enqueuer.Enqueue(task, asynq.MaxRetry(2), asynq.Queue("default")); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "failed to enqueue newsletter generation")
	}

	return c.JSON(http.StatusAccepted, map[string]string{"message": "newsletter generation queued"})
}
