package repositories

import (
	"abstractions/db/models"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RepositoryFileRepository struct {
	db *gorm.DB
}

func NewRepositoryFileRepository(db *gorm.DB) *RepositoryFileRepository {
	return &RepositoryFileRepository{db: db}
}

// Upsert inserts or updates a RepositoryFile row by (repo_id, path).
func (r *RepositoryFileRepository) Upsert(ctx context.Context, fi *models.RepositoryFile) error {
	return r.db.WithContext(ctx).
		Where(models.RepositoryFile{RepoID: fi.RepoID, Path: fi.Path}).
		Assign(models.RepositoryFile{
			GitSHA:      fi.GitSHA,
			Size:        fi.Size,
			Tags:        fi.Tags,
			ProcessedAt: fi.ProcessedAt,
		}).
		FirstOrCreate(fi).Error
}

// GetByRepoAndPath returns the RepositoryFile row for a given repo and file path.
func (r *RepositoryFileRepository) GetByRepoAndPath(ctx context.Context, repoID uuid.UUID, filePath string) (*models.RepositoryFile, error) {
	var fi models.RepositoryFile
	if err := r.db.WithContext(ctx).Where("repo_id = ? AND path = ?", repoID, filePath).First(&fi).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &fi, nil
}

// ListByRepoID returns all indexed files for a repository.
func (r *RepositoryFileRepository) ListByRepoID(ctx context.Context, repoID uuid.UUID) ([]*models.RepositoryFile, error) {
	var files []*models.RepositoryFile
	if err := r.db.WithContext(ctx).Where("repo_id = ?", repoID).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

// DeleteByRepoID removes all repository file entries for a repository.
func (r *RepositoryFileRepository) DeleteByRepoID(ctx context.Context, repoID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("repo_id = ?", repoID).Delete(&models.RepositoryFile{}).Error
}

// DeleteByPath removes a single repository file entry.
func (r *RepositoryFileRepository) DeleteByPath(ctx context.Context, repoID uuid.UUID, filePath string) error {
	return r.db.WithContext(ctx).Where("repo_id = ? AND path = ?", repoID, filePath).Delete(&models.RepositoryFile{}).Error
}

// UpdateMetadata updates language, tags, and processed_at for a given file after chunk processing.
func (r *RepositoryFileRepository) UpdateMetadata(ctx context.Context, repoID uuid.UUID, path string, language string, tags models.StringArray, processedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.RepositoryFile{}).
		Where("repo_id = ? AND path = ?", repoID, path).
		Updates(map[string]interface{}{
			"language":     language,
			"tags":         tags,
			"processed_at": processedAt,
		}).Error
}
