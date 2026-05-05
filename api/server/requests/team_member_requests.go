package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// AddTeamMemberRequest is the request body for POST /repositories/:repo_id/members.
type AddTeamMemberRequest struct {
	Email string `json:"email"`
}

func (r AddTeamMemberRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.Email),
	)
}
