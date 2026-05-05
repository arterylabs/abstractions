package models_test

import (
	"abstractions/db/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationConfiguration(t *testing.T) {
	integration := &models.Integration{
		ID:          uuid.New(),
		Type:        "slack",
		WorkspaceID: uuid.New(),
		Token:       "test-token",
		CreatedBy:   uuid.New(),
		UpdatedBy:   uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	t.Run("SetConfiguration with Slack config", func(t *testing.T) {
		slackConfig := &models.SlackConfiguration{
			AwayMessage: "I'm currently away, but I'll get back to you soon!",
		}

		err := integration.SetConfiguration(slackConfig)
		require.NoError(t, err)
		assert.NotNil(t, integration.Configuration)

		// Verify the configuration was stored correctly
		awayMessage, exists := integration.Configuration["away_message"]
		assert.True(t, exists)
		assert.Equal(t, "I'm currently away, but I'll get back to you soon!", awayMessage)
	})

	t.Run("GetConfiguration retrieves Slack config", func(t *testing.T) {
		var retrievedConfig models.SlackConfiguration
		err := integration.GetConfiguration(&retrievedConfig)
		require.NoError(t, err)

		assert.Equal(t, "I'm currently away, but I'll get back to you soon!", retrievedConfig.AwayMessage)
	})

	t.Run("SetConfiguration with empty away message", func(t *testing.T) {
		slackConfig := &models.SlackConfiguration{
			AwayMessage: "",
		}

		err := integration.SetConfiguration(slackConfig)
		require.NoError(t, err)
		assert.NotNil(t, integration.Configuration)

		var retrievedConfig models.SlackConfiguration
		err = integration.GetConfiguration(&retrievedConfig)
		require.NoError(t, err)
		assert.Empty(t, retrievedConfig.AwayMessage)
	})

	t.Run("SetConfiguration with nil clears config", func(t *testing.T) {
		err := integration.SetConfiguration(nil)
		require.NoError(t, err)
		assert.Nil(t, integration.Configuration)
	})

	t.Run("GetConfiguration with nil config does nothing", func(t *testing.T) {
		var retrievedConfig models.SlackConfiguration
		err := integration.GetConfiguration(&retrievedConfig)
		require.NoError(t, err)
		// Should remain zero values since configuration is nil
		assert.Empty(t, retrievedConfig.AwayMessage)
	})

	t.Run("Update Slack configuration", func(t *testing.T) {
		// First set a configuration
		slackConfig := &models.SlackConfiguration{
			AwayMessage: "Original message",
		}
		err := integration.SetConfiguration(slackConfig)
		require.NoError(t, err)

		// Update the configuration
		updatedSlackConfig := &models.SlackConfiguration{
			AwayMessage: "Updated away message",
		}
		err = integration.SetConfiguration(updatedSlackConfig)
		require.NoError(t, err)

		// Verify the update
		var retrievedConfig models.SlackConfiguration
		err = integration.GetConfiguration(&retrievedConfig)
		require.NoError(t, err)
		assert.Equal(t, "Updated away message", retrievedConfig.AwayMessage)
	})
}
