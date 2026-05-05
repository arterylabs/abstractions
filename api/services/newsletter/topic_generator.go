package newsletter

import (
	ai_client "abstractions/clients/ai"
	"abstractions/db/models"
	"abstractions/repositories"
	"context"
	"encoding/json"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Topic represents a newsletter topic candidate returned by the generator.
type Topic struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	RelatedFiles []string `json:"relatedFiles"`
	RecencyScore float64  `json:"recencyScore"`
}

const (
	recentActivityWindow = 7 * 24 * time.Hour // files modified in the last 7 days
	recentNewsletterCap  = 4                  // number of past newsletters to de-duplicate against
	topicLimit           = 5                  // maximum topics returned
)

// TopicGenerator builds newsletter topic proposals from the indexed repository.
type TopicGenerator struct {
	aiProvider     ai_client.AIProvider
	fileIndexRepo  *repositories.FileIndexRepository
	newsletterRepo *repositories.NewsletterRepository
	repositoryRepo *repositories.RepositoryRepository
}

func NewTopicGenerator(
	aiProvider ai_client.AIProvider,
	fileIndexRepo *repositories.FileIndexRepository,
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
		if fi.LastProcessedAt != nil && fi.LastProcessedAt.After(cutoff) {
			recentFiles[fi.FilePath] = true
		}
		if len(fi.DomainTags) > 0 {
			fileSummaries[fi.FilePath] = strings.Join(fi.DomainTags, ", ")
		}
	}

	if len(fileSummaries) == 0 {
		return nil, nil
	}

	// Load last N newsletters to extract topics already covered
	recentTopicTitles := g.loadRecentTopicTitles(ctx, repoID)

	// Ask AI for topic proposals
	aiTopics, err := g.aiProvider.GenerateTopics(ctx, &ai_client.TopicGenerationRequest{
		RepoFullName:  repo.FullName,
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
		// Skip topic if it closely matches a recently sent newsletter topic
		if isDuplicateTopic(t.Title, recentTopicTitles) {
			continue
		}
		score := computeRecencyScore(t.RelatedFiles, recentFiles)
		candidates = append(candidates, Topic{
			Title:        t.Title,
			Description:  t.Description,
			RelatedFiles: t.RelatedFiles,
			RecencyScore: score,
		})
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].RecencyScore > candidates[j].RecencyScore
	})

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
	// Most-recent first; cap at recentNewsletterCap
	if len(newsletters) > recentNewsletterCap {
		newsletters = newsletters[:recentNewsletterCap]
	}
	var titles []string
	for _, nl := range newsletters {
		for _, t := range extractTopicTitles(nl) {
			titles = append(titles, strings.ToLower(t))
		}
	}
	return titles
}

// extractTopicTitles deserialises the Topics JSONBRaw field of a newsletter row.
func extractTopicTitles(nl *models.Newsletter) []string {
	if nl.Topics == nil {
		return nil
	}
	var items []struct {
		Title string `json:"title"`
	}
	if err := json.Unmarshal(nl.Topics, &items); err != nil {
		return nil
	}
	titles := make([]string, 0, len(items))
	for _, item := range items {
		if item.Title != "" {
			titles = append(titles, item.Title)
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

// computeRecencyScore returns the fraction of a topic's related files that were
// recently modified (last 7 days).  Returns 0 if the topic has no related files.
func computeRecencyScore(relatedFiles []string, recentFiles map[string]bool) float64 {
	if len(relatedFiles) == 0 {
		return 0
	}
	var count int
	for _, f := range relatedFiles {
		if recentFiles[f] {
			count++
		}
	}
	return float64(count) / float64(len(relatedFiles))
}
