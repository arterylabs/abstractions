package requests

type UpdateNewsletterRequest struct {
	HTMLContent string `json:"htmlContent" validate:"required"`
}
