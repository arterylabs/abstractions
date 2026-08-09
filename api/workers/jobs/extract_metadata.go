package jobs

import (
	"abstractions/config"
	"abstractions/logger"
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

// ExtractMetadataJob is retained for queue compatibility. Metadata extraction now
// happens inline during SummarizeFileJob (language, tags, symbols, service metadata).
// The handler acknowledges and drops legacy tasks without side effects.
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
	// Drain any legacy tasks that may still be on the queue.
	var payload map[string]interface{}
	_ = json.Unmarshal(t.Payload(), &payload)
	logger.Info(ctx, fmt.Sprintf("extract_metadata: noop (handled by summarize_file) payloadKeys=%d", len(payload)))
	return nil
}
