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
//  1. GenerateTopics  — finds meaningful topics from the indexed file summaries
//  2. GenerateNewsletter — produces per-topic content entries
//  3. Assemble — renders the HTML email
//  4. Create newsletter row in the database
//  5. Create newsletter_topics row with the top topic
//  6. SendNewsletter — delivers to all configured recipients
//  7. Updates repository.newsletter_sent_at on success
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
		return fmt.Errorf("generate_newsletter: topics: %w", err)
	}
	if len(topics) == 0 {
		logger.Info(ctx, "generate_newsletter: no topics repoId="+repoID.String())
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
		return fmt.Errorf("generate_newsletter: content: %w", err)
	}

	// Step 3 — assemble HTML
	assembler := &newsletter_service.Assembler{}
	html, err := assembler.Assemble(repo.Name, entries, "")
	if err != nil {
		return fmt.Errorf("generate_newsletter: assemble: %w", err)
	}

	// Step 4 — persist top-scored topic row
	// Use only the top-scored topic (topics are already sorted by recency score).
	topTopic := topics[0]
	relatedFilesJSON, err := json.Marshal(topTopic.RelatedFiles)
	if err != nil {
		return fmt.Errorf("generate_newsletter: marshal related_files: %w", err)
	}

	topic := &models.NewsletterTopic{
		ID:           uuid.New(),
		Title:        topTopic.Title,
		Description:  topTopic.Description,
		RelatedFiles: models.JSONBRaw(relatedFilesJSON),
	}
	if err := j.topicRepo.Create(ctx, topic); err != nil {
		return fmt.Errorf("generate_newsletter: create topic: %w", err)
	}

	// Step 5 — persist newsletter row linked by topic_id
	nl := &models.Newsletter{
		ID:          uuid.New(),
		RepoID:      repoID,
		TopicID:     &topic.ID,
		HTMLContent: html,
	}
	if err := j.newsletterRepo.Create(ctx, nl); err != nil {
		return fmt.Errorf("generate_newsletter: create newsletter: %w", err)
	}

	// Step 6 — deliver to recipients
	delivery := newsletter_service.NewDeliveryService(j.config, j.newsletterRepo, j.repositoryRepo, repositories.NewWorkspaceRepository(j.db))
	if err := delivery.SendNewsletter(ctx, repoID, nl.ID); err != nil {
		return fmt.Errorf("generate_newsletter: send: %w", err)
	}

	// Step 7 — stamp newsletter_sent_at
	now := time.Now()
	repo.NewsletterSentAt = &now
	if err := j.repositoryRepo.Update(ctx, repo); err != nil {
		logger.Info(ctx, fmt.Sprintf("generate_newsletter: update newsletter_sent_at error: %v", err))
	}

	logger.Info(ctx, "generate_newsletter: done repoId="+repoID.String())
	return nil
}
