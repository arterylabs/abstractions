package models_test

import (
	"abstractions/db/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testConfig struct {
	Message string `json:"message"`
}

func TestIntegrationConfiguration(t *testing.T) {
	integration := &models.Integration{
		ID:          uuid.New(),
		Type:        "github",
		WorkspaceID: uuid.New(),
		Token:       "test-token",
		CreatedBy:   uuid.New(),
		UpdatedBy:   uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	t.Run("SetConfiguration stores config", func(t *testing.T) {
		cfg := &testConfig{Message: "hello"}

		err := integration.SetConfiguration(cfg)
		require.NoError(t, err)
		assert.NotNil(t, integration.Configuration)

		msg, exists := integration.Configuration["message"]
		assert.True(t, exists)
		assert.Equal(t, "hello", msg)
	})

	t.Run("GetConfiguration retrieves config", func(t *testing.T) {
		var retrieved testConfig
		err := integration.GetConfiguration(&retrieved)
		require.NoError(t, err)
		assert.Equal(t, "hello", retrieved.Message)
	})

	t.Run("SetConfiguration with empty value", func(t *testing.T) {
		cfg := &testConfig{Message: ""}

		err := integration.SetConfiguration(cfg)
		require.NoError(t, err)
		assert.NotNil(t, integration.Configuration)

		var retrieved testConfig
		err = integration.GetConfiguration(&retrieved)
		require.NoError(t, err)
		assert.Empty(t, retrieved.Message)
	})

	t.Run("SetConfiguration with nil clears config", func(t *testing.T) {
		err := integration.SetConfiguration(nil)
		require.NoError(t, err)
		assert.Nil(t, integration.Configuration)
	})

	t.Run("GetConfiguration with nil config does nothing", func(t *testing.T) {
		var retrieved testConfig
		err := integration.GetConfiguration(&retrieved)
		require.NoError(t, err)
		assert.Empty(t, retrieved.Message)
	})

	t.Run("Update configuration", func(t *testing.T) {
		err := integration.SetConfiguration(&testConfig{Message: "original"})
		require.NoError(t, err)

		err = integration.SetConfiguration(&testConfig{Message: "updated"})
		require.NoError(t, err)

		var retrieved testConfig
		err = integration.GetConfiguration(&retrieved)
		require.NoError(t, err)
		assert.Equal(t, "updated", retrieved.Message)
	})
}
