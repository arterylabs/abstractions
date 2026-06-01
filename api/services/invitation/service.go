package invitation_service

import (
	"errors"

	"abstractions/repositories"
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
	db             *gorm.DB
	invitationRepo *repositories.InvitationRepository
	workspaceRepo  *repositories.WorkspaceRepository
	userRepo       *repositories.UserRepository
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db:             db,
		invitationRepo: repositories.NewInvitationRepository(db),
		workspaceRepo:  repositories.NewWorkspaceRepository(db),
		userRepo:       repositories.NewUserRepository(db),
	}
}
