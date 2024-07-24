package stores

import (
	"github.com/YasiruR/connector/core"
	"github.com/YasiruR/connector/core/dsp/negotiation"
	"github.com/YasiruR/connector/core/errors"
	"github.com/YasiruR/connector/core/pkg"
	"github.com/YasiruR/connector/core/protocols/odrl"
)

const (
	negotiationCollection  = `negotiation-collection`
	assigneeCollection     = `assignee-collection`
	callbackAddrCollection = `callbackAddr-collection`
)

// ContractNegotiation stores any ongoing activities related to Contract Negotiation Protocol
type ContractNegotiation struct {
	store        pkg.Collection
	assignees    pkg.Collection
	callbackAddr pkg.Collection
}

func NewContractNegotiationStore(plugins core.Plugins) *ContractNegotiation {
	plugins.Log.Info("initialized contract negotiation store")
	return &ContractNegotiation{
		store:        plugins.Database.NewCollection(),
		assignees:    plugins.Database.NewCollection(),
		callbackAddr: plugins.Database.NewCollection(),
	}
}

func (cn *ContractNegotiation) Set(cnId string, val negotiation.Negotiation) {
	_ = cn.store.Set(cnId, val)
}

func (cn *ContractNegotiation) Negotiation(cnId string) (negotiation.Negotiation, error) {
	val, err := cn.store.Get(cnId)
	if err != nil {
		return negotiation.Negotiation{}, errors.QueryFailed(negotiationCollection, `Get`, err)
	}
	return val.(negotiation.Negotiation), nil
}

func (cn *ContractNegotiation) UpdateState(cnId string, s negotiation.State) error {
	neg, err := cn.Negotiation(cnId)
	if err != nil {
		return errors.QueryFailed(negotiationCollection, `Get`, err)
	}

	neg.State = s
	cn.Set(cnId, neg)
	return nil
}

func (cn *ContractNegotiation) State(cnId string) (negotiation.State, error) {
	neg, err := cn.Negotiation(cnId)
	if err != nil {
		return ``, errors.QueryFailed(negotiationCollection, `Get`, err)
	}
	return neg.State, nil
}

func (cn *ContractNegotiation) SetAssignee(cnId string, a odrl.Assignee) {
	_ = cn.assignees.Set(cnId, a)
}

func (cn *ContractNegotiation) Assignee(cnId string) (odrl.Assignee, error) {
	val, err := cn.assignees.Get(cnId)
	if err != nil {
		return ``, errors.QueryFailed(assigneeCollection, `get`, err)
	}
	return val.(odrl.Assignee), nil
}

func (cn *ContractNegotiation) SetCallbackAddr(cnId, addr string) {
	_ = cn.callbackAddr.Set(cnId, addr)
}

func (cn *ContractNegotiation) CallbackAddr(cnId string) (string, error) {
	addr, err := cn.callbackAddr.Get(cnId)
	if err != nil {
		return ``, errors.QueryFailed(callbackAddrCollection, `get`, err)
	}
	return addr.(string), nil
}

//func (cn *ContractNegotiation) SetAssigner(cnId string, a odrl.Assigner) {
//	_ = cn.assigners.Set(cnId, a)
//}
//
//func (cn *ContractNegotiation) Assigner(cnId string) (odrl.Assigner, error) {
//
//}
