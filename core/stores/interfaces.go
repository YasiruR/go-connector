package stores

import (
	"github.com/YasiruR/connector/boot/config"
	"github.com/YasiruR/connector/core/dsp/negotiation"
	"github.com/YasiruR/connector/core/protocols/dcat"
	"github.com/YasiruR/connector/core/protocols/odrl"
)

const (
	TypeCatalog             = `Catalog`
	TypeContractNegotiation = `ContractNegotiation`
	TypePolicy              = `Policy`
)

// Catalog stores Datasets as per the DCAT profile recommended by IDSA
type Catalog interface {
	Init(cfg config.Config) error
	Get() (dcat.Catalog, error)
	AddDataset(id string, val dcat.Dataset)
	Dataset(id string) (dcat.Dataset, error)
}

// ContractNegotiation includes get and set methods for attributes required
// in Negotiation Protocol ('cnId' refers to Contract Negotiation ID)
type ContractNegotiation interface {
	Set(cnId string, val negotiation.Negotiation)
	Get(cnId string) (negotiation.Negotiation, error)
	UpdateState(cnId string, s negotiation.State) error
	State(cnId string) (negotiation.State, error)
	SetAssignee(cnId string, a odrl.Assignee)
	Assignee(cnId string) (odrl.Assignee, error)
	SetCallbackAddr(cnId, addr string)
	CallbackAddr(cnId string) (string, error)
}

type Policy interface {
	SetOffer(id string, val odrl.Offer)
	Offer(id string) (odrl.Offer, error)
}
