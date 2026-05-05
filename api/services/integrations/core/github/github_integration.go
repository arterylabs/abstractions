package github_integration

import (
	"abstractions/config"
	"abstractions/services/integrations/enums"
	"errors"
)

type GitHubIntegration struct{}

var IntegrationTypeGitHub enums.IntegrationTypeEnum = "github"

func New() *GitHubIntegration {
	return &GitHubIntegration{}
}

func (g *GitHubIntegration) GetEndpoints() []enums.EndpointDetails {
	return []enums.EndpointDetails{
		{
			Path:    "/api/v1/integrations/github/oauth/callback",
			Method:  "GET",
			Public:  false,
			Webhook: false,
		},
	}
}

func (g *GitHubIntegration) GetAuthType() enums.AuthType {
	return enums.AuthTypeOAuth2
}

func (g *GitHubIntegration) GetIcon() string {
	return config.Environment.ClientURL + "/assets/img/github.svg"
}

func (g *GitHubIntegration) GetName() string {
	return "GitHub"
}

func (g *GitHubIntegration) GetIntegrationType() enums.IntegrationTypeEnum {
	return IntegrationTypeGitHub
}

func (g *GitHubIntegration) GetDescription() string {
	return "Connect your GitHub account to enable CodeBrief to read your repositories and generate automated newsletters for your engineering team."
}

// GetRequiredScopes returns the OAuth scopes needed: repo gives read access to private repos
// and the ability to create webhooks on repos where the user has admin access.
func (g *GitHubIntegration) GetRequiredScopes() []string {
	return []string{"repo"}
}

func (g *GitHubIntegration) ValidateConfiguration(config map[string]interface{}) error {
	if config == nil {
		return errors.New("configuration cannot be nil")
	}
	return nil
}
