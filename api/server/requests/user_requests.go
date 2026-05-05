package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

const (
	minPathLength = 8
)

type BasicAuth struct {
	Email    string `json:"email" validate:"required" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required" example:"11111111"`
}

func (ba BasicAuth) Validate() error {
	return validation.ValidateStruct(&ba,
		validation.Field(&ba.Email, is.Email),
		validation.Field(&ba.Password, validation.Length(minPathLength, 0)),
	)
}

type LoginRequest struct {
	BasicAuth
}

type SignupRequest struct {
	BasicAuth
	Name string `json:"name" validate:"required" example:"John Doe"`
}

func (rr SignupRequest) Validate() error {
	err := rr.BasicAuth.Validate()
	if err != nil {
		return err
	}

	return validation.ValidateStruct(&rr,
		validation.Field(&rr.Name, validation.Required),
	)
}

type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

func (ver VerifyEmailRequest) Validate() error {
	return validation.ValidateStruct(&ver,
		validation.Field(&ver.Token, validation.Required),
	)
}

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (rpr ResetPasswordRequest) Validate() error {
	return validation.ValidateStruct(&rpr,
		validation.Field(&rpr.Email, validation.Required, is.Email),
	)
}

type ResetPasswordConfirmRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required"`
	ID       string `json:"id" validate:"required"`
}

func (rpc ResetPasswordConfirmRequest) Validate() error {
	return validation.ValidateStruct(&rpc,
		validation.Field(&rpc.Token, validation.Required),
		validation.Field(&rpc.Password, validation.Required, validation.Length(minPathLength, 0)),
		validation.Field(&rpc.Token, validation.Required),
	)
}

type UpdateUserRequest struct {
	Name string `json:"name"`
}

func (uur UpdateUserRequest) Validate() error {
	return validation.ValidateStruct(&uur,
		validation.Field(&uur.Name, validation.Required),
	)
}

type UpdateUserPasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func (uur UpdateUserPasswordRequest) Validate() error {
	return validation.ValidateStruct(&uur,
		validation.Field(&uur.NewPassword, validation.Required),
	)
}

type SwitchWorkspaceRequest struct {
	WorkspaceID uuid.UUID `json:"workspace_id" validate:"required"`
}

func (r *SwitchWorkspaceRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.WorkspaceID, validation.Required, is.UUID),
	)
}
