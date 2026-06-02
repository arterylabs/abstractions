package jobs

import (
	"abstractions/config"
	"abstractions/db/models"
	"abstractions/logger"
	"abstractions/repositories"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

const NewsletterSchedulerJobID = "newsletter_scheduler"

// NewsletterSchedulerJob queries repos whose send time matches the current hour and
// enqueues a GenerateNewsletterJob for each qualifying repo.
// It is designed to be enqueued every hour by the asynq Scheduler.
type NewsletterSchedulerJob struct {
	db             *gorm.DB
	config         *config.Config
	enqueuer       JobEnqueuer
	repositoryRepo *repositories.RepositoryRepository
	newsletterRepo *repositories.NewsletterRepository
}

func NewNewsletterSchedulerJob(db *gorm.DB, cfg *config.Config, enqueuer JobEnqueuer) *NewsletterSchedulerJob {
	return &NewsletterSchedulerJob{
		db:             db,
		config:         cfg,
		enqueuer:       enqueuer,
		repositoryRepo: repositories.NewRepositoryRepository(db),
		newsletterRepo: repositories.NewNewsletterRepository(db),
	}
}

func (j *NewsletterSchedulerJob) ID() string {
	return NewsletterSchedulerJobID
}

// NewNewsletterSchedulerTask builds the asynq task to be enqueued by the periodic scheduler.
func NewNewsletterSchedulerTask() (*asynq.Task, error) {
	return asynq.NewTask(NewsletterSchedulerJobID, nil), nil
}

// Handler processes a NewsletterSchedulerJob task.
// It finds all indexed repos matching the current send_day + send_hour,
// applies the 2-hour double-send guard, and enqueues GenerateNewsletterJob for each.
func (j *NewsletterSchedulerJob) Handler(ctx context.Context, t *asynq.Task) error {
	trigger := time.Now().UTC()

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
		// Double-send guard: skip if a newsletter was already sent within the last 2 hours.
		if repo.NewsletterSentAt != nil && time.Since(*repo.NewsletterSentAt) < 2*time.Hour {
			logger.Info(ctx, fmt.Sprintf("newsletter_scheduler: skip repoId=%s (recent send)", repo.ID))
			continue
		}

		// Create a queued newsletter record so the status is visible immediately.
		nl := &models.Newsletter{
			ID:     uuid.New(),
			RepoID: repo.ID,
			Status: models.NewsletterStatusQueued,
		}
		if err := j.newsletterRepo.Create(ctx, nl); err != nil {
			logger.Info(ctx, fmt.Sprintf("newsletter_scheduler: create newsletter error repoId=%s: %v", repo.ID, err))
			continue
		}

		if err := j.enqueuer.EnqueueGenerateNewsletter(repo.ID, nl.ID); err != nil {
			logger.Info(ctx, fmt.Sprintf("newsletter_scheduler: enqueue error repoId=%s: %v", repo.ID, err))
		}
	}

	return nil
}
