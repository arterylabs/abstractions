package repositories

import (
	"abstractions/db/models"
	"abstractions/utils/pagination_utils"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkspaceRepository struct {
	db *gorm.DB
}

// AdminWorkspaceData represents enhanced workspace data for admin views
type AdminWorkspaceData struct {
	models.Workspace
	CreatedByName *string              `gorm:"column:created_by_name"`
	UpdatedByName *string              `gorm:"column:updated_by_name"`
	DeletedByName *string              `gorm:"column:deleted_by_name"`
	UsersCount    int                  `gorm:"column:users_count"`
	Subscription  *models.Subscription `gorm:"embedded;embeddedPrefix:subscription_"`
}

func NewWorkspaceRepository(db *gorm.DB) *WorkspaceRepository {
	return &WorkspaceRepository{
		db: db,
	}
}

func (r *WorkspaceRepository) Create(ctx context.Context, workspace *models.Workspace) error {
	return r.db.WithContext(ctx).Create(workspace).Error
}

func (r *WorkspaceRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Workspace, error) {
	var workspace models.Workspace
	if err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&workspace).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("workspace_repository: get_by_id: %w", err)
	}
	return &workspace, nil
}

func (r *WorkspaceRepository) GetByIDIncludingDeleted(ctx context.Context, id uuid.UUID) (*models.Workspace, error) {
	var workspace models.Workspace
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&workspace).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("workspace_repository: get_by_id_including_deleted: %w", err)
	}
	return &workspace, nil
}

func (r *WorkspaceRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Workspace, error) {
	var workspaces []models.Workspace
	if err := r.db.WithContext(ctx).Where("id IN ? AND deleted_at IS NULL", ids).Find(&workspaces).Error; err != nil {
		return nil, fmt.Errorf("workspace_repository: get_by_ids: %w", err)
	}
	return workspaces, nil
}

func (r *WorkspaceRepository) GetByIDsIncludingDeleted(ctx context.Context, ids []uuid.UUID) ([]models.Workspace, error) {
	var workspaces []models.Workspace
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&workspaces).Error; err != nil {
		return nil, fmt.Errorf("workspace_repository: get_by_ids_including_deleted: %w", err)
	}
	return workspaces, nil
}

func (r *WorkspaceRepository) Update(ctx context.Context, workspace *models.Workspace) error {
	return r.db.WithContext(ctx).Save(workspace).Error
}

func (r *WorkspaceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Workspace{}, "id = ?", id).Error
}

func (r *WorkspaceRepository) List(ctx context.Context) ([]models.Workspace, error) {
	var workspaces []models.Workspace
	if err := r.db.WithContext(ctx).Find(&workspaces).Error; err != nil {
		return nil, fmt.Errorf("workspace_repository: list: %w", err)
	}
	return workspaces, nil
}

func (r *WorkspaceRepository) AddUserToWorkspace(ctx context.Context, relationship *models.UserWorkspaceRelationship) error {
	return r.db.WithContext(ctx).Create(relationship).Error
}

func (r *WorkspaceRepository) RemoveUserFromWorkspace(ctx context.Context, userID, workspaceID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.UserWorkspaceRelationship{}, "user_id = ? AND workspace_id = ?", userID, workspaceID).Error
}

func (r *WorkspaceRepository) GetWorkspaceUsers(ctx context.Context, workspaceID uuid.UUID) ([]models.User, error) {
	var users []models.User
	if err := r.db.WithContext(ctx).
		Table("users").
		Joins("JOIN user_workspace_relationships ON users.id = user_workspace_relationships.user_id").
		Where("user_workspace_relationships.workspace_id = ?", workspaceID).
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("workspace_repository: get_workspace_users: %w", err)
	}
	return users, nil
}

func (r *WorkspaceRepository) GetUserWorkspaces(ctx context.Context, userID uuid.UUID) ([]models.Workspace, error) {
	var workspaces []models.Workspace
	if err := r.db.WithContext(ctx).
		Table("workspaces").
		Joins("JOIN user_workspace_relationships ON workspaces.id = user_workspace_relationships.workspace_id").
		Where("user_workspace_relationships.user_id = ? AND workspaces.deleted_at IS NULL", userID).
		Order("workspaces.created_at ASC").
		Find(&workspaces).Error; err != nil {
		return nil, fmt.Errorf("workspace_repository: get_user_workspaces: %w", err)
	}
	return workspaces, nil
}

func (r *WorkspaceRepository) GetUserWorkspacesWithRole(ctx context.Context, userID uuid.UUID) ([]models.WorkspaceWithRole, error) {
	var workspaces []models.WorkspaceWithRole
	if err := r.db.WithContext(ctx).
		Table("workspaces").
		Select("workspaces.*, user_workspace_relationships.role as role").
		Joins("JOIN user_workspace_relationships ON workspaces.id = user_workspace_relationships.workspace_id").
		Where("user_workspace_relationships.user_id = ? AND workspaces.deleted_at IS NULL", userID).
		Order("workspaces.created_at ASC").
		Find(&workspaces).Error; err != nil {
		return nil, fmt.Errorf("workspace_repository: get_user_workspaces_with_role: %w", err)
	}
	return workspaces, nil
}

func (r *WorkspaceRepository) UpdateUserRole(ctx context.Context, userID, workspaceID uuid.UUID, newRole string) error {
	return r.db.WithContext(ctx).
		Model(&models.UserWorkspaceRelationship{}).
		Where("user_id = ? AND workspace_id = ?", userID, workspaceID).
		Update("role", newRole).
		Error
}

func (r *WorkspaceRepository) GetUserRole(ctx context.Context, userID, workspaceID uuid.UUID) (string, error) {
	var relationship models.UserWorkspaceRelationship
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND workspace_id = ?", userID, workspaceID).
		First(&relationship).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", fmt.Errorf("workspace_repository: get_user_role: %w", err)
	}
	return relationship.Role, nil
}

func (r *WorkspaceRepository) Paginate(ctx context.Context, page int, perPage int, query string) ([]models.Workspace, error) {
	var workspaces []models.Workspace
	q := r.db.WithContext(ctx).Offset((page - 1) * perPage).Limit(perPage)
	if query != "" {
		q = q.Where("name LIKE ?", "%"+query+"%")
	}
	if err := q.Find(&workspaces).Error; err != nil {
		return nil, fmt.Errorf("workspace_repository: paginate: %w", err)
	}
	return workspaces, nil
}

func (r *WorkspaceRepository) PaginateWithCount(ctx context.Context, page int, perPage int, query string) ([]models.Workspace, int, error) {
	var workspaces []models.Workspace
	var totalCount int64

	// Build the base query
	baseQuery := r.db.WithContext(ctx).Model(&models.Workspace{})
	if query != "" {
		baseQuery = baseQuery.Where("name LIKE ?", "%"+query+"%")
	}

	// Get total count first
	if err := baseQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("workspace_repository: paginate_with_count: count: %w", err)
	}

	// Apply pagination and fetch results
	if err := baseQuery.Offset((page - 1) * perPage).Limit(perPage).Find(&workspaces).Error; err != nil {
		return nil, 0, fmt.Errorf("workspace_repository: paginate_with_count: find: %w", err)
	}

	return workspaces, int(totalCount), nil
}

// GetWorkspaceUsersPaginated returns a paginated list of users in a workspace with search and ordering
func (r *WorkspaceRepository) GetWorkspaceUsersPaginated(ctx context.Context, workspaceID uuid.UUID, paginationParams pagination_utils.PaginationParams) ([]models.UserWithRole, int, error) {
	var users []models.UserWithRole
	var totalCount int64

	query := r.db.WithContext(ctx).
		Table("users").
		Select("users.*, user_workspace_relationships.role as role").
		Joins("JOIN user_workspace_relationships ON users.id = user_workspace_relationships.user_id").
		Where("user_workspace_relationships.workspace_id = ?", workspaceID)

	// Apply search if provided
	if paginationParams.Query != "" {
		query = query.Where("users.name LIKE ? OR users.email LIKE ?", "%"+paginationParams.Query+"%", "%"+paginationParams.Query+"%")
	}

	// Get total count before pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("workspace_repository: get_workspace_users_paginated: count: %w", err)
	}

	// Apply ordering
	if paginationParams.Sort != "" {
		sortOrder := pagination_utils.ParseSortOrder(paginationParams.SortOrder)
		switch paginationParams.Sort {
		case "name":
			query = query.Order("users.name " + sortOrder.String())
		case "email":
			query = query.Order("users.email " + sortOrder.String())
		case "created_at":
			query = query.Order("users.created_at " + sortOrder.String())
		default:
			return nil, 0, errors.New("invalid sort field")
		}
	} else {
		// Default ordering by created_at if no sort field is provided
		query = query.Order("users.created_at DESC")
	}

	// Apply pagination
	if err := query.Offset(paginationParams.GetOffset()).Limit(paginationParams.GetLimit()).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("workspace_repository: get_workspace_users_paginated: find: %w", err)
	}

	return users, int(totalCount), nil
}

func (r *WorkspaceRepository) Count(ctx context.Context) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Workspace{}).Count(&count).Error
	return int(count), err
}

// GetHistoricalCreationCounts returns daily workspace creation counts for the specified number of days
func (r *WorkspaceRepository) GetHistoricalCreationCounts(ctx context.Context, days int) (map[string]int, error) {
	type result struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}

	var results []result

	// Calculate the start date
	startDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as count
		FROM workspaces 
		WHERE DATE(created_at) >= ? 
		GROUP BY DATE(created_at) 
		ORDER BY date ASC
	`, startDate).Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("workspace_repository: get_historical_creation_counts: %w", err)
	}

	// Convert to map for easier lookup
	counts := make(map[string]int)
	for _, r := range results {
		counts[r.Date] = r.Count
	}

	return counts, nil
}

// GetAdminWorkspaceDataPaginated returns a paginated list of workspaces with additional admin data
func (r *WorkspaceRepository) GetAdminWorkspaceDataPaginated(ctx context.Context, paginationParams pagination_utils.PaginationParams) ([]AdminWorkspaceData, int, error) {
	var workspaces []AdminWorkspaceData
	var totalCount int64

	query := r.db.WithContext(ctx).
		Table("workspaces").
		Select("workspaces.*, u1.name as created_by_name, u2.name as updated_by_name, u3.name as deleted_by_name, (SELECT COUNT(*) FROM user_workspace_relationships WHERE workspace_id = workspaces.id) as users_count").
		Joins("LEFT JOIN users u1 ON workspaces.created_by = u1.id").
		Joins("LEFT JOIN users u2 ON workspaces.updated_by = u2.id").
		Joins("LEFT JOIN users u3 ON workspaces.deleted_by = u3.id")

	// Apply search if provided
	if paginationParams.Query != "" {
		query = query.Where("workspaces.name LIKE ?", "%"+paginationParams.Query+"%")
	}

	// Get total count before pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("workspace_repository: get_admin_workspace_data_paginated: count: %w", err)
	}

	// Apply ordering
	if paginationParams.Sort != "" {
		sortOrder := pagination_utils.ParseSortOrder(paginationParams.SortOrder)
		switch paginationParams.Sort {
		case "name":
			query = query.Order("workspaces.name " + sortOrder.String())
		case "created_at":
			query = query.Order("workspaces.created_at " + sortOrder.String())
		case "updated_at":
			query = query.Order("workspaces.updated_at " + sortOrder.String())
		case "users_count":
			query = query.Order("users_count " + sortOrder.String())
		default:
			return nil, 0, errors.New("invalid sort field")
		}
	} else {
		// Default ordering by created_at if no sort field is provided
		query = query.Order("workspaces.created_at DESC")
	}

	// Apply pagination
	if err := query.Offset(paginationParams.GetOffset()).Limit(paginationParams.GetLimit()).Find(&workspaces).Error; err != nil {
		return nil, 0, fmt.Errorf("workspace_repository: get_admin_workspace_data_paginated: find: %w", err)
	}

	return workspaces, int(totalCount), nil
}

// PaginateWithCountForAdmin returns paginated workspaces with enhanced data for admin views
func (r *WorkspaceRepository) PaginateWithCountForAdmin(ctx context.Context, page int, perPage int, query string) ([]AdminWorkspaceData, int, error) {
	var workspaces []AdminWorkspaceData
	var totalCount int64

	// Build the base query with JOINs for user names, user count, and subscription data
	baseQuery := r.db.WithContext(ctx).Table("workspaces").
		Select(`workspaces.*, 
		        created_by_users.name as created_by_name,
		        updated_by_users.name as updated_by_name,
		        deleted_by_users.name as deleted_by_name,
		        COALESCE(user_counts.count, 0) as users_count,
		        subscriptions.id as subscription_id,
		        subscriptions.plan_id as subscription_plan_id,
		        subscriptions.status as subscription_status,
		        subscriptions.expires_at as subscription_expires_at,
		        subscriptions.limit as subscription_limit`).
		Joins("LEFT JOIN users as created_by_users ON workspaces.created_by = created_by_users.id").
		Joins("LEFT JOIN users as updated_by_users ON workspaces.updated_by = updated_by_users.id").
		Joins("LEFT JOIN users as deleted_by_users ON workspaces.deleted_by = deleted_by_users.id").
		Joins(`LEFT JOIN (
			SELECT workspace_id, COUNT(*) as count 
			FROM user_workspace_relationships 
			GROUP BY workspace_id
		) as user_counts ON workspaces.id = user_counts.workspace_id`).
		Joins(`LEFT JOIN subscriptions ON workspaces.id = subscriptions.workspace_id 
		       AND subscriptions.status = 'active' 
		       AND subscriptions.expires_at > NOW()`)

	// Apply search filter if provided
	if query != "" {
		searchTerm := "%" + query + "%"
		baseQuery = baseQuery.Where("workspaces.name LIKE ? OR workspaces.billing_email LIKE ?", searchTerm, searchTerm)
	}

	// Get total count
	if err := baseQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("workspace_repository: paginate_with_count_for_admin: count: %w", err)
	}

	// Apply pagination and ordering
	offset := (page - 1) * perPage
	if err := baseQuery.
		Order("workspaces.created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&workspaces).Error; err != nil {
		return nil, 0, fmt.Errorf("workspace_repository: paginate_with_count_for_admin: find: %w", err)
	}

	return workspaces, int(totalCount), nil
}

// GetAdminWorkspaceDataByID returns a single workspace with admin data by ID
func (r *WorkspaceRepository) GetAdminWorkspaceDataByID(ctx context.Context, id uuid.UUID) (*AdminWorkspaceData, error) {
	var workspace AdminWorkspaceData

	query := r.db.WithContext(ctx).
		Table("workspaces").
		Select(`workspaces.*, 
		        created_by_users.name as created_by_name,
		        updated_by_users.name as updated_by_name,
		        deleted_by_users.name as deleted_by_name,
		        COALESCE(user_counts.count, 0) as users_count,
		        subscriptions.id as subscription_id,
		        subscriptions.plan_id as subscription_plan_id,
		        subscriptions.status as subscription_status,
		        subscriptions.expires_at as subscription_expires_at,
		        subscriptions.limit as subscription_limit`).
		Joins("LEFT JOIN users as created_by_users ON workspaces.created_by = created_by_users.id").
		Joins("LEFT JOIN users as updated_by_users ON workspaces.updated_by = updated_by_users.id").
		Joins("LEFT JOIN users as deleted_by_users ON workspaces.deleted_by = deleted_by_users.id").
		Joins(`LEFT JOIN (
			SELECT workspace_id, COUNT(*) as count 
			FROM user_workspace_relationships 
			GROUP BY workspace_id
		) as user_counts ON workspaces.id = user_counts.workspace_id`).
		Joins(`LEFT JOIN subscriptions ON workspaces.id = subscriptions.workspace_id 
		       AND subscriptions.status = 'active' 
		       AND subscriptions.expires_at > NOW()`).
		Where("workspaces.id = ?", id)

	if err := query.First(&workspace).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("workspace_repository: get_admin_workspace_data_by_id: %w", err)
	}

	return &workspace, nil
}

// BulkUpdate updates multiple workspaces in the database
func (r *WorkspaceRepository) BulkUpdate(ctx context.Context, workspaces []models.Workspace) error {
	return r.db.WithContext(ctx).Save(&workspaces).Error
}
