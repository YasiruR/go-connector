package negotiation

import (
	"fmt"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/dsp/negotiation"
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

func (h *Handler) HandleContractAgreement(consumerPid string, ca negotiation.ContractAgreement) (negotiation.Ack, error) {
	// validate agreement (e.g. consumerPid, target)

	h.agrStore.Set(ca.ConsPId, ca.Agreement)
	h.log.Trace(fmt.Sprintf("stored contract agreement (id: %s) for negotation (id: %s)", ca.Agreement.Id, ca.ConsPId))

	if err := h.negStore.UpdateState(ca.ConsPId, negotiation.StateAgreed); err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	h.negStore.SetCallbackAddr(ca.ConsPId, ca.CallbackAddr)
	neg, err := h.negStore.Negotiation(ca.ConsPId)
	if err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	h.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", ca.ConsPId, negotiation.StateAgreed))
	return negotiation.Ack(neg), nil
}
