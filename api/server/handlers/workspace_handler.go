package handlers

import (
	"abstractions/server"
	"abstractions/server/requests"
	"abstractions/server/responses"
	workspace_service "abstractions/services/workspace"
	"abstractions/utils/context_utils"
	"abstractions/utils/pagination_utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
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
