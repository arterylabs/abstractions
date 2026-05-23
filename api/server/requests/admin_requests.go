package requests

import (
	"abstractions/db/models"
	"abstractions/utils/analytics_utils"
	"abstractions/utils/pagination_utils"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type AdminGetUsers struct {
	pagination_utils.PaginationParams
	IncludeDeleted bool `query:"includeDeleted"`
}

func (ba AdminGetUsers) Validate() error {
	ba.PaginationParams.Validate()
	return validation.ValidateStruct(&ba)
}

type AdminUpdateUser struct {
	ID    uuid.UUID `json:"id"`
	Email *string   `json:"email"`
	Name  *string   `json:"name"`
}

type AdminUpdateUsers []AdminUpdateUser

func (ba AdminUpdateUsers) Validate() error {
	for _, user := range ba {
		if err := user.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type AdminCreateUser struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (ba AdminCreateUser) Validate() error {
	return validation.ValidateStruct(&ba,
		validation.Field(&ba.Email, validation.Required),
		validation.Field(&ba.Name, validation.Required),
	)
}

func (ba AdminUpdateUser) Validate() error {
	return validation.ValidateStruct(&ba,
		validation.Field(&ba.ID, validation.Required),
	)
}

type AdminVerifyUsers struct {
	IDs []uuid.UUID `json:"userIds"`
}

func (ba AdminVerifyUsers) Validate() error {
	return validation.ValidateStruct(&ba,
		validation.Field(&ba.IDs, validation.Required),
	)
}

type AdminDeleteUsers struct {
	IDs []uuid.UUID `json:"userIds"`
}

func (ba AdminDeleteUsers) Validate() error {
	return validation.ValidateStruct(&ba,
		validation.Field(&ba.IDs, validation.Required),
	)
}

type AdminRecoverUsers struct {
	IDs []uuid.UUID `json:"userIds"`
}

func (ba AdminRecoverUsers) Validate() error {
	return validation.ValidateStruct(&ba,
		validation.Field(&ba.IDs, validation.Required),
	)
}

type AdminDeleteWorkspaces struct {
	IDs []uuid.UUID `json:"workspaceIds"`
}

func (ba AdminDeleteWorkspaces) Validate() error {
	return validation.ValidateStruct(&ba,
		validation.Field(&ba.IDs, validation.Required),
	)
}

type AdminRecoverWorkspaces struct {
	IDs []uuid.UUID `json:"workspaceIds"`
}

func (ba AdminRecoverWorkspaces) Validate() error {
	return validation.ValidateStruct(&ba,
		validation.Field(&ba.IDs, validation.Required),
	)
}

type AdminCreateSubscription struct {
	WorkspaceIDs []uuid.UUID               `json:"workspaceIds"`
	PlanID       models.PlanID             `json:"planId"`
	Status       models.SubscriptionStatus `json:"status"`
	ExpiresAt    time.Time                 `json:"expiresAt"`
	Limit        int                       `json:"limit"`
}

func (ba AdminCreateSubscription) Validate() error {
	return validation.ValidateStruct(&ba,
		validation.Field(&ba.WorkspaceIDs, validation.Required),
		validation.Field(&ba.PlanID, validation.Required),
		validation.Field(&ba.Status, validation.Required),
		validation.Field(&ba.ExpiresAt, validation.Required),
		validation.Field(&ba.Limit, validation.Required),
	)
}

type AdminUpdateSubscriptions struct {
	WorkspaceIDs []uuid.UUID                `json:"workspaceIds"`
	PlanID       *models.PlanID             `json:"planId"`
	Status       *models.SubscriptionStatus `json:"status"`
	ExpiresAt    *time.Time                 `json:"expiresAt"`
	Limit        *int                       `json:"limit"`
}

func (ba AdminUpdateSubscriptions) Validate() error {
	return validation.ValidateStruct(&ba,
		validation.Field(&ba.WorkspaceIDs, validation.Required),
	)
}

type AdminDisableSubscriptions struct {
	WorkspaceIDs   []uuid.UUID `json:"workspaceIds"`
	DisabledReason string      `json:"disabledReason"`
}

func (ba AdminDisableSubscriptions) Validate() error {
	return validation.ValidateStruct(&ba,
		validation.Field(&ba.WorkspaceIDs, validation.Required),
		validation.Field(&ba.DisabledReason, validation.Required),
	)
}

type AdminImpersonateUser struct {
	ID uuid.UUID `json:"id"`
}

func (ba AdminImpersonateUser) Validate() error {
	return validation.ValidateStruct(&ba,
		validation.Field(&ba.ID, validation.Required),
	)
}

type AdminGetHistoricalDataRequest struct {
	analytics_utils.AnalyticsParams
}

func (r AdminGetHistoricalDataRequest) Validate() error {
	r.SetDefaults()
	return r.AnalyticsParams.Validate()
}
