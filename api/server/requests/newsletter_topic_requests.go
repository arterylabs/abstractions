package requests

// CreateNewsletterTopicRequest is the body for POST /repositories/:repo_id/newsletter-topics.
type CreateNewsletterTopicRequest struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	RelatedFiles []string `json:"relatedFiles"`
}

func (r *CreateNewsletterTopicRequest) Validate() string {
	if r.Title == "" {
		return "title is required"
	}
	if len(r.Title) > 200 {
		return "title must be 200 characters or fewer"
	}
	return ""
}
