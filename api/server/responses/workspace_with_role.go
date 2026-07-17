package responses

import (
	"abstractions/db/models"
	"time"

	"github.com/google/uuid"
)

type WorkspaceWithRole struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Logo         *string   `json:"logo,omitempty"`
	Website      *string   `json:"website,omitempty"`
	About        *string   `json:"about,omitempty"`
	BillingEmail string    `json:"billingEmail"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	UserRole     string    `json:"userRole"`
}

func NewWorkspaceWithRoleResponse(workspace *models.WorkspaceWithRole) *WorkspaceWithRole {
	if workspace == nil {
		return nil
	}

	return &WorkspaceWithRole{
		ID:           workspace.ID,
		Name:         workspace.Name,
		Logo:         workspace.Logo,
		Website:      workspace.Website,
		About:        workspace.About,
		BillingEmail: workspace.BillingEmail,
		CreatedAt:    workspace.CreatedAt,
		UpdatedAt:    workspace.UpdatedAt,
		UserRole:     workspace.UserRole,
	}
}

func NewWorkspaceWithRoleResponses(workspaces []models.WorkspaceWithRole) []WorkspaceWithRole {
	responses := make([]WorkspaceWithRole, len(workspaces))
	for i, workspace := range workspaces {
		responses[i] = *NewWorkspaceWithRoleResponse(&workspace)
	}
	return responses
}
