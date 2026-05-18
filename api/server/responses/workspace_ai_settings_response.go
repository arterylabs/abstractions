package responses

type WorkspaceAISettingsResponse struct {
	AnthropicAPIKeyLastFour *string `json:"anthropicApiKeyLastFour"`
	OpenAIAPIKeyLastFour    *string `json:"openAIApiKeyLastFour"`
	IndexingModel           *string `json:"indexingModel"`
	NewsletterModel         *string `json:"newsletterModel"`
	EmbeddingsModel         *string `json:"embeddingsModel"`
}
