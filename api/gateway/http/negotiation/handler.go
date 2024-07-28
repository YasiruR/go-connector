package negotiation

import (
	"github.com/YasiruR/connector/api/gateway/http/middleware"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/gateway/http/negotiation"
	"github.com/YasiruR/connector/domain/core"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/models/odrl"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/stores"
	"github.com/gorilla/mux"
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
	var req negotiation.ContractRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.ParseRequestFailed(negotiation.RequestContractEndpoint, err), http.StatusBadRequest)
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
		middleware.WriteError(w, errors.DSPControllerFailed(core.RoleConsumer, `RequestContract`, err), http.StatusBadRequest)
		return
	}

	middleware.WriteAck(w, negotiation.ContractRequestResponse{Id: negId}, http.StatusOK)
}

func (h *Handler) AgreeContract(w http.ResponseWriter, r *http.Request) {
	var req negotiation.AgreeContractRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.ParseRequestFailed(negotiation.AgreeContractEndpoint, err), http.StatusBadRequest)
		return
	}

	agrId, err := h.provider.AgreeContract(req.OfferId, req.NegotiationId)
	if err != nil {
		middleware.WriteError(w, errors.DSPControllerFailed(core.RoleProvider, `AgreeContract`, err), http.StatusBadRequest)
		return
	}

	middleware.WriteAck(w, negotiation.ContractAgreementResponse{Id: agrId}, http.StatusOK)
}

func (h *Handler) GetAgreement(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	agreementId, ok := params[negotiation.ParamAgreementId]
	if !ok {
		middleware.WriteError(w, errors.PathParamNotFound(negotiation.GetAgreementEndpoint, negotiation.ParamConsumerPid), http.StatusBadRequest)
		return
	}

	agr, err := h.agrStore.Get(agreementId)
	if err != nil {
		middleware.WriteError(w, errors.StoreFailed(stores.TypeAgreement, `Get`, err), http.StatusBadRequest)
		return
	}

	middleware.WriteAck(w, agr, http.StatusOK)
}

func (h *Handler) VerifyAgreement(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	consumerPid, ok := params[negotiation.ParamConsumerPid]
	if !ok {
		middleware.WriteError(w, errors.PathParamNotFound(negotiation.VerifyAgreementEndpoint, negotiation.ParamConsumerPid), http.StatusBadRequest)
		return
	}

	if err := h.consumer.VerifyAgreement(consumerPid); err != nil {
		middleware.WriteError(w, errors.DSPControllerFailed(core.RoleConsumer, `VerifyAgreement`, err), http.StatusBadRequest)
		return
	}

	middleware.WriteAck(w, nil, http.StatusOK)
}

func (h *Handler) FinalizeContract(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	providerPid, ok := params[negotiation.ParamProviderPid]
	if !ok {
		middleware.WriteError(w, errors.PathParamNotFound(negotiation.FinalizeContractEndpoint, negotiation.ParamProviderPid), http.StatusBadRequest)
		return
	}

	if err := h.provider.FinalizeContract(providerPid); err != nil {
		middleware.WriteError(w, errors.DSPControllerFailed(core.RoleProvider, `FinalizeContract`, err), http.StatusBadRequest)
		return
	}

	middleware.WriteAck(w, nil, http.StatusOK)
}
