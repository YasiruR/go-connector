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

func (h *Handler) HandleNegotiationsRequest(providerPid string) (negotiation.Ack, error) {
	ack, err := h.cnStore.GetNegotiation(providerPid)
	if err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `GetNegotiation`, err)
	}

	return negotiation.Ack(ack), nil
}

func (h *Handler) HandleContractRequest(cr negotiation.ContractRequest) (ack negotiation.Ack, err error) {
	// return error message if offerId is invalid

	// return error message if callbackAddress is invalid

	// associate with existing contract negotiation if providerPid exists and create a new contract
	// negotiation if otherwise
	var cn negotiation.Negotiation
	provPId := cr.ProvPId
	if provPId != `` {
		cn, err = h.cnStore.GetNegotiation(provPId)
		if err != nil {
			return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `GetNegotiation`, err)
		}

		if cn.State != negotiation.StateOffered {
			return negotiation.Ack{}, errors.IncompatibleValues(`state`, string(cn.State), string(negotiation.StateOffered))
		}

		if cn.ConsPId != cr.ConsPId {
			return negotiation.Ack{}, errors.IncompatibleValues(`consumerPid`, cn.ConsPId, cr.ConsPId)
		}

		cn.State = negotiation.StateRequested
		cn.Type = negotiation.MsgTypeNegotiationAck
		h.log.Debug("a valid contract negotiation exists", cn.ProvPId)
	} else {
		provPId, err = h.urn.NewURN()
		if err != nil {
			return negotiation.Ack{}, errors.URNFailed(`providerPid`, `NewURN`, err)
		}

		cn = negotiation.Negotiation{
			Ctx:     core.Context,
			Type:    negotiation.MsgTypeNegotiationAck,
			ConsPId: cr.ConsPId,
			ProvPId: provPId,
			State:   negotiation.StateRequested,
		}
	}

	h.cnStore.Set(provPId, cn)
	h.cnStore.SetAssignee(provPId, cr.Offer.Assignee)
	h.cnStore.SetCallbackAddr(provPId, cr.CallbackAddr)
	h.log.Trace(fmt.Sprintf("stored contract negotiation (assigner: %s, assignee: %s)", cr.Offer.Assigner, cr.Offer.Assignee), cn)
	h.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", provPId, negotiation.StateRequested))
	return negotiation.Ack(cn), nil
}

func (h *Handler) HandleAcceptOffer(providerPid string) (negotiation.Ack, error) {
	cn, err := h.cnStore.GetNegotiation(providerPid)
	if err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `GetNegotiation`, err)
	}

	if cn.State != negotiation.StateOffered {
		return negotiation.Ack{}, errors.IncompatibleValues(`state`, string(cn.State), string(negotiation.StateOffered))
	}

	if err = h.cnStore.UpdateState(providerPid, negotiation.StateAccepted); err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	h.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", providerPid, negotiation.StateAccepted))
	return negotiation.Ack(cn), nil
}

func (h *Handler) HandleAgreementVerification(providerPid string) (negotiation.Ack, error) {
	cn, err := h.cnStore.GetNegotiation(providerPid)
	if err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `GetNegotiation`, err)
	}

	if err = h.cnStore.UpdateState(providerPid, negotiation.StateVerified); err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	cn.State = negotiation.StateVerified
	cn.Type = negotiation.MsgTypeNegotiationAck // todo check if all stored negotiations have negotiations msg type
	h.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", providerPid, negotiation.StateVerified))
	return negotiation.Ack(cn), nil
}

func (h *Handler) HandleTermination(ct negotiation.ContractTermination) (negotiation.Ack, error) {
	cn, err := h.cnStore.GetNegotiation(ct.ProvPId)
	if err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `GetNegotiation`, err)
	}

	// can clear stores instead of this
	if err = h.cnStore.UpdateState(ct.ProvPId, negotiation.StateTerminated); err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	cn.State = negotiation.StateTerminated
	cn.Type = negotiation.MsgTypeNegotiationAck
	h.log.Info("consumer terminated the negotiation flow", ct.ProvPId)
	return negotiation.Ack(cn), nil
}

// validate
// - state, consPid
