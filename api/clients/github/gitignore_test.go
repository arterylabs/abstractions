package github_client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGitIgnore_BasicPatterns(t *testing.T) {
	m := ParseGitIgnore(`
# comment
node_modules/
*.log
build/
`)
	assert.True(t, m.Match("node_modules/foo.js"))
	assert.True(t, m.Match("app.log"))
	assert.True(t, m.Match("build"))
	assert.True(t, m.Match("build/out.js"))
	assert.False(t, m.Match("src/main.go"))
}

func TestParseGitIgnore_Negation(t *testing.T) {
	m := ParseGitIgnore("*.log\n!keep.log\n")
	assert.True(t, m.Match("drop.log"))
	assert.False(t, m.Match("keep.log"))
}

func TestShouldIgnoreWithMatcher(t *testing.T) {
	m := ParseGitIgnore("generated/\n")
	assert.True(t, ShouldIgnoreWithMatcher("generated/foo.go", 100, m))
	assert.False(t, ShouldIgnoreWithMatcher("src/foo.go", 100, m))
	assert.True(t, ShouldIgnoreWithMatcher("node_modules/x.js", 100, m))
}
