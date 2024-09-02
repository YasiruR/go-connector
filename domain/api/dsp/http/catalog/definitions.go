package catalog

import "github.com/YasiruR/connector/domain/api"

// Message types
const (
	MsgTypRequest        = `dspace:CatalogRequestMessage`
	MsgTypDatasetRequest = `dspace:DatasetRequestMessage`
)

// Endpoints
const (
	RequestEndpoint        = `/catalog/request`
	RequestDatasetEndpoint = `/catalog/datasets{` + api.ParamPid + `}`
)
