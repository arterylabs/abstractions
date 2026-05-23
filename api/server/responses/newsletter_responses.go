package responses

import (
	"abstractions/db/models"
	"time"
)

// NewsletterSummaryResponse is returned in the list endpoint (no HTML body).
type NewsletterSummaryResponse struct {
	ID             string  `json:"id"`
	RepoID         string  `json:"repoId"`
	Subject        string  `json:"subject"`
	SentAt         *string `json:"sentAt"`
	RecipientCount int     `json:"recipientCount"`
	CreatedAt      string  `json:"createdAt"`
}

// NewsletterDetailResponse is returned by the single-newsletter endpoint.
type NewsletterDetailResponse struct {
	ID             string                    `json:"id"`
	RepoID         string                    `json:"repoId"`
	TopicID        *string                   `json:"topicId"`
	Topic          *NewsletterTopicResponse  `json:"topic"`
	Topics         []NewsletterTopicResponse `json:"topics"`
	HTMLContent    string                    `json:"htmlContent"`
	SentAt         *string                   `json:"sentAt"`
	RecipientCount int                       `json:"recipientCount"`
	CreatedAt      string                    `json:"createdAt"`
}

func fmtTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.UTC().Format(time.RFC3339)
	return &s
}

func firstTopicTitle(topic *models.NewsletterTopic) string {
	if topic == nil {
		return ""
	}
	return topic.Title
}

func newsletterTopics(topic *models.NewsletterTopic) []NewsletterTopicResponse {
	if topic == nil {
		return []NewsletterTopicResponse{}
	}
	return []NewsletterTopicResponse{NewNewsletterTopicResponse(topic)}
}

func NewNewsletterSummaryResponse(n *models.Newsletter) NewsletterSummaryResponse {
	return NewsletterSummaryResponse{
		ID:             n.ID.String(),
		RepoID:         n.RepoID.String(),
		Subject:        firstTopicTitle(n.Topic),
		SentAt:         fmtTime(n.SentAt),
		RecipientCount: n.RecipientCount,
		CreatedAt:      n.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func NewNewsletterSummaryResponses(newsletters []*models.Newsletter) []NewsletterSummaryResponse {
	result := make([]NewsletterSummaryResponse, len(newsletters))
	for i, n := range newsletters {
		result[i] = NewNewsletterSummaryResponse(n)
	}
	return result
}

func NewNewsletterDetailResponse(n *models.Newsletter) NewsletterDetailResponse {
	var topicID *string
	if n.TopicID != nil {
		s := n.TopicID.String()
		topicID = &s
	}
	topics := newsletterTopics(n.Topic)
	var topic *NewsletterTopicResponse
	if len(topics) > 0 {
		topic = &topics[0]
	}
	return NewsletterDetailResponse{
		ID:             n.ID.String(),
		RepoID:         n.RepoID.String(),
		TopicID:        topicID,
		Topic:          topic,
		Topics:         topics,
		HTMLContent:    n.HTMLContent,
		SentAt:         fmtTime(n.SentAt),
		RecipientCount: n.RecipientCount,
		CreatedAt:      n.CreatedAt.UTC().Format(time.RFC3339),
	}
}
