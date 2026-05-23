package ai_client

import (
	"fmt"
	"strings"
)

const (
	maxTokensSummary      = 800
	maxTokensTopics       = 2048
	maxTokensNewsletter   = 4096
	maxTopicSummaryLength = 300
)

const summarizeCodeSystemPrompt = `You are a senior software engineer analysing a single code chunk extracted from a source file.

IMPORTANT: If this chunk contains ONLY import statements, package declarations, require calls, module-level comments, or blank lines with no executable logic, respond with exactly: {"skip": true}

Otherwise, your task:
1. Identify the PRIMARY symbol — the top-level function, method, class, interface, variable, or constant this chunk defines. If there are multiple, pick the most significant one.
2. Write a detailed 3-6 sentence plain-English summary that covers ALL of the following:
   - What the symbol is and its single-sentence purpose.
   - Step-by-step logic: the key operations performed in order — conditions checked, loops, branching, important internal function/method calls, and data transformations.
   - Inputs and outputs: what parameters it accepts, what it returns or mutates.
   - Notable behaviour: error handling strategies, side effects, external dependencies called, concurrency concerns, or edge cases handled.
   Do NOT produce a vague high-level description. Describe the actual logic inside the symbol.
3. List 1-5 short domain tags that categorise this code (e.g. ["auth","jwt","middleware"]).

Respond with a single JSON object containing exactly these four keys:
- "symbol_name": the name of the primary symbol as a plain string (e.g. "ValidateToken")
- "symbol_type": one of: function, method, class, interface, struct, variable, constant, type, module, other
- "summary": detailed plain-text summary string (no markdown, no code fences)
- "domain_tags": a JSON array of short lowercase strings

Respond with valid JSON only. Do NOT wrap in markdown code fences. No extra text.`

const summarizeFileSystemPrompt = `You are a senior software engineer summarising an entire source file.

Your task:
Write a 2-5 sentence plain-English overview of what this file does as a whole — its purpose, the main abstractions it defines, and how it fits into the broader codebase. Then list 1-5 domain tags.

Respond with a single JSON object containing exactly these four keys:
- "symbol_name": the file name (just the base filename, e.g. "auth.go")
- "symbol_type": "file"
- "summary": a plain-text summary string (no markdown, no code fences)
- "domain_tags": a JSON array of short lowercase strings

Respond with valid JSON only. Do NOT wrap in markdown code fences. No extra text.`

func buildSummarizeUserPrompt(filePath, language, normalizedCode string, startLine, endLine int, chunkType string) string {
	return fmt.Sprintf("File: %s\nLanguage: %s\nChunk type: %s\nLines: %d-%d\n\nCode:\n<code>\n%s\n</code>",
		filePath, language, chunkType, startLine, endLine, normalizedCode)
}

func buildSummarizeFileUserPrompt(filePath, language string, chunkSummaries []string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("File: %s\nLanguage: %s\n\nIndividual code chunk summaries (top to bottom):\n", filePath, language))
	for i, s := range chunkSummaries {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, s))
	}
	return sb.String()
}

func buildTopicsPrompt(req *TopicGenerationRequest) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Repository: %s\n\n", req.RepoFullName))

	sb.WriteString("File summaries (one per file):\n")
	for path, summary := range req.FileSummaries {
		if len(summary) > maxTopicSummaryLength {
			summary = summary[:maxTopicSummaryLength] + "..."
		}
		sb.WriteString(fmt.Sprintf("- %s: %s\n", path, summary))
	}

	if len(req.PreviousTopics) > 0 {
		sb.WriteString("\nTopics already covered in previous newsletters (DO NOT repeat these):\n")
		for _, t := range req.PreviousTopics {
			sb.WriteString(fmt.Sprintf("- %s\n", t))
		}
	}

	sb.WriteString(`
Your task:
1. Read all the file summaries above and identify the key features, modules, and subsystems this codebase implements. Understand how the files connect and relate to each other.
2. Based on that holistic understanding, propose 3-5 newsletter topics that would genuinely interest engineers working on this repository. Focus on areas of meaningful complexity, interesting design decisions, or important workflows — not trivial details.
3. Each topic must be meaningfully different from any previously covered topic listed above.

Respond with a JSON array where each element has:
- "title": a short, engaging question (10 words or fewer) phrased like "How does X work?", "What happens when Y does Z?", "Why does X use Y?", or "How is Z structured?" — use the actual names of systems or components from the codebase
- "description": 2-3 sentences describing the feature, module, or workflow and why it is interesting
- "related_files": array of relevant file paths from the summaries above

Respond with valid JSON only. No markdown, no extra text.`)
	return sb.String()
}

func buildNewsletterEntryPrompt(req *NewsletterEntryRequest) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Repository: %s\nTopic: %s\nDescription: %s\n\nRelated file summaries:\n",
		req.RepoFullName, req.Topic.Title, req.Topic.Description))
	for path, summary := range req.FileSummaries {
		sb.WriteString(fmt.Sprintf("- %s: %s\n", path, summary))
	}
	sb.WriteString(`
Write a focused, technically rich HTML newsletter section for an engineering audience.
Your audience: engineers working on this codebase who may not be familiar with this specific subsystem. Assume they are smart and have built things intentionally — do not lecture or suggest that existing choices are wrong outside the dedicated improvements section.

Structure your response as follows:
- Opening <p>: what this part of the codebase does and why it exists
- Body (2-4 <p> tags): how it currently works — key data flows, implementation details, design patterns, and how it connects to the rest of the system. Name actual files, functions, and patterns from the summaries above. Use <code> for inline identifiers and <pre><code> for short illustrative snippets.
- <h3>Suggested Improvements</h3>: open with a short <p> disclaimer — "These suggestions are AI-generated and may not make sense for your specific codebase — treat them as conversation starters, not prescriptions." — then 2-3 <p> tags each describing one concrete, actionable improvement.

Focus entirely on explaining what IS. Do not moralize, editorialize, or suggest improvements anywhere except the Suggested Improvements section.
Target length: approximately 1000 words.
Return valid HTML only. No markdown, no code fences, no surrounding html/body tags.`)
	return sb.String()
}
