package usecase

import (
	"context"
	"fmt"
	"leenwood/yandex-http/config"
	"leenwood/yandex-http/internal/domain/url"
	"leenwood/yandex-http/internal/domain/url/sqliteRepository"
	"leenwood/yandex-http/internal/usecase/dto"
)

type UrlUseCaseInterface interface {
	CreateShortUrl(request dto.CreateShortUrlUseCaseRequest) (dto.CreateShortUrlResponse, error)
	CreateShortUrlWithCustomId(request dto.CreateShortUrlWithCustomIdRequest) (dto.CreateShortUrlResponse, error)
}

type UrlUseCase struct {
	r url.RepositoryInterface
	c config.Config
}

func NewUrlUseCase(ctx context.Context, config config.Config) (*UrlUseCase, error) {
	repository, err := sqliteRepository.NewRepository(ctx, config.Database)
	if err != nil {
		return nil, err
	}
	return &UrlUseCase{r: repository}, nil
}

func (us *UrlUseCase) CreateShortUrl(request dto.CreateShortUrlUseCaseRequest) (dto.CreateShortUrlResponse, error) {
	model := dto.CreateShortUrlResponse{}
	// Проверяем, существует ли уже запись с таким OriginalUrl
	existingUrl, err := us.r.FindByUrl(request.Url)

	if err != nil {
		return model, err
	}

	if existingUrl != nil {
		returnUrl := fmt.Sprintf("%s:%s/%s", us.c.App.Hostname, us.c.App.Port, existingUrl.Id)
		model.Url = returnUrl
		model.ClickCount = existingUrl.ClickCount
		return model, nil
	}
	result, err := us.r.Save(request.Url, "")

	if err != nil {
		return model, err
	}

	returnUrl := fmt.Sprintf("%s:%s/%s", us.c.App.Hostname, us.c.App.Port, result.Id)
	model.Url = returnUrl
	model.ClickCount = result.ClickCount
	return model, nil
}

func (us *UrlUseCase) CreateShortUrlWithCustomId(request dto.CreateShortUrlWithCustomIdRequest) (dto.CreateShortUrlResponse, error) {
	model := dto.CreateShortUrlResponse{}
	// Проверяем, существует ли уже запись с таким OriginalUrl
	existingUrl, err := us.r.FindByUrl(request.Url)

	if err != nil {
		return model, err
	}

	if existingUrl != nil {
		returnUrl := fmt.Sprintf("%s:%s/%s", us.c.App.Hostname, us.c.App.Port, existingUrl.Id)
		model.Url = returnUrl
		model.ClickCount = existingUrl.ClickCount
		return model, nil
	}
	result, err := us.r.Save(request.Url, request.Id)

	if err != nil {
		return model, err
	}

	returnUrl := fmt.Sprintf("%s:%s/%s", us.c.App.Hostname, us.c.App.Port, result.Id)
	model.Url = returnUrl
	model.ClickCount = result.ClickCount
	return model, nil
}
