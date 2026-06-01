package workspace_service

import (
	"abstractions/clients/app_notifications"
	"abstractions/db/models"
	"abstractions/repositories"
	authorization_service "abstractions/services/authorization"
	"abstractions/utils/pagination_utils"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrWorkspaceNotFound       = errors.New("workspace not found")
	ErrInsufficientPermissions = errors.New("insufficient permissions to perform this action")
	ErrCannotModifyOwnRole     = errors.New("cannot modify your own role")
)

type Service struct {
	db            *gorm.DB
	workspaceRepo *repositories.WorkspaceRepository
	userRepo      *repositories.UserRepository
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db:            db,
		workspaceRepo: repositories.NewWorkspaceRepository(db),
		userRepo:      repositories.NewUserRepository(db),
	}
}

func (s *Service) Create(ctx context.Context, name string, ownerID uuid.UUID) (*models.Workspace, error) {
	// Create user repository to fetch user email for billing_email
	user, err := s.userRepo.GetByID(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("workspace_service: create_workspace: get_user: %w", err)
	}

	workspace := &models.Workspace{
		Name:         name,
		BillingEmail: user.Email, // Set billing email to creator's email
		CreatedBy:    ownerID,    // Set created by user ID
		UpdatedBy:    ownerID,    // Set updated by user ID
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	if err := s.workspaceRepo.Create(ctx, workspace); err != nil {
		return nil, fmt.Errorf("workspace_service: create_workspace: create: %w", err)
	}

	// Add user relationship
	relationship := &models.UserWorkspaceRelationship{
		UserID:      ownerID,
		WorkspaceID: workspace.ID,
		Role:        authorization_service.DEFAULT_ROLE.String(), // Default role for all users
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := s.workspaceRepo.AddUserToWorkspace(ctx, relationship); err != nil {
		return nil, fmt.Errorf("workspace_service: create_workspace: add_user: %w", err)
	}

	go app_notifications.SendWorkspaceCreatedNotification(context.Background(), app_notifications.WorkspaceCreatedParams{
		WorkspaceName: workspace.Name,
		OwnerEmail:    user.Email,
	})

	return workspace, nil
}

func (s *Service) GetWorkspace(ctx context.Context, id uuid.UUID) (*models.Workspace, error) {
	workspace, err := s.workspaceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("workspace_service: get_workspace: %w", err)
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}
	return workspace, nil
}

func (s *Service) UpdateWorkspace(ctx context.Context, id uuid.UUID, params UpdateWorkspaceParams, userID uuid.UUID) error {
	workspace, err := s.workspaceRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("workspace_service: update_workspace: get_workspace: %w", err)
	}
	if workspace == nil {
		return ErrWorkspaceNotFound
	}

	workspace.Name = params.Name
	if params.Logo != "" {
		workspace.Logo = &params.Logo
	}
	if params.Website != "" {
		workspace.Website = &params.Website
	}
	if params.About != "" {
		workspace.About = &params.About
	}
	if params.BillingEmail != "" {
		workspace.BillingEmail = params.BillingEmail
	}
	workspace.UpdatedBy = userID // Set updated by user ID
	workspace.UpdatedAt = time.Now()

	return s.workspaceRepo.Update(ctx, workspace)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	workspace, err := s.workspaceRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("workspace_service: delete_workspace: get_workspace: %w", err)
	}
	if workspace == nil {
		return ErrWorkspaceNotFound
	}

	now := time.Now()
	workspace.DeletedAt = &now
	workspace.DeletedBy = &userID
	workspace.UpdatedAt = now

	return s.workspaceRepo.Update(ctx, workspace)
}

func (s *Service) AddUserToWorkspace(ctx context.Context, workspaceID, userID uuid.UUID, role string, requestingUserID uuid.UUID) error {
	requestingUserRole, godMode, err := s.requestingUserRole(ctx, s.workspaceRepo, s.userRepo, workspaceID, requestingUserID, "add_user_to_workspace", "get_user_role")
	if err != nil {
		return err
	}
	if !godMode && !authorization_service.CanGrantRole(requestingUserRole, role) {
		return ErrInsufficientPermissions
	}

	relationship := &models.UserWorkspaceRelationship{
		UserID:      userID,
		WorkspaceID: workspaceID,
		Role:        role,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	return s.workspaceRepo.AddUserToWorkspace(ctx, relationship)
}

func (s *Service) RemoveUserFromWorkspace(ctx context.Context, workspaceID, userID, requestingUserID uuid.UUID) error {
	// Prevent users from removing themselves
	if userID == requestingUserID {
		return ErrCannotModifyOwnRole
	}

	requestingUserRole, godMode, err := s.requestingUserRole(ctx, s.workspaceRepo, s.userRepo, workspaceID, requestingUserID, "remove_user_from_workspace", "get_requesting_user_role")
	if err != nil {
		return err
	}
	if !godMode {
		targetUserRole, err := s.workspaceUserRole(ctx, s.workspaceRepo, workspaceID, userID, "remove_user_from_workspace", "get_target_user_role")
		if err != nil {
			return err
		}
		if !authorization_service.CanModifyUserRole(requestingUserRole, targetUserRole) {
			return ErrInsufficientPermissions
		}
	}

	return s.workspaceRepo.RemoveUserFromWorkspace(ctx, userID, workspaceID)
}

func (s *Service) UpdateUserRole(ctx context.Context, workspaceID, userID uuid.UUID, newRole string, requestingUserID uuid.UUID) error {
	// Prevent users from modifying their own role
	if userID == requestingUserID {
		return ErrCannotModifyOwnRole
	}

	requestingUserRole, godMode, err := s.requestingUserRole(ctx, s.workspaceRepo, s.userRepo, workspaceID, requestingUserID, "update_user_role", "get_requesting_user_role")
	if err != nil {
		return err
	}
	if !godMode {
		targetUserRole, err := s.workspaceUserRole(ctx, s.workspaceRepo, workspaceID, userID, "update_user_role", "get_target_user_role")
		if err != nil {
			return err
		}
		if !authorization_service.CanModifyUserRole(requestingUserRole, targetUserRole) {
			return ErrInsufficientPermissions
		}
		if !authorization_service.CanGrantRole(requestingUserRole, newRole) {
			return ErrInsufficientPermissions
		}
	}

	return s.workspaceRepo.UpdateUserRole(ctx, userID, workspaceID, newRole)
}

func (s *Service) requestingUserRole(ctx context.Context, workspaceRepo *repositories.WorkspaceRepository, userRepo *repositories.UserRepository, workspaceID, requestingUserID uuid.UUID, operation, roleStep string) (string, bool, error) {
	requestingUser, err := userRepo.GetByID(ctx, requestingUserID)
	if err != nil {
		return "", false, fmt.Errorf("workspace_service: %s: get_requesting_user: %w", operation, err)
	}
	if requestingUser.GodMode {
		return "", true, nil
	}

	role, err := s.workspaceUserRole(ctx, workspaceRepo, workspaceID, requestingUserID, operation, roleStep)
	if err != nil {
		return "", false, err
	}
	return role, false, nil
}

func (s *Service) workspaceUserRole(ctx context.Context, workspaceRepo *repositories.WorkspaceRepository, workspaceID, userID uuid.UUID, operation, step string) (string, error) {
	role, err := workspaceRepo.GetUserRole(ctx, userID, workspaceID)
	if err != nil {
		return "", fmt.Errorf("workspace_service: %s: %s: %w", operation, step, err)
	}
	if role == "" {
		return "", ErrInsufficientPermissions
	}
	return role, nil
}

func (s *Service) GetWorkspaceUsers(ctx context.Context, workspaceID, requestingUserID uuid.UUID) ([]models.User, error) {
	return s.workspaceRepo.GetWorkspaceUsers(ctx, workspaceID)
}

func (s *Service) GetUserWorkspaces(ctx context.Context, userID uuid.UUID) ([]models.Workspace, error) {
	return s.workspaceRepo.GetUserWorkspaces(ctx, userID)
}

func (s *Service) GetUserWorkspacesWithRole(ctx context.Context, userID uuid.UUID) ([]models.WorkspaceWithRole, error) {
	return s.workspaceRepo.GetUserWorkspacesWithRole(ctx, userID)
}

func (s *Service) Paginate(ctx context.Context, page int, perPage int, query string) ([]models.Workspace, error) {
	return s.workspaceRepo.Paginate(ctx, page, perPage, query)
}

func (s *Service) PaginateWithCount(ctx context.Context, page int, perPage int, query string) ([]models.Workspace, int, error) {
	return s.workspaceRepo.PaginateWithCount(ctx, page, perPage, query)
}

func (s *Service) Count(ctx context.Context) (int, error) {
	return s.workspaceRepo.Count(ctx)
}

// GetHistoricalCreationCounts returns daily workspace creation counts for the specified number of days
func (s *Service) GetHistoricalCreationCounts(ctx context.Context, days int) (map[string]int, error) {
	return s.workspaceRepo.GetHistoricalCreationCounts(ctx, days)
}

// GetWorkspaceUsersPaginated returns a paginated list of users in a workspace
func (s *Service) GetWorkspaceUsersPaginated(ctx context.Context, workspaceID uuid.UUID, paginationParams pagination_utils.PaginationParams) ([]models.UserWithRole, int, error) {
	// Validate and normalize params
	paginationParams.Validate()
	return s.workspaceRepo.GetWorkspaceUsersPaginated(ctx, workspaceID, paginationParams)
}

// PaginateWithCountForAdmin returns paginated workspaces with enhanced data for admin views
func (s *Service) PaginateWithCountForAdmin(ctx context.Context, page int, perPage int, query string, includeDeleted bool) ([]repositories.AdminWorkspaceData, int, error) {
	return s.workspaceRepo.PaginateWithCountForAdmin(ctx, page, perPage, query, includeDeleted)
}

func (s *Service) BulkDeleteWorkspaces(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) error {
	workspaces, err := s.workspaceRepo.GetByIDs(ctx, ids)
	if err != nil {
		return fmt.Errorf("workspace_service: bulk_delete_workspaces: get_workspaces: %w", err)
	}

	now := time.Now()
	for idx := range workspaces {
		workspaces[idx].DeletedAt = &now
		workspaces[idx].DeletedBy = &userID
		workspaces[idx].UpdatedAt = now
	}

	return s.workspaceRepo.BulkUpdate(ctx, workspaces)
}

func (s *Service) BulkRecoverWorkspaces(ctx context.Context, ids []uuid.UUID, userID uuid.UUID) error {
	workspaces, err := s.workspaceRepo.GetByIDsIncludingDeleted(ctx, ids)
	if err != nil {
		return fmt.Errorf("workspace_service: bulk_recover_workspaces: get_workspaces: %w", err)
	}

	now := time.Now()
	for idx := range workspaces {
		workspaces[idx].DeletedAt = nil
		workspaces[idx].DeletedBy = nil
		workspaces[idx].UpdatedBy = userID
		workspaces[idx].UpdatedAt = now
	}

	return s.workspaceRepo.BulkUpdate(ctx, workspaces)
}

// GetAdminWorkspaceData returns a single workspace with admin data
func (s *Service) GetAdminWorkspaceData(ctx context.Context, id uuid.UUID) (*repositories.AdminWorkspaceData, error) {
	return s.workspaceRepo.GetAdminWorkspaceDataByID(ctx, id)
}
