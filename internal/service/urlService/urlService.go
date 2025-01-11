package urlservice

import (
	"context"
	"leenwood/yandex-http/internal/domain/url"
	"leenwood/yandex-http/internal/domain/url/postgresRepository"
	"time"
)

type UrlService struct {
	r url.RepositoryInterface
}

func NewService(ctx context.Context) (*UrlService, error) {
	r, err := postgresRepository.NewRepository(ctx, "")
	if err != nil {
		return nil, err
	}
	return &UrlService{
		r: r,
	}, nil
}

func (s *UrlService) Create(OriginalUrl string) (*url.Url, error) {

	// Проверяем, существует ли уже запись с таким OriginalUrl
	existingUrl, err := s.r.FindByUrl(OriginalUrl)

	if err != nil {
		return nil, err
	}

	if existingUrl != nil {
		return existingUrl, nil
	}

	model := &url.Url{}
	model.OriginalUrl = OriginalUrl
	model.Date = time.Now()
	model.ClickCount = 0
	result, err := s.r.Save(model)

	if err != nil {
		return nil, err
	}

	return result, nil
}
