package catalog

import (
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/dsp/http/catalog"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/models/odrl"
	"github.com/YasiruR/connector/domain/pkg"
)

const consCatCollection = `consumer-catalog`

type ConsumerCatalog struct {
	urn   pkg.URNService
	store pkg.Collection
}

func NewConsumerCatalog(plugins domain.Plugins) *ConsumerCatalog {
	return &ConsumerCatalog{urn: plugins.URNService, store: plugins.Database.NewCollection()}
}

func (c *ConsumerCatalog) AddCatalog(res catalog.Response) {
	_ = c.store.Set(res.DspaceParticipantID, res)
}

func (c *ConsumerCatalog) Catalog(providerId string) (catalog.Response, error) {
	val, err := c.store.Get(providerId)
	if err != nil {
		return catalog.Response{}, errors.QueryFailed(consCatCollection, `Get`, err)
	}

	return val.(catalog.Response), nil
}

func (c *ConsumerCatalog) Offer(offerId string) (ofr odrl.Offer, err error) {
	// should optimize this function
	// - expensive AllCatalogs function
	// - nested loops
	cats, err := c.AllCatalogs()
	if err != nil {
		return odrl.Offer{}, errors.StoreFailed(consCatCollection, `AllCatalogs`, err)
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

	return odrl.Offer{}, errors.InvalidKey(offerId)
}

//func (c *ConsumerCatalog) ConnectorEndpoints(providerId string) ([]string, error) {
//	cat, err := c.Catalog(providerId)
//	if err != nil {
//		return nil, errors.StoreFailed(stores.TypeOffer, `Catalog`, err)
//	}
//
//	var endpoints []string
//	for _, svc := range cat.DcatService {
//		if svc.EndpointDescription == core.ServiceConnector {
//			endpoints = append(endpoints, svc.EndpointURL)
//		}
//	}
//
//	if len(endpoints) == 0 {
//		return nil, fmt.Errorf(`no connector endpoints found for provider %s`, providerId)
//	}
//
//	return endpoints, nil
//}

func (c *ConsumerCatalog) AllCatalogs() ([]catalog.Response, error) {
	vals, err := c.store.GetAll()
	if err != nil {
		return nil, errors.QueryFailed(consCatCollection, `GetAll`, err)
	}

	res := make([]catalog.Response, len(vals))
	for _, val := range vals {
		res = append(res, val.(catalog.Response))
	}

	return res, nil
}
