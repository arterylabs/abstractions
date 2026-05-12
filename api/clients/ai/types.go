package ai_client

// SummarizeResult holds the AI-generated summary and extracted domain tags for a code chunk.
type SummarizeResult struct {
	Summary    string
	DomainTags []string
	SymbolName string
	SymbolType string
}

// TopicGenerationRequest describes a newsletter topic generation job.
type TopicGenerationRequest struct {
	// RepoFullName is the GitHub "owner/repo" identifier.
	RepoFullName string
	// FileSummaries maps file paths to their per-chunk summaries.
	FileSummaries map[string]string
	// PreviousTopics is a list of topic titles already covered in recent newsletters.
	// The AI should avoid generating duplicates.
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
