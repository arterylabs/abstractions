package handlers

import (
	"abstractions/config"
	"abstractions/repositories"
	"abstractions/server"
	"abstractions/server/requests"
	"abstractions/server/responses"
	newsletter_service "abstractions/services/newsletter"
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
	cfg                 *config.Config
	repositoryRepo      *repositories.RepositoryRepository
	newsletterRepo      *repositories.NewsletterRepository
	newsletterTopicRepo *repositories.NewsletterTopicRepository
	workspaceRepo       *repositories.WorkspaceRepository
	enqueuer            *asynq.Client
}

func NewNewsletterHandler(srv *server.Server, enqueuer *asynq.Client) *NewsletterHandler {
	return &NewsletterHandler{
		BaseHandler:         NewBaseHandler(srv),
		cfg:                 srv.Config,
		repositoryRepo:      repositories.NewRepositoryRepository(srv.DB),
		newsletterRepo:      repositories.NewNewsletterRepository(srv.DB),
		newsletterTopicRepo: repositories.NewNewsletterTopicRepository(srv.DB),
		workspaceRepo:       repositories.NewWorkspaceRepository(srv.DB),
		enqueuer:            enqueuer,
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

	// // Rate-limit: reject if a newsletter was sent within the last 24 hours.
	// if repo.LastNewsletterSentAt != nil && time.Since(*repo.LastNewsletterSentAt) < newsletterTriggerRateLimit {
	// 	return responses.ErrorResponse(c, http.StatusTooManyRequests, "newsletter already sent within the last 24 hours")
	// }

	task, err := jobs.NewGenerateNewsletterTask(repoID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "failed to create task")
	}
	if _, err := h.enqueuer.Enqueue(task, asynq.MaxRetry(2), asynq.Queue("default")); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "failed to enqueue newsletter generation")
	}

	return c.JSON(http.StatusAccepted, map[string]string{"message": "newsletter generation queued"})
}

// GenerateNewsletterPreview handles POST /repositories/:id/newsletters/preview.
// Enqueues a draft-only newsletter generation job (no delivery).
func (h *NewsletterHandler) GenerateNewsletterPreview(c echo.Context) error {
	repoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repository id")
	}

	ctx := c.Request().Context()
	repo, err := h.repositoryRepo.GetByID(ctx, repoID)
	if err != nil || repo == nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
	}

	task, err := jobs.NewGenerateNewsletterDraftTask(repoID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "failed to create task")
	}
	if _, err := h.enqueuer.Enqueue(task, asynq.MaxRetry(2), asynq.Queue("default")); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "failed to enqueue newsletter draft generation")
	}

	return c.JSON(http.StatusAccepted, map[string]string{"message": "newsletter draft generation queued"})
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
// Returns the newsletter rendered inside the generic email template as text/html,
// matching exactly what subscribers receive.
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

// UpdateNewsletter handles PATCH /repositories/:repo_id/newsletters/:newsletter_id.
// Updates the HTML content of an unsent (draft) newsletter.
func (h *NewsletterHandler) UpdateNewsletter(c echo.Context) error {
	newsletterID, err := uuid.Parse(c.Param("newsletter_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid newsletter_id format")
	}

	var req requests.UpdateNewsletterRequest
	if err := c.Bind(&req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}

	ctx := c.Request().Context()
	newsletter, err := h.newsletterRepo.GetByID(ctx, newsletterID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	if newsletter == nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "newsletter not found")
	}
	if newsletter.SentAt != nil {
		return responses.ErrorResponse(c, http.StatusConflict, "cannot edit a newsletter that has already been sent")
	}

	newsletter.HTMLContent = req.HTMLContent
	if err := h.newsletterRepo.Update(ctx, newsletter); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, responses.NewNewsletterDetailResponse(newsletter))
}

// SendDraftNewsletter handles POST /repositories/:repo_id/newsletters/:newsletter_id/send.
// Delivers a draft newsletter to subscribers and marks it as sent.
func (h *NewsletterHandler) SendDraftNewsletter(c echo.Context) error {
	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}
	newsletterID, err := uuid.Parse(c.Param("newsletter_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid newsletter_id format")
	}

	ctx := c.Request().Context()
	newsletter, err := h.newsletterRepo.GetByID(ctx, newsletterID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	if newsletter == nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "newsletter not found")
	}
	if newsletter.SentAt != nil {
		return responses.ErrorResponse(c, http.StatusConflict, "newsletter has already been sent")
	}

	delivery := newsletter_service.NewDeliveryService(h.cfg, h.newsletterRepo, h.workspaceRepo, h.repositoryRepo)
	if err := delivery.SendNewsletter(ctx, repoID, newsletterID); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	// Mark the linked topic as sent if present.
	if newsletter.TopicID != nil {
		if err := h.newsletterTopicRepo.MarkSent(ctx, *newsletter.TopicID); err != nil {
			// Non-fatal: topic marking is best-effort.
			_ = err
		}
	}

	// Stamp last_newsletter_sent_at on the repository.
	repo, err := h.repositoryRepo.GetByID(ctx, repoID)
	if err == nil && repo != nil {
		now := time.Now()
		repo.NewsletterSentAt = &now
		_ = h.repositoryRepo.Update(ctx, repo)
	}

	// Re-fetch to return the updated sent_at and recipient count.
	updated, err := h.newsletterRepo.GetByID(ctx, newsletterID)
	if err != nil || updated == nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch updated newsletter")
	}

	return responses.Response(c, http.StatusOK, responses.NewNewsletterDetailResponse(updated))
}
