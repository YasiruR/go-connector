package consumer

import (
	"github.com/YasiruR/connector/core"
	"github.com/YasiruR/connector/core/dsp/catalog"
	"github.com/YasiruR/connector/core/dsp/negotiation"
	internalCatalog "github.com/YasiruR/connector/dsp/consumer/catalog"
	internalNegotiation "github.com/YasiruR/connector/dsp/consumer/negotiation"
)

// validator should verify states before transitioning into next, signatures, authorization

type Service struct {
	catCtrl    catalog.Consumer
	negCtrl    negotiation.ConsumerController
	negHandler negotiation.ConsumerHandler
}

func NewService(port int, stores core.Stores, plugins core.Plugins) *Service {
	return &Service{
		catCtrl:    internalCatalog.NewController(plugins.Client),
		negCtrl:    internalNegotiation.NewController(port, stores, plugins),
		negHandler: internalNegotiation.NewHandler(stores, plugins),
	}
}

func (s *Service) RequestCatalog(endpoint string) (catalog.Response, error) {
	return s.catCtrl.RequestCatalog(endpoint)
}

func (s *Service) RequestDataset(id, endpoint string) (catalog.DatasetResponse, error) {
	return s.catCtrl.RequestDataset(id, endpoint)
}

func (s *Service) RequestContract(offerId, providerEndpoint, providerPid, target, assigner, assignee, action string) (negotiationId string, err error) {
	return s.negCtrl.RequestContract(offerId, providerEndpoint, providerPid, target, assigner, assignee, action)
}

func (s *Service) AcceptContract() {
	s.negCtrl.AcceptContract()
}

func (s *Service) VerifyAgreement(consumerPid string) error {
	return s.negCtrl.VerifyAgreement(consumerPid)
}

func (s *Service) TerminateContract() {
	s.negCtrl.TerminateContract()
}

func (s *Service) HandleContractAgreement(consumerPid string, ca negotiation.ContractAgreement) (negotiation.Ack, error) {
	return s.negHandler.HandleContractAgreement(consumerPid, ca)
}
