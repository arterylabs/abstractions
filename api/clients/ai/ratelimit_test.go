package ai_client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- header parsing ---

func TestParseAnthropicRateLimit_RetryAfterHeader_TakesPriority(t *testing.T) {
	h := http.Header{}
	h.Set("retry-after", "2")
	h.Set("anthropic-ratelimit-requests-reset", time.Now().Add(time.Hour).Format(time.RFC3339))

	info := parseAnthropicRateLimit(h)

	require.True(t, info.ok)
	assert.Equal(t, 2*time.Second, info.retryAfter)
}

func TestParseAnthropicRateLimit_FallsBackToResetTimestamp(t *testing.T) {
	h := http.Header{}
	reset := time.Now().Add(5 * time.Second)
	h.Set("anthropic-ratelimit-tokens-reset", reset.Format(time.RFC3339))

	info := parseAnthropicRateLimit(h)

	require.True(t, info.ok)
	// Allow slack for the RFC3339 second-granularity round trip and test execution time.
	assert.InDelta(t, 5*time.Second, info.retryAfter, float64(2*time.Second))
}

func TestParseAnthropicRateLimit_NoHeaders_NotOk(t *testing.T) {
	info := parseAnthropicRateLimit(http.Header{})
	assert.False(t, info.ok)
}

func TestParseOpenAIRateLimit_RetryAfterMs_TakesPriority(t *testing.T) {
	h := http.Header{}
	h.Set("retry-after-ms", "1500")
	h.Set("retry-after", "9")

	info := parseOpenAIRateLimit(h)

	require.True(t, info.ok)
	assert.Equal(t, 1500*time.Millisecond, info.retryAfter)
}

func TestParseOpenAIRateLimit_FallsBackToResetDuration(t *testing.T) {
	h := http.Header{}
	h.Set("x-ratelimit-reset-requests", "6s")
	h.Set("x-ratelimit-reset-tokens", "1m30s")

	info := parseOpenAIRateLimit(h)

	require.True(t, info.ok)
	// The longer of the two reset windows wins, since that's the binding constraint.
	assert.Equal(t, 90*time.Second, info.retryAfter)
}

// xAI Grok is OpenAI-compatible but sends no Retry-After header at all (verified July 2026) — the
// caller must fall back to exponential backoff in this case.
func TestParseOpenAIRateLimit_NoHeaders_NotOk(t *testing.T) {
	info := parseOpenAIRateLimit(http.Header{})
	assert.False(t, info.ok)
}

// --- RateLimitError plumbing ---

func TestIsRateLimitError_UnwrapsWrappedError(t *testing.T) {
	base := &RateLimitError{Provider: ProviderAnthropic, RetryAfter: time.Second}
	wrapped := fmt.Errorf("summarize_file: %w", fmt.Errorf("SummarizeCode: %w", base))

	assert.True(t, IsRateLimitError(wrapped))
}

func TestIsRateLimitError_OrdinaryError_ReturnsFalse(t *testing.T) {
	assert.False(t, IsRateLimitError(errors.New("boom")))
}

// --- end-to-end retry behaviour through the real transports ---

// countingRateLimitServer returns an httptest.Server that responds with statusOnFirstN for the
// first n requests (attaching headers via setHeaders) and successBody with a 200 afterward.
func countingRateLimitServer(t *testing.T, n int, statusOnFirstN int, setHeaders func(http.Header), successBody string) (*httptest.Server, *int32) {
	t.Helper()
	var mu sync.Mutex
	count := 0
	requestCount := new(int32)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		count++
		current := count
		mu.Unlock()
		*requestCount = int32(current)

		w.Header().Set("Content-Type", "application/json")
		if current <= n {
			if setHeaders != nil {
				setHeaders(w.Header())
			}
			w.WriteHeader(statusOnFirstN)
			_, _ = w.Write([]byte(`{"error":{"type":"rate_limit_error","message":"rate limited"}}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(successBody))
	}))
	return server, requestCount
}

func TestAnthropicComplete_RetryAfterHeader_SucceedsOnRetry(t *testing.T) {
	server, requestCount := countingRateLimitServer(t, 1, http.StatusTooManyRequests, func(h http.Header) {
		h.Set("retry-after", "0")
	}, anthropicSuccessBody("hello"))
	defer server.Close()

	transport := &rewriteTransport{target: server.URL, inner: server.Client().Transport}
	provider := newAnthropicProvider("test-key", &http.Client{Transport: transport})

	result, err := provider.complete(context.Background(), completionRequest{model: "claude-sonnet-5", userPrompt: "hi", maxTokens: 10})

	require.NoError(t, err)
	assert.Equal(t, "hello", result)
	assert.EqualValues(t, 2, *requestCount)
}

func TestAnthropicComplete_PersistentRateLimit_ReturnsRateLimitError(t *testing.T) {
	server, requestCount := countingRateLimitServer(t, 1000, http.StatusTooManyRequests, func(h http.Header) {
		h.Set("retry-after", "0")
	}, anthropicSuccessBody("unused"))
	defer server.Close()

	transport := &rewriteTransport{target: server.URL, inner: server.Client().Transport}
	provider := newAnthropicProvider("test-key", &http.Client{Transport: transport})

	_, err := provider.complete(context.Background(), completionRequest{model: "claude-sonnet-5", userPrompt: "hi", maxTokens: 10})

	require.Error(t, err)
	assert.True(t, IsRateLimitError(err))
	// Retries are capped by the Anthropic registry descriptor's MaxRetries.
	descriptor, found := ProviderByID(ProviderAnthropic)
	require.True(t, found)
	assert.EqualValues(t, descriptor.MaxRetries, *requestCount)
}

func TestOpenAIComplete_RetryAfterMs_SucceedsOnRetry(t *testing.T) {
	successBody := `{"choices":[{"message":{"content":"hi there"}}]}`
	server, requestCount := countingRateLimitServer(t, 1, http.StatusTooManyRequests, func(h http.Header) {
		h.Set("retry-after-ms", "0")
	}, successBody)
	defer server.Close()

	transport := &rewriteTransport{target: server.URL, inner: server.Client().Transport}
	provider := newOpenAIProvider("test-key", &http.Client{Transport: transport})

	result, err := provider.complete(context.Background(), completionRequest{model: "gpt-4.1-mini", userPrompt: "hi", maxTokens: 10})

	require.NoError(t, err)
	assert.Equal(t, "hi there", result)
	assert.EqualValues(t, 2, *requestCount)
}

// xAI Grok sends no Retry-After of any kind on a 429 — the transport must fall back to
// exponential backoff (see backoffWithJitter) rather than failing immediately or waiting forever.
func TestGrokComplete_NoRetryAfterHeader_FallsBackToBackoffAndSucceeds(t *testing.T) {
	successBody := `{"choices":[{"message":{"content":"grok says hi"}}]}`
	server, requestCount := countingRateLimitServer(t, 1, http.StatusTooManyRequests, nil, successBody)
	defer server.Close()

	transport := &rewriteTransport{target: server.URL, inner: server.Client().Transport}
	provider := newOpenAICompatibleProvider(grokAPIBase, "test-key", &http.Client{Transport: transport})

	result, err := provider.complete(context.Background(), completionRequest{model: "grok-4", userPrompt: "hi", maxTokens: 10})

	require.NoError(t, err)
	assert.Equal(t, "grok says hi", result)
	assert.EqualValues(t, 2, *requestCount)
}

func TestProviderForBaseURL_Grok(t *testing.T) {
	assert.Equal(t, ProviderGrok, providerForBaseURL(grokAPIBase))
	assert.Equal(t, ProviderOpenAI, providerForBaseURL(openAIAPIBase))
}

// --- proactive throttle ---

func TestNewProviderLimiter_ZeroLimits_DisablesThrottling(t *testing.T) {
	pl := newProviderLimiter(0, 0)
	err := pl.wait(context.Background(), 1_000_000)
	assert.NoError(t, err)
}

func TestProviderLimiter_TokenBucket_ClampsToBurst(t *testing.T) {
	pl := newProviderLimiter(0, 100) // burst == 100 tokens
	// Requesting far more than the bucket can ever hold must not block forever — it should be
	// clamped to the burst size and succeed immediately (bucket starts full).
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := pl.wait(ctx, 10_000)
	assert.NoError(t, err)
}

func TestEstimateTokens_AddsMaxTokensAndPromptEstimate(t *testing.T) {
	got := estimateTokens(100, 400)
	assert.Equal(t, 200, got) // 100 + 400/4
}

func TestParseRetryAfterSeconds_InvalidValue_NotOk(t *testing.T) {
	_, ok := parseRetryAfterSeconds("not-a-number")
	assert.False(t, ok)
}

func TestParseRetryAfterSeconds_ValidValue(t *testing.T) {
	d, ok := parseRetryAfterSeconds(strconv.Itoa(3))
	require.True(t, ok)
	assert.Equal(t, 3*time.Second, d)
}
