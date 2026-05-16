package jobs

import (
	"abstractions/db/models"
	"abstractions/repositories"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// ErrRepositoryNotActive is returned by LoadRepositoryContext when the repository
// does not exist or has been soft-deleted.
var ErrRepositoryNotActive = errors.New("repository is deleted or not found")

// RepositoryJobContext holds the validated repository model for use within a job.
type RepositoryJobContext struct {
	Repository *models.Repository
}

// LoadRepositoryContext fetches the repository by ID and validates it is active.
// GetByID already filters deleted_at IS NULL, so a nil result means deleted or missing.
// All repository-based jobs should call this instead of querying GetByID directly.
func LoadRepositoryContext(ctx context.Context, repoID uuid.UUID, repo *repositories.RepositoryRepository) (*RepositoryJobContext, error) {
	r, err := repo.GetByID(ctx, repoID)
	if err != nil {
		return nil, fmt.Errorf("load_repository_context: %w", err)
	}
	if r == nil {
		return nil, ErrRepositoryNotActive
	}
	return &RepositoryJobContext{Repository: r}, nil
}
