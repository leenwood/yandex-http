package urlservice

import "leenwood/yandex-http/internal/domain/url"

type URLService interface {
	Create(url string) (*url.Url, error)
	CreateWithCustomId(url, id string) (*url.Url, error)
}
