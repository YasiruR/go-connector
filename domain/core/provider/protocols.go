package provider

import (
	"github.com/YasiruR/connector/domain/api/dsp/http/catalog"
	"github.com/YasiruR/connector/domain/api/dsp/http/negotiation"
	"github.com/YasiruR/connector/domain/api/dsp/http/transfer"
)

type CatalogHandler interface {
	HandleCatalogRequest(filter any) (catalog.Response, error)
	HandleDatasetRequest(id string) (catalog.DatasetResponse, error)
}

type NegotiationController interface {
	// OfferContract sends an Offer to the consumer. providerPid and consumerAddr parameters are
	// mutually exclusive. Former should be given when the Provider is responding to a Contract
	// Request by a Consumer, whereas the latter when the Provider is the initiator of the flow.
	OfferContract(offerId, providerPid, consumerAddr string) (cnId string, err error)
	AgreeContract(offerId, negotiationId string) (agreementId string, err error)
	FinalizeContract(providerPid string) error
}

type NegotiationHandler interface {
	HandleNegotiationsRequest(providerPid string) (negotiation.Ack, error)
	HandleContractRequest(cr negotiation.ContractRequest) (negotiation.Ack, error)
	HandleAcceptOffer(providerPid string) (negotiation.Ack, error)
	HandleAgreementVerification(providerPid string) (negotiation.Ack, error)
}

type TransferController interface {
	StartTransfer(tpId, sourceEndpoint string) error
	SuspendTransfer(tpId, code string, reasons []interface{}) error
	CompleteTransfer(tpId string) error
}

type TransferHandler interface {
	//  should be idempotent for multiple transfer requests
	HandleTransferRequest(tr transfer.Request) (transfer.Ack, error)
	HandleTransferSuspension(sr transfer.SuspendRequest) (transfer.Ack, error)
	HandleTransferCompletion(cr transfer.CompleteRequest) (transfer.Ack, error)
}
