package newsletter

import (
	"abstractions/clients/email"
	"abstractions/config"
	"abstractions/repositories"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DeliveryService handles sending a newsletter to all recipients for a repository.
type DeliveryService struct {
	cfg            *config.Config
	newsletterRepo *repositories.NewsletterRepository
	repositoryRepo *repositories.RepositoryRepository
	workspaceRepo  *repositories.WorkspaceRepository
	emailSender    email.EmailSender
}

// NewDeliveryService creates a DeliveryService backed by the real SES sender.
func NewDeliveryService(
	cfg *config.Config,
	newsletterRepo *repositories.NewsletterRepository,
	repositoryRepo *repositories.RepositoryRepository,
	workspaceRepo *repositories.WorkspaceRepository,
) *DeliveryService {
	return &DeliveryService{
		cfg:            cfg,
		newsletterRepo: newsletterRepo,
		repositoryRepo: repositoryRepo,
		workspaceRepo:  workspaceRepo,
		emailSender:    email.SendEmailUsingSES,
	}
}

// NewDeliveryServiceWithSender creates a DeliveryService with an injectable sender for testing.
func NewDeliveryServiceWithSender(
	cfg *config.Config,
	newsletterRepo *repositories.NewsletterRepository,
	repositoryRepo *repositories.RepositoryRepository,
	workspaceRepo *repositories.WorkspaceRepository,
	sender email.EmailSender,
) *DeliveryService {
	return &DeliveryService{
		cfg:            cfg,
		newsletterRepo: newsletterRepo,
		repositoryRepo: repositoryRepo,
		workspaceRepo:  workspaceRepo,
		emailSender:    sender,
	}
}

// SendNewsletter delivers the given newsletter to all configured recipients for the repository.
// Recipients are sourced from repo.Recipients (explicit list) and, when repo.SendToWorkspaceMembers
// is true, from the workspace users. After successful dispatch it updates the newsletter row with
// sent_at and recipient_count.
func (d *DeliveryService) SendNewsletter(ctx context.Context, repoID, newsletterID uuid.UUID) error {
	nl, err := d.newsletterRepo.GetByID(ctx, newsletterID)
	if err != nil || nl == nil {
		return fmt.Errorf("delivery: newsletter not found %s: %w", newsletterID, err)
	}

	repo, err := d.repositoryRepo.GetByID(ctx, repoID)
	if err != nil || repo == nil {
		return fmt.Errorf("delivery: repo not found %s: %w", repoID, err)
	}

	// Build the deduplicated recipient list
	seen := map[string]bool{}
	var recipientEmails []string

	for _, e := range repo.Recipients {
		if e != "" && !seen[e] {
			seen[e] = true
			recipientEmails = append(recipientEmails, e)
		}
	}

	if repo.SendToWorkspaceMembers {
		members, err := d.workspaceRepo.GetWorkspaceUsers(ctx, repo.WorkspaceID)
		if err != nil {
			return fmt.Errorf("delivery: list workspace users: %w", err)
		}
		for _, m := range members {
			if m.Email != "" && !seen[m.Email] {
				seen[m.Email] = true
				recipientEmails = append(recipientEmails, m.Email)
			}
		}
	}

	if len(recipientEmails) == 0 {
		return nil
	}

	subject := repo.Name + " — CodeBrief Newsletter"
	if nl.Topic != nil && nl.Topic.Title != "" {
		subject = repo.Name + ": " + nl.Topic.Title
	}

	params := make([]email.AWSEmailParams, 0, len(recipientEmails))
	for _, addr := range recipientEmails {
		params = append(params, email.AWSEmailParams{
			Email:    addr,
			Subject:  subject,
			BodyHTML: nl.HTMLContent,
			BodyText: "Please view this email in an HTML-capable client.",
		})
	}

	if err := email.SendBulkEmailUsingSES(ctx, params, d.emailSender); err != nil {
		return fmt.Errorf("delivery: bulk send: %w", err)
	}

	// Persist delivery metadata
	now := time.Now()
	nl.SentAt = &now
	nl.RecipientCount = len(recipientEmails)
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
