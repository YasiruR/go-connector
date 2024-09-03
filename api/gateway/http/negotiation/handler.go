package negotiation

import (
	defaultErr "errors"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/gateway/http/negotiation"
	"github.com/YasiruR/connector/domain/core"
	coreErr "github.com/YasiruR/connector/domain/errors/core"
	"github.com/YasiruR/connector/domain/errors/external"
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
		if defaultErr.Is(err, coreErr.TypeUnmarshalError) {
			middleware.WriteError(w, external.InvalidReqBody(`request contract`,
				err), http.StatusBadRequest)
			return
		}
		middleware.WriteError(w, external.ParseError(`request contract`, err), http.StatusBadRequest)
		return
	}

	cnId, err := h.consumer.RequestContract(req.ConsumerPId, req.ProviderEndpoint, req.OfferId, req.Constraints)
	if err != nil {
		middleware.WriteError(w, coreErr.DSPControllerFailed(core.RoleConsumer, `RequestContract`, err),
			http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, negotiation.ContractRequestResponse{Id: cnId}, http.StatusOK); err != nil {
		middleware.WriteError(w, external.WriteAckError(`request contract`,
			err), http.StatusInternalServerError)
	}
}

func (h *Handler) OfferContract(w http.ResponseWriter, r *http.Request) {
	var req negotiation.OfferRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		if defaultErr.Is(err, coreErr.TypeUnmarshalError) {
			middleware.WriteError(w, external.InvalidReqBody(`offer contract`,
				err), http.StatusBadRequest)
			return
		}
		middleware.WriteError(w, external.ParseError(`offer contract`, err), http.StatusBadRequest)
		return
	}

	cnId, err := h.provider.OfferContract(req.OfferId, req.ProviderPid, req.ConsumerAddr)
	if err != nil {
		middleware.WriteError(w, coreErr.DSPControllerFailed(core.RoleProvider, `OfferContract`, err),
			http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, negotiation.ContractRequestResponse{Id: cnId}, http.StatusOK); err != nil {
		middleware.WriteError(w, external.WriteAckError(`offer contract`,
			err), http.StatusInternalServerError)
	}
}

func (h *Handler) AcceptOffer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	consumerPid, ok := params[negotiation.ParamConsumerPid]
	if !ok {
		middleware.WriteError(w, external.PathParamError(negotiation.ParamConsumerPid), http.StatusBadRequest)
		return
	}

	if err := h.consumer.AcceptOffer(consumerPid); err != nil {
		middleware.WriteError(w, coreErr.DSPControllerFailed(core.RoleConsumer, `AcceptOffer`, err),
			http.StatusBadRequest)
		return
	}

	if err := middleware.WriteAck(w, nil, http.StatusOK); err != nil {
		middleware.WriteError(w, external.WriteAckError(`accept offer`,
			err), http.StatusInternalServerError)
	}
}

func (h *Handler) AgreeContract(w http.ResponseWriter, r *http.Request) {
	var req negotiation.AgreeContractRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		if defaultErr.Is(err, coreErr.TypeUnmarshalError) {
			middleware.WriteError(w, external.InvalidReqBody(`agree contract`,
				err), http.StatusBadRequest)
			return
		}
		middleware.WriteError(w, external.ParseError(`agree contract`, err), http.StatusBadRequest)
		return
	}

	agrId, err := h.provider.AgreeContract(req.OfferId, req.NegotiationId)
	if err != nil {
		middleware.WriteError(w, coreErr.DSPControllerFailed(core.RoleProvider, `AgreeContract`, err),
			http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, negotiation.ContractAgreementResponse{Id: agrId}, http.StatusOK); err != nil {
		middleware.WriteError(w, external.WriteAckError(`accept offer`,
			err), http.StatusInternalServerError)
	}
}

func (h *Handler) GetAgreement(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	agreementId, ok := params[negotiation.ParamAgreementId]
	if !ok {
		middleware.WriteError(w, external.PathParamError(negotiation.ParamConsumerPid), http.StatusBadRequest)
		return
	}

	agr, err := h.agrStore.Agreement(agreementId)
	if err != nil {
		middleware.WriteError(w, external.InvalidKeyError(stores.TypeAgreement,
			`agreement id`, err), http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, agr, http.StatusOK); err != nil {
		middleware.WriteError(w, external.WriteAckError(`get agreement`,
			err), http.StatusInternalServerError)
	}
}

func (h *Handler) VerifyAgreement(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	consumerPid, ok := params[negotiation.ParamConsumerPid]
	if !ok {
		middleware.WriteError(w, external.PathParamError(negotiation.ParamConsumerPid), http.StatusBadRequest)
		return
	}

	if err := h.consumer.VerifyAgreement(consumerPid); err != nil {
		middleware.WriteError(w, coreErr.DSPControllerFailed(core.RoleConsumer, `VerifyAgreement`, err),
			http.StatusBadRequest)
		return
	}

	if err := middleware.WriteAck(w, nil, http.StatusOK); err != nil {
		middleware.WriteError(w, external.WriteAckError(`verify agreement`,
			err), http.StatusInternalServerError)
	}
}

func (h *Handler) FinalizeContract(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	providerPid, ok := params[negotiation.ParamProviderPid]
	if !ok {
		middleware.WriteError(w, external.PathParamError(negotiation.ParamProviderPid), http.StatusBadRequest)
		return
	}

	if err := h.provider.FinalizeContract(providerPid); err != nil {
		middleware.WriteError(w, coreErr.DSPControllerFailed(core.RoleProvider, `FinalizeContract`, err),
			http.StatusBadRequest)
		return
	}

	if err := middleware.WriteAck(w, nil, http.StatusOK); err != nil {
		middleware.WriteError(w, external.WriteAckError(`get agreement`,
			err), http.StatusInternalServerError)
	}
}

func (h *Handler) TerminateContract(w http.ResponseWriter, r *http.Request) {
	var req negotiation.TerminateContractRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		if defaultErr.Is(err, coreErr.TypeUnmarshalError) {
			middleware.WriteError(w, external.InvalidReqBody(`contract termination`,
				err), http.StatusBadRequest)
			return
		}
		middleware.WriteError(w, external.ParseError(`terminate contract`, err), http.StatusBadRequest)
		return
	}

	if req.ProviderPid == `` && req.ConsumerPid != `` {
		if err := h.consumer.TerminateContract(req.ConsumerPid, req.Code, req.Reasons); err != nil {
			middleware.WriteError(w, coreErr.DSPControllerFailed(core.RoleConsumer, `TerminateContract`,
				err), http.StatusBadRequest)
			return
		}

		if err := middleware.WriteAck(w, nil, http.StatusOK); err != nil {
			middleware.WriteError(w, external.WriteAckError(`terminate contract`,
				err), http.StatusInternalServerError)
		}
		return
	}

	if req.ProviderPid != `` && req.ConsumerPid == `` {

	}

	middleware.WriteError(w, external.IncompatibleReqBody(
		`only one of consumer and provider process IDs should be provided`), http.StatusBadRequest)
}
