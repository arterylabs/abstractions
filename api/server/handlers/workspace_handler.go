package handlers

import (
	ai_client "abstractions/clients/ai"
	"abstractions/clients/storage"
	"abstractions/server"
	"abstractions/server/requests"
	"abstractions/server/responses"
	github_integration "abstractions/services/integrations/core/github"
	workspace_service "abstractions/services/workspace"
	workspace_configuration_service "abstractions/services/workspace_configuration"
	"abstractions/utils/context_utils"
	"abstractions/utils/pagination_utils"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type WorkspaceHandler struct {
	*BaseHandler
}

func NewWorkspaceHandler(server *server.Server) *WorkspaceHandler {
	return &WorkspaceHandler{
		BaseHandler: NewBaseHandler(server),
	}
}

// GetUserWorkspaces returns a list of workspaces that the authenticated user has access to
func (h *WorkspaceHandler) GetUserWorkspaces(c echo.Context) error {
	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return h.internalError(c, err)
	}
	ctx := c.Request().Context()
	workspaces, err := h.Services(c).Workspace.GetUserWorkspacesWithRole(ctx, user.ID)
	if err != nil {
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusOK, responses.NewWorkspaceWithRoleResponses(workspaces))
}

// CreateWorkspace creates a new workspace with the authenticated user as the owner
func (h *WorkspaceHandler) CreateWorkspace(c echo.Context) error {
	req := new(requests.CreateWorkspace)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return h.internalError(c, err)
	}

	ctx := c.Request().Context()

	workspace, err := h.Services(c).Workspace.Create(ctx, req.Name, user.ID)
	if err != nil {
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusCreated, responses.NewWorkspaceResponse(workspace))
}

// UpdateWorkspace updates an existing workspace
func (h *WorkspaceHandler) UpdateWorkspace(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}

	req := new(requests.UpdateWorkspace)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return h.internalError(c, err)
	}

	ctx := c.Request().Context()

	params := workspace_service.UpdateWorkspaceParams{
		Name:         req.Name,
		Logo:         req.Logo,
		Website:      req.Website,
		About:        req.About,
		BillingEmail: req.BillingEmail,
	}
	if err := h.Services(c).Workspace.UpdateWorkspace(ctx, workspaceUUID, params, user.ID); err != nil {
		if err == workspace_service.ErrWorkspaceNotFound {
			return responses.ErrorResponse(c, http.StatusNotFound, "workspace not found")
		}
		return h.internalError(c, err)
	}

	// Get the updated workspace
	workspace, err := h.Services(c).Workspace.GetWorkspace(ctx, workspaceUUID)
	if err != nil {
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusOK, responses.NewWorkspaceResponse(workspace))
}

// DeleteWorkspace soft deletes a workspace by setting its deleted_at timestamp
func (h *WorkspaceHandler) DeleteWorkspace(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return h.internalError(c, err)
	}

	ctx := c.Request().Context()

	if err := h.Services(c).Workspace.Delete(ctx, workspaceUUID, user.ID); err != nil {
		if err == workspace_service.ErrWorkspaceNotFound {
			return responses.ErrorResponse(c, http.StatusNotFound, "workspace not found")
		}
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusNoContent, nil)
}

// GetWorkspaceUsers returns a list of users in the workspace
func (h *WorkspaceHandler) GetWorkspaceUsers(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}

	paginationParams := pagination_utils.PaginationParams{}
	if err := c.Bind(&paginationParams); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	paginationParams.Validate()

	ctx := c.Request().Context()

	users, total, err := h.Services(c).Workspace.GetWorkspaceUsersPaginated(ctx, workspaceUUID, paginationParams)
	if err != nil {
		return h.internalError(c, err)
	}

	response := responses.NewPaginatedWorkspaceUsersResponse(users, paginationParams.Page, paginationParams.PerPage, total)

	return responses.Response(c, http.StatusOK, response)
}

// RemoveWorkspaceUser removes a user from the workspace
func (h *WorkspaceHandler) RemoveWorkspaceUser(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}

	req := new(requests.RemoveWorkspaceUser)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return h.internalError(c, err)
	}

	ctx := c.Request().Context()

	if err := h.Services(c).Workspace.RemoveUserFromWorkspace(ctx, workspaceUUID, req.UserID, user.ID); err != nil {
		switch err {
		case workspace_service.ErrInsufficientPermissions:
			return responses.ErrorResponse(c, http.StatusForbidden, "You don't have permission to remove this user")
		case workspace_service.ErrCannotModifyOwnRole:
			return responses.ErrorResponse(c, http.StatusBadRequest, "You cannot remove yourself from the workspace")
		default:
			return h.internalError(c, err)
		}
	}

	return responses.Response(c, http.StatusNoContent, nil)
}

// UpdateWorkspaceUserRole updates the role of a user in the workspace
func (h *WorkspaceHandler) UpdateWorkspaceUserRole(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}

	req := new(requests.UpdateWorkspaceUserRole)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return h.internalError(c, err)
	}

	ctx := c.Request().Context()

	if err := h.Services(c).Workspace.UpdateUserRole(ctx, workspaceUUID, req.UserID, req.Role, user.ID); err != nil {
		switch err {
		case workspace_service.ErrInsufficientPermissions:
			return responses.ErrorResponse(c, http.StatusForbidden, "You don't have permission to assign this role or modify this user")
		case workspace_service.ErrCannotModifyOwnRole:
			return responses.ErrorResponse(c, http.StatusBadRequest, "You cannot modify your own role")
		default:
			return h.internalError(c, err)
		}
	}

	return responses.Response(c, http.StatusOK, nil)
}

// UploadLogo uploads a new logo for the workspace
func (h *WorkspaceHandler) UploadLogo(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "file is required")
	}
	src, err := file.Open()
	if err != nil {
		return h.internalError(c, err)
	}
	defer src.Close()

	ctx := c.Request().Context()

	// Get workspace first to verify it exists
	workspace, err := h.Services(c).Workspace.GetWorkspace(ctx, workspaceUUID)
	if err != nil {
		return h.internalError(c, err)
	}
	if workspace == nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "workspace not found")
	}

	// Upload the logo
	workspace, err = h.Services(c).Workspace.UpdateWorkspaceLogo(ctx, workspaceUUID, src)
	if err != nil {
		if errors.Is(err, storage.ErrFileTooLarge) || errors.Is(err, storage.ErrMIMENotAllowed) {
			return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		return responses.ErrorResponse(c, http.StatusInternalServerError, "upload failed")
	}

	return responses.Response(c, http.StatusOK, responses.NewWorkspaceResponse(workspace))
}

func (h *WorkspaceHandler) GetRoles(c echo.Context) error {

	// this endpoint is scoped to the workspace because in future we might allow custom roles per workspace. right now it just returns the default roles
	return responses.Response(c, http.StatusOK, h.Services(c).Authorization.GetAllRoles(c.Request().Context()))
}

// GetWorkspaceAISettings returns the AI configuration for the workspace with secret values masked.
func (h *WorkspaceHandler) GetWorkspaceAISettings(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}

	ctx := c.Request().Context()
	cfgSvc := h.Services(c).WorkspaceConfiguration

	// Collect the masked last-four of every provider's key, keyed by provider ID.
	apiKeyLastFour := make(map[string]*string, len(ai_client.Providers))
	for _, provider := range ai_client.Providers {
		lastFour, err := cfgSvc.GetLastFour(ctx, workspaceUUID, provider.ConfigKey)
		if err != nil {
			return h.internalError(c, err)
		}
		apiKeyLastFour[string(provider.ID)] = lastFour
	}

	indexingModel, err := cfgSvc.Get(ctx, workspaceUUID, workspace_configuration_service.KeyIndexingModel)
	if err != nil {
		return h.internalError(c, err)
	}
	newsletterModel, err := cfgSvc.Get(ctx, workspaceUUID, workspace_configuration_service.KeyNewsletterModel)
	if err != nil {
		return h.internalError(c, err)
	}
	embeddingsModel, err := cfgSvc.Get(ctx, workspaceUUID, workspace_configuration_service.KeyEmbeddingsModel)
	if err != nil {
		return h.internalError(c, err)
	}

	summarizeCodePrompt, err := cfgSvc.Get(ctx, workspaceUUID, workspace_configuration_service.KeySystemPromptSummarizeCode)
	if err != nil {
		return h.internalError(c, err)
	}
	summarizeFilePrompt, err := cfgSvc.Get(ctx, workspaceUUID, workspace_configuration_service.KeySystemPromptSummarizeFile)
	if err != nil {
		return h.internalError(c, err)
	}
	topicGeneratorPrompt, err := cfgSvc.Get(ctx, workspaceUUID, workspace_configuration_service.KeySystemPromptTopicGenerator)
	if err != nil {
		return h.internalError(c, err)
	}
	newsletterEntryPrompt, err := cfgSvc.Get(ctx, workspaceUUID, workspace_configuration_service.KeySystemPromptNewsletterEntry)
	if err != nil {
		return h.internalError(c, err)
	}

	indexingModelStr := indexingModel
	newsletterModelStr := newsletterModel
	embeddingsModelStr := embeddingsModel
	defaults := ai_client.DefaultSystemPrompts()

	return responses.Response(c, http.StatusOK, responses.WorkspaceAISettingsResponse{
		APIKeyLastFour:  apiKeyLastFour,
		IndexingModel:   &indexingModelStr,
		NewsletterModel: &newsletterModelStr,
		EmbeddingsModel: &embeddingsModelStr,
		SystemPrompts: responses.SystemPromptsResponse{
			SummarizeCode:   optionalNonEmptyString(summarizeCodePrompt),
			SummarizeFile:   optionalNonEmptyString(summarizeFilePrompt),
			TopicGenerator:  optionalNonEmptyString(topicGeneratorPrompt),
			NewsletterEntry: optionalNonEmptyString(newsletterEntryPrompt),
		},
		DefaultSystemPrompts: responses.DefaultSystemPromptsResponse{
			SummarizeCode:   defaults.SummarizeCode,
			SummarizeFile:   defaults.SummarizeFile,
			TopicGenerator:  defaults.TopicGenerator,
			NewsletterEntry: defaults.NewsletterEntry,
		},
	})
}

// optionalNonEmptyString returns a pointer to s when s is non-empty after trim; otherwise nil.
func optionalNonEmptyString(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return &s
}

// GetAIModels returns the master list of supported AI models.
func (h *WorkspaceHandler) GetAIModels(c echo.Context) error {
	models := make([]responses.AIModelResponse, len(ai_client.Models))
	for i, m := range ai_client.Models {
		models[i] = responses.AIModelResponse{
			ID:                     m.ID,
			DisplayName:            m.DisplayName,
			Provider:               string(m.Provider),
			SupportsTextGeneration: m.SupportsTextGeneration,
			SupportsEmbeddings:     m.SupportsEmbeddings,
			PricePerMillionTokens:  m.PricePerMillionTokens,
			ContextWindow:          m.ContextWindow,
		}
	}
	return responses.Response(c, http.StatusOK, models)
}

// GetAIProviders returns the master list of supported AI providers so the frontend can render
// provider-specific UI (icons, key placeholders, key rows) without hardcoding provider identity.
func (h *WorkspaceHandler) GetAIProviders(c echo.Context) error {
	providers := make([]responses.AIProviderResponse, len(ai_client.Providers))
	for i, p := range ai_client.Providers {
		providers[i] = responses.AIProviderResponse{
			ID:                 string(p.ID),
			DisplayName:        p.DisplayName,
			IconPath:           p.IconPath,
			KeyPlaceholder:     p.KeyPlaceholder,
			SupportsEmbeddings: p.SupportsEmbeddings,
		}
	}
	return responses.Response(c, http.StatusOK, providers)
}

// UpdateWorkspaceAISettings updates AI configuration keys and model preferences for the workspace.
func (h *WorkspaceHandler) UpdateWorkspaceAISettings(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}

	req := new(requests.UpdateWorkspaceAISettings)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return h.internalError(c, err)
	}

	ctx := c.Request().Context()
	cfgSvc := h.Services(c).WorkspaceConfiguration

	type keyValue struct {
		key   string
		value *string
	}
	var fields []keyValue

	// Validate and stage each provider API key from the registry. A non-empty value is validated
	// against the provider before it is persisted; an empty value clears the stored key.
	for _, provider := range ai_client.Providers {
		value, provided := req.APIKeys[string(provider.ID)]
		if !provided || value == nil {
			continue
		}
		if *value != "" && provider.Validate != nil {
			if err := provider.Validate(*value, nil); err != nil {
				return responses.ErrorResponse(c, http.StatusUnprocessableEntity, "Invalid "+provider.DisplayName+" API key")
			}
		}
		fields = append(fields, keyValue{provider.ConfigKey, value})
	}

	fields = append(fields,
		keyValue{workspace_configuration_service.KeyIndexingModel, req.IndexingModel},
		keyValue{workspace_configuration_service.KeyNewsletterModel, req.NewsletterModel},
		keyValue{workspace_configuration_service.KeyEmbeddingsModel, req.EmbeddingsModel},
	)

	if req.SystemPrompts != nil {
		fields = append(fields,
			keyValue{workspace_configuration_service.KeySystemPromptSummarizeCode, req.SystemPrompts.SummarizeCode},
			keyValue{workspace_configuration_service.KeySystemPromptSummarizeFile, req.SystemPrompts.SummarizeFile},
			keyValue{workspace_configuration_service.KeySystemPromptTopicGenerator, req.SystemPrompts.TopicGenerator},
			keyValue{workspace_configuration_service.KeySystemPromptNewsletterEntry, req.SystemPrompts.NewsletterEntry},
		)
	}

	for _, f := range fields {
		if f.value == nil {
			continue
		}
		// Whitespace-only system prompt enhancements are treated as clear (default only).
		value := *f.value
		if strings.HasPrefix(f.key, "system_prompt_") && strings.TrimSpace(value) == "" {
			value = ""
		}
		if err := cfgSvc.Set(ctx, workspaceUUID, f.key, value, user.ID); err != nil {
			return h.internalError(c, err)
		}
	}

	// Return updated settings
	return h.GetWorkspaceAISettings(c)
}

// GetWorkspaceHealth returns the health status of third-party connections for the workspace.
func (h *WorkspaceHandler) GetWorkspaceHealth(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}

	ctx := c.Request().Context()
	cfgSvc := h.Services(c).WorkspaceConfiguration

	// Determine which providers are required by the workspace's configured models.
	requiredProviders := map[ai_client.Provider]bool{}
	for _, modelKey := range []string{
		workspace_configuration_service.KeyIndexingModel,
		workspace_configuration_service.KeyNewsletterModel,
		workspace_configuration_service.KeyEmbeddingsModel,
	} {
		modelID, err := cfgSvc.Get(ctx, workspaceUUID, modelKey)
		if err != nil {
			return h.internalError(c, err)
		}
		if info, ok := ai_client.ModelByID(modelID); ok {
			requiredProviders[info.Provider] = true
		}
	}

	// Report configured/valid/required status for every provider in the registry.
	providerHealth := make(map[string]responses.WorkspaceProviderHealth, len(ai_client.Providers))
	for _, provider := range ai_client.Providers {
		key, err := cfgSvc.Get(ctx, workspaceUUID, provider.ConfigKey)
		if err != nil {
			return h.internalError(c, err)
		}
		status := responses.WorkspaceProviderHealth{
			Configured: key != "",
			Required:   requiredProviders[provider.ID],
		}
		if status.Configured && provider.Validate != nil {
			if err := provider.Validate(key, nil); err != nil {
				status.Error = err.Error()
			} else {
				status.Valid = true
			}
		}
		providerHealth[string(provider.ID)] = status
	}

	return responses.Response(c, http.StatusOK, responses.WorkspaceHealthResponse{
		Providers: providerHealth,
	})
}

// GetGettingStarted returns the completion status of each onboarding step for the workspace.
func (h *WorkspaceHandler) GetGettingStarted(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "invalid workspace_id format")
	}

	ctx := c.Request().Context()
	svc := h.Services(c)

	// AI is considered set up once any provider has a key configured.
	hasKey := false
	for _, provider := range ai_client.Providers {
		if lastFour, _ := svc.WorkspaceConfiguration.GetLastFour(ctx, workspaceUUID, provider.ConfigKey); lastFour != nil {
			hasKey = true
			break
		}
	}

	indexingModel, _ := svc.WorkspaceConfiguration.Get(ctx, workspaceUUID, workspace_configuration_service.KeyIndexingModel)
	newsletterModel, _ := svc.WorkspaceConfiguration.Get(ctx, workspaceUUID, workspace_configuration_service.KeyNewsletterModel)
	hasModel := indexingModel != "" || newsletterModel != ""

	aiConfigured := hasKey && hasModel

	githubIntegration, _ := svc.Integration.GetIntegrationByType(ctx, workspaceUUID, github_integration.IntegrationTypeGitHub)

	repos, _ := svc.Repository.ListRepositories(ctx, workspaceUUID)
	githubConnected := githubIntegration != nil && len(repos) > 0
	newsletterSent := false
	for _, repo := range repos {
		if repo.NewsletterSentAt != nil {
			newsletterSent = true
			break
		}
	}

	return responses.Response(c, http.StatusOK, responses.WorkspaceGettingStartedResponse{
		AIConfigured:    aiConfigured,
		GitHubConnected: githubConnected,
		NewsletterSent:  newsletterSent,
	})
}
