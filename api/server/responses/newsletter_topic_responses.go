package responses

import (
	"abstractions/db/models"
	"encoding/json"
	"time"
)

type NewsletterTopicResponse struct {
	ID           string   `json:"id"`
	RepoID       string   `json:"repoId"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	RelatedFiles []string `json:"relatedFiles"`
	SentAt       *string  `json:"sentAt"`
	CreatedAt    string   `json:"createdAt"`
}

func NewNewsletterTopicResponse(t *models.NewsletterTopic) NewsletterTopicResponse {
	var relatedFiles []string
	if len(t.RelatedFiles) > 0 {
		_ = json.Unmarshal(t.RelatedFiles, &relatedFiles)
	}
	if relatedFiles == nil {
		relatedFiles = []string{}
	}

	var sentAt *string
	if t.SentAt != nil {
		s := t.SentAt.UTC().Format(time.RFC3339)
		sentAt = &s
	}

	return NewsletterTopicResponse{
		ID:           t.ID.String(),
		RepoID:       t.RepoID.String(),
		Title:        t.Title,
		Description:  t.Description,
		RelatedFiles: relatedFiles,
		SentAt:       sentAt,
		CreatedAt:    t.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func NewNewsletterTopicResponses(topics []*models.NewsletterTopic) []NewsletterTopicResponse {
	result := make([]NewsletterTopicResponse, len(topics))
	for i, t := range topics {
		result[i] = NewNewsletterTopicResponse(t)
	}
	return result
}
