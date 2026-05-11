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
	aiProvider     ai_client.AIProvider
	fileIndexRepo  *repositories.RepositoryFileRepository
	newsletterRepo *repositories.NewsletterRepository
	repositoryRepo *repositories.RepositoryRepository
}

func NewTopicGenerator(
	aiProvider ai_client.AIProvider,
	fileIndexRepo *repositories.RepositoryFileRepository,
	newsletterRepo *repositories.NewsletterRepository,
	repositoryRepo *repositories.RepositoryRepository,
) *TopicGenerator {
	return &TopicGenerator{
		aiProvider:     aiProvider,
		fileIndexRepo:  fileIndexRepo,
		newsletterRepo: newsletterRepo,
		repositoryRepo: repositoryRepo,
	}
}

// GenerateTopics produces up to 5 ranked newsletter topic proposals for the given repository.
func (g *TopicGenerator) GenerateTopics(ctx context.Context, repoID uuid.UUID) ([]Topic, error) {
	repo, err := g.repositoryRepo.GetByID(ctx, repoID)
	if err != nil || repo == nil {
		return nil, err
	}

	// Load all indexed files
	files, err := g.fileIndexRepo.ListByRepoID(ctx, repoID)
	if err != nil {
		return nil, err
	}

	// Build file summaries map (domain tags joined) and track recently active files
	cutoff := time.Now().Add(-recentActivityWindow)
	recentFiles := map[string]bool{}
	fileSummaries := make(map[string]string, len(files))

	for _, fi := range files {
		if fi.ProcessedAt != nil && fi.ProcessedAt.After(cutoff) {
			recentFiles[fi.Path] = true
		}
		if len(fi.Tags) > 0 {
			fileSummaries[fi.Path] = strings.Join(fi.Tags, ", ")
		}
	}

	if len(fileSummaries) == 0 {
		return nil, nil
	}

	// Load last N newsletters to extract topics already covered
	recentTopicTitles := g.loadRecentTopicTitles(ctx, repoID)

	// Ask AI for topic proposals
	aiTopics, err := g.aiProvider.GenerateTopics(ctx, &ai_client.TopicGenerationRequest{
		RepoFullName:  repo.Name,
		FileSummaries: fileSummaries,
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
// Each newsletter stores its topic via the newsletter_topics relation (Topic field).
func (g *TopicGenerator) loadRecentTopicTitles(ctx context.Context, repoID uuid.UUID) []string {
	newsletters, err := g.newsletterRepo.ListByRepoID(ctx, repoID)
	if err != nil || len(newsletters) == 0 {
		return nil
	}
	// Most-recent first; cap at recentNewsletterCap
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
