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
	// ContextWindow is the maximum number of tokens (input + output) the model supports.
	ContextWindow int
}

// Models is the canonical list of AI models available in the system.
// Keep this list up-to-date as providers release or retire models.
var Models = []ModelInfo{
	// --- Anthropic — Claude 4 ---
	{
		ID:                     "claude-opus-4-7",
		DisplayName:            "Claude Opus 4.7",
		Provider:               ProviderAnthropic,
		SupportsTextGeneration: true,
		SupportsEmbeddings:     false,
		PricePerMillionTokens:  15.00,
		ContextWindow:          200_000,
	},
	{
		ID:                     "claude-sonnet-4-6",
		DisplayName:            "Claude Sonnet 4.6",
		Provider:               ProviderAnthropic,
		SupportsTextGeneration: true,
		SupportsEmbeddings:     false,
		PricePerMillionTokens:  3.00,
		ContextWindow:          200_000,
	},
	{
		ID:                     "claude-haiku-4-5-20251001",
		DisplayName:            "Claude Haiku 4.5",
		Provider:               ProviderAnthropic,
		SupportsTextGeneration: true,
		SupportsEmbeddings:     false,
		PricePerMillionTokens:  0.80,
		ContextWindow:          200_000,
	},

	// --- Anthropic — Claude 3.5 ---
	{
		ID:                     "claude-3-5-sonnet-20241022",
		DisplayName:            "Claude 3.5 Sonnet",
		Provider:               ProviderAnthropic,
		SupportsTextGeneration: true,
		SupportsEmbeddings:     false,
		PricePerMillionTokens:  3.00,
		ContextWindow:          200_000,
	},
	{
		ID:                     "claude-3-5-haiku-20241022",
		DisplayName:            "Claude 3.5 Haiku",
		Provider:               ProviderAnthropic,
		SupportsTextGeneration: true,
		SupportsEmbeddings:     false,
		PricePerMillionTokens:  0.80,
		ContextWindow:          200_000,
	},

	// --- Anthropic — Claude 3 ---
	{
		ID:                     "claude-3-opus-20240229",
		DisplayName:            "Claude 3 Opus",
		Provider:               ProviderAnthropic,
		SupportsTextGeneration: true,
		SupportsEmbeddings:     false,
		PricePerMillionTokens:  15.00,
		ContextWindow:          200_000,
	},
	{
		ID:                     "claude-3-haiku-20240307",
		DisplayName:            "Claude 3 Haiku",
		Provider:               ProviderAnthropic,
		SupportsTextGeneration: true,
		SupportsEmbeddings:     false,
		PricePerMillionTokens:  0.25,
		ContextWindow:          200_000,
	},

	// --- OpenAI — GPT-4.1 ---
	{
		ID:                     "gpt-4.1",
		DisplayName:            "GPT-4.1",
		Provider:               ProviderOpenAI,
		SupportsTextGeneration: true,
		SupportsEmbeddings:     false,
		PricePerMillionTokens:  2.00,
		ContextWindow:          1_047_576,
	},
	{
		ID:                     "gpt-4.1-mini",
		DisplayName:            "GPT-4.1 Mini",
		Provider:               ProviderOpenAI,
		SupportsTextGeneration: true,
		SupportsEmbeddings:     false,
		PricePerMillionTokens:  0.40,
		ContextWindow:          1_047_576,
	},
	{
		ID:                     "gpt-4.1-nano",
		DisplayName:            "GPT-4.1 Nano",
		Provider:               ProviderOpenAI,
		SupportsTextGeneration: true,
		SupportsEmbeddings:     false,
		PricePerMillionTokens:  0.10,
		ContextWindow:          1_047_576,
	},

	// --- OpenAI — GPT-4o ---
	{
		ID:                     "gpt-4o",
		DisplayName:            "GPT-4o",
		Provider:               ProviderOpenAI,
		SupportsTextGeneration: true,
		SupportsEmbeddings:     false,
		PricePerMillionTokens:  2.50,
		ContextWindow:          128_000,
	},
	{
		ID:                     "gpt-4o-mini",
		DisplayName:            "GPT-4o Mini",
		Provider:               ProviderOpenAI,
		SupportsTextGeneration: true,
		SupportsEmbeddings:     false,
		PricePerMillionTokens:  0.15,
		ContextWindow:          128_000,
	},

	// --- OpenAI — Reasoning ---
	{
		ID:                     "o1",
		DisplayName:            "o1",
		Provider:               ProviderOpenAI,
		SupportsTextGeneration: true,
		SupportsEmbeddings:     false,
		PricePerMillionTokens:  15.00,
		ContextWindow:          200_000,
	},
	{
		ID:                     "o3",
		DisplayName:            "o3",
		Provider:               ProviderOpenAI,
		SupportsTextGeneration: true,
		SupportsEmbeddings:     false,
		PricePerMillionTokens:  10.00,
		ContextWindow:          200_000,
	},
	{
		ID:                     "o3-mini",
		DisplayName:            "o3 Mini",
		Provider:               ProviderOpenAI,
		SupportsTextGeneration: true,
		SupportsEmbeddings:     false,
		PricePerMillionTokens:  1.10,
		ContextWindow:          200_000,
	},
	{
		ID:                     "o4-mini",
		DisplayName:            "o4 Mini",
		Provider:               ProviderOpenAI,
		SupportsTextGeneration: true,
		SupportsEmbeddings:     false,
		PricePerMillionTokens:  1.10,
		ContextWindow:          200_000,
	},

	// --- OpenAI — Legacy ---
	{
		ID:                     "gpt-4-turbo",
		DisplayName:            "GPT-4 Turbo",
		Provider:               ProviderOpenAI,
		SupportsTextGeneration: true,
		SupportsEmbeddings:     false,
		PricePerMillionTokens:  10.00,
		ContextWindow:          128_000,
	},

	// --- OpenAI — embeddings ---
	{
		ID:                     "text-embedding-3-small",
		DisplayName:            "Text Embedding 3 Small",
		Provider:               ProviderOpenAI,
		SupportsTextGeneration: false,
		SupportsEmbeddings:     true,
		PricePerMillionTokens:  0.02,
		ContextWindow:          8_191,
	},
	{
		ID:                     "text-embedding-3-large",
		DisplayName:            "Text Embedding 3 Large",
		Provider:               ProviderOpenAI,
		SupportsTextGeneration: false,
		SupportsEmbeddings:     true,
		PricePerMillionTokens:  0.13,
		ContextWindow:          8_191,
	},
	{
		ID:                     "text-embedding-ada-002",
		DisplayName:            "Text Embedding Ada 002",
		Provider:               ProviderOpenAI,
		SupportsTextGeneration: false,
		SupportsEmbeddings:     true,
		PricePerMillionTokens:  0.10,
		ContextWindow:          8_191,
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
