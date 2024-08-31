package negotiation

import (
	"fmt"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/gateway/http/negotiation"
	"github.com/YasiruR/connector/domain/core"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/stores"
	"github.com/YasiruR/connector/pkg/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	provider core.Provider
	consumer core.Consumer
	agrStore stores.AgreementStore
	log      pkg.Log
}

func NewHandler(roles domain.Roles, stores domain.Stores, log pkg.Log) *Handler {
	return &Handler{
		provider: roles.Provider,
		consumer: roles.Consumer,
		agrStore: stores.AgreementStore,
		log:      log,
	}
}

func (h *Handler) RequestContract(w http.ResponseWriter, r *http.Request) {
	var req negotiation.ContractRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.ParseRequestFailed(negotiation.RequestContractEndpoint, err), http.StatusBadRequest)
		return
	}

	// todo store offer details when fetched (store dataset id as the target) and provide only the constraints here
	//ofr := odrl.Offer{
	//	Id:   req.OfferId,
	//	Type: odrl.TypeOffer,
	//	//Target:      odrl.Target(req.OdrlTarget),
	//	Assigner:    odrl.Assigner(req.Assigner),
	//	Assignee:    odrl.Assignee(req.Assignee),
	//	Permissions: []odrl.Rule{{Action: odrl.Action(req.Action)}}, // should handle constraints
	//}

	cnId, err := h.consumer.RequestContract(req.ConsumerPId, req.ProviderEndpoint, req.OfferId, req.Constraints)
	if err != nil {
		middleware.WriteError(w, errors.DSPControllerFailed(core.RoleConsumer, `RequestContract`, err), http.StatusBadRequest)
		return
	}

	middleware.WriteAck(w, negotiation.ContractRequestResponse{Id: cnId}, http.StatusOK)
}

func (h *Handler) OfferContract(w http.ResponseWriter, r *http.Request) {
	var req negotiation.OfferRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.ParseRequestFailed(negotiation.OfferContractEndpoint, err), http.StatusBadRequest)
		return
	}

	cnId, err := h.provider.OfferContract(req.OfferId, req.ProviderPid, req.ConsumerAddr)
	if err != nil {
		middleware.WriteError(w, errors.DSPControllerFailed(core.RoleProvider, `OfferContract`, err), http.StatusBadRequest)
		return
	}

	middleware.WriteAck(w, negotiation.ContractRequestResponse{Id: cnId}, http.StatusOK)
}

func (h *Handler) AcceptOffer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	consumerPid, ok := params[negotiation.ParamConsumerPid]
	if !ok {
		middleware.WriteError(w, errors.PathParamNotFound(negotiation.AcceptOfferEndpoint, negotiation.ParamConsumerPid), http.StatusBadRequest)
		return
	}

	if err := h.consumer.AcceptOffer(consumerPid); err != nil {
		middleware.WriteError(w, errors.DSPControllerFailed(core.RoleConsumer, `AcceptOffer`, err), http.StatusBadRequest)
		return
	}

	middleware.WriteAck(w, nil, http.StatusOK)
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

	agr, err := h.agrStore.Agreement(agreementId)
	if err != nil {
		middleware.WriteError(w, errors.StoreFailed(stores.TypeAgreement, `Agreement`, err), http.StatusBadRequest)
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

func (h *Handler) TerminateContract(w http.ResponseWriter, r *http.Request) {
	var req negotiation.TerminateContractRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.ParseRequestFailed(negotiation.TerminateContractEndpoint, err), http.StatusBadRequest)
		return
	}

	if req.ProviderPid == `` && req.ConsumerPid != `` {
		if err := h.consumer.TerminateContract(req.ConsumerPid, req.Code, req.Reasons); err != nil {
			middleware.WriteError(w, errors.DSPControllerFailed(core.RoleConsumer, `TerminateContract`, err), http.StatusBadRequest)
			return
		}
		middleware.WriteAck(w, nil, http.StatusOK)
		return
	}

	if req.ProviderPid != `` && req.ConsumerPid == `` {

	}

	middleware.WriteError(w, errors.InvalidRequestBody(negotiation.TerminateContractEndpoint,
		fmt.Errorf(`only one of consumerPid and providerPid should be provided`)), http.StatusBadRequest)
}
