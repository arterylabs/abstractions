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

// UpdateNewsletterTopicRequest is the body for PUT /repositories/:repo_id/newsletter-topics/:topic_id.
type UpdateNewsletterTopicRequest struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	RelatedFiles []string `json:"relatedFiles"`
}

func (r *UpdateNewsletterTopicRequest) Validate() string {
	if r.Title == "" {
		return "title is required"
	}
	if len(r.Title) > 200 {
		return "title must be 200 characters or fewer"
	}
	return ""
}

// ReorderNewsletterTopicRequest is the body for PUT /repositories/:repo_id/newsletter-topics/:topic_id/reorder.
type ReorderNewsletterTopicRequest struct {
	Direction string `json:"direction"`
}

func (r *ReorderNewsletterTopicRequest) Validate() string {
	switch r.Direction {
	case "up", "down", "top":
		return ""
	default:
		return "direction must be one of: up, down, top"
	}
}
