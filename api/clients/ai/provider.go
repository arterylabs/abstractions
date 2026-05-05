package ai_client

import "context"

// AIProvider is the task-focused interface for all AI operations.
// Concrete implementations (Anthropic, OpenAI, ...) must implement all methods.
// Callers never import provider packages directly — they use factory.NewAIProvider().
type AIProvider interface {
	// SummarizeCode produces a concise summary and domain tags for a single code chunk.
	// filePath and language are used to build context-aware prompts.
	// chunkContent must never be logged by any implementation.
	SummarizeCode(ctx context.Context, filePath, language, chunkContent string) (*SummarizeResult, error)

	// GenerateTopics analyses a set of per-file summaries and proposes newsletter topics.
	GenerateTopics(ctx context.Context, req *TopicGenerationRequest) ([]Topic, error)

	// GenerateNewsletterEntry writes the HTML/Markdown content for one newsletter topic.
	GenerateNewsletterEntry(ctx context.Context, req *NewsletterEntryRequest) (*NewsletterEntry, error)

	// GenerateEmbedding returns a vector embedding for the given text.
	// The returned slice length is provider-dependent (typically 1536 for OpenAI-compatible).
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
}
