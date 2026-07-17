package repositories

import (
	"abstractions/db/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TestCleanupRepository struct {
	DB *gorm.DB
}

func NewTestCleanupRepository(db *gorm.DB) *TestCleanupRepository {
	return &TestCleanupRepository{DB: db}
}

func (r *TestCleanupRepository) VerifyUnverifiedUserByEmail(ctx context.Context, email string) error {
	return r.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("email = ? AND verified_at IS NULL", email).
		Update("verified_at", gorm.Expr("NOW()")).Error
}

func (r *TestCleanupRepository) ListUsersByEmailSuffix(ctx context.Context, suffix string) ([]models.User, error) {
	var users []models.User
	err := r.DB.WithContext(ctx).Unscoped().Where("email LIKE ?", "%"+suffix).Find(&users).Error
	return users, err
}

func (r *TestCleanupRepository) ListWorkspaceRelationshipsByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]models.UserWorkspaceRelationship, error) {
	var relationships []models.UserWorkspaceRelationship
	err := r.DB.WithContext(ctx).Where("user_id IN ?", userIDs).Find(&relationships).Error
	return relationships, err
}

func (r *TestCleanupRepository) DeleteWorkspaceData(ctx context.Context, workspaceIDs []uuid.UUID) error {
	db := r.DB.WithContext(ctx)
	if err := db.Unscoped().Where("workspace_id IN ?", workspaceIDs).Delete(&models.Repository{}).Error; err != nil {
		return err
	}
	if err := db.Unscoped().Where("workspace_id IN ?", workspaceIDs).Delete(&models.Integration{}).Error; err != nil {
		return err
	}
	if err := db.Unscoped().Where("workspace_id IN ?", workspaceIDs).Delete(&models.Invitation{}).Error; err != nil {
		return err
	}
	if err := db.Where("workspace_id IN ?", workspaceIDs).Delete(&models.UserWorkspaceRelationship{}).Error; err != nil {
		return err
	}
	return db.Unscoped().Where("id IN ?", workspaceIDs).Delete(&models.Workspace{}).Error
}

func (r *TestCleanupRepository) DeleteUserData(ctx context.Context, userIDs []uuid.UUID) error {
	db := r.DB.WithContext(ctx)
	if err := db.Where("user_id IN ?", userIDs).Delete(&models.Session{}).Error; err != nil {
		return err
	}
	if err := db.Where("user_id IN ?", userIDs).Delete(&models.Token{}).Error; err != nil {
		return err
	}
	if err := db.Where("user_id IN ?", userIDs).Delete(&models.Authenticator{}).Error; err != nil {
		return err
	}
	if err := db.Where("user_id IN ?", userIDs).Delete(&models.UserWorkspaceRelationship{}).Error; err != nil {
		return err
	}
	return db.Unscoped().Where("id IN ?", userIDs).Delete(&models.User{}).Error
}
