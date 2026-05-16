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

const GenerateTopicsJobID = "generate_topics"

type GenerateTopicsPayload struct {
	RepositoryID string `json:"repositoryId"`
}

type GenerateTopicsJob struct {
	db     *gorm.DB
	config *config.Config
}

func NewGenerateTopicsJob(db *gorm.DB, cfg *config.Config) *GenerateTopicsJob {
	return &GenerateTopicsJob{db: db, config: cfg}
}

func (j *GenerateTopicsJob) ID() string {
	return GenerateTopicsJobID
}

func NewGenerateTopicsTask(repoID uuid.UUID) (*asynq.Task, error) {
	payload := GenerateTopicsPayload{RepositoryID: repoID.String()}
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(GenerateTopicsJobID, b), nil
}

// Handler generates topics for a repo and inserts them into the newsletter_topics queue.
// Skips topics with duplicate (repo_id, title) to avoid re-runs creating duplicates.
func (j *GenerateTopicsJob) Handler(ctx context.Context, t *asynq.Task) error {
	var payload GenerateTopicsPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("generate_topics: decode payload: %w", err)
	}

	repoID, err := uuid.Parse(payload.RepositoryID)
	if err != nil {
		return fmt.Errorf("generate_topics: invalid repo id: %w", err)
	}

	logger.Info(ctx, "generate_topics: starting repoId="+repoID.String())

	if _, err := LoadRepositoryContext(ctx, repoID, repositories.NewRepositoryRepository(j.db)); err != nil {
		return fmt.Errorf("generate_topics: %w", err)
	}

	aiProvider := ai_client.NewAIProvider()
	topicSvc := newsletter_service.NewTopicGenerator(
		aiProvider,
		repositories.NewRepositoryFileRepository(j.db),
		repositories.NewChunkEmbeddingRepository(j.db),
		repositories.NewNewsletterRepository(j.db),
		repositories.NewRepositoryRepository(j.db),
	)

	topics, err := topicSvc.GenerateTopics(ctx, repoID)
	if err != nil {
		return fmt.Errorf("generate_topics: generate: %w", err)
	}
	if len(topics) == 0 {
		logger.Info(ctx, "generate_topics: no topics repoId="+repoID.String())
		return nil
	}

	topicRepo := repositories.NewNewsletterTopicRepository(j.db)
	topicModels := make([]*models.NewsletterTopic, 0, len(topics))
	for _, tp := range topics {
		relatedFilesJSON, err := json.Marshal(tp.RelatedFiles)
		if err != nil {
			continue
		}
		topicModels = append(topicModels, &models.NewsletterTopic{
			RepoID:       repoID,
			Title:        tp.Title,
			Description:  tp.Description,
			RelatedFiles: models.JSONBRaw(relatedFilesJSON),
		})
	}

	if err := topicRepo.CreateBatch(ctx, topicModels); err != nil {
		return fmt.Errorf("generate_topics: insert batch: %w", err)
	}

	logger.Info(ctx, fmt.Sprintf("generate_topics: inserted %d topics repoId=%s", len(topicModels), repoID.String()))
	return nil
}
