package responses

import (
	"abstractions/db/models"
	"encoding/json"
	"time"
)

// TopicResponse is the topic metadata embedded in newsletter responses.
type TopicResponse struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	RelatedFiles []string `json:"relatedFiles"`
}

// NewsletterSummaryResponse is returned in the list endpoint (no HTML body).
type NewsletterSummaryResponse struct {
	ID             string  `json:"id"`
	RepoID         string  `json:"repoId"`
	SentAt         *string `json:"sentAt"`
	RecipientCount int     `json:"recipientCount"`
	OpenCount      int     `json:"openCount"`
	ClickCount     int     `json:"clickCount"`
	OpenRate       float64 `json:"openRate"`
	ClickRate      float64 `json:"clickRate"`
	CreatedAt      string  `json:"createdAt"`
}

// NewsletterDetailResponse is returned by the single-newsletter endpoint.
type NewsletterDetailResponse struct {
	ID             string         `json:"id"`
	RepoID         string         `json:"repoId"`
	Topic          *TopicResponse `json:"topic"`
	HTMLContent    string         `json:"htmlContent"`
	SentAt         *string        `json:"sentAt"`
	RecipientCount int            `json:"recipientCount"`
	OpenCount      int            `json:"openCount"`
	ClickCount     int            `json:"clickCount"`
	OpenRate       float64        `json:"openRate"`
	ClickRate      float64        `json:"clickRate"`
	CreatedAt      string         `json:"createdAt"`
}

// UpcomingNewsletterResponse describes the next scheduled send time and the most recent topic.
type UpcomingNewsletterResponse struct {
	NextSendAt string         `json:"nextSendAt"`
	Topic      *TopicResponse `json:"topic"`
}

func rateOf(count, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(count) / float64(total)
}

func fmtTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.UTC().Format(time.RFC3339)
	return &s
}

func newTopicResponse(t *models.NewsletterTopic) *TopicResponse {
	if t == nil {
		return nil
	}
	var relatedFiles []string
	if len(t.RelatedFiles) > 0 {
		_ = json.Unmarshal(t.RelatedFiles, &relatedFiles)
	}
	if relatedFiles == nil {
		relatedFiles = []string{}
	}
	return &TopicResponse{
		ID:           t.ID.String(),
		Title:        t.Title,
		Description:  t.Description,
		RelatedFiles: relatedFiles,
	}
}

func NewTopicResponse(t *models.NewsletterTopic) *TopicResponse {
	return newTopicResponse(t)
}

func NewNewsletterSummaryResponse(n *models.Newsletter) NewsletterSummaryResponse {
	return NewsletterSummaryResponse{
		ID:             n.ID.String(),
		RepoID:         n.RepoID.String(),
		SentAt:         fmtTime(n.SentAt),
		RecipientCount: n.RecipientCount,
		OpenCount:      n.OpenCount,
		ClickCount:     n.ClickCount,
		OpenRate:       rateOf(n.OpenCount, n.RecipientCount),
		ClickRate:      rateOf(n.ClickCount, n.RecipientCount),
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
	return NewsletterDetailResponse{
		ID:             n.ID.String(),
		RepoID:         n.RepoID.String(),
		Topic:          newTopicResponse(n.Topic),
		HTMLContent:    n.HTMLContent,
		SentAt:         fmtTime(n.SentAt),
		RecipientCount: n.RecipientCount,
		OpenCount:      n.OpenCount,
		ClickCount:     n.ClickCount,
		OpenRate:       rateOf(n.OpenCount, n.RecipientCount),
		ClickRate:      rateOf(n.ClickCount, n.RecipientCount),
		CreatedAt:      n.CreatedAt.UTC().Format(time.RFC3339),
	}
}
