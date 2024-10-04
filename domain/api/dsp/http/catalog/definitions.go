package catalog

import "github.com/YasiruR/go-connector/domain/api"

// Message types
const (
	MsgTypRequest        = `dspace:CatalogRequestMessage`
	MsgTypDatasetRequest = `dspace:DatasetRequestMessage`
	MsgTypeError         = `dspace:CatalogError`
)

// Endpoints
const (
	RequestEndpoint        = `/catalog/request`
	RequestDatasetEndpoint = `/catalog/datasets/{` + api.ParamPid + `}`
)
