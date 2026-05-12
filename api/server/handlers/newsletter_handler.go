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
	newsletterRepo *repositories.NewsletterRepository
	enqueuer       *asynq.Client
}

func NewNewsletterHandler(srv *server.Server, enqueuer *asynq.Client) *NewsletterHandler {
	return &NewsletterHandler{
		BaseHandler:    NewBaseHandler(srv),
		repositoryRepo: repositories.NewRepositoryRepository(srv.DB),
		newsletterRepo: repositories.NewNewsletterRepository(srv.DB),
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
	if repo.NewsletterSentAt != nil && time.Since(*repo.NewsletterSentAt) < newsletterTriggerRateLimit {
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

// ListNewsletters handles GET /repositories/:repo_id/newsletters.
// Returns newsletter summaries with open_rate and click_rate.
func (h *NewsletterHandler) ListNewsletters(c echo.Context) error {
	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}

	newsletters, err := h.newsletterRepo.ListByRepoID(c.Request().Context(), repoID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, responses.NewNewsletterSummaryResponses(newsletters))
}

// GetNewsletter handles GET /repositories/:repo_id/newsletters/:newsletter_id.
// Returns the full newsletter including HTML content and topics.
func (h *NewsletterHandler) GetNewsletter(c echo.Context) error {
	newsletterID, err := uuid.Parse(c.Param("newsletter_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid newsletter_id format")
	}

	newsletter, err := h.newsletterRepo.GetByID(c.Request().Context(), newsletterID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	if newsletter == nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "newsletter not found")
	}

	return responses.Response(c, http.StatusOK, responses.NewNewsletterDetailResponse(newsletter))
}

// GetEmailContent handles GET /repositories/:repo_id/newsletters/:newsletter_id/email-content.
// Returns the stored full HTML email as text/html for previewing in an iframe.
func (h *NewsletterHandler) GetEmailContent(c echo.Context) error {
	newsletterID, err := uuid.Parse(c.Param("newsletter_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid newsletter_id format")
	}

	newsletter, err := h.newsletterRepo.GetByID(c.Request().Context(), newsletterID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	if newsletter == nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "newsletter not found")
	}

	return c.HTML(http.StatusOK, newsletter.HTMLContent)
}
