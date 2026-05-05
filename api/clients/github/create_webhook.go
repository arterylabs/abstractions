package github_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CreateWebhook creates a GitHub webhook on the given repo and returns the webhook ID.
// payloadURL is the public HTTPS URL that GitHub will POST push events to.
// secret is a per-repo random string used to sign the HMAC-SHA256 payload signature.
func CreateWebhook(ctx context.Context, accessToken, owner, repo, payloadURL, secret string) (int64, error) {
	type webhookConfig struct {
		URL         string `json:"url"`
		ContentType string `json:"content_type"`
		Secret      string `json:"secret"`
		InsecureSSL string `json:"insecure_ssl"`
	}
	type createWebhookRequest struct {
		Name   string        `json:"name"`
		Active bool          `json:"active"`
		Events []string      `json:"events"`
		Config webhookConfig `json:"config"`
	}
	type createWebhookResponse struct {
		ID int64 `json:"id"`
	}

	reqBody := createWebhookRequest{
		Name:   "web",
		Active: true,
		Events: []string{"push", "delete"},
		Config: webhookConfig{
			URL:         payloadURL,
			ContentType: "json",
			Secret:      secret,
			InsecureSSL: "0",
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return 0, fmt.Errorf("github_client: marshal webhook request: %w", err)
	}

	url := fmt.Sprintf("%s/repos/%s/%s/hooks", githubAPIBase, owner, repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return 0, fmt.Errorf("github_client: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("github_client: do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("github_client: read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("github_client: unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	var result createWebhookResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return 0, fmt.Errorf("github_client: unmarshal response: %w", err)
	}

	return result.ID, nil
}

// httpGetURL performs a plain GET request and returns the response (caller must close Body).
func httpGetURL(url string) (*http.Response, error) {
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
