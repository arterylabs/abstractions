package requests

import (
	authorization_service "abstractions/services/authorization"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type SendInvitationsRequest struct {
	Invitations []InvitationInput `json:"invitations" validate:"required,dive"`
}

type InvitationInput struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (r SendInvitationsRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Invitations, validation.Required, validation.Length(1, 0)),
	)
}

func (i InvitationInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Email, validation.Required, emailRule),
		validation.Field(&i.Role, validation.Required, validation.In(authorization_service.AllRoleIDsInterface()...)), // Add valid roles here
	)
}
