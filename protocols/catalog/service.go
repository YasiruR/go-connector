package catalog

type Service interface {
	GetCatalog(filter any) error
	GetDataset(id string)
}
