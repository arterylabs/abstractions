package handlers

import (
	"abstractions/server"
	"abstractions/server/responses"
	newsletter_service "abstractions/services/newsletter"
	repository_service "abstractions/services/repository"
	"abstractions/workers/jobs"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// NewsletterHandler handles newsletter-related HTTP endpoints.
type NewsletterHandler struct {
	*BaseHandler
	jobClient *jobs.Client
}

func NewNewsletterHandler(srv *server.Server, jobClient *jobs.Client) *NewsletterHandler {
	return &NewsletterHandler{
		BaseHandler: NewBaseHandler(srv),
		jobClient:   jobClient,
	}
}

// TriggerNewsletter handles POST /workspaces/:workspace_id/repositories/:repo_id/trigger_newsletter.
// Rate-limited: at most once every 24 hours per repository.
func (h *NewsletterHandler) TriggerNewsletter(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("workspace_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}

	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}

	repo, nl, err := h.Services(c).Newsletter.QueueNewsletter(c.Request().Context(), workspaceID, repoID)
	if err != nil {
		if errors.Is(err, repository_service.ErrRepositoryNotFound) {
			return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
		}
		return h.internalError(c, err)
	}

	if err := h.jobClient.EnqueueGenerateNewsletter(repo.ID, nl.ID); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "failed to enqueue newsletter generation")
	}

	return c.JSON(http.StatusAccepted, map[string]string{"message": "newsletter generation queued"})
}

// ListNewsletters handles GET /workspaces/:workspace_id/repositories/:repo_id/newsletters.
func (h *NewsletterHandler) ListNewsletters(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("workspace_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}

	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}

	newsletters, err := h.Services(c).Newsletter.ListNewsletters(c.Request().Context(), workspaceID, repoID)
	if err != nil {
		if errors.Is(err, repository_service.ErrRepositoryNotFound) {
			return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
		}
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusOK, responses.NewNewsletterSummaryResponses(newsletters))
}

// GetNewsletter handles GET /workspaces/:workspace_id/repositories/:repo_id/newsletters/:newsletter_id.
// Returns the full newsletter including HTML content and topics.
func (h *NewsletterHandler) GetNewsletter(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("workspace_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}

	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}

	newsletterID, err := uuid.Parse(c.Param("newsletter_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid newsletter_id format")
	}

	newsletter, err := h.Services(c).Newsletter.GetNewsletter(c.Request().Context(), workspaceID, repoID, newsletterID)
	if err != nil {
		if errors.Is(err, repository_service.ErrRepositoryNotFound) {
			return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
		}
		if errors.Is(err, newsletter_service.ErrNewsletterNotFound) {
			return responses.ErrorResponse(c, http.StatusNotFound, "newsletter not found")
		}
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusOK, responses.NewNewsletterDetailResponse(newsletter))
}

// GetEmailContent handles GET /workspaces/:workspace_id/repositories/:repo_id/newsletters/:newsletter_id/email_content.
// Returns the newsletter rendered inside the generic email template as text/html,
// matching exactly what subscribers receive.
func (h *NewsletterHandler) GetEmailContent(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("workspace_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}

	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}

	newsletterID, err := uuid.Parse(c.Param("newsletter_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid newsletter_id format")
	}

	newsletter, err := h.Services(c).Newsletter.GetNewsletter(c.Request().Context(), workspaceID, repoID, newsletterID)
	if err != nil {
		if errors.Is(err, repository_service.ErrRepositoryNotFound) {
			return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
		}
		if errors.Is(err, newsletter_service.ErrNewsletterNotFound) {
			return responses.ErrorResponse(c, http.StatusNotFound, "newsletter not found")
		}
		return h.internalError(c, err)
	}

	return c.HTML(http.StatusOK, newsletter.HTMLContent)
}

// ResendNewsletter handles POST /workspaces/:workspace_id/repositories/:repo_id/newsletters/:newsletter_id/resend.
// It is only allowed when the newsletter has status=error.
// The same newsletter record is reset to queued (clearing the error) and re-enqueued.
func (h *NewsletterHandler) ResendNewsletter(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("workspace_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}

	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}

	newsletterID, err := uuid.Parse(c.Param("newsletter_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid newsletter_id format")
	}

	repo, nl, err := h.Services(c).Newsletter.ResetNewsletterForResend(c.Request().Context(), workspaceID, repoID, newsletterID)
	if err != nil {
		if errors.Is(err, repository_service.ErrRepositoryNotFound) {
			return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
		}
		if errors.Is(err, newsletter_service.ErrNewsletterNotFound) {
			return responses.ErrorResponse(c, http.StatusNotFound, "newsletter not found")
		}
		if errors.Is(err, newsletter_service.ErrNewsletterNotResendable) {
			return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		return h.internalError(c, err)
	}

	if err := h.jobClient.EnqueueGenerateNewsletter(repo.ID, nl.ID); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "failed to enqueue newsletter generation")
	}

	return c.JSON(http.StatusAccepted, map[string]string{"message": "newsletter re-queued for sending"})
}
