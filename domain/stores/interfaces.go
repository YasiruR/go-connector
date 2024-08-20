package stores

import (
	"github.com/YasiruR/connector/domain/api/dsp/http/negotiation"
	"github.com/YasiruR/connector/domain/api/dsp/http/transfer"
	"github.com/YasiruR/connector/domain/models/dcat"
	"github.com/YasiruR/connector/domain/models/odrl"
)

const (
	TypeCatalog             = `catalog`
	TypeContractNegotiation = `contract-negotiation`
	TypePolicy              = `policy`
	TypeAgreement           = `agreement`
	TypeTransfer            = `transfer`
)

// CatalogStore stores Datasets as per the DCAT profile recommended by IDSA
type CatalogStore interface {
	Catalog() (dcat.Catalog, error)
	AddDataset(id string, val dcat.Dataset)
	Dataset(id string) (dcat.Dataset, error)
}

// ContractNegotiationStore includes get and set methods for attributes required
// in Negotiation Protocol ('cnId' refers to Contract Negotiation ID)
type ContractNegotiationStore interface {
	AddNegotiation(cnId string, val negotiation.Negotiation)
	Negotiation(cnId string) (negotiation.Negotiation, error)
	UpdateState(cnId string, s negotiation.State) error
	State(cnId string) (negotiation.State, error)
	SetParticipants(cnId, callbackAddr string, assigner odrl.Assigner, assignee odrl.Assignee)
	Assignee(cnId string) (odrl.Assignee, error)
	Assigner(cnId string) (odrl.Assigner, error)
	CallbackAddr(cnId string) (string, error)
}

type PolicyStore interface {
	AddOffer(id string, val odrl.Offer)
	Offer(id string) (odrl.Offer, error)
}

type AgreementStore interface {
	// AddAgreement stores contract agreement with agreement ID as the key
	AddAgreement(id string, val odrl.Agreement)
	// Agreement retrieves contract agreement by agreement ID
	Agreement(id string) (odrl.Agreement, error)
	AgreementByNegotiationID(cnId string) (odrl.Agreement, error)
}

type TransferStore interface {
	AddProcess(tpId string, val transfer.Process)
	Process(id string) (transfer.Process, error)
	SetCallbackAddr(tpId, addr string)
	CallbackAddr(tpId string) (string, error)
	UpdateState(tpId string, s transfer.State) error
}
