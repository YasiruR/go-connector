package stores

import (
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/boot"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/models/dcat"
	"github.com/YasiruR/connector/domain/pkg"
)

// Catalog stores Datasets and Data Services which can be shared through a connector
type Catalog struct {
	meta  dcat.CatalogMetadata
	urn   pkg.URNService
	store pkg.Collection
}

func NewCatalog(plugins domain.Plugins) *Catalog {
	plugins.Log.Info("initialized catalog store")
	return &Catalog{
		urn:   plugins.URNService,
		store: plugins.Database.NewCollection(),
	}
}

func (c *Catalog) Init(cfg boot.Config) error {
	catId, err := c.urn.NewURN()
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
		svcId, err := c.urn.NewURN()
		if err != nil {
			return errors.PkgFailed(pkg.TypeURN, `New`, err)
		}

		svcs = append(svcs, dcat.AccessService{
			ID:                  svcId,
			Type:                dcat.TypeDataService,
			EndpointURL:         e,
			EndpointDescription: "", // should be considered in later versions
		})
	}

	c.meta = dcat.CatalogMetadata{
		ID:             catId,
		Type:           dcat.TypeCatalog,
		DctTitle:       cfg.Catalog.Title,
		DctDescription: descs,
		DcatKeyword:    kws,
		DcatService:    svcs,
	}

	return nil
}

func (c *Catalog) Get() (dcat.Catalog, error) {
	vals, err := c.store.GetAll()
	if err != nil {
		return dcat.Catalog{}, errors.QueryFailed(`dataset`, `GetAll`, err)
	}

	var cat dcat.Catalog
	cat.CatalogMetadata = c.meta

	for _, val := range vals {
		cat.DcatDataset = append(cat.DcatDataset, val.(dcat.Dataset))
	}

	return cat, nil
}

func (c *Catalog) AddDataset(id string, val dcat.Dataset) {
	_ = c.store.Set(id, val)
}

func (c *Catalog) Dataset(id string) (dcat.Dataset, error) {
	val, err := c.store.Get(id)
	if err != nil {
		return dcat.Dataset{}, errors.QueryFailed(`dataset`, `Get`, err)
	}
	return val.(dcat.Dataset), nil
}
