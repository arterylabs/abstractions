package jobs

import (
	ai_client "abstractions/clients/ai"
	"abstractions/config"
	"abstractions/db/models"
	"abstractions/logger"
	"abstractions/repositories"
	newsletter_service "abstractions/services/newsletter"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

const GenerateNewsletterDraftJobID = "generate_newsletter_draft"

// GenerateNewsletterDraftJob runs steps 1–4 of the newsletter pipeline and persists a
// draft newsletter row (SentAt = nil) without delivering it:
//  1. GenerateTopics  — finds meaningful topics from the indexed file summaries
//  2. GenerateNewsletter — produces per-topic content entries
//  3. Assemble — renders the HTML email body
//  4. Creates a draft newsletter row (SentAt = nil) — no delivery
type GenerateNewsletterDraftJob struct {
	db             *gorm.DB
	config         *config.Config
	repositoryRepo *repositories.RepositoryRepository
	newsletterRepo *repositories.NewsletterRepository
	topicRepo      *repositories.NewsletterTopicRepository
}

func NewGenerateNewsletterDraftJob(db *gorm.DB, cfg *config.Config) *GenerateNewsletterDraftJob {
	return &GenerateNewsletterDraftJob{
		db:             db,
		config:         cfg,
		repositoryRepo: repositories.NewRepositoryRepository(db),
		newsletterRepo: repositories.NewNewsletterRepository(db),
		topicRepo:      repositories.NewNewsletterTopicRepository(db),
	}
}

func (j *GenerateNewsletterDraftJob) ID() string {
	return GenerateNewsletterDraftJobID
}

// NewGenerateNewsletterDraftTask builds the asynq task for the given repository ID.
func NewGenerateNewsletterDraftTask(repoID uuid.UUID) (*asynq.Task, error) {
	payload := GenerateNewsletterPayload{RepositoryID: repoID.String()}
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(GenerateNewsletterDraftJobID, b), nil
}

// Handler processes a GenerateNewsletterDraftJob task.
func (j *GenerateNewsletterDraftJob) Handler(ctx context.Context, t *asynq.Task) error {
	var payload GenerateNewsletterPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("generate_newsletter_draft: decode payload: %w", err)
	}

	repoID, err := uuid.Parse(payload.RepositoryID)
	if err != nil {
		return fmt.Errorf("generate_newsletter_draft: invalid repo id: %w", err)
	}

	logger.Info(ctx, "generate_newsletter_draft: starting repoId="+repoID.String())

	repo, err := j.repositoryRepo.GetByID(ctx, repoID)
	if err != nil || repo == nil {
		return fmt.Errorf("generate_newsletter_draft: repo not found: %w", err)
	}

	// Step 1 — generate topics
	aiProvider := ai_client.NewAIProvider()
	topicSvc := newsletter_service.NewTopicGenerator(
		aiProvider,
		repositories.NewRepositoryFileRepository(j.db),
		j.newsletterRepo,
		j.repositoryRepo,
	)
	topics, err := topicSvc.GenerateTopics(ctx, repoID)
	if err != nil {
		return fmt.Errorf("generate_newsletter_draft: topics: %w", err)
	}
	if len(topics) == 0 {
		logger.Info(ctx, "generate_newsletter_draft: no topics available repoId="+repoID.String())
		return nil
	}

	// Step 2 — generate content
	contentSvc := newsletter_service.NewContentGenerator(
		aiProvider,
		repositories.NewChunkEmbeddingRepository(j.db),
		j.repositoryRepo,
	)
	entries, err := contentSvc.GenerateNewsletter(ctx, repoID, topics)
	if err != nil {
		return fmt.Errorf("generate_newsletter_draft: content: %w", err)
	}

	// Step 3 — assemble HTML
	assembler := &newsletter_service.Assembler{}
	html, err := assembler.Assemble(repo.Name, entries, "")
	if err != nil {
		return fmt.Errorf("generate_newsletter_draft: assemble: %w", err)
	}

	// Step 4 — persist top-scored topic and draft newsletter row (SentAt remains nil)
	topTopic := topics[0]
	relatedFilesJSON, err := json.Marshal(topTopic.RelatedFiles)
	if err != nil {
		return fmt.Errorf("generate_newsletter_draft: marshal related_files: %w", err)
	}

	topic := &models.NewsletterTopic{
		ID:           uuid.New(),
		Title:        topTopic.Title,
		Description:  topTopic.Description,
		RelatedFiles: models.JSONBRaw(relatedFilesJSON),
	}
	if err := j.topicRepo.Create(ctx, topic); err != nil {
		return fmt.Errorf("generate_newsletter_draft: create topic: %w", err)
	}

	nl := &models.Newsletter{
		ID:          uuid.New(),
		RepoID:      repoID,
		TopicID:     &topic.ID,
		HTMLContent: html,
	}
	if err := j.newsletterRepo.Create(ctx, nl); err != nil {
		return fmt.Errorf("generate_newsletter_draft: create newsletter: %w", err)
	}

	logger.Info(ctx, "generate_newsletter_draft: done repoId="+repoID.String()+" newsletterId="+nl.ID.String())
	return nil
}
