package ai_client

// Provider identifies which API a model belongs to.
type Provider string

const (
	ProviderAnthropic Provider = "anthropic"
	ProviderOpenAI    Provider = "openai"
)

// ModelInfo describes a single AI model available for use in the system.
type ModelInfo struct {
	// ID is the exact string passed to the API (e.g. "claude-sonnet-4-6").
	ID string
	// DisplayName is the human-readable label shown in the UI.
	DisplayName string
	// Provider indicates which vendor owns this model.
	Provider Provider
	// SupportsTextGeneration indicates the model can be used for summarisation,
	// topic generation, and newsletter writing.
	SupportsTextGeneration bool
	// SupportsEmbeddings indicates the model can produce vector embeddings.
	SupportsEmbeddings bool
	// PricePerMillionTokens is the input price in USD per one million tokens.
	PricePerMillionTokens float64
}

// Models is the canonical list of AI models available in the system.
// Keep this list up-to-date as providers release or retire models.
var Models = []ModelInfo{
	// --- Anthropic ---
	{
		ID:                     "claude-opus-4-7",
		DisplayName:            "Claude Opus 4.7",
		Provider:               ProviderAnthropic,
		SupportsTextGeneration: true,
		SupportsEmbeddings:     false,
		PricePerMillionTokens:  15.00,
	},
	{
		ID:                     "claude-sonnet-4-6",
		DisplayName:            "Claude Sonnet 4.6",
		Provider:               ProviderAnthropic,
		SupportsTextGeneration: true,
		SupportsEmbeddings:     false,
		PricePerMillionTokens:  3.00,
	},
	{
		ID:                     "claude-haiku-4-5-20251001",
		DisplayName:            "Claude Haiku 4.5",
		Provider:               ProviderAnthropic,
		SupportsTextGeneration: true,
		SupportsEmbeddings:     false,
		PricePerMillionTokens:  0.80,
	},

	// --- OpenAI — text generation ---
	{
		ID:                     "gpt-4o",
		DisplayName:            "GPT-4o",
		Provider:               ProviderOpenAI,
		SupportsTextGeneration: true,
		SupportsEmbeddings:     false,
		PricePerMillionTokens:  2.50,
	},
	{
		ID:                     "gpt-4o-mini",
		DisplayName:            "GPT-4o Mini",
		Provider:               ProviderOpenAI,
		SupportsTextGeneration: true,
		SupportsEmbeddings:     false,
		PricePerMillionTokens:  0.15,
	},
	{
		ID:                     "gpt-4-turbo",
		DisplayName:            "GPT-4 Turbo",
		Provider:               ProviderOpenAI,
		SupportsTextGeneration: true,
		SupportsEmbeddings:     false,
		PricePerMillionTokens:  10.00,
	},

	// --- OpenAI — embeddings ---
	{
		ID:                     "text-embedding-3-small",
		DisplayName:            "Text Embedding 3 Small",
		Provider:               ProviderOpenAI,
		SupportsTextGeneration: false,
		SupportsEmbeddings:     true,
		PricePerMillionTokens:  0.02,
	},
	{
		ID:                     "text-embedding-3-large",
		DisplayName:            "Text Embedding 3 Large",
		Provider:               ProviderOpenAI,
		SupportsTextGeneration: false,
		SupportsEmbeddings:     true,
		PricePerMillionTokens:  0.13,
	},
	{
		ID:                     "text-embedding-ada-002",
		DisplayName:            "Text Embedding Ada 002",
		Provider:               ProviderOpenAI,
		SupportsTextGeneration: false,
		SupportsEmbeddings:     true,
		PricePerMillionTokens:  0.10,
	},
}

// TextGenerationModels returns all models that support text generation.
func TextGenerationModels() []ModelInfo {
	var out []ModelInfo
	for _, m := range Models {
		if m.SupportsTextGeneration {
			out = append(out, m)
		}
	}
	return out
}

// EmbeddingModels returns all models that support generating embeddings.
func EmbeddingModels() []ModelInfo {
	var out []ModelInfo
	for _, m := range Models {
		if m.SupportsEmbeddings {
			out = append(out, m)
		}
	}
	return out
}

// ModelByID looks up a model by its API identifier. Returns the zero value and false if not found.
func ModelByID(id string) (ModelInfo, bool) {
	for _, m := range Models {
		if m.ID == id {
			return m, true
		}
	}
	return ModelInfo{}, false
}
