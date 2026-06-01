package repositories

import (
	"abstractions/db/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepositoryQ interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]models.User, error)
	GetByIDsIncludingDeleted(ctx context.Context, ids []uuid.UUID) ([]models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	ExistsByEmail(ctx context.Context, email string) bool
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID, deletedBy *uuid.UUID) error
	BulkUpdate(ctx context.Context, users []models.User) error
	Paginate(ctx context.Context, page int, perPage int, query string) ([]models.User, error)
	PaginateWithCount(ctx context.Context, page int, perPage int, query string) ([]models.User, int, error)
	Count(ctx context.Context) (int, error)
	CountActiveUsers(ctx context.Context) (int, error)
	// PaginateWithCountForAdmin returns paginated users with enhanced data for admin views
	PaginateWithCountForAdmin(ctx context.Context, page int, perPage int, query string, includeDeleted bool) ([]AdminUserData, int, error)
	// GetAdminUserDataByID returns a single user with admin data by ID
	GetAdminUserDataByID(ctx context.Context, id uuid.UUID) (*AdminUserData, error)
}

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	usr := models.User{}
	if err := r.DB.WithContext(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&usr).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &usr, nil
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) bool {
	var count int64
	r.DB.WithContext(ctx).Model(&models.User{}).Where("email = ? AND deleted_at IS NULL", email).Count(&count)
	return count > 0
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	usr := models.User{}
	if err := r.DB.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&usr).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &usr, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	return r.DB.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) Paginate(ctx context.Context, page int, perPage int, query string) ([]models.User, error) {
	users := []models.User{}
	q := r.DB.WithContext(ctx).Where("deleted_at IS NULL").Offset((page - 1) * perPage).Limit(perPage)
	if query != "" {
		q = q.Where("email = ? or email LIKE ?", query, "%"+query+"%")
	}
	res := q.Find(&users)

	return users, res.Error
}

func (r *UserRepository) PaginateWithCount(ctx context.Context, page int, perPage int, query string) ([]models.User, int, error) {
	var users []models.User
	var totalCount int64

	// Build the base query
	baseQuery := r.DB.WithContext(ctx).Model(&models.User{}).Where("deleted_at IS NULL")
	if query != "" {
		baseQuery = baseQuery.Where("email = ? or email LIKE ?", query, "%"+query+"%")
	}

	// Get total count first
	if err := baseQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("user_repository: paginate_with_count: count: %w", err)
	}

	// Apply pagination and fetch results
	if err := baseQuery.Offset((page - 1) * perPage).Limit(perPage).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("user_repository: paginate_with_count: find: %w", err)
	}

	return users, int(totalCount), nil
}

func (r *UserRepository) BulkUpdate(ctx context.Context, users []models.User) error {
	now := time.Now()
	for i := range users {
		users[i].UpdatedAt = now
	}
	return r.DB.WithContext(ctx).Save(&users).Error
}

func (r *UserRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]models.User, error) {
	users := []models.User{}
	res := r.DB.WithContext(ctx).Where("id IN ? AND deleted_at IS NULL", ids).Find(&users)
	return users, res.Error
}

func (r *UserRepository) GetByIDsIncludingDeleted(ctx context.Context, ids []uuid.UUID) ([]models.User, error) {
	users := []models.User{}
	res := r.DB.WithContext(ctx).Where("id IN ?", ids).Find(&users)
	return users, res.Error
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	if user.Email == "" {
		return gorm.ErrInvalidData
	}
	if user.Name == "" {
		return gorm.ErrInvalidData
	}
	return r.DB.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) Count(ctx context.Context) (int, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.User{}).Count(&count).Error
	return int(count), err
}

func (r *UserRepository) CountActiveUsers(ctx context.Context) (int, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.User{}).Where("deleted_at IS NULL").Count(&count).Error
	return int(count), err
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy *uuid.UUID) error {
	user, err := r.GetByID(ctx, id)
	if err != nil || user == nil {
		return err
	}
	now := time.Now()
	user.DeletedAt = &now
	user.DeletedBy = deletedBy
	return r.Update(ctx, user)
}

func (r *UserRepository) PermanentlyDelete(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
}

// GetHistoricalCreationCounts returns daily user creation counts for the specified number of days
func (r *UserRepository) GetHistoricalCreationCounts(ctx context.Context, days int) (map[string]int, error) {
	type result struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}

	var results []result

	// Calculate the start date
	startDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	err := r.DB.WithContext(ctx).Raw(`
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as count
		FROM users 
		WHERE DATE(created_at) >= ? 
		GROUP BY DATE(created_at) 
		ORDER BY date ASC
	`, startDate).Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("user_repository: get_historical_creation_counts: %w", err)
	}

	// Convert to map for easier lookup
	counts := make(map[string]int)
	for _, r := range results {
		counts[r.Date] = r.Count
	}

	return counts, nil
}

// AdminUserData represents enhanced user data for admin views
type AdminUserData struct {
	models.User
	DefaultWorkspaceName *string `gorm:"column:default_workspace_name"`
	DeletedByName        *string `gorm:"column:deleted_by_name"`
	WorkspacesCount      int     `gorm:"column:workspaces_count"`
}

// PaginateWithCountForAdmin returns paginated users with enhanced data for admin views
func (r *UserRepository) PaginateWithCountForAdmin(ctx context.Context, page int, perPage int, query string, includeDeleted bool) ([]AdminUserData, int, error) {
	var users []AdminUserData
	var totalCount int64

	// Build the base query with JOINs for workspace and deleted-by user names, and workspace count
	baseQuery := r.DB.WithContext(ctx).Table("users").
		Select(`users.*,
		        default_workspaces.name as default_workspace_name,
		        deleted_by_users.name as deleted_by_name,
		        COALESCE(workspace_counts.count, 0) as workspaces_count`).
		Joins("LEFT JOIN workspaces as default_workspaces ON users.default_workspace_id = default_workspaces.id").
		Joins("LEFT JOIN users as deleted_by_users ON users.deleted_by = deleted_by_users.id").
		Joins(`LEFT JOIN (
			SELECT user_workspace_relationships.user_id, COUNT(*) as count
			FROM user_workspace_relationships
			GROUP BY user_workspace_relationships.user_id
		) as workspace_counts ON users.id = workspace_counts.user_id`)

	if !includeDeleted {
		baseQuery = baseQuery.Where("users.deleted_at IS NULL")
	}

	if query != "" {
		baseQuery = baseQuery.Where("users.email = ? or users.email LIKE ? or users.name LIKE ?",
			query, "%"+query+"%", "%"+query+"%")
	}

	// Get total count first
	if err := baseQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("user_repository: paginate_with_count_for_admin: count: %w", err)
	}

	// Apply pagination and fetch results
	if err := baseQuery.Offset((page - 1) * perPage).Limit(perPage).Scan(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("user_repository: paginate_with_count_for_admin: find: %w", err)
	}

	return users, int(totalCount), nil
}

// GetAdminUserDataByID returns a single user with admin data by ID
func (r *UserRepository) GetAdminUserDataByID(ctx context.Context, id uuid.UUID) (*AdminUserData, error) {
	var user AdminUserData

	// Build the same query as PaginateWithCountForAdmin but for a single user
	query := r.DB.WithContext(ctx).Table("users").
		Select(`users.*, 
		        default_workspaces.name as default_workspace_name,
		        deleted_by_users.name as deleted_by_name,
		        COALESCE(workspace_counts.count, 0) as workspaces_count`).
		Joins("LEFT JOIN workspaces as default_workspaces ON users.default_workspace_id = default_workspaces.id").
		Joins("LEFT JOIN users as deleted_by_users ON users.deleted_by = deleted_by_users.id").
		Joins(`LEFT JOIN (
			SELECT user_workspace_relationships.user_id, COUNT(*) as count 
			FROM user_workspace_relationships 
			GROUP BY user_workspace_relationships.user_id
		) as workspace_counts ON users.id = workspace_counts.user_id`).
		Where("users.id = ?", id)

	if err := query.First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("user_repository: get_admin_user_data_by_id: %w", err)
	}

	return &user, nil
}
