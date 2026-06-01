package analytics_service

import (
	"abstractions/server/responses"
	"abstractions/utils/analytics_utils"
	"context"
	"fmt"
)

type analyticsDataFetcher func(context.Context, analytics_utils.AnalyticsParams) ([]responses.AnalyticsDataPoint, int, error)

// GetAnalyticsData retrieves analytics data based on the provided parameters
func (s *Service) GetAnalyticsData(ctx context.Context, params analytics_utils.AnalyticsParams) (*responses.AnalyticsDataResponse, error) {
	params.SetDefaults()

	if err := params.Validate(); err != nil {
		return nil, err
	}

	fetchers := map[analytics_utils.Metric]analyticsDataFetcher{
		analytics_utils.MetricUsersCreated:        s.getUsersCreatedData,
		analytics_utils.MetricWorkspacesCreated:   s.getWorkspacesCreatedData,
		analytics_utils.MetricRepositoriesCreated: s.getRepositoriesCreatedData,
		analytics_utils.MetricNewslettersCreated:  s.getNewslettersCreatedData,
	}

	fetchData, ok := fetchers[params.Metric]
	if !ok {
		return nil, fmt.Errorf("unsupported metric: %s", params.Metric)
	}

	dataPoints, totalCount, err := fetchData(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("analytics_service: get_analytics_data: %s: %w", params.Metric, err)
	}

	response := &responses.AnalyticsDataResponse{
		DataPoints:  dataPoints,
		Metric:      string(params.Metric),
		Dimension:   string(params.Dimension),
		Granularity: string(params.Granularity),
		StartDate:   params.FromDate.Format("2006-01-02"),
		EndDate:     params.ToDate.Format("2006-01-02"),
		TotalCount:  totalCount,
	}

	return response, nil
}
