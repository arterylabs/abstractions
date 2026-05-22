package ai_client

import "abstractions/config"

// AIProviderParams holds the configuration for NewAIProviderWithKeys.
// Empty key fields fall back to system environment values; empty model fields fall back to built-in defaults.
type AIProviderParams struct {
	AnthropicKey    string
	OpenAIKey       string
	IndexingModel   string
	NewsletterModel string
	EmbeddingsModel string
}

// NewAIProvider returns the configured AI provider using system environment defaults.
// Anthropic is used for text generation; OpenAI is used for embeddings (set OPENAI_API_KEY).
func NewAIProvider() AIProvider {
	return NewAIProviderWithKeys(AIProviderParams{})
}

// NewAIProviderWithKeys returns a provider using the given keys and model names. Empty strings fall
// back to the system environment values and built-in model defaults respectively.
func NewAIProviderWithKeys(p AIProviderParams) AIProvider {
	if p.AnthropicKey == "" && config.Environment != nil && !config.Environment.DisableDefaultAIUsage {
		p.AnthropicKey = config.Environment.AnthropicAPIKey
	}
	if p.OpenAIKey == "" && config.Environment != nil && !config.Environment.DisableDefaultAIUsage {
		p.OpenAIKey = config.Environment.OpenAIAPIKey
	}
	return NewAnthropicProviderWithModels(AnthropicProviderParams{
		APIKey:          p.AnthropicKey,
		OpenAIKey:       p.OpenAIKey,
		IndexingModel:   p.IndexingModel,
		NewsletterModel: p.NewsletterModel,
		EmbeddingsModel: p.EmbeddingsModel,
	})
}
