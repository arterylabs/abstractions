package analytics_service

import (
	"abstractions/server/responses"
	"abstractions/utils/analytics_utils"
	"fmt"
	"time"
)

func (s *Service) getUsersCreatedData(params analytics_utils.AnalyticsParams) ([]responses.AnalyticsDataPoint, int, error) {
	type result struct {
		Date  time.Time `json:"date"`
		Count int       `json:"count"`
	}

	var results []result

	// Build the SQL query based on granularity
	truncFormat := params.GetDateTruncFormat()

	query := fmt.Sprintf(`
		SELECT
			DATE_TRUNC('%s', created_at) as date,
			COUNT(*) as count
		FROM users
		WHERE created_at >= ? AND created_at < ?
		GROUP BY DATE_TRUNC('%s', created_at)
		ORDER BY date ASC
	`, truncFormat, truncFormat)

	err := s.DB.Raw(query, params.FromDate.Time, params.ToDateExclusive()).Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	// Generate all dates in the range and fill gaps with zero
	allDates := params.GenerateDateRange()
	dataMap := make(map[string]int)

	// Fill the map with actual data
	totalCount := 0
	for _, r := range results {
		label := params.FormatDateForResponse(r.Date)
		dataMap[label] = r.Count
		totalCount += r.Count
	}

	// Create data points for all dates, including zeros
	var dataPoints []responses.AnalyticsDataPoint
	for _, date := range allDates {
		label := params.FormatDateForResponse(date)
		count := dataMap[label] // Will be 0 if not found

		dataPoints = append(dataPoints, responses.AnalyticsDataPoint{
			Label: label,
			Value: date.Unix(),
			Count: count,
		})
	}

	return dataPoints, totalCount, nil
}

func (s *Service) getRepositoriesCreatedData(params analytics_utils.AnalyticsParams) ([]responses.AnalyticsDataPoint, int, error) {
	type result struct {
		Date  time.Time `json:"date"`
		Count int       `json:"count"`
	}

	var results []result
	truncFormat := params.GetDateTruncFormat()

	query := fmt.Sprintf(`
		SELECT
			DATE_TRUNC('%s', created_at) as date,
			COUNT(*) as count
		FROM repositories
		WHERE created_at >= ? AND created_at < ? AND deleted_at IS NULL
		GROUP BY DATE_TRUNC('%s', created_at)
		ORDER BY date ASC
	`, truncFormat, truncFormat)

	err := s.DB.Raw(query, params.FromDate.Time, params.ToDateExclusive()).Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	allDates := params.GenerateDateRange()
	dataMap := make(map[string]int)
	totalCount := 0
	for _, r := range results {
		label := params.FormatDateForResponse(r.Date)
		dataMap[label] = r.Count
		totalCount += r.Count
	}

	var dataPoints []responses.AnalyticsDataPoint
	for _, date := range allDates {
		label := params.FormatDateForResponse(date)
		dataPoints = append(dataPoints, responses.AnalyticsDataPoint{
			Label: label,
			Value: date.Unix(),
			Count: dataMap[label],
		})
	}

	return dataPoints, totalCount, nil
}

func (s *Service) getNewslettersCreatedData(params analytics_utils.AnalyticsParams) ([]responses.AnalyticsDataPoint, int, error) {
	type result struct {
		Date  time.Time `json:"date"`
		Count int       `json:"count"`
	}

	var results []result
	truncFormat := params.GetDateTruncFormat()

	query := fmt.Sprintf(`
		SELECT
			DATE_TRUNC('%s', created_at) as date,
			COUNT(*) as count
		FROM newsletters
		WHERE created_at >= ? AND created_at < ?
		GROUP BY DATE_TRUNC('%s', created_at)
		ORDER BY date ASC
	`, truncFormat, truncFormat)

	err := s.DB.Raw(query, params.FromDate.Time, params.ToDateExclusive()).Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	allDates := params.GenerateDateRange()
	dataMap := make(map[string]int)
	totalCount := 0
	for _, r := range results {
		label := params.FormatDateForResponse(r.Date)
		dataMap[label] = r.Count
		totalCount += r.Count
	}

	var dataPoints []responses.AnalyticsDataPoint
	for _, date := range allDates {
		label := params.FormatDateForResponse(date)
		dataPoints = append(dataPoints, responses.AnalyticsDataPoint{
			Label: label,
			Value: date.Unix(),
			Count: dataMap[label],
		})
	}

	return dataPoints, totalCount, nil
}

func (s *Service) getWorkspacesCreatedData(params analytics_utils.AnalyticsParams) ([]responses.AnalyticsDataPoint, int, error) {
	type result struct {
		Date  time.Time `json:"date"`
		Count int       `json:"count"`
	}

	var results []result

	// Build the SQL query based on granularity
	truncFormat := params.GetDateTruncFormat()

	query := fmt.Sprintf(`
		SELECT
			DATE_TRUNC('%s', created_at) as date,
			COUNT(*) as count
		FROM workspaces
		WHERE created_at >= ? AND created_at < ?
		GROUP BY DATE_TRUNC('%s', created_at)
		ORDER BY date ASC
	`, truncFormat, truncFormat)

	err := s.DB.Raw(query, params.FromDate.Time, params.ToDateExclusive()).Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	// Generate all dates in the range and fill gaps with zero
	allDates := params.GenerateDateRange()
	dataMap := make(map[string]int)

	// Fill the map with actual data
	totalCount := 0
	for _, r := range results {
		label := params.FormatDateForResponse(r.Date)
		dataMap[label] = r.Count
		totalCount += r.Count
	}

	// Create data points for all dates, including zeros
	var dataPoints []responses.AnalyticsDataPoint
	for _, date := range allDates {
		label := params.FormatDateForResponse(date)
		count := dataMap[label] // Will be 0 if not found

		dataPoints = append(dataPoints, responses.AnalyticsDataPoint{
			Label: label,
			Value: date.Unix(),
			Count: count,
		})
	}

	return dataPoints, totalCount, nil
}
