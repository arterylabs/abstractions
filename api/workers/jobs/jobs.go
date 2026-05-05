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

func AddHandlers(mux *asynq.ServeMux, db *gorm.DB, config *config.Config) {
	queueEmailJob := NewQueueEmailJob(db, config)
	mux.HandleFunc(queueEmailJob.ID(), queueEmailJob.Handler)

	downloadEmailsJob := NewDownloadEmailsJob(db, config)
	mux.HandleFunc(downloadEmailsJob.ID(), downloadEmailsJob.Handler)

	extractMetadataJob := NewExtractMetadataJob(db, config)
	mux.HandleFunc(extractMetadataJob.ID(), extractMetadataJob.Handler)
}
