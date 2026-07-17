package test_service

import (
	"abstractions/repositories"
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	db   *gorm.DB
	repo *repositories.TestCleanupRepository
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db, repo: repositories.NewTestCleanupRepository(db)}
}

// VerifyTestUser sets verified_at to now for a single unverified test user by email.
// Returns an error if the email is not on the test.arterylabs.com domain.
func (s *Service) VerifyTestUser(ctx context.Context, email string) error {
	if !isTestEmail(email) {
		return fmt.Errorf("test_service: verify test user: not a test email")
	}

	if err := s.repo.VerifyUnverifiedUserByEmail(ctx, email); err != nil {
		return fmt.Errorf("test_service: verify test user: %w", err)
	}

	return nil
}

func isTestEmail(email string) bool {
	return len(email) > len("@test.arterylabs.com") &&
		email[len(email)-len("@test.arterylabs.com"):] == "@test.arterylabs.com"
}

// CleanupTestUsers deletes all users with emails on the test.arterylabs.com domain
// and all associated data (workspaces, repositories, newsletters, sessions, etc.).
// Returns the number of users deleted.
func (s *Service) CleanupTestUsers(ctx context.Context) (int, error) {
	testUsers, err := s.repo.ListUsersByEmailSuffix(ctx, "@test.arterylabs.com")
	if err != nil {
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
	relationships, err := s.repo.ListWorkspaceRelationshipsByUserIDs(ctx, userIDs)
	if err != nil {
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
		if err := s.deleteWorkspaces(ctx, workspaceIDs); err != nil {
			return 0, err
		}
	}

	if err := s.deleteUsers(ctx, userIDs); err != nil {
		return 0, err
	}

	return len(testUsers), nil
}

func (s *Service) deleteWorkspaces(ctx context.Context, workspaceIDs []uuid.UUID) error {
	// Repositories cascade to repository_files, file_chunk_embeddings, newsletters,
	// and newsletter_topics via ON DELETE CASCADE in the database schema.
	if err := s.repo.DeleteWorkspaceData(ctx, workspaceIDs); err != nil {
		return fmt.Errorf("test_service: delete repositories: %w", err)
	}
	return nil
}

func (s *Service) deleteUsers(ctx context.Context, userIDs []uuid.UUID) error {
	if err := s.repo.DeleteUserData(ctx, userIDs); err != nil {
		return fmt.Errorf("test_service: delete sessions: %w", err)
	}
	return nil
}
