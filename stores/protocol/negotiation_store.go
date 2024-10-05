package protocol

import (
	"github.com/YasiruR/go-connector/domain"
	"github.com/YasiruR/go-connector/domain/api/dsp/http/negotiation"
	"github.com/YasiruR/go-connector/domain/models/odrl"
	"github.com/YasiruR/go-connector/domain/pkg"
	"github.com/YasiruR/go-connector/domain/stores"
)

const (
	collNegotiation  = `negotiation`
	collAssignee     = `assignee`
	collAssigner     = `assigner`
	collCallbackAddr = `callbackAddr`
)

// ContractNegotiation stores any ongoing activities related to Contract Negotiation Protocol
type ContractNegotiation struct {
	negotiations pkg.Collection
	assignees    pkg.Collection
	assigners    pkg.Collection
	callbackAddr pkg.Collection
}

func NewContractNegotiationStore(plugins domain.Plugins) *ContractNegotiation {
	plugins.Log.Info("initialized contract negotiation store")
	return &ContractNegotiation{
		negotiations: plugins.Store.NewCollection(),
		assignees:    plugins.Store.NewCollection(),
		assigners:    plugins.Store.NewCollection(),
		callbackAddr: plugins.Store.NewCollection(),
	}
}

func (cn *ContractNegotiation) AddNegotiation(cnId string, val negotiation.Negotiation) {
	_ = cn.negotiations.Set(cnId, val)
}

func (cn *ContractNegotiation) Negotiation(cnId string) (negotiation.Negotiation, error) {
	val, err := cn.negotiations.Get(cnId)
	if err != nil {
		return negotiation.Negotiation{}, stores.QueryFailed(collNegotiation, `Get`, err)
	}

	if val == nil {
		return negotiation.Negotiation{}, stores.InvalidKey(cnId)
	}

	return val.(negotiation.Negotiation), nil
}

func (cn *ContractNegotiation) UpdateState(cnId string, s negotiation.State) error {
	neg, err := cn.Negotiation(cnId)
	if err != nil {
		// todo should I handle invalid key error here?
		return stores.QueryFailed(collNegotiation, `Get`, err)
	}

	//switch s {
	//case negotiation.StateRequested:
	//	if neg.State != negotiation.StateOffered && neg.State != `` {
	//		return negotiation.Negotiation{}, errors.IncompatibleValues(`state`, string(neg.State), string(negotiation.StateOffered)+" or null")
	//	}
	//case negotiation.StateOffered:
	//	if neg.State != negotiation.StateRequested && neg.State != `` {
	//		return negotiation.Negotiation{}, errors.IncompatibleValues(`state`, string(neg.State), string(negotiation.StateRequested)+" or null")
	//	}
	//case negotiation.StateAccepted:
	//	if neg.State != negotiation.StateOffered {
	//		return negotiation.Negotiation{}, errors.IncompatibleValues(`state`, string(neg.State), string(negotiation.StateOffered))
	//	}
	//case negotiation.StateAgreed:
	//	if neg.State != negotiation.StateAccepted && neg.State != negotiation.StateRequested {
	//		return negotiation.Negotiation{}, errors.IncompatibleValues(`state`, string(neg.State), string(negotiation.StateAccepted)+
	//			" or "+string(negotiation.StateRequested))
	//	}
	//case negotiation.StateVerified:
	//	if neg.State != negotiation.StateAgreed {
	//		return negotiation.Negotiation{}, errors.IncompatibleValues(`state`, string(neg.State), string(negotiation.StateAgreed))
	//	}
	//case negotiation.StateFinalized:
	//	if neg.State != negotiation.StateVerified {
	//		return negotiation.Negotiation{}, errors.IncompatibleValues(`state`, string(neg.State), string(negotiation.StateVerified))
	//	}
	//default:
	//	return negotiation.Negotiation{}, errors.IncompatibleValues(`state`, string(neg.State), "valid state value")
	//}

	neg.State = s
	cn.AddNegotiation(cnId, neg)
	return nil
}

func (cn *ContractNegotiation) State(cnId string) (negotiation.State, error) {
	neg, err := cn.Negotiation(cnId)
	if err != nil {
		// todo should I handle invalid key error here?
		return ``, stores.QueryFailed(collNegotiation, `Get`, err)
	}
	return neg.State, nil
}

func (cn *ContractNegotiation) SetParticipants(cnId, callbackAddr string, assigner odrl.Assigner, assignee odrl.Assignee) {
	_ = cn.assigners.Set(cnId, assigner)
	_ = cn.assignees.Set(cnId, assignee)
	_ = cn.callbackAddr.Set(cnId, callbackAddr)
}

func (cn *ContractNegotiation) Assignee(cnId string) (odrl.Assignee, error) {
	val, err := cn.assignees.Get(cnId)
	if err != nil {
		return ``, stores.QueryFailed(collAssignee, `get`, err)
	}

	if val == nil {
		return ``, stores.InvalidKey(cnId)
	}

	return val.(odrl.Assignee), nil
}

func (cn *ContractNegotiation) Assigner(cnId string) (odrl.Assigner, error) {
	val, err := cn.assigners.Get(cnId)
	if err != nil {
		return ``, stores.QueryFailed(collAssigner, `get`, err)
	}

	if val == nil {
		return ``, stores.InvalidKey(cnId)
	}

	return val.(odrl.Assigner), nil
}

func (cn *ContractNegotiation) CallbackAddr(cnId string) (string, error) {
	addr, err := cn.callbackAddr.Get(cnId)
	if err != nil {
		return ``, stores.QueryFailed(collCallbackAddr, `get`, err)
	}

	if addr == nil {
		return ``, stores.InvalidKey(cnId)
	}

	return addr.(string), nil
}
