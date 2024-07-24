package consumer

import (
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/protocols/dsp/catalog"
	"github.com/YasiruR/connector/domain/protocols/dsp/negotiation"
	internalCatalog "github.com/YasiruR/connector/dsp/consumer/catalog"
	internalNegotiation "github.com/YasiruR/connector/dsp/consumer/negotiation"
)

// validator should verify states before transitioning into next, signatures, authorization

type Consumer struct {
	catCtrl    catalog.Controller
	negCtrl    negotiation.ControllerConsumer
	negHandler negotiation.HandlerConsumer
}

func New(port int, stores domain.Stores, plugins domain.Plugins) *Consumer {
	return &Consumer{
		catCtrl:    internalCatalog.NewController(plugins.Client),
		negCtrl:    internalNegotiation.NewController(port, stores, plugins),
		negHandler: internalNegotiation.NewHandler(stores, plugins),
	}
}

func (c *Consumer) RequestCatalog(endpoint string) (catalog.Response, error) {
	return c.catCtrl.RequestCatalog(endpoint)
}

func (c *Consumer) RequestDataset(id, endpoint string) (catalog.DatasetResponse, error) {
	return c.catCtrl.RequestDataset(id, endpoint)
}

func (c *Consumer) RequestContract(offerId, providerEndpoint, providerPid, target, assigner, assignee, action string) (negotiationId string, err error) {
	return c.negCtrl.RequestContract(offerId, providerEndpoint, providerPid, target, assigner, assignee, action)
}

func (c *Consumer) AcceptContract() {
	c.negCtrl.AcceptContract()
}

func (c *Consumer) VerifyAgreement(consumerPid string) error {
	return c.negCtrl.VerifyAgreement(consumerPid)
}

func (c *Consumer) TerminateContract() {
	c.negCtrl.TerminateContract()
}

func (c *Consumer) HandleContractAgreement(ca negotiation.ContractAgreement) (negotiation.Ack, error) {
	return c.negHandler.HandleContractAgreement(ca)
}

func (c *Consumer) HandleFinalizedEvent(consumerPid string) (negotiation.Ack, error) {
	return c.negHandler.HandleFinalizedEvent(consumerPid)
}
