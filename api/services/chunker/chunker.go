package chunker

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"abstractions/logger"
)

// Chunk represents a semantically meaningful segment of a source file.
type Chunk struct {
	Content   string
	StartLine int
	EndLine   int
	Type      string // "declaration", "block", or "window"
}

// ChunkFile splits a source file into semantic chunks.
// For supported languages it uses top-level declaration patterns.
// For unsupported languages it falls back to blank-line-separated blocks (max 100 lines).
// Any resulting chunk > 200 lines is further split into overlapping 150-line windows with 30-line overlap.
// Note: chunkContent is never logged; only filePath and chunkCount are logged.
func ChunkFile(filePath, content, language string) []Chunk {
	if language == "" {
		language = detectLanguage(filePath)
	}

	lines := splitLines(content)

	var chunks []Chunk
	if isSupported(language) {
		chunks = chunkByDeclarations(lines, language)
	} else {
		chunks = chunkByBlankLines(lines)
	}

	chunks = splitOversizedChunks(chunks, lines)

	logger.Debug(context.Background(), "chunker: filePath="+filePath+" chunkCount="+fmt.Sprint(len(chunks)))

	return chunks
}

// DetectLanguage returns a normalised language name from a file path extension.
// Exported so it can be used independently (e.g., in the indexing pipeline).
func DetectLanguage(filePath string) string {
	return detectLanguage(filePath)
}

// ---- internal helpers ----

func detectLanguage(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	// Programming languages with declaration-based chunking support
	case ".go":
		return LangGo
	case ".ts", ".tsx":
		return LangTypeScript
	case ".js", ".jsx", ".mjs", ".cjs":
		return LangJavaScript
	case ".py":
		return LangPython
	case ".rb":
		return LangRuby
	case ".java":
		return LangJava
	case ".rs":
		return LangRust
	// Additional programming languages (blank-line fallback chunking)
	case ".c", ".h":
		return LangC
	case ".cpp", ".cc", ".cxx", ".hpp":
		return LangCpp
	case ".cs":
		return LangCSharp
	case ".php":
		return LangPHP
	case ".swift":
		return LangSwift
	case ".kt", ".kts":
		return LangKotlin
	case ".scala":
		return LangScala
	// Config and data formats (blank-line fallback chunking)
	case ".json":
		return LangJSON
	case ".toml":
		return LangTOML
	case ".yaml", ".yml":
		return LangYAML
	case ".xml":
		return LangXML
	case ".html", ".htm":
		return LangHTML
	case ".css":
		return LangCSS
	case ".scss", ".sass":
		return LangSCSS
	case ".sql":
		return LangSQL
	case ".md", ".mdx":
		return LangMarkdown
	case ".graphql", ".gql":
		return LangGraphQL
	case ".proto":
		return LangProtobuf
	case ".sh", ".bash", ".zsh":
		return LangShell
	case ".tf", ".tfvars":
		return LangTerraform
	case ".env":
		return LangDotenv
	}

	// Basename-based detection for extension-less files and well-known config filenames.
	base := filepath.Base(filePath)
	switch base {
	case "Dockerfile":
		return LangDockerfile
	case "Makefile":
		return LangMakefile
	case "Gemfile", "Rakefile", "Vagrantfile":
		return LangRuby
	case "Procfile":
		return LangShell
	case ".gitignore", ".dockerignore":
		return LangGitignore
	case ".editorconfig":
		return LangINI
	}

	return LangUnknown
}

func isSupported(lang string) bool {
	for _, l := range SupportedLanguages {
		if l == lang {
			return true
		}
	}
	return false
}

// chunkByDeclarations groups lines around top-level declaration boundaries.
func chunkByDeclarations(lines []string, language string) []Chunk {
	pattern, ok := topLevelPatterns[language]
	if !ok {
		return chunkByBlankLines(lines)
	}

	var chunks []Chunk
	start := 0

	for i := 1; i < len(lines); i++ {
		if pattern.MatchString(lines[i]) && i > start {
			chunk := buildChunk(lines, start, i-1, "declaration")
			if chunk != nil {
				chunks = append(chunks, *chunk)
			}
			start = i
		}
	}
	// last chunk
	if start < len(lines) {
		chunk := buildChunk(lines, start, len(lines)-1, "declaration")
		if chunk != nil {
			chunks = append(chunks, *chunk)
		}
	}
	return chunks
}

// chunkByBlankLines groups lines into blocks separated by blank lines, max 100 lines.
func chunkByBlankLines(lines []string) []Chunk {
	const maxLines = 100

	var chunks []Chunk
	start := -1

	flush := func(end int) {
		if start < 0 {
			return
		}
		chunk := buildChunk(lines, start, end, "block")
		if chunk != nil {
			chunks = append(chunks, *chunk)
		}
		start = -1
	}

	for i, line := range lines {
		isBlank := strings.TrimSpace(line) == ""
		if isBlank {
			flush(i - 1)
			continue
		}
		if start < 0 {
			start = i
		}
		// enforce max-lines boundary
		if i-start+1 >= maxLines {
			flush(i)
		}
	}
	flush(len(lines) - 1)
	return chunks
}

// splitOversizedChunks breaks any chunk > 200 lines into overlapping 150-line windows.
func splitOversizedChunks(chunks []Chunk, lines []string) []Chunk {
	const (
		maxChunkLines   = 200
		windowSize      = 150
		windowStepLines = windowSize - 30 // 120 — 30-line overlap
	)

	result := make([]Chunk, 0, len(chunks))
	for _, c := range chunks {
		// Convert to 0-indexed for line-slice operations
		startIdx := c.StartLine - 1
		endIdx := c.EndLine - 1
		size := endIdx - startIdx + 1

		if size <= maxChunkLines {
			result = append(result, c)
			continue
		}
		// split into overlapping windows
		for wStart := startIdx; wStart <= endIdx; wStart += windowStepLines {
			wEnd := wStart + windowSize - 1
			if wEnd > endIdx {
				wEnd = endIdx
			}
			window := buildChunk(lines, wStart, wEnd, "window")
			if window != nil {
				result = append(result, *window)
			}
			if wEnd == endIdx {
				break
			}
		}
	}
	return result
}

// buildChunk assembles a Chunk from the given [start, end] line range (0-indexed).
// Returns nil when the range is empty or contains only whitespace.
func buildChunk(lines []string, start, end int, chunkType string) *Chunk {
	if start > end || start < 0 || end >= len(lines) {
		return nil
	}
	segment := strings.Join(lines[start:end+1], "\n")
	if strings.TrimSpace(segment) == "" {
		return nil
	}
	return &Chunk{
		Content:   segment,
		StartLine: start + 1, // 1-indexed for human-readable output
		EndLine:   end + 1,
		Type:      chunkType,
	}
}

func splitLines(content string) []string {
	// Normalise line endings then split
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	return strings.Split(content, "\n")
}

// IsImportChunk reports whether a chunk consists entirely of import/package declarations
// and comments with no executable logic. Such chunks produce no meaningful summary and
// should be skipped during the summarisation pipeline.
func IsImportChunk(chunk Chunk, language string) bool {
	patterns, ok := importLinePatterns[language]
	if !ok {
		return false
	}

	lines := splitLines(chunk.Content)
	substantiveLines := 0
	importLines := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		substantiveLines++

		for _, prefix := range patterns {
			if strings.HasPrefix(trimmed, prefix) {
				importLines++
				break
			}
		}
	}

	// Require at least one substantive line and all of them must match import patterns.
	return substantiveLines > 0 && importLines == substantiveLines
}
