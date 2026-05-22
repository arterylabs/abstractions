package responses

type WorkspaceConnectionStatus struct {
	Configured bool   `json:"configured"`
	Valid      bool   `json:"valid"`
	Error      string `json:"error,omitempty"`
}

type WorkspaceHealthResponse struct {
	Anthropic         WorkspaceConnectionStatus `json:"anthropic"`
	OpenAI            WorkspaceConnectionStatus `json:"openAI"`
	AnthropicRequired bool                      `json:"anthropicRequired"`
	OpenAIRequired    bool                      `json:"openAIRequired"`
}
