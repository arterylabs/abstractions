package jobs

import (
	app_notifications "abstractions/clients/app_notifications"
	email_client "abstractions/clients/email"
	"abstractions/config"
	"abstractions/db/models"
	"abstractions/logger"
	"abstractions/repositories"
	newsletter_service "abstractions/services/newsletter"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

const GenerateNewsletterJobID = "generate_newsletter"

// GenerateNewsletterPayload is the JSON payload for a GenerateNewsletterJob task.
type GenerateNewsletterPayload struct {
	RepositoryID string `json:"repositoryId"`
	// NewsletterID is the ID of the newsletter record (status=queued) created before enqueueing.
	NewsletterID string `json:"newsletterId"`
}

// GenerateNewsletterJob runs the full end-to-end newsletter pipeline for a single repository:
//  1. Pull the next queued topic (or backfill the queue if empty)
//  2. Generate newsletter content for that single topic
//  3. Assemble and persist the newsletter row
//  4. Send the newsletter and stamp the topic as sent
//  5. Update repository.newsletter_sent_at on success
type GenerateNewsletterJob struct {
	db             *gorm.DB
	config         *config.Config
	repositoryRepo *repositories.RepositoryRepository
	newsletterRepo *repositories.NewsletterRepository
	topicRepo      *repositories.NewsletterTopicRepository
	userRepo       *repositories.UserRepository
}

func NewGenerateNewsletterJob(db *gorm.DB, cfg *config.Config) *GenerateNewsletterJob {
	return &GenerateNewsletterJob{
		db:             db,
		config:         cfg,
		repositoryRepo: repositories.NewRepositoryRepository(db),
		newsletterRepo: repositories.NewNewsletterRepository(db),
		topicRepo:      repositories.NewNewsletterTopicRepository(db),
		userRepo:       repositories.NewUserRepository(db),
	}
}

func (j *GenerateNewsletterJob) ID() string {
	return GenerateNewsletterJobID
}

// NewGenerateNewsletterTask builds the asynq task for the given repository and newsletter IDs.
func NewGenerateNewsletterTask(repoID, newsletterID uuid.UUID) (*asynq.Task, error) {
	payload := GenerateNewsletterPayload{RepositoryID: repoID.String(), NewsletterID: newsletterID.String()}
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(GenerateNewsletterJobID, b), nil
}

// Handler processes a GenerateNewsletterJob task.
func (j *GenerateNewsletterJob) Handler(ctx context.Context, t *asynq.Task) error {
	var payload GenerateNewsletterPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("generate_newsletter: decode payload: %w", err)
	}

	repoID, err := uuid.Parse(payload.RepositoryID)
	if err != nil {
		return fmt.Errorf("generate_newsletter: invalid repo id: %w", err)
	}

	newsletterID, err := uuid.Parse(payload.NewsletterID)
	if err != nil {
		return fmt.Errorf("generate_newsletter: invalid newsletter id: %w", err)
	}

	ctx = logger.SetMetadata(ctx, map[string]interface{}{
		"repoId":       repoID.String(),
		"newsletterId": newsletterID.String(),
	})

	logger.Info(ctx, "generate_newsletter: starting")

	// Load the newsletter record that was created when the job was enqueued.
	nl, err := j.newsletterRepo.GetByID(ctx, newsletterID)
	if err != nil {
		return fmt.Errorf("generate_newsletter: load newsletter: %w", err)
	}
	if nl == nil {
		return fmt.Errorf("generate_newsletter: newsletter %s not found", newsletterID)
	}

	// Transition to processing so the UI reflects activity.
	nl.Status = models.NewsletterStatusProcessing
	if err := j.newsletterRepo.Update(ctx, nl); err != nil {
		return fmt.Errorf("generate_newsletter: set processing: %w", err)
	}

	// repoName and workspaceName are updated once the repo is loaded so that markError
	// can include them in the Slack alert regardless of where in the pipeline it fires.
	repoName := repoID.String()
	workspaceName := ""

	// markError stamps the newsletter with error status and the error message, then returns the error.
	markError := func(wrapped error) error {
		nl.Status = models.NewsletterStatusError
		msg := wrapped.Error()
		nl.ErrorMessage = &msg
		_ = j.newsletterRepo.Update(ctx, nl)
		go app_notifications.SendNewsletterFailedNotification(context.Background(), app_notifications.NewsletterFailedParams{
			RepoName:      repoName,
			WorkspaceName: workspaceName,
			ErrorMessage:  msg,
		})
		return wrapped
	}

	repoCtx, err := LoadRepositoryContext(ctx, repoID, j.repositoryRepo)
	if err != nil {
		return markError(fmt.Errorf("generate_newsletter: %w", err))
	}
	repo := repoCtx.Repository
	repoName = repo.Name
	workspaceName = repo.WorkspaceID.String()

	ctx = logger.SetMetadata(ctx, map[string]interface{}{
		"workspaceId": repo.WorkspaceID.String(),
		"repoName":    repo.Name,
	})

	// Resolve the name of the user who configured this repository for the "why receiving" footer.
	reason := buildNewsletterReason(ctx, j.userRepo, repo)

	aiProvider, err := newWorkspaceAIProvider(ctx, j.db, repo.WorkspaceID)
	if err != nil {
		logger.Error(ctx, err)
		return markError(fmt.Errorf("generate_newsletter: %w", err))
	}
	// Use the pre-assigned topic when available; otherwise fall back to the next queued topic.
	var topic *models.NewsletterTopic
	if nl.TopicID != nil {
		topic, err = j.topicRepo.GetByID(ctx, *nl.TopicID)
		if err != nil {
			return markError(fmt.Errorf("generate_newsletter: load pre-assigned topic: %w", err))
		}
	} else {
		topic, err = j.topicRepo.GetNextForRepo(ctx, repoID)
		if err != nil {
			return markError(fmt.Errorf("generate_newsletter: get next topic: %w", err))
		}
	}

	if topic == nil {
		topicSvc := newsletter_service.NewTopicGenerator(
			aiProvider,
			repositories.NewRepositoryFileRepository(j.db),
			repositories.NewChunkEmbeddingRepository(j.db),
			j.newsletterRepo,
			j.repositoryRepo,
		)
		generatedTopics, err := topicSvc.GenerateTopics(ctx, repoID)
		if err != nil {
			return markError(fmt.Errorf("generate_newsletter: topics: %w", err))
		}

		topicModels := make([]*models.NewsletterTopic, 0, len(generatedTopics))
		for _, generatedTopic := range generatedTopics {
			relatedFilesJSON, err := json.Marshal(generatedTopic.RelatedFiles)
			if err != nil {
				continue
			}
			topicModels = append(topicModels, &models.NewsletterTopic{
				RepoID:       repoID,
				Title:        generatedTopic.Title,
				Description:  generatedTopic.Description,
				RelatedFiles: models.JSONBRaw(relatedFilesJSON),
			})
		}
		if err := j.topicRepo.CreateBatch(ctx, topicModels); err != nil {
			return markError(fmt.Errorf("generate_newsletter: create topic batch: %w", err))
		}

		topic, err = j.topicRepo.GetNextForRepo(ctx, repoID)
		if err != nil {
			return markError(fmt.Errorf("generate_newsletter: refresh next topic: %w", err))
		}
	}

	if topic == nil {
		logger.Info(ctx, "generate_newsletter: no topics available")
		nl.Status = models.NewsletterStatusError
		noTopicsMsg := "no topics available to generate newsletter"
		nl.ErrorMessage = &noTopicsMsg
		_ = j.newsletterRepo.Update(ctx, nl)
		return nil
	}

	var relatedFiles []string
	if len(topic.RelatedFiles) > 0 {
		_ = json.Unmarshal(topic.RelatedFiles, &relatedFiles)
	}

	contentSvc := newsletter_service.NewContentGenerator(
		aiProvider,
		repositories.NewChunkEmbeddingRepository(j.db),
		j.repositoryRepo,
	)
	entries, err := contentSvc.GenerateNewsletter(ctx, repoID, []newsletter_service.Topic{{
		Title:        topic.Title,
		Description:  topic.Description,
		RelatedFiles: relatedFiles,
	}})
	if err != nil {
		return markError(fmt.Errorf("generate_newsletter: content: %w", err))
	}

	repoURL := "https://github.com/" + repo.Name
	assembler := &newsletter_service.Assembler{}
	innerHTML, err := assembler.Assemble(repo.Name, repoURL, relatedFiles, repo.IndexedAt, entries)
	if err != nil {
		return markError(fmt.Errorf("generate_newsletter: assemble: %w", err))
	}

	subject := topic.Title
	fullHTML, _, err := email_client.BuildEmail(email_client.SendGenericEmailParams{
		Title:    subject,
		BodyHTML: innerHTML,
		Reason:   reason,
	})
	if err != nil {
		return markError(fmt.Errorf("generate_newsletter: build email: %w", err))
	}

	// Populate the pre-created newsletter record with the generated content.
	nl.TopicID = &topic.ID
	nl.HTMLContent = fullHTML
	if err := j.newsletterRepo.Update(ctx, nl); err != nil {
		return markError(fmt.Errorf("generate_newsletter: update newsletter content: %w", err))
	}

	delivery := newsletter_service.NewDeliveryService(j.config, j.newsletterRepo, repositories.NewWorkspaceRepository(j.db), j.repositoryRepo)
	if err := delivery.SendNewsletter(ctx, repoID, nl.ID, subject); err != nil {
		return markError(fmt.Errorf("generate_newsletter: send: %w", err))
	}
	if err := j.topicRepo.MarkSent(ctx, topic.ID); err != nil {
		logger.Info(ctx, fmt.Sprintf("generate_newsletter: mark topic sent error: %v", err))
	}

	now := time.Now()
	repo.NewsletterSentAt = &now
	if err := j.repositoryRepo.Update(ctx, repo); err != nil {
		logger.Info(ctx, fmt.Sprintf("generate_newsletter: update newsletter_sent_at error: %v", err))
	}

	logger.Info(ctx, "generate_newsletter: done")
	return nil
}

// buildNewsletterReason returns the footer "why you're receiving this" message.
// It looks up the user who created the repository; falls back to a generic message on any error.
func buildNewsletterReason(ctx context.Context, userRepo *repositories.UserRepository, repo *models.Repository) string {
	fallback := `You're receiving this because someone configured these newsletters on <a href="https://useabstractions.com">Abstractions</a>.`
	creator, err := userRepo.GetByID(ctx, repo.CreatedBy)
	if err != nil || creator == nil {
		return fallback
	}
	return fmt.Sprintf(`You're receiving this because %s configured these newsletters on <a href="https://useabstractions.com">Abstractions</a>.`, creator.Name)
}
