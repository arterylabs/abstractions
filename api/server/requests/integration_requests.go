package requests

import (
	"abstractions/db/models"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type IntegrationSetupRequest struct {
	RedirectTo   *string `json:"redirectTo"`
	ResourceType *string `json:"resourceType"`
	ResourceID   *string `json:"resourceId"`
}

func (r IntegrationSetupRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.RedirectTo, validation.Length(1, 500)),
	)
}

type UpdateIntegrationRequest struct {
	ServiceTitle  *string           `json:"service_title"`
	Configuration models.JSONBField `json:"configuration"`
}

func (r UpdateIntegrationRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ServiceTitle, validation.Length(1, 255)),
	)
}
