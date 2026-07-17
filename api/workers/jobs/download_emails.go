package jobs

import (
	"abstractions/config"
	"context"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type DownloadEmailsJob struct {
	db     *gorm.DB
	config *config.Config
}

func NewDownloadEmailsJob(db *gorm.DB, config *config.Config) *DownloadEmailsJob {
	return &DownloadEmailsJob{db: db, config: config}
}

func (j *DownloadEmailsJob) ID() string {
	return "download_emails"
}

func (j *DownloadEmailsJob) Handler(ctx context.Context, t *asynq.Task) error {
	return nil
}
