package chunker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSupportedLanguages_AllHavePatterns(t *testing.T) {
	for _, lang := range SupportedLanguages {
		_, ok := topLevelPatterns[lang]
		assert.True(t, ok, "SupportedLanguages entry %q has no topLevelPatterns entry", lang)
	}
}

// Sample pattern coverage; full per-language suites are redundant with ChunkFile_* tests.
func TestTopLevelPatterns_SampleLanguages(t *testing.T) {
	cases := []struct {
		lang    string
		content string
	}{
		{"go", "func Foo() {}\n"},
		{"typescript", "export function foo() {}\n"},
		{"python", "def foo():\n    pass\n"},
		{"rust", "fn main() {}\n"},
	}
	for _, tc := range cases {
		t.Run(tc.lang, func(t *testing.T) {
			re, ok := topLevelPatterns[tc.lang]
			require.True(t, ok)
			assert.True(t, re.MatchString(tc.content), "expected pattern match for %s", tc.lang)
		})
	}
}

func TestIsImportChunk(t *testing.T) {
	assert.True(t, IsImportChunk(Chunk{Content: "import (\n\t\"fmt\"\n)\n"}, "go"))
	assert.False(t, IsImportChunk(Chunk{Content: "func main() {}\n"}, "go"))
	assert.True(t, IsImportChunk(Chunk{Content: "import { x } from './x';\n"}, "typescript"))
	assert.False(t, IsImportChunk(Chunk{Content: "import foo\n"}, "unknown"))
}
