package stores

import (
	"github.com/YasiruR/connector/domain/api/dsp/http/negotiation"
	"github.com/YasiruR/connector/domain/api/dsp/http/transfer"
	"github.com/YasiruR/connector/domain/boot"
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

// Catalog stores Datasets as per the DCAT profile recommended by IDSA
type Catalog interface {
	Init(cfg boot.Config) error
	Get() (dcat.Catalog, error)
	AddDataset(id string, val dcat.Dataset)
	Dataset(id string) (dcat.Dataset, error)
}

// ContractNegotiation includes get and set methods for attributes required
// in Negotiation Protocol ('cnId' refers to Contract Negotiation ID)
type ContractNegotiation interface {
	Set(cnId string, val negotiation.Negotiation)
	Negotiation(cnId string) (negotiation.Negotiation, error)
	UpdateState(cnId string, s negotiation.State) error
	State(cnId string) (negotiation.State, error)
	SetAssignee(cnId string, a odrl.Assignee)
	Assignee(cnId string) (odrl.Assignee, error)
	SetCallbackAddr(cnId, addr string)
	CallbackAddr(cnId string) (string, error)
}

type Policy interface {
	SetOffer(id string, val odrl.Offer)
	GetOffer(id string) (odrl.Offer, error)
}

type Agreement interface {
	// Set stores contract agreement with agreement ID as the key
	Set(id string, val odrl.Agreement)
	// Get retrieves contract agreement by agreement ID
	Get(id string) (odrl.Agreement, error)
	GetByNegotiationID(cnId string) (odrl.Agreement, error)
}

type Transfer interface {
	Set(tpId string, val transfer.Process)
	GetProcess(id string) (transfer.Process, error)
	SetCallbackAddr(tpId, addr string)
	CallbackAddr(tpId string) (string, error)
	UpdateState(tpId string, s transfer.State) error
}
