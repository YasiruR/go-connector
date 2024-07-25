package provider

import (
	internalCatalog "github.com/YasiruR/connector/core/provider/catalog"
	internalNegotiation "github.com/YasiruR/connector/core/provider/negotiation"
	"github.com/YasiruR/connector/domain"
	catalog2 "github.com/YasiruR/connector/domain/api/dsp/http/catalog"
	negotiation2 "github.com/YasiruR/connector/domain/api/dsp/http/negotiation"
	"github.com/YasiruR/connector/domain/core/provider"
)

type Provider struct {
	catCtrl    provider.CatalogService
	negCtrl    provider.NegotiationController
	negHandler provider.NegotiationHandler
}

func New(port int, stores domain.Stores, plugins domain.Plugins) *Provider {
	return &Provider{
		catCtrl:    internalCatalog.NewHandler(stores.Catalog, plugins.Log),
		negCtrl:    internalNegotiation.NewController(port, stores, plugins),
		negHandler: internalNegotiation.NewHandler(stores.ContractNegotiation, plugins),
	}
}

func (p *Provider) HandleCatalogRequest(filter any) (catalog2.Response, error) {
	return p.catCtrl.HandleCatalogRequest(filter)
}

func (p *Provider) HandleDatasetRequest(id string) (catalog2.DatasetResponse, error) {
	return p.catCtrl.HandleDatasetRequest(id)
}

func (p *Provider) OfferContract() {
	p.negCtrl.OfferContract()
}

func (p *Provider) AgreeContract(offerId, negotiationId string) (agreementId string, err error) {
	return p.negCtrl.AgreeContract(offerId, negotiationId)
}

func (p *Provider) FinalizeContract(providerPid string) error {
	return p.negCtrl.FinalizeContract(providerPid)
}

func (p *Provider) HandleNegotiationsRequest(providerPid string) (negotiation2.Ack, error) {
	return p.negHandler.HandleNegotiationsRequest(providerPid)
}

func (p *Provider) HandleContractRequest(cr negotiation2.ContractRequest) (negotiation2.Ack, error) {
	return p.negHandler.HandleContractRequest(cr)
}

func (p *Provider) HandleAgreementVerification(providerPid string) (negotiation2.Ack, error) {
	return p.negHandler.HandleAgreementVerification(providerPid)
}
