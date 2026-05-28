package newsletter

import (
	"abstractions/db/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// ---- ClassifyFile ----

func TestClassifyFile_CodeFiles_AreHighPriority(t *testing.T) {
	cases := []string{
		"cmd/main.go",
		"services/auth/service.go",
		"src/components/Button.tsx",
		"src/api/index.ts",
		"lib/utils.js",
		"lib/app.jsx",
		"app/models/user.py",
		"app/controllers/users_controller.rb",
		"src/main.rs",
		"internal/server.java",
		"pkg/handler.kt",
	}
	for _, path := range cases {
		t.Run(path, func(t *testing.T) {
			assert.Equal(t, FilePriorityHigh, ClassifyFile(path), "expected high priority for %s", path)
		})
	}
}

func TestClassifyFile_DocFiles_AreLowPriority(t *testing.T) {
	cases := []string{
		"README.md",
		"docs/guide.md",
		"CHANGELOG.mdx",
		"docs/setup.rst",
		"notes.txt",
		"CONTRIBUTING.md",
	}
	for _, path := range cases {
		t.Run(path, func(t *testing.T) {
			assert.Equal(t, FilePriorityLow, ClassifyFile(path), "expected low priority for %s", path)
		})
	}
}

func TestClassifyFile_TestFiles_AreLowPriority(t *testing.T) {
	cases := []string{
		"services/auth/service_test.go",
		"src/components/Button.test.tsx",
		"src/api/client.test.ts",
		"src/utils/format.spec.ts",
		"tests/test_models.py",
		"app/models/user_test.rb",
		"app/controllers/users_spec.rb",
		"e2e/login.spec.js",
	}
	for _, path := range cases {
		t.Run(path, func(t *testing.T) {
			assert.Equal(t, FilePriorityLow, ClassifyFile(path), "expected low priority for %s", path)
		})
	}
}

// ---- PartitionSummaries ----

func TestPartitionSummaries_SplitsCorrectly(t *testing.T) {
	summaries := map[string]string{
		"cmd/main.go":           "entry point",
		"README.md":             "project readme",
		"services/auth.go":      "auth service",
		"services/auth_test.go": "auth tests",
		"src/App.tsx":           "root component",
		"src/App.test.tsx":      "app component tests",
	}

	code, lowPri := PartitionSummaries(summaries)

	assert.Equal(t, 3, len(code), "expected 3 code files")
	assert.Contains(t, code, "cmd/main.go")
	assert.Contains(t, code, "services/auth.go")
	assert.Contains(t, code, "src/App.tsx")

	assert.Equal(t, 3, len(lowPri), "expected 3 low-priority files")
	assert.Contains(t, lowPri, "README.md")
	assert.Contains(t, lowPri, "services/auth_test.go")
	assert.Contains(t, lowPri, "src/App.test.tsx")
}

func TestPartitionSummaries_EmptyInput(t *testing.T) {
	code, lowPri := PartitionSummaries(map[string]string{})
	assert.Empty(t, code)
	assert.Empty(t, lowPri)
}

// ---- prioritisedFileSummaries ----

func TestPrioritisedFileSummaries_CodeFilesFirstWhenOverLimit(t *testing.T) {
	summaries := map[string]string{}
	// 3 code files
	summaries["a/main.go"] = "main"
	summaries["b/service.go"] = "service"
	summaries["c/handler.go"] = "handler"
	// 2 doc files
	summaries["README.md"] = "readme"
	summaries["docs/setup.md"] = "setup"

	// Limit to 3 — should keep the 3 code files and drop the docs
	result := prioritisedFileSummaries(summaries, 3)
	assert.Len(t, result, 3)
	assert.Contains(t, result, "a/main.go")
	assert.Contains(t, result, "b/service.go")
	assert.Contains(t, result, "c/handler.go")
	assert.NotContains(t, result, "README.md")
	assert.NotContains(t, result, "docs/setup.md")
}

func TestPrioritisedFileSummaries_LowPriFilledWhenSpaceRemains(t *testing.T) {
	summaries := map[string]string{
		"a/main.go":    "main",
		"README.md":    "readme",
		"CHANGELOG.md": "changelog",
	}
	// Limit of 3 — 1 code file leaves 2 slots for docs
	result := prioritisedFileSummaries(summaries, 3)
	assert.Len(t, result, 3)
	assert.Contains(t, result, "a/main.go")
	assert.Contains(t, result, "README.md")
	assert.Contains(t, result, "CHANGELOG.md")
}

func TestPrioritisedFileSummaries_NoCapWhenUnderLimit(t *testing.T) {
	summaries := map[string]string{
		"a/main.go": "main",
		"README.md": "readme",
	}
	result := prioritisedFileSummaries(summaries, 200)
	assert.Equal(t, summaries, result)
}

// ---- IsDocsTopic ----

func TestIsDocsTopic(t *testing.T) {
	assert.True(t, IsDocsTopic("Improving Documentation Coverage", "We look at the README and guides."))
	assert.True(t, IsDocsTopic("README overhaul", "The markdown files were restructured."))
	assert.True(t, IsDocsTopic("Writing a changelog", "A guide for contributors."))
	assert.False(t, IsDocsTopic("How the auth middleware works", "JWT validation and session handling."))
	assert.False(t, IsDocsTopic("Database migration strategy", "Schema evolution in production."))
}

// ---- IsTestingTopic ----

func TestIsTestingTopic(t *testing.T) {
	assert.True(t, IsTestingTopic("Testing strategy in the API", "Unit and integration tests."))
	assert.True(t, IsTestingTopic("End-to-end coverage with Playwright", "E2E test suite overview."))
	assert.True(t, IsTestingTopic("Jest configuration", "Frontend unit test setup."))
	assert.False(t, IsTestingTopic("How the job queue works", "Asynq-based background processing."))
	assert.False(t, IsTestingTopic("Auth token refresh flow", "JWT signing and rotation."))
}

// ---- FilterChunksByPriority ----

func makeChunk(path string) models.ChunkEmbedding {
	return models.ChunkEmbedding{
		ID:           uuid.New(),
		RepoFilePath: path,
	}
}

func TestFilterChunksByPriority_RemovesLowPriByDefault(t *testing.T) {
	chunks := []models.ChunkEmbedding{
		makeChunk("services/auth.go"),
		makeChunk("README.md"),
		makeChunk("services/auth_test.go"),
		makeChunk("src/App.tsx"),
	}

	filtered := FilterChunksByPriority(chunks, false, false)
	assert.Len(t, filtered, 2)
	assert.Equal(t, "services/auth.go", filtered[0].RepoFilePath)
	assert.Equal(t, "src/App.tsx", filtered[1].RepoFilePath)
}

func TestFilterChunksByPriority_IncludesMarkdownWhenRequested(t *testing.T) {
	chunks := []models.ChunkEmbedding{
		makeChunk("services/auth.go"),
		makeChunk("README.md"),
		makeChunk("services/auth_test.go"),
	}

	filtered := FilterChunksByPriority(chunks, true, false)
	paths := make([]string, len(filtered))
	for i, c := range filtered {
		paths[i] = c.RepoFilePath
	}
	assert.Contains(t, paths, "services/auth.go")
	assert.Contains(t, paths, "README.md")
	assert.NotContains(t, paths, "services/auth_test.go")
}

func TestFilterChunksByPriority_IncludesTestsWhenRequested(t *testing.T) {
	chunks := []models.ChunkEmbedding{
		makeChunk("services/auth.go"),
		makeChunk("README.md"),
		makeChunk("services/auth_test.go"),
	}

	filtered := FilterChunksByPriority(chunks, false, true)
	paths := make([]string, len(filtered))
	for i, c := range filtered {
		paths[i] = c.RepoFilePath
	}
	assert.Contains(t, paths, "services/auth.go")
	assert.NotContains(t, paths, "README.md")
	assert.Contains(t, paths, "services/auth_test.go")
}

func TestFilterChunksByPriority_IncludesAllWhenBothTrue(t *testing.T) {
	chunks := []models.ChunkEmbedding{
		makeChunk("services/auth.go"),
		makeChunk("README.md"),
		makeChunk("services/auth_test.go"),
	}

	filtered := FilterChunksByPriority(chunks, true, true)
	assert.Len(t, filtered, 3)
}

func TestFilterChunksByPriority_EmptyInput(t *testing.T) {
	filtered := FilterChunksByPriority(nil, false, false)
	assert.Empty(t, filtered)
}
