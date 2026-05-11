package jobs

import (
	"abstractions/config"
	"abstractions/logger"
	"abstractions/repositories"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

const NewsletterSchedulerJobID = "newsletter_scheduler"

// NewsletterSchedulerPayload carries the scheduler invocation time so the handler
// can compute the matching send_day / send_hour without relying on wall clock drift.
type NewsletterSchedulerPayload struct {
	TriggerTime string `json:"triggerTime"` // RFC3339
}

// NewsletterSchedulerJob queries repos whose send time matches the current hour and
// enqueues a GenerateNewsletterJob for each qualifying repo.
// It is designed to be enqueued every hour by the asynq Scheduler.
type NewsletterSchedulerJob struct {
	db             *gorm.DB
	config         *config.Config
	enqueuer       TaskEnqueuer
	repositoryRepo *repositories.RepositoryRepository
}

func NewNewsletterSchedulerJob(db *gorm.DB, cfg *config.Config, enqueuer TaskEnqueuer) *NewsletterSchedulerJob {
	return &NewsletterSchedulerJob{
		db:             db,
		config:         cfg,
		enqueuer:       enqueuer,
		repositoryRepo: repositories.NewRepositoryRepository(db),
	}
}

func (j *NewsletterSchedulerJob) ID() string {
	return NewsletterSchedulerJobID
}

// NewNewsletterSchedulerTask builds the asynq task to be enqueued by the periodic scheduler.
func NewNewsletterSchedulerTask() (*asynq.Task, error) {
	payload := NewsletterSchedulerPayload{TriggerTime: time.Now().UTC().Format(time.RFC3339)}
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(NewsletterSchedulerJobID, b), nil
}

// Handler processes a NewsletterSchedulerJob task.
// It finds all indexed repos matching the current send_day + send_hour,
// applies the 20-hour double-send guard, and enqueues GenerateNewsletterJob for each.
func (j *NewsletterSchedulerJob) Handler(ctx context.Context, t *asynq.Task) error {
	var payload NewsletterSchedulerPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("newsletter_scheduler: decode payload: %w", err)
	}

	trigger := time.Now().UTC()
	if payload.TriggerTime != "" {
		if parsed, err := time.Parse(time.RFC3339, payload.TriggerTime); err == nil {
			trigger = parsed.UTC()
		}
	}

	// send_day uses ISO weekday: Mon=1 … Sun=7
	sendDay := int(trigger.Weekday())
	if sendDay == 0 {
		sendDay = 7 // Sunday → 7
	}
	sendHour := trigger.Hour()

	repos, err := j.repositoryRepo.GetDueForNewsletter(ctx, sendDay, sendHour)
	if err != nil {
		return fmt.Errorf("newsletter_scheduler: query repos: %w", err)
	}

	logger.Info(ctx, fmt.Sprintf("newsletter_scheduler: day=%d hour=%d eligible=%d", sendDay, sendHour, len(repos)))

	for _, repo := range repos {
		// Double-send guard: skip if a newsletter was already sent within the last 20 hours.
		if repo.NewsletterSentAt != nil && time.Since(*repo.NewsletterSentAt) < 20*time.Hour {
			logger.Info(ctx, fmt.Sprintf("newsletter_scheduler: skip repoId=%s (recent send)", repo.ID))
			continue
		}

		task, err := NewGenerateNewsletterTask(repo.ID)
		if err != nil {
			logger.Info(ctx, fmt.Sprintf("newsletter_scheduler: build task error repoId=%s: %v", repo.ID, err))
			continue
		}
		if _, err := j.enqueuer.Enqueue(task, asynq.MaxRetry(2), asynq.Queue("default")); err != nil {
			logger.Info(ctx, fmt.Sprintf("newsletter_scheduler: enqueue error repoId=%s: %v", repo.ID, err))
		}
	}

	return nil
}
