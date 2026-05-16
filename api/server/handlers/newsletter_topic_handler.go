package handlers

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/server"
	"abstractions/server/requests"
	"abstractions/server/responses"
	"abstractions/utils/context_utils"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// NewsletterTopicHandler handles newsletter topic management endpoints.
type NewsletterTopicHandler struct {
	*BaseHandler
	newsletterTopicRepo *repositories.NewsletterTopicRepository
}

func NewNewsletterTopicHandler(srv *server.Server) *NewsletterTopicHandler {
	return &NewsletterTopicHandler{
		BaseHandler:         NewBaseHandler(srv),
		newsletterTopicRepo: repositories.NewNewsletterTopicRepository(srv.DB),
	}
}

// ListTopics handles GET /repositories/:repo_id/newsletter-topics.
// Returns all topics for a repository (pending first by position, then sent most-recent-first).
func (h *NewsletterTopicHandler) ListTopics(c echo.Context) error {
	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}

	topics, err := h.newsletterTopicRepo.ListByRepoID(c.Request().Context(), repoID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, responses.NewNewsletterTopicResponses(topics))
}

// CreateTopic handles POST /repositories/:repo_id/newsletter-topics.
// Manually adds a topic to the repository's newsletter queue.
func (h *NewsletterTopicHandler) CreateTopic(c echo.Context) error {
	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}

	var req requests.CreateNewsletterTopicRequest
	if err := c.Bind(&req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}
	if msg := req.Validate(); msg != "" {
		return responses.ErrorResponse(c, http.StatusBadRequest, msg)
	}

	relatedFiles := req.RelatedFiles
	if relatedFiles == nil {
		relatedFiles = []string{}
	}
	relatedJSON, err := json.Marshal(relatedFiles)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "failed to encode related files")
	}

	ctx := c.Request().Context()
	maxPos, err := h.newsletterTopicRepo.GetMaxPendingPosition(ctx, repoID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	topic := &models.NewsletterTopic{
		ID:           uuid.New(),
		RepoID:       repoID,
		Title:        req.Title,
		Description:  req.Description,
		RelatedFiles: models.JSONBRaw(relatedJSON),
		Position:     maxPos + 1,
	}

	if user, err := context_utils.GetUserFromContext(c); err == nil {
		topic.CreatedBy = &user.ID
	}
	if err := h.newsletterTopicRepo.Create(ctx, topic); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusCreated, responses.NewNewsletterTopicResponse(topic))
}

// UpdateTopic handles PUT /repositories/:repo_id/newsletter-topics/:topic_id.
// Only unsent topics can be updated.
func (h *NewsletterTopicHandler) UpdateTopic(c echo.Context) error {
	topicID, err := uuid.Parse(c.Param("topic_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid topic_id format")
	}

	ctx := c.Request().Context()
	topic, err := h.newsletterTopicRepo.GetByID(ctx, topicID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	if topic == nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "topic not found")
	}
	if topic.SentAt != nil {
		return responses.ErrorResponse(c, http.StatusConflict, "cannot update a topic that has already been sent")
	}

	var req requests.UpdateNewsletterTopicRequest
	if err := c.Bind(&req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}
	if msg := req.Validate(); msg != "" {
		return responses.ErrorResponse(c, http.StatusBadRequest, msg)
	}

	relatedFiles := req.RelatedFiles
	if relatedFiles == nil {
		relatedFiles = []string{}
	}
	relatedJSON, err := json.Marshal(relatedFiles)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "failed to encode related files")
	}

	if err := h.newsletterTopicRepo.UpdateTopic(ctx, topicID, req.Title, req.Description, models.JSONBRaw(relatedJSON)); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	topic.Title = req.Title
	topic.Description = req.Description
	topic.RelatedFiles = models.JSONBRaw(relatedJSON)
	return responses.Response(c, http.StatusOK, responses.NewNewsletterTopicResponse(topic))
}

// ReorderTopic handles PUT /repositories/:repo_id/newsletter-topics/:topic_id/reorder.
// Moves a pending topic up, down, or to the top of the queue.
func (h *NewsletterTopicHandler) ReorderTopic(c echo.Context) error {
	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}
	topicID, err := uuid.Parse(c.Param("topic_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid topic_id format")
	}

	ctx := c.Request().Context()
	topic, err := h.newsletterTopicRepo.GetByID(ctx, topicID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	if topic == nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "topic not found")
	}
	if topic.SentAt != nil {
		return responses.ErrorResponse(c, http.StatusConflict, "cannot reorder a topic that has already been sent")
	}

	var req requests.ReorderNewsletterTopicRequest
	if err := c.Bind(&req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}
	if msg := req.Validate(); msg != "" {
		return responses.ErrorResponse(c, http.StatusBadRequest, msg)
	}

	if err := h.newsletterTopicRepo.ReorderTopic(ctx, repoID, topicID, req.Direction); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteTopic handles DELETE /repositories/:repo_id/newsletter-topics/:topic_id.
// Only pending (unsent) topics can be deleted.
func (h *NewsletterTopicHandler) DeleteTopic(c echo.Context) error {
	topicID, err := uuid.Parse(c.Param("topic_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid topic_id format")
	}

	ctx := c.Request().Context()
	topic, err := h.newsletterTopicRepo.GetByID(ctx, topicID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	if topic == nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "topic not found")
	}
	if topic.SentAt != nil {
		return responses.ErrorResponse(c, http.StatusConflict, "cannot delete a topic that has already been sent")
	}

	if err := h.newsletterTopicRepo.Delete(ctx, topicID); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
