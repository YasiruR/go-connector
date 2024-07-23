package catalog

type Provider interface {
	HandleCatalogRequest(filter any) (Response, error)
	HandleDatasetRequest(id string) (DatasetResponse, error)
}

type Consumer interface {
	RequestCatalog(endpoint string) (Response, error) // endpoint should be generic
	RequestDataset()
}
