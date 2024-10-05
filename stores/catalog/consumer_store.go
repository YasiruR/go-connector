package catalog

import (
	"github.com/YasiruR/go-connector/domain"
	"github.com/YasiruR/go-connector/domain/api/dsp/http/catalog"
	"github.com/YasiruR/go-connector/domain/errors"
	"github.com/YasiruR/go-connector/domain/models/odrl"
	"github.com/YasiruR/go-connector/domain/pkg"
	"github.com/YasiruR/go-connector/domain/stores"
)

const (
	collConsumerCatalog = `consumer-catalog`
	collOfferedData     = `offer-dataset`
)

type ConsumerCatalog struct {
	urn      pkg.URNService
	coll     pkg.Collection
	ofrdData pkg.Collection // key: offer id, value: dataset id
}

func NewConsumerCatalog(plugins domain.Plugins) *ConsumerCatalog {
	return &ConsumerCatalog{
		urn:      plugins.URNService,
		coll:     plugins.Store.NewCollection(),
		ofrdData: plugins.Store.NewCollection(),
	}
}

func (c *ConsumerCatalog) AddCatalog(res catalog.Response) {
	_ = c.coll.Set(res.DspaceParticipantID, res)
	for _, ds := range res.DcatDataset {
		for _, ofr := range ds.OdrlHasPolicy {
			_ = c.ofrdData.Set(ofr.Id, ds.ID)
		}
	}
}

func (c *ConsumerCatalog) Catalog(providerId string) (catalog.Response, error) {
	val, err := c.coll.Get(providerId)
	if err != nil {
		return catalog.Response{}, stores.QueryFailed(collConsumerCatalog, `Get`, err)
	}

	if val == nil {
		return catalog.Response{}, stores.InvalidKey(providerId)
	}

	return val.(catalog.Response), nil
}

func (c *ConsumerCatalog) Offer(offerId string) (ofr odrl.Offer, err error) {
	// should optimize this function
	// - expensive AllCatalogs function
	// - nested loops
	cats, err := c.AllCatalogs()
	if err != nil {
		return odrl.Offer{}, errors.StoreFailed(collConsumerCatalog, `AllCatalogs`, err)
	}

	for _, cat := range cats {
		for _, ds := range cat.DcatDataset {
			for _, ofr = range ds.OdrlHasPolicy {
				if offerId == ofr.Id {
					ofr.Target = odrl.Target(ds.ID)
					return ofr, nil
				}
			}
		}
	}

	return odrl.Offer{}, stores.InvalidKey(offerId)
}

func (c *ConsumerCatalog) AllCatalogs() ([]catalog.Response, error) {
	vals, err := c.coll.GetAll()
	if err != nil {
		return nil, stores.QueryFailed(collConsumerCatalog, `GetAll`, err)
	}

	res := make([]catalog.Response, len(vals))
	for _, val := range vals {
		res = append(res, val.(catalog.Response))
	}

	return res, nil
}

func (c *ConsumerCatalog) DatasetByOfferId(offerId string) (id string, err error) {
	val, err := c.ofrdData.Get(offerId)
	if err != nil {
		return ``, stores.QueryFailed(collOfferedData, `Get`, err)
	}

	return val.(string), nil
}
