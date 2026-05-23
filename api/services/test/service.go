package test_service

import (
	"abstractions/db/models"
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// CleanupTestUsers deletes all users with emails on the test.arterylabs.com domain
// and all associated data (workspaces, repositories, newsletters, sessions, etc.).
// Returns the number of users deleted.
func (s *Service) CleanupTestUsers(ctx context.Context) (int, error) {
	db := s.db.WithContext(ctx)

	var testUsers []models.User
	if err := db.Unscoped().Where("email LIKE ?", "%@test.arterylabs.com").Find(&testUsers).Error; err != nil {
		return 0, fmt.Errorf("test_service: find test users: %w", err)
	}

	if len(testUsers) == 0 {
		return 0, nil
	}

	userIDs := make([]uuid.UUID, len(testUsers))
	for i, u := range testUsers {
		userIDs[i] = u.ID
	}

	// Collect workspace IDs associated with any test user.
	var relationships []models.UserWorkspaceRelationship
	if err := db.Where("user_id IN ?", userIDs).Find(&relationships).Error; err != nil {
		return 0, fmt.Errorf("test_service: find workspace relationships: %w", err)
	}

	workspaceIDs := make([]uuid.UUID, 0, len(relationships))
	seen := map[uuid.UUID]struct{}{}
	for _, r := range relationships {
		if _, ok := seen[r.WorkspaceID]; !ok {
			workspaceIDs = append(workspaceIDs, r.WorkspaceID)
			seen[r.WorkspaceID] = struct{}{}
		}
	}

	if len(workspaceIDs) > 0 {
		if err := s.deleteWorkspaces(ctx, db, workspaceIDs); err != nil {
			return 0, err
		}
	}

	if err := s.deleteUsers(ctx, db, userIDs); err != nil {
		return 0, err
	}

	return len(testUsers), nil
}

func (s *Service) deleteWorkspaces(ctx context.Context, db *gorm.DB, workspaceIDs []uuid.UUID) error {
	// Repositories cascade to repository_files, file_chunk_embeddings, newsletters,
	// and newsletter_topics via ON DELETE CASCADE in the database schema.
	if err := db.Unscoped().Where("workspace_id IN ?", workspaceIDs).Delete(&models.Repository{}).Error; err != nil {
		return fmt.Errorf("test_service: delete repositories: %w", err)
	}

	if err := db.Unscoped().Where("workspace_id IN ?", workspaceIDs).Delete(&models.Integration{}).Error; err != nil {
		return fmt.Errorf("test_service: delete integrations: %w", err)
	}

	if err := db.Unscoped().Where("workspace_id IN ?", workspaceIDs).Delete(&models.Invitation{}).Error; err != nil {
		return fmt.Errorf("test_service: delete invitations: %w", err)
	}

	if err := db.Unscoped().Where("workspace_id IN ?", workspaceIDs).Delete(&models.Subscription{}).Error; err != nil {
		return fmt.Errorf("test_service: delete subscriptions: %w", err)
	}

	if err := db.Where("workspace_id IN ?", workspaceIDs).Delete(&models.UserWorkspaceRelationship{}).Error; err != nil {
		return fmt.Errorf("test_service: delete workspace user relationships: %w", err)
	}

	if err := db.Unscoped().Where("id IN ?", workspaceIDs).Delete(&models.Workspace{}).Error; err != nil {
		return fmt.Errorf("test_service: delete workspaces: %w", err)
	}

	return nil
}

func (s *Service) deleteUsers(ctx context.Context, db *gorm.DB, userIDs []uuid.UUID) error {
	if err := db.Where("user_id IN ?", userIDs).Delete(&models.Session{}).Error; err != nil {
		return fmt.Errorf("test_service: delete sessions: %w", err)
	}

	if err := db.Where("user_id IN ?", userIDs).Delete(&models.Token{}).Error; err != nil {
		return fmt.Errorf("test_service: delete tokens: %w", err)
	}

	if err := db.Where("user_id IN ?", userIDs).Delete(&models.Authenticator{}).Error; err != nil {
		return fmt.Errorf("test_service: delete authenticators: %w", err)
	}

	// Remove any remaining workspace relationships (e.g. workspaces owned by non-test users).
	if err := db.Where("user_id IN ?", userIDs).Delete(&models.UserWorkspaceRelationship{}).Error; err != nil {
		return fmt.Errorf("test_service: delete remaining user workspace relationships: %w", err)
	}

	if err := db.Unscoped().Where("id IN ?", userIDs).Delete(&models.User{}).Error; err != nil {
		return fmt.Errorf("test_service: delete users: %w", err)
	}

	return nil
}
