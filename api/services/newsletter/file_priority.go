package newsletter

import (
	"path/filepath"
	"strings"

	"abstractions/db/models"
)

// FilePriority classifies a file's relevance to newsletter generation.
// High-priority files are actual source code; Low-priority files are
// documentation (markdown, rst, txt) and test files.
type FilePriority int

const (
	FilePriorityHigh FilePriority = iota
	FilePriorityLow
)

// lowPriorityExtensions are file extensions that belong to docs or config files
// that should only surface in newsletters explicitly about documentation.
var lowPriorityExtensions = map[string]struct{}{
	".md":  {},
	".mdx": {},
	".rst": {},
	".txt": {},
}

// ClassifyFile returns FilePriorityLow for documentation and test files,
// and FilePriorityHigh for everything else (actual source code).
func ClassifyFile(path string) FilePriority {
	base := strings.ToLower(filepath.Base(path))
	ext := strings.ToLower(filepath.Ext(path))

	// Documentation files by extension
	if _, ok := lowPriorityExtensions[ext]; ok {
		return FilePriorityLow
	}

	// Go test files: foo_test.go
	if strings.HasSuffix(base, "_test.go") {
		return FilePriorityLow
	}

	// JavaScript / TypeScript test and spec files
	// e.g. foo.test.ts, foo.spec.tsx, foo.test.js, foo.spec.jsx
	for _, pattern := range []string{".test.ts", ".test.tsx", ".test.js", ".test.jsx",
		".spec.ts", ".spec.tsx", ".spec.js", ".spec.jsx"} {
		if strings.HasSuffix(base, pattern) {
			return FilePriorityLow
		}
	}

	// Python test files: test_foo.py or foo_test.py
	if ext == ".py" && (strings.HasPrefix(base, "test_") || strings.HasSuffix(strings.TrimSuffix(base, ".py"), "_test")) {
		return FilePriorityLow
	}

	// Ruby test files: foo_test.rb or foo_spec.rb
	if ext == ".rb" && (strings.HasSuffix(strings.TrimSuffix(base, ".rb"), "_test") || strings.HasSuffix(strings.TrimSuffix(base, ".rb"), "_spec")) {
		return FilePriorityLow
	}

	return FilePriorityHigh
}

// PartitionSummaries splits a file-summaries map into high-priority (source code)
// and low-priority (docs + tests) subsets.
func PartitionSummaries(summaries map[string]string) (code map[string]string, lowPri map[string]string) {
	code = make(map[string]string)
	lowPri = make(map[string]string)
	for path, summary := range summaries {
		if ClassifyFile(path) == FilePriorityHigh {
			code[path] = summary
		} else {
			lowPri[path] = summary
		}
	}
	return code, lowPri
}

// docsKeywords are terms that signal a topic is specifically about documentation.
var docsKeywords = []string{
	"doc", "docs", "documentation", "readme", "changelog", "contributing",
	"guide", "tutorial", "wiki", "mdx", "markdown",
}

// testingKeywords are terms that signal a topic is specifically about testing.
var testingKeywords = []string{
	"test", "tests", "testing", "spec", "specs", "unit test", "integration test",
	"end-to-end", "e2e", "playwright", "jest", "pytest", "rspec", "testify",
	"coverage", "mock", "mocking", "fixture", "test suite",
}

// IsDocsTopic returns true when the topic title or description is about documentation.
func IsDocsTopic(title, description string) bool {
	haystack := strings.ToLower(title + " " + description)
	for _, kw := range docsKeywords {
		if strings.Contains(haystack, kw) {
			return true
		}
	}
	return false
}

// IsTestingTopic returns true when the topic title or description is about testing.
func IsTestingTopic(title, description string) bool {
	haystack := strings.ToLower(title + " " + description)
	for _, kw := range testingKeywords {
		if strings.Contains(haystack, kw) {
			return true
		}
	}
	return false
}

// FilterChunksByPriority removes low-priority (docs/test) chunks from the slice
// unless the caller signals that docs or testing content is wanted.
// When both includeMarkdown and includeTests are false, only high-priority chunks
// (actual source code) are returned.
func FilterChunksByPriority(chunks []models.ChunkEmbedding, includeMarkdown, includeTests bool) []models.ChunkEmbedding {
	if includeMarkdown && includeTests {
		return chunks
	}

	filtered := make([]models.ChunkEmbedding, 0, len(chunks))
	for _, c := range chunks {
		priority := ClassifyFile(c.RepoFilePath)
		if priority == FilePriorityHigh {
			filtered = append(filtered, c)
			continue
		}

		ext := strings.ToLower(filepath.Ext(c.RepoFilePath))
		_, isDoc := lowPriorityExtensions[ext]
		if isDoc && includeMarkdown {
			filtered = append(filtered, c)
			continue
		}
		if !isDoc && includeTests {
			// Low-priority non-doc file is a test file
			filtered = append(filtered, c)
		}
	}
	return filtered
}
