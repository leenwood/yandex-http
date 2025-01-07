package urlservice

import (
	"crypto/rand"
	"leenwood/yandex-http/internal/domain/url"
	"leenwood/yandex-http/internal/domain/url/repository"
	"time"
)

type UrlService struct {
	r url.RepositoryInterface
}

func NewService() (*UrlService, error) {
	r := repository.NewRepository()
	return &UrlService{
		r: r,
	}, nil
}

func (s *UrlService) Create(OriginalUrl string) url.Url {
	// TODO Добавить проверку на наличие уникальной ссылки
	uuid := make([]byte, 6)
	rand.Read(uuid)
	// TODO Добавить проверку на занятость юида
	var model url.Url
	model.ShortUrl = string(uuid)
	model.OriginalUrl = OriginalUrl
	model.Date = time.Now()
	result := s.r.Save(model)
	if result {
		return model
	} else {
		return url.Url{}
	}
}
