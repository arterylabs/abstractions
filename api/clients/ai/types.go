package ai_client

// SummarizeResult holds the AI-generated summary and extracted domain tags for a code chunk.
type SummarizeResult struct {
	// Summary is a plain-text description of what the code does.
	Summary string
	// DomainTags are short domain labels (e.g. ["auth", "jwt"]).
	DomainTags []string
	// SymbolName is the name of the primary symbol in the chunk (e.g. "ValidateToken").
	// Empty for file-level summaries.
	SymbolName string
	// SymbolType is the kind of symbol (e.g. "function", "class", "method", "variable", "file").
	SymbolType string
}

// TopicGenerationRequest describes a newsletter topic generation job.
type TopicGenerationRequest struct {
	// RepoFullName is the GitHub "owner/repo" identifier.
	RepoFullName string
	// FileSummaries maps file paths to their holistic file-level summaries.
	FileSummaries map[string]string
	// PreviousTopics lists topic titles from recently sent newsletters so the AI
	// can avoid proposing topics that have already been covered.
	PreviousTopics []string
}

// Topic is a single generated newsletter topic.
type Topic struct {
	Title        string
	Description  string
	RelatedFiles []string
}

// NewsletterEntryRequest provides context for generating one newsletter section.
type NewsletterEntryRequest struct {
	RepoFullName string
	Topic        Topic
	// FileSummaries provides the raw summaries for the files related to this topic.
	FileSummaries map[string]string
}

// NewsletterEntry is the generated HTML/Markdown content for one topic section.
type NewsletterEntry struct {
	Title   string
	Content string
}
