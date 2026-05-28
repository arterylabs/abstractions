package analytics_service

import (
	"abstractions/server/responses"
	"abstractions/utils/analytics_utils"
	"fmt"
	"time"
)

type creationCountRow struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

type creationMetricSource struct {
	table           string
	additionalWhere string
}

func (s *Service) getCreationCountData(params analytics_utils.AnalyticsParams, source creationMetricSource) ([]responses.AnalyticsDataPoint, int, error) {
	truncFormat := params.GetDateTruncFormat()

	whereClause := "created_at >= ? AND created_at < ?"
	if source.additionalWhere != "" {
		whereClause += " AND " + source.additionalWhere
	}

	query := fmt.Sprintf(`
		SELECT
			DATE_TRUNC('%s', created_at) as date,
			COUNT(*) as count
		FROM %s
		WHERE %s
		GROUP BY DATE_TRUNC('%s', created_at)
		ORDER BY date ASC
	`, truncFormat, source.table, whereClause, truncFormat)

	var bucketedCounts []creationCountRow
	err := s.DB.Raw(query, params.FromDate.Time, params.ToDateExclusive()).Scan(&bucketedCounts).Error
	if err != nil {
		return nil, 0, fmt.Errorf("analytics_service: get_creation_count_data: query: %w", err)
	}

	countByLabel := make(map[string]int, len(bucketedCounts))
	totalCount := 0
	for _, bucket := range bucketedCounts {
		label := params.FormatDateForResponse(bucket.Date)
		countByLabel[label] = bucket.Count
		totalCount += bucket.Count
	}

	dateRange := params.GenerateDateRange()
	dataPoints := make([]responses.AnalyticsDataPoint, 0, len(dateRange))
	for _, date := range dateRange {
		label := params.FormatDateForResponse(date)

		dataPoints = append(dataPoints, responses.AnalyticsDataPoint{
			Label: label,
			Value: date.Unix(),
			Count: countByLabel[label],
		})
	}

	return dataPoints, totalCount, nil
}

func (s *Service) getUsersCreatedData(params analytics_utils.AnalyticsParams) ([]responses.AnalyticsDataPoint, int, error) {
	return s.getCreationCountData(params, creationMetricSource{table: "users"})
}

func (s *Service) getRepositoriesCreatedData(params analytics_utils.AnalyticsParams) ([]responses.AnalyticsDataPoint, int, error) {
	return s.getCreationCountData(params, creationMetricSource{
		table:           "repositories",
		additionalWhere: "deleted_at IS NULL",
	})
}

func (s *Service) getNewslettersCreatedData(params analytics_utils.AnalyticsParams) ([]responses.AnalyticsDataPoint, int, error) {
	return s.getCreationCountData(params, creationMetricSource{table: "newsletters"})
}

func (s *Service) getWorkspacesCreatedData(params analytics_utils.AnalyticsParams) ([]responses.AnalyticsDataPoint, int, error) {
	return s.getCreationCountData(params, creationMetricSource{table: "workspaces"})
}
