package workspace_configuration_service

import ai_client "abstractions/clients/ai"

// Model preference keys. Provider API-key config keys are not hardcoded here — they are sourced from
// the AI provider registry (ai_client.Providers) so adding a provider requires no change in this
// package.
const (
	KeyIndexingModel   = "indexing_model"
	KeyNewsletterModel = "newsletter_model"
	KeyEmbeddingsModel = "embeddings_model"

	// System prompt enhancements. Empty/unset means only the built-in default is used; non-empty
	// values are appended on top of the default in ai_client (never replace it).
	// These keys are intentionally not in envDefaults so Get returns "" when unset.
	KeySystemPromptSummarizeCode   = "system_prompt_summarize_code"
	KeySystemPromptSummarizeFile   = "system_prompt_summarize_file"
	KeySystemPromptTopicGenerator  = "system_prompt_topic_generator"
	KeySystemPromptNewsletterEntry = "system_prompt_newsletter_entry"
)

// EncryptedKeys is the set of configuration keys whose values must be encrypted at rest. Every
// provider's API key (from the registry) is encrypted.
var EncryptedKeys = func() map[string]bool {
	keys := make(map[string]bool, len(ai_client.Providers))
	for _, p := range ai_client.Providers {
		keys[p.ConfigKey] = true
	}
	return keys
}()
