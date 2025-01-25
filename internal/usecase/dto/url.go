package dto

import "time"

type PaginationRequest struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

type UrlInfoResponse struct {
	Id          string    `json:"id"`
	OriginalUrl string    `json:"original_url"`
	ShortUrl    string    `json:"short_url"`
	CountClick  uint64    `json:"count_click"`
	CreatedDate time.Time `json:"created_date"`
}

type UrlClickRequest struct {
	Id string
}

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
