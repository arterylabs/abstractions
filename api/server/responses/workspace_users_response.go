package responses

import (
	"abstractions/db/models"
	"time"
)

type PaginatedWorkspaceUsersResponse struct {
	Results  []WorkspaceUser `json:"results"`
	PageInfo PageInfo        `json:"pageInfo"`
}

type WorkspaceUser struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func NewPaginatedWorkspaceUsersResponse(users []models.UserWithRole, page, perPage, totalRecords int) *PaginatedWorkspaceUsersResponse {
	userResponses := make([]WorkspaceUser, len(users))
	for i, user := range users {
		userResponses[i] = WorkspaceUser{
			ID:        user.ID.String(),
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		}
	}

	return &PaginatedWorkspaceUsersResponse{
		Results:  userResponses,
		PageInfo: GetPageInfo(page, perPage, len(users), totalRecords),
	}
}
