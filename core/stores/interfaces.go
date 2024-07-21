package stores

import (
	"github.com/YasiruR/connector/core/dsp/negotiation"
	"github.com/YasiruR/connector/core/protocols/dcat"
	"github.com/YasiruR/connector/core/protocols/odrl"
)

const (
	TypeContractNegotiation = `ContractNegotiation`
	TypePolicy              = `Policy`
	TypeDataset             = `Dataset`
)

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

type Dataset interface {
	Set(id string, val dcat.Dataset)
	Get(id string) (dcat.Dataset, error)
}
