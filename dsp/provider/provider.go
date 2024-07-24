package provider

import (
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/dsp/catalog"
	"github.com/YasiruR/connector/domain/dsp/negotiation"
	internalCatalog "github.com/YasiruR/connector/dsp/provider/catalog"
	internalNegotiation "github.com/YasiruR/connector/dsp/provider/negotiation"
)

type Provider struct {
	catCtrl    catalog.Handler
	negCtrl    negotiation.ControllerProvider
	negHandler negotiation.HandlerProvider
}

func New(port int, stores domain.Stores, plugins domain.Plugins) *Provider {
	return &Provider{
		catCtrl:    internalCatalog.NewHandler(stores.Catalog, plugins.Log),
		negCtrl:    internalNegotiation.NewController(port, stores, plugins),
		negHandler: internalNegotiation.NewHandler(stores.ContractNegotiation, plugins),
	}
}

func (p *Provider) HandleCatalogRequest(filter any) (catalog.Response, error) {
	return p.catCtrl.HandleCatalogRequest(filter)
}

func (p *Provider) HandleDatasetRequest(id string) (catalog.DatasetResponse, error) {
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

func (p *Provider) HandleNegotiationsRequest(providerPid string) (negotiation.Ack, error) {
	return p.negHandler.HandleNegotiationsRequest(providerPid)
}

func (p *Provider) HandleContractRequest(cr negotiation.ContractRequest) (negotiation.Ack, error) {
	return p.negHandler.HandleContractRequest(cr)
}

func (p *Provider) HandleAgreementVerification(providerPid string) (negotiation.Ack, error) {
	return p.negHandler.HandleAgreementVerification(providerPid)
}
