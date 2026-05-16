package responses

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"encoding/json"
)

// NewsletterTopicResponse is the DTO for a single newsletter topic.
type NewsletterTopicResponse struct {
	ID                string   `json:"id"`
	RepoID            string   `json:"repoId"`
	Title             string   `json:"title"`
	Description       string   `json:"description"`
	RelatedFiles      []string `json:"relatedFiles"`
	Position          int      `json:"position"`
	SentAt            *string  `json:"sentAt"`
	RequestedByUserID *string  `json:"requestedByUserId"`
	RequestedByName   *string  `json:"requestedByName"`
	RequestedByAvatar *string  `json:"requestedByAvatar"`
	NewsletterID      *string  `json:"newsletterId"`
	CreatedAt         string   `json:"createdAt"`
}

func newTopicResponse(t *models.NewsletterTopic, newsletterID *string, creatorName *string, creatorAvatar *string) NewsletterTopicResponse {
	var relatedFiles []string
	if t.RelatedFiles != nil {
		_ = json.Unmarshal(t.RelatedFiles, &relatedFiles)
	}
	if relatedFiles == nil {
		relatedFiles = []string{}
	}

	var requestedByUserID *string
	if t.CreatedBy != nil {
		s := t.CreatedBy.String()
		requestedByUserID = &s
	}

	return NewsletterTopicResponse{
		ID:                t.ID.String(),
		RepoID:            t.RepoID.String(),
		Title:             t.Title,
		Description:       t.Description,
		RelatedFiles:      relatedFiles,
		Position:          t.Position,
		SentAt:            fmtTime(t.SentAt),
		RequestedByUserID: requestedByUserID,
		RequestedByName:   creatorName,
		RequestedByAvatar: creatorAvatar,
		NewsletterID:      newsletterID,
		CreatedAt:         t.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
}

// NewNewsletterTopicResponse builds a response from a plain topic model (no newsletter join or user info).
func NewNewsletterTopicResponse(t *models.NewsletterTopic) NewsletterTopicResponse {
	return newTopicResponse(t, nil, nil, nil)
}

// NewNewsletterTopicResponseWithNewsletter builds a response from an enriched topic that includes
// the linked newsletter ID and creator info.
func NewNewsletterTopicResponseWithNewsletter(t *repositories.NewsletterTopicWithNewsletterID) NewsletterTopicResponse {
	var newsletterID *string
	if t.NewsletterID != nil {
		s := t.NewsletterID.String()
		newsletterID = &s
	}
	return newTopicResponse(&t.NewsletterTopic, newsletterID, t.CreatorName, t.CreatorAvatar)
}

func NewNewsletterTopicResponses(topics []*repositories.NewsletterTopicWithNewsletterID) []NewsletterTopicResponse {
	result := make([]NewsletterTopicResponse, len(topics))
	for i, t := range topics {
		result[i] = NewNewsletterTopicResponseWithNewsletter(t)
	}
	return result
}
