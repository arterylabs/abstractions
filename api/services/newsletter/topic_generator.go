package newsletter

import (
	ai_client "abstractions/clients/ai"
	"abstractions/repositories"
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Topic represents a newsletter topic candidate returned by the generator.
type Topic struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	RelatedFiles []string `json:"relatedFiles"`
}

const (
	recentActivityWindow = 7 * 24 * time.Hour // files modified in the last 7 days
	recentNewsletterCap  = 4                  // number of past newsletters to de-duplicate against
	topicLimit           = 5                  // maximum topics returned
)

// TopicGenerator builds newsletter topic proposals from the indexed repository.
type TopicGenerator struct {
	aiProvider         ai_client.AIProvider
	fileIndexRepo      *repositories.RepositoryFileRepository
	chunkEmbeddingRepo *repositories.ChunkEmbeddingRepository
	newsletterRepo     *repositories.NewsletterRepository
	repositoryRepo     *repositories.RepositoryRepository
}

func NewTopicGenerator(
	aiProvider ai_client.AIProvider,
	fileIndexRepo *repositories.RepositoryFileRepository,
	chunkEmbeddingRepo *repositories.ChunkEmbeddingRepository,
	newsletterRepo *repositories.NewsletterRepository,
	repositoryRepo *repositories.RepositoryRepository,
) *TopicGenerator {
	return &TopicGenerator{
		aiProvider:         aiProvider,
		fileIndexRepo:      fileIndexRepo,
		chunkEmbeddingRepo: chunkEmbeddingRepo,
		newsletterRepo:     newsletterRepo,
		repositoryRepo:     repositoryRepo,
	}
}

// GenerateTopics produces up to 5 ranked newsletter topic proposals for the given repository.
func (g *TopicGenerator) GenerateTopics(ctx context.Context, repoID uuid.UUID) ([]Topic, error) {
	repo, err := g.repositoryRepo.GetByID(ctx, repoID)
	if err != nil || repo == nil {
		return nil, err
	}

	// Primary source: file-level chunk summaries (chunk_index = -1)
	fileSummaries := make(map[string]string)
	if g.chunkEmbeddingRepo != nil {
		chunkSummaries, err := g.chunkEmbeddingRepo.ListFileSummaries(ctx, repoID)
		if err == nil {
			for _, cs := range chunkSummaries {
				if cs.Summary != "" {
					fileSummaries[cs.RepoFilePath] = cs.Summary
				}
			}
		}
	}

	// Fallback: use domain tags from file index if no chunk summaries exist
	if len(fileSummaries) == 0 {
		files, err := g.fileIndexRepo.ListByRepoID(ctx, repoID)
		if err != nil {
			return nil, err
		}
		for _, fi := range files {
			if len(fi.Tags) > 0 {
				fileSummaries[fi.Path] = strings.Join(fi.Tags, ", ")
			}
		}
	}

	if len(fileSummaries) == 0 {
		return nil, nil
	}

	// Load last N newsletters to extract topics already covered
	recentTopicTitles := g.loadRecentTopicTitles(ctx, repoID)

	// Ask AI for topic proposals, passing previous topics to avoid repetition
	aiTopics, err := g.aiProvider.GenerateTopics(ctx, &ai_client.TopicGenerationRequest{
		RepoFullName:   repo.Name,
		FileSummaries:  fileSummaries,
		PreviousTopics: recentTopicTitles,
	})
	if err != nil {
		return nil, err
	}

	// Score, de-duplicate, and rank
	candidates := make([]Topic, 0, len(aiTopics))
	for _, t := range aiTopics {
		if t.Title == "" {
			continue
		}
		if isDuplicateTopic(t.Title, recentTopicTitles) {
			continue
		}
		candidates = append(candidates, Topic{
			Title:        t.Title,
			Description:  t.Description,
			RelatedFiles: t.RelatedFiles,
		})
	}

	if len(candidates) > topicLimit {
		candidates = candidates[:topicLimit]
	}
	return candidates, nil
}

// loadRecentTopicTitles returns normalised topic titles from the last N newsletters.
func (g *TopicGenerator) loadRecentTopicTitles(ctx context.Context, repoID uuid.UUID) []string {
	newsletters, err := g.newsletterRepo.ListByRepoID(ctx, repoID)
	if err != nil || len(newsletters) == 0 {
		return nil
	}
	if len(newsletters) > recentNewsletterCap {
		newsletters = newsletters[:recentNewsletterCap]
	}
	var titles []string
	for _, nl := range newsletters {
		if nl.Topic != nil && nl.Topic.Title != "" {
			titles = append(titles, strings.ToLower(nl.Topic.Title))
		}
	}
	return titles
}

// isDuplicateTopic returns true when title (lowercased) already appears in recentTitles.
func isDuplicateTopic(title string, recentTitles []string) bool {
	norm := strings.ToLower(title)
	for _, rt := range recentTitles {
		if strings.Contains(norm, rt) || strings.Contains(rt, norm) {
			return true
		}
	}
	return false
}
