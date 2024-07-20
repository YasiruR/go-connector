package stores

import (
	"fmt"
	"github.com/YasiruR/connector/core/dsp/negotiation"
	"github.com/YasiruR/connector/core/pkg"
)

type ContractNegotiation struct {
	db pkg.Database
}

func NewContractNegotiationStore(db pkg.Database) *ContractNegotiation {
	return &ContractNegotiation{db: db}
}

func (cn *ContractNegotiation) Set(cnId string, val negotiation.Negotiation) {
	_ = cn.db.Set(cnId, val)
}

func (cn *ContractNegotiation) Get(cnId string) (negotiation.Negotiation, error) {
	val, err := cn.db.Get(cnId)
	if err != nil {
		return negotiation.Negotiation{}, fmt.Errorf("database failure - %w", err)
	}
	return val.(negotiation.Negotiation), nil
}

func (cn *ContractNegotiation) GetState(cnId string) (negotiation.State, error) {
	cnAck, err := cn.Get(cnId)
	if err != nil {
		return ``, fmt.Errorf("fetching contract negotiation failed - %w", err)
	}
	return cnAck.State, nil
}
