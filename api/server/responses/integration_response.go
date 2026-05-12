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
	ServiceTitle    *string                   `json:"service_title"`
	ServiceIcon     *string                   `json:"service_icon"`
	WorkspaceID     string                    `json:"workspace_id"`
	TokenExpiresAt  *string                   `json:"token_expires_at"`
	Scopes          *string                   `json:"scopes"`
	ServiceTenantID *string                   `json:"service_tenant_id"`
	ServiceUserID   *string                   `json:"service_user_id"`
	ResourceType    *string                   `json:"resourceType"`
	ResourceID      *string                   `json:"resourceId"`
	Configuration   map[string]interface{}    `json:"configuration"`
	CreatedAt       string                    `json:"created_at"`
	UpdatedAt       string                    `json:"updated_at"`
}

func NewIntegrationResponse(integration *models.Integration) IntegrationResponse {
	var tokenExpiresAt *string
	if integration.TokenExpiresAt != nil {
		formatted := integration.TokenExpiresAt.UTC().Format(time.RFC3339)
		tokenExpiresAt = &formatted
	}

	return IntegrationResponse{
		ID:              integration.ID.String(),
		Type:            integration.Type,
		ServiceTitle:    integration.ServiceTitle,
		ServiceIcon:     integration.ServiceIcon,
		WorkspaceID:     integration.WorkspaceID.String(),
		TokenExpiresAt:  tokenExpiresAt,
		Scopes:          integration.Scopes,
		ServiceTenantID: integration.ServiceTenantID,
		ServiceUserID:   integration.ServiceUserID,
		ResourceType:    integration.ResourceType,
		ResourceID: func() *string {
			if integration.ResourceID == nil {
				return nil
			}
			s := integration.ResourceID.String()
			return &s
		}(),
		Configuration: integration.Configuration,
		CreatedAt:     integration.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     integration.UpdatedAt.UTC().Format(time.RFC3339),
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
