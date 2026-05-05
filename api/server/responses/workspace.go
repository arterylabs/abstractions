package responses

import (
	"abstractions/db/models"
	"time"

	"github.com/google/uuid"
)

type Workspace struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Logo         *string   `json:"logo,omitempty"`
	Website      *string   `json:"website,omitempty"`
	About        *string   `json:"about,omitempty"`
	BillingEmail string    `json:"billingEmail"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func NewWorkspaceResponse(workspace *models.Workspace) *Workspace {
	if workspace == nil {
		return nil
	}

	return &Workspace{
		ID:           workspace.ID,
		Name:         workspace.Name,
		Logo:         workspace.Logo,
		Website:      workspace.Website,
		About:        workspace.About,
		BillingEmail: workspace.BillingEmail,
		CreatedAt:    workspace.CreatedAt,
		UpdatedAt:    workspace.UpdatedAt,
	}
}
