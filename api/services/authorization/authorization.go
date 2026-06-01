package authorization_service

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrUnauthorized = errors.New("user does not have required permissions")
	ErrNoRole       = errors.New("user has no role in workspace")
)

type Service struct {
	db            *gorm.DB
	workspaceRepo *repositories.WorkspaceRepository
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db:            db,
		workspaceRepo: repositories.NewWorkspaceRepository(db),
	}
}

// Validate checks if a user has the required permissions in a workspace
func (s *Service) Validate(ctx context.Context, user *models.User, workspaceID uuid.UUID, requiredPermissions []Permission) error {
	// Super admin (god mode) bypass all permission checks
	if user.GodMode {
		return nil
	}

	// Get user's role in the workspace
	role, err := s.workspaceRepo.GetUserRole(ctx, user.ID, workspaceID)
	if err != nil {
		return fmt.Errorf("authorization_service: validate: get_user_role: %w", err)
	}
	if role == "" {
		return ErrNoRole
	}

	roleDetails := GetRoleByID(role)
	if roleDetails == nil {
		return ErrUnauthorized
	}

	// Check if user has all required permissions
	permissionSet := make(map[Permission]bool)
	for _, p := range roleDetails.Permissions {
		permissionSet[p] = true
	}

	for _, requiredPerm := range requiredPermissions {
		if !permissionSet[requiredPerm] {
			return ErrUnauthorized
		}
	}

	return nil
}

// HasPermission checks if a user has a specific permission in a workspace
func (s *Service) HasPermission(ctx context.Context, user *models.User, workspaceID uuid.UUID, permission string) bool {
	if user.GodMode {
		return true
	}

	role, err := s.workspaceRepo.GetUserRole(ctx, user.ID, workspaceID)
	if err != nil || role == "" {
		return false
	}

	roleDetails := GetRoleByID(role)
	if roleDetails == nil {
		return false
	}

	for _, p := range roleDetails.Permissions {
		if p.String() == permission {
			return true
		}
	}

	return false
}

func (s *Service) GetAllRoles(ctx context.Context) []Role {
	_ = ctx
	return ALL_ROLES
}
