package usecase

import (
	"context"
	"fmt"
	"leenwood/yandex-http/config"
	urlservice "leenwood/yandex-http/internal/service/urlService"
	"leenwood/yandex-http/internal/usecase/dto"
)

type UrlUseCaseInterface interface {
	CreateShortUrl(request dto.CreateShortUrlRequest) (dto.CreateShortUrlResponse, error)
	CreateShortUrlWithCustomId(request dto.CreateShortUrlWithCustomIdRequest) (dto.CreateShortUrlResponse, error)
}

type UrlUseCase struct {
	s urlservice.URLService
	c config.Config
}

func NewUrlUseCase(ctx context.Context, config config.Config) (*UrlUseCase, error) {
	service, err := urlservice.NewService(ctx)
	if err != nil {
		return nil, err
	}
	return &UrlUseCase{s: service}, nil
}

func (us *UrlUseCase) CreateShortUrl(request dto.CreateShortUrlRequest) (dto.CreateShortUrlResponse, error) {
	model := dto.CreateShortUrlResponse{}

	result, err := us.s.Create(request.Url)
	if err != nil {
		return model, err
	}

	url := fmt.Sprintf("%s:%s/%s", us.c.App.Hostname, us.c.App.Port, result.Id)
	model.Url = url
	model.ClickCount = result.ClickCount
	return model, nil
}

func (us *UrlUseCase) CreateShortUrlWithCustomId(request dto.CreateShortUrlWithCustomIdRequest) (dto.CreateShortUrlResponse, error) {
	model := dto.CreateShortUrlResponse{}

	result, err := us.s.CreateWithCustomId(request.Url, request.Id)
	if err != nil {
		return model, err
	}

	url := fmt.Sprintf("%s:%s/%s", us.c.App.Hostname, us.c.App.Port, result.Id)
	model.Url = url
	model.ClickCount = result.ClickCount
	return model, nil
}
