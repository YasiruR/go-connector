package catalog

import (
	"github.com/YasiruR/connector/domain/models/dcat"
)

type Response struct {
	Context             string `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	DspaceParticipantID string `json:"dspace:participantId"` // provider of the catalog
	dcat.Catalog
}

type DatasetResponse struct {
	Context string `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	dcat.Dataset
}

type ErrorResponse struct {
	Context      string `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	Type         string `json:"@type" default:"dspace:CatalogError"`
	DspaceCode   string `json:"dspace:code"`
	DspaceReason []struct {
		Value    string `json:"@value"`
		Language string `json:"@language"`
	} `json:"dspace:reason"`
}
