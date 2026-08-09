package handlers

import (
	"abstractions/db/models"
	"abstractions/logger"
	"abstractions/repositories"
	"abstractions/server"
	"abstractions/utils/context_utils"
	"abstractions/workers/jobs"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// GitHubWebhookHandler receives GitHub push events and re-indexes affected repositories.
type GitHubWebhookHandler struct {
	*BaseHandler
	jobClient *jobs.Client
}

func NewGitHubWebhookHandler(server *server.Server, jobClient *jobs.Client) *GitHubWebhookHandler {
	return &GitHubWebhookHandler{
		BaseHandler: NewBaseHandler(server),
		jobClient:   jobClient,
	}
}

type githubPushPayload struct {
	Ref        string `json:"ref"`
	Deleted    bool   `json:"deleted"`
	Repository struct {
		ID       int64  `json:"id"`
		FullName string `json:"full_name"`
		Name     string `json:"name"`
		// DefaultBranch is the repo default branch name (without refs/heads/).
		DefaultBranch string `json:"default_branch"`
	} `json:"repository"`
	Commits []struct {
		Added    []string `json:"added"`
		Removed  []string `json:"removed"`
		Modified []string `json:"modified"`
	} `json:"commits"`
}

// HandlePush processes GitHub webhook deliveries.
// Only push events to the default branch trigger re-indexing.
func (h *GitHubWebhookHandler) HandlePush(c echo.Context) error {
	ctx := c.Request().Context()
	event := c.Request().Header.Get("X-GitHub-Event")
	if event == "ping" {
		return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
	}
	if event != "" && event != "push" {
		return c.JSON(http.StatusOK, map[string]string{"message": "ignored event"})
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to read body")
	}

	if secret := h.Server().Config.Environment.GitHubWebhookSecret; secret != "" {
		sig := c.Request().Header.Get("X-Hub-Signature-256")
		if !verifyGitHubSignature(secret, body, sig) {
			logger.Warn(ctx, "github_webhook: invalid signature")
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid signature")
		}
	}

	var payload githubPushPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid json payload")
	}

	if payload.Deleted {
		return c.JSON(http.StatusOK, map[string]string{"message": "branch deleted ignored"})
	}

	serviceID := fmt.Sprintf("%d", payload.Repository.ID)
	db := context_utils.GetTransactionFromContext(c)
	repoRepo := repositories.NewRepositoryRepository(db)
	repos, err := repoRepo.GetByServiceID(ctx, models.RepositoryServiceTypeGitHub, serviceID)
	if err != nil {
		logger.Error(ctx, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "lookup failed")
	}
	if len(repos) == 0 {
		// Also try matching by full name if service id differs (legacy rows).
		return c.JSON(http.StatusOK, map[string]string{"message": "no matching repository"})
	}

	// Only process pushes to the default branch (or refs/heads/<default>).
	defaultBranch := payload.Repository.DefaultBranch
	expectedRef := "refs/heads/" + defaultBranch
	if defaultBranch == "" {
		// Fall back to first connected repo's configured branch.
		defaultBranch = repos[0].DefaultBranch
		expectedRef = "refs/heads/" + defaultBranch
	}
	if payload.Ref != "" && defaultBranch != "" && payload.Ref != expectedRef {
		return c.JSON(http.StatusOK, map[string]string{"message": "non-default branch ignored"})
	}

	// Collect removed paths for targeted deletes; always also re-index for SHA reconciliation.
	removed := map[string]struct{}{}
	for _, commit := range payload.Commits {
		for _, p := range commit.Removed {
			removed[p] = struct{}{}
		}
	}

	enqueued := 0
	for _, repo := range repos {
		if repo.Status == models.RepositoryStatusIndexing {
			logger.Info(ctx, fmt.Sprintf("github_webhook: skip already indexing repoId=%s", repo.ID.String()))
			continue
		}
		for path := range removed {
			if err := h.jobClient.EnqueueDeleteFile(repo.ID, path); err != nil {
				logger.Error(ctx, fmt.Errorf("github_webhook: enqueue delete %s: %w", path, err))
			}
		}
		if err := h.jobClient.EnqueueIndexRepository(repo.ID, repo.ServiceID); err != nil {
			logger.Error(ctx, fmt.Errorf("github_webhook: enqueue index: %w", err))
			continue
		}
		enqueued++
		logger.Info(ctx, fmt.Sprintf("github_webhook: enqueued reindex repoId=%s fullName=%s", repo.ID.String(), repo.Name))
	}

	return c.JSON(http.StatusAccepted, map[string]interface{}{
		"message":  "processed",
		"enqueued": enqueued,
	})
}

func verifyGitHubSignature(secret string, body []byte, signatureHeader string) bool {
	if signatureHeader == "" || !strings.HasPrefix(signatureHeader, "sha256=") {
		return false
	}
	sigHex := strings.TrimPrefix(signatureHeader, "sha256=")
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(sigHex))
}
