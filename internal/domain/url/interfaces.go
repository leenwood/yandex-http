package url

type RepositoryInterface interface {
	FindById(id int) Url
	FindByGuid(guid string) Url
	Save(model Url) bool
}
