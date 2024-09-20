package negotiation

import (
	defaultErr "errors"
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
			if defaultErr.Is(err, stores.TypeInvalidKey) {
				return negotiation.Ack{}, errors.Negotiation(``, co.ConsPId,
					errors.InvalidKey(stores.TypeContractNegotiation, `contract negotiation id`, err))
			}
			return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
		}

		if co.ProvPId != cn.ProvPId {
			return negotiation.Ack{}, errors.Negotiation(cn.ProvPId, co.ConsPId,
				errors.InvalidValue(`providerPid`, cn.ProvPId, co.ProvPId))
		}

		if cn.State != negotiation.StateRequested {
			return negotiation.Ack{}, errors.Negotiation(co.ProvPId, co.ConsPId,
				errors.StateError(`offer contract`, string(cn.State)))
		}

		h.log.Trace("a contract negotiation already exists for the contract offer", "id: "+co.ConsPId)
	} else {
		consumerPid, err := h.urn.NewURN()
		if err != nil {
			return negotiation.Ack{}, errors.PkgError(pkg.TypeURN, `NewURN`, err, `contract negotiation id`)
		}

		cn.Ctx = core.Context
		cn.ConsPId = consumerPid
		cn.Type = negotiation.MsgTypeNegotiation
	}

	cn.ProvPId = co.ProvPId
	cn.State = negotiation.StateOffered
	h.cnStore.AddNegotiation(cn.ConsPId, cn)
	h.cnStore.SetParticipants(cn.ConsPId, co.CallbackAddr, co.Offer.Assigner, co.Offer.Assignee)

	h.log.Trace(fmt.Sprintf("consumer updated callback address for contract negotiation (id: %s, address: %s)",
		cn.ConsPId, co.CallbackAddr))
	h.log.Debug(fmt.Sprintf("consumer handler updated negotiation state (id: %s, state: %s)", cn.ConsPId,
		negotiation.StateOffered))
	cn.Type = negotiation.MsgTypeNegotiationAck
	return negotiation.Ack(cn), nil
}

func (h *Handler) HandleContractAgreement(ca negotiation.ContractAgreement) (negotiation.Ack, error) {
	// validate agreement (e.g. consumerPid, providerPid, target)

	cn, err := h.cnStore.Negotiation(ca.ConsPId)
	if err != nil {
		if defaultErr.Is(err, stores.TypeInvalidKey) {
			return negotiation.Ack{}, errors.Negotiation(``, ca.ConsPId,
				errors.InvalidKey(stores.TypeContractNegotiation, `contract negotiation id`, err))
		}
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	if cn.State != negotiation.StateRequested && cn.State != negotiation.StateAccepted {
		return negotiation.Ack{}, errors.Negotiation(ca.ProvPId, ca.ConsPId,
			errors.StateError(`agree contract`, string(cn.State)))
	}

	h.agrStore.AddAgreement(ca.Agreement.Id, ca.Agreement)
	h.log.Trace(fmt.Sprintf("consumer stored contract agreement (id: %s) for negotation (id: %s)",
		ca.Agreement.Id, ca.ConsPId))

	if err := h.cnStore.UpdateState(ca.ConsPId, negotiation.StateAgreed); err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	cn.Type = negotiation.MsgTypeNegotiationAck
	h.log.Debug(fmt.Sprintf("consumer handler updated negotiation state (id: %s, state: %s)",
		ca.ConsPId, negotiation.StateAgreed))
	return negotiation.Ack(cn), nil
}

func (h *Handler) HandleFinalizedEvent(e negotiation.ContractNegotiationEvent) (negotiation.Ack, error) {
	if err := h.cnStore.UpdateState(e.ConsPId, negotiation.StateFinalized); err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	cn, err := h.cnStore.Negotiation(e.ConsPId)
	if err != nil {
		if defaultErr.Is(err, stores.TypeInvalidKey) {
			return negotiation.Ack{}, errors.Negotiation(``, e.ConsPId,
				errors.InvalidKey(stores.TypeContractNegotiation, `contract negotiation id`, err))
		}
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	h.log.Info(fmt.Sprintf("consumer handler updated negotiation state (id: %s, state: %s)",
		e.ConsPId, negotiation.StateFinalized))
	return negotiation.Ack(cn), nil
}

func (h *Handler) validAgreement(agr negotiation.ContractAgreement) bool {
	cn, err := h.cnStore.Negotiation(agr.ConsPId)
	if err != nil {
		// todo handle invalid key error
		return false
	}

	if agr.ProvPId != cn.ProvPId {
		return false
	}

	// todo save catalog first, then set assigner, and this
	//if agr.Agreement.Assigner !=
	return true
}
