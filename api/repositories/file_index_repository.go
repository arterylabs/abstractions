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

type RepositoryFileRepository struct {
	db *gorm.DB
}

func NewRepositoryFileRepository(db *gorm.DB) *RepositoryFileRepository {
	return &RepositoryFileRepository{db: db}
}

// Upsert inserts or updates a RepositoryFile row by (repo_id, path).
// GitSHA, Size, and Language are refreshed on conflict; Tags and ProcessedAt are
// managed exclusively by UpdateMetadata (called from SummarizeFileJob) and must
// not be overwritten here, otherwise re-indexing would clear previously computed tags.
func (r *RepositoryFileRepository) Upsert(ctx context.Context, fi *models.RepositoryFile) error {
	return r.db.WithContext(ctx).
		Where(models.RepositoryFile{RepoID: fi.RepoID, Path: fi.Path}).
		Assign(models.RepositoryFile{
			GitSHA:   fi.GitSHA,
			Size:     fi.Size,
			Language: fi.Language,
		}).
		FirstOrCreate(fi).Error
}

// GetByRepoAndPath returns the RepositoryFile row for a given repo and file path.
func (r *RepositoryFileRepository) GetByRepoAndPath(ctx context.Context, repoID uuid.UUID, path string) (*models.RepositoryFile, error) {
	var fi models.RepositoryFile
	if err := r.db.WithContext(ctx).Where("repo_id = ? AND path = ?", repoID, path).First(&fi).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("repository_file_repository: get_by_repo_and_path: %w", err)
	}
	return &fi, nil
}

// ListByRepoID returns all indexed files for a repository.
func (r *RepositoryFileRepository) ListByRepoID(ctx context.Context, repoID uuid.UUID) ([]*models.RepositoryFile, error) {
	var files []*models.RepositoryFile
	if err := r.db.WithContext(ctx).Where("repo_id = ?", repoID).Find(&files).Error; err != nil {
		return nil, fmt.Errorf("repository_file_repository: list_by_repo_id: %w", err)
	}
	return files, nil
}

// DeleteByRepoID removes all repository file entries for a repository.
func (r *RepositoryFileRepository) DeleteByRepoID(ctx context.Context, repoID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("repo_id = ?", repoID).Delete(&models.RepositoryFile{}).Error
}

// DeleteByPath removes a single repository file entry.
func (r *RepositoryFileRepository) DeleteByPath(ctx context.Context, repoID uuid.UUID, path string) error {
	return r.db.WithContext(ctx).Where("repo_id = ? AND path = ?", repoID, path).Delete(&models.RepositoryFile{}).Error
}

// LanguageStat holds the file count for a single language.
type LanguageStat struct {
	Language  string
	FileCount int
}

// GetLanguageStats returns per-language file counts for a repository, ordered by count desc.
func (r *RepositoryFileRepository) GetLanguageStats(ctx context.Context, repoID uuid.UUID) ([]LanguageStat, error) {
	type row struct {
		Language  string
		FileCount int
	}
	var rows []row
	if err := r.db.WithContext(ctx).
		Model(&models.RepositoryFile{}).
		Select("language, COUNT(*) AS file_count").
		Where("repo_id = ? AND deleted_at IS NULL", repoID).
		Group("language").
		Order("file_count DESC").
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("repository_file_repository: get_language_stats: %w", err)
	}
	stats := make([]LanguageStat, len(rows))
	for i, r := range rows {
		stats[i] = LanguageStat{Language: r.Language, FileCount: r.FileCount}
	}
	return stats, nil
}

// UpdateServiceMetadata stores the raw GitHub Contents API metadata for a file.
// It is a non-critical write — callers should log and continue on error rather than failing the job.
func (r *RepositoryFileRepository) UpdateServiceMetadata(ctx context.Context, repoID uuid.UUID, path string, metadata models.JSONB) error {
	return r.db.WithContext(ctx).
		Model(&models.RepositoryFile{}).
		Where("repo_id = ? AND path = ?", repoID, path).
		Update("service_metadata", metadata).Error
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
