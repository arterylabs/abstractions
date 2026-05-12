package ai_client

import "abstractions/config"

// NewAIProvider returns the configured AI provider.
// Anthropic is used for text generation; OpenAI is used for embeddings (set OPENAI_API_KEY).
func NewAIProvider() AIProvider {
	apiKey := ""
	openAIKey := ""
	if config.Environment != nil {
		apiKey = config.Environment.AnthropicAPIKey
		openAIKey = config.Environment.OpenAIAPIKey
	}
	return NewAnthropicProvider(apiKey, openAIKey)
}
