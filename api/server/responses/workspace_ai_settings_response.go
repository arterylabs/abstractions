package responses

// SystemPromptsResponse holds workspace system prompt enhancements. Each field is the stored
// enhancement text appended on top of the built-in default, or null when none is set.
type SystemPromptsResponse struct {
	SummarizeCode   *string `json:"summarizeCode"`
	SummarizeFile   *string `json:"summarizeFile"`
	TopicGenerator  *string `json:"topicGenerator"`
	NewsletterEntry *string `json:"newsletterEntry"`
}

// DefaultSystemPromptsResponse holds the built-in system prompts (read-only reference).
// Workspace enhancements are layered on top of these and never replace them.
type DefaultSystemPromptsResponse struct {
	SummarizeCode   string `json:"summarizeCode"`
	SummarizeFile   string `json:"summarizeFile"`
	TopicGenerator  string `json:"topicGenerator"`
	NewsletterEntry string `json:"newsletterEntry"`
}

// WorkspaceAISettingsResponse returns a workspace's AI configuration with secret keys masked.
// APIKeyLastFour is keyed by provider ID; the value is the last 4 characters of the stored key, or
// nil when no key is configured for that provider.
type WorkspaceAISettingsResponse struct {
	APIKeyLastFour       map[string]*string           `json:"apiKeyLastFour"`
	IndexingModel        *string                      `json:"indexingModel"`
	NewsletterModel      *string                      `json:"newsletterModel"`
	EmbeddingsModel      *string                      `json:"embeddingsModel"`
	SystemPrompts        SystemPromptsResponse        `json:"systemPrompts"`
	DefaultSystemPrompts DefaultSystemPromptsResponse `json:"defaultSystemPrompts"`
}
