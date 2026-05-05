package handlers

import (
	github_client "abstractions/clients/github"
	"abstractions/logger"
	"abstractions/server"
	"abstractions/server/responses"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
)

// webhookTaskEnqueuer enqueues asynq tasks (implemented by *asynq.Client).
type webhookTaskEnqueuer interface {
	Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)
}

// WebhookHandler handles inbound webhook events from external services.
type WebhookHandler struct {
	*BaseHandler
	enqueuer webhookTaskEnqueuer
}

func NewWebhookHandler(server *server.Server, enqueuer webhookTaskEnqueuer) *WebhookHandler {
	return &WebhookHandler{
		BaseHandler: NewBaseHandler(server),
		enqueuer:    enqueuer,
	}
}

// ---- task type constants and payloads (reused by workers/jobs in task 4.2) ----

const UpdateFileJobID = "update_file"
const DeleteFileJobID = "delete_file"

type UpdateFilePayload struct {
	RepositoryID string `json:"repositoryId"`
	FilePath     string `json:"filePath"`
	SHA          string `json:"sha"`
}

type DeleteFilePayload struct {
	RepositoryID string `json:"repositoryId"`
	FilePath     string `json:"filePath"`
}

// ---- GitHub webhook push event structs ----

type githubCommit struct {
	Added    []string `json:"added"`
	Modified []string `json:"modified"`
	Removed  []string `json:"removed"`
}

type githubPushPayload struct {
	After   string         `json:"after"` // HEAD commit SHA after push
	Commits []githubCommit `json:"commits"`
}

// ---- handler ----

// ReceiveGitHubWebhook is a public endpoint that receives GitHub webhook deliveries.
// It verifies the HMAC-SHA256 signature using the repo's stored webhook_secret, then
// enqueues incremental update or delete jobs for changed files.
func (h *WebhookHandler) ReceiveGitHubWebhook(c echo.Context) error {
	ctx := c.Request().Context()

	// Read raw body before any parsing (required for HMAC verification)
	rawBody, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "failed to read request body")
	}

	// Look up the repository by webhook_id from the X-GitHub-Hook-ID header
	hookIDStr := c.Request().Header.Get("X-GitHub-Hook-ID")
	if hookIDStr == "" {
		return responses.ErrorResponse(c, http.StatusBadRequest, "missing X-GitHub-Hook-ID header")
	}
	hookID, err := strconv.ParseInt(hookIDStr, 10, 64)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid X-GitHub-Hook-ID header")
	}

	repo, err := h.Services(c).Repository.GetByWebhookID(ctx, hookID)
	if err != nil || repo == nil {
		logger.Info(ctx, fmt.Sprintf("webhook: repo not found for webhook_id=%d", hookID))
		// Return 200 so GitHub doesn't keep retrying for removed/unknown webhooks
		return c.NoContent(http.StatusOK)
	}

	// Verify HMAC-SHA256 signature
	signature := c.Request().Header.Get("X-Hub-Signature-256")
	webhookSecret := ""
	if repo.WebhookSecret != nil {
		webhookSecret = *repo.WebhookSecret
	}
	if !verifyGitHubSignature(rawBody, webhookSecret, signature) {
		logger.Info(ctx, fmt.Sprintf("webhook: invalid signature for repo %s", repo.ID))
		return responses.ErrorResponse(c, http.StatusUnauthorized, "invalid webhook signature")
	}

	event := c.Request().Header.Get("X-GitHub-Event")
	logger.Info(ctx, fmt.Sprintf("webhook: event=%s repoId=%s", event, repo.ID))

	switch event {
	case "push":
		var push githubPushPayload
		if err := json.Unmarshal(rawBody, &push); err != nil {
			return responses.ErrorResponse(c, http.StatusBadRequest, "failed to parse push event")
		}
		h.enqueuePushFileJobs(ctx, repo.ID.String(), push)
	}

	// Always return 200 quickly
	return c.NoContent(http.StatusOK)
}

// enqueuePushFileJobs deduplicates modified/added files across all commits, filters ignored paths,
// and enqueues one UpdateFileJob per file and one DeleteFileJob per removed file.
func (h *WebhookHandler) enqueuePushFileJobs(ctx context.Context, repoID string, push githubPushPayload) {
	updatedFiles := make(map[string]struct{})
	deletedFiles := make(map[string]struct{})

	for _, commit := range push.Commits {
		for _, path := range commit.Added {
			if !github_client.ShouldIgnore(path, 0) {
				updatedFiles[path] = struct{}{}
			}
		}
		for _, path := range commit.Modified {
			if !github_client.ShouldIgnore(path, 0) {
				updatedFiles[path] = struct{}{}
			}
		}
		for _, path := range commit.Removed {
			deletedFiles[path] = struct{}{}
			delete(updatedFiles, path) // removed overrides modified
		}
	}

	for path := range updatedFiles {
		b, err := json.Marshal(UpdateFilePayload{
			RepositoryID: repoID,
			FilePath:     path,
			SHA:          push.After,
		})
		if err != nil {
			continue
		}
		if _, err := h.enqueuer.Enqueue(asynq.NewTask(UpdateFileJobID, b), asynq.MaxRetry(3), asynq.Queue("default")); err != nil {
			logger.Info(ctx, fmt.Sprintf("webhook: enqueue update_file error %s: %v", path, err))
		}
	}

	for path := range deletedFiles {
		b, err := json.Marshal(DeleteFilePayload{RepositoryID: repoID, FilePath: path})
		if err != nil {
			continue
		}
		if _, err := h.enqueuer.Enqueue(asynq.NewTask(DeleteFileJobID, b), asynq.MaxRetry(3), asynq.Queue("default")); err != nil {
			logger.Info(ctx, fmt.Sprintf("webhook: enqueue delete_file error %s: %v", path, err))
		}
	}
}

// verifyGitHubSignature validates the X-Hub-Signature-256 header using HMAC-SHA256 with a
// timing-safe comparison to prevent timing attacks.
func verifyGitHubSignature(body []byte, secret, signature string) bool {
	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}
