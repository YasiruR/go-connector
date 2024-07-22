package catalog

import (
	"github.com/YasiruR/connector/core/protocols/dcat"
	"github.com/YasiruR/connector/core/protocols/odrl"
)

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

type Response struct {
	Context             string               `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	ID                  string               `json:"@id"`
	Type                string               `json:"@type" default:"dcat:Catalog"`
	DctTitle            string               `json:"dct:title"`
	DctDescription      []dcat.Description   `json:"dct:description"`
	DspaceParticipantID string               `json:"dspace:participantId"` // provider of the catalog
	DcatKeyword         []dcat.Keyword       `json:"dcat:keyword"`
	DcatService         []dcat.AccessService `json:"dcat:service"`
	DcatDataset         []dcat.Dataset       `json:"dcatDataset"`
}

type DatasetResponse struct {
	Context          string              `json:"@context" default:"https://w3id.org/dspace/2024/1/context.json"`
	ID               string              `json:"@id"`
	Type             string              `json:"@type" default:"dcat:Dataset"`
	DctTitle         string              `json:"dct:title"`
	DctDescription   []dcat.Description  `json:"dct:description"`
	DcatKeyword      []dcat.Keyword      `json:"dcat:keyword"`
	OdrlHasPolicy    []odrl.Offer        `json:"odrl:hasPolicy"`
	DcatDistribution []dcat.Distribution `json:"dcat:distribution"`
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
