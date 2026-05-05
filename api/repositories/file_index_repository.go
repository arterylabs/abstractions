package repositories

import (
	"abstractions/db/models"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileIndexRepository struct {
	db *gorm.DB
}

func NewFileIndexRepository(db *gorm.DB) *FileIndexRepository {
	return &FileIndexRepository{db: db}
}

// Upsert inserts or updates a FileIndex row by (repo_id, file_path).
func (r *FileIndexRepository) Upsert(ctx context.Context, fi *models.FileIndex) error {
	return r.db.WithContext(ctx).
		Where(models.FileIndex{RepoID: fi.RepoID, FilePath: fi.FilePath}).
		Assign(models.FileIndex{
			GitSHA:          fi.GitSHA,
			DomainTags:      fi.DomainTags,
			LastProcessedAt: fi.LastProcessedAt,
			PineconeIDs:     fi.PineconeIDs,
		}).
		FirstOrCreate(fi).Error
}

// GetByRepoAndPath returns the FileIndex row for a given repo and file path.
func (r *FileIndexRepository) GetByRepoAndPath(ctx context.Context, repoID uuid.UUID, filePath string) (*models.FileIndex, error) {
	var fi models.FileIndex
	if err := r.db.WithContext(ctx).Where("repo_id = ? AND file_path = ?", repoID, filePath).First(&fi).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &fi, nil
}

// ListByRepoID returns all indexed files for a repository.
func (r *FileIndexRepository) ListByRepoID(ctx context.Context, repoID uuid.UUID) ([]*models.FileIndex, error) {
	var files []*models.FileIndex
	if err := r.db.WithContext(ctx).Where("repo_id = ?", repoID).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

// DeleteByRepoID removes all file index entries for a repository.
func (r *FileIndexRepository) DeleteByRepoID(ctx context.Context, repoID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("repo_id = ?", repoID).Delete(&models.FileIndex{}).Error
}

// DeleteByPath removes a single file index entry.
func (r *FileIndexRepository) DeleteByPath(ctx context.Context, repoID uuid.UUID, filePath string) error {
	return r.db.WithContext(ctx).Where("repo_id = ? AND file_path = ?", repoID, filePath).Delete(&models.FileIndex{}).Error
}
