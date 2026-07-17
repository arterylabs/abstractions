package email

import (
	"abstractions/logger"
	"context"
	"sync"
	"time"
)

const (
	bulkBatchSize  = 14              // emails per batch (stays within SES 14 msgs/s default quota)
	bulkBatchDelay = 1 * time.Second // pause between batches
)

// EmailSender is a function-level abstraction over SendEmailUsingSES, allowing injection in tests.
type EmailSender func(ctx context.Context, params AWSEmailParams) error

// SendBulkEmailUsingSES sends params in concurrent batches of bulkBatchSize, with a
// bulkBatchDelay pause between batches to respect provider rate limits. Despite the name
// (kept for backwards compatibility), it sends through the configured provider dispatcher
// (SendEmail) by default, not necessarily SES.
func SendBulkEmailUsingSES(ctx context.Context, params []AWSEmailParams, sender ...EmailSender) error {
	send := SendEmail
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
				// Individual failures are already logged/suppressed inside the sender;
				// a failed send must not stop the rest of the batch.
				if err := send(ctx, p); err != nil {
					logger.Error(ctx, "email: bulk send failed for "+p.Email, err)
				}
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
