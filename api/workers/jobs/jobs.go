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
	queueEmailJob := NewQueueEmailJob(db, config)
	mux.HandleFunc(queueEmailJob.ID(), queueEmailJob.Handler)

	downloadEmailsJob := NewDownloadEmailsJob(db, config)
	mux.HandleFunc(downloadEmailsJob.ID(), downloadEmailsJob.Handler)

	extractMetadataJob := NewExtractMetadataJob(db, config)
	mux.HandleFunc(extractMetadataJob.ID(), extractMetadataJob.Handler)

	summarizeFileJob := NewSummarizeFileJob(db, config, enqueuer)
	mux.HandleFunc(summarizeFileJob.ID(), summarizeFileJob.Handler)

	indexRepositoryJob := NewIndexRepositoryJob(db, config, enqueuer)
	mux.HandleFunc(indexRepositoryJob.ID(), indexRepositoryJob.Handler)

	newsletterSchedulerJob := NewNewsletterSchedulerJob(db, config, enqueuer)
	mux.HandleFunc(newsletterSchedulerJob.ID(), newsletterSchedulerJob.Handler)

	generateNewsletterJob := NewGenerateNewsletterJob(db, config)
	mux.HandleFunc(generateNewsletterJob.ID(), generateNewsletterJob.Handler)

	generateTopicsJob := NewGenerateTopicsJob(db, config)
	mux.HandleFunc(generateTopicsJob.ID(), generateTopicsJob.Handler)
}
