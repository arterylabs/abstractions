package ai_client

import "abstractions/config"

const (
	defaultIndexingModel   = "claude-haiku-4-5-20251001"
	defaultNewsletterModel = "claude-sonnet-4-6"
	defaultEmbeddingsModel = "text-embedding-3-small"

	// DefaultIndexingModel is the default model used for code summarisation when no model is
	// configured at the workspace level. Exported so callers can reason about required providers.
	DefaultIndexingModel = defaultIndexingModel
	// DefaultNewsletterModel is the default model used for topic generation and newsletter writing.
	DefaultNewsletterModel = defaultNewsletterModel
	// DefaultEmbeddingsModel is the default model used for generating vector embeddings.
	DefaultEmbeddingsModel = defaultEmbeddingsModel
)

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
func NewAIProvider() AIProvider {
	return NewAIProviderWithKeys(AIProviderParams{})
}

// NewAIProviderWithKeys returns a compositeProvider wired with the given keys and model names.
// Empty strings fall back to environment values and built-in model defaults.
//
// The composite routes each operation to the correct provider based on the model's Provider field:
// Anthropic models go to anthropicProvider, OpenAI models go to openAIProvider.
// To add a new provider, implement its HTTP transport in a new file and add a routing case
// to compositeProvider.complete().
func NewAIProviderWithKeys(p AIProviderParams) AIProvider {
	if p.AnthropicKey == "" && config.Environment != nil && !config.Environment.DisableDefaultAIUsage {
		p.AnthropicKey = config.Environment.AnthropicAPIKey
	}
	if p.OpenAIKey == "" && config.Environment != nil && !config.Environment.DisableDefaultAIUsage {
		p.OpenAIKey = config.Environment.OpenAIAPIKey
	}
	if p.IndexingModel == "" {
		p.IndexingModel = defaultIndexingModel
	}
	if p.NewsletterModel == "" {
		p.NewsletterModel = defaultNewsletterModel
	}
	if p.EmbeddingsModel == "" {
		p.EmbeddingsModel = defaultEmbeddingsModel
	}

	var openai *openAIProvider
	if p.OpenAIKey != "" {
		openai = newOpenAIProvider(p.OpenAIKey, nil)
	}

	return &compositeProvider{
		anthropic:       newAnthropicProvider(p.AnthropicKey, nil),
		openai:          openai,
		indexingModel:   p.IndexingModel,
		newsletterModel: p.NewsletterModel,
		embeddingsModel: p.EmbeddingsModel,
	}
}
