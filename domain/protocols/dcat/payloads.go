package dcat

import "github.com/YasiruR/connector/domain/protocols/odrl"

// namespace prefix reference: https://www.w3.org/TR/vocab-dcat-2/#normative-namespaces

const (
	TypeCatalog      = `dcat:Catalog`
	TypeDataset      = `dcat:Dataset`
	TypeDistribution = `dcat:Distribution`
	TypeDataService  = `dcat:DataService`
	LanguageEnglish  = `en`
)

type Keyword string

type Catalog struct {
	CatalogMetadata
	DcatDataset []Dataset `json:"dcat:dataset"`
}

type CatalogMetadata struct {
	ID             string          `json:"@id"`
	Type           string          `json:"@type" default:"dcat:Catalog"`
	DctTitle       string          `json:"dct:title"`
	DctDescription []Description   `json:"dct:description"`
	DcatKeyword    []Keyword       `json:"dcat:keyword"`
	DcatService    []AccessService `json:"dcat:service"`
}

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
	ID                  string `json:"@id"`
	Type                string `json:"@type"`
	EndpointURL         string `json:"dcat:endpointURL"`
	EndpointDescription string `json:"dcat:endpointDescription"`
}
