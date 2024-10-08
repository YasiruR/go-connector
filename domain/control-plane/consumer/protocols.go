package consumer

import (
	"github.com/YasiruR/go-connector/domain/api/dsp/http/catalog"
	"github.com/YasiruR/go-connector/domain/api/dsp/http/negotiation"
	"github.com/YasiruR/go-connector/domain/api/dsp/http/transfer"
	"github.com/YasiruR/go-connector/domain/models"
)

type CatalogController interface {
	RequestCatalog(endpoint string) (catalog.Response, error) // endpoint should be generic
	RequestDataset(id, endpoint string) (catalog.DatasetResponse, error)
}

type NegotiationController interface {
	RequestContract(consumerPid, providerAddr, offerId string, constraints map[string]string) (cnId string, err error)
	AcceptOffer(consumerPid string) error
	VerifyAgreement(consumerPid string) error
	TerminateContract(consumerPid, code string, reasons []string) error
}

type NegotiationHandler interface {
	HandleContractOffer(co negotiation.ContractOffer) (ack negotiation.Ack, err error)
	HandleContractAgreement(ca negotiation.ContractAgreement) (negotiation.Ack, error)
	HandleFinalizedEvent(e negotiation.ContractNegotiationEvent) (negotiation.Ack, error)
}

type TransferController interface {
	GetProviderProcess(tpId string) (transfer.Process, error)
	RequestTransfer(transferType, agreementId, providerAddr string, sinkDb models.Database) (tpId string, err error)
	SuspendTransfer(tpId, code string, reasons []interface{}) error
	StartTransfer(tpId string) error
	CompleteTransfer(tpId string) error
	TerminateTransfer(tpId, code string, reasons []interface{}) error
}

type TransferHandler interface {
	HandleTransferStart(sr transfer.StartRequest) (transfer.Ack, error)
	HandleTransferSuspension(sr transfer.SuspendRequest) (transfer.Ack, error)
	HandleTransferCompletion(cr transfer.CompleteRequest) (transfer.Ack, error)
	HandleTransferTermination(tr transfer.TerminateRequest) (transfer.Ack, error)
}
