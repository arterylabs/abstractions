package requests

type UpdateWorkspaceAISettings struct {
	AnthropicAPIKey *string `json:"anthropicApiKey"`
	OpenAIAPIKey    *string `json:"openAIApiKey"`
	IndexingModel   *string `json:"indexingModel"`
	NewsletterModel *string `json:"newsletterModel"`
	EmbeddingsModel *string `json:"embeddingsModel"`
}

func (r *UpdateWorkspaceAISettings) Validate() error {
	return nil
}
