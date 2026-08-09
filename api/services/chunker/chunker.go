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
	Type      string // "declaration", "block", "window", or "heading"
	// Context is optional surrounding context (imports, parent type) prepended for summarization only.
	// It is not part of Content line ranges.
	Context string
}

// ChunkFile splits a source file into semantic chunks.
// For supported languages it uses top-level declaration patterns (or headings for markdown).
// For unsupported languages it falls back to blank-line-separated blocks (max 100 lines).
// Any resulting chunk > 200 lines is further split into overlapping 150-line windows with 30-line overlap.
func ChunkFile(filePath, content, language string) []Chunk {
	if language == "" {
		language = detectLanguage(filePath)
	}

	lines := splitLines(content)

	var chunks []Chunk
	switch {
	case language == LangMarkdown:
		chunks = chunkByMarkdownHeadings(lines)
	case isSupported(language):
		chunks = chunkByDeclarations(lines, language)
	default:
		chunks = chunkByBlankLines(lines)
	}

	chunks = splitOversizedChunks(chunks, lines)
	chunks = attachImportContext(chunks, lines, language)

	logger.Debug(context.Background(), "chunker: filePath="+filePath+" chunkCount="+fmt.Sprint(len(chunks)))

	return chunks
}

// DetectLanguage returns a normalised language name from a file path extension.
func DetectLanguage(filePath string) string {
	return detectLanguage(filePath)
}

// ExtractImportBlock returns a condensed import/package header for context (max ~40 lines).
func ExtractImportBlock(content, language string) string {
	lines := splitLines(content)
	return importBlockFromLines(lines, language)
}

func detectLanguage(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
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

func chunkByDeclarations(lines []string, language string) []Chunk {
	pattern, ok := topLevelPatterns[language]
	if !ok {
		return chunkByBlankLines(lines)
	}

	var chunks []Chunk
	start := 0

	for i := 1; i < len(lines); i++ {
		if pattern.MatchString(lines[i]) && i > start {
			// Include leading decorators/attributes that sit immediately above the declaration.
			chunkStart := start
			chunk := buildChunk(lines, chunkStart, i-1, "declaration")
			if chunk != nil {
				chunks = append(chunks, *chunk)
			}
			start = i
		}
	}
	if start < len(lines) {
		chunk := buildChunk(lines, start, len(lines)-1, "declaration")
		if chunk != nil {
			chunks = append(chunks, *chunk)
		}
	}
	return chunks
}

// chunkByMarkdownHeadings splits markdown on ATX headings (# … ######).
func chunkByMarkdownHeadings(lines []string) []Chunk {
	heading := regexpMustCompileMarkdownHeading()
	var chunks []Chunk
	start := 0

	for i := 1; i < len(lines); i++ {
		if heading.MatchString(lines[i]) && i > start {
			chunk := buildChunk(lines, start, i-1, "heading")
			if chunk != nil {
				chunks = append(chunks, *chunk)
			}
			start = i
		}
	}
	if start < len(lines) {
		chunk := buildChunk(lines, start, len(lines)-1, "heading")
		if chunk != nil {
			chunks = append(chunks, *chunk)
		}
	}
	if len(chunks) == 0 {
		return chunkByBlankLines(lines)
	}
	return chunks
}

func regexpMustCompileMarkdownHeading() *regexpMarkdown {
	return &regexpMarkdown{}
}

// Thin wrapper so we don't import regexp only for one pattern in this file path —
// use the package-level pattern instead.
type regexpMarkdown struct{}

func (r *regexpMarkdown) MatchString(s string) bool {
	if len(s) == 0 {
		return false
	}
	i := 0
	for i < len(s) && s[i] == '#' {
		i++
	}
	return i > 0 && i <= 6 && (i == len(s) || s[i] == ' ' || s[i] == '\t')
}

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
		if i-start+1 >= maxLines {
			flush(i)
		}
	}
	flush(len(lines) - 1)
	return chunks
}

func splitOversizedChunks(chunks []Chunk, lines []string) []Chunk {
	const (
		maxChunkLines   = 200
		windowSize      = 150
		windowStepLines = windowSize - 30
	)

	result := make([]Chunk, 0, len(chunks))
	for _, c := range chunks {
		startIdx := c.StartLine - 1
		endIdx := c.EndLine - 1
		size := endIdx - startIdx + 1

		if size <= maxChunkLines {
			result = append(result, c)
			continue
		}
		for wStart := startIdx; wStart <= endIdx; wStart += windowStepLines {
			wEnd := wStart + windowSize - 1
			if wEnd > endIdx {
				wEnd = endIdx
			}
			window := buildChunk(lines, wStart, wEnd, "window")
			if window != nil {
				window.Context = c.Context
				result = append(result, *window)
			}
			if wEnd == endIdx {
				break
			}
		}
	}
	return result
}

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
		StartLine: start + 1,
		EndLine:   end + 1,
		Type:      chunkType,
	}
}

func splitLines(content string) []string {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	return strings.Split(content, "\n")
}

// IsImportChunk reports whether a chunk consists entirely of import/package declarations
// and comments with no executable logic.
func IsImportChunk(chunk Chunk, language string) bool {
	patterns, ok := importLinePatterns[language]
	if !ok {
		return isCommentOrBlankChunk(chunk)
	}

	lines := splitLines(chunk.Content)
	substantiveLines := 0
	importLines := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if isCommentLine(trimmed) {
			// Pure comment lines don't count against import purity.
			continue
		}
		// Go multi-line import parentheses.
		if language == LangGo && (trimmed == "(" || trimmed == ")" || isGoImportPathLine(trimmed)) {
			substantiveLines++
			importLines++
			continue
		}
		substantiveLines++

		matched := false
		for _, prefix := range patterns {
			if prefix != "" && strings.HasPrefix(trimmed, prefix) {
				matched = true
				break
			}
		}
		if matched {
			importLines++
		}
	}

	if substantiveLines == 0 {
		// Comment-only or blank chunk — skip during summarisation.
		return true
	}
	return importLines == substantiveLines
}

func isCommentOrBlankChunk(chunk Chunk) bool {
	lines := splitLines(chunk.Content)
	hasSubstantive := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || isCommentLine(trimmed) {
			continue
		}
		hasSubstantive = true
		break
	}
	return !hasSubstantive
}

func isCommentLine(trimmed string) bool {
	for _, p := range commentLinePrefixes {
		if strings.HasPrefix(trimmed, p) {
			return true
		}
	}
	return false
}

func isGoImportPathLine(trimmed string) bool {
	// e.g. `"fmt"` or `alias "path"` inside import (
	if strings.HasPrefix(trimmed, "\"") && strings.HasSuffix(trimmed, "\"") {
		return true
	}
	// named import: foo "bar"
	parts := strings.Fields(trimmed)
	return len(parts) == 2 && strings.HasPrefix(parts[1], "\"") && strings.HasSuffix(parts[1], "\"")
}

func importBlockFromLines(lines []string, language string) string {
	patterns, ok := importLinePatterns[language]
	if !ok || len(patterns) == 0 {
		return ""
	}
	const maxLines = 40
	var out []string
	inGoParen := false
	for i, line := range lines {
		if i >= maxLines*2 && len(out) > 0 {
			break
		}
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if len(out) > 0 && !inGoParen {
				break
			}
			continue
		}
		if isCommentLine(trimmed) {
			if len(out) == 0 {
				continue
			}
			// Keep comments only while still in import section.
			if inGoParen {
				out = append(out, line)
			}
			continue
		}
		if language == LangGo {
			if strings.HasPrefix(trimmed, "package ") {
				out = append(out, line)
				continue
			}
			if strings.HasPrefix(trimmed, "import ") {
				out = append(out, line)
				if strings.Contains(trimmed, "(") && !strings.Contains(trimmed, ")") {
					inGoParen = true
				}
				continue
			}
			if inGoParen {
				out = append(out, line)
				if trimmed == ")" {
					inGoParen = false
				}
				continue
			}
			break
		}
		matched := false
		for _, prefix := range patterns {
			if prefix != "" && strings.HasPrefix(trimmed, prefix) {
				matched = true
				break
			}
		}
		if matched {
			out = append(out, line)
			if len(out) >= maxLines {
				break
			}
			continue
		}
		break
	}
	return strings.Join(out, "\n")
}

// attachImportContext sets Context on non-import chunks to a condensed import header.
func attachImportContext(chunks []Chunk, lines []string, language string) []Chunk {
	header := importBlockFromLines(lines, language)
	if header == "" {
		return chunks
	}
	for i := range chunks {
		if IsImportChunk(chunks[i], language) {
			continue
		}
		chunks[i].Context = header
	}
	return chunks
}

// ContentForSummary returns the text sent to the LLM for a chunk, including optional context.
func ContentForSummary(chunk Chunk) string {
	if chunk.Context == "" {
		return chunk.Content
	}
	return "// file context (imports/package):\n" + chunk.Context + "\n\n// chunk:\n" + chunk.Content
}
