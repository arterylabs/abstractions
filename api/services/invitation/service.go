package invitation_service

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrInvalidInvitation = errors.New("invitation is invalid")
	ErrInvitationExpired = errors.New("invitation has expired")
)

type InvitationInput struct {
	Email string
	Role  string
}

type Service struct {
	DB *gorm.DB
}

func NewInvitationService(db *gorm.DB) *Service {
	return &Service{DB: db}
}
