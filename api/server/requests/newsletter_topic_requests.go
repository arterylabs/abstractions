package requests

type CreateNewsletterTopicRequest struct {
	Title        string   `json:"title" validate:"required,max=200"`
	Description  string   `json:"description"`
	RelatedFiles []string `json:"relatedFiles"`
}
