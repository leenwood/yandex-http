package handlers

import (
	"context"
	"encoding/json"
	"leenwood/yandex-http/config"
	"leenwood/yandex-http/internal/usecase"
	"leenwood/yandex-http/internal/usecase/dto"
	"net/http"
)

type UrlHandler struct {
	us usecase.UrlUseCaseInterface
}

func NewUrlHandler(ctx context.Context, cfg config.Config) (*UrlHandler, error) {
	us, err := usecase.NewUrlUseCase(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &UrlHandler{us: us}, nil
}

func (uh *UrlHandler) CreateShortUrl(res http.ResponseWriter, req *http.Request) {
	var urlShort, id string

	err := req.ParseForm()
	if err != nil {
		http.Error(res, "failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}
	for k, v := range req.Form {
		switch k {
		case "url":
			urlShort = v[0]
		case "id":
			id = v[0]
		}
	}
	if urlShort == "" {
		http.Error(res, "url parameter is required", http.StatusBadRequest)
		return
	}

	var response dto.CreateShortUrlResponse

	if id != "" {
		request := dto.CreateShortUrlWithCustomIdRequest{
			Url: urlShort,
			Id:  id,
		}
		response, err = uh.us.CreateShortUrlWithCustomId(request)
	} else {
		request := dto.CreateShortUrlRequest{Url: urlShort}
		response, err = uh.us.CreateShortUrl(request)
	}

	if err != nil {
		http.Error(res, "failed to create short URL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Кодируем и отправляем ответ
	res.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(res).Encode(response); err != nil {
		http.Error(res, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}
