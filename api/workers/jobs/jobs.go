package jobs

import (
	"abstractions/config"
	"context"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type Job interface {
	ID() string
	Handler(ctx context.Context, t *asynq.Task) error
}

func AddHandlers(mux *asynq.ServeMux, db *gorm.DB, config *config.Config, enqueuer JobEnqueuer) {
	summarizeFileJob := NewSummarizeFileJob(db, config, enqueuer)
	mux.HandleFunc(summarizeFileJob.ID(), summarizeFileJob.Handler)

	indexRepositoryJob := NewIndexRepositoryJob(db, config, enqueuer)
	mux.HandleFunc(indexRepositoryJob.ID(), indexRepositoryJob.Handler)

	deleteFileJob := NewDeleteFileJob(db, config)
	mux.HandleFunc(deleteFileJob.ID(), deleteFileJob.Handler)

	newsletterSchedulerJob := NewNewsletterSchedulerJob(db, config, enqueuer)
	mux.HandleFunc(newsletterSchedulerJob.ID(), newsletterSchedulerJob.Handler)

	generateNewsletterJob := NewGenerateNewsletterJob(db, config)
	mux.HandleFunc(generateNewsletterJob.ID(), generateNewsletterJob.Handler)

	generateTopicsJob := NewGenerateTopicsJob(db, config)
	mux.HandleFunc(generateTopicsJob.ID(), generateTopicsJob.Handler)
}
