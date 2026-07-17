package github_client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const githubAPIBase = "https://api.github.com"

// ListUserRepos fetches the authenticated user's repositories (first page, up to 100).
func ListUserRepos(ctx context.Context, accessToken string) ([]GitHubRepo, error) {
	url := githubAPIBase + "/user/repos?per_page=100&sort=updated&type=all"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("github_client: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github_client: do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("github_client: read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github_client: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var repos []GitHubRepo
	if err := json.Unmarshal(body, &repos); err != nil {
		return nil, fmt.Errorf("github_client: unmarshal response: %w", err)
	}
	return repos, nil
}

// GitHubRepo is a minimal representation of a GitHub repository from the /user/repos endpoint.
type GitHubRepo struct {
	ID            int64  `json:"id"`
	FullName      string `json:"full_name"`
	DefaultBranch string `json:"default_branch"`
	Private       bool   `json:"private"`
	Description   string `json:"description"`
	HTMLURL       string `json:"html_url"`
}
