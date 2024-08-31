package catalog

import (
	"fmt"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/boot"
	"github.com/YasiruR/connector/domain/core"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/models/dcat"
	"github.com/YasiruR/connector/domain/pkg"
)

// ProviderCatalog stores Datasets and Data Services which can be shared through a connector
type ProviderCatalog struct {
	meta  dcat.CatalogMetadata
	urn   pkg.URNService
	store pkg.Collection
}

func NewProviderCatalog(cfg boot.Config, plugins domain.Plugins) *ProviderCatalog {
	c := &ProviderCatalog{
		urn:   plugins.URNService,
		store: plugins.Database.NewCollection(),
	}

	if err := c.init(cfg); err != nil {
		plugins.Log.Fatal(fmt.Sprintf("init catalog store failed: %v", err))
	}

	plugins.Log.Info("initialized catalog store", "catalog id: "+c.meta.ID)
	return c
}

// may not need to init if only a consumer
func (p *ProviderCatalog) init(cfg boot.Config) error {
	catId, err := p.urn.NewURN()
	if err != nil {
		return errors.PkgFailed(pkg.TypeURN, `New`, err)
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
			return errors.PkgFailed(pkg.TypeURN, `New`, err)
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
	vals, err := p.store.GetAll()
	if err != nil {
		return dcat.Catalog{}, errors.QueryFailed(`dataset`, `GetAll`, err)
	}

	var cat dcat.Catalog
	cat.CatalogMetadata = p.meta

	for _, val := range vals {
		cat.DcatDataset = append(cat.DcatDataset, val.(dcat.Dataset))
		fmt.Println("FETCHED DS: ", val.(dcat.Dataset))
	}

	return cat, nil
}

func (p *ProviderCatalog) AddDataset(id string, val dcat.Dataset) {
	_ = p.store.Set(id, val)
}

func (p *ProviderCatalog) Dataset(id string) (dcat.Dataset, error) {
	val, err := p.store.Get(id)
	if err != nil {
		return dcat.Dataset{}, errors.QueryFailed(`dataset`, `Get`, err)
	}
	return val.(dcat.Dataset), nil
}
