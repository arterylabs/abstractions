package repositories

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type CreationCountRow struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

type CreationMetricSource struct {
	Table           string
	AdditionalWhere string
}

type AnalyticsRepository struct {
	DB *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) *AnalyticsRepository {
	return &AnalyticsRepository{DB: db}
}

func (r *AnalyticsRepository) GetCreationCounts(ctx context.Context, truncFormat string, source CreationMetricSource, fromDate, toDate time.Time) ([]CreationCountRow, error) {
	whereClause := "created_at >= ? AND created_at < ?"
	if source.AdditionalWhere != "" {
		whereClause += " AND " + source.AdditionalWhere
	}

	query := fmt.Sprintf(`
		SELECT
			DATE_TRUNC('%s', created_at) as date,
			COUNT(*) as count
		FROM %s
		WHERE %s
		GROUP BY DATE_TRUNC('%s', created_at)
		ORDER BY date ASC
	`, truncFormat, source.Table, whereClause, truncFormat)

	var bucketedCounts []CreationCountRow
	err := r.DB.WithContext(ctx).Raw(query, fromDate, toDate).Scan(&bucketedCounts).Error
	return bucketedCounts, err
}
