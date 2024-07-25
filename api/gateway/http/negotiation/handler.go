package negotiation

import (
	"encoding/json"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/dsp/http/negotiation"
	domainNegotiation "github.com/YasiruR/connector/domain/api/gateway/http/negotiation"
	"github.com/YasiruR/connector/domain/core"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/models/odrl"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/stores"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type Handler struct {
	provider core.Provider
	consumer core.Consumer
	agrStore stores.Agreement
	log      pkg.Log
}

func NewHandler(roles domain.Roles, stores domain.Stores, log pkg.Log) *Handler {
	return &Handler{
		provider: roles.Provider,
		consumer: roles.Consumer,
		agrStore: stores.Agreement,
		log:      log,
	}
}

func (h *Handler) RequestContract(w http.ResponseWriter, r *http.Request) {
	body, err := h.readBody(domainNegotiation.RequestContractEndpoint, w, r)
	if err != nil {
		return
	}

	var req domainNegotiation.ContractRequest
	if err = json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Error(errors.UnmarshalError(domainNegotiation.RequestContractEndpoint, err))
		return
	}

	ofr := odrl.Offer{
		Id:          req.OfferId,
		Type:        odrl.TypeOffer,
		Target:      odrl.Target(req.OdrlTarget),
		Assigner:    odrl.Assigner(req.Assigner),
		Assignee:    odrl.Assignee(req.Assignee),
		Permissions: []odrl.Rule{{Action: odrl.Action(req.Action)}}, // should handle constraints
	}

	negId, err := h.consumer.RequestContract(req.ProviderEndpoint, ofr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Error(errors.DSPFailed(core.RoleConsumer, `RequestContract`, err))
		return
	}

	h.sendAck(w, domainNegotiation.RequestContractEndpoint, domainNegotiation.ContractRequestResponse{Id: negId}, http.StatusOK)
}

func (h *Handler) AgreeContract(w http.ResponseWriter, r *http.Request) {
	body, err := h.readBody(domainNegotiation.AgreeContractEndpoint, w, r)
	if err != nil {
		return
	}

	var req domainNegotiation.AgreeContractRequest
	if err = json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Error(errors.UnmarshalError(domainNegotiation.AgreeContractEndpoint, err))
		return
	}

	agrId, err := h.provider.AgreeContract(req.OfferId, req.NegotiationId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Error(errors.DSPFailed(core.RoleProvider, `AgreeContract`, err))
		return
	}

	h.sendAck(w, domainNegotiation.AgreeContractEndpoint, domainNegotiation.ContractAgreementResponse{Id: agrId}, http.StatusOK)
}

func (h *Handler) GetAgreement(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	agreementId, ok := params[domainNegotiation.ParamAgreementId]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Error(errors.PathParamNotFound(domainNegotiation.GetAgreementEndpoint, negotiation.ParamConsumerPid))
		return
	}

	agr, err := h.agrStore.Get(agreementId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Error(errors.StoreFailed(stores.TypeAgreement, `Get`, err))
		return
	}

	h.sendAck(w, domainNegotiation.GetAgreementEndpoint, agr, http.StatusOK)
}

func (h *Handler) VerifyAgreement(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	consumerPid, ok := params[domainNegotiation.ParamConsumerPid]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Error(errors.PathParamNotFound(domainNegotiation.VerifyAgreementEndpoint, negotiation.ParamConsumerPid))
		return
	}

	if err := h.consumer.VerifyAgreement(consumerPid); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Error(errors.DSPFailed(core.RoleConsumer, `VerifyAgreement`, err))
		return
	}

	h.sendAck(w, domainNegotiation.VerifyAgreementEndpoint, nil, http.StatusOK)
}

func (h *Handler) FinalizeContract(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	providerPid, ok := params[domainNegotiation.ParamProviderPid]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Error(errors.PathParamNotFound(domainNegotiation.FinalizeContractEndpoint, negotiation.ParamProviderId))
		return
	}

	if err := h.provider.FinalizeContract(providerPid); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Error(errors.DSPFailed(core.RoleProvider, `FinalizeContract`, err))
		return
	}

	h.sendAck(w, domainNegotiation.FinalizeContractEndpoint, nil, http.StatusOK)
}

func (h *Handler) readBody(endpoint string, w http.ResponseWriter, r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		err = errors.InvalidRequestBody(endpoint, err)
		w.WriteHeader(http.StatusBadRequest)
		r.Body.Close()
		return nil, err
	}
	defer r.Body.Close()

	return body, nil
}

func (h *Handler) sendAck(w http.ResponseWriter, receivedEndpoint string, data any, code int) {
	body, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(errors.MarshalError(receivedEndpoint, err))
		return
	}

	w.WriteHeader(code)
	if _, err = w.Write(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(errors.WriteBodyError(receivedEndpoint, err))
	}
}
