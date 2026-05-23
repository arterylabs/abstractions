package analytics_service

import (
	"abstractions/server/responses"
	"abstractions/utils/analytics_utils"
	"fmt"
)

// GetAnalyticsData retrieves analytics data based on the provided parameters
func (s *Service) GetAnalyticsData(params analytics_utils.AnalyticsParams) (*responses.AnalyticsDataResponse, error) {
	params.SetDefaults()

	if err := params.Validate(); err != nil {
		return nil, err
	}

	var dataPoints []responses.AnalyticsDataPoint
	var totalCount int

	switch params.Metric {
	case analytics_utils.MetricUsersCreated:
		data, total, err := s.getUsersCreatedData(params)
		if err != nil {
			return nil, fmt.Errorf("analytics_service: get_analytics_data: get_users_created: %w", err)
		}
		dataPoints = data
		totalCount = total

	case analytics_utils.MetricWorkspacesCreated:
		data, total, err := s.getWorkspacesCreatedData(params)
		if err != nil {
			return nil, fmt.Errorf("analytics_service: get_analytics_data: get_workspaces_created: %w", err)
		}
		dataPoints = data
		totalCount = total

	case analytics_utils.MetricRepositoriesCreated:
		data, total, err := s.getRepositoriesCreatedData(params)
		if err != nil {
			return nil, fmt.Errorf("analytics_service: get_analytics_data: get_repositories_created: %w", err)
		}
		dataPoints = data
		totalCount = total

	case analytics_utils.MetricNewslettersCreated:
		data, total, err := s.getNewslettersCreatedData(params)
		if err != nil {
			return nil, fmt.Errorf("analytics_service: get_analytics_data: get_newsletters_created: %w", err)
		}
		dataPoints = data
		totalCount = total

	default:
		return nil, fmt.Errorf("unsupported metric: %s", params.Metric)
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
