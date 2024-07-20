package stores

import (
	"github.com/YasiruR/connector/core/dsp/negotiation"
	"github.com/YasiruR/connector/core/protocols/dcat"
	"github.com/YasiruR/connector/core/protocols/odrl"
)

type ContractNegotiation interface {
	Set(cnId string, val negotiation.Negotiation)
	Get(cnId string) (negotiation.Negotiation, error)
	GetState(cnId string) (negotiation.State, error)
}

type Policy interface {
	SetOffer(id string, val odrl.Offer)
	GetOffer(id string) (odrl.Offer, error)
}

type Dataset interface {
	Set(id string, val dcat.Dataset)
	Get(id string) (dcat.Dataset, error)
}
