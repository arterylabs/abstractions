package handlers

import (
	"abstractions/config"
	"abstractions/db/models"
	"abstractions/logger"
	"abstractions/server"
	"abstractions/server/requests"
	"abstractions/server/responses"
	integration_service "abstractions/services/integrations"
	github_integration "abstractions/services/integrations/core/github"
	slack_integration "abstractions/services/integrations/core/slack"
	"abstractions/services/integrations/core/slack/slack_client"
	"abstractions/services/integrations/enums"
	"abstractions/utils/context_utils"
	"abstractions/utils/pagination_utils"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type IntegrationHandler struct {
	*BaseHandler
}

func NewIntegrationHandler(server *server.Server) *IntegrationHandler {
	return &IntegrationHandler{
		BaseHandler: NewBaseHandler(server),
	}
}

// IntegrationSetup godoc
//
//	@Summary		Setup integration
//	@Description	Returns OAuth authorization URL for integration setup
//	@ID			integration-setup
//	@Tags			Integrations
//	@Accept			json
//	@Produce		json
//	@Param			workspace_id	path		string	true	"Workspace ID"
//	@Param			type		path		string	true	"Integration type"
//	@Param			params		body		requests.IntegrationSetupRequest	false	"Setup parameters"
//	@Success		200		{object}	responses.IntegrationSetupResponse
//	@Failure		400		{object}	responses.Error
//	@Failure		404		{object}	responses.Error
//	@Router			/workspaces/{workspace_id}/integrations/{type}/setup [post]
func (h *IntegrationHandler) IntegrationSetup(c echo.Context) error {
	workspaceIDStr := c.Param("workspace_id")
	integrationType, err := enums.FromString(c.Param("type"))
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid integration type")
	}

	// Bind and validate request body
	req := new(requests.IntegrationSetupRequest)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Parse workspace ID
	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid workspace ID")
	}

	// Validate integration type
	var oauthConfig *config.OAuthConfig
	switch integrationType {
	case slack_integration.IntegrationTypeSlack:
		oauthConfig = config.OAuthProviders.Slack
	case github_integration.IntegrationTypeGitHub:
		oauthConfig = config.OAuthProviders.GitHub
	default:
		return responses.ErrorResponse(c, http.StatusBadRequest, "Unsupported integration type")
	}

	if oauthConfig == nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "Integration configuration not found")
	}

	// Verify user has access to the workspace
	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "User not found")
	}

	workspaceService := h.Services(c).Workspace
	_, err = workspaceService.GetWorkspace(c.Request().Context(), workspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusForbidden, "Access denied to workspace")
	}

	// Create OAuth state with integration metadata
	redirectTo := "/integrations" // Default fallback
	if req.RedirectTo != nil && *req.RedirectTo != "" {
		redirectTo = *req.RedirectTo
	}

	// When connecting a new repository, override the redirect to the
	// repository connection page regardless of what was passed.
	if req.CreateRepository != nil && *req.CreateRepository {
		redirectTo = "/repositories/connect"
	}

	state := config.NewOAuthState(integrationType, redirectTo, &workspaceID, &user.ID, nil)

	// Attach resource context when provided (e.g. repository-level integrations)
	if req.ResourceType != nil && req.ResourceID != nil {
		resourceUUID, err := uuid.Parse(*req.ResourceID)
		if err != nil {
			return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid resource ID")
		}
		state.ResourceType = req.ResourceType
		state.ResourceID = &resourceUUID
	}

	// Propagate the createRepository flag so the callback knows to redirect
	// to the repository connection flow after the integration is created.
	if req.CreateRepository != nil {
		state.CreateRepository = req.CreateRepository
	}

	authURL := oauthConfig.AuthCodeURL(state.ToString())

	return responses.Response(c, http.StatusOK, responses.IntegrationSetupResponse{
		AuthorizeURL: authURL,
	})
}

// IntegrationCallback godoc
//
//	@Summary		Integration OAuth callback
//	@Description	Handles OAuth callback and creates integration
//	@ID			integration-callback
//	@Tags			Integrations
//	@Accept			json
//	@Produce		json
//	@Param			type	path		string	true	"Integration type"
//	@Param			code	query		string	true	"OAuth authorization code"
//	@Param			state	query		string	true	"OAuth state"
//	@Success		302		{string}	string	"Redirect to frontend"
//	@Failure		400		{object}	responses.Error
//	@Router			/integrations/callback/{type} [get]
func (h *IntegrationHandler) IntegrationCallback(c echo.Context) error {
	// Default redirect URL used when state cannot be parsed
	defaultRedirect := func(errorParam string) string {
		return config.Environment.ClientURL + "/integrations" + "?" + errorParam
	}

	integrationType, err := enums.FromString(c.Param("type"))
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, defaultRedirect("error=invalid_type"))
	}
	code := c.QueryParam("code")
	stateStr := c.QueryParam("state")

	if code == "" || stateStr == "" {
		return c.Redirect(http.StatusTemporaryRedirect, defaultRedirect("error=missing_parameters"))
	}

	// Parse OAuth state
	oauthState, err := config.ParseOAuthState(stateStr)
	if err != nil {
		logger.Error(c.Request().Context(), "Invalid OAuth state", err)
		return c.Redirect(http.StatusTemporaryRedirect, defaultRedirect("error=invalid_state"))
	}

	if oauthState.WorkspaceID == nil || oauthState.UserID == nil {
		logger.Error(c.Request().Context(), "OAuth state missing workspace or user ID")
		return c.Redirect(http.StatusTemporaryRedirect, defaultRedirect("error=invalid_state"))
	}

	// Helper to build redirect URL using state.RedirectTo
	getRedirectURL := func(errorParam string) string {
		redirectPath := oauthState.RedirectTo
		if redirectPath == "" {
			redirectPath = "/integrations"
		}

		unescapedRedirectTo, err := url.QueryUnescape(redirectPath)
		if err != nil || unescapedRedirectTo == "" {
			unescapedRedirectTo = "/integrations"
		}
		return config.Environment.ClientURL + unescapedRedirectTo + "?" + errorParam
	}

	// Validate integration type
	var oauthConfig *config.OAuthConfig
	switch integrationType {
	case slack_integration.IntegrationTypeSlack:
		oauthConfig = config.OAuthProviders.Slack
	case github_integration.IntegrationTypeGitHub:
		oauthConfig = config.OAuthProviders.GitHub
	default:
		return c.Redirect(http.StatusTemporaryRedirect, getRedirectURL("error=unsupported_type"))
	}

	// Exchange code for token
	token, err := oauthConfig.Exchange(c.Request().Context(), code)
	if err != nil {
		logger.Error(c.Request().Context(), "Failed to exchange code for token", err)
		return c.Redirect(http.StatusTemporaryRedirect, getRedirectURL("error=token_exchange_failed"))
	}

	// Get integration-specific data based on type
	var serviceTitle, serviceIcon, serviceTenantID, serviceUserID *string
	var configuration models.JSONBField

	switch integrationType {
	case slack_integration.IntegrationTypeSlack:
		// Get Slack team info
		slackData, err := slack_client.GetTeamInfo(token.AccessToken)
		if err != nil {
			logger.Error(c.Request().Context(), "Failed to get Slack team info", err)
		} else {
			serviceTitle = &slackData.Team.Name
			serviceIcon = slackData.Team.Icon.Image132
			serviceTenantID = &slackData.Team.ID
			slackAuthTest, err := slack_client.AuthTest(token.AccessToken)
			if err != nil {
				logger.Error(c.Request().Context(), "Failed to get Slack auth test", err)
			} else {
				serviceUserID = &slackAuthTest.UserID
			}
		}
	case github_integration.IntegrationTypeGitHub:
		// Fetch GitHub user info to populate service metadata
		githubUserInfo, err := fetchGitHubUserInfo(token.AccessToken)
		if err != nil {
			logger.Error(c.Request().Context(), "Failed to get GitHub user info", err)
		} else {
			serviceTitle = &githubUserInfo.Login
			if githubUserInfo.AvatarURL != "" {
				serviceIcon = &githubUserInfo.AvatarURL
			}
			idStr := fmt.Sprintf("%d", githubUserInfo.ID)
			serviceTenantID = &idStr
			serviceUserID = &githubUserInfo.Login
		}
	}

	// Create integration via service (service handles encryption)
	integrationSvc := h.Services(c).Integration
	var scopesPtr *string
	if len(oauthConfig.Scopes) > 0 {
		s := strings.Join(oauthConfig.Scopes, ",")
		scopesPtr = &s
	}

	createReq := &integration_service.CreateIntegrationRequest{
		WorkspaceID:  *oauthState.WorkspaceID,
		Type:         integrationType,
		ResourceType: oauthState.ResourceType,
		ResourceID:   oauthState.ResourceID,
		ServiceTitle: serviceTitle,
		ServiceIcon:  serviceIcon,
		Token:        token.AccessToken,
		RefreshToken: func() *string {
			if token.RefreshToken == "" {
				return nil
			}
			v := token.RefreshToken
			return &v
		}(),
		TokenExpiresAt:  &token.Expiry,
		Scopes:          scopesPtr,
		ServiceTenantID: serviceTenantID,
		ServiceUserID:   serviceUserID,
		Configuration:   configuration,
		CreatedBy:       *oauthState.UserID,
	}

	_, err = integrationSvc.CreateIntegration(c.Request().Context(), createReq)
	if err != nil {
		logger.Error(c.Request().Context(), "Failed to create integration", err)
		return c.Redirect(http.StatusTemporaryRedirect, getRedirectURL("error=creation_failed"))
	}

	// Success redirect — when createRepository is set, always send the user to
	// the repository connection page so they can pick a repo to connect.
	redirectPath := oauthState.RedirectTo
	if oauthState.CreateRepository != nil && *oauthState.CreateRepository {
		redirectPath = "/repositories/connect"
	}
	if redirectPath == "" {
		redirectPath = "/integrations"
	}
	unescapedRedirectTo, err := url.QueryUnescape(redirectPath)
	if err != nil || unescapedRedirectTo == "" {
		logger.Error(c.Request().Context(), "Failed to unescape redirectTo", err)
		unescapedRedirectTo = "/integrations"
	}

	return c.Redirect(http.StatusTemporaryRedirect,
		config.Environment.ClientURL+unescapedRedirectTo+"?success=true&type="+integrationType.String())
}

// GetIntegrations godoc
//
//	@Summary		Get workspace integrations
//	@Description	Get paginated list of integrations for a workspace
//	@ID			get-integrations
//	@Tags			Integrations
//	@Accept			json
//	@Produce		json
//	@Param			workspace_id	path		string	true	"Workspace ID"
//	@Param			page			query		int		false	"Page number"	default(1)
//	@Param			per_page		query		int		false	"Items per page"	default(10)
//	@Param			sort			query		string	false	"Sort field"
//	@Param			sort_order		query		string	false	"Sort order (asc/desc)"	default(desc)
//	@Param			query			query		string	false	"Search query"
//	@Param			resource_type	query		string	false	"Filter by resource type (e.g. repositories)"
//	@Param			resource_id		query		string	false	"Filter by resource ID (UUID)"
//	@Success		200			{object}	responses.PaginatedIntegrationsResponse
//	@Failure		400			{object}	responses.Error
//	@Failure		403			{object}	responses.Error
//	@Router			/workspaces/{workspace_id}/integrations [get]
func (h *IntegrationHandler) GetIntegrations(c echo.Context) error {
	workspaceIDStr := c.Param("workspace_id")
	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid workspace ID")
	}

	// Parse pagination parameters
	paginationParams := pagination_utils.PaginationParams{}
	if err := c.Bind(&paginationParams); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	paginationParams.Validate()

	// Optional resource filter
	var resourceType *string
	var resourceID *uuid.UUID
	if rt := c.QueryParam("resource_type"); rt != "" {
		resourceType = &rt
	}
	if ridStr := c.QueryParam("resource_id"); ridStr != "" {
		rid, err := uuid.Parse(ridStr)
		if err != nil {
			return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid resource ID")
		}
		resourceID = &rid
	}

	integrationService := h.Services(c).Integration
	var integrations []*models.Integration
	var total int
	if resourceType != nil || resourceID != nil {
		integrations, total, err = integrationService.GetIntegrationsPaginatedByResource(c.Request().Context(), workspaceID, paginationParams, resourceType, resourceID)
	} else {
		integrations, total, err = integrationService.GetIntegrationsPaginated(c.Request().Context(), workspaceID, paginationParams)
	}
	if err != nil {
		logger.Error(c.Request().Context(), "Failed to get integrations", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve integrations")
	}

	return responses.Response(c, http.StatusOK, responses.NewPaginatedIntegrationsResponse(integrations, paginationParams.Page, paginationParams.PerPage, total))
}

// DeleteIntegration godoc
//
//	@Summary		Delete integration
//	@Description	Delete an integration by ID
//	@ID			delete-integration
//	@Tags			Integrations
//	@Accept			json
//	@Produce		json
//	@Param			workspace_id	path		string	true	"Workspace ID"
//	@Param			integration_id	path		string	true	"Integration ID"
//	@Success		200			{object}	responses.Data
//	@Failure		400			{object}	responses.Error
//	@Failure		403			{object}	responses.Error
//	@Failure		404			{object}	responses.Error
//	@Router			/workspaces/{workspace_id}/integrations/{integration_id} [delete]
func (h *IntegrationHandler) DeleteIntegration(c echo.Context) error {
	workspaceIDStr := c.Param("workspace_id")
	integrationIDStr := c.Param("integration_id")

	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid workspace ID")
	}

	integrationID, err := uuid.Parse(integrationIDStr)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid integration ID")
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "User not found")
	}

	integrationService := h.Services(c).Integration

	// Verify integration exists and belongs to workspace
	integration, err := integrationService.GetIntegration(c.Request().Context(), integrationID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "Integration not found")
	}

	if integration.WorkspaceID != workspaceID {
		return responses.ErrorResponse(c, http.StatusForbidden, "Integration does not belong to this workspace")
	}

	err = integrationService.DeleteIntegration(c.Request().Context(), integrationID, user.ID)
	if err != nil {
		logger.Error(c.Request().Context(), "Failed to delete integration", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete integration")
	}

	return responses.Response(c, http.StatusOK, map[string]interface{}{
		"message": "Integration deleted successfully",
	})
}

// UpdateIntegration godoc
//
//	@Summary		Update integration
//	@Description	Update integration configuration
//	@ID			update-integration
//	@Tags			Integrations
//	@Accept			json
//	@Produce		json
//	@Param			workspace_id	path		string				true	"Workspace ID"
//	@Param			integration_id	path		string				true	"Integration ID"
//	@Param			params			body		requests.UpdateIntegrationRequest	true	"Integration update parameters"
//	@Success		200			{object}	responses.Data
//	@Failure		400			{object}	responses.Error
//	@Failure		403			{object}	responses.Error
//	@Failure		404			{object}	responses.Error
//	@Router			/workspaces/{workspace_id}/integrations/{integration_id} [put]
func (h *IntegrationHandler) UpdateIntegration(c echo.Context) error {
	workspaceIDStr := c.Param("workspace_id")
	integrationIDStr := c.Param("integration_id")

	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid workspace ID")
	}

	integrationID, err := uuid.Parse(integrationIDStr)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid integration ID")
	}

	req := new(requests.UpdateIntegrationRequest)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "User not found")
	}

	integrationService := h.Services(c).Integration

	// Verify integration exists and belongs to workspace
	integration, err := integrationService.GetIntegration(c.Request().Context(), integrationID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "Integration not found")
	}

	if integration.WorkspaceID != workspaceID {
		return responses.ErrorResponse(c, http.StatusForbidden, "Integration does not belong to this workspace")
	}

	// Update integration
	if req.ServiceTitle != nil {
		integration.ServiceTitle = req.ServiceTitle
	}
	if req.Configuration != nil {
		integration.Configuration = req.Configuration
	}
	integration.UpdatedBy = user.ID
	integration.UpdatedAt = time.Now().UTC()

	err = integrationService.UpdateIntegration(c.Request().Context(), integration)
	if err != nil {
		logger.Error(c.Request().Context(), "Failed to update integration", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Failed to update integration")
	}

	return responses.Response(c, http.StatusOK, map[string]interface{}{
		"message": "Integration updated successfully",
	})
}

// GetIntegrationByType godoc
//
//	@Summary		Get integration by type
//	@Description	Get integration by type for a workspace
//	@ID			get-integration-by-type
//	@Tags			Integrations
//	@Accept			json
//	@Produce		json
//	@Param			workspace_id	path		string	true	"Workspace ID"
//	@Param			type			path		string	true	"Integration type (e.g., slack)"
//	@Success		200			{object}	responses.IntegrationResponse
//	@Failure		400			{object}	responses.Error
//	@Failure		403			{object}	responses.Error
//	@Failure		404			{object}	responses.Error
//	@Router			/workspaces/{workspace_id}/integrations/{type} [get]
func (h *IntegrationHandler) GetIntegrationByType(c echo.Context) error {
	workspaceIDStr := c.Param("workspace_id")
	typeStr := c.Param("type")

	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid workspace ID")
	}

	integrationType, err := enums.FromString(typeStr)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid integration type")
	}

	integrationService := h.Services(c).Integration
	integration, err := integrationService.GetIntegrationByType(c.Request().Context(), workspaceID, integrationType)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "Integration not found for this type and workspace")
	}

	return responses.Response(c, http.StatusOK, responses.NewIntegrationResponse(integration))
}

// GetIntegrationTypes godoc
//
//	@Summary		List supported integration types
//	@Description	Returns all supported integration types in the product
//	@ID			get-integration-types
//	@Tags			Integrations
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	responses.IntegrationTypeResponse
//	@Router			/integrations/types [get]
func (h *IntegrationHandler) GetIntegrationTypes(c echo.Context) error {
	// Use service layer registry but only expose limited fields (no endpoints/webhooks)
	types := integration_service.GetAllIntegrationTypes()

	results := make([]responses.IntegrationTypeResponse, 0, len(types))
	for _, t := range types {
		results = append(results, responses.IntegrationTypeResponse{
			Type:        t.GetIntegrationType(),
			Name:        t.GetName(),
			Icon:        t.GetIcon(),
			Description: t.GetDescription(),
		})
	}

	return responses.Response(c, http.StatusOK, results)
}

// githubUserInfo holds the subset of fields returned by the GitHub /user endpoint.
type githubUserInfo struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
}

// fetchGitHubUserInfo calls the GitHub API to retrieve the authenticated user's profile.
func fetchGitHubUserInfo(accessToken string) (*githubUserInfo, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info githubUserInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}
	return &info, nil
}
