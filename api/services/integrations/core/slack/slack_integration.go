package slack_integration

import (
	"abstractions/config"
	"abstractions/services/integrations/enums"
	"errors"
)

type SlackIntegration struct{}

var IntegrationTypeSlack enums.IntegrationTypeEnum = "slack"

func New() *SlackIntegration {
	return &SlackIntegration{}
}

func (s *SlackIntegration) GetEndpoints() []enums.EndpointDetails {
	return []enums.EndpointDetails{
		{
			Path:    "/api/v1/integrations/slack/oauth/callback",
			Method:  "GET",
			Public:  false,
			Webhook: false,
		},
		{
			Path:    "/api/v1/integrations/slack/webhooks/events",
			Method:  "POST",
			Public:  true,
			Webhook: true,
		},
	}
}

func (s *SlackIntegration) GetAuthType() enums.AuthType {
	return enums.AuthTypeOAuth2
}

func (s *SlackIntegration) GetIcon() string {
	return config.Environment.ClientURL + "/assets/img/slack.svg"
}

func (s *SlackIntegration) GetName() string {
	return "Slack"
}

func (s *SlackIntegration) GetIntegrationType() enums.IntegrationTypeEnum {
	return IntegrationTypeSlack
}

func (s *SlackIntegration) GetDescription() string {
	return "Connect with your Slack workspace to send messages, notifications, and receive updates from your customer interactions."
}

func (s *SlackIntegration) GetRequiredScopes() []string {
	return []string{
		"channels:read",
		"chat:write",
		"users:read",
		"team:read",
		"incoming-webhook",
	}
}

func (s *SlackIntegration) ValidateConfiguration(config map[string]interface{}) error {
	// Validate required configuration fields for Slack
	if config == nil {
		return errors.New("configuration cannot be nil")
	}

	// Check for required webhook URL if using incoming webhooks
	if webhookURL, exists := config["webhook_url"]; exists {
		if str, ok := webhookURL.(string); !ok || str == "" {
			return errors.New("webhook_url must be a non-empty string")
		}
	}

	// Check for channel configuration
	if channel, exists := config["default_channel"]; exists {
		if str, ok := channel.(string); !ok || str == "" {
			return errors.New("default_channel must be a non-empty string")
		}
	}

	return nil
}
