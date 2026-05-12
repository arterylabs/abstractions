package jobs

import (
	ai_client "abstractions/clients/ai"
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
}

func NewGenerateNewsletterJob(db *gorm.DB, cfg *config.Config) *GenerateNewsletterJob {
	return &GenerateNewsletterJob{
		db:             db,
		config:         cfg,
		repositoryRepo: repositories.NewRepositoryRepository(db),
		newsletterRepo: repositories.NewNewsletterRepository(db),
		topicRepo:      repositories.NewNewsletterTopicRepository(db),
	}
}

func (j *GenerateNewsletterJob) ID() string {
	return GenerateNewsletterJobID
}

// NewGenerateNewsletterTask builds the asynq task for the given repository ID.
func NewGenerateNewsletterTask(repoID uuid.UUID) (*asynq.Task, error) {
	payload := GenerateNewsletterPayload{RepositoryID: repoID.String()}
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

	logger.Info(ctx, "generate_newsletter: starting repoId="+repoID.String())

	repo, err := j.repositoryRepo.GetByID(ctx, repoID)
	if err != nil || repo == nil {
		return fmt.Errorf("generate_newsletter: repo not found: %w", err)
	}

	aiProvider := ai_client.NewAIProvider()
	topic, err := j.topicRepo.GetNextForRepo(ctx, repoID)
	if err != nil {
		return fmt.Errorf("generate_newsletter: get next topic: %w", err)
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
			return fmt.Errorf("generate_newsletter: topics: %w", err)
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
			return fmt.Errorf("generate_newsletter: create topic batch: %w", err)
		}

		topic, err = j.topicRepo.GetNextForRepo(ctx, repoID)
		if err != nil {
			return fmt.Errorf("generate_newsletter: refresh next topic: %w", err)
		}
	}

	if topic == nil {
		logger.Info(ctx, "generate_newsletter: no topics repoId="+repoID.String())
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
		return fmt.Errorf("generate_newsletter: content: %w", err)
	}

	assembler := &newsletter_service.Assembler{}
	innerHTML, err := assembler.Assemble(repo.Name, entries)
	if err != nil {
		return fmt.Errorf("generate_newsletter: assemble: %w", err)
	}

	subject := repo.Name + ": " + topic.Title
	fullHTML, _, err := email_client.BuildEmail(email_client.SendGenericEmailParams{
		Title:    subject,
		BodyHTML: innerHTML,
		Reason:   "You are receiving this because you are subscribed to the " + repo.Name + " engineering newsletter.",
	})
	if err != nil {
		return fmt.Errorf("generate_newsletter: build email: %w", err)
	}

	nl := &models.Newsletter{
		ID:          uuid.New(),
		RepoID:      repoID,
		TopicID:     &topic.ID,
		HTMLContent: fullHTML,
	}
	if err := j.newsletterRepo.Create(ctx, nl); err != nil {
		return fmt.Errorf("generate_newsletter: create newsletter: %w", err)
	}

	delivery := newsletter_service.NewDeliveryService(j.config, j.newsletterRepo, j.repositoryRepo, repositories.NewWorkspaceRepository(j.db))
	if err := delivery.SendNewsletter(ctx, repoID, nl.ID); err != nil {
		return fmt.Errorf("generate_newsletter: send: %w", err)
	}
	if err := j.topicRepo.MarkSent(ctx, topic.ID); err != nil {
		logger.Info(ctx, fmt.Sprintf("generate_newsletter: mark topic sent error: %v", err))
	}

	now := time.Now()
	repo.NewsletterSentAt = &now
	if err := j.repositoryRepo.Update(ctx, repo); err != nil {
		logger.Info(ctx, fmt.Sprintf("generate_newsletter: update newsletter_sent_at error: %v", err))
	}

	logger.Info(ctx, "generate_newsletter: done repoId="+repoID.String())
	return nil
}
