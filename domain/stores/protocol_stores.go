package stores

import (
	"github.com/YasiruR/go-connector/domain/api/dsp/http/negotiation"
	"github.com/YasiruR/go-connector/domain/api/dsp/http/transfer"
	"github.com/YasiruR/go-connector/domain/models/odrl"
)

// ContractNegotiationStore includes get and set methods for attributes required
// in Negotiation Protocol, such as process information, states and participants
// as defined by IDSA standards ('cnId' refers to Contract Negotiation ID).
type ContractNegotiationStore interface {
	AddNegotiation(cnId string, val negotiation.Negotiation)
	Negotiation(cnId string) (negotiation.Negotiation, error)
	UpdateState(cnId string, s negotiation.State) error
	State(cnId string) (negotiation.State, error)
	SetParticipants(cnId, callbackAddr string, assigner odrl.Assigner, assignee odrl.Assignee)
	Assignee(cnId string) (odrl.Assignee, error)
	Assigner(cnId string) (odrl.Assigner, error)
	CallbackAddr(cnId string) (string, error)
}

// TransferStore includes get and set methods for attributes required
// in Transfer Protocol, such as process information, states and participants
// as defined by IDSA standards ('cnId' refers to Contract Negotiation ID).
type TransferStore interface {
	AddProcess(tpId string, val transfer.Process)
	Process(id string) (transfer.Process, error)
	SetCallbackAddr(tpId, addr string)
	CallbackAddr(tpId string) (string, error)
	UpdateState(tpId string, s transfer.State) error
}
