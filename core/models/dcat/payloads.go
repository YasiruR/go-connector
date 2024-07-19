package dcat

import "github.com/YasiruR/connector/core/models/odrl"

// namespace prefix reference: https://www.w3.org/TR/vocab-dcat-2/#normative-namespaces

type Keyword string

type Dataset struct {
	ID               string         `json:"@id"`
	Type             string         `json:"@type" default:"dcat:Dataset"`
	DctTitle         string         `json:"dct:title"`
	DctDescription   []Description  `json:"dct:description"`
	DcatKeyword      []Keyword      `json:"dcat:keyword"`
	OdrlHasPolicy    []odrl.Offer   `json:"odrl:hasPolicy"`
	DcatDistribution []Distribution `json:"dcat:distribution"`
}

type Description struct {
	Value    string `json:"@value"`
	Language string `json:"@language"`
}

type Distribution struct {
	Type              string          `json:"@type"`
	DctFormat         string          `json:"dct:format"`
	DcatAccessService []AccessService `json:"dcat:accessService"`
}

type AccessService struct {
	ID              string `json:"@id"`
	Type            string `json:"@type"`
	DcatEndpointURL string `json:"dcat:endpointURL"`
}
