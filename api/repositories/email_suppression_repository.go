package repositories

import (
	"abstractions/db/models"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type EmailSuppressionRepositoryQ interface {
	GetActiveByEmail(ctx context.Context, email string) (*models.EmailSuppression, error)
	Upsert(ctx context.Context, email string, reason string) error
	Delete(ctx context.Context, email string) error
}

type EmailSuppressionRepository struct {
	DB *gorm.DB
}

func NewEmailSuppressionRepository(db *gorm.DB) *EmailSuppressionRepository {
	return &EmailSuppressionRepository{DB: db}
}

// GetActiveByEmail returns the active (non-deleted) suppression for an address, or nil.
func (r *EmailSuppressionRepository) GetActiveByEmail(ctx context.Context, email string) (*models.EmailSuppression, error) {
	suppression := models.EmailSuppression{}
	if err := r.DB.WithContext(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&suppression).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("email_suppression_repository: get_active_by_email: %w", err)
	}
	return &suppression, nil
}

// Upsert suppresses an address. If an active suppression already exists it is a no-op.
func (r *EmailSuppressionRepository) Upsert(ctx context.Context, email string, reason string) error {
	existing, err := r.GetActiveByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("email_suppression_repository: upsert: %w", err)
	}
	if existing != nil {
		return nil
	}
	suppression := &models.EmailSuppression{
		Email:  email,
		Reason: reason,
	}
	if err := r.DB.WithContext(ctx).Create(suppression).Error; err != nil {
		return fmt.Errorf("email_suppression_repository: upsert: create: %w", err)
	}
	return nil
}

// Delete soft-deletes the active suppression for an address (manual unblock).
func (r *EmailSuppressionRepository) Delete(ctx context.Context, email string) error {
	now := time.Now()
	err := r.DB.WithContext(ctx).
		Model(&models.EmailSuppression{}).
		Where("email = ? AND deleted_at IS NULL", email).
		Updates(map[string]interface{}{"deleted_at": now, "updated_at": now}).Error
	if err != nil {
		return fmt.Errorf("email_suppression_repository: delete: %w", err)
	}
	return nil
}
