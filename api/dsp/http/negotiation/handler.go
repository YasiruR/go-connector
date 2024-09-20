package negotiation

import (
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api"
	"github.com/YasiruR/connector/domain/api/dsp/http/negotiation"
	"github.com/YasiruR/connector/domain/core"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/pkg/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	provider core.Provider
	consumer core.Consumer
	log      pkg.Log
}

func NewHandler(roles domain.Roles, log pkg.Log) *Handler {
	return &Handler{
		provider: roles.Provider,
		consumer: roles.Consumer,
		log:      log,
	}
}

func (h *Handler) GetNegotiation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	providerPid, ok := params[api.ParamProviderPid]
	if !ok {
		middleware.WriteError(w, errors.Negotiation(``, ``,
			errors.PathParamNotFound(api.ParamProviderPid)), http.StatusBadRequest)
		return
	}

	neg, err := h.provider.HandleNegotiationsRequest(providerPid)
	if err != nil {
		middleware.WriteError(w, errors.DSPHandlerFailed(core.RoleProvider, negotiation.RequestEndpoint,
			err), http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, neg, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Negotiation(neg.ProvPId, neg.ConsPId,
			errors.WriteAckError(`get negotiations`, err)), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleContractRequest(w http.ResponseWriter, r *http.Request) {
	var endpoint string
	_, ok := mux.Vars(r)[api.ParamConsumerPid]
	if ok {
		endpoint = negotiation.ContractOfferToRequestEndpoint
	} else {
		endpoint = negotiation.ContractOfferEndpoint
	}

	var req negotiation.ContractRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.Negotiation(``, ``,
			errors.InvalidReqBody(`contract request`, err)), http.StatusBadRequest)
		return
	}

	ack, err := h.provider.HandleContractRequest(req)
	if err != nil {
		middleware.WriteError(w, errors.DSPHandlerFailed(core.RoleProvider, endpoint, err), http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, ack, http.StatusCreated); err != nil {
		middleware.WriteError(w, errors.Negotiation(ack.ProvPId, ack.ConsPId,
			errors.WriteAckError(`contract request`, err)), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleContractOffer(w http.ResponseWriter, r *http.Request) {
	var endpoint string
	_, ok := mux.Vars(r)[api.ParamConsumerPid]
	if ok {
		endpoint = negotiation.ContractOfferToRequestEndpoint
	} else {
		endpoint = negotiation.ContractOfferEndpoint
	}

	var req negotiation.ContractOffer
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.Negotiation(``, ``,
			errors.InvalidReqBody(`contract offer`, err)), http.StatusBadRequest)
		return
	}

	ack, err := h.consumer.HandleContractOffer(req)
	if err != nil {
		middleware.WriteError(w, errors.DSPHandlerFailed(core.RoleConsumer, endpoint, err), http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, ack, http.StatusCreated); err != nil {
		middleware.WriteError(w, errors.Negotiation(ack.ProvPId, ack.ConsPId,
			errors.WriteAckError(`contract offer`, err)), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleContractAgreement(w http.ResponseWriter, r *http.Request) {
	var req negotiation.ContractAgreement
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.Negotiation(``, ``,
			errors.InvalidReqBody(`contract agreement`, err)), http.StatusBadRequest)
		return
	}

	ack, err := h.consumer.HandleContractAgreement(req)
	if err != nil {
		middleware.WriteError(w, errors.DSPHandlerFailed(core.RoleConsumer, negotiation.ContractAgreementEndpoint,
			err), http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, ack, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Negotiation(ack.ProvPId, ack.ConsPId,
			errors.WriteAckError(`contract agreement`, err)), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleAgreementVerification(w http.ResponseWriter, r *http.Request) {
	var req negotiation.ContractVerification
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.Negotiation(``, ``,
			errors.InvalidReqBody(`agreement verification`, err)), http.StatusBadRequest)
		return
	}

	ack, err := h.provider.HandleAgreementVerification(req)
	if err != nil {
		middleware.WriteError(w, errors.DSPHandlerFailed(core.RoleProvider,
			negotiation.AgreementVerificationEndpoint, err), http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, ack, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Negotiation(ack.ProvPId, ack.ConsPId,
			errors.WriteAckError(`agreement verification`, err)), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleNegotiationEvent(w http.ResponseWriter, r *http.Request) {
	var req negotiation.ContractNegotiationEvent
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.Negotiation(``, ``,
			errors.InvalidReqBody(`negotiation event`, err)), http.StatusBadRequest)
		return
	}

	var ack negotiation.Ack
	var err error

	switch req.EventType {
	case negotiation.EventAccepted:
		ack, err = h.provider.HandleAcceptOffer(req)
		if err != nil {
			middleware.WriteError(w, errors.DSPHandlerFailed(core.RoleProvider, negotiation.EventsEndpoint,
				err), http.StatusBadRequest)
			return
		}
	case negotiation.EventFinalized:
		ack, err = h.consumer.HandleFinalizedEvent(req)
		if err != nil {
			middleware.WriteError(w, errors.DSPHandlerFailed(core.RoleConsumer, negotiation.EventsEndpoint,
				err), http.StatusBadRequest)
			return
		}
	default:
		middleware.WriteError(w, errors.Negotiation(``, ``,
			errors.IncorrectReqValues(`unsupported event type`)), http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, ack, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Negotiation(ack.ProvPId, ack.ConsPId,
			errors.WriteAckError(`negotiation event`, err)), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleTermination(w http.ResponseWriter, r *http.Request) {
	var req negotiation.ContractTermination
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.Negotiation(``, ``,
			errors.InvalidReqBody(`contract termination`, err)), http.StatusBadRequest)
		return
	}

	ack, err := h.provider.HandleContractTermination(req)
	if err != nil {
		middleware.WriteError(w, errors.DSPHandlerFailed(core.RoleProvider, negotiation.TerminateEndpoint,
			err), http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, ack, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Negotiation(ack.ProvPId, ack.ConsPId,
			errors.WriteAckError(`contract termination`, err)), http.StatusInternalServerError)
	}
}
