package jobs

import (
	"abstractions/config"
	"context"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type QueueEmailJob struct {
	db     *gorm.DB
	config *config.Config
}

func NewQueueEmailJob(db *gorm.DB, config *config.Config) *QueueEmailJob {
	return &QueueEmailJob{db: db, config: config}
}

func (j *QueueEmailJob) ID() string {
	return "import_email"
}

func (j *QueueEmailJob) Handler(ctx context.Context, t *asynq.Task) error {
	return nil
}
