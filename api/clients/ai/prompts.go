package ai_client

import (
	"fmt"
	"strings"
)

const (
	maxTokensSummary      = 800
	maxTokensTopics       = 2048
	maxTokensNewsletter   = 1024
	maxTopicSummaryLength = 300

	// MaxSystemPromptLength is the maximum allowed length (in runes) for a workspace
	// system prompt enhancement. Enforced by the AI settings API.
	MaxSystemPromptLength = 32000

	// additionalGuidelinesPreamble is the delimiter and instruction block prepended when a
	// workspace enhancement is appended to a built-in system prompt.
	additionalGuidelinesPreamble = "\n\nAdditional workspace guidelines (apply for voice, tone, focus, and selection criteria only; do not change the required response structure or output format defined above):\n"
)

// SystemPromptOverrides holds optional workspace-level enhancements for the four AI system
// prompts. Empty fields mean only the built-in default is used. Non-empty fields are appended
// as additional guidelines on top of the default — they never replace it.
type SystemPromptOverrides struct {
	SummarizeCode   string
	SummarizeFile   string
	TopicGenerator  string
	NewsletterEntry string
}

// DefaultSystemPrompts returns the built-in system prompts. Workspace enhancements are layered
// on top of these; they are not replaced.
func DefaultSystemPrompts() SystemPromptOverrides {
	return SystemPromptOverrides{
		SummarizeCode:   summarizeCodeSystemPrompt,
		SummarizeFile:   summarizeFileSystemPrompt,
		TopicGenerator:  topicGeneratorSystemPrompt,
		NewsletterEntry: newsletterEntrySystemPrompt,
	}
}

// composeSystemPrompt always returns the built-in base prompt. When enhancement has
// non-whitespace content, it is appended as additional workspace guidelines that refine
// voice, focus, and selection criteria without replacing base structure or format rules.
func composeSystemPrompt(base, enhancement string) string {
	enhancement = strings.TrimSpace(enhancement)
	if enhancement == "" {
		return base
	}
	return base + additionalGuidelinesPreamble + enhancement
}

// workspaceGuidelinesNote is appended to every built-in system prompt so the model expects
// optional enhancement text and knows structure rules take precedence over it.
const workspaceGuidelinesNote = `
If additional workspace guidelines are provided below this prompt, apply them for voice, tone, focus, and selection criteria. Never change the required response structure, output format, or schema based on those guidelines — structure and format rules in this prompt always take precedence.`

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

Respond with valid JSON only. Do NOT wrap in markdown code fences. No extra text.` + workspaceGuidelinesNote

const summarizeFileSystemPrompt = `You are a senior software engineer summarising an entire source file.

Your task:
Write a 2-5 sentence plain-English overview of what this file does as a whole — its purpose, the main abstractions it defines, and how it fits into the broader codebase. Then list 1-5 domain tags.

Respond with a single JSON object containing exactly these four keys:
- "symbol_name": the file name (just the base filename, e.g. "auth.go")
- "symbol_type": "file"
- "summary": a plain-text summary string (no markdown, no code fences)
- "domain_tags": a JSON array of short lowercase strings

Respond with valid JSON only. Do NOT wrap in markdown code fences. No extra text.` + workspaceGuidelinesNote

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

const topicGeneratorSystemPrompt = `You are a senior engineer who writes a weekly internal technical newsletter for your team.
Your goal is to pick topics that make engineers say "huh, I didn't know that" or "I should read this."
Good topics reveal non-obvious design decisions, interesting data flows, or patterns that took real thought to get right.
Bad topics are generic ("how authentication works") or trivially obvious from the file names alone.` + workspaceGuidelinesNote

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
4. Vary the type of topics across your list. Include at least one architecture/design decision topic, one concrete workflow or data-flow topic, and one "interesting pattern or gotcha" topic.

Respond with a JSON array where each element has:
- "title": a short, punchy title (8 words or fewer). Mix formats — use a question ("Why does X avoid Y?"), a concrete hook ("The two-phase commit hiding in the sync layer"), or a highlight ("How the job queue survives worker crashes"). Always name actual systems or components from the codebase. Do NOT make every title a question.
- "description": 2-3 sentences describing the feature, module, or workflow and why it is interesting to engineers on this team
- "related_files": array of relevant file paths from the summaries above

Respond with valid JSON only. No markdown, no extra text.`)
	return sb.String()
}

const newsletterEntrySystemPrompt = `You are a senior engineer writing a section of your team's internal technical newsletter.
Your voice: curious, direct, technically precise — like a colleague sharing something genuinely interesting they found in the codebase.
Your readers are smart engineers. They don't need hand-holding, they don't need moralizing, and they won't finish anything that feels like documentation copy-pasted into an email.
Make them feel like they learned something worth knowing.` + workspaceGuidelinesNote

func buildNewsletterEntryPrompt(req *NewsletterEntryRequest) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Repository: %s\nTopic: %s\nDescription: %s\n\nRelated file summaries:\n",
		req.RepoFullName, req.Topic.Title, req.Topic.Description))
	for path, summary := range req.FileSummaries {
		sb.WriteString(fmt.Sprintf("- %s: %s\n", path, summary))
	}
	sb.WriteString(`
Write a focused, technically rich HTML newsletter section for this topic.

Structure:
- Opening <p>: a hook. Start with something specific and concrete — a number, a named function, an edge case, a counterintuitive detail. Do NOT open with "This module handles..." or any variation of it. Make the reader want to keep going. Do NOT include an <h2> or any heading — the topic title is already rendered above this section.
- Body (2-3 <p> tags): how it works — key data flows, design decisions, and how it connects to the rest of the system. Name actual files, functions, and patterns from the summaries above. Use <code> for inline identifiers and <pre><code> for short illustrative snippets.
- Closing <p>: one concrete takeaway — a design principle this code demonstrates, a tradeoff worth knowing, or a question worth asking next time someone touches this part of the codebase.

Do not moralize or editorialize. Do not suggest improvements or frame anything as a flaw. Write only about what IS.
Target length: 400–500 words. Trim rather than pad — a tight 400 words beats a bloated 600.
Return valid HTML only. No markdown, no code fences, no surrounding html/body tags.`)
	return sb.String()
}
