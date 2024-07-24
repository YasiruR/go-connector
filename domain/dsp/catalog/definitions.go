package catalog

//const ParamDatasetId = `datasetId`

// endpoints
const (
	RequestEndpoint = `/catalog/request`
	//RequestDatasetEndpoint = `/catalog/datasets/{` + ParamDatasetId + `}`
	RequestDatasetEndpoint = `/catalog/datasets`
)

const (
	TypeCatalogRequest = `dspace:CatalogRequestMessage`
	TypeDatasetRequest = `dspace:DatasetRequestMessage`
)
