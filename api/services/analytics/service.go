package analytics_service

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrInvalidDateRange     = errors.New("invalid date range")
	ErrInvalidAnalyticsType = errors.New("invalid analytics type")
	ErrDataNotAvailable     = errors.New("analytics data not available")
)

type Service struct {
	DB *gorm.DB
}

func NewAnalyticsService(db *gorm.DB) *Service {
	return &Service{DB: db}
}
