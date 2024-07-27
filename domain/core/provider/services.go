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
	OfferContract()
	AgreeContract(offerId, negotiationId string) (agreementId string, err error)
	FinalizeContract(providerPid string) error
}

type NegotiationHandler interface {
	HandleNegotiationsRequest(providerPid string) (negotiation.Ack, error)
	HandleContractRequest(cr negotiation.ContractRequest) (negotiation.Ack, error)
	HandleAgreementVerification(providerPid string) (negotiation.Ack, error)
}

type TransferController interface {
	StartTransfer(tpId, sourceEndpoint string) error
}

type TransferHandler interface {
	HandleTransferRequest(tr transfer.Request) (transfer.Ack, error)
}
