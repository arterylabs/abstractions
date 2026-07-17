package responses

// WorkspaceProviderHealth describes the connection status of a single AI provider for a workspace.
// Required is true when at least one of the workspace's configured models belongs to this provider.
type WorkspaceProviderHealth struct {
	Configured bool   `json:"configured"`
	Valid      bool   `json:"valid"`
	Error      string `json:"error,omitempty"`
	Required   bool   `json:"required"`
}

// WorkspaceHealthResponse reports the health of every AI provider, keyed by provider ID.
type WorkspaceHealthResponse struct {
	Providers map[string]WorkspaceProviderHealth `json:"providers"`
}
