package indexer

import (
	"fmt"
	"strings"
)

// BuildEmbedText constructs the hybrid embedding input used for vector search.
// Including path, symbol, language, and tags improves retrieval for identifier-heavy queries.
func BuildEmbedText(filePath, language, symbolName, symbolType, summary string, tags []string) string {
	var b strings.Builder
	b.WriteString("path: ")
	b.WriteString(filePath)
	if language != "" {
		b.WriteString("\nlanguage: ")
		b.WriteString(language)
	}
	if symbolName != "" {
		b.WriteString("\nsymbol: ")
		b.WriteString(symbolName)
		if symbolType != "" {
			b.WriteString(fmt.Sprintf(" (%s)", symbolType))
		}
	} else if symbolType != "" {
		b.WriteString("\nsymbol_type: ")
		b.WriteString(symbolType)
	}
	if len(tags) > 0 {
		b.WriteString("\ntags: ")
		b.WriteString(strings.Join(tags, ", "))
	}
	if summary != "" {
		b.WriteString("\nsummary: ")
		b.WriteString(summary)
	}
	return b.String()
}

// BuildSearchDocument is the keyword-search document stored alongside the embedding.
func BuildSearchDocument(filePath, language, symbolName, symbolType, summary string, tags []string) string {
	parts := []string{filePath, language, symbolName, symbolType, summary}
	if len(tags) > 0 {
		parts = append(parts, strings.Join(tags, " "))
	}
	return strings.Join(filterEmpty(parts), " ")
}

func filterEmpty(parts []string) []string {
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			out = append(out, p)
		}
	}
	return out
}
