package ai_client

import "context"

// AIProvider is the task-focused interface for all AI operations.
// Concrete implementations (Anthropic, OpenAI, ...) must implement all methods.
// Callers never import provider packages directly — they use factory.NewAIProvider().
type AIProvider interface {
	// SummarizeCode produces a concise plain-text summary, domain tags, and symbol metadata for a
	// single code chunk. filePath, language, startLine, endLine, and chunkType are used to build
	// context-aware prompts. chunkContent must never be logged by any implementation.
	SummarizeCode(ctx context.Context, filePath, language, chunkContent string, startLine, endLine int, chunkType string) (*SummarizeResult, error)

	// SummarizeFile produces a holistic plain-text summary for an entire file using the already-
	// generated per-chunk summaries. Raw file content is never passed to this method.
	SummarizeFile(ctx context.Context, filePath, language string, chunkSummaries []string) (*SummarizeResult, error)

	// GenerateTopics analyses a set of per-file summaries and proposes newsletter topics.
	GenerateTopics(ctx context.Context, req *TopicGenerationRequest) ([]Topic, error)

	// GenerateNewsletterEntry writes the HTML/Markdown content for one newsletter topic.
	GenerateNewsletterEntry(ctx context.Context, req *NewsletterEntryRequest) (*NewsletterEntry, error)

	// GenerateEmbedding returns a vector embedding for the given text.
	// The returned slice length is provider-dependent (typically 1536 for OpenAI-compatible).
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)

	// GenerateEmbeddings returns vector embeddings for a batch of texts in a single API call.
	// The returned slice has the same length and order as texts. Prefer this over calling
	// GenerateEmbedding in a loop to reduce round-trip overhead.
	GenerateEmbeddings(ctx context.Context, texts []string) ([][]float32, error)
}
