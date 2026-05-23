package responses

type WorkspaceGettingStartedResponse struct {
	AIConfigured    bool `json:"aiConfigured"`
	GitHubConnected bool `json:"githubConnected"`
	NewsletterSent  bool `json:"newsletterSent"`
}
