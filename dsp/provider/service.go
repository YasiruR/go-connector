package provider

import (
	"github.com/YasiruR/connector/core"
	"github.com/YasiruR/connector/core/dsp/catalog"
	"github.com/YasiruR/connector/core/dsp/negotiation"
	catalog2 "github.com/YasiruR/connector/dsp/provider/catalog"
	negotiation2 "github.com/YasiruR/connector/dsp/provider/negotiation"
)

type Service struct {
	catCtrl    catalog.Provider
	negCtrl    negotiation.ProviderController
	negHandler negotiation.ProviderHandler
}

func NewService(port int, stores core.Stores, plugins core.Plugins) *Service {
	return &Service{
		catCtrl:    catalog2.NewHandler(stores.Catalog, plugins.Log),
		negCtrl:    negotiation2.NewController(port, stores, plugins),
		negHandler: negotiation2.NewHandler(stores.ContractNegotiation, plugins),
	}
}

func (s *Service) HandleCatalogRequest(filter any) (catalog.Response, error) {
	return s.catCtrl.HandleCatalogRequest(filter)
}

func (s *Service) HandleDatasetRequest(id string) (catalog.DatasetResponse, error) {
	return s.catCtrl.HandleDatasetRequest(id)
}

func (s *Service) OfferContract() {
	s.negCtrl.OfferContract()
}

func (s *Service) AgreeContract(offerId, negotiationId string) (agreementId string, err error) {
	return s.negCtrl.AgreeContract(offerId, negotiationId)
}

func (s *Service) FinalizeContract(providerPid string) error {
	return s.negCtrl.FinalizeContract(providerPid)
}

func (s *Service) HandleNegotiationsRequest(providerPid string) (negotiation.Ack, error) {
	return s.negHandler.HandleNegotiationsRequest(providerPid)
}

func (s *Service) HandleContractRequest(cr negotiation.ContractRequest) (negotiation.Ack, error) {
	return s.negHandler.HandleContractRequest(cr)
}

func (s *Service) HandleAgreementVerification(providerPid string) (negotiation.Ack, error) {
	return s.negHandler.HandleAgreementVerification(providerPid)
}
