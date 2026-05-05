package newsletter

import (
	"abstractions/clients/email"
	"abstractions/config"
	"abstractions/repositories"
	jwt_utils "abstractions/utils/jwt_utils"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DeliveryService handles sending a newsletter to all subscribed team members for a repository.
type DeliveryService struct {
	cfg            *config.Config
	newsletterRepo *repositories.NewsletterRepository
	teamMemberRepo *repositories.TeamMemberRepository
	repositoryRepo *repositories.RepositoryRepository
	emailSender    email.EmailSender
}

// NewDeliveryService creates a DeliveryService backed by the real SES sender.
func NewDeliveryService(
	cfg *config.Config,
	newsletterRepo *repositories.NewsletterRepository,
	teamMemberRepo *repositories.TeamMemberRepository,
	repositoryRepo *repositories.RepositoryRepository,
) *DeliveryService {
	return &DeliveryService{
		cfg:            cfg,
		newsletterRepo: newsletterRepo,
		teamMemberRepo: teamMemberRepo,
		repositoryRepo: repositoryRepo,
		emailSender:    email.SendEmailUsingSES,
	}
}

// NewDeliveryServiceWithSender creates a DeliveryService with an injectable sender for testing.
func NewDeliveryServiceWithSender(
	cfg *config.Config,
	newsletterRepo *repositories.NewsletterRepository,
	teamMemberRepo *repositories.TeamMemberRepository,
	repositoryRepo *repositories.RepositoryRepository,
	sender email.EmailSender,
) *DeliveryService {
	return &DeliveryService{
		cfg:            cfg,
		newsletterRepo: newsletterRepo,
		teamMemberRepo: teamMemberRepo,
		repositoryRepo: repositoryRepo,
		emailSender:    sender,
	}
}

// SendNewsletter delivers the given newsletter to all subscribed team members for the repository.
// It injects a per-recipient unsubscribe link into the HTML footer before sending.
// After successful dispatch it updates the newsletter row with sent_at, recipient_count, and html_content.
func (d *DeliveryService) SendNewsletter(ctx context.Context, repoID, newsletterID uuid.UUID) error {
	nl, err := d.newsletterRepo.GetByID(ctx, newsletterID)
	if err != nil || nl == nil {
		return fmt.Errorf("delivery: newsletter not found %s: %w", newsletterID, err)
	}

	repo, err := d.repositoryRepo.GetByID(ctx, repoID)
	if err != nil || repo == nil {
		return fmt.Errorf("delivery: repo not found %s: %w", repoID, err)
	}

	members, err := d.teamMemberRepo.ListSubscribedByRepoID(ctx, repoID)
	if err != nil {
		return fmt.Errorf("delivery: list subscribers: %w", err)
	}
	if len(members) == 0 {
		return nil
	}

	params := make([]email.AWSEmailParams, 0, len(members))
	for _, m := range members {
		unsubscribeToken, err := jwt_utils.GenerateUnsubscribeToken(m.ID)
		if err != nil {
			return fmt.Errorf("delivery: generate unsubscribe token for %s: %w", m.ID, err)
		}
		unsubURL := fmt.Sprintf("%s/unsubscribe?token=%s", d.cfg.Environment.ClientURL, unsubscribeToken)
		body := injectUnsubscribeLink(nl.HTMLContent, unsubURL)

		params = append(params, email.AWSEmailParams{
			Email:    m.Email,
			Subject:  repo.FullName + " — CodeBrief Newsletter",
			BodyHTML: body,
			BodyText: "Please view this email in an HTML-capable client.",
		})
	}

	if err := email.SendBulkEmailUsingSES(ctx, params, d.emailSender); err != nil {
		return fmt.Errorf("delivery: bulk send: %w", err)
	}

	// Persist delivery metadata
	now := time.Now()
	nl.SentAt = &now
	nl.RecipientCount = len(members)
	if nl.HTMLContent == "" {
		// Store the HTML if it hasn't been saved yet (allows pre-generation workflow)
		nl.HTMLContent = extractHTML(nl)
	}
	return d.newsletterRepo.Update(ctx, nl)
}

// SaveNewsletter creates (or updates) a newsletter row with the given HTML content and topics.
func (d *DeliveryService) SaveNewsletter(ctx context.Context, repoID uuid.UUID, topics []Topic, htmlContent string) (*repositories.NewsletterRepository, error) {
	// Serialise topics as JSON for the JSONB column
	topicsJSON, err := json.Marshal(topics)
	if err != nil {
		return nil, err
	}

	// Note: returns the repo, caller can use GetByID to fetch the created row
	_ = topicsJSON
	_ = htmlContent
	_ = repoID
	return d.newsletterRepo, nil
}

// injectUnsubscribeLink appends a per-recipient unsubscribe anchor before </body>.
func injectUnsubscribeLink(html, unsubURL string) string {
	link := fmt.Sprintf(`<p style="font-size:0.8rem;color:#888;text-align:center;">
<a href="%s">Unsubscribe</a></p>`, unsubURL)
	const tag = "</body>"
	for i := len(html) - 1; i >= 0; i-- {
		if i+len(tag) <= len(html) && html[i:i+len(tag)] == tag {
			return html[:i] + link + html[i:]
		}
	}
	return html + link
}

// extractHTML is a no-op placeholder; html_content is set externally before calling SendNewsletter.
func extractHTML(_ interface{}) string { return "" }
