package consumer

import (
	"github.com/YasiruR/connector/domain/api/dsp/http/catalog"
	"github.com/YasiruR/connector/domain/api/dsp/http/negotiation"
	"github.com/YasiruR/connector/domain/models/odrl"
)

type CatalogController interface {
	RequestCatalog(endpoint string) (catalog.Response, error) // endpoint should be generic
	RequestDataset(id, endpoint string) (catalog.DatasetResponse, error)
}

type NegotiationController interface {
	// change endpoint to generic
	RequestContract(providerEndpoint string, ofr odrl.Offer) (cnId string, err error)
	AcceptContract()
	VerifyAgreement(consumerPid string) error
	TerminateContract()
}

type NegotiationHandler interface {
	HandleContractAgreement(ca negotiation.ContractAgreement) (negotiation.Ack, error)
	HandleFinalizedEvent(consumerPid string) (negotiation.Ack, error)
}

type TransferController interface {
	RequestTransfer(transferType, agreementId, sinkEndpoint, providerEndpoint string) (tpId string, err error)
}
