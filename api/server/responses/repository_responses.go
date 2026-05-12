package responses

import (
	github_client "abstractions/clients/github"
	"abstractions/db/models"
	"time"
)

type RepositoryResponse struct {
	ID                     string   `json:"id"`
	WorkspaceID            string   `json:"workspaceId"`
	UserID                 string   `json:"userId"`
	ServiceType            string   `json:"serviceType"`
	ServiceID              string   `json:"serviceId"`
	Name                   string   `json:"name"`
	Description            string   `json:"description"`
	DefaultBranch          string   `json:"defaultBranch"`
	Status                 string   `json:"status"`
	IndexedAt              *string  `json:"indexedAt"`
	SendDay                int      `json:"sendDay"`
	SendHour               int      `json:"sendHour"`
	NewsletterEnabled      bool     `json:"newsletterEnabled"`
	LastNewsletterSentAt   *string  `json:"lastNewsletterSentAt"`
	Emails                 []string `json:"emails"`
	SendToWorkspaceMembers bool     `json:"sendToWorkspaceMembers"`
	CreatedAt              string   `json:"createdAt"`
	UpdatedAt              string   `json:"updatedAt"`
}

func NewRepositoryResponse(repo *models.Repository) RepositoryResponse {
	var indexedAt *string
	if repo.IndexedAt != nil {
		formatted := repo.IndexedAt.UTC().Format(time.RFC3339)
		indexedAt = &formatted
	}

	var lastNewsletterSentAt *string
	if repo.NewsletterSentAt != nil {
		formatted := repo.NewsletterSentAt.UTC().Format(time.RFC3339)
		lastNewsletterSentAt = &formatted
	}

	emails := []string(repo.Recipients)
	if emails == nil {
		emails = []string{}
	}

	return RepositoryResponse{
		ID:                     repo.ID.String(),
		WorkspaceID:            repo.WorkspaceID.String(),
		UserID:                 repo.CreatedBy.String(),
		ServiceType:            repo.ServiceType,
		ServiceID:              repo.ServiceID,
		Name:                   repo.Name,
		Description:            "",
		DefaultBranch:          repo.DefaultBranch,
		Status:                 string(repo.Status),
		IndexedAt:              indexedAt,
		SendDay:                repo.SendDay,
		SendHour:               repo.SendHour,
		NewsletterEnabled:      repo.NewsletterEnabled,
		LastNewsletterSentAt:   lastNewsletterSentAt,
		Emails:                 emails,
		SendToWorkspaceMembers: repo.SendToWorkspaceMembers,
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
