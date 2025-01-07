package repository

import "leenwood/yandex-http/internal/domain/url"

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) FindById(id int) url.Url {
	var model url.Url

	return model
}

func (r *Repository) FindByGuid(guid string) url.Url {
	var model url.Url

	return model
}

func (r *Repository) Save(model url.Url) bool {
	return true
}
