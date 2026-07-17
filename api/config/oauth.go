package config

import (
	"abstractions/services/integrations/enums"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// ErrInvalidOAuthState is returned when an OAuth state string is malformed or
// fails signature verification (i.e. it was tampered with or not issued by us).
var ErrInvalidOAuthState = errors.New("invalid oauth state")

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
	Google     *OAuthConfig
	GitHub     *OAuthConfig
	GitHubAuth *OAuthConfig
}

var OAuthProviders *Providers

func LoadOAuthProviders() (*Providers, error) {
	OAuthProviders = &Providers{
		Google: &OAuthConfig{
			Config: oauth2.Config{
				ClientID:     Environment.GoogleClientID,
				ClientSecret: Environment.GoogleClientSecret,
				RedirectURL:  Environment.APIURL + "/oauth/callback/google",
				Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
				Endpoint:     google.Endpoint,
			},
		},

		// GitHub is used for the workspace integration flow (repo scope only).
		GitHub: &OAuthConfig{
			Config: oauth2.Config{
				ClientID:     Environment.GitHubClientID,
				ClientSecret: Environment.GitHubClientSecret,
				RedirectURL:  Environment.APIURL + "/oauth/callback/github",
				Scopes:       []string{"repo"},
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://github.com/login/oauth/authorize",
					TokenURL: "https://github.com/login/oauth/access_token",
				},
			},
		},

		// GitHubAuth is used for the login/signup flow (includes user:email scope).
		GitHubAuth: &OAuthConfig{
			Config: oauth2.Config{
				ClientID:     Environment.GitHubClientID,
				ClientSecret: Environment.GitHubClientSecret,
				RedirectURL:  Environment.APIURL + "/oauth/callback/github",
				Scopes:       []string{"user:email", "repo"},
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://github.com/login/oauth/authorize",
					TokenURL: "https://github.com/login/oauth/access_token",
				},
			},
		},
	}

	return OAuthProviders, nil
}

// OAuthFlowAuth identifies a login/signup OAuth flow.
const OAuthFlowAuth = "auth"

// OAuthFlowIntegration identifies a workspace integration OAuth flow.
const OAuthFlowIntegration = "integration"

type OAuthState struct {
	Provider     string     `json:"provider"`
	Flow         string     `json:"flow"`
	Value        *string    `json:"value"`
	Timestamp    string     `json:"timestamp"`
	WorkspaceID  *uuid.UUID `json:"workspace_id,omitempty"`
	UserID       *uuid.UUID `json:"user_id,omitempty"`
	RedirectTo   string     `json:"redirect_to"`
	ResourceType *string    `json:"resource_type,omitempty"`
	ResourceID   *uuid.UUID `json:"resource_id,omitempty"`
	// SourcePage is the page that initiated the OAuth flow ("login" or "signup").
	// Used to redirect errors back to the correct page.
	SourcePage string `json:"source_page,omitempty"`
	// CreateRepository signals that after this integration is created the user
	// should be directed to the repository connection flow. We cannot use
	// ResourceID here because the repository does not exist yet.
	CreateRepository *bool `json:"create_repository,omitempty"`
}

func NewOAuthState(provider enums.IntegrationTypeEnum, flow string, redirectTo string, workspaceID, userID *uuid.UUID, value *string) *OAuthState {
	return &OAuthState{
		Provider:    provider.String(),
		Flow:        flow,
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

// oauthStateMAC computes the HMAC-SHA256 of the payload using the server signing
// key, so a state string cannot be forged or tampered with by a client.
func oauthStateMAC(payload []byte) []byte {
	mac := hmac.New(sha256.New, []byte(Environment.SigningKey))
	mac.Write(payload)
	return mac.Sum(nil)
}

// ToString serializes the state as "<base64(json)>.<base64(hmac)>". The signature
// binds the payload to the server secret; ParseOAuthState rejects anything whose
// signature does not verify.
func (os *OAuthState) ToString() string {
	br, err := json.Marshal(os)
	if err != nil {
		return ""
	}
	payload := base64.RawURLEncoding.EncodeToString(br)
	sig := base64.RawURLEncoding.EncodeToString(oauthStateMAC(br))
	return payload + "." + sig
}

func (os *OAuthState) FromString(str string) error {
	payloadPart, sigPart, ok := strings.Cut(str, ".")
	if !ok {
		return ErrInvalidOAuthState
	}

	payload, err := base64.RawURLEncoding.DecodeString(payloadPart)
	if err != nil {
		return ErrInvalidOAuthState
	}

	sig, err := base64.RawURLEncoding.DecodeString(sigPart)
	if err != nil {
		return ErrInvalidOAuthState
	}

	// Constant-time comparison so we don't leak signature bytes via timing.
	if !hmac.Equal(sig, oauthStateMAC(payload)) {
		return ErrInvalidOAuthState
	}

	return json.Unmarshal(payload, os)
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

type GitHubProfileResponse struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
}

type GitHubEmailResponse struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}
