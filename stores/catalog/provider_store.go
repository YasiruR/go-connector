package catalog

import (
	"fmt"
	"github.com/YasiruR/go-connector/domain"
	"github.com/YasiruR/go-connector/domain/boot"
	"github.com/YasiruR/go-connector/domain/core"
	"github.com/YasiruR/go-connector/domain/errors"
	"github.com/YasiruR/go-connector/domain/models/dcat"
	"github.com/YasiruR/go-connector/domain/pkg"
	"github.com/YasiruR/go-connector/domain/stores"
)

const collProviderCatalog = `provider-catalog`

// ProviderCatalog stores Datasets and Data Services which can be shared through a connector
type ProviderCatalog struct {
	meta dcat.CatalogMetadata
	urn  pkg.URNService
	coll pkg.Collection
}

func NewProviderCatalog(cfg boot.Config, plugins domain.Plugins) *ProviderCatalog {
	c := &ProviderCatalog{
		urn:  plugins.URNService,
		coll: plugins.Store.NewCollection(),
	}

	if err := c.init(cfg); err != nil {
		plugins.Log.Fatal(fmt.Sprintf("init catalog store failed: %v", err))
	}

	plugins.Log.Info(fmt.Sprintf("initialized %s store", collProviderCatalog), `catalog ID: `+c.meta.ID)
	return c
}

// may not need to init if only a consumer
func (p *ProviderCatalog) init(cfg boot.Config) error {
	catId, err := p.urn.NewURN()
	if err != nil {
		return errors.PkgError(pkg.TypeURN, `NewURN`, err, `catalog id`)
	}

	var kws []dcat.Keyword
	for _, key := range cfg.Catalog.Keywords {
		kws = append(kws, dcat.Keyword(key))
	}

	var descs []dcat.Description
	for _, desc := range cfg.Catalog.Descriptions {
		descs = append(descs, dcat.Description{Value: desc, Language: dcat.LanguageEnglish})
	}

	var svcs []dcat.AccessService
	for _, e := range cfg.Catalog.AccessServices {
		svcId, err := p.urn.NewURN()
		if err != nil {
			return errors.PkgError(pkg.TypeURN, `NewURN`, err, `service id`)
		}

		svcs = append(svcs, dcat.AccessService{
			ID:                  svcId,
			Type:                dcat.TypeDataService,
			EndpointURL:         e,
			EndpointDescription: core.ServiceConnector, // should be considered in later versions
		})
	}

	p.meta = dcat.CatalogMetadata{
		ID:             catId,
		Type:           dcat.TypeCatalog,
		DctTitle:       cfg.Catalog.Title,
		DctDescription: descs,
		DcatKeyword:    kws,
		DcatService:    svcs,
	}

	return nil
}

func (p *ProviderCatalog) Catalog() (dcat.Catalog, error) {
	vals, err := p.coll.GetAll()
	if err != nil {
		return dcat.Catalog{}, stores.QueryFailed(collProviderCatalog, `GetAll`, err)
	}

	var cat dcat.Catalog
	cat.CatalogMetadata = p.meta

	for _, val := range vals {
		cat.DcatDataset = append(cat.DcatDataset, val.(dcat.Dataset))
	}

	return cat, nil
}

func (p *ProviderCatalog) AddDataset(id string, val dcat.Dataset) {
	_ = p.coll.Set(id, val)
}

func (p *ProviderCatalog) Dataset(id string) (dcat.Dataset, error) {
	val, err := p.coll.Get(id)
	if err != nil {
		return dcat.Dataset{}, stores.QueryFailed(collProviderCatalog, `Get`, err)
	}

	if val == nil {
		return dcat.Dataset{}, stores.InvalidKey(id)
	}

	return val.(dcat.Dataset), nil
}

func (p *ProviderCatalog) DatasetByOfferId(offerId string) (dcat.Dataset, error) {
	vals, err := p.coll.GetAll()
	if err != nil {
		return dcat.Dataset{}, stores.QueryFailed(collProviderCatalog, `GetAll`, err)
	}

	for _, val := range vals {
		ds, ok := val.(dcat.Dataset)
		if ok {
			for _, offer := range ds.OdrlHasPolicy {
				if offer.Id == offerId {
					return ds, nil
				}
			}
		}
	}

	return dcat.Dataset{}, stores.InvalidKey(offerId)
}
