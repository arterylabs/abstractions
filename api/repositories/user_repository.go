package repositories

import (
	"abstractions/db/models"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepositoryQ interface {
	GetByID(id uuid.UUID) (*models.User, error)
	GetByIDs(ids []uuid.UUID) ([]models.User, error)
	GetByIDsIncludingDeleted(ids []uuid.UUID) ([]models.User, error)
	GetByEmail(email string) (*models.User, error)
	ExistsByEmail(email string) bool
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uuid.UUID) error
	BulkUpdate(users []models.User) error
	Paginate(page int, perPage int, query string) ([]models.User, error)
	PaginateWithCount(page int, perPage int, query string) ([]models.User, int, error)
	Count() (int, error)
	CountActiveUsers() (int, error)
	// PaginateWithCountForAdmin returns paginated users with enhanced data for admin views
	PaginateWithCountForAdmin(page int, perPage int, query string, includeDeleted bool) ([]AdminUserData, int, error)
	// GetAdminUserDataByID returns a single user with admin data by ID
	GetAdminUserDataByID(id uuid.UUID) (*AdminUserData, error)
}

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (userRepository *UserRepository) GetByEmail(email string) (*models.User, error) {
	usr := models.User{}
	res := userRepository.DB.Where("email = ?", email).Find(&usr)
	if res.Error != nil || usr.ID == uuid.Nil {
		// If no user is found, return nil
		if res.Error == nil {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, res.Error
	}
	return &usr, nil
}

func (userRepository *UserRepository) ExistsByEmail(email string) bool {
	var count int64
	userRepository.DB.Model(&models.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

func (userRepository *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	usr := models.User{}
	res := userRepository.DB.Where("id = ?", id).Find(&usr)
	return &usr, res.Error
}

func (userRepository *UserRepository) Update(user *models.User) error {
	return userRepository.DB.Save(user).Error
}

func (userRepository *UserRepository) Paginate(page int, perPage int, query string) ([]models.User, error) {
	users := []models.User{}
	q := userRepository.DB.Offset((page - 1) * perPage).Limit(perPage)
	if query != "" {
		q = q.Where("email = ? or email LIKE ?", query, "%"+query+"%")
	}
	res := q.Find(&users)

	return users, res.Error
}

func (userRepository *UserRepository) PaginateWithCount(page int, perPage int, query string) ([]models.User, int, error) {
	var users []models.User
	var totalCount int64

	// Build the base query
	baseQuery := userRepository.DB.Model(&models.User{})
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

func (userRepository *UserRepository) BulkUpdate(users []models.User) error {
	return userRepository.DB.Save(&users).Error
}

func (userRepository *UserRepository) GetByIDs(ids []uuid.UUID) ([]models.User, error) {
	users := []models.User{}
	res := userRepository.DB.Table("users").Where("id IN (?)", ids).Find(&users)
	return users, res.Error
}

func (userRepository *UserRepository) GetByIDsIncludingDeleted(ids []uuid.UUID) ([]models.User, error) {
	users := []models.User{}
	res := userRepository.DB.Unscoped().Table("users").Where("id IN (?)", ids).Find(&users)
	return users, res.Error
}

func (userRepository *UserRepository) Create(user *models.User) error {
	if user.Email == "" {
		return gorm.ErrInvalidData
	}
	if user.Name == "" {
		return gorm.ErrInvalidData
	}
	return userRepository.DB.Create(user).Error
}

func (userRepository *UserRepository) Count() (int, error) {
	var count int64
	err := userRepository.DB.Model(&models.User{}).Count(&count).Error
	return int(count), err
}

func (userRepository *UserRepository) CountActiveUsers() (int, error) {
	var count int64
	err := userRepository.DB.Model(&models.User{}).Where("deleted_at IS NULL").Count(&count).Error
	return int(count), err
}

func (userRepository *UserRepository) Delete(id uuid.UUID) error {
	return userRepository.DB.Delete(&models.User{}, id).Error
}

// GetHistoricalCreationCounts returns daily user creation counts for the specified number of days
func (userRepository *UserRepository) GetHistoricalCreationCounts(days int) (map[string]int, error) {
	type result struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}

	var results []result

	// Calculate the start date
	startDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	err := userRepository.DB.Raw(`
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
func (userRepository *UserRepository) PaginateWithCountForAdmin(page int, perPage int, query string, includeDeleted bool) ([]AdminUserData, int, error) {
	var users []AdminUserData
	var totalCount int64

	// Build the base query with JOINs for workspace and deleted-by user names, and workspace count
	baseQuery := userRepository.DB.Table("users").
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
func (userRepository *UserRepository) GetAdminUserDataByID(id uuid.UUID) (*AdminUserData, error) {
	var user AdminUserData

	// Build the same query as PaginateWithCountForAdmin but for a single user
	query := userRepository.DB.Table("users").
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
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("user_repository: get_admin_user_data_by_id: %w", err)
	}

	return &user, nil
}
