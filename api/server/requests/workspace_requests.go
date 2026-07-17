package requests

import (
	"abstractions/utils/pagination_utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type AdminGetWorkspaces struct {
	pagination_utils.PaginationParams
	IncludeDeleted bool `query:"includeDeleted"`
}

func (ba AdminGetWorkspaces) Validate() error {
	ba.PaginationParams.Validate()
	return validation.ValidateStruct(&ba)
}

type AdminCreateWorkspace struct {
	Name    string `json:"name" validate:"required"`
	Logo    string `json:"logo"`
	Website string `json:"website"`
	About   string `json:"about"`
}

type AdminUpdateWorkspace struct {
	Name    string `json:"name" validate:"required"`
	Logo    string `json:"logo"`
	Website string `json:"website"`
	About   string `json:"about"`
}

type AdminGetWorkspaceUsers struct {
	pagination_utils.PaginationParams
}

func (ba AdminGetWorkspaceUsers) Validate() error {
	ba.PaginationParams.Validate()
	return validation.ValidateStruct(&ba)
}

type AdminAddUserToWorkspace struct {
	UserID uuid.UUID `json:"userId" validate:"required"`
	Role   string    `json:"role" validate:"required"`
}

type AdminUpdateWorkspaceUserRole struct {
	Role string `json:"role" validate:"required"`
}

// Regular user workspace request structs
type CreateWorkspace struct {
	Name    string `json:"name"`
	Logo    string `json:"logo"`
	Website string `json:"website"`
	About   string `json:"about"`
}

func (r CreateWorkspace) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
	)
}

type UpdateWorkspace struct {
	Name         string `json:"name"`
	Logo         string `json:"logo"`
	Website      string `json:"website"`
	About        string `json:"about"`
	BillingEmail string `json:"billingEmail"`
}

func (r UpdateWorkspace) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.BillingEmail, validation.Required, emailRule),
	)
}

type UpdateWorkspaceUserRole struct {
	UserID uuid.UUID `json:"userId"`
	Role   string    `json:"role"`
}

func (r UpdateWorkspaceUserRole) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UserID, validation.Required),
		validation.Field(&r.Role, validation.Required),
	)
}

type RemoveWorkspaceUser struct {
	UserID uuid.UUID `json:"userId"`
}

func (r RemoveWorkspaceUser) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UserID, validation.Required),
	)
}
