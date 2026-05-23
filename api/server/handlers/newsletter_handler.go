package handlers

import (
	"abstractions/config"
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/server"
	"abstractions/server/responses"
	"abstractions/workers/jobs"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// NewsletterHandler handles newsletter-related HTTP endpoints.
type NewsletterHandler struct {
	*BaseHandler
	cfg                 *config.Config
	repositoryRepo      *repositories.RepositoryRepository
	newsletterRepo      *repositories.NewsletterRepository
	newsletterTopicRepo *repositories.NewsletterTopicRepository
	workspaceRepo       *repositories.WorkspaceRepository
	jobClient           *jobs.Client
}

func NewNewsletterHandler(srv *server.Server, jobClient *jobs.Client) *NewsletterHandler {
	return &NewsletterHandler{
		BaseHandler:         NewBaseHandler(srv),
		cfg:                 srv.Config,
		repositoryRepo:      repositories.NewRepositoryRepository(srv.DB),
		newsletterRepo:      repositories.NewNewsletterRepository(srv.DB),
		newsletterTopicRepo: repositories.NewNewsletterTopicRepository(srv.DB),
		workspaceRepo:       repositories.NewWorkspaceRepository(srv.DB),
		jobClient:           jobClient,
	}
}

// TriggerNewsletter handles POST /workspaces/:workspace_id/repositories/:repo_id/trigger-newsletter.
// Rate-limited: at most once every 24 hours per repository.
func (h *NewsletterHandler) TriggerNewsletter(c echo.Context) error {
	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}

	// // Rate-limit: reject if a newsletter was sent within the last 24 hours.
	// if repo.LastNewsletterSentAt != nil && time.Since(*repo.LastNewsletterSentAt) < newsletterTriggerRateLimit {
	// 	return responses.ErrorResponse(c, http.StatusTooManyRequests, "newsletter already sent within the last 24 hours")
	// }

	// Create the newsletter record immediately so it appears in the list as queued.
	nl := &models.Newsletter{
		ID:     uuid.New(),
		RepoID: repoID,
		Status: models.NewsletterStatusQueued,
	}
	if err := h.newsletterRepo.Create(c.Request().Context(), nl); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "failed to create newsletter record")
	}

	if err := h.jobClient.EnqueueGenerateNewsletter(repoID, nl.ID); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "failed to enqueue newsletter generation")
	}

	return c.JSON(http.StatusAccepted, map[string]string{"message": "newsletter generation queued"})
}

// ListNewsletters handles GET /workspaces/:workspace_id/repositories/:repo_id/newsletters.
func (h *NewsletterHandler) ListNewsletters(c echo.Context) error {
	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}

	newsletters, err := h.newsletterRepo.ListByRepoID(c.Request().Context(), repoID)
	if err != nil {
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusOK, responses.NewNewsletterSummaryResponses(newsletters))
}

// GetNewsletter handles GET /workspaces/:workspace_id/repositories/:repo_id/newsletters/:newsletter_id.
// Returns the full newsletter including HTML content and topics.
func (h *NewsletterHandler) GetNewsletter(c echo.Context) error {
	newsletterID, err := uuid.Parse(c.Param("newsletter_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid newsletter_id format")
	}

	newsletter, err := h.newsletterRepo.GetByID(c.Request().Context(), newsletterID)
	if err != nil {
		return h.internalError(c, err)
	}
	if newsletter == nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "newsletter not found")
	}

	return responses.Response(c, http.StatusOK, responses.NewNewsletterDetailResponse(newsletter))
}

// GetEmailContent handles GET /workspaces/:workspace_id/repositories/:repo_id/newsletters/:newsletter_id/email-content.
// Returns the newsletter rendered inside the generic email template as text/html,
// matching exactly what subscribers receive.
func (h *NewsletterHandler) GetEmailContent(c echo.Context) error {
	newsletterID, err := uuid.Parse(c.Param("newsletter_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid newsletter_id format")
	}

	newsletter, err := h.newsletterRepo.GetByID(c.Request().Context(), newsletterID)
	if err != nil {
		return h.internalError(c, err)
	}
	if newsletter == nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "newsletter not found")
	}

	return c.HTML(http.StatusOK, newsletter.HTMLContent)
}

// ResendNewsletter handles POST /workspaces/:workspace_id/repositories/:repo_id/newsletters/:newsletter_id/resend.
// It is only allowed when the newsletter has status=error.
// The same newsletter record is reset to queued (clearing the error) and re-enqueued.
func (h *NewsletterHandler) ResendNewsletter(c echo.Context) error {
	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}

	newsletterID, err := uuid.Parse(c.Param("newsletter_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid newsletter_id format")
	}

	nl, err := h.newsletterRepo.GetByID(c.Request().Context(), newsletterID)
	if err != nil {
		return h.internalError(c, err)
	}
	if nl == nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "newsletter not found")
	}
	if nl.Status != models.NewsletterStatusError {
		return responses.ErrorResponse(c, http.StatusBadRequest, "only newsletters with status 'error' can be resent")
	}

	// Reset the existing record so it goes back through the pipeline.
	nl.Status = models.NewsletterStatusQueued
	nl.ErrorMessage = nil
	nl.HTMLContent = ""
	nl.TopicID = nil
	nl.SentAt = nil
	nl.RecipientCount = 0
	if err := h.newsletterRepo.Update(c.Request().Context(), nl); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "failed to reset newsletter status")
	}

	if err := h.jobClient.EnqueueGenerateNewsletter(repoID, nl.ID); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "failed to enqueue newsletter generation")
	}

	return c.JSON(http.StatusAccepted, map[string]string{"message": "newsletter re-queued for sending"})
}
