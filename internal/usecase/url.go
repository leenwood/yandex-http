package usecase

import (
	"context"
	"fmt"
	"leenwood/yandex-http/config"
	"leenwood/yandex-http/internal/domain/url"
	"leenwood/yandex-http/internal/domain/url/postgresRepository"
	"leenwood/yandex-http/internal/usecase/dto"
	"strings"
)

type UrlUseCaseInterface interface {
	CreateShortUrl(request dto.CreateShortUrlUseCaseRequest) (dto.CreateShortUrlResponse, error)
	CreateShortUrlWithCustomId(request dto.CreateShortUrlWithCustomIdRequest) (dto.CreateShortUrlResponse, error)
	GetUrlList(pagination dto.PaginationRequest) ([]dto.UrlInfoResponse, error)
	ClickUrl(request dto.UrlClickRequest) (string, error)
}

type UrlUseCase struct {
	r url.RepositoryInterface
	c config.Config
}

func NewUrlUseCase(ctx context.Context, config config.Config) (*UrlUseCase, error) {
	repository, err := postgresRepository.NewRepository(ctx, config.Database)
	if err != nil {
		return nil, err
	}
	return &UrlUseCase{r: repository, c: config}, nil
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

func (us *UrlUseCase) GetUrlList(pagination dto.PaginationRequest) ([]dto.UrlInfoResponse, error) {
	urlsRepositoryInfo, err := us.r.FindAll(pagination.Page, pagination.Limit)
	if err != nil {
		return nil, err
	}

	return us.transformSliceToUrlInfo(urlsRepositoryInfo), nil
}

func (us *UrlUseCase) ClickUrl(request dto.UrlClickRequest) (string, error) {
	urlRepository, err := us.r.FindById(request.Id)
	if err != nil {
		return "", err
	}
	urlRepository.ClickCount++

	urlRepository, err = us.r.Update(urlRepository)

	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(urlRepository.OriginalUrl, "http://") && !strings.HasPrefix(urlRepository.OriginalUrl, "https://") {
		urlRepository.OriginalUrl = "https://" + urlRepository.OriginalUrl
	}

	return urlRepository.OriginalUrl, nil
}

func (us *UrlUseCase) transformSliceToUrlInfo(urls []*url.Url) []dto.UrlInfoResponse {
	var result []dto.UrlInfoResponse
	for i := range urls {
		result = append(result, us.transformToUrlInfo(urls[i]))
	}
	return result
}

func (us *UrlUseCase) transformToUrlInfo(repositoryUrl *url.Url) dto.UrlInfoResponse {
	return dto.UrlInfoResponse{
		Id:          repositoryUrl.Id,
		OriginalUrl: repositoryUrl.OriginalUrl,
		ShortUrl:    fmt.Sprintf("%s:%s/%s", us.c.App.Hostname, us.c.App.Port, repositoryUrl.Id),
		CountClick:  repositoryUrl.ClickCount,
		CreatedDate: repositoryUrl.CreatedDate,
	}
}
