package catalog

// Data models required for the Catalog Protocol as specified by
// https://docs.internationaldataspaces.org/ids-knowledgebase/v/dataspace-protocol/catalog/catalog.protocol

type Request struct {
	Context      string   `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type         string   `json:"@type" default:"dspace:CatalogRequestMessage"`
	DspaceFilter []string `json:"dspace:filter"` // optional and implementation-specific
}

type DatasetRequest struct {
	Context   string `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type      string `json:"@type" default:"dspace:DatasetRequestMessage"`
	DatasetId string `json:"dspace:dataset"`
}
