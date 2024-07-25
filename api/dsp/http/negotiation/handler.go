package negotiation

import (
	"encoding/json"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/dsp/http/catalog"
	negotiation2 "github.com/YasiruR/connector/domain/api/dsp/http/negotiation"
	"github.com/YasiruR/connector/domain/core"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/gorilla/mux"
	"io"
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
	providerPid, ok := params[negotiation2.ParamProviderId]
	if !ok {
		h.sendError(w, errors.PathParamNotFound(negotiation2.RequestEndpoint, negotiation2.ParamProviderId), http.StatusBadRequest)
		return
	}

	neg, err := h.provider.HandleNegotiationsRequest(providerPid)
	if err != nil {
		h.sendError(w, errors.HandlerFailed(negotiation2.RequestEndpoint, core.RoleProvider, err), http.StatusBadRequest)
		return
	}

	h.sendAck(w, negotiation2.RequestEndpoint, neg, http.StatusOK)
}

func (h *Handler) HandleContractRequest(w http.ResponseWriter, r *http.Request) {
	body, err := h.readBody(catalog.RequestEndpoint, w, r)
	if err != nil {
		return
	}

	var req negotiation2.ContractRequest
	if err = json.Unmarshal(body, &req); err != nil {
		h.sendError(w, errors.UnmarshalError(negotiation2.ContractRequestEndpoint, err), http.StatusBadRequest)
		return
	}

	ack, err := h.provider.HandleContractRequest(req)
	if err != nil {
		h.sendError(w, errors.HandlerFailed(negotiation2.ContractRequestEndpoint, core.RoleProvider, err), http.StatusBadRequest)
		return
	}

	h.sendAck(w, negotiation2.ContractRequestEndpoint, ack, http.StatusCreated)
}

func (h *Handler) HandleContractAgreement(w http.ResponseWriter, r *http.Request) {
	body, err := h.readBody(negotiation2.ContractAgreementEndpoint, w, r)
	if err != nil {
		return
	}

	var req negotiation2.ContractAgreement
	if err = json.Unmarshal(body, &req); err != nil {
		h.sendError(w, errors.UnmarshalError(negotiation2.ContractAgreementEndpoint, err), http.StatusBadRequest)
		return
	}

	ack, err := h.consumer.HandleContractAgreement(req)
	if err != nil {
		h.sendError(w, errors.HandlerFailed(negotiation2.ContractAgreementEndpoint, core.RoleConsumer, err), http.StatusBadRequest)
		return
	}

	h.sendAck(w, negotiation2.ContractAgreementEndpoint, ack, http.StatusOK)
}

func (h *Handler) HandleAgreementVerification(w http.ResponseWriter, r *http.Request) {
	body, err := h.readBody(negotiation2.AgreementVerificationEndpoint, w, r)
	if err != nil {
		return
	}

	var req negotiation2.ContractVerification
	if err = json.Unmarshal(body, &req); err != nil {
		h.sendError(w, errors.UnmarshalError(negotiation2.AgreementVerificationEndpoint, err), http.StatusBadRequest)
		return
	}

	ack, err := h.provider.HandleAgreementVerification(req.ProvPId)
	if err != nil {
		h.sendError(w, errors.HandlerFailed(negotiation2.AgreementVerificationEndpoint, core.RoleProvider, err), http.StatusBadRequest)
		return
	}

	h.sendAck(w, negotiation2.AgreementVerificationEndpoint, ack, http.StatusOK)
}

func (h *Handler) HandleEventConsumer(w http.ResponseWriter, r *http.Request) {
	body, err := h.readBody(negotiation2.EventConsumerEndpoint, w, r)
	if err != nil {
		return
	}

	var req negotiation2.ContractNegotiationEvent
	if err = json.Unmarshal(body, &req); err != nil {
		h.sendError(w, errors.UnmarshalError(negotiation2.EventConsumerEndpoint, err), http.StatusBadRequest)
		return
	}

	switch req.EventType {
	case negotiation2.EventFinalized:
		ack, err := h.consumer.HandleFinalizedEvent(req.ConsPId)
		if err != nil {
			h.sendError(w, errors.HandlerFailed(req.ConsPId, core.RoleConsumer, err), http.StatusBadRequest)
			return
		}
		h.sendAck(w, negotiation2.EventConsumerEndpoint, ack, http.StatusOK)
	default:
		h.sendError(w, errors.InvalidRequestBody(negotiation2.EventConsumerEndpoint, err), http.StatusBadRequest)
	}
}

func (h *Handler) readBody(endpoint string, w http.ResponseWriter, r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		err = errors.InvalidRequestBody(endpoint, err)
		w.WriteHeader(http.StatusBadRequest)
		h.log.Error(err)
		r.Body.Close()
		return nil, err
	}
	defer r.Body.Close()

	return body, nil
}

func (h *Handler) sendAck(w http.ResponseWriter, receivedEndpoint string, data any, code int) {
	body, err := json.Marshal(data)
	if err != nil {
		h.sendError(w, errors.MarshalError(receivedEndpoint, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	if _, err = w.Write(body); err != nil {
		h.sendError(w, errors.WriteBodyError(receivedEndpoint, err), http.StatusInternalServerError)
	}
}

func (h *Handler) sendError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	h.log.Error(errors.APIFailed(`dsp`, err))
}
