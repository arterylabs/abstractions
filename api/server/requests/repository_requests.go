package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type ConnectRepositoryRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
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
	SendDay                int       `json:"sendDay"`
	SendHour               int       `json:"sendHour"`
	NewsletterEnabled      *bool     `json:"newsletterEnabled"`
	Emails                 *[]string `json:"emails"`
	SendToWorkspaceMembers *bool     `json:"sendToWorkspaceMembers"`
}

func (r UpdateRepositorySettingsRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.SendDay, validation.Required, validation.Min(0), validation.Max(7)),
		validation.Field(&r.SendHour, validation.Min(0), validation.Max(23)),
		validation.Field(&r.Emails, validation.By(validateEmails)),
	)
}

// validateEmails runs is.EmailFormat over each address. Emails is a
// pointer field (nil means "leave recipients unchanged"), and
// validation.Each requires an already-dereferenced slice.
func validateEmails(value interface{}) error {
	emails, _ := value.(*[]string)
	if emails == nil {
		return nil
	}
	return validation.Validate(*emails, validation.Each(is.EmailFormat))
}

// UpdateRepositoryDescriptionRequest is the request body for PATCH /workspaces/:workspace_id/repositories/:repo_id.
type UpdateRepositoryDescriptionRequest struct {
	Description string `json:"description"`
}

func (r UpdateRepositoryDescriptionRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Description, validation.Length(0, 1000)),
	)
}
