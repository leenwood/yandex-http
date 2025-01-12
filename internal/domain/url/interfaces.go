package url

type RepositoryInterface interface {
	FindById(id string) (*Url, error)
	FindByUrl(url string) (*Url, error)
	Save(originalUrl string, shortUuid string) (*Url, error)
}
