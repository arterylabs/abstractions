package responses

import (
	github_client "abstractions/clients/github"
	"abstractions/db/models"
	"time"
)

type RepositoryResponse struct {
	ID                     string   `json:"id"`
	WorkspaceID            string   `json:"workspaceId"`
	CreatedBy              string   `json:"createdBy"`
	Name                   string   `json:"name"`
	ServiceType            string   `json:"serviceType"`
	ServiceID              string   `json:"serviceId"`
	DefaultBranch          string   `json:"defaultBranch"`
	Status                 string   `json:"status"`
	IndexedAt              *string  `json:"indexedAt"`
	SendDay                int      `json:"sendDay"`
	SendHour               int      `json:"sendHour"`
	Recipients             []string `json:"recipients"`
	SendToWorkspaceMembers bool     `json:"sendToWorkspaceMembers"`
	NewsletterSentAt       *string  `json:"newsletterSentAt"`
	CreatedAt              string   `json:"createdAt"`
	UpdatedAt              string   `json:"updatedAt"`
}

func NewRepositoryResponse(repo *models.Repository) RepositoryResponse {
	var indexedAt *string
	if repo.IndexedAt != nil {
		formatted := repo.IndexedAt.UTC().Format(time.RFC3339)
		indexedAt = &formatted
	}

	var newsletterSentAt *string
	if repo.NewsletterSentAt != nil {
		formatted := repo.NewsletterSentAt.UTC().Format(time.RFC3339)
		newsletterSentAt = &formatted
	}

	recipients := []string(repo.Recipients)
	if recipients == nil {
		recipients = []string{}
	}

	return RepositoryResponse{
		ID:                     repo.ID.String(),
		WorkspaceID:            repo.WorkspaceID.String(),
		CreatedBy:              repo.CreatedBy.String(),
		Name:                   repo.Name,
		ServiceType:            repo.ServiceType,
		ServiceID:              repo.ServiceID,
		DefaultBranch:          repo.DefaultBranch,
		Status:                 string(repo.Status),
		IndexedAt:              indexedAt,
		SendDay:                repo.SendDay,
		SendHour:               repo.SendHour,
		Recipients:             recipients,
		SendToWorkspaceMembers: repo.SendToWorkspaceMembers,
		NewsletterSentAt:       newsletterSentAt,
		CreatedAt:              repo.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:              repo.UpdatedAt.UTC().Format(time.RFC3339),
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
