package analytics_service

import (
	"abstractions/repositories"
	"abstractions/server/responses"
	"abstractions/utils/analytics_utils"
	"context"
	"fmt"
)

func (s *Service) getCreationCountData(ctx context.Context, params analytics_utils.AnalyticsParams, source repositories.CreationMetricSource) ([]responses.AnalyticsDataPoint, int, error) {
	truncFormat := params.GetDateTruncFormat()

	bucketedCounts, err := s.repo.GetCreationCounts(ctx, truncFormat, source, params.FromDate.Time, params.ToDateExclusive())
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

func (s *Service) getUsersCreatedData(ctx context.Context, params analytics_utils.AnalyticsParams) ([]responses.AnalyticsDataPoint, int, error) {
	return s.getCreationCountData(ctx, params, repositories.CreationMetricSource{Table: "users"})
}

func (s *Service) getRepositoriesCreatedData(ctx context.Context, params analytics_utils.AnalyticsParams) ([]responses.AnalyticsDataPoint, int, error) {
	return s.getCreationCountData(ctx, params, repositories.CreationMetricSource{
		Table:           "repositories",
		AdditionalWhere: "deleted_at IS NULL",
	})
}

func (s *Service) getNewslettersCreatedData(ctx context.Context, params analytics_utils.AnalyticsParams) ([]responses.AnalyticsDataPoint, int, error) {
	return s.getCreationCountData(ctx, params, repositories.CreationMetricSource{Table: "newsletters"})
}

func (s *Service) getWorkspacesCreatedData(ctx context.Context, params analytics_utils.AnalyticsParams) ([]responses.AnalyticsDataPoint, int, error) {
	return s.getCreationCountData(ctx, params, repositories.CreationMetricSource{Table: "workspaces"})
}
