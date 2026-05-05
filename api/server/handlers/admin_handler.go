package handlers

import (
	"abstractions/server"
	"abstractions/server/requests"
	"abstractions/server/responses"
	session_service "abstractions/services/session"
	subscription_service "abstractions/services/subscription"
	user_service "abstractions/services/user"
	workspace_service "abstractions/services/workspace"
	"abstractions/utils/context_utils"
	"abstractions/utils/ip_utils"
	"abstractions/utils/pagination_utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type AdminHandler struct {
	*BaseHandler
}

func NewAdminHandler(server *server.Server) *AdminHandler {
	return &AdminHandler{
		BaseHandler: NewBaseHandler(server),
	}
}

func (h *AdminHandler) GetUsers(c echo.Context) error {
	req := new(requests.AdminGetUsers)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	users, totalCount, err := h.Services(c).User.PaginateWithCountForAdmin(req.Page, req.PerPage, req.Query)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	response := responses.NewPaginatedAdminUsersResponse(users, req.Page, req.PerPage, totalCount)
	return responses.Response(c, http.StatusOK, response)
}

func (h *AdminHandler) GetUser(c echo.Context) error {
	userID := c.Param("id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	adminUserData, err := h.Services(c).User.GetAdminUserData(userUUID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	if adminUserData == nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "user not found")
	}

	return responses.Response(c, http.StatusOK, responses.NewAdminUserResponse(adminUserData))
}

func (h *AdminHandler) GetStats(c echo.Context) error {
	usersCount, err := h.Services(c).User.Count()
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	activeUsersCount, err := h.Services(c).User.CountActiveUsers()
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	subscriptionsCount, err := h.Services(c).Subscription.Count()
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	activeSubscriptionsCount, err := h.Services(c).Subscription.CountActiveSubscriptions()
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	workspacesCount, err := h.Services(c).Workspace.Count(c.Request().Context())
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, responses.AdminStatsResponse{
		UsersCount:               usersCount,
		ActiveUsersCount:         activeUsersCount,
		SubscriptionsCount:       subscriptionsCount,
		ActiveSubscriptionsCount: activeSubscriptionsCount,
		WorkspacesCount:          workspacesCount,
	})
}

func (h *AdminHandler) GetHistoricalData(c echo.Context) error {
	req := new(requests.AdminGetHistoricalDataRequest)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	response, err := h.Services(c).Analytics.GetAnalyticsData(req.AnalyticsParams)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, response)
}

func (h *AdminHandler) UpdateUser(c echo.Context) error {
	req := new(requests.AdminUpdateUser)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := h.Services(c).User.Update(req.ID, user_service.UpdateUserParams{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, user)
}

func (h *AdminHandler) CreateUser(c echo.Context) error {
	req := new(requests.AdminCreateUser)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := h.Services(c).User.CreateVerifiedUser(user_service.CreateVerifiedUserParams{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, user)
}

func (h *AdminHandler) CreateSubscription(c echo.Context) error {
	req := new(requests.AdminCreateSubscription)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	subs, err := h.Services(c).Subscription.CreateBulkSubscription(req.WorkspaceIDs, subscription_service.CreateSubscriptionParams{
		PlanID:    req.PlanID,
		Status:    req.Status,
		ExpiresAt: req.ExpiresAt,
		Limit:     req.Limit,
	})
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, responses.NewAdminSubscriptionsResponse(subs))
}

func (h *AdminHandler) UpdateSubscription(c echo.Context) error {
	req := new(requests.AdminUpdateSubscriptions)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	db := context_utils.GetTransactionFromContext(c)
	subscriptionService := subscription_service.NewSubscriptionService(db)
	err := subscriptionService.UpdateSubscriptionsForWorkspaces(req.WorkspaceIDs, subscription_service.UpdateSubscriptionParams{
		PlanID:    req.PlanID,
		Status:    req.Status,
		ExpiresAt: req.ExpiresAt,
		Limit:     req.Limit,
	})
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, "Subscriptions updated successfully")
}

func (h *AdminHandler) DisableSubscriptions(c echo.Context) error {
	req := new(requests.AdminDisableSubscriptions)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	db := context_utils.GetTransactionFromContext(c)
	subscriptionService := subscription_service.NewSubscriptionService(db)
	err := subscriptionService.DisableSubscriptions(req.WorkspaceIDs, req.DisabledReason)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, "Subscriptions disabled successfully")
}

func (h *AdminHandler) VerifyUsers(c echo.Context) error {
	req := new(requests.AdminVerifyUsers)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	db := context_utils.GetTransactionFromContext(c)
	userService := user_service.NewUserService(db)
	err := userService.BulkVerifyUsers(req.IDs)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, responses.AdminVerifyUsersResponse{VerifiedIDs: req.IDs})
}

func (h *AdminHandler) DeleteUsers(c echo.Context) error {
	req := new(requests.AdminDeleteUsers)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	db := context_utils.GetTransactionFromContext(c)
	userService := user_service.NewUserService(db)
	err = userService.BulkDeleteUsers(req.IDs, user.ID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, responses.AdminDeleteUsersResponse{DeletedIDs: req.IDs})
}

func (h *AdminHandler) RecoverUsers(c echo.Context) error {
	req := new(requests.AdminRecoverUsers)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	db := context_utils.GetTransactionFromContext(c)
	userService := user_service.NewUserService(db)
	err := userService.BulkRecoverUsers(req.IDs)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, responses.AdminRecoverUsersResponse{RecoveredIDs: req.IDs})
}

func (h *AdminHandler) ImpersonateUser(c echo.Context) error {
	req := new(requests.AdminImpersonateUser)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Get the admin user who is performing the impersonation
	adminUser, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	db := context_utils.GetTransactionFromContext(c)
	userService := user_service.NewUserService(db)

	// Create session with impersonation tracking
	sessionParams := &user_service.CreateSessionParams{
		UserAgent:      c.Request().UserAgent(),
		IP:             ip_utils.GetUserIPEcho(c),
		ImpersonatedBy: &adminUser.ID,
	}

	session, err := userService.CreateUserSession(req.ID, sessionParams)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, responses.AdminImpersonateUserResponse{Token: session.Token, ExpiresAt: session.ExpiresAt})
}

func (h *AdminHandler) GetUserWorkspacesWithRole(c echo.Context) error {
	userID := c.Param("user_id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	db := context_utils.GetTransactionFromContext(c)
	workspaceService := workspace_service.NewWorkspaceService(db)
	ctx := c.Request().Context()

	workspaces, err := workspaceService.GetUserWorkspacesWithRole(ctx, userUUID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, responses.NewWorkspaceWithRoleResponses(workspaces))
}

func (h *AdminHandler) GetWorkspaces(c echo.Context) error {
	req := new(requests.AdminGetWorkspaces)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	db := context_utils.GetTransactionFromContext(c)
	workspaceService := workspace_service.NewWorkspaceService(db)
	ctx := c.Request().Context()
	workspaces, totalCount, err := workspaceService.PaginateWithCountForAdmin(ctx, req.Page, req.PerPage, req.Query)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	response := responses.NewPaginatedAdminWorkspacesResponse(workspaces, req.Page, req.PerPage, totalCount)
	return responses.Response(c, http.StatusOK, response)
}

func (h *AdminHandler) CreateWorkspace(c echo.Context) error {
	req := new(requests.AdminCreateWorkspace)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	db := context_utils.GetTransactionFromContext(c)
	workspaceService := workspace_service.NewWorkspaceService(db)
	ctx := c.Request().Context()
	workspace, err := workspaceService.CreateWorkspace(ctx, req.Name, user.ID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusCreated, responses.NewWorkspaceResponse(workspace))
}

func (h *AdminHandler) GetWorkspace(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	db := context_utils.GetTransactionFromContext(c)
	workspaceService := workspace_service.NewWorkspaceService(db)
	ctx := c.Request().Context()
	adminWorkspaceData, err := workspaceService.GetAdminWorkspaceData(ctx, workspaceUUID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	if adminWorkspaceData == nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "workspace not found")
	}

	return responses.Response(c, http.StatusOK, responses.NewAdminWorkspaceResponse(adminWorkspaceData))
}

func (h *AdminHandler) UpdateWorkspace(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	req := new(requests.AdminUpdateWorkspace)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	db := context_utils.GetTransactionFromContext(c)
	workspaceService := workspace_service.NewWorkspaceService(db)
	ctx := c.Request().Context()
	params := workspace_service.UpdateWorkspaceParams{
		Name:    req.Name,
		Logo:    req.Logo,
		Website: req.Website,
		About:   req.About,
	}
	if err := workspaceService.UpdateWorkspace(ctx, workspaceUUID, params, user.ID); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, nil)
}

func (h *AdminHandler) GetWorkspaceUsers(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	req := new(requests.AdminGetWorkspaceUsers)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	db := context_utils.GetTransactionFromContext(c)
	workspaceService := workspace_service.NewWorkspaceService(db)
	ctx := c.Request().Context()

	paginationParams := pagination_utils.PaginationParams{
		Page:      req.Page,
		PerPage:   req.PerPage,
		Query:     req.Query,
		Sort:      req.Sort,
		SortOrder: req.SortOrder,
	}

	users, totalCount, err := workspaceService.GetWorkspaceUsersPaginated(ctx, workspaceUUID, paginationParams)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	response := responses.NewPaginatedWorkspaceUsersResponse(users, req.Page, req.PerPage, totalCount)
	return responses.Response(c, http.StatusOK, response)
}

func (h *AdminHandler) AddUserToWorkspace(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	req := new(requests.AdminAddUserToWorkspace)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	db := context_utils.GetTransactionFromContext(c)
	workspaceService := workspace_service.NewWorkspaceService(db)
	ctx := c.Request().Context()
	if err := workspaceService.AddUserToWorkspace(ctx, workspaceUUID, req.UserID, req.Role, user.ID); err != nil {
		switch err {
		case workspace_service.ErrInsufficientPermissions:
			return responses.ErrorResponse(c, http.StatusForbidden, "You don't have permission to assign this role")
		default:
			return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
	}

	return responses.Response(c, http.StatusCreated, nil)
}

func (h *AdminHandler) UpdateWorkspaceUserRole(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	userID := c.Param("userId")

	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	req := new(requests.AdminUpdateWorkspaceUserRole)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	db := context_utils.GetTransactionFromContext(c)
	workspaceService := workspace_service.NewWorkspaceService(db)
	ctx := c.Request().Context()
	if err := workspaceService.UpdateUserRole(ctx, workspaceUUID, userUUID, req.Role, user.ID); err != nil {
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

func (h *AdminHandler) RemoveUserFromWorkspace(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	userID := c.Param("userId")

	workspaceUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	db := context_utils.GetTransactionFromContext(c)
	workspaceService := workspace_service.NewWorkspaceService(db)
	ctx := c.Request().Context()
	if err := workspaceService.RemoveUserFromWorkspace(ctx, workspaceUUID, userUUID, user.ID); err != nil {
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

func (h *AdminHandler) BulkDeleteWorkspaces(c echo.Context) error {
	req := new(requests.AdminDeleteWorkspaces)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	db := context_utils.GetTransactionFromContext(c)
	workspaceService := workspace_service.NewWorkspaceService(db)
	ctx := c.Request().Context()

	err = workspaceService.BulkDeleteWorkspaces(ctx, req.IDs, user.ID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, responses.AdminDeleteWorkspacesResponse{DeletedIDs: req.IDs})
}

func (h *AdminHandler) BulkRecoverWorkspaces(c echo.Context) error {
	req := new(requests.AdminRecoverWorkspaces)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	if err := req.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	db := context_utils.GetTransactionFromContext(c)
	workspaceService := workspace_service.NewWorkspaceService(db)
	ctx := c.Request().Context()

	err = workspaceService.BulkRecoverWorkspaces(ctx, req.IDs, user.ID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, responses.AdminRecoverWorkspacesResponse{RecoveredIDs: req.IDs})
}

func (h *AdminHandler) ExitImpersonation(c echo.Context) error {

	// get user from context and delete all its impersonation sessions
	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	db := context_utils.GetTransactionFromContext(c)
	sessionService := session_service.NewSessionService(db)

	err = sessionService.InvalidateAllImpersonationSessionsForAdmin(user.ID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	// return success response
	return responses.MessageResponse(c, http.StatusOK, "Exited impersonation successfully")
}
