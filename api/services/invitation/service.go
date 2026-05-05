package invitation_service

import (
	"errors"

	"gorm.io/gorm"
)

var (
	// ErrInvalidInvitation indicates that the invitation is not valid for the requested operation
	ErrInvalidInvitation = errors.New("invalid invitation")
	// ErrInvitationExpired indicates that the invitation has expired
	ErrInvitationExpired = errors.New("invitation has expired")
)

type InvitationInput struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required"`
}

type Service struct {
	DB *gorm.DB
}

func NewInvitationService(db *gorm.DB) *Service {
	return &Service{
		DB: db,
	}
}
