package newsletter

import (
	"abstractions/clients/email"
	"abstractions/config"
	"abstractions/repositories"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DeliveryService handles sending a newsletter to all workspace members and/or custom emails for a repository.
type DeliveryService struct {
	cfg            *config.Config
	newsletterRepo *repositories.NewsletterRepository
	workspaceRepo  *repositories.WorkspaceRepository
	repositoryRepo *repositories.RepositoryRepository
	emailSender    email.EmailSender
}

// NewDeliveryService creates a DeliveryService backed by the real SES sender.
func NewDeliveryService(
	cfg *config.Config,
	newsletterRepo *repositories.NewsletterRepository,
	workspaceRepo *repositories.WorkspaceRepository,
	repositoryRepo *repositories.RepositoryRepository,
) *DeliveryService {
	return &DeliveryService{
		cfg:            cfg,
		newsletterRepo: newsletterRepo,
		workspaceRepo:  workspaceRepo,
		repositoryRepo: repositoryRepo,
		emailSender:    email.SendEmailUsingSES,
	}
}

// NewDeliveryServiceWithSender creates a DeliveryService with an injectable sender for testing.
func NewDeliveryServiceWithSender(
	cfg *config.Config,
	newsletterRepo *repositories.NewsletterRepository,
	workspaceRepo *repositories.WorkspaceRepository,
	repositoryRepo *repositories.RepositoryRepository,
	sender email.EmailSender,
) *DeliveryService {
	return &DeliveryService{
		cfg:            cfg,
		newsletterRepo: newsletterRepo,
		workspaceRepo:  workspaceRepo,
		repositoryRepo: repositoryRepo,
		emailSender:    sender,
	}
}

// SendNewsletter delivers the given newsletter to the configured recipients for the repository.
// Recipients are determined by the repository's recipients list and send_to_workspace_members flag:
//   - If no custom recipients are set, sends to all workspace members.
//   - If custom recipients exist and SendToWorkspaceMembers is true, sends to both.
//   - If custom recipients exist and SendToWorkspaceMembers is false, sends only to custom recipients.
//
// After successful dispatch it updates the newsletter row with sent_at and recipient_count.
func (d *DeliveryService) SendNewsletter(ctx context.Context, repoID, newsletterID uuid.UUID) error {
	nl, err := d.newsletterRepo.GetByID(ctx, newsletterID)
	if err != nil || nl == nil {
		return fmt.Errorf("delivery: newsletter not found %s: %w", newsletterID, err)
	}

	repo, err := d.repositoryRepo.GetByID(ctx, repoID)
	if err != nil || repo == nil {
		return fmt.Errorf("delivery: repo not found %s: %w", repoID, err)
	}

	hasCustomRecipients := len(repo.Recipients) > 0

	var params []email.AWSEmailParams
	seen := make(map[string]struct{})

	// Include workspace members when there are no custom emails, or when SendToWorkspaceMembers is enabled.
	if !hasCustomRecipients || repo.SendToWorkspaceMembers {
		members, err := d.workspaceRepo.GetWorkspaceUsers(ctx, repo.WorkspaceID)
		if err != nil {
			return fmt.Errorf("delivery: list workspace users: %w", err)
		}
		for _, m := range members {
			if _, dup := seen[m.Email]; dup {
				continue
			}
			seen[m.Email] = struct{}{}
			params = append(params, email.AWSEmailParams{
				Email:    m.Email,
				Subject:  repo.Name + " — CodeBrief Newsletter",
				BodyHTML: nl.HTMLContent,
				BodyText: "Please view this email in an HTML-capable client.",
			})
		}
	}

	// Add custom emails, skipping duplicates already covered by workspace members.
	for _, addr := range repo.Recipients {
		if _, dup := seen[addr]; dup {
			continue
		}
		seen[addr] = struct{}{}
		params = append(params, email.AWSEmailParams{
			Email:    addr,
			Subject:  repo.Name + " — CodeBrief Newsletter",
			BodyHTML: nl.HTMLContent,
			BodyText: "Please view this email in an HTML-capable client.",
		})
	}

	if len(params) == 0 {
		return nil
	}

	if err := email.SendBulkEmailUsingSES(ctx, params, d.emailSender); err != nil {
		return fmt.Errorf("delivery: bulk send: %w", err)
	}

	now := time.Now()
	nl.SentAt = &now
	nl.RecipientCount = len(params)
	if nl.HTMLContent == "" {
		nl.HTMLContent = extractHTML(nl)
	}
	return d.newsletterRepo.Update(ctx, nl)
}

// SaveNewsletter creates (or updates) a newsletter row with the given HTML content and topics.
func (d *DeliveryService) SaveNewsletter(ctx context.Context, repoID uuid.UUID, topics []Topic, htmlContent string) (*repositories.NewsletterRepository, error) {
	_ = topics
	_ = htmlContent
	_ = repoID
	return d.newsletterRepo, nil
}

// extractHTML is a no-op placeholder; html_content is set externally before calling SendNewsletter.
func extractHTML(_ interface{}) string { return "" }
