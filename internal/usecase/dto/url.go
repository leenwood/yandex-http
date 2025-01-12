package dto

type CreateShortUrlRequest struct {
	Url string `json:"url"`
}

type CreateShortUrlWithCustomIdRequest struct {
	Url string `json:"url"`
	Id  string `json:"id"`
}

type CreateShortUrlResponse struct {
	Url        string `json:"url"`
	ClickCount uint64 `json:"click_count"`
}
