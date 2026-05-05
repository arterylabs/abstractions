package tests

import (
	analytics_service "abstractions/services/analytics"
	"abstractions/tests/testutils"
	"abstractions/utils/analytics_utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AnalyticsServiceTestSuite struct {
	testutils.TestSuite
	seeder           *testutils.Seeder
	analyticsService *analytics_service.Service
}

func TestAnalyticsService(t *testing.T) {
	suite.Run(t, new(AnalyticsServiceTestSuite))
}

func (s *AnalyticsServiceTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.analyticsService = analytics_service.NewAnalyticsService(s.DB)
}

func (s *AnalyticsServiceTestSuite) TestNewAnalyticsService() {
	service := analytics_service.NewAnalyticsService(s.DB)
	assert.NotNil(s.T(), service)
	assert.Equal(s.T(), s.DB, service.DB)
}

func (s *AnalyticsServiceTestSuite) TestGetAnalyticsDataUsersCreated() {
	// Setup test data - create users on the same day but different times
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())

	// Create test users
	user1 := s.seeder.CreateUser()
	user2 := s.seeder.CreateUser()
	user3 := s.seeder.CreateUser()

	// Update creation dates to be on the same day but different times
	s.DB.Model(&user1).Update("created_at", today.Add(-2*time.Hour))
	s.DB.Model(&user2).Update("created_at", today.Add(-1*time.Hour))
	s.DB.Model(&user3).Update("created_at", today)

	// Test getting users created analytics for today only
	fromDate := today.Truncate(24 * time.Hour)
	toDate := today.Truncate(24 * time.Hour).Add(24*time.Hour - time.Nanosecond) // End of day

	params := analytics_utils.AnalyticsParams{
		FromDate:    analytics_utils.NewDateParam(fromDate),
		ToDate:      analytics_utils.NewDateParam(toDate),
		Metric:      analytics_utils.MetricUsersCreated,
		Dimension:   analytics_utils.DimensionDate,
		Granularity: analytics_utils.GranularityDay,
	}

	response, err := s.analyticsService.GetAnalyticsData(params)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), response)
	assert.Equal(s.T(), string(analytics_utils.MetricUsersCreated), response.Metric)
	assert.Equal(s.T(), string(analytics_utils.DimensionDate), response.Dimension)
	assert.Equal(s.T(), string(analytics_utils.GranularityDay), response.Granularity)
	// Should have exactly 3 users created today
	assert.Equal(s.T(), 3, response.TotalCount)
	assert.Len(s.T(), response.DataPoints, 1) // 1 day of data points
}

func (s *AnalyticsServiceTestSuite) TestGetAnalyticsDataWorkspacesCreated() {
	// Setup test data - create workspaces on the same day
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())

	// Create test workspaces
	workspace1 := s.seeder.CreateWorkspace()
	workspace2 := s.seeder.CreateWorkspace()

	// Update creation dates to be on the same day
	s.DB.Model(&workspace1).Update("created_at", today.Add(-1*time.Hour))
	s.DB.Model(&workspace2).Update("created_at", today)

	// Test getting workspaces created analytics for today only
	fromDate := today.Truncate(24 * time.Hour)
	toDate := today.Truncate(24 * time.Hour).Add(24*time.Hour - time.Nanosecond) // End of day

	params := analytics_utils.AnalyticsParams{
		FromDate:    analytics_utils.NewDateParam(fromDate),
		ToDate:      analytics_utils.NewDateParam(toDate),
		Metric:      analytics_utils.MetricWorkspacesCreated,
		Dimension:   analytics_utils.DimensionDate,
		Granularity: analytics_utils.GranularityDay,
	}

	response, err := s.analyticsService.GetAnalyticsData(params)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), response)
	assert.Equal(s.T(), string(analytics_utils.MetricWorkspacesCreated), response.Metric)
	// Should have exactly 2 workspaces created today
	assert.Equal(s.T(), 2, response.TotalCount)
	assert.Len(s.T(), response.DataPoints, 1) // 1 day of data points
}

func (s *AnalyticsServiceTestSuite) TestGetAnalyticsDataWithEmptyResults() {
	// Test with no data in the specified date range
	now := time.Now()
	futureDate := now.AddDate(0, 0, 10)

	params := analytics_utils.AnalyticsParams{
		FromDate:    analytics_utils.NewDateParam(futureDate),
		ToDate:      analytics_utils.NewDateParam(futureDate.AddDate(0, 0, 1)),
		Metric:      analytics_utils.MetricUsersCreated,
		Dimension:   analytics_utils.DimensionDate,
		Granularity: analytics_utils.GranularityDay,
	}

	response, err := s.analyticsService.GetAnalyticsData(params)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), response)
	assert.Equal(s.T(), 0, response.TotalCount)
	assert.Len(s.T(), response.DataPoints, 2) // Still returns data points with zero counts

	// Verify all data points have zero count
	for _, dataPoint := range response.DataPoints {
		assert.Equal(s.T(), 0, dataPoint.Count)
	}
}

func (s *AnalyticsServiceTestSuite) TestGetAnalyticsDataWithInvalidMetric() {
	params := analytics_utils.AnalyticsParams{
		FromDate:    analytics_utils.NewDateParam(time.Now().AddDate(0, 0, -1)),
		ToDate:      analytics_utils.NewDateParam(time.Now()),
		Metric:      "invalid_metric",
		Dimension:   analytics_utils.DimensionDate,
		Granularity: analytics_utils.GranularityDay,
	}

	response, err := s.analyticsService.GetAnalyticsData(params)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), response)
	// The error should be related to the invalid metric - either validation or unsupported metric
	assert.NotEmpty(s.T(), err.Error())
}

func (s *AnalyticsServiceTestSuite) TestGetAnalyticsDataWithValidation() {
	// Test with invalid date range (ToDate before FromDate)
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	params := analytics_utils.AnalyticsParams{
		FromDate:    analytics_utils.NewDateParam(now),
		ToDate:      analytics_utils.NewDateParam(yesterday), // Invalid: ToDate before FromDate
		Metric:      analytics_utils.MetricUsersCreated,
		Dimension:   analytics_utils.DimensionDate,
		Granularity: analytics_utils.GranularityDay,
	}

	response, err := s.analyticsService.GetAnalyticsData(params)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), response)
}

func (s *AnalyticsServiceTestSuite) TestGetAnalyticsDataWithDifferentGranularities() {
	// Setup test data
	now := time.Now()
	user := s.seeder.CreateUser()
	s.DB.Model(&user).Update("created_at", now)

	testCases := []struct {
		granularity analytics_utils.Granularity
		dimension   analytics_utils.Dimension
	}{
		{analytics_utils.GranularityHour, analytics_utils.DimensionDate},
		{analytics_utils.GranularityDay, analytics_utils.DimensionDate},
		{analytics_utils.GranularityWeek, analytics_utils.DimensionDate},
		{analytics_utils.GranularityMonth, analytics_utils.DimensionDate},
	}

	for _, tc := range testCases {
		s.T().Run(string(tc.granularity), func(t *testing.T) {
			params := analytics_utils.AnalyticsParams{
				FromDate:    analytics_utils.NewDateParam(now.AddDate(0, 0, -1)),
				ToDate:      analytics_utils.NewDateParam(now),
				Metric:      analytics_utils.MetricUsersCreated,
				Dimension:   tc.dimension,
				Granularity: tc.granularity,
			}

			response, err := s.analyticsService.GetAnalyticsData(params)

			assert.NoError(t, err)
			assert.NotNil(t, response)
			assert.Equal(t, string(tc.granularity), response.Granularity)
			assert.Equal(t, string(tc.dimension), response.Dimension)
		})
	}
}

func (s *AnalyticsServiceTestSuite) TestGetAnalyticsDataDateFormatting() {
	// Create a user for testing
	now := time.Now()
	user := s.seeder.CreateUser()
	s.DB.Model(&user).Update("created_at", now)

	params := analytics_utils.AnalyticsParams{
		FromDate:    analytics_utils.NewDateParam(now.Truncate(24 * time.Hour)),
		ToDate:      analytics_utils.NewDateParam(now.Truncate(24 * time.Hour)),
		Metric:      analytics_utils.MetricUsersCreated,
		Dimension:   analytics_utils.DimensionDate,
		Granularity: analytics_utils.GranularityDay,
	}

	response, err := s.analyticsService.GetAnalyticsData(params)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), response)

	// Verify date formatting in response
	expectedStartDate := now.Truncate(24 * time.Hour).Format("2006-01-02")
	expectedEndDate := now.Truncate(24 * time.Hour).Format("2006-01-02")

	assert.Equal(s.T(), expectedStartDate, response.StartDate)
	assert.Equal(s.T(), expectedEndDate, response.EndDate)

	// Verify data point labels are properly formatted
	for _, dataPoint := range response.DataPoints {
		assert.NotEmpty(s.T(), dataPoint.Label)
		assert.NotNil(s.T(), dataPoint.Value)
	}
}

func (s *AnalyticsServiceTestSuite) TestGetAnalyticsDataDefaultParameters() {
	// Create a user for testing
	user := s.seeder.CreateUser()
	assert.NotNil(s.T(), user)

	// Test with minimal parameters - should use defaults
	params := analytics_utils.AnalyticsParams{
		Metric: analytics_utils.MetricUsersCreated,
	}

	response, err := s.analyticsService.GetAnalyticsData(params)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), response)
	assert.Equal(s.T(), string(analytics_utils.MetricUsersCreated), response.Metric)
	assert.Equal(s.T(), string(analytics_utils.DimensionDate), response.Dimension)
	assert.Equal(s.T(), string(analytics_utils.GranularityDay), response.Granularity)
	assert.NotEmpty(s.T(), response.StartDate)
	assert.NotEmpty(s.T(), response.EndDate)
}
