package handlers

import (
	ai_client "abstractions/clients/ai"
	"abstractions/server"
	"abstractions/server/requests"
	"abstractions/server/responses"
	workspace_service "abstractions/services/workspace"
	workspace_configuration_service "abstractions/services/workspace_configuration"
	"abstractions/utils/context_utils"
	"abstractions/utils/pagination_utils"
	"net/http"

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
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	ctx := c.Request().Context()
	workspaces, err := h.Services(c).Workspace.GetUserWorkspacesWithRole(ctx, user.ID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
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
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	ctx := c.Request().Context()

	workspace, err := h.Services(c).Workspace.CreateWorkspace(ctx, req.Name, user.ID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
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
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
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
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	// Get the updated workspace
	workspace, err := h.Services(c).Workspace.GetWorkspace(ctx, workspaceUUID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
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
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	ctx := c.Request().Context()

	if err := h.Services(c).Workspace.DeleteWorkspace(ctx, workspaceUUID, user.ID); err != nil {
		if err == workspace_service.ErrWorkspaceNotFound {
			return responses.ErrorResponse(c, http.StatusNotFound, "workspace not found")
		}
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
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
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
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
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	ctx := c.Request().Context()

	if err := h.Services(c).Workspace.RemoveUserFromWorkspace(ctx, workspaceUUID, req.UserID, user.ID); err != nil {
		switch err {
		case workspace_service.ErrInsufficientPermissions:
			return responses.ErrorResponse(c, http.StatusForbidden, "You don't have permission to remove this user")
		case workspace_service.ErrCannotModifyOwnRole:
			return responses.ErrorResponse(c, http.StatusBadRequest, "You cannot remove yourself from the workspace")
		default:
			return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
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
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	ctx := c.Request().Context()

	if err := h.Services(c).Workspace.UpdateUserRole(ctx, workspaceUUID, req.UserID, req.Role, user.ID); err != nil {
		switch err {
		case workspace_service.ErrInsufficientPermissions:
			return responses.ErrorResponse(c, http.StatusForbidden, "You don't have permission to assign this role or modify this user")
		case workspace_service.ErrCannotModifyOwnRole:
			return responses.ErrorResponse(c, http.StatusBadRequest, "You cannot modify your own role")
		default:
			return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
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
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	ctx := c.Request().Context()

	// Get workspace first to verify it exists
	workspace, err := h.Services(c).Workspace.GetWorkspace(ctx, workspaceUUID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	if workspace == nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "workspace not found")
	}

	// Upload the logo
	workspace, err = h.Services(c).Workspace.UpdateWorkspaceLogo(ctx, workspaceUUID, src)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, responses.NewWorkspaceResponse(workspace))
}

func (h *WorkspaceHandler) GetRoles(c echo.Context) error {

	// this endpoint is scoped to the workspace because in future we might allow custom roles per workspace. right now it just returns the default roles
	return responses.Response(c, http.StatusOK, h.Services(c).Authorization.GetAllRoles())
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

	anthropicLastFour, err := cfgSvc.GetLastFour(ctx, workspaceUUID, workspace_configuration_service.KeyAnthropicAPIKey)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	openAILastFour, err := cfgSvc.GetLastFour(ctx, workspaceUUID, workspace_configuration_service.KeyOpenAIAPIKey)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	indexingModel, err := cfgSvc.Get(ctx, workspaceUUID, workspace_configuration_service.KeyIndexingModel)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	newsletterModel, err := cfgSvc.Get(ctx, workspaceUUID, workspace_configuration_service.KeyNewsletterModel)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	embeddingsModel, err := cfgSvc.Get(ctx, workspaceUUID, workspace_configuration_service.KeyEmbeddingsModel)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	indexingModelStr := indexingModel
	newsletterModelStr := newsletterModel
	embeddingsModelStr := embeddingsModel

	return responses.Response(c, http.StatusOK, responses.WorkspaceAISettingsResponse{
		AnthropicAPIKeyLastFour: anthropicLastFour,
		OpenAIAPIKeyLastFour:    openAILastFour,
		IndexingModel:           &indexingModelStr,
		NewsletterModel:         &newsletterModelStr,
		EmbeddingsModel:         &embeddingsModelStr,
	})
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
		}
	}
	return responses.Response(c, http.StatusOK, models)
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

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	ctx := c.Request().Context()
	cfgSvc := h.Services(c).WorkspaceConfiguration

	if req.AnthropicAPIKey != nil && *req.AnthropicAPIKey != "" {
		if err := ai_client.ValidateAnthropicAPIKey(*req.AnthropicAPIKey, nil); err != nil {
			return responses.ErrorResponse(c, http.StatusUnprocessableEntity, "Invalid Anthropic API key")
		}
	}
	if req.OpenAIAPIKey != nil && *req.OpenAIAPIKey != "" {
		if err := ai_client.ValidateOpenAIAPIKey(*req.OpenAIAPIKey, nil); err != nil {
			return responses.ErrorResponse(c, http.StatusUnprocessableEntity, "Invalid OpenAI API key")
		}
	}

	type keyValue struct {
		key   string
		value *string
	}
	fields := []keyValue{
		{workspace_configuration_service.KeyAnthropicAPIKey, req.AnthropicAPIKey},
		{workspace_configuration_service.KeyOpenAIAPIKey, req.OpenAIAPIKey},
		{workspace_configuration_service.KeyIndexingModel, req.IndexingModel},
		{workspace_configuration_service.KeyNewsletterModel, req.NewsletterModel},
		{workspace_configuration_service.KeyEmbeddingsModel, req.EmbeddingsModel},
	}
	for _, f := range fields {
		if f.value == nil {
			continue
		}
		if err := cfgSvc.Set(ctx, workspaceUUID, f.key, *f.value, user.ID); err != nil {
			return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
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

	anthropicKey, err := cfgSvc.Get(ctx, workspaceUUID, workspace_configuration_service.KeyAnthropicAPIKey)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	openAIKey, err := cfgSvc.Get(ctx, workspaceUUID, workspace_configuration_service.KeyOpenAIAPIKey)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	anthropicStatus := responses.WorkspaceConnectionStatus{Configured: anthropicKey != ""}
	if anthropicStatus.Configured {
		if err := ai_client.ValidateAnthropicAPIKey(anthropicKey, nil); err != nil {
			anthropicStatus.Error = err.Error()
		} else {
			anthropicStatus.Valid = true
		}
	}

	openAIStatus := responses.WorkspaceConnectionStatus{Configured: openAIKey != ""}
	if openAIStatus.Configured {
		if err := ai_client.ValidateOpenAIAPIKey(openAIKey, nil); err != nil {
			openAIStatus.Error = err.Error()
		} else {
			openAIStatus.Valid = true
		}
	}

	var anthropicRequired, openAIRequired bool
	for _, modelKey := range []string{
		workspace_configuration_service.KeyIndexingModel,
		workspace_configuration_service.KeyNewsletterModel,
		workspace_configuration_service.KeyEmbeddingsModel,
	} {
		modelID, err := cfgSvc.Get(ctx, workspaceUUID, modelKey)
		if err != nil {
			return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		if info, ok := ai_client.ModelByID(modelID); ok {
			if info.Provider == ai_client.ProviderAnthropic {
				anthropicRequired = true
			}
			if info.Provider == ai_client.ProviderOpenAI {
				openAIRequired = true
			}
		}
	}

	return responses.Response(c, http.StatusOK, responses.WorkspaceHealthResponse{
		Anthropic:         anthropicStatus,
		OpenAI:            openAIStatus,
		AnthropicRequired: anthropicRequired,
		OpenAIRequired:    openAIRequired,
	})
}
