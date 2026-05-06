package responses

import (
	"abstractions/db/models"
	"encoding/json"
	"time"
)

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
	ID             string          `json:"id"`
	RepoID         string          `json:"repoId"`
	Topics         json.RawMessage `json:"topics"`
	HTMLContent    string          `json:"htmlContent"`
	SentAt         *string         `json:"sentAt"`
	RecipientCount int             `json:"recipientCount"`
	OpenCount      int             `json:"openCount"`
	ClickCount     int             `json:"clickCount"`
	OpenRate       float64         `json:"openRate"`
	ClickRate      float64         `json:"clickRate"`
	CreatedAt      string          `json:"createdAt"`
}

// UpcomingNewsletterResponse describes the next scheduled send time and pre-computed topics.
type UpcomingNewsletterResponse struct {
	NextSendAt string          `json:"nextSendAt"`
	Topics     json.RawMessage `json:"topics"`
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
	topics := json.RawMessage(n.Topics)
	if topics == nil {
		topics = json.RawMessage("[]")
	}
	return NewsletterDetailResponse{
		ID:             n.ID.String(),
		RepoID:         n.RepoID.String(),
		Topics:         topics,
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
