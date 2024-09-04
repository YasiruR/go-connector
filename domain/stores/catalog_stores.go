package stores

import (
	"github.com/YasiruR/connector/domain/api/dsp/http/catalog"
	"github.com/YasiruR/connector/domain/models/dcat"
	"github.com/YasiruR/connector/domain/models/odrl"
)

/*
	Stores required to maintain a catalog either at provider or consumer side
	are defined here.
*/

// ProviderCatalog stores Datasets as per the DCAT profile recommended by IDSA.
// Current implementation supports only a single catalog per a provider.
type ProviderCatalog interface {
	Catalog() (dcat.Catalog, error)
	AddDataset(id string, val dcat.Dataset)
	Dataset(id string) (dcat.Dataset, error)
}

// ConsumerCatalog stores catalogs received by providers and therefore, it may
// include multiple catalogs as opposed to ProviderCatalog.
type ConsumerCatalog interface {
	AddCatalog(res catalog.Response)
	Catalog(providerId string) (catalog.Response, error)
	Offer(offerId string) (ofr odrl.Offer, err error)
	AllCatalogs() ([]catalog.Response, error)
}
