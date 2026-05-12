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
	DB *gorm.DB
}

func NewAuthorizationService(db *gorm.DB) *Service {
	return &Service{
		DB: db,
	}
}

// Validate checks if a user has the required permissions in a workspace
func (s *Service) Validate(ctx context.Context, user *models.User, workspaceID uuid.UUID, requiredPermissions []Permission) error {
	// Super admin (god mode) bypass all permission checks
	if user.GodMode {
		return nil
	}

	workspaceRepo := repositories.NewWorkspaceRepository(s.DB)

	// Get user's role in the workspace
	role, err := workspaceRepo.GetUserRole(ctx, user.ID, workspaceID)
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

	workspaceRepo := repositories.NewWorkspaceRepository(s.DB)

	role, err := workspaceRepo.GetUserRole(ctx, user.ID, workspaceID)
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

func (s *Service) GetAllRoles() []Role {
	return ALL_ROLES
}
