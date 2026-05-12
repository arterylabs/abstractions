package handlers

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"abstractions/server"
	"abstractions/server/requests"
	"abstractions/server/responses"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type NewsletterTopicHandler struct {
	*BaseHandler
	topicRepo *repositories.NewsletterTopicRepository
}

func NewNewsletterTopicHandler(srv *server.Server) *NewsletterTopicHandler {
	return &NewsletterTopicHandler{
		BaseHandler: NewBaseHandler(srv),
		topicRepo:   repositories.NewNewsletterTopicRepository(srv.DB),
	}
}

// ListNewsletterTopics handles GET /repositories/:repo_id/newsletter-topics
func (h *NewsletterTopicHandler) ListNewsletterTopics(c echo.Context) error {
	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}

	topics, err := h.topicRepo.ListByRepoID(c.Request().Context(), repoID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, responses.NewNewsletterTopicResponses(topics))
}

// CreateNewsletterTopic handles POST /repositories/:repo_id/newsletter-topics
func (h *NewsletterTopicHandler) CreateNewsletterTopic(c echo.Context) error {
	repoID, err := uuid.Parse(c.Param("repo_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid repo_id format")
	}

	var req requests.CreateNewsletterTopicRequest
	if err := c.Bind(&req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	relatedFilesJSON, _ := json.Marshal(req.RelatedFiles)
	topic := &models.NewsletterTopic{
		RepoID:       repoID,
		Title:        req.Title,
		Description:  req.Description,
		RelatedFiles: models.JSONBRaw(relatedFilesJSON),
	}
	if err := h.topicRepo.Create(c.Request().Context(), topic); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusCreated, responses.NewNewsletterTopicResponse(topic))
}

// DeleteNewsletterTopic handles DELETE /repositories/:repo_id/newsletter-topics/:topic_id.
// Returns 409 if the topic has already been sent.
func (h *NewsletterTopicHandler) DeleteNewsletterTopic(c echo.Context) error {
	topicID, err := uuid.Parse(c.Param("topic_id"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid topic_id format")
	}

	if err := h.topicRepo.Delete(c.Request().Context(), topicID); err != nil {
		if errors.Is(err, repositories.ErrTopicAlreadySent) {
			return responses.ErrorResponse(c, http.StatusConflict, "topic has already been sent")
		}
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
