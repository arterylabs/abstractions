package handlers

import (
	"abstractions/server"
	"abstractions/server/requests"
	"abstractions/server/responses"
	newsletter_service "abstractions/services/newsletter"
	repository_service "abstractions/services/repository"
	"abstractions/utils/context_utils"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// NewsletterTopicHandler handles newsletter topic management endpoints.
type NewsletterTopicHandler struct {
	*BaseHandler
}

func NewNewsletterTopicHandler(srv *server.Server) *NewsletterTopicHandler {
	return &NewsletterTopicHandler{
		BaseHandler: NewBaseHandler(srv),
	}
}

// ListNewsletterTopics handles GET /workspaces/:workspace_id/repositories/:repo_id/newsletter_topics.
// Returns all topics for a repository (pending first by position, then sent most-recent-first).
func (h *NewsletterTopicHandler) ListNewsletterTopics(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("workspace_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}
	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}

	topics, err := h.Services(c).Newsletter.ListNewsletterTopics(c.Request().Context(), workspaceID, repoID)
	if err != nil {
		if errors.Is(err, repository_service.ErrRepositoryNotFound) {
			return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
		}
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusOK, responses.NewNewsletterTopicResponses(topics))
}

// CreateNewsletterTopic handles POST /workspaces/:workspace_id/repositories/:repo_id/newsletter_topics.
// Manually adds a topic to the repository's newsletter queue.
func (h *NewsletterTopicHandler) CreateNewsletterTopic(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("workspace_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}
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

	var createdBy *uuid.UUID
	if user, err := context_utils.GetUserFromContext(c); err == nil {
		createdBy = &user.ID
	}
	topic, err := h.Services(c).Newsletter.CreateNewsletterTopic(c.Request().Context(), workspaceID, repoID, req.Title, req.Description, relatedFiles, createdBy)
	if err != nil {
		if errors.Is(err, repository_service.ErrRepositoryNotFound) {
			return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
		}
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusCreated, responses.NewNewsletterTopicResponse(topic))
}

// UpdateNewsletterTopic handles PUT /workspaces/:workspace_id/repositories/:repo_id/newsletter_topics/:topic_id.
// Only unsent topics can be updated.
func (h *NewsletterTopicHandler) UpdateNewsletterTopic(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("workspace_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}
	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}
	topicID, err := uuid.Parse(c.Param("topic_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid topic_id format")
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

	topic, err := h.Services(c).Newsletter.UpdateNewsletterTopic(c.Request().Context(), workspaceID, repoID, topicID, req.Title, req.Description, relatedFiles)
	if err != nil {
		return h.newsletterTopicError(c, err, "update")
	}
	return responses.Response(c, http.StatusOK, responses.NewNewsletterTopicResponse(topic))
}

// ReorderNewsletterTopic handles PUT /workspaces/:workspace_id/repositories/:repo_id/newsletter_topics/:topic_id/reorder.
// Moves a pending topic up, down, or to the top of the queue.
func (h *NewsletterTopicHandler) ReorderNewsletterTopic(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("workspace_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}
	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}
	topicID, err := uuid.Parse(c.Param("topic_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid topic_id format")
	}

	var req requests.ReorderNewsletterTopicRequest
	if err := c.Bind(&req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}
	if msg := req.Validate(); msg != "" {
		return responses.ErrorResponse(c, http.StatusBadRequest, msg)
	}

	if err := h.Services(c).Newsletter.ReorderNewsletterTopic(c.Request().Context(), workspaceID, repoID, topicID, req.Direction); err != nil {
		return h.newsletterTopicError(c, err, "reorder")
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteNewsletterTopic handles DELETE /workspaces/:workspace_id/repositories/:repo_id/newsletter_topics/:topic_id.
// Only pending (unsent) topics can be deleted.
func (h *NewsletterTopicHandler) DeleteNewsletterTopic(c echo.Context) error {
	workspaceID, err := uuid.Parse(c.Param("workspace_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}
	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}
	topicID, err := uuid.Parse(c.Param("topic_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid topic_id format")
	}

	if err := h.Services(c).Newsletter.DeleteNewsletterTopic(c.Request().Context(), workspaceID, repoID, topicID); err != nil {
		return h.newsletterTopicError(c, err, "delete")
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *NewsletterTopicHandler) newsletterTopicError(c echo.Context, err error, action string) error {
	if errors.Is(err, repository_service.ErrRepositoryNotFound) {
		return responses.ErrorResponse(c, http.StatusNotFound, "repository not found")
	}
	if errors.Is(err, newsletter_service.ErrTopicNotFound) {
		return responses.ErrorResponse(c, http.StatusNotFound, "topic not found")
	}
	if errors.Is(err, newsletter_service.ErrTopicAlreadySent) {
		return responses.ErrorResponse(c, http.StatusConflict, "cannot "+action+" a topic that has already been sent")
	}
	return h.internalError(c, err)
}
