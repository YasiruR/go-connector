package stores

import (
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/dsp/http/negotiation"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/models/odrl"
	"github.com/YasiruR/connector/domain/pkg"
)

const (
	negotiationCollection  = `negotiation`
	assigneeCollection     = `assignee`
	assignerCollection     = `assigner`
	callbackAddrCollection = `callbackAddr`
)

// ContractNegotiation stores any ongoing activities related to Contract Negotiation Protocol
type ContractNegotiation struct {
	store        pkg.Collection
	assignees    pkg.Collection
	assigners    pkg.Collection
	callbackAddr pkg.Collection
}

func NewContractNegotiationStore(plugins domain.Plugins) *ContractNegotiation {
	plugins.Log.Info("initialized contract negotiation store")
	return &ContractNegotiation{
		store:        plugins.Database.NewCollection(),
		assignees:    plugins.Database.NewCollection(),
		assigners:    plugins.Database.NewCollection(),
		callbackAddr: plugins.Database.NewCollection(),
	}
}

func (cn *ContractNegotiation) AddNegotiation(cnId string, val negotiation.Negotiation) {
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
		return ``, errors.QueryFailed(negotiationCollection, `Get`, err)
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
		return ``, errors.QueryFailed(assigneeCollection, `get`, err)
	}
	return val.(odrl.Assignee), nil
}

func (cn *ContractNegotiation) Assigner(cnId string) (odrl.Assigner, error) {
	val, err := cn.assigners.Get(cnId)
	if err != nil {
		return ``, errors.QueryFailed(assignerCollection, `get`, err)
	}
	return val.(odrl.Assigner), nil
}

func (cn *ContractNegotiation) CallbackAddr(cnId string) (string, error) {
	addr, err := cn.callbackAddr.Get(cnId)
	if err != nil {
		return ``, errors.QueryFailed(callbackAddrCollection, `get`, err)
	}
	return addr.(string), nil
}
