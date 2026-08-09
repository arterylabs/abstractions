package github_client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldIgnore(t *testing.T) {
	cases := []struct {
		path   string
		size   int
		ignore bool
	}{
		{"node_modules/lodash/index.js", 100, true},
		{"frontend/node_modules/react/index.js", 100, true},
		{"package-lock.json", 1024, true},
		{"yarn.lock", 1024, true},
		{"go.sum", 1024, true},
		{"dist/bundle.min.js", 100, true},
		{"dist/bundle.js.map", 100, true},
		{"public/logo.png", 10000, true},
		{"db/migrations/000001_create_users.up.sql", 100, true},
		{"large_file.go", 600 * 1024, true},
		{"dist/index.js", 500, true},
		{".env", 100, true},
		{".env.local", 100, true},
		{"cmd/main.go", 2048, false},
		{"src/components/Button.tsx", 1024, false},
		{"api/views.py", 4096, false},
	}

	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			assert.Equal(t, tc.ignore, ShouldIgnore(tc.path, tc.size))
		})
	}
}

func TestFilterFileTree_RemovesIgnoredFiles(t *testing.T) {
	input := []FileTreeItem{
		{Path: "cmd/main.go", SHA: "a1", Size: 1024},
		{Path: "node_modules/lodash/index.js", SHA: "b1", Size: 500},
		{Path: "yarn.lock", SHA: "c1", Size: 20000},
		{Path: "dist/bundle.min.js", SHA: "d1", Size: 200},
		{Path: "public/logo.png", SHA: "e1", Size: 8000},
		{Path: "large.go", SHA: "f1", Size: 600 * 1024},
		{Path: "db/migrations/000001.up.sql", SHA: "g1", Size: 300},
		{Path: "services/user_service.go", SHA: "h1", Size: 3000},
	}

	result := FilterFileTree(input)

	assert.Len(t, result, 2)
	paths := make([]string, len(result))
	for i, item := range result {
		paths[i] = item.Path
	}
	assert.Contains(t, paths, "cmd/main.go")
	assert.Contains(t, paths, "services/user_service.go")
}
