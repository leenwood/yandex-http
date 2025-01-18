package dto

type CreateShortUrlRequest struct {
	Url string `form:"url" binding:"required"`
	Id  string `form:"id"`
}

type CreateShortUrlUseCaseRequest struct {
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
