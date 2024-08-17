package negotiation

import (
	"fmt"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/dsp/http/negotiation"
	"github.com/YasiruR/connector/domain/core"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/stores"
)

type Handler struct {
	cnStore  stores.ContractNegotiation
	agrStore stores.Agreement
	urn      pkg.URNService
	log      pkg.Log
}

func NewHandler(stores domain.Stores, plugins domain.Plugins) *Handler {
	return &Handler{
		cnStore:  stores.ContractNegotiation,
		agrStore: stores.Agreement,
		urn:      plugins.URNService,
		log:      plugins.Log,
	}
}

func (h *Handler) HandleContractOffer(co negotiation.ContractOffer) (ack negotiation.Ack, err error) {
	var cn negotiation.Negotiation
	if co.ConsPId != `` {
		// validate the given consumerPid
		cn, err = h.cnStore.GetNegotiation(co.ConsPId)
		if err != nil {
			return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `GetNegotiation`, err)
		}
		h.log.Trace("a contract negotiation already exists for the contract offer", co.ConsPId)
	} else {
		consumerPid, err := h.urn.NewURN()
		if err != nil {
			return negotiation.Ack{}, errors.PkgFailed(pkg.TypeURN, `NewURN`, err)
		}

		cn.Ctx = core.Context
		cn.ConsPId = consumerPid
		cn.Type = negotiation.MsgTypeNegotiationAck
	}

	cn.ProvPId = co.ProvPId
	cn.State = negotiation.StateOffered
	h.cnStore.Set(cn.ConsPId, cn)
	h.cnStore.SetCallbackAddr(cn.ConsPId, co.CallbackAddr)
	h.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", cn.ConsPId, negotiation.StateOffered))
	return negotiation.Ack(cn), nil
}

func (h *Handler) HandleContractAgreement(ca negotiation.ContractAgreement) (negotiation.Ack, error) {
	// validate agreement (e.g. consumerPid, target)

	h.agrStore.Set(ca.ConsPId, ca.Agreement)
	h.log.Trace(fmt.Sprintf("stored contract agreement (id: %s) for negotation (id: %s)", ca.Agreement.Id, ca.ConsPId))
	h.cnStore.SetCallbackAddr(ca.ConsPId, ca.CallbackAddr)

	if err := h.cnStore.UpdateState(ca.ConsPId, negotiation.StateAgreed); err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	cn, err := h.cnStore.GetNegotiation(ca.ConsPId)
	if err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `GetNegotiation`, err)
	}

	h.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", ca.ConsPId, negotiation.StateAgreed))
	return negotiation.Ack(cn), nil
}

func (h *Handler) HandleFinalizedEvent(consumerPid string) (negotiation.Ack, error) {
	if err := h.cnStore.UpdateState(consumerPid, negotiation.StateFinalized); err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	cn, err := h.cnStore.GetNegotiation(consumerPid)
	if err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `GetNegotiation`, err)
	}

	h.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", consumerPid, negotiation.StateFinalized))
	return negotiation.Ack(cn), nil
}
