package ai_client

import (
	cache_client "abstractions/clients/cache"
	"abstractions/config"
	"log"
	"time"
)

const (
	defaultIndexingModel   = "claude-haiku-4-5-20251001"
	defaultNewsletterModel = "claude-sonnet-4-6"
	defaultEmbeddingsModel = "text-embedding-3-small"
	defaultPromptCacheTTL  = 7 * 24 * time.Hour

	// DefaultIndexingModel is the default model used for code summarisation when no model is
	// configured at the workspace level. Exported so callers can reason about required providers.
	DefaultIndexingModel = defaultIndexingModel
	// DefaultNewsletterModel is the default model used for topic generation and newsletter writing.
	DefaultNewsletterModel = defaultNewsletterModel
	// DefaultEmbeddingsModel is the default model used for generating vector embeddings.
	DefaultEmbeddingsModel = defaultEmbeddingsModel
)

// AIProviderParams holds the configuration for NewAIProviderWithKeys.
// Keys maps a Provider to its API key; a provider absent from the map (or with an empty value)
// falls back to the system environment value from its registry descriptor. Empty model fields fall
// back to built-in defaults.
// SystemPrompts are optional workspace enhancements for the four generation system prompts;
// empty fields use the built-in defaults only. Non-empty fields are appended on top of the
// defaults (never replace them) so response structure rules stay intact.
// Set Cache to inject a custom cache implementation (useful for tests or custom TTL requirements).
// If Cache is nil, the factory attempts to connect to Redis using environment configuration.
type AIProviderParams struct {
	Keys            map[Provider]string
	IndexingModel   string
	NewsletterModel string
	EmbeddingsModel string
	SystemPrompts   SystemPromptOverrides
	Cache           cache_client.Client
}

// NewAIProvider returns the configured AI provider using system environment defaults.
func NewAIProvider() AIProvider {
	return NewAIProviderWithKeys(AIProviderParams{})
}

// NewAIProviderWithKeys returns a compositeProvider wired from the provider registry using the
// given keys and model names. A provider without a key in p.Keys falls back to its registry
// EnvKey(); empty model fields fall back to built-in defaults.
//
// The composite routes each text operation to the provider registered for the model's Provider
// field via the textProviders map. To add a new provider, append a descriptor to Providers in
// registry.go and its model rows to models.go — this factory needs no changes.
func NewAIProviderWithKeys(p AIProviderParams) AIProvider {
	if p.IndexingModel == "" {
		p.IndexingModel = defaultIndexingModel
	}
	if p.NewsletterModel == "" {
		p.NewsletterModel = defaultNewsletterModel
	}
	if p.EmbeddingsModel == "" {
		p.EmbeddingsModel = defaultEmbeddingsModel
	}

	var cacheClient cache_client.Client
	if p.Cache != nil {
		cacheClient = p.Cache
	} else if config.Environment != nil && config.Environment.RedisHost != "" && config.Environment.RedisPort != 0 {
		redisClient, err := cache_client.NewRedisClientFromEnvironment()
		if err != nil {
			log.Printf("ai: prompt cache disabled: %v", err)
		} else {
			cacheClient = redisClient
		}
	}

	// Resolve the effective key for each provider (explicit override, else env fallback) and build
	// a text transport for every provider that has one available.
	textProviders := make(map[Provider]textProvider, len(Providers))
	var embeddings *openAIProvider
	for _, descriptor := range Providers {
		key := p.Keys[descriptor.ID]
		if key == "" && descriptor.EnvKey != nil {
			key = descriptor.EnvKey()
		}
		if key == "" {
			continue
		}
		if descriptor.NewTextProvider != nil {
			textProviders[descriptor.ID] = descriptor.NewTextProvider(key, nil)
		}
		if descriptor.SupportsEmbeddings && embeddings == nil {
			embeddings = newOpenAIProvider(key, nil)
		}
	}

	return &compositeProvider{
		textProviders:   textProviders,
		embeddings:      embeddings,
		indexingModel:   p.IndexingModel,
		newsletterModel: p.NewsletterModel,
		embeddingsModel: p.EmbeddingsModel,
		systemPrompts:   p.SystemPrompts,
		cache:           cacheClient,
		cacheTTL:        defaultPromptCacheTTL,
	}
}
