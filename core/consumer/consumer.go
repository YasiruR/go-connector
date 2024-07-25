package consumer

import (
	"github.com/YasiruR/connector/core/consumer/catalog"
	"github.com/YasiruR/connector/core/consumer/negotiation"
	"github.com/YasiruR/connector/core/consumer/transfer"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/core/consumer"
)

// validator should verify states before transitioning into next, signatures, authorization

type Consumer struct {
	consumer.CatalogController
	consumer.NegotiationController
	consumer.NegotiationHandler
	consumer.TransferController
}

func New(port int, stores domain.Stores, plugins domain.Plugins) *Consumer {
	return &Consumer{
		CatalogController:     catalog.NewController(plugins.Client),
		NegotiationController: negotiation.NewController(port, stores, plugins),
		NegotiationHandler:    negotiation.NewHandler(stores, plugins),
		TransferController:    transfer.NewController(port, stores, plugins),
	}
}

//type Consumer struct {
//	catCtrl    consumer.CatalogController
//	negCtrl    consumer.NegotiationController
//	negHandler consumer.NegotiationHandler
//}
//
//func New(port int, stores domain.Stores, plugins domain.Plugins) *Consumer {
//	return &Consumer{
//		catCtrl:    internalCatalog.NewController(plugins.Client),
//		negCtrl:    internalNegotiation.NewController(port, stores, plugins),
//		negHandler: internalNegotiation.NewHandler(stores, plugins),
//	}
//}
//
//func (c *Consumer) RequestCatalog(endpoint string) (catalog2.Response, error) {
//	return c.catCtrl.RequestCatalog(endpoint)
//}
//
//func (c *Consumer) RequestDataset(id, endpoint string) (catalog2.DatasetResponse, error) {
//	return c.catCtrl.RequestDataset(id, endpoint)
//}
//
//func (c *Consumer) RequestContract(providerEndpoint string, ofr odrl.Offer) (negotiationId string, err error) {
//	return c.negCtrl.RequestContract(providerEndpoint, ofr)
//}
//
//func (c *Consumer) AcceptContract() {
//	c.negCtrl.AcceptContract()
//}
//
//func (c *Consumer) VerifyAgreement(consumerPid string) error {
//	return c.negCtrl.VerifyAgreement(consumerPid)
//}
//
//func (c *Consumer) TerminateContract() {
//	c.negCtrl.TerminateContract()
//}
//
//func (c *Consumer) HandleContractAgreement(ca negotiation2.ContractAgreement) (negotiation2.Ack, error) {
//	return c.negHandler.HandleContractAgreement(ca)
//}
//
//func (c *Consumer) HandleFinalizedEvent(consumerPid string) (negotiation2.Ack, error) {
//	return c.negHandler.HandleFinalizedEvent(consumerPid)
//}
