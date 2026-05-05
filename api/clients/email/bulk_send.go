package email

import (
	"context"
	"sync"
	"time"
)

const (
	bulkBatchSize  = 14              // emails per batch (stays within SES 14 msgs/s default quota)
	bulkBatchDelay = 1 * time.Second // pause between batches
)

// EmailSender is a function-level abstraction over SendEmailUsingSES, allowing injection in tests.
type EmailSender func(ctx context.Context, params AWSEmailParams)

// SendBulkEmailUsingSES sends params in concurrent batches of bulkBatchSize, with a
// bulkBatchDelay pause between batches to respect SES rate limits.
func SendBulkEmailUsingSES(ctx context.Context, params []AWSEmailParams, sender ...EmailSender) error {
	send := SendEmailUsingSES
	if len(sender) > 0 && sender[0] != nil {
		send = sender[0]
	}

	for i := 0; i < len(params); i += bulkBatchSize {
		end := i + bulkBatchSize
		if end > len(params) {
			end = len(params)
		}
		batch := params[i:end]

		var wg sync.WaitGroup
		for _, p := range batch {
			wg.Add(1)
			go func(p AWSEmailParams) {
				defer wg.Done()
				send(ctx, p)
			}(p)
		}
		wg.Wait()

		// Pause between batches (skip after the last batch)
		if end < len(params) {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(bulkBatchDelay):
			}
		}
	}
	return nil
}
