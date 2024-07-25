package negotiation

import (
	"fmt"
	"github.com/YasiruR/connector/domain"
	negotiation2 "github.com/YasiruR/connector/domain/api/dsp/http/negotiation"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/stores"
)

type Handler struct {
	negStore stores.ContractNegotiation
	agrStore stores.Agreement
	log      pkg.Log
}

func NewHandler(stores domain.Stores, plugins domain.Plugins) *Handler {
	return &Handler{
		negStore: stores.ContractNegotiation,
		agrStore: stores.Agreement,
		log:      plugins.Log,
	}
}

func (h *Handler) HandleContractAgreement(ca negotiation2.ContractAgreement) (negotiation2.Ack, error) {
	// validate agreement (e.g. consumerPid, target)

	h.agrStore.Set(ca.ConsPId, ca.Agreement)
	h.log.Trace(fmt.Sprintf("stored contract agreement (id: %s) for negotation (id: %s)", ca.Agreement.Id, ca.ConsPId))
	h.negStore.SetCallbackAddr(ca.ConsPId, ca.CallbackAddr)

	if err := h.negStore.UpdateState(ca.ConsPId, negotiation2.StateAgreed); err != nil {
		return negotiation2.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	neg, err := h.negStore.Negotiation(ca.ConsPId)
	if err != nil {
		return negotiation2.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	h.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", ca.ConsPId, negotiation2.StateAgreed))
	return negotiation2.Ack(neg), nil
}

func (h *Handler) HandleFinalizedEvent(consumerPid string) (negotiation2.Ack, error) {
	if err := h.negStore.UpdateState(consumerPid, negotiation2.StateFinalized); err != nil {
		return negotiation2.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	neg, err := h.negStore.Negotiation(consumerPid)
	if err != nil {
		return negotiation2.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	h.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", consumerPid, negotiation2.StateFinalized))
	return negotiation2.Ack(neg), nil
}
