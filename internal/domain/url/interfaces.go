package url

type RepositoryInterface interface {
	FindById(id string) (*Url, error)
	FindByUrl(url string) (*Url, error)
	Save(model *Url) (*Url, error)
}
