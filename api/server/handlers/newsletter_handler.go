package handlers

import (
	ai_client "abstractions/clients/ai"
	"abstractions/repositories"
	"abstractions/server"
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

// GetUpcoming handles GET /repositories/:repo_id/newsletters/upcoming.
// Returns the next scheduled send time and any pre-computed topics from the most recent newsletter.
func (h *NewsletterHandler) GetUpcoming(c echo.Context) error {
	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}

	ctx := c.Request().Context()
	repo, err := h.repositoryRepo.GetByID(ctx, repoID)
	if err != nil || repo == nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
	}

	nextSend := nextSendTime(repo.SendDay, repo.SendHour)

	// Use topic from the most recent newsletter if available.
	var topic *responses.TopicResponse
	newsletters, err := h.newsletterRepo.ListByRepoID(ctx, repoID)
	if err == nil && len(newsletters) > 0 {
		topic = responses.NewTopicResponse(newsletters[0].Topic)
	}

	return responses.Response(c, http.StatusOK, responses.UpcomingNewsletterResponse{
		NextSendAt: nextSend.UTC().Format(time.RFC3339),
		Topic:      topic,
	})
}

// RegenerateTopics handles POST /repositories/:repo_id/newsletters/topics/regenerate.
// Re-runs topic generation synchronously and returns the new topics.
func (h *NewsletterHandler) RegenerateTopics(c echo.Context) error {
	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}

	ctx := c.Request().Context()

	aiProvider := ai_client.NewAIProvider()
	generator := newsletter_service.NewTopicGenerator(
		aiProvider,
		repositories.NewRepositoryFileRepository(h.Server().DB),
		repositories.NewChunkEmbeddingRepository(h.Server().DB),
		repositories.NewNewsletterRepository(h.Server().DB),
		repositories.NewRepositoryRepository(h.Server().DB),
	)

	topics, err := generator.GenerateTopics(ctx, repoID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "topic generation failed: "+err.Error())
	}

	return responses.Response(c, http.StatusOK, map[string]interface{}{"topics": topics})
}

// nextSendTime returns the next UTC wall-clock moment matching sendDay (ISO 1=Mon…7=Sun) at sendHour:00.
func nextSendTime(sendDay, sendHour int) time.Time {
	now := time.Now().UTC()
	// Convert ISO weekday (Mon=1, Sun=7) to Go time.Weekday (Sun=0, Mon=1, …, Sat=6).
	targetWeekday := time.Weekday(sendDay % 7) // Mon=1→1, … Sat=6→6, Sun=7→0

	candidate := time.Date(now.Year(), now.Month(), now.Day(), sendHour, 0, 0, 0, time.UTC)
	for i := 0; i <= 7; i++ {
		if candidate.Weekday() == targetWeekday && candidate.After(now) {
			return candidate
		}
		candidate = candidate.Add(24 * time.Hour)
	}
	return candidate
}
