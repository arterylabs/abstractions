package indexer

import (
	"path/filepath"
	"strings"
)

// FilePriority classifies a file's value for indexing order and retrieval.
type FilePriority int

const (
	// PriorityHigh is primary source code — index first.
	PriorityHigh FilePriority = iota
	// PriorityMedium is config / infra that is useful but secondary.
	PriorityMedium
	// PriorityLow is docs and tests — index last or on a low queue.
	PriorityLow
)

var lowPriorityExtensions = map[string]struct{}{
	".md":  {},
	".mdx": {},
	".rst": {},
	".txt": {},
}

var mediumPriorityExtensions = map[string]struct{}{
	".json": {},
	".yaml": {},
	".yml":  {},
	".toml": {},
	".xml":  {},
	".env":  {},
	".sql":  {},
	".tf":   {},
}

var highPriorityPathPrefixes = []string{
	"cmd/", "src/", "app/", "lib/", "pkg/", "internal/", "api/", "server/",
	"services/", "packages/", "core/", "backend/", "frontend/src/",
}

// ClassifyFile returns indexing priority for a repository path.
func ClassifyFile(path string) FilePriority {
	base := strings.ToLower(filepath.Base(path))
	ext := strings.ToLower(filepath.Ext(path))

	if _, ok := lowPriorityExtensions[ext]; ok {
		return PriorityLow
	}

	// Go / JS / Python / Ruby test files
	if strings.HasSuffix(base, "_test.go") {
		return PriorityLow
	}
	for _, pattern := range []string{".test.ts", ".test.tsx", ".test.js", ".test.jsx",
		".spec.ts", ".spec.tsx", ".spec.js", ".spec.jsx"} {
		if strings.HasSuffix(base, pattern) {
			return PriorityLow
		}
	}
	if ext == ".py" && (strings.HasPrefix(base, "test_") || strings.HasSuffix(strings.TrimSuffix(base, ".py"), "_test")) {
		return PriorityLow
	}
	if ext == ".rb" && (strings.HasSuffix(strings.TrimSuffix(base, ".rb"), "_test") || strings.HasSuffix(strings.TrimSuffix(base, ".rb"), "_spec")) {
		return PriorityLow
	}

	// Test directories
	lower := strings.ToLower(path)
	if strings.Contains(lower, "/__tests__/") || strings.Contains(lower, "/testdata/") ||
		strings.HasPrefix(lower, "test/") || strings.HasPrefix(lower, "tests/") ||
		strings.Contains(lower, "/e2e/") {
		return PriorityLow
	}

	if _, ok := mediumPriorityExtensions[ext]; ok {
		return PriorityMedium
	}

	return PriorityHigh
}

// SortKey returns a sort key for enqueue order (lower = earlier).
func SortKey(path string) int {
	pri := ClassifyFile(path)
	base := 0
	switch pri {
	case PriorityHigh:
		base = 0
	case PriorityMedium:
		base = 1000
	default:
		base = 2000
	}

	lower := strings.ToLower(path)
	for i, prefix := range highPriorityPathPrefixes {
		if strings.HasPrefix(lower, prefix) || strings.Contains(lower, "/"+prefix) {
			return base + i
		}
	}
	return base + 500
}
