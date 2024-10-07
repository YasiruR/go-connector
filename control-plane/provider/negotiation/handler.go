package negotiation

import (
	defaultErr "errors"
	"fmt"
	"github.com/YasiruR/go-connector/domain"
	"github.com/YasiruR/go-connector/domain/api/dsp/http/negotiation"
	"github.com/YasiruR/go-connector/domain/boot"
	"github.com/YasiruR/go-connector/domain/control-plane"
	"github.com/YasiruR/go-connector/domain/errors"
	"github.com/YasiruR/go-connector/domain/models/odrl"
	"github.com/YasiruR/go-connector/domain/pkg"
	"github.com/YasiruR/go-connector/domain/stores"
)

type Handler struct {
	assignerId  string
	cnStore     stores.ContractNegotiationStore
	policyStore stores.OfferStore
	urn         pkg.URNService
	log         pkg.Log
}

func NewHandler(cfg boot.Config, stores domain.Stores, plugins domain.Plugins) *Handler {
	return &Handler{
		assignerId:  cfg.DataSpace.AssignerId,
		cnStore:     stores.ContractNegotiationStore,
		policyStore: stores.OfferStore,
		urn:         plugins.URNService,
		log:         plugins.Log,
	}
}

func (h *Handler) HandleNegotiationsRequest(providerPid string) (negotiation.Ack, error) {
	ack, err := h.cnStore.Negotiation(providerPid)
	if err != nil {
		if defaultErr.Is(err, stores.TypeInvalidKey) {
			return negotiation.Ack{}, errors.Negotiation(providerPid, ``,
				errors.InvalidKey(stores.TypeContractNegotiation, `contract negotiation id`, err))
		}
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	return negotiation.Ack(ack), nil
}

func (h *Handler) HandleContractRequest(cr negotiation.ContractRequest) (ack negotiation.Ack, err error) {
	// associate with existing contract negotiation if providerPid exists and create a new contract
	// negotiation if otherwise
	var cn negotiation.Negotiation
	provPId := cr.ProvPId
	if provPId != `` {
		cn, err = h.cnStore.Negotiation(provPId)
		if err != nil {
			if defaultErr.Is(err, stores.TypeInvalidKey) {
				return negotiation.Ack{}, errors.Negotiation(provPId, cr.ConsPId,
					errors.InvalidKey(stores.TypeContractNegotiation, `contract negotiation id`, err))
			}
			return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
		}

		if cn.State != negotiation.StateOffered {
			return negotiation.Ack{}, errors.Negotiation(cn.ProvPId,
				cn.ConsPId, errors.StateError(`request contract`, string(cn.State)))
		}

		if cn.ConsPId != cr.ConsPId {
			return negotiation.Ack{}, errors.Negotiation(cn.ProvPId, cn.ConsPId,
				errors.InvalidValue(`consumerPid`, cn.ConsPId, cr.ConsPId))
		}

		cn.State = negotiation.StateRequested
		h.log.Debug("a valid contract negotiation exists", cn.ProvPId)
	} else {
		provPId, err = h.urn.NewURN()
		if err != nil {
			return negotiation.Ack{}, errors.PkgError(pkg.TypeURN, `NewURN`, err, `contract negotiation id`)
		}

		cn = negotiation.Negotiation{
			Ctx:     control_plane.Context,
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

	// store (new or updated) contract negotiation, assignee, assigner and its callback address
	h.cnStore.AddNegotiation(provPId, cn)
	h.cnStore.SetParticipants(provPId, cr.CallbackAddr, cr.Offer.Assigner, cr.Offer.Assignee)

	h.log.Trace(fmt.Sprintf("provider stored contract negotiation (id: %s, assigner: %s, assignee: %s, address: %s)",
		provPId, cr.Offer.Assigner, cr.Offer.Assignee, cr.CallbackAddr))
	h.log.Debug(fmt.Sprintf("provider handler updated negotiation state (id: %s, state: %s)",
		provPId, negotiation.StateRequested))

	cn.Type = negotiation.MsgTypeNegotiationAck
	return negotiation.Ack(cn), nil
}

func (h *Handler) HandleAcceptOffer(e negotiation.ContractNegotiationEvent) (negotiation.Ack, error) {
	// validate other attributes (consPid)

	cn, err := h.cnStore.Negotiation(e.ProvPId)
	if err != nil {
		if defaultErr.Is(err, stores.TypeInvalidKey) {
			return negotiation.Ack{}, errors.Negotiation(e.ProvPId, e.ConsPId,
				errors.InvalidKey(stores.TypeContractNegotiation, `contract negotiation id`, err))
		}
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	if cn.State != negotiation.StateOffered {
		return negotiation.Ack{}, errors.Negotiation(cn.ProvPId,
			cn.ConsPId, errors.StateError(`accept offer`, string(cn.State)))
	}

	if err = h.cnStore.UpdateState(e.ProvPId, negotiation.StateAccepted); err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	cn.State = negotiation.StateAccepted
	cn.Type = negotiation.MsgTypeNegotiationAck
	h.log.Debug(fmt.Sprintf("provider handler updated negotiation state (id: %s, state: %s)",
		e.ProvPId, negotiation.StateAccepted))
	return negotiation.Ack(cn), nil
}

func (h *Handler) HandleAgreementVerification(cv negotiation.ContractVerification) (negotiation.Ack, error) {
	// validate message (must contain consumerPid, providerPid)

	cn, err := h.cnStore.Negotiation(cv.ProvPId)
	if err != nil {
		if defaultErr.Is(err, stores.TypeInvalidKey) {
			return negotiation.Ack{}, errors.Negotiation(cv.ProvPId, cv.ConsPId,
				errors.InvalidKey(stores.TypeContractNegotiation, `contract negotiation id`, err))
		}
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	if cn.State != negotiation.StateAgreed {
		return negotiation.Ack{}, errors.Negotiation(cn.ProvPId,
			cn.ConsPId, errors.StateError(`verify agreement`, string(cn.State)))
	}

	if err = h.cnStore.UpdateState(cv.ProvPId, negotiation.StateVerified); err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	cn.State = negotiation.StateVerified
	cn.Type = negotiation.MsgTypeNegotiationAck
	h.log.Debug(fmt.Sprintf("provider handler updated negotiation state (id: %s, state: %s)",
		cv.ProvPId, negotiation.StateVerified))
	return negotiation.Ack(cn), nil
}

func (h *Handler) HandleContractTermination(ct negotiation.ContractTermination) (negotiation.Ack, error) {
	cn, err := h.cnStore.Negotiation(ct.ProvPId)
	if err != nil {
		if defaultErr.Is(err, stores.TypeInvalidKey) {
			return negotiation.Ack{}, errors.Negotiation(ct.ProvPId, ct.ConsPId,
				errors.InvalidKey(stores.TypeContractNegotiation, `contract negotiation id`, err))
		}
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
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
	storedOfr, err := h.policyStore.Offer(receivedOfr.Id)
	if err != nil {
		// todo handle invalid key error
		h.log.Debug(errors.StoreFailed(stores.TypeOffer, `Offer`, err))
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
