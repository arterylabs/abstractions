package tests

import (
	analytics_service "abstractions/services/analytics"
	"abstractions/tests/testutils"
	"abstractions/utils/analytics_utils"
	"context"
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
	s.analyticsService = analytics_service.NewService(s.DB)
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

	response, err := s.analyticsService.GetAnalyticsData(context.Background(), params)

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

	response, err := s.analyticsService.GetAnalyticsData(context.Background(), params)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), response)
	assert.Equal(s.T(), string(analytics_utils.MetricWorkspacesCreated), response.Metric)
	// Should have exactly 2 workspaces created today
	assert.Equal(s.T(), 2, response.TotalCount)
	assert.Len(s.T(), response.DataPoints, 1) // 1 day of data points
}

func (s *AnalyticsServiceTestSuite) TestGetAnalyticsDataWithInvalidMetric() {
	params := analytics_utils.AnalyticsParams{
		FromDate:    analytics_utils.NewDateParam(time.Now().AddDate(0, 0, -1)),
		ToDate:      analytics_utils.NewDateParam(time.Now()),
		Metric:      "invalid_metric",
		Dimension:   analytics_utils.DimensionDate,
		Granularity: analytics_utils.GranularityDay,
	}

	response, err := s.analyticsService.GetAnalyticsData(context.Background(), params)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), response)
	assert.NotEmpty(s.T(), err.Error())
}
