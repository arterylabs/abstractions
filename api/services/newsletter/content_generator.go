package newsletter

import (
	ai_client "abstractions/clients/ai"
	"abstractions/repositories"
	"context"
	"fmt"

	"github.com/google/uuid"
)

const topKChunks = 12

// NewsletterEntry holds the generated HTML content for a single topic section.
type NewsletterEntry struct {
	TopicTitle          string
	Content             string // HTML, ≤300 words
	FilePathsReferenced []string
}

// ContentGenerator builds newsletter entries from topic proposals and chunk embeddings.
type ContentGenerator struct {
	aiProvider ai_client.AIProvider
	chunkRepo  *repositories.ChunkEmbeddingRepository
	repoRepo   *repositories.RepositoryRepository
}

func NewContentGenerator(
	aiProvider ai_client.AIProvider,
	chunkRepo *repositories.ChunkEmbeddingRepository,
	repoRepo *repositories.RepositoryRepository,
) *ContentGenerator {
	return &ContentGenerator{
		aiProvider: aiProvider,
		chunkRepo:  chunkRepo,
		repoRepo:   repoRepo,
	}
}

// GenerateNewsletter produces a NewsletterEntry for each topic.
// For each topic it:
//  1. Embeds the topic title
//  2. Queries chunk_embeddings for the topK=12 most similar chunks
//  3. Builds a FileSummaries map from the retrieved chunks
//  4. Calls AIProvider.GenerateNewsletterEntry
func (g *ContentGenerator) GenerateNewsletter(ctx context.Context, repoID uuid.UUID, topics []Topic) ([]NewsletterEntry, error) {
	repo, err := g.repoRepo.GetByID(ctx, repoID)
	if err != nil || repo == nil {
		return nil, fmt.Errorf("content_generator: repo not found: %w", err)
	}

	entries := make([]NewsletterEntry, 0, len(topics))

	for _, topic := range topics {
		embedding, err := g.aiProvider.GenerateEmbedding(ctx, topic.Title)
		if err != nil {
			return nil, fmt.Errorf("content_generator: embed topic %q: %w", topic.Title, err)
		}

		chunks, err := g.chunkRepo.QuerySimilar(ctx, repoID, embedding, topKChunks)
		if err != nil {
			return nil, fmt.Errorf("content_generator: query similar for topic %q: %w", topic.Title, err)
		}

		// Remove low-priority (docs/test) chunks unless the topic is specifically
		// about documentation or testing.
		includeMarkdown := IsDocsTopic(topic.Title, topic.Description)
		includeTests := IsTestingTopic(topic.Title, topic.Description)
		chunks = FilterChunksByPriority(chunks, includeMarkdown, includeTests)

		// Build summaries map and collect unique file paths from the retrieved chunks
		fileSummaries := make(map[string]string, len(chunks))
		fileSet := map[string]struct{}{}
		for _, c := range chunks {
			fileSummaries[c.RepoFilePath] = c.Summary
			fileSet[c.RepoFilePath] = struct{}{}
		}

		filePaths := make([]string, 0, len(fileSet))
		for p := range fileSet {
			filePaths = append(filePaths, p)
		}

		aiEntry, err := g.aiProvider.GenerateNewsletterEntry(ctx, &ai_client.NewsletterEntryRequest{
			RepoFullName: repo.Name,
			Topic: ai_client.Topic{
				Title:        topic.Title,
				Description:  topic.Description,
				RelatedFiles: topic.RelatedFiles,
			},
			FileSummaries: fileSummaries,
		})
		if err != nil {
			return nil, fmt.Errorf("content_generator: generate entry for topic %q: %w", topic.Title, err)
		}

		entries = append(entries, NewsletterEntry{
			TopicTitle:          topic.Title,
			Content:             aiEntry.Content,
			FilePathsReferenced: filePaths,
		})
	}

	return entries, nil
}
