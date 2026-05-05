package jobs

import (
	"abstractions/config"
	"context"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type ExtractMetadataJob struct {
	db     *gorm.DB
	config *config.Config
}

func NewExtractMetadataJob(db *gorm.DB, config *config.Config) *ExtractMetadataJob {
	return &ExtractMetadataJob{db: db, config: config}
}

func (j *ExtractMetadataJob) ID() string {
	return "extract_metadata"
}

func (j *ExtractMetadataJob) Handler(ctx context.Context, t *asynq.Task) error {
	return nil
}
