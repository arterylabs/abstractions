package requests

import (
	"abstractions/db/models"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type IntegrationSetupRequest struct {
	RedirectTo   *string `json:"redirectTo"`
	ResourceType *string `json:"resourceType"`
	ResourceID   *string `json:"resourceId"`
	// CreateRepository signals that this OAuth flow is for connecting a new
	// repository. No repository UUID exists yet, so ResourceID must not be used.
	CreateRepository *bool `json:"createRepository"`
}

func (r IntegrationSetupRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.RedirectTo, validation.Length(1, 500)),
		validation.Field(&r.ResourceType, validation.Length(1, 100)),
		validation.Field(&r.ResourceID, is.UUID),
	)
}

type UpdateIntegrationRequest struct {
	ServiceTitle  *string           `json:"serviceTitle"`
	Configuration models.JSONBField `json:"configuration"`
}

func (r UpdateIntegrationRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ServiceTitle, validation.Length(1, 255)),
	)
}
