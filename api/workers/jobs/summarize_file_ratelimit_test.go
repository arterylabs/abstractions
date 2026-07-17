package jobs

import (
	ai_client "abstractions/clients/ai"
	"abstractions/services/chunker"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// stubAIProvider is a minimal ai_client.AIProvider used to exercise summarizeFileChunks and
// generateEmbeddings in isolation, without the DB-backed job harness used by
// tests/summarize_file_job_test.go. Each method can be scripted to return a canned error so the
// rate-limit propagation path can be verified.
type stubAIProvider struct {
	summarizeCodeErr error
	summarizeFileErr error
	embeddingsErr    error
}

func (s *stubAIProvider) SummarizeCode(_ context.Context, _, _, _ string, _, _ int, _ string) (*ai_client.SummarizeResult, error) {
	if s.summarizeCodeErr != nil {
		return nil, s.summarizeCodeErr
	}
	return &ai_client.SummarizeResult{Summary: "does a thing"}, nil
}

func (s *stubAIProvider) SummarizeFile(_ context.Context, _, _ string, _ []string) (*ai_client.SummarizeResult, error) {
	if s.summarizeFileErr != nil {
		return nil, s.summarizeFileErr
	}
	return &ai_client.SummarizeResult{Summary: "file summary"}, nil
}

func (s *stubAIProvider) GenerateTopics(_ context.Context, _ *ai_client.TopicGenerationRequest) ([]ai_client.Topic, error) {
	return nil, nil
}

func (s *stubAIProvider) GenerateNewsletterEntry(_ context.Context, _ *ai_client.NewsletterEntryRequest) (*ai_client.NewsletterEntry, error) {
	return nil, nil
}

func (s *stubAIProvider) GenerateEmbedding(_ context.Context, _ string) ([]float32, error) {
	return []float32{0.1}, nil
}

func (s *stubAIProvider) GenerateEmbeddings(_ context.Context, texts []string) ([][]float32, error) {
	if s.embeddingsErr != nil {
		return nil, s.embeddingsErr
	}
	out := make([][]float32, len(texts))
	for i := range out {
		out[i] = []float32{0.1}
	}
	return out, nil
}

func oneChunk() []chunker.Chunk {
	return []chunker.Chunk{{Content: "func foo() {}", StartLine: 1, EndLine: 1, Type: "declaration"}}
}

// A rate limit hit on a single chunk must propagate as an error, not be swallowed like an
// ordinary per-chunk failure — otherwise the file would be marked processed with missing
// summaries and never retried (the bug this feature fixes).
func TestSummarizeFileChunks_RateLimitOnChunk_PropagatesError(t *testing.T) {
	provider := &stubAIProvider{
		summarizeCodeErr: &ai_client.RateLimitError{Provider: ai_client.ProviderAnthropic, RetryAfter: time.Second},
	}

	_, err := summarizeFileChunks(context.Background(), provider, "main.go", "go", oneChunk())

	assert.Error(t, err)
	assert.True(t, ai_client.IsRateLimitError(err))
}

// An ordinary (non-rate-limit) per-chunk error is still skipped so one bad chunk doesn't block
// the rest of the file — this is the pre-existing best-effort behaviour and must be unchanged.
func TestSummarizeFileChunks_OrdinaryChunkError_SkipsChunk(t *testing.T) {
	provider := &stubAIProvider{summarizeCodeErr: errors.New("boom")}

	results, err := summarizeFileChunks(context.Background(), provider, "main.go", "go", oneChunk())

	assert.NoError(t, err)
	assert.Empty(t, results.chunks)
}

func TestSummarizeFileChunks_RateLimitOnFileSummary_PropagatesError(t *testing.T) {
	provider := &stubAIProvider{
		summarizeFileErr: &ai_client.RateLimitError{Provider: ai_client.ProviderOpenAI, RetryAfter: time.Second},
	}

	_, err := summarizeFileChunks(context.Background(), provider, "main.go", "go", oneChunk())

	assert.Error(t, err)
	assert.True(t, ai_client.IsRateLimitError(err))
}

func TestGenerateEmbeddings_RateLimitError_Propagates(t *testing.T) {
	provider := &stubAIProvider{
		embeddingsErr: &ai_client.RateLimitError{Provider: ai_client.ProviderGrok, RetryAfter: time.Second},
	}

	_, err := generateEmbeddings(context.Background(), provider, []string{"summary text"})

	assert.Error(t, err)
	assert.True(t, ai_client.IsRateLimitError(err))
}

func TestGenerateEmbeddings_OrdinaryError_ReturnsNilNil(t *testing.T) {
	provider := &stubAIProvider{embeddingsErr: errors.New("boom")}

	vecs, err := generateEmbeddings(context.Background(), provider, []string{"summary text"})

	assert.NoError(t, err)
	assert.Nil(t, vecs)
}
