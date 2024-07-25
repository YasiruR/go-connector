package provider

import (
	"github.com/YasiruR/connector/core/provider/catalog"
	"github.com/YasiruR/connector/core/provider/negotiation"
	"github.com/YasiruR/connector/core/provider/transfer"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/core/provider"
)

type Provider struct {
	provider.CatalogHandler
	provider.NegotiationController
	provider.NegotiationHandler
	provider.TransferHandler
}

func New(port int, stores domain.Stores, plugins domain.Plugins) *Provider {
	return &Provider{
		CatalogHandler:        catalog.NewHandler(stores.Catalog, plugins.Log),
		NegotiationController: negotiation.NewController(port, stores, plugins),
		NegotiationHandler:    negotiation.NewHandler(stores.ContractNegotiation, plugins),
		TransferHandler:       transfer.NewHandler(stores, plugins),
	}
}

//type Provider struct {
//	catCtrl    provider.CatalogHandler
//	negCtrl    provider.NegotiationController
//	negHandler provider.NegotiationHandler
//	trHandler  provider.TransferHandler
//}
//
//func New(port int, stores domain.Stores, plugins domain.Plugins) *Provider {
//	return &Provider{
//		catCtrl:    internalCatalog.NewHandler(stores.Catalog, plugins.Log),
//		negCtrl:    internalNegotiation.NewController(port, stores, plugins),
//		negHandler: internalNegotiation.NewHandler(stores.ContractNegotiation, plugins),
//	}
//}
//
//func (p *Provider) HandleCatalogRequest(filter any) (catalog2.Response, error) {
//	return p.catCtrl.HandleCatalogRequest(filter)
//}
//
//func (p *Provider) HandleDatasetRequest(id string) (catalog2.DatasetResponse, error) {
//	return p.catCtrl.HandleDatasetRequest(id)
//}
//
//func (p *Provider) OfferContract() {
//	p.negCtrl.OfferContract()
//}
//
//func (p *Provider) AgreeContract(offerId, negotiationId string) (agreementId string, err error) {
//	return p.negCtrl.AgreeContract(offerId, negotiationId)
//}
//
//func (p *Provider) FinalizeContract(providerPid string) error {
//	return p.negCtrl.FinalizeContract(providerPid)
//}
//
//func (p *Provider) HandleNegotiationsRequest(providerPid string) (negotiation2.Ack, error) {
//	return p.negHandler.HandleNegotiationsRequest(providerPid)
//}
//
//func (p *Provider) HandleContractRequest(cr negotiation2.ContractRequest) (negotiation2.Ack, error) {
//	return p.negHandler.HandleContractRequest(cr)
//}
//
//func (p *Provider) HandleAgreementVerification(providerPid string) (negotiation2.Ack, error) {
//	return p.negHandler.HandleAgreementVerification(providerPid)
//}
//
//func (p *Provider)
