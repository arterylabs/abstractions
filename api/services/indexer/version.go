// Package indexer holds shared constants and helpers for the repository indexing pipeline.
package indexer

// Version is bumped when chunking, summarization prompts, or embedding text format change
// in a way that should invalidate SHA-skipped re-index of unchanged files.
const Version = "2"

// MaxFilesPerIndex is a soft cap on how many files are processed in a single index run.
// High-priority source files are preferred when the tree is larger than this.
const MaxFilesPerIndex = 8000

// MaxProcessAttempts is how many times SummarizeFileJob may fail for a file before
// the file is marked processed-with-error so the repository can finish indexing.
const MaxProcessAttempts = 3

// MaxConcurrentChunkSummaries limits parallel LLM calls within a single file.
const MaxConcurrentChunkSummaries = 4
