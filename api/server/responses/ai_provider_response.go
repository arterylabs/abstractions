package responses

// AIProviderResponse describes a single AI provider so the frontend can render provider-specific
// UI (icons, key input placeholders, key rows) dynamically without hardcoding provider identity.
type AIProviderResponse struct {
	ID                 string `json:"id"`
	DisplayName        string `json:"displayName"`
	IconPath           string `json:"iconPath"`
	KeyPlaceholder     string `json:"keyPlaceholder"`
	SupportsEmbeddings bool   `json:"supportsEmbeddings"`
}
