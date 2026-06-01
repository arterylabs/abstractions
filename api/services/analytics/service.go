package analytics_service

import (
	"abstractions/repositories"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrInvalidDateRange     = errors.New("invalid date range")
	ErrInvalidAnalyticsType = errors.New("invalid analytics type")
	ErrDataNotAvailable     = errors.New("analytics data not available")
)

type Service struct {
	db   *gorm.DB
	repo *repositories.AnalyticsRepository
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db, repo: repositories.NewAnalyticsRepository(db)}
}
