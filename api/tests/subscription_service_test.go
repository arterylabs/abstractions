package tests

import (
	"abstractions/db/models"
	subscription_service "abstractions/services/subscription"
	"abstractions/tests/testutils"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SubscriptionTestSuite struct {
	testutils.TestSuite
	seeder              *testutils.Seeder
	subscriptionService *subscription_service.Service
}

func TestSubscriptionService(t *testing.T) {
	suite.Run(t, new(SubscriptionTestSuite))
}

func (s *SubscriptionTestSuite) SetupTest() {
	s.TestSuite.SetupTest()
	s.seeder = testutils.NewSeeder(s.DB)
	s.subscriptionService = subscription_service.NewService(s.DB)
}

func (s *SubscriptionTestSuite) TestCreateBulkSubscription() {
	workspace1 := s.seeder.CreateWorkspace(testutils.WithName("Workspace 1"))
	workspace2 := s.seeder.CreateWorkspace(testutils.WithName("Workspace 2"))

	workspaceIDs := []uuid.UUID{workspace1.ID, workspace2.ID}
	params := subscription_service.CreateSubscriptionParams{
		PlanID:    models.PlanIDPro,
		Status:    models.SubscriptionStatusActive,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		Limit:     1000,
	}

	// Test creating bulk subscriptions
	subs, err := s.subscriptionService.CreateBulkSubscription(context.Background(), workspaceIDs, params)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), subs)
	assert.Len(s.T(), *subs, 2)

	// Verify subscriptions were created properly
	for _, sub := range *subs {
		assert.Equal(s.T(), params.PlanID, sub.PlanID)
		assert.Equal(s.T(), params.Status, sub.Status)
		assert.Equal(s.T(), params.Limit, sub.Limit)
	}
}

func (s *SubscriptionTestSuite) TestGetActiveSubscription() {
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))

	// Create an active subscription
	params := subscription_service.CreateSubscriptionParams{
		PlanID:    models.PlanIDPro,
		Status:    models.SubscriptionStatusActive,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		Limit:     1000,
	}

	subs, err := s.subscriptionService.CreateBulkSubscription(context.Background(), []uuid.UUID{workspace.ID}, params)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), subs)

	// Test getting active subscription
	activeSub, err := s.subscriptionService.GetActiveSubscriptionByWorkspaceID(context.Background(), workspace.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), activeSub)
	assert.Equal(s.T(), models.SubscriptionStatusActive, activeSub.Status)
}

func (s *SubscriptionTestSuite) TestDisableSubscription() {
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))

	// Create a subscription
	params := subscription_service.CreateSubscriptionParams{
		PlanID:    models.PlanIDPro,
		Status:    models.SubscriptionStatusActive,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		Limit:     1000,
	}

	subs, err := s.subscriptionService.CreateBulkSubscription(context.Background(), []uuid.UUID{workspace.ID}, params)
	assert.NoError(s.T(), err)
	sub := (*subs)[0]

	// Test disabling subscription
	err = s.subscriptionService.DisableSubscription(context.Background(), sub.ID, "test reason")
	assert.NoError(s.T(), err)

	// Verify subscription is disabled
	subscription, err := s.subscriptionService.GetActiveSubscriptionByWorkspaceID(context.Background(), workspace.ID)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), subscription)

}

func (s *SubscriptionTestSuite) TestUpdateSubscription() {
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))

	// Create initial subscription
	params := subscription_service.CreateSubscriptionParams{
		PlanID:    models.PlanIDPro,
		Status:    models.SubscriptionStatusActive,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		Limit:     1000,
	}

	_, err := s.subscriptionService.CreateBulkSubscription(context.Background(), []uuid.UUID{workspace.ID}, params)
	assert.NoError(s.T(), err)

	// Update subscription
	newPlanID := models.PlanIDPro
	newStatus := models.SubscriptionStatusActive
	newLimit := 2000
	newExpiresAt := time.Now().Add(60 * 24 * time.Hour)

	updateParams := subscription_service.UpdateSubscriptionParams{
		PlanID:    &newPlanID,
		Status:    &newStatus,
		Limit:     &newLimit,
		ExpiresAt: &newExpiresAt,
	}

	err = s.subscriptionService.UpdateSubscriptionsForWorkspaces(context.Background(), []uuid.UUID{workspace.ID}, updateParams)
	assert.NoError(s.T(), err)

	// Verify update
	subscriptions, err := s.subscriptionService.GetSubscriptionsByWorkspaceID(context.Background(), workspace.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), newPlanID, subscriptions[0].PlanID)
	assert.Equal(s.T(), newLimit, subscriptions[0].Limit)
}

func (s *SubscriptionTestSuite) TestExpiredSubscription() {
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Test Workspace"))

	// Create expired subscription
	params := subscription_service.CreateSubscriptionParams{
		PlanID:    models.PlanIDPro,
		Status:    models.SubscriptionStatusActive,
		ExpiresAt: time.Now().Add(-24 * time.Hour), // Expired
		Limit:     1000,
	}

	subs, err := s.subscriptionService.CreateBulkSubscription(context.Background(), []uuid.UUID{workspace.ID}, params)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), subs)

	// Test getting active subscription should return nil for expired subscription
	activeSub, err := s.subscriptionService.GetActiveSubscriptionByWorkspaceID(context.Background(), workspace.ID)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), activeSub)
}

func (s *SubscriptionTestSuite) TestSubscriptionCounts() {
	workspace1 := s.seeder.CreateWorkspace(testutils.WithName("Workspace 1"))
	workspace2 := s.seeder.CreateWorkspace(testutils.WithName("Workspace 2"))

	// Create active and inactive subscriptions
	activeParams := subscription_service.CreateSubscriptionParams{
		PlanID:    models.PlanIDPro,
		Status:    models.SubscriptionStatusActive,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		Limit:     1000,
	}

	inactiveParams := subscription_service.CreateSubscriptionParams{
		PlanID:    models.PlanIDPro,
		Status:    models.SubscriptionStatusTrialing,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		Limit:     1000,
	}

	s.subscriptionService.CreateBulkSubscription(context.Background(), []uuid.UUID{workspace1.ID}, activeParams)
	s.subscriptionService.CreateBulkSubscription(context.Background(), []uuid.UUID{workspace2.ID}, inactiveParams)

	// Test counting active subscriptions
	activeCount, err := s.subscriptionService.CountActiveSubscriptions(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, activeCount)

	// Test counting all subscriptions
	totalCount, err := s.subscriptionService.Count(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 2, totalCount)
}

// TestEnableSubscription verifies that EnableSubscription re-activates a previously disabled
// subscription so that it becomes visible again as the active subscription for the workspace.
func (s *SubscriptionTestSuite) TestEnableSubscription() {
	workspace := s.seeder.CreateWorkspace(testutils.WithName("Enable Sub WS"))

	subs, err := s.subscriptionService.CreateBulkSubscription(context.Background(), []uuid.UUID{workspace.ID}, subscription_service.CreateSubscriptionParams{
		PlanID:    models.PlanIDPro,
		Status:    models.SubscriptionStatusActive,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		Limit:     500,
	})
	assert.NoError(s.T(), err)
	sub := (*subs)[0]

	// Disable it first.
	err = s.subscriptionService.DisableSubscription(context.Background(), sub.ID, "billing issue")
	assert.NoError(s.T(), err)

	// Confirm it is no longer active.
	active, err := s.subscriptionService.GetActiveSubscriptionByWorkspaceID(context.Background(), workspace.ID)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), active)

	// Re-enable.
	err = s.subscriptionService.EnableSubscription(context.Background(), sub.ID)
	assert.NoError(s.T(), err)

	// Should now be the active subscription again.
	active, err = s.subscriptionService.GetActiveSubscriptionByWorkspaceID(context.Background(), workspace.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), active)
	assert.Equal(s.T(), models.SubscriptionStatusActive, active.Status)
	assert.Nil(s.T(), active.DisabledAt)
	assert.Nil(s.T(), active.DisabledReason)
}

// TestDisableSubscriptions_Bulk verifies that DisableSubscriptions can disable multiple
// subscriptions in a single call.
func (s *SubscriptionTestSuite) TestDisableSubscriptions_Bulk() {
	ws1 := s.seeder.CreateWorkspace(testutils.WithName("Bulk Disable WS 1"))
	ws2 := s.seeder.CreateWorkspace(testutils.WithName("Bulk Disable WS 2"))

	subs1, err := s.subscriptionService.CreateBulkSubscription(context.Background(), []uuid.UUID{ws1.ID}, subscription_service.CreateSubscriptionParams{
		PlanID:    models.PlanIDPro,
		Status:    models.SubscriptionStatusActive,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		Limit:     100,
	})
	assert.NoError(s.T(), err)

	subs2, err := s.subscriptionService.CreateBulkSubscription(context.Background(), []uuid.UUID{ws2.ID}, subscription_service.CreateSubscriptionParams{
		PlanID:    models.PlanIDPro,
		Status:    models.SubscriptionStatusActive,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		Limit:     200,
	})
	assert.NoError(s.T(), err)

	ids := []uuid.UUID{(*subs1)[0].ID, (*subs2)[0].ID}
	err = s.subscriptionService.DisableSubscriptions(context.Background(), ids, "bulk disable test")
	assert.NoError(s.T(), err)

	// Both workspaces should no longer have an active subscription.
	for _, ws := range []*models.Workspace{ws1, ws2} {
		active, aerr := s.subscriptionService.GetActiveSubscriptionByWorkspaceID(context.Background(), ws.ID)
		assert.NoError(s.T(), aerr)
		assert.Nil(s.T(), active, "workspace %s should have no active subscription after bulk disable", ws.ID)
	}
}
