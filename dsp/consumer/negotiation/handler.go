package negotiation

import (
	"fmt"
	"github.com/YasiruR/connector/core"
	"github.com/YasiruR/connector/core/dsp/negotiation"
	"github.com/YasiruR/connector/core/errors"
	"github.com/YasiruR/connector/core/pkg"
	"github.com/YasiruR/connector/core/stores"
)

type Handler struct {
	negStore stores.ContractNegotiation
	agrStore stores.Agreement
	log      pkg.Log
}

func NewHandler(stores core.Stores, plugins core.Plugins) *Handler {
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
