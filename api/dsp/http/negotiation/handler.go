package negotiation

import (
	"github.com/YasiruR/connector/api/dsp/http/middleware"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/dsp/http/negotiation"
	"github.com/YasiruR/connector/domain/core"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/pkg"
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
	providerPid, ok := params[negotiation.ParamProviderId]
	if !ok {
		middleware.WriteError(w, errors.PathParamNotFound(negotiation.RequestEndpoint, negotiation.ParamProviderId), http.StatusBadRequest)
		return
	}

	neg, err := h.provider.HandleNegotiationsRequest(providerPid)
	if err != nil {
		middleware.WriteError(w, errors.DSPHandlerFailed(core.RoleProvider, negotiation.RequestEndpoint, err), http.StatusBadRequest)
		return
	}

	middleware.WriteAck(w, neg, http.StatusOK)
}

func (h *Handler) HandleContractRequest(w http.ResponseWriter, r *http.Request) {
	var req negotiation.ContractRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.ParseRequestFailed(negotiation.ContractRequestEndpoint, err), http.StatusBadRequest)
		return
	}

	ack, err := h.provider.HandleContractRequest(req)
	if err != nil {
		middleware.WriteError(w, errors.DSPHandlerFailed(core.RoleProvider, negotiation.ContractRequestEndpoint, err), http.StatusBadRequest)
		return
	}

	middleware.WriteAck(w, ack, http.StatusCreated)
}

func (h *Handler) HandleContractOffer(w http.ResponseWriter, r *http.Request) {
	var endpoint string
	_, ok := mux.Vars(r)[negotiation.ParamConsumerPid]
	if ok {
		endpoint = negotiation.ContractOfferToRequestEndpoint
	} else {
		endpoint = negotiation.ContractOfferEndpoint
	}

	var req negotiation.ContractOffer
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.ParseRequestFailed(endpoint, err), http.StatusBadRequest)
		return
	}

	ack, err := h.consumer.HandleContractOffer(req)
	if err != nil {
		middleware.WriteError(w, errors.DSPHandlerFailed(core.RoleConsumer, endpoint, err), http.StatusBadRequest)
		return
	}

	middleware.WriteAck(w, ack, http.StatusCreated)
}

func (h *Handler) HandleContractAgreement(w http.ResponseWriter, r *http.Request) {
	var req negotiation.ContractAgreement
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.ParseRequestFailed(negotiation.ContractAgreementEndpoint, err), http.StatusBadRequest)
		return
	}

	ack, err := h.consumer.HandleContractAgreement(req)
	if err != nil {
		middleware.WriteError(w, errors.DSPHandlerFailed(core.RoleConsumer, negotiation.ContractAgreementEndpoint, err), http.StatusBadRequest)
		return
	}

	middleware.WriteAck(w, ack, http.StatusOK)
}

func (h *Handler) HandleAgreementVerification(w http.ResponseWriter, r *http.Request) {
	var req negotiation.ContractVerification
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.ParseRequestFailed(negotiation.AgreementVerificationEndpoint, err), http.StatusBadRequest)
		return
	}

	ack, err := h.provider.HandleAgreementVerification(req.ProvPId)
	if err != nil {
		middleware.WriteError(w, errors.DSPHandlerFailed(core.RoleProvider, negotiation.AgreementVerificationEndpoint, err), http.StatusBadRequest)
		return
	}

	middleware.WriteAck(w, ack, http.StatusOK)
}

func (h *Handler) HandleEventConsumer(w http.ResponseWriter, r *http.Request) {
	var req negotiation.ContractNegotiationEvent
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.ParseRequestFailed(negotiation.EventConsumerEndpoint, err), http.StatusBadRequest)
		return
	}

	switch req.EventType {
	case negotiation.EventFinalized:
		ack, err := h.consumer.HandleFinalizedEvent(req.ConsPId)
		if err != nil {
			middleware.WriteError(w, errors.DSPHandlerFailed(core.RoleConsumer, req.ConsPId, err), http.StatusBadRequest)
			return
		}

		middleware.WriteAck(w, ack, http.StatusOK)
	default:
		middleware.WriteError(w, errors.InvalidRequestBody(negotiation.EventConsumerEndpoint, nil), http.StatusBadRequest)
	}
}
