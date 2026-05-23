package responses

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"time"

	"github.com/google/uuid"
)

type AdminUserResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
}

type AdminSubscriptionsResponse struct {
	ID          uuid.UUID                 `json:"id"`
	WorkspaceID uuid.UUID                 `json:"workspaceId"`
	PlanID      models.PlanID             `json:"planId"`
	Status      models.SubscriptionStatus `json:"status"`
	ExpiresAt   time.Time                 `json:"expiresAt"`
	Limit       int                       `json:"limit"`
}

func NewAdminSubscriptionsResponse(subscriptions *[]models.Subscription) []AdminSubscriptionsResponse {
	response := []AdminSubscriptionsResponse{}
	for _, subscription := range *subscriptions {
		response = append(response, AdminSubscriptionsResponse{
			ID:          subscription.ID,
			WorkspaceID: subscription.WorkspaceID,
			PlanID:      subscription.PlanID,
			Status:      subscription.Status,
			ExpiresAt:   subscription.ExpiresAt,
			Limit:       subscription.Limit,
		})
	}
	return response
}

type AdminVerifyUsersResponse struct {
	VerifiedIDs []uuid.UUID `json:"verifiedIds"`
}

type AdminDeleteUsersResponse struct {
	DeletedIDs []uuid.UUID `json:"deletedIds"`
}

type AdminRecoverUsersResponse struct {
	RecoveredIDs []uuid.UUID `json:"recoveredIds"`
}

type AdminDeleteWorkspacesResponse struct {
	DeletedIDs []uuid.UUID `json:"deletedIds"`
}

type AdminRecoverWorkspacesResponse struct {
	RecoveredIDs []uuid.UUID `json:"recoveredIds"`
}

type AdminDisableSubscriptionsResponse struct {
	DisabledIDs []uuid.UUID `json:"disabledIds"`
}

type AdminImpersonateUserResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type AdminStatsResponse struct {
	UsersCount        int `json:"usersCount"`
	ActiveUsersCount  int `json:"activeUsersCount"`
	WorkspacesCount   int `json:"workspacesCount"`
	RepositoriesCount int `json:"repositoriesCount"`
	NewslettersCount  int `json:"newslettersCount"`
}

// Standard response structure for analytics chart data
type AnalyticsDataPoint struct {
	Label string      `json:"label"` // The formatted date/time label
	Value interface{} `json:"value"` // The actual timestamp for sorting
	Count int         `json:"count"` // The count for this data point
}

type AnalyticsDataResponse struct {
	DataPoints  []AnalyticsDataPoint `json:"dataPoints"`
	Metric      string               `json:"metric"`      // e.g., "users_created", "workspaces_created"
	Dimension   string               `json:"dimension"`   // e.g., "date", "hour", "week"
	Granularity string               `json:"granularity"` // e.g., "day", "hour", "week"
	StartDate   string               `json:"startDate"`
	EndDate     string               `json:"endDate"`
	TotalCount  int                  `json:"totalCount"` // Sum of all counts
}

type AdminUser struct {
	ID                   string  `json:"id"`
	Name                 string  `json:"name"`
	Email                string  `json:"email"`
	Avatar               *string `json:"avatar"`
	DefaultWorkspaceID   *string `json:"defaultWorkspaceId"`
	DefaultWorkspaceName *string `json:"defaultWorkspaceName"`
	WorkspacesCount      int     `json:"workspacesCount"`
	CreatedAt            string  `json:"createdAt"`
	UpdatedAt            string  `json:"updatedAt"`
	VerifiedAt           *string `json:"verifiedAt"`
	DeletedAt            *string `json:"deletedAt"`
	DeletedBy            *string `json:"deletedBy"`
	DeletedByName        *string `json:"deletedByName"`
	GodMode              bool    `json:"godMode"`
}

type PaginatedAdminUsersResponse struct {
	Results  []AdminUser `json:"results"`
	PageInfo PageInfo    `json:"pageInfo"`
}

type AdminWorkspace struct {
	ID            string             `json:"id"`
	Name          string             `json:"name"`
	Logo          *string            `json:"logo"`
	Website       *string            `json:"website"`
	About         *string            `json:"about"`
	BillingEmail  string             `json:"billingEmail"`
	CreatedAt     string             `json:"createdAt"`
	UpdatedAt     string             `json:"updatedAt"`
	DeletedAt     *string            `json:"deletedAt"`
	DeletedBy     *string            `json:"deletedBy"`
	DeletedByName *string            `json:"deletedByName"`
	CreatedByName *string            `json:"createdByName"`
	UpdatedByName *string            `json:"updatedByName"`
	UsersCount    int                `json:"usersCount"`
	Subscription  *AdminSubscription `json:"subscription"`
}

type AdminSubscription struct {
	ID        string                    `json:"id"`
	PlanID    models.PlanID             `json:"planId"`
	Status    models.SubscriptionStatus `json:"status"`
	ExpiresAt string                    `json:"expiresAt"`
	Limit     int                       `json:"limit"`
}

type PaginatedAdminWorkspacesResponse struct {
	Results  []AdminWorkspace `json:"results"`
	PageInfo PageInfo         `json:"pageInfo"`
}

func NewPaginatedAdminUsersResponse(users []repositories.AdminUserData, page, perPage, totalRecords int) *PaginatedAdminUsersResponse {
	userResponses := make([]AdminUser, len(users))
	for i, userData := range users {
		user := userData.User
		userResponse := AdminUser{
			ID:              user.ID.String(),
			Name:            user.Name,
			Email:           user.Email,
			Avatar:          user.Avatar,
			WorkspacesCount: userData.WorkspacesCount,
			CreatedAt:       user.CreatedAt.Format(time.RFC3339),
			UpdatedAt:       user.UpdatedAt.Format(time.RFC3339),
			GodMode:         user.GodMode,
		}

		// Handle default workspace ID and name
		if user.DefaultWorkspaceID != nil {
			workspaceIdStr := user.DefaultWorkspaceID.String()
			userResponse.DefaultWorkspaceID = &workspaceIdStr
			userResponse.DefaultWorkspaceName = userData.DefaultWorkspaceName
		}

		// Handle verified status
		if user.VerifiedAt != nil {
			verifiedAtStr := user.VerifiedAt.Format(time.RFC3339)
			userResponse.VerifiedAt = &verifiedAtStr
		}

		// Handle deletion info
		if user.DeletedAt != nil {
			deletedAtStr := user.DeletedAt.Format(time.RFC3339)
			userResponse.DeletedAt = &deletedAtStr

			if user.DeletedBy != nil {
				deletedByStr := user.DeletedBy.String()
				userResponse.DeletedBy = &deletedByStr
				userResponse.DeletedByName = userData.DeletedByName
			}
		}

		userResponses[i] = userResponse
	}

	return &PaginatedAdminUsersResponse{
		Results:  userResponses,
		PageInfo: GetPageInfo(page, perPage, len(users), totalRecords),
	}
}

func NewPaginatedAdminWorkspacesResponse(workspaces []repositories.AdminWorkspaceData, page, perPage, totalRecords int) *PaginatedAdminWorkspacesResponse {
	workspaceResponses := make([]AdminWorkspace, len(workspaces))
	for i, workspaceData := range workspaces {
		workspace := workspaceData.Workspace

		workspaceResponse := AdminWorkspace{
			ID:            workspace.ID.String(),
			Name:          workspace.Name,
			Logo:          workspace.Logo,
			Website:       workspace.Website,
			About:         workspace.About,
			BillingEmail:  workspace.BillingEmail,
			CreatedAt:     workspace.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     workspace.UpdatedAt.Format(time.RFC3339),
			CreatedByName: workspaceData.CreatedByName,
			UpdatedByName: workspaceData.UpdatedByName,
			UsersCount:    workspaceData.UsersCount,
		}

		// Handle deleted status
		if workspace.DeletedAt != nil {
			deletedAtStr := workspace.DeletedAt.Format(time.RFC3339)
			workspaceResponse.DeletedAt = &deletedAtStr
			if workspace.DeletedBy != nil {
				deletedByStr := workspace.DeletedBy.String()
				workspaceResponse.DeletedBy = &deletedByStr
			}
			if workspaceData.DeletedByName != nil {
				workspaceResponse.DeletedByName = workspaceData.DeletedByName
			}
		}

		// Handle subscription data
		if workspaceData.Subscription != nil && workspaceData.Subscription.ID != uuid.Nil {
			workspaceResponse.Subscription = &AdminSubscription{
				ID:        workspaceData.Subscription.ID.String(),
				PlanID:    workspaceData.Subscription.PlanID,
				Status:    workspaceData.Subscription.Status,
				ExpiresAt: workspaceData.Subscription.ExpiresAt.Format(time.RFC3339),
				Limit:     workspaceData.Subscription.Limit,
			}
		}

		workspaceResponses[i] = workspaceResponse
	}

	return &PaginatedAdminWorkspacesResponse{
		Results:  workspaceResponses,
		PageInfo: GetPageInfo(page, perPage, len(workspaces), totalRecords),
	}
}

func NewAdminWorkspaceResponse(workspaceData *repositories.AdminWorkspaceData) *AdminWorkspace {
	if workspaceData == nil {
		return nil
	}

	workspace := workspaceData.Workspace

	workspaceResponse := AdminWorkspace{
		ID:            workspace.ID.String(),
		Name:          workspace.Name,
		Logo:          workspace.Logo,
		Website:       workspace.Website,
		About:         workspace.About,
		BillingEmail:  workspace.BillingEmail,
		CreatedAt:     workspace.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     workspace.UpdatedAt.Format(time.RFC3339),
		CreatedByName: workspaceData.CreatedByName,
		UpdatedByName: workspaceData.UpdatedByName,
		UsersCount:    workspaceData.UsersCount,
	}

	// Handle deleted status
	if workspace.DeletedAt != nil {
		deletedAtStr := workspace.DeletedAt.Format(time.RFC3339)
		workspaceResponse.DeletedAt = &deletedAtStr

		if workspace.DeletedBy != nil {
			deletedByStr := workspace.DeletedBy.String()
			workspaceResponse.DeletedBy = &deletedByStr
		}
		if workspaceData.DeletedByName != nil {
			workspaceResponse.DeletedByName = workspaceData.DeletedByName
		}
	}

	// Handle subscription data
	if workspaceData.Subscription != nil && workspaceData.Subscription.ID != uuid.Nil {
		workspaceResponse.Subscription = &AdminSubscription{
			ID:        workspaceData.Subscription.ID.String(),
			PlanID:    workspaceData.Subscription.PlanID,
			Status:    workspaceData.Subscription.Status,
			ExpiresAt: workspaceData.Subscription.ExpiresAt.Format(time.RFC3339),
			Limit:     workspaceData.Subscription.Limit,
		}
	}

	return &workspaceResponse
}

// NewAdminUserResponse converts AdminUserData to AdminUser response format
func NewAdminUserResponse(userData *repositories.AdminUserData) *AdminUser {
	if userData == nil {
		return nil
	}

	user := userData.User
	userResponse := AdminUser{
		ID:              user.ID.String(),
		Name:            user.Name,
		Email:           user.Email,
		Avatar:          user.Avatar,
		WorkspacesCount: userData.WorkspacesCount,
		CreatedAt:       user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       user.UpdatedAt.Format(time.RFC3339),
		GodMode:         user.GodMode,
	}

	// Handle default workspace ID and name
	if user.DefaultWorkspaceID != nil {
		workspaceIdStr := user.DefaultWorkspaceID.String()
		userResponse.DefaultWorkspaceID = &workspaceIdStr
		userResponse.DefaultWorkspaceName = userData.DefaultWorkspaceName
	}

	// Handle verified status
	if user.VerifiedAt != nil {
		verifiedAtStr := user.VerifiedAt.Format(time.RFC3339)
		userResponse.VerifiedAt = &verifiedAtStr
	}

	// Handle deletion info
	if user.DeletedAt != nil {
		deletedAtStr := user.DeletedAt.Format(time.RFC3339)
		userResponse.DeletedAt = &deletedAtStr

		if user.DeletedBy != nil {
			deletedByStr := user.DeletedBy.String()
			userResponse.DeletedBy = &deletedByStr
			userResponse.DeletedByName = userData.DeletedByName
		}
	}

	return &userResponse
}
