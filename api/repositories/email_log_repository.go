package repositories

import (
	"abstractions/db/models"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type EmailLogRepository struct {
	DB *gorm.DB
}

func NewEmailLogRepository(db *gorm.DB) *EmailLogRepository {
	return &EmailLogRepository{DB: db}
}

func (r *EmailLogRepository) Create(ctx context.Context, log *models.EmailLog) error {
	if log.Email == "" {
		return gorm.ErrInvalidData
	}
	if log.Status == "" {
		return gorm.ErrInvalidData
	}
	return r.DB.WithContext(ctx).Create(log).Error
}

// CountSoftFailuresByEmail returns the number of transient (soft) send failures logged
// for an address, used to decide when the failure threshold has been crossed.
func (r *EmailLogRepository) CountSoftFailuresByEmail(ctx context.Context, email string) (int, error) {
	var count int64
	err := r.DB.WithContext(ctx).
		Model(&models.EmailLog{}).
		Where("email = ? AND status = ? AND error_type = ?", email, models.EmailLogStatusFailed, models.EmailErrorTypeSoft).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("email_log_repository: count_soft_failures_by_email: %w", err)
	}
	return int(count), nil
}
