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

type ContractNegotiation interface {
	Set(cnId string, val negotiation.Negotiation)
	Get(cnId string) (negotiation.Negotiation, error)
	State(cnId string) (negotiation.State, error)
	SetAssignee(cnId string, a odrl.Assignee)
	SetAssigner(cnId string, a odrl.Assigner)
}

type Policy interface {
	SetOffer(id string, val odrl.Offer)
	Offer(id string) (odrl.Offer, error)
}

type Dataset interface {
	Set(id string, val dcat.Dataset)
	Get(id string) (dcat.Dataset, error)
}
