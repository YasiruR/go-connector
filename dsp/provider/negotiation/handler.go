package negotiation

import (
	"fmt"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/protocols/dsp"
	negotiation2 "github.com/YasiruR/connector/domain/protocols/dsp/negotiation"
	"github.com/YasiruR/connector/domain/stores"
)

type Handler struct {
	cnStore stores.ContractNegotiation
	urn     pkg.URNService
	log     pkg.Log
}

func NewHandler(cnStore stores.ContractNegotiation, plugins domain.Plugins) *Handler {
	return &Handler{
		cnStore: cnStore,
		urn:     plugins.URNService,
		log:     plugins.Log,
	}
}

func (h *Handler) HandleNegotiationsRequest(providerPid string) (negotiation2.Ack, error) {
	ack, err := h.cnStore.Negotiation(providerPid)
	if err != nil {
		return negotiation2.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	return negotiation2.Ack(ack), nil
}

func (h *Handler) HandleContractRequest(cr negotiation2.ContractRequest) (ack negotiation2.Ack, err error) {
	// return error message if offerId is invalid

	// return error message if callbackAddress is invalid

	// associate with existing contract negotiation if providerPid exists and create a new contract
	// negotiation if otherwise
	var cn negotiation2.Negotiation
	provPId := cr.ProvPId
	if provPId != `` {
		cn, err = h.cnStore.Negotiation(provPId)
		if err != nil {
			return negotiation2.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
		}

		if cn.State != negotiation2.StateOffered {
			return negotiation2.Ack{}, errors.IncompatibleValues(`state`, string(cn.State), string(negotiation2.StateOffered))
		}

		if cn.ConsPId != cr.ConsPId {
			return negotiation2.Ack{}, errors.IncompatibleValues(`consumerPid`, cn.ConsPId, cr.ConsPId)
		}

		cn.State = negotiation2.StateRequested
		cn.Type = negotiation2.TypeNegotiationAck
		h.log.Trace("a valid contract negotiation exists", cn.ProvPId)
	} else {
		provPId, err = h.urn.NewURN()
		if err != nil {
			return negotiation2.Ack{}, errors.URNFailed(`providerPid`, `NewURN`, err)
		}

		cn = negotiation2.Negotiation{
			Ctx:     dsp.Context,
			Type:    negotiation2.TypeNegotiationAck,
			ConsPId: cr.ConsPId,
			ProvPId: provPId,
			State:   negotiation2.StateRequested,
		}
	}

	h.cnStore.Set(provPId, cn)
	h.cnStore.SetAssignee(provPId, cr.Offer.Assignee)
	h.cnStore.SetCallbackAddr(provPId, cr.CallbackAddr)
	h.log.Trace(fmt.Sprintf("stored contract negotiation (assigner: %s, assignee: %s)", cr.Offer.Assigner, cr.Offer.Assignee), cn)
	h.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", provPId, negotiation2.StateRequested))
	return negotiation2.Ack(cn), nil
}

func (h *Handler) HandleAgreementVerification(providerPid string) (negotiation2.Ack, error) {
	cn, err := h.cnStore.Negotiation(providerPid)
	if err != nil {
		return negotiation2.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	if err = h.cnStore.UpdateState(providerPid, negotiation2.StateVerified); err != nil {
		return negotiation2.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	cn.State = negotiation2.StateVerified
	cn.Type = negotiation2.TypeNegotiationAck
	h.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", providerPid, negotiation2.StateVerified))
	return negotiation2.Ack(cn), nil
}
