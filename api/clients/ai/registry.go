package ai_client

import (
	"abstractions/config"
	"context"
	"net/http"
)

// textProvider is the transport-level interface for text generation. Each concrete provider
// (anthropicProvider, openAIProvider — the latter also serving OpenAI-compatible APIs like Grok)
// implements it so compositeProvider can route by model without knowing the concrete type.
type textProvider interface {
	complete(ctx context.Context, req completionRequest) (string, error)
}

// grokAPIBase is the OpenAI-compatible base URL for xAI's Grok API.
const grokAPIBase = "https://api.x.ai/v1"

// ProviderDescriptor is the single source of truth for everything provider-specific: how to read
// its API key from the environment and workspace configuration, how to validate that key, and how
// to build its transport. Adding a new provider is a matter of appending one descriptor here (plus
// its model rows in models.go) — no other layer hardcodes provider identity.
type ProviderDescriptor struct {
	// ID is the stable provider identifier, matching the Provider field on ModelInfo.
	ID Provider
	// DisplayName is the human-readable label shown in the UI and error messages.
	DisplayName string
	// ConfigKey is the workspace-configuration key that stores this provider's API key.
	ConfigKey string
	// KeyPlaceholder is the example key format shown in the settings input (e.g. "sk-ant-…").
	KeyPlaceholder string
	// IconPath is the frontend asset path for the provider's icon.
	IconPath string
	// SupportsEmbeddings indicates whether the provider can generate vector embeddings.
	SupportsEmbeddings bool
	// EnvKey returns the system-level API key for this provider, honouring DisableDefaultAIUsage.
	// Returns "" when no system default is available.
	EnvKey func() string
	// Validate checks whether an API key is accepted by the provider. Pass nil httpClient to use
	// the default client.
	Validate func(apiKey string, httpClient *http.Client) error
	// NewTextProvider builds the text-generation transport for this provider. Nil when the provider
	// cannot generate text.
	NewTextProvider func(apiKey string, httpClient *http.Client) textProvider
	// RequestsPerMin and TokensPerMin are conservative default rate limits used to pace outbound
	// calls to this provider (see ratelimit.go: limiterFor). They are deliberately below the
	// published limits for a typical low tier — if a workspace has a higher tier, requests simply
	// queue briefly rather than erroring. Zero disables throttling on that dimension.
	RequestsPerMin int
	TokensPerMin   int
	// MaxRetries is how many times doWithRetry will wait-and-retry a rate-limited (429/529)
	// response before giving up and returning a *RateLimitError.
	MaxRetries int
}

// envKeyFor returns a closure that reads the given field from config.Environment, honouring
// DisableDefaultAIUsage and a nil environment (as in tests).
func envKeyFor(field func(*config.ENV) string) func() string {
	return func() string {
		if config.Environment == nil || config.Environment.DisableDefaultAIUsage {
			return ""
		}
		return field(config.Environment)
	}
}

// Providers is the canonical registry of AI providers. Order controls display order in the UI.
var Providers = []ProviderDescriptor{
	{
		ID:                 ProviderAnthropic,
		DisplayName:        "Anthropic",
		ConfigKey:          "anthropic_api_key",
		KeyPlaceholder:     "sk-ant-…",
		IconPath:           "/assets/img/claude.svg",
		SupportsEmbeddings: false,
		EnvKey:             envKeyFor(func(e *config.ENV) string { return e.AnthropicAPIKey }),
		Validate:           ValidateAnthropicAPIKey,
		NewTextProvider: func(apiKey string, httpClient *http.Client) textProvider {
			return newAnthropicProvider(apiKey, httpClient)
		},
		RequestsPerMin: 50,
		TokensPerMin:   40_000,
		MaxRetries:     5,
	},
	{
		ID:                 ProviderOpenAI,
		DisplayName:        "OpenAI",
		ConfigKey:          "openai_api_key",
		KeyPlaceholder:     "sk-…",
		IconPath:           "/assets/img/openai.svg",
		SupportsEmbeddings: true,
		EnvKey:             envKeyFor(func(e *config.ENV) string { return e.OpenAIAPIKey }),
		Validate:           ValidateOpenAIAPIKey,
		NewTextProvider: func(apiKey string, httpClient *http.Client) textProvider {
			return newOpenAIProvider(apiKey, httpClient)
		},
		RequestsPerMin: 500,
		TokensPerMin:   200_000,
		MaxRetries:     5,
	},
	{
		ID:                 ProviderGrok,
		DisplayName:        "Grok (xAI)",
		ConfigKey:          "grok_api_key",
		KeyPlaceholder:     "xai-…",
		IconPath:           "/assets/img/grok.svg",
		SupportsEmbeddings: false,
		EnvKey:             envKeyFor(func(e *config.ENV) string { return e.XaiAPIKey }),
		Validate:           ValidateGrokAPIKey,
		NewTextProvider: func(apiKey string, httpClient *http.Client) textProvider {
			return newOpenAICompatibleProvider(grokAPIBase, apiKey, httpClient)
		},
		// xAI does not send Retry-After (see ratelimit.go: parseOpenAIRateLimit falls back to
		// exponential backoff for this provider), so lean a bit more on proactive throttling.
		RequestsPerMin: 60,
		TokensPerMin:   60_000,
		MaxRetries:     5,
	},
}

// ProviderByID looks up a provider descriptor by its ID. Returns the zero value and false if not
// found.
func ProviderByID(id Provider) (ProviderDescriptor, bool) {
	for _, p := range Providers {
		if p.ID == id {
			return p, true
		}
	}
	return ProviderDescriptor{}, false
}
