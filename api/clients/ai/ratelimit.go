package ai_client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimitError is returned by a transport's complete/embed call when a provider keeps
// responding with a rate-limit status (429, or 529 "overloaded" for Anthropic) past MaxRetries or
// past maxRetryWait. Callers that loop over many AI calls (e.g. workers/jobs/summarize_file.go)
// should check IsRateLimitError and propagate the error so the enclosing asynq job is retried
// later, instead of silently dropping the unit of work.
type RateLimitError struct {
	Provider   Provider
	RetryAfter time.Duration
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("ai: %s: rate limited, retry after %s", e.Provider, e.RetryAfter)
}

// IsRateLimitError reports whether err is, or wraps (via %w), a *RateLimitError.
func IsRateLimitError(err error) bool {
	var rle *RateLimitError
	return errors.As(err, &rle)
}

// rateLimitInfo is the outcome of inspecting a 429/529 response for a "how long to wait" signal.
type rateLimitInfo struct {
	retryAfter time.Duration
	ok         bool // true when a concrete wait duration was found in the response
}

// parseAnthropicRateLimit reads Anthropic's rate-limit signals from response headers, per
// https://platform.claude.com/docs/en/api/rate-limits (verified July 2026):
//   - retry-after: seconds to wait, present on 429s and aligned to the limit that was hit.
//   - anthropic-ratelimit-{requests,tokens}-reset: RFC3339 timestamps; fall back to the latest of
//     these minus now when retry-after is absent.
func parseAnthropicRateLimit(h http.Header) rateLimitInfo {
	if seconds := h.Get("retry-after"); seconds != "" {
		if d, ok := parseRetryAfterSeconds(seconds); ok {
			return rateLimitInfo{retryAfter: d, ok: true}
		}
	}

	var latest time.Time
	for _, key := range []string{"anthropic-ratelimit-requests-reset", "anthropic-ratelimit-tokens-reset"} {
		raw := h.Get(key)
		if raw == "" {
			continue
		}
		if ts, err := time.Parse(time.RFC3339, raw); err == nil && ts.After(latest) {
			latest = ts
		}
	}
	if latest.IsZero() {
		return rateLimitInfo{}
	}
	wait := time.Until(latest)
	if wait < 0 {
		wait = 0
	}
	return rateLimitInfo{retryAfter: wait, ok: true}
}

// parseOpenAIRateLimit reads OpenAI's (and xAI Grok's OpenAI-compatible) rate-limit signals,
// verified July 2026:
//   - retry-after-ms: milliseconds to wait, the most precise signal, present on 429s.
//   - retry-after: seconds, fallback.
//   - x-ratelimit-reset-{requests,tokens}: Go-duration-shaped strings (e.g. "6s", "1m30s").
//
// xAI Grok does not send Retry-After at all — in that case ok=false and the caller falls back to
// exponential backoff.
func parseOpenAIRateLimit(h http.Header) rateLimitInfo {
	if ms := h.Get("retry-after-ms"); ms != "" {
		if n, err := strconv.ParseInt(ms, 10, 64); err == nil && n >= 0 {
			return rateLimitInfo{retryAfter: time.Duration(n) * time.Millisecond, ok: true}
		}
	}
	if seconds := h.Get("retry-after"); seconds != "" {
		if d, ok := parseRetryAfterSeconds(seconds); ok {
			return rateLimitInfo{retryAfter: d, ok: true}
		}
	}

	var longest time.Duration
	found := false
	for _, key := range []string{"x-ratelimit-reset-requests", "x-ratelimit-reset-tokens"} {
		raw := h.Get(key)
		if raw == "" {
			continue
		}
		if d, err := time.ParseDuration(raw); err == nil {
			found = true
			if d > longest {
				longest = d
			}
		}
	}
	if !found {
		return rateLimitInfo{}
	}
	return rateLimitInfo{retryAfter: longest, ok: true}
}

// parseRetryAfterSeconds parses a Retry-After header value expressed as an integer number of
// seconds (both Anthropic and OpenAI use this form; the HTTP-date form is not used by either).
func parseRetryAfterSeconds(raw string) (time.Duration, bool) {
	n, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil || n < 0 {
		return 0, false
	}
	return time.Duration(n * float64(time.Second)), true
}

// retryConfig bounds how long doWithRetry will block on rate-limit responses before giving up
// and returning a *RateLimitError to the caller.
type retryConfig struct {
	maxRetries int
	maxWait    time.Duration // longest single wait we'll sleep for before giving up
}

const defaultMaxRetryWait = 60 * time.Second

// isRateLimitStatus reports whether an HTTP status code indicates the request should be retried
// after a wait: 429 (Too Many Requests) from any provider, or 529 (Anthropic-specific overloaded)
// and 503 (generic overloaded/unavailable, seen from OpenAI-compatible endpoints under load).
func isRateLimitStatus(status int) bool {
	return status == http.StatusTooManyRequests || status == 529 || status == http.StatusServiceUnavailable
}

// doWithRetry sends the request built by buildReq, applying the provider's proactive rate limiter
// before every attempt. On a rate-limit-shaped response it waits using the signal parse extracts
// (or exponential backoff with jitter when the provider gives no signal, e.g. Grok), then retries
// up to cfg.maxRetries times. estimatedTokens is used to draw down the provider's token bucket and
// should be a rough estimate (max output tokens + prompt length / 4 is fine — this only needs to
// keep us under the limit, not be exact).
//
// Returns the final HTTP status and body once a non-rate-limit response is received. If retries
// are exhausted, or a wait would exceed cfg.maxWait, it returns a *RateLimitError instead —
// callers that loop over many units of work (see workers/jobs/summarize_file.go) should propagate
// this so the enclosing job is retried later rather than silently dropping work.
func doWithRetry(
	ctx context.Context,
	client *http.Client,
	provider Provider,
	estimatedTokens int,
	buildReq func() (*http.Request, error),
	parse func(http.Header) rateLimitInfo,
	cfg retryConfig,
) (int, []byte, error) {
	if cfg.maxRetries <= 0 {
		cfg.maxRetries = 1
	}
	if cfg.maxWait <= 0 {
		cfg.maxWait = defaultMaxRetryWait
	}

	limiter := limiterFor(provider)

	var lastRetryAfter time.Duration
	for attempt := 0; attempt < cfg.maxRetries; attempt++ {
		if err := limiter.wait(ctx, estimatedTokens); err != nil {
			return 0, nil, fmt.Errorf("ai: %s: rate limiter: %w", provider, err)
		}

		req, err := buildReq()
		if err != nil {
			return 0, nil, err
		}

		resp, err := client.Do(req)
		if err != nil {
			return 0, nil, err
		}
		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return 0, nil, readErr
		}

		if !isRateLimitStatus(resp.StatusCode) {
			return resp.StatusCode, body, nil
		}

		info := parse(resp.Header)
		wait := info.retryAfter
		if !info.ok {
			wait = backoffWithJitter(attempt)
		}
		lastRetryAfter = wait

		if wait > cfg.maxWait || attempt == cfg.maxRetries-1 {
			return resp.StatusCode, body, &RateLimitError{Provider: provider, RetryAfter: wait}
		}

		timer := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			timer.Stop()
			return 0, nil, ctx.Err()
		case <-timer.C:
		}
	}

	return 0, nil, &RateLimitError{Provider: provider, RetryAfter: lastRetryAfter}
}

// backoffWithJitter is used when a provider gives no concrete wait signal (xAI Grok sends no
// Retry-After). Doubles from 1s per attempt, capped at defaultMaxRetryWait, with +/-25% jitter to
// avoid every concurrent worker retrying in lockstep.
func backoffWithJitter(attempt int) time.Duration {
	base := time.Second * time.Duration(math.Pow(2, float64(attempt)))
	if base > defaultMaxRetryWait {
		base = defaultMaxRetryWait
	}
	jitter := time.Duration(rand.Int63n(int64(base) / 2))
	return base/2 + jitter
}

// --- Proactive per-provider throttling ---

// providerLimiter paces outbound requests to a single provider so we stay under its published
// RPM/TPM limits and rarely hit a 429 in the first place. Requests and tokens are throttled
// independently since either can be the binding constraint.
type providerLimiter struct {
	requests *rate.Limiter
	tokens   *rate.Limiter
}

func newProviderLimiter(requestsPerMin, tokensPerMin int) *providerLimiter {
	pl := &providerLimiter{}
	if requestsPerMin > 0 {
		pl.requests = rate.NewLimiter(rate.Limit(float64(requestsPerMin)/60.0), max(1, requestsPerMin))
	}
	if tokensPerMin > 0 {
		pl.tokens = rate.NewLimiter(rate.Limit(float64(tokensPerMin)/60.0), max(1, tokensPerMin))
	}
	return pl
}

// wait blocks until both the request bucket and the token bucket (sized to estimatedTokens) have
// capacity, or ctx is done.
func (pl *providerLimiter) wait(ctx context.Context, estimatedTokens int) error {
	if pl.requests != nil {
		if err := pl.requests.WaitN(ctx, 1); err != nil {
			return err
		}
	}
	if pl.tokens != nil && estimatedTokens > 0 {
		// Burst size is capped to the bucket's max (tokensPerMin), so clamp requests that ask for
		// more than the bucket can ever hold — otherwise WaitN would return an error forever.
		burst := pl.tokens.Burst()
		n := estimatedTokens
		if burst > 0 && n > burst {
			n = burst
		}
		if err := pl.tokens.WaitN(ctx, n); err != nil {
			return err
		}
	}
	return nil
}

// limiters is the process-wide registry of per-provider throttles, shared across every
// compositeProvider/transport instance. This must be process-wide rather than per-transport
// because factory.NewAIProviderWithKeys builds a fresh transport for every job
// (see workers/jobs/*.go newWorkspaceAIProvider) — limiters living on the transport would reset
// per job and provide no real throttling.
var (
	limitersMu sync.Mutex
	limiters   = map[Provider]*providerLimiter{}
)

// limiterFor returns the shared limiter for provider, building it from the provider's registry
// descriptor on first use.
func limiterFor(provider Provider) *providerLimiter {
	limitersMu.Lock()
	defer limitersMu.Unlock()

	if pl, ok := limiters[provider]; ok {
		return pl
	}

	rpm, tpm := 0, 0
	if descriptor, found := ProviderByID(provider); found {
		rpm, tpm = descriptor.RequestsPerMin, descriptor.TokensPerMin
	}
	pl := newProviderLimiter(rpm, tpm)
	limiters[provider] = pl
	return pl
}

// estimateTokens gives a rough token estimate for throttling purposes: max output tokens plus a
// ~4-chars-per-token estimate of input length. It only needs to keep the token bucket roughly
// honest, not be exact — the provider's own response is the source of truth for real usage.
func estimateTokens(maxTokens int, promptChars int) int {
	return maxTokens + promptChars/4
}
