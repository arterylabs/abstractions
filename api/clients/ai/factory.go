package ai_client

import "abstractions/config"

// NewAIProvider returns the configured AI provider.
// Anthropic is the default and current provider. To switch providers, change the
// implementation returned here — callers depend only on the AIProvider interface.
func NewAIProvider() AIProvider {
	apiKey := ""
	if config.Environment != nil {
		apiKey = config.Environment.AnthropicAPIKey
	}
	return NewAnthropicProvider(apiKey)
}
