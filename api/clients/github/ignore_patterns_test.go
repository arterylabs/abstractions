package github_client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldIgnore_NodeModules(t *testing.T) {
	assert.True(t, ShouldIgnore("node_modules/lodash/index.js", 100))
}

func TestShouldIgnore_NestedNodeModules(t *testing.T) {
	assert.True(t, ShouldIgnore("frontend/node_modules/react/index.js", 100))
}

func TestShouldIgnore_LockFile(t *testing.T) {
	assert.True(t, ShouldIgnore("package-lock.json", 1024))
	assert.True(t, ShouldIgnore("yarn.lock", 1024))
	assert.True(t, ShouldIgnore("go.sum", 1024))
	assert.True(t, ShouldIgnore("Gemfile.lock", 1024))
}

func TestShouldIgnore_MinifiedFiles(t *testing.T) {
	assert.True(t, ShouldIgnore("dist/bundle.min.js", 100))
	assert.True(t, ShouldIgnore("public/style.min.css", 100))
}

func TestShouldIgnore_SourceMaps(t *testing.T) {
	assert.True(t, ShouldIgnore("dist/bundle.js.map", 100))
}

func TestShouldIgnore_Images(t *testing.T) {
	assert.True(t, ShouldIgnore("public/logo.png", 10000))
	assert.True(t, ShouldIgnore("assets/banner.jpg", 10000))
	assert.True(t, ShouldIgnore("icons/favicon.ico", 500))
}

func TestShouldIgnore_MigrationSQL(t *testing.T) {
	assert.True(t, ShouldIgnore("db/migrations/000001_create_users.up.sql", 100))
	assert.True(t, ShouldIgnore("db/migrations/000001_create_users.down.sql", 100))
}

func TestShouldIgnore_LargeFile(t *testing.T) {
	assert.True(t, ShouldIgnore("large_file.go", 600*1024))
}

func TestShouldIgnore_ValidGoFile(t *testing.T) {
	assert.False(t, ShouldIgnore("cmd/main.go", 2048))
}

func TestShouldIgnore_ValidTypeScriptFile(t *testing.T) {
	assert.False(t, ShouldIgnore("src/components/Button.tsx", 1024))
}

func TestShouldIgnore_ValidPythonFile(t *testing.T) {
	assert.False(t, ShouldIgnore("api/views.py", 4096))
}

func TestShouldIgnore_DistDir(t *testing.T) {
	assert.True(t, ShouldIgnore("dist/index.js", 500))
}

func TestShouldIgnore_EnvFile(t *testing.T) {
	assert.True(t, ShouldIgnore(".env", 100))
	assert.True(t, ShouldIgnore(".env.local", 100))
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

	// Only cmd/main.go and services/user_service.go should survive.
	assert.Len(t, result, 2)
	paths := make([]string, len(result))
	for i, item := range result {
		paths[i] = item.Path
	}
	assert.Contains(t, paths, "cmd/main.go")
	assert.Contains(t, paths, "services/user_service.go")
}

func TestFilterFileTree_AllValidFiles(t *testing.T) {
	input := []FileTreeItem{
		{Path: "main.go", SHA: "a", Size: 100},
		{Path: "README.md", SHA: "b", Size: 200},
		{Path: "Makefile", SHA: "c", Size: 50},
	}
	result := FilterFileTree(input)
	assert.Len(t, result, 3)
}

func TestFilterFileTree_EmptyInput(t *testing.T) {
	result := FilterFileTree([]FileTreeItem{})
	assert.Nil(t, result)
}
