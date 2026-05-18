package workspace_configuration_service

const (
	KeyAnthropicAPIKey = "anthropic_api_key"
	KeyOpenAIAPIKey    = "openai_api_key"
	KeyIndexingModel   = "indexing_model"
	KeyNewsletterModel = "newsletter_model"
	KeyEmbeddingsModel = "embeddings_model"
)

// EncryptedKeys is the set of configuration keys whose values must be encrypted at rest.
var EncryptedKeys = map[string]bool{
	KeyAnthropicAPIKey: true,
	KeyOpenAIAPIKey:    true,
}
