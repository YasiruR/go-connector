package stores

import (
	"github.com/YasiruR/connector/core/dsp/negotiation"
	"github.com/YasiruR/connector/core/errors"
	"github.com/YasiruR/connector/core/pkg"
	"github.com/YasiruR/connector/core/protocols/odrl"
)

type ContractNegotiation struct {
	store     pkg.DataStore
	assignees pkg.DataStore
	assigners pkg.DataStore
}

func NewContractNegotiationStore(db pkg.Database) *ContractNegotiation {
	return &ContractNegotiation{
		store:     db.NewDataStore(),
		assignees: db.NewDataStore(),
		assigners: db.NewDataStore(),
	}
}

func (cn *ContractNegotiation) Set(cnId string, val negotiation.Negotiation) {
	_ = cn.store.Set(cnId, val)
}

func (cn *ContractNegotiation) Get(cnId string) (negotiation.Negotiation, error) {
	val, err := cn.store.Get(cnId)
	if err != nil {
		return negotiation.Negotiation{}, errors.QueryFailed(`get`, err)
	}
	return val.(negotiation.Negotiation), nil
}

func (cn *ContractNegotiation) State(cnId string) (negotiation.State, error) {
	cnAck, err := cn.Get(cnId)
	if err != nil {
		return ``, errors.QueryFailed(`get`, err)
	}
	return cnAck.State, nil
}

func (cn *ContractNegotiation) SetAssignee(cnId string, a odrl.Assignee) {
	_ = cn.assignees.Set(cnId, a)
}

//func (cn *ContractNegotiation) Assignee(cnId string) (odrl.Assignee, error) {
//	return ``, nil
//}

func (cn *ContractNegotiation) SetAssigner(cnId string, a odrl.Assigner) {
	_ = cn.assigners.Set(cnId, a)
}

//func (cn *ContractNegotiation) Assigner(cnId string) (odrl.Assigner, error) {
//
//}
