package jobs

import (
	"errors"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

// Queue name constants used across all job enqueue calls.
const (
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"
)

// Client wraps asynq.Client and exposes job-specific enqueue methods.
// Code outside the workers package should use this type instead of importing asynq directly.
type Client struct {
	c *asynq.Client
}

// NewClient creates a Client connected to Redis at the given address (host:port).
func NewClient(redisAddr string) *Client {
	return &Client{c: asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})}
}

// Close releases the underlying Redis connection.
func (c *Client) Close() error {
	return c.c.Close()
}

// EnqueueIndexRepository enqueues an index-repository job.
// serviceID is used to build a stable, unique task ID so duplicate submissions are silently
// dropped rather than treated as errors.
func (c *Client) EnqueueIndexRepository(repositoryID uuid.UUID, serviceID string) error {
	task, err := NewIndexRepositoryTask(repositoryID)
	if err != nil {
		return err
	}
	_, err = c.c.Enqueue(task,
		asynq.MaxRetry(3),
		asynq.Queue(QueueDefault),
		asynq.TaskID(IndexRepositoryJobID+":"+serviceID),
	)
	if errors.Is(err, asynq.ErrTaskIDConflict) {
		// Task is already queued — treat as a no-op.
		return nil
	}
	return err
}

// EnqueueGenerateNewsletter enqueues a full newsletter generation and delivery job.
func (c *Client) EnqueueGenerateNewsletter(repoID uuid.UUID) error {
	task, err := NewGenerateNewsletterTask(repoID)
	if err != nil {
		return err
	}
	_, err = c.c.Enqueue(task, asynq.MaxRetry(2), asynq.Queue(QueueDefault))
	return err
}

// EnqueueGenerateNewsletterDraft enqueues a draft-only newsletter generation job (no delivery).
func (c *Client) EnqueueGenerateNewsletterDraft(repoID uuid.UUID) error {
	task, err := NewGenerateNewsletterDraftTask(repoID)
	if err != nil {
		return err
	}
	_, err = c.c.Enqueue(task, asynq.MaxRetry(2), asynq.Queue(QueueDefault))
	return err
}

// EnqueueSummarizeFile enqueues a summarize-file job for a single repository file.
// Used internally by IndexRepositoryJob when sub-enqueueing per-file work.
func (c *Client) EnqueueSummarizeFile(repositoryID uuid.UUID, filePath, gitSHA string) error {
	task, err := NewSummarizeFileTask(repositoryID, filePath, gitSHA)
	if err != nil {
		return err
	}
	_, err = c.c.Enqueue(task, asynq.MaxRetry(3), asynq.Queue(QueueDefault))
	return err
}

// EnqueueGenerateTopics enqueues a topic-generation job immediately after repository indexing.
// Used internally by IndexRepositoryJob after a repository finishes indexing.
func (c *Client) EnqueueGenerateTopics(repoID uuid.UUID) error {
	task, err := NewGenerateTopicsTask(repoID)
	if err != nil {
		return err
	}
	_, err = c.c.Enqueue(task, asynq.MaxRetry(2), asynq.Queue(QueueDefault))
	return err
}
