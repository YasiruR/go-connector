package stores

import (
	"fmt"
	"github.com/YasiruR/connector/core/dsp/negotiation"
	"github.com/YasiruR/connector/core/pkg"
	"github.com/YasiruR/connector/pkg/store/memory"
)

type ContractNegotiation struct {
	store pkg.Store
}

func NewContractNegotiationStore() *ContractNegotiation {
	return &ContractNegotiation{store: memory.NewStore()}
}

func (cn *ContractNegotiation) Set(cnId string, val negotiation.Ack) {
	_ = cn.store.Set(cnId, val)
}

func (cn *ContractNegotiation) Get(cnId string) (negotiation.Ack, error) {
	val, err := cn.store.Get(cnId)
	if err != nil {
		return negotiation.Ack{}, fmt.Errorf("store failure - %w", err)
	}
	return val.(negotiation.Ack), nil
}

func (cn *ContractNegotiation) GetState(cnId string) (negotiation.State, error) {
	cnAck, err := cn.Get(cnId)
	if err != nil {
		return ``, fmt.Errorf("fetching contract negotiation failed - %w", err)
	}
	return cnAck.State, nil
}
