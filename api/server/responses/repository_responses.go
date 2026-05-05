package responses

import (
	github_client "abstractions/clients/github"
	"abstractions/db/models"
	"time"
)

type RepositoryResponse struct {
	ID                   string  `json:"id"`
	WorkspaceID          string  `json:"workspaceId"`
	UserID               string  `json:"userId"`
	GitHubRepoID         int64   `json:"githubRepoId"`
	FullName             string  `json:"fullName"`
	DefaultBranch        string  `json:"defaultBranch"`
	WebhookID            *int64  `json:"webhookId"`
	Status               string  `json:"status"`
	IndexedAt            *string `json:"indexedAt"`
	SendDay              int     `json:"sendDay"`
	SendHour             int     `json:"sendHour"`
	LastNewsletterSentAt *string `json:"lastNewsletterSentAt"`
	CreatedAt            string  `json:"createdAt"`
	UpdatedAt            string  `json:"updatedAt"`
}

func NewRepositoryResponse(repo *models.Repository) RepositoryResponse {
	var indexedAt *string
	if repo.IndexedAt != nil {
		formatted := repo.IndexedAt.UTC().Format(time.RFC3339)
		indexedAt = &formatted
	}

	var lastNewsletterSentAt *string
	if repo.LastNewsletterSentAt != nil {
		formatted := repo.LastNewsletterSentAt.UTC().Format(time.RFC3339)
		lastNewsletterSentAt = &formatted
	}

	return RepositoryResponse{
		ID:                   repo.ID.String(),
		WorkspaceID:          repo.WorkspaceID.String(),
		UserID:               repo.UserID.String(),
		GitHubRepoID:         repo.GitHubRepoID,
		FullName:             repo.FullName,
		DefaultBranch:        repo.DefaultBranch,
		WebhookID:            repo.WebhookID,
		Status:               string(repo.Status),
		IndexedAt:            indexedAt,
		SendDay:              repo.SendDay,
		SendHour:             repo.SendHour,
		LastNewsletterSentAt: lastNewsletterSentAt,
		CreatedAt:            repo.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:            repo.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func NewRepositoryResponses(repos []*models.Repository) []RepositoryResponse {
	result := make([]RepositoryResponse, len(repos))
	for i, r := range repos {
		result[i] = NewRepositoryResponse(r)
	}
	return result
}

// GitHubRepoResponse is a minimal representation of a GitHub repo for the listing endpoint.
type GitHubRepoResponse struct {
	ID            int64  `json:"id"`
	FullName      string `json:"fullName"`
	DefaultBranch string `json:"defaultBranch"`
	Private       bool   `json:"private"`
	Description   string `json:"description"`
	HTMLURL       string `json:"htmlUrl"`
}

func NewGitHubRepoResponse(r github_client.GitHubRepo) GitHubRepoResponse {
	return GitHubRepoResponse{
		ID:            r.ID,
		FullName:      r.FullName,
		DefaultBranch: r.DefaultBranch,
		Private:       r.Private,
		Description:   r.Description,
		HTMLURL:       r.HTMLURL,
	}
}

func NewGitHubRepoResponses(repos []github_client.GitHubRepo) []GitHubRepoResponse {
	result := make([]GitHubRepoResponse, len(repos))
	for i, r := range repos {
		result[i] = NewGitHubRepoResponse(r)
	}
	return result
}
