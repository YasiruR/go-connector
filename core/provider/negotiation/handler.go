package negotiation

import (
	"fmt"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/dsp/http/negotiation"
	"github.com/YasiruR/connector/domain/core"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/models/odrl"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/stores"
)

type Handler struct {
	assignerId  string
	cnStore     stores.ContractNegotiation
	policyStore stores.Policy
	urn         pkg.URNService
	log         pkg.Log
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
		h.log.Debug("a valid contract negotiation exists", cn.ProvPId)
	} else {
		provPId, err = h.urn.NewURN()
		if err != nil {
			return negotiation.Ack{}, errors.URNFailed(`providerPid`, `NewURN`, err)
		}

		cn = negotiation.Negotiation{
			Ctx:     core.Context,
			Type:    negotiation.MsgTypeNegotiation,
			ConsPId: cr.ConsPId,
			ProvPId: provPId,
			State:   negotiation.StateRequested,
		}
	}

	// return error message if the offer is invalid
	if !h.validOffer(cr.Offer) {
		return negotiation.Ack{}, fmt.Errorf("received an invalid offer")
	}

	// return error message if callback address is invalid
	if !h.validAddress(cr.CallbackAddr) {
		return negotiation.Ack{}, fmt.Errorf("received an invalid callback address")
	}

	// store (new or updated) contract negotiation, assignee and its callback address
	h.cnStore.Set(provPId, cn)
	h.cnStore.SetAssignee(provPId, cr.Offer.Assignee)
	h.cnStore.SetCallbackAddr(provPId, cr.CallbackAddr)
	h.log.Trace("stored contract negotiation", "id: "+provPId, "assigner: "+cr.Offer.Assigner,
		"assignee: "+cr.Offer.Assignee, "address: "+cr.CallbackAddr)
	h.log.Debug("updated negotiation state", "id: "+provPId, "state: "+negotiation.StateRequested)

	cn.Type = negotiation.MsgTypeNegotiationAck
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

// todo move policy engine validator as a package
func (h *Handler) validOffer(receivedOfr odrl.Offer) bool {
	storedOfr, err := h.policyStore.GetOffer(receivedOfr.Id)
	if err != nil {
		h.log.Debug(errors.StoreFailed(stores.TypePolicy, `GetOffer`, err))
		return false
	}

	if receivedOfr.Assigner != storedOfr.Assigner {
		h.log.Debug("assigner in the received offer did not match with the stored offer",
			"received:"+receivedOfr.Assigner, "stored:"+storedOfr.Assigner)
		return false
	}

	// validate rules
	return true
}

func (h *Handler) validAddress(addr string) bool {
	return true
}
