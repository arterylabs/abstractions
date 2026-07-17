package tests

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"abstractions/config"
	github_integration "abstractions/services/integrations/core/github"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestOAuthStateRoundTrip verifies that a legitimately-generated state string
// parses back into an equivalent state. This must keep passing after integrity
// protection is added.
func TestOAuthStateRoundTrip(t *testing.T) {
	config.Environment = &config.ENV{SigningKey: "test_signing_key_for_jwt_tokens_123456"}

	ws := uuid.New()
	uid := uuid.New()
	state := config.NewOAuthState(github_integration.IntegrationTypeGitHub, config.OAuthFlowIntegration, "/integrations", &ws, &uid, nil)

	parsed, err := config.ParseOAuthState(state.ToString())
	assert.NoError(t, err)
	assert.NotNil(t, parsed)
	assert.Equal(t, config.OAuthFlowIntegration, parsed.Flow)
	if assert.NotNil(t, parsed.WorkspaceID) {
		assert.Equal(t, ws, *parsed.WorkspaceID)
	}
	if assert.NotNil(t, parsed.UserID) {
		assert.Equal(t, uid, *parsed.UserID)
	}
}

// TestOAuthStateTamperedIsRejected verifies that an attacker cannot take a
// server-issued state, swap in a different workspace_id, and have it accepted.
// Because the state currently carries no signature, a tampered payload parses
// cleanly — so this test fails until integrity protection is added.
func TestOAuthStateTamperedIsRejected(t *testing.T) {
	config.Environment = &config.ENV{SigningKey: "test_signing_key_for_jwt_tokens_123456"}

	victimWS := uuid.New()
	attackerWS := uuid.New()
	uid := uuid.New()

	original := config.NewOAuthState(github_integration.IntegrationTypeGitHub, config.OAuthFlowIntegration, "/integrations", &victimWS, &uid, nil)
	encoded := original.ToString()

	// Attacker decodes the state, rewrites the workspace_id to a workspace they
	// control (or the victim's), and re-encodes it the naive way. The signature
	// (if any) is not recomputed with the server secret.
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		// A signed encoding may not be raw base64 JSON; craft the forged payload directly.
		raw = nil
	}

	var payload map[string]interface{}
	if raw != nil {
		_ = json.Unmarshal(raw, &payload)
	}
	if payload == nil {
		payload = map[string]interface{}{
			"provider": github_integration.IntegrationTypeGitHub.String(),
			"flow":     config.OAuthFlowIntegration,
			"user_id":  uid.String(),
		}
	}
	payload["workspace_id"] = attackerWS.String()

	forgedJSON, err := json.Marshal(payload)
	assert.NoError(t, err)
	forged := base64.StdEncoding.EncodeToString(forgedJSON)

	_, err = config.ParseOAuthState(forged)
	assert.Error(t, err, "a tampered/forged OAuth state must be rejected by ParseOAuthState")
}
