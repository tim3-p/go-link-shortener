package models

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

type UserUrl struct {
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}

type UserUrlsResponse struct {
	UserUrls []UserUrl
}
