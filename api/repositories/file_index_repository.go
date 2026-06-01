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
	if err := r.db.WithContext(ctx).Where("repo_id = ? AND path = ? AND deleted_at IS NULL", repoID, path).First(&fi).Error; err != nil {
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
	if err := r.db.WithContext(ctx).Where("repo_id = ? AND deleted_at IS NULL", repoID).Find(&files).Error; err != nil {
		return nil, fmt.Errorf("repository_file_repository: list_by_repo_id: %w", err)
	}
	return files, nil
}

// DeleteByRepoID removes all repository file entries for a repository.
func (r *RepositoryFileRepository) DeleteByRepoID(ctx context.Context, repoID uuid.UUID) error {
	files, err := r.ListByRepoID(ctx, repoID)
	if err != nil {
		return err
	}
	now := time.Now()
	for _, file := range files {
		file.DeletedAt = &now
		if err := r.db.WithContext(ctx).Save(file).Error; err != nil {
			return fmt.Errorf("repository_file_repository: delete_by_repo_id: save: %w", err)
		}
	}
	return nil
}

// DeleteByPath removes a single repository file entry.
func (r *RepositoryFileRepository) DeleteByPath(ctx context.Context, repoID uuid.UUID, path string) error {
	file, err := r.GetByRepoAndPath(ctx, repoID, path)
	if err != nil || file == nil {
		return err
	}
	now := time.Now()
	file.DeletedAt = &now
	return r.db.WithContext(ctx).Save(file).Error
}

// LanguageStat holds the file count for a single language.
type LanguageStat struct {
	Language  string
	FileCount int
}

// ProcessingCounts holds total repository file count and how many files were summarized.
type ProcessingCounts struct {
	TotalFiles      int
	SummarizedFiles int
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

// GetProcessingCounts returns total files and summarized files for a repository.
func (r *RepositoryFileRepository) GetProcessingCounts(ctx context.Context, repoID uuid.UUID) (ProcessingCounts, error) {
	type row struct {
		TotalFiles      int
		SummarizedFiles int
	}
	var counts row
	if err := r.db.WithContext(ctx).
		Model(&models.RepositoryFile{}).
		Select(
			"COUNT(*) AS total_files, COALESCE(SUM(CASE WHEN processed_at IS NOT NULL THEN 1 ELSE 0 END), 0) AS summarized_files",
		).
		Where("repo_id = ? AND deleted_at IS NULL", repoID).
		Scan(&counts).Error; err != nil {
		return ProcessingCounts{}, fmt.Errorf("repository_file_repository: get_processing_counts: %w", err)
	}

	return ProcessingCounts{
		TotalFiles:      counts.TotalFiles,
		SummarizedFiles: counts.SummarizedFiles,
	}, nil
}

// UpdateServiceMetadata stores the raw GitHub Contents API metadata for a file.
// It is a non-critical write — callers should log and continue on error rather than failing the job.
func (r *RepositoryFileRepository) UpdateServiceMetadata(ctx context.Context, repoID uuid.UUID, path string, metadata models.JSONB) error {
	file, err := r.GetByRepoAndPath(ctx, repoID, path)
	if err != nil || file == nil {
		return err
	}
	file.ServiceMetadata = metadata
	file.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(file).Error
}

// UpdateMetadata updates language, tags, and processed_at for a given file after chunk processing.
func (r *RepositoryFileRepository) UpdateMetadata(ctx context.Context, repoID uuid.UUID, path string, language string, tags models.StringArray, processedAt time.Time) error {
	file, err := r.GetByRepoAndPath(ctx, repoID, path)
	if err != nil || file == nil {
		return err
	}
	file.Language = language
	file.Tags = tags
	file.ProcessedAt = &processedAt
	file.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(file).Error
}

// ResetProcessingState clears derived metadata when a file has changed and must
// be fully re-processed.
func (r *RepositoryFileRepository) ResetProcessingState(ctx context.Context, repoID uuid.UUID, path string) error {
	file, err := r.GetByRepoAndPath(ctx, repoID, path)
	if err != nil || file == nil {
		return err
	}
	file.Tags = models.StringArray{}
	file.ProcessedAt = nil
	file.ServiceMetadata = nil
	file.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(file).Error
}

func (r *RepositoryFileRepository) PermanentlyDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.RepositoryFile{}, "id = ?", id).Error
}
