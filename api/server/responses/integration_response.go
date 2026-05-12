package responses

import (
	"abstractions/db/models"
	"abstractions/services/integrations/enums"
	"time"
)

type PaginatedIntegrationsResponse struct {
	Results  []IntegrationResponse `json:"results"`
	PageInfo PageInfo              `json:"pageInfo"`
}

type IntegrationResponse struct {
	ID              string                    `json:"id"`
	Type            enums.IntegrationTypeEnum `json:"type"`
	ServiceTitle    *string                   `json:"serviceTitle"`
	ServiceIcon     *string                   `json:"serviceIcon"`
	WorkspaceID     string                    `json:"workspaceId"`
	ResourceType    *string                   `json:"resourceType"`
	ResourceID      *string                   `json:"resourceId"`
	TokenExpiresAt  *string                   `json:"tokenExpiresAt"`
	Scopes          *string                   `json:"scopes"`
	ServiceTenantID *string                   `json:"serviceTenantId"`
	ServiceUserID   *string                   `json:"serviceUserId"`
	Configuration   map[string]interface{}    `json:"configuration"`
	CreatedAt       string                    `json:"createdAt"`
	UpdatedAt       string                    `json:"updatedAt"`
}

func NewIntegrationResponse(integration *models.Integration) IntegrationResponse {
	var tokenExpiresAt *string
	if integration.TokenExpiresAt != nil {
		formatted := integration.TokenExpiresAt.UTC().Format(time.RFC3339)
		tokenExpiresAt = &formatted
	}

	var resourceIDStr *string
	if integration.ResourceID != nil {
		s := integration.ResourceID.String()
		resourceIDStr = &s
	}

	return IntegrationResponse{
		ID:              integration.ID.String(),
		Type:            integration.Type,
		ServiceTitle:    integration.ServiceTitle,
		ServiceIcon:     integration.ServiceIcon,
		WorkspaceID:     integration.WorkspaceID.String(),
		ResourceType:    integration.ResourceType,
		ResourceID:      resourceIDStr,
		TokenExpiresAt:  tokenExpiresAt,
		Scopes:          integration.Scopes,
		ServiceTenantID: integration.ServiceTenantID,
		ServiceUserID:   integration.ServiceUserID,
		Configuration:   integration.Configuration,
		CreatedAt:       integration.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:       integration.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func NewIntegrationResponses(integrations []*models.Integration) []IntegrationResponse {
	responses := make([]IntegrationResponse, len(integrations))
	for i, integration := range integrations {
		responses[i] = NewIntegrationResponse(integration)
	}
	return responses
}

func NewPaginatedIntegrationsResponse(integrations []*models.Integration, page, perPage, totalRecords int) *PaginatedIntegrationsResponse {
	return &PaginatedIntegrationsResponse{
		Results:  NewIntegrationResponses(integrations),
		PageInfo: GetPageInfo(page, perPage, len(integrations), totalRecords),
	}
}

type IntegrationSetupResponse struct {
	AuthorizeURL string `json:"authorize_url"`
}

// IntegrationTypeResponse represents metadata for a supported integration type.
// Endpoints and webhooks are intentionally omitted.
type IntegrationTypeResponse struct {
	Type        enums.IntegrationTypeEnum `json:"type"`
	Name        string                    `json:"name"`
	Icon        string                    `json:"icon"`
	Description string                    `json:"description"`
}
