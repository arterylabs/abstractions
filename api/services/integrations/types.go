package integration_service

import (
	github_integration "abstractions/services/integrations/core/github"
	slack_integration "abstractions/services/integrations/core/slack"
	"abstractions/services/integrations/enums"
)

var integrationTypesList = map[enums.IntegrationTypeEnum]IntegrationType{}

// Initialize the integration types map with supported integrations
func init() {
	slackIntegration := slack_integration.New()
	integrationTypesList[slackIntegration.GetIntegrationType()] = slackIntegration

	githubIntegration := github_integration.New()
	integrationTypesList[githubIntegration.GetIntegrationType()] = githubIntegration
}

// GetAllSupportedTypes returns a list of all supported integration types
func GetAllSupportedTypes() []enums.IntegrationTypeEnum {
	types := []enums.IntegrationTypeEnum{}
	for integrationType := range integrationTypesList {
		types = append(types, integrationType)
	}
	return types
}

// IntegrationType defines the interface that all integrations must implement
type IntegrationType interface {
	// GetEndpoints returns the list of API endpoints this integration supports
	GetEndpoints() []enums.EndpointDetails
	// GetAuthType returns the authentication type required
	GetAuthType() enums.AuthType
	// GetIntegrationType returns the integration type enum
	GetIntegrationType() enums.IntegrationTypeEnum
	// GetIcon returns the icon identifier for UI display
	GetIcon() string
	// GetName returns the human-readable name
	GetName() string
	// GetDescription returns a description of what this integration does
	GetDescription() string
	// GetRequiredScopes returns the OAuth scopes required for this integration
	GetRequiredScopes() []string
	// ValidateConfiguration validates the integration configuration
	ValidateConfiguration(config map[string]interface{}) error
}

// ErrUnsupportedIntegrationType is returned when an integration type is not supported
type ErrUnsupportedIntegrationType struct {
	Type enums.IntegrationTypeEnum
}

func (e *ErrUnsupportedIntegrationType) Error() string {
	return "unsupported integration type: " + string(e.Type)
}

// GetIntegrationType returns the integration implementation for a given type
func GetIntegrationType(integrationType enums.IntegrationTypeEnum) (IntegrationType, error) {
	integration, exists := integrationTypesList[integrationType]
	if !exists {
		return nil, &ErrUnsupportedIntegrationType{Type: integrationType}
	}
	return integration, nil
}

// GetAllIntegrationTypes returns all registered integration types
func GetAllIntegrationTypes() []IntegrationType {
	integrations := make([]IntegrationType, 0, len(integrationTypesList))
	for _, integration := range integrationTypesList {
		integrations = append(integrations, integration)
	}
	return integrations
}

// IsValidType checks if the integration type is supported
func IsValidType(integrationType enums.IntegrationTypeEnum) bool {
	_, exists := integrationTypesList[integrationType]
	return exists
}
