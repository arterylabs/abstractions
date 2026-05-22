package ai_client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockAnthropicServer returns an httptest.Server that responds with the given body and status code.
func mockAnthropicServer(statusCode int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(body))
	}))
}

// anthropicSuccessBody builds a minimal Anthropic Messages API response body.
func anthropicSuccessBody(text string) string {
	resp := anthropicResponse{
		Content: []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		}{{Type: "text", Text: text}},
	}
	b, _ := json.Marshal(resp)
	return string(b)
}

// providerWithServer creates an AnthropicProvider that points at the given test server URL.
func providerWithServer(server *httptest.Server) *AnthropicProvider {
	client := server.Client()
	p := NewAnthropicProviderWithHTTPClient("test-key", "test-openai-key", client)
	// Override the Anthropic base URL by using the server URL directly.
	// We can't change the const, so we test via the complete method indirectly through public methods,
	// using a custom transport that rewrites the host.
	p.httpClient = &http.Client{
		Transport: &rewriteTransport{target: server.URL, inner: server.Client().Transport},
	}
	return p
}

// rewriteTransport rewrites all requests to point at the target server URL (ignoring the original host).
type rewriteTransport struct {
	target string
	inner  http.RoundTripper
}

func (t *rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Replace scheme+host with the test server URL.
	req = req.Clone(req.Context())
	req.URL.Scheme = "http"
	req.URL.Host = strings.TrimPrefix(t.target, "http://")
	req.URL.Path = "/" // simplify: ignore original path for the mock
	return t.inner.RoundTrip(req)
}

// --- Tests ---

func TestAnthropicProvider_SummarizeCode_ParsesJSON(t *testing.T) {
	responseText := `{"summary":"Handles user authentication via JWT.","domain_tags":["auth","jwt"]}`
	server := mockAnthropicServer(http.StatusOK, anthropicSuccessBody(responseText))
	defer server.Close()

	p := providerWithServer(server)
	result, err := p.SummarizeCode(context.Background(), "api/auth.go", "Go", "func ValidateToken(...) {}", 1, 1, "function")
	require.NoError(t, err)
	assert.Equal(t, "Handles user authentication via JWT.", result.Summary)
	assert.Equal(t, []string{"auth", "jwt"}, result.DomainTags)
}

func TestAnthropicProvider_SummarizeCode_FallsBackOnBadJSON(t *testing.T) {
	server := mockAnthropicServer(http.StatusOK, anthropicSuccessBody("This code validates tokens."))
	defer server.Close()

	p := providerWithServer(server)
	result, err := p.SummarizeCode(context.Background(), "api/auth.go", "Go", "func ValidateToken(...) {}", 1, 1, "function")
	require.NoError(t, err)
	// Fallback: raw text becomes the summary.
	assert.Contains(t, result.Summary, "This code validates tokens.")
}

func TestAnthropicProvider_SummarizeCode_ReturnsErrorOnAPIError(t *testing.T) {
	errBody := `{"error":{"type":"authentication_error","message":"Invalid API key"}}`
	server := mockAnthropicServer(http.StatusOK, errBody)
	defer server.Close()

	p := providerWithServer(server)
	_, err := p.SummarizeCode(context.Background(), "api/auth.go", "Go", "func ValidateToken(...) {}", 1, 1, "function")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "authentication_error")
}

func TestAnthropicProvider_GenerateTopics_ParsesJSON(t *testing.T) {
	topicsText := `[{"title":"JWT Auth","description":"New token validation.","related_files":["api/auth.go"]}]`
	server := mockAnthropicServer(http.StatusOK, anthropicSuccessBody(topicsText))
	defer server.Close()

	p := providerWithServer(server)
	req := &TopicGenerationRequest{
		RepoFullName:  "octocat/myrepo",
		FileSummaries: map[string]string{"api/auth.go": "Handles JWT validation."},
	}
	topics, err := p.GenerateTopics(context.Background(), req)
	require.NoError(t, err)
	require.Len(t, topics, 1)
	assert.Equal(t, "JWT Auth", topics[0].Title)
	assert.Equal(t, []string{"api/auth.go"}, topics[0].RelatedFiles)
}

func TestAnthropicProvider_GenerateNewsletterEntry_ReturnsHTMLContent(t *testing.T) {
	htmlContent := "<h2>Auth Refactor</h2><p>We improved JWT handling.</p>"
	server := mockAnthropicServer(http.StatusOK, anthropicSuccessBody(htmlContent))
	defer server.Close()

	p := providerWithServer(server)
	req := &NewsletterEntryRequest{
		RepoFullName:  "octocat/myrepo",
		Topic:         Topic{Title: "Auth Refactor", Description: "JWT improvements"},
		FileSummaries: map[string]string{"api/auth.go": "Validates tokens"},
	}
	entry, err := p.GenerateNewsletterEntry(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "Auth Refactor", entry.Title)
	assert.Equal(t, htmlContent, entry.Content)
}

func TestAnthropicProvider_GenerateEmbedding_ReturnsCorrectDimension(t *testing.T) {
	server := mockAnthropicServer(http.StatusOK, `{"data":[{"embedding":[0.1,0.2,0.3]}]}`)
	defer server.Close()

	p := providerWithServer(server)
	vec, err := p.GenerateEmbedding(context.Background(), "some text")
	require.NoError(t, err)
	assert.Len(t, vec, 3)
}

func TestBuildSummarizeUserPrompt_ContainsFilePathAndLanguage(t *testing.T) {
	prompt := buildSummarizeUserPrompt("cmd/main.go", "Go", "package main", 1, 1, "file")
	assert.Contains(t, prompt, "cmd/main.go")
	assert.Contains(t, prompt, "Go")
}

func TestBuildSummarizeUserPrompt_ContainsChunkContent(t *testing.T) {
	prompt := buildSummarizeUserPrompt("file.go", "Go", "func hello() {}", 3, 5, "function")
	assert.Contains(t, prompt, "func hello() {}")
}

func TestNormalizeCodeContent_StripsCommonIndent(t *testing.T) {
	code := "    func foo() {\n        return 1\n    }"
	result := normalizeCodeContent(code)
	assert.Equal(t, "func foo() {\n    return 1\n}", result)
}

func TestNormalizeCodeContent_PreservesRelativeIndent(t *testing.T) {
	// Python-style: remove only common prefix, preserve structure
	code := "    def foo():\n        x = 1\n        return x"
	result := normalizeCodeContent(code)
	assert.Equal(t, "def foo():\n    x = 1\n    return x", result)
}

func TestNormalizeCodeContent_CollapsesExtraBlankLines(t *testing.T) {
	code := "func a() {}\n\n\n\nfunc b() {}"
	result := normalizeCodeContent(code)
	assert.Equal(t, "func a() {}\n\n\nfunc b() {}", result)
}

func TestBuildTopicsPrompt_TruncatesLongSummaries(t *testing.T) {
	longSummary := strings.Repeat("x", maxTopicSummaryLength+100)
	req := &TopicGenerationRequest{
		RepoFullName:  "owner/repo",
		FileSummaries: map[string]string{"main.go": longSummary},
	}
	prompt := buildTopicsPrompt(req)
	assert.Contains(t, prompt, "...")
	assert.NotContains(t, prompt, longSummary)
}

func TestBuildTopicsPrompt_ContainsRepoAndSummaries(t *testing.T) {
	req := &TopicGenerationRequest{
		RepoFullName:  "owner/repo",
		FileSummaries: map[string]string{"main.go": "entry point"},
	}
	prompt := buildTopicsPrompt(req)
	assert.Contains(t, prompt, "owner/repo")
	assert.Contains(t, prompt, "main.go")
	assert.Contains(t, prompt, "entry point")
}

func TestFactory_ReturnsAnthropicProvider(t *testing.T) {
	provider := NewAIProvider()
	assert.NotNil(t, provider)
	_, ok := provider.(*AnthropicProvider)
	assert.True(t, ok, "NewAIProvider should return *AnthropicProvider")
}
