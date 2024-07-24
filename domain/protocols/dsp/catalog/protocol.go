package catalog

type Handler interface {
	HandleCatalogRequest(filter any) (Response, error)
	HandleDatasetRequest(id string) (DatasetResponse, error)
}

type Controller interface {
	RequestCatalog(endpoint string) (Response, error) // endpoint should be generic
	RequestDataset(id, endpoint string) (DatasetResponse, error)
}
