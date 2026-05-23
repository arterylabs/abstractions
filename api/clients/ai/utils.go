package ai_client

import (
	"encoding/json"
	"math"
	"strings"
)

// contentBudgetChars returns the maximum number of characters that variable content may occupy
// in a prompt, given the model's context window. Returns math.MaxInt when the model is unknown.
func contentBudgetChars(modelID string, maxOutputTokens, templateTokens int) int {
	modelInfo, ok := ModelByID(modelID)
	if !ok || modelInfo.ContextWindow == 0 {
		return math.MaxInt
	}
	availableTokens := modelInfo.ContextWindow - maxOutputTokens - templateTokens
	if availableTokens <= 0 {
		return 0
	}
	return availableTokens * 4
}

func truncateContent(content string, maxChars int) string {
	if len(content) <= maxChars {
		return content
	}
	return content[:maxChars] + "\n[... truncated to fit model context window ...]"
}

func truncateSummarySlice(summaries []string, totalBudgetChars int) []string {
	if len(summaries) == 0 {
		return summaries
	}
	charsPerEntry := totalBudgetChars / len(summaries)
	if charsPerEntry <= 0 {
		return summaries[:0]
	}
	truncated := make([]string, len(summaries))
	for i, s := range summaries {
		truncated[i] = truncateContent(s, charsPerEntry)
	}
	return truncated
}

func truncateSummaryMap(summaries map[string]string, totalBudgetChars int) map[string]string {
	if len(summaries) == 0 {
		return summaries
	}
	charsPerEntry := totalBudgetChars / len(summaries)
	if charsPerEntry <= 0 {
		return map[string]string{}
	}
	out := make(map[string]string, len(summaries))
	for filePath, summary := range summaries {
		out[filePath] = truncateContent(summary, charsPerEntry)
	}
	return out
}

// normalizeCodeContent strips excess leading indentation and collapses blank-line runs before
// code is embedded in a prompt. It preserves relative indentation so Python/YAML/Makefile
// semantics are not broken.
func normalizeCodeContent(code string) string {
	lines := strings.Split(strings.ReplaceAll(strings.ReplaceAll(code, "\r\n", "\n"), "\r", "\n"), "\n")

	// Find the common leading whitespace prefix across all non-blank lines.
	commonPrefix := ""
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		leadingWhitespace := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
		if commonPrefix == "" {
			commonPrefix = leadingWhitespace
			continue
		}
		sharedLength := len(commonPrefix)
		if len(leadingWhitespace) < sharedLength {
			sharedLength = len(leadingWhitespace)
		}
		for i := 0; i < sharedLength; i++ {
			if commonPrefix[i] != leadingWhitespace[i] {
				sharedLength = i
				break
			}
		}
		commonPrefix = commonPrefix[:sharedLength]
		if commonPrefix == "" {
			break
		}
	}

	result := make([]string, 0, len(lines))
	consecutiveBlankLines := 0
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			consecutiveBlankLines++
			if consecutiveBlankLines <= 2 {
				result = append(result, "")
			}
			continue
		}
		consecutiveBlankLines = 0
		trimmed := line
		if commonPrefix != "" && strings.HasPrefix(line, commonPrefix) {
			trimmed = line[len(commonPrefix):]
		}
		result = append(result, strings.TrimRight(trimmed, " \t"))
	}

	for len(result) > 0 && result[0] == "" {
		result = result[1:]
	}
	for len(result) > 0 && result[len(result)-1] == "" {
		result = result[:len(result)-1]
	}
	return strings.Join(result, "\n")
}

// stripMarkdownFences removes leading/trailing markdown code fences that models sometimes add
// despite being instructed not to.
func stripMarkdownFences(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		newline := strings.Index(s, "\n")
		if newline >= 0 {
			s = s[newline+1:]
		}
	}
	if strings.HasSuffix(s, "```") {
		s = s[:strings.LastIndex(s, "```")]
	}
	return strings.TrimSpace(s)
}

// --- Response parsers ---

type summarizeJSONResponse struct {
	Summary    string   `json:"summary"`
	DomainTags []string `json:"domain_tags"`
	SymbolName string   `json:"symbol_name"`
	SymbolType string   `json:"symbol_type"`
	Skip       bool     `json:"skip"`
}

func parseSummarizeResult(rawResponse string) *SummarizeResult {
	cleaned := stripMarkdownFences(rawResponse)
	var parsed summarizeJSONResponse
	if err := json.Unmarshal([]byte(cleaned), &parsed); err != nil {
		return &SummarizeResult{Summary: strings.TrimSpace(rawResponse)}
	}
	if parsed.Skip {
		return nil
	}
	return &SummarizeResult{
		Summary:    parsed.Summary,
		DomainTags: parsed.DomainTags,
		SymbolName: parsed.SymbolName,
		SymbolType: parsed.SymbolType,
	}
}

type topicJSONResponse struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	RelatedFiles []string `json:"related_files"`
}

func parseTopics(rawResponse string) []Topic {
	var parsed []topicJSONResponse
	if err := json.Unmarshal([]byte(rawResponse), &parsed); err != nil {
		return nil
	}
	topics := make([]Topic, len(parsed))
	for i, t := range parsed {
		topics[i] = Topic{Title: t.Title, Description: t.Description, RelatedFiles: t.RelatedFiles}
	}
	return topics
}
