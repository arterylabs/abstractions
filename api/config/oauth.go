package config

import (
	"abstractions/services/integrations/enums"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuthConfig struct {
	oauth2.Config
}

func (c *OAuthConfig) GetAuthURL() string {
	return c.Endpoint.AuthURL
}

func (c *OAuthConfig) GetTokenURL() string {
	return c.Endpoint.TokenURL
}

type Providers struct {
	Google *OAuthConfig
	Slack  *OAuthConfig
}

var OAuthProviders *Providers

func LoadOAuthProviders() (*Providers, error) {
	OAuthProviders = &Providers{
		Google: &OAuthConfig{
			Config: oauth2.Config{
				ClientID:     Environment.GoogleClientID,
				ClientSecret: Environment.GoogleClientSecret,
				RedirectURL:  Environment.APIURL + "/auth/google/callback",
				Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
				Endpoint:     google.Endpoint,
			},
		},

		// Add other providers here
		Slack: &OAuthConfig{
			Config: oauth2.Config{
				ClientID:     Environment.SlackClientID,
				ClientSecret: Environment.SlackClientSecret,
				RedirectURL:  Environment.APIURL + "/integrations/callback/slack",
				Scopes:       []string{"channels:history", "channels:read", "chat:write", "chat:write.public", "commands", "files:read", "groups:history", "groups:read", "incoming-webhook", "links:read", "links:write", "reactions:read", "reactions:write", "users:read", "users:read.email", "team:read"},
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://slack.com/oauth/v2/authorize",
					TokenURL: "https://slack.com/api/oauth.v2.access",
				},
			},
		},
	}

	return OAuthProviders, nil
}

type OAuthState struct {
	Provider    string     `json:"provider"`
	Value       *string    `json:"value"`
	Timestamp   string     `json:"timestamp"`
	WorkspaceID *uuid.UUID `json:"workspace_id,omitempty"`
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	RedirectTo  string     `json:"redirect_to"`
}

func NewOAuthState(provider enums.IntegrationTypeEnum, redirectTo string, workspaceID, userID *uuid.UUID, value *string) *OAuthState {
	return &OAuthState{
		Provider:    provider.String(),
		Timestamp:   time.Now().Format(time.RFC3339),
		RedirectTo:  redirectTo,
		WorkspaceID: workspaceID,
		UserID:      userID,
		Value:       value,
	}
}

func (os *OAuthState) SetValue(value *string) {
	os.Value = value
}

func (os *OAuthState) ToString() string {
	br, err := json.Marshal(os)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(br)
}

func (os *OAuthState) FromString(str string) error {
	b, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, os)
}

func ParseOAuthState(str string) (*OAuthState, error) {
	os := &OAuthState{}
	err := os.FromString(str)
	if err != nil {
		return nil, err
	}
	return os, nil
}

type GoogleProfileResponse struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}
