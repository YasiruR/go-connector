package stores

import "github.com/YasiruR/connector/core/dsp/negotiation"

type ContractNegotiation interface {
	Set(cnId string, val negotiation.Ack)
	Get(cnId string) (negotiation.Ack, error)
	GetState(cnId string) (negotiation.State, error)
}
