package responses

type AIModelResponse struct {
	ID                     string  `json:"id"`
	DisplayName            string  `json:"displayName"`
	Provider               string  `json:"provider"`
	SupportsTextGeneration bool    `json:"supportsTextGeneration"`
	SupportsEmbeddings     bool    `json:"supportsEmbeddings"`
	PricePerMillionTokens  float64 `json:"pricePerMillionTokens"`
	ContextWindow          int     `json:"contextWindow"`
}
