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
	cnStore  stores.ContractNegotiationStore
	agrStore stores.AgreementStore
	urn      pkg.URNService
	log      pkg.Log
}

func NewHandler(stores domain.Stores, plugins domain.Plugins) *Handler {
	return &Handler{
		cnStore:  stores.ContractNegotiationStore,
		agrStore: stores.AgreementStore,
		urn:      plugins.URNService,
		log:      plugins.Log,
	}
}

func (h *Handler) HandleContractOffer(co negotiation.ContractOffer) (ack negotiation.Ack, err error) {
	var cn negotiation.Negotiation
	if co.ConsPId != `` {
		// validate the given consumerPid
		cn, err = h.cnStore.Negotiation(co.ConsPId)
		if err != nil {
			return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
		}

		if cn.State != negotiation.StateRequested {
			return negotiation.Ack{}, errors.IncompatibleValues(`state`, string(cn.State), string(negotiation.StateRequested))
		}

		if co.ProvPId != cn.ProvPId {
			return negotiation.Ack{}, errors.IncompatibleValues(`providerPid`, cn.ConsPId, co.ConsPId)
		}

		h.log.Trace("a contract negotiation already exists for the contract offer", "id: "+co.ConsPId)
	} else {
		consumerPid, err := h.urn.NewURN()
		if err != nil {
			return negotiation.Ack{}, errors.PkgFailed(pkg.TypeURN, `NewURN`, err)
		}

		cn.Ctx = core.Context
		cn.ConsPId = consumerPid
		cn.Type = negotiation.MsgTypeNegotiation
	}

	cn.ProvPId = co.ProvPId
	cn.State = negotiation.StateOffered
	h.cnStore.AddNegotiation(cn.ConsPId, cn)
	h.cnStore.SetParticipants(cn.ConsPId, co.CallbackAddr, co.Offer.Assigner, co.Offer.Assignee)

	h.log.Trace(fmt.Sprintf("updated callback address for contract negotiation (id: %s, address: %s)", cn.ConsPId, co.CallbackAddr))
	h.log.Debug(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", cn.ConsPId, negotiation.StateOffered))
	cn.Type = negotiation.MsgTypeNegotiationAck
	return negotiation.Ack(cn), nil
}

func (h *Handler) HandleContractAgreement(ca negotiation.ContractAgreement) (negotiation.Ack, error) {
	// validate agreement (e.g. consumerPid, providerPid, target)

	cn, err := h.cnStore.Negotiation(ca.ConsPId)
	if err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	if cn.State != negotiation.StateRequested && cn.State != negotiation.StateAccepted {
		return negotiation.Ack{}, errors.IncompatibleValues(`state`, string(cn.State),
			string(negotiation.StateRequested)+" or "+string(negotiation.StateAccepted))
	}

	h.agrStore.AddAgreement(ca.Agreement.Id, ca.Agreement)
	h.log.Trace(fmt.Sprintf("stored contract agreement (id: %s) for negotation (id: %s)",
		ca.Agreement.Id, ca.ConsPId))

	if err := h.cnStore.UpdateState(ca.ConsPId, negotiation.StateAgreed); err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	cn.Type = negotiation.MsgTypeNegotiationAck
	h.log.Debug(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", ca.ConsPId, negotiation.StateAgreed))
	return negotiation.Ack(cn), nil
}

func (h *Handler) HandleFinalizedEvent(consumerPid string) (negotiation.Ack, error) {
	if err := h.cnStore.UpdateState(consumerPid, negotiation.StateFinalized); err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	cn, err := h.cnStore.Negotiation(consumerPid)
	if err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	h.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", consumerPid, negotiation.StateFinalized))
	return negotiation.Ack(cn), nil
}

func (h *Handler) validAgreement(agr negotiation.ContractAgreement) bool {
	cn, err := h.cnStore.Negotiation(agr.ConsPId)
	if err != nil {
		return false
	}

	if agr.ProvPId != cn.ProvPId {
		return false
	}

	// todo save catalog first, then set assigner, and this
	//if agr.Agreement.Assigner !=
	return true
}
