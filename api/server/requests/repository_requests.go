package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ConnectRepositoryRequest struct {
	Name          string `json:"name"`
	ServiceType   string `json:"serviceType"`
	ServiceID     string `json:"serviceId"`
	DefaultBranch string `json:"defaultBranch"`
}

func (r ConnectRepositoryRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.ServiceID, validation.Required),
		validation.Field(&r.DefaultBranch, validation.Length(0, 255)),
	)
}

// UpdateRepositorySettingsRequest is the request body for PUT /workspaces/:workspace_id/repositories/:repo_id.
type UpdateRepositorySettingsRequest struct {
	SendDay           int   `json:"sendDay"`
	SendHour          int   `json:"sendHour"`
	NewsletterEnabled *bool `json:"newsletterEnabled"`
}

func (r UpdateRepositorySettingsRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.SendDay, validation.Required, validation.Min(0), validation.Max(7)),
		validation.Field(&r.SendHour, validation.Min(0), validation.Max(23)),
	)
}
