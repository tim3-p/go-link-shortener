package models

// Struct for Task channel
type Task struct {
	URLs   []string
	UserID string
}

// Request for new shorten URL
type ShortenRequest struct {
	URL string `json:"url"`
}

// Response for new shorten URL
type ShortenResponse struct {
	Result string `json:"result"`
}

// Response for user URLs
type UserURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// Request for batch operation
type ShortenBatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// Response for batch operation
type ShortenBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
