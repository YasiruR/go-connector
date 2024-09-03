package transfer

import (
	defaultErr "errors"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api"
	"github.com/YasiruR/connector/domain/api/gateway/http/transfer"
	"github.com/YasiruR/connector/domain/core"
	coreErr "github.com/YasiruR/connector/domain/errors/core"
	"github.com/YasiruR/connector/domain/errors/external"
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

func (h *Handler) GetProviderProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tpId, ok := vars[api.ParamConsumerPid]
	if !ok {
		middleware.WriteError(w, external.PathParamError(api.ParamConsumerPid), http.StatusBadRequest)
		return
	}

	tp, err := h.consumer.GetProviderProcess(tpId)
	if err != nil {
		middleware.WriteError(w, coreErr.DSPControllerFailed(core.RoleConsumer, `GetProviderProcess`, err), http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, tp, http.StatusOK); err != nil {
		middleware.WriteError(w, external.WriteAckError(`get provider process`,
			err), http.StatusInternalServerError)
	}
}

func (h *Handler) RequestTransfer(w http.ResponseWriter, r *http.Request) {
	var req transfer.Request
	if err := middleware.ParseRequest(r, &req); err != nil {
		if defaultErr.Is(err, coreErr.TypeUnmarshalError) {
			middleware.WriteError(w, external.InvalidReqBody(`request transfer`,
				err), http.StatusBadRequest)
			return
		}
		middleware.WriteError(w, external.ParseError(`request transfer`, err), http.StatusBadRequest)
		return
	}

	trId, err := h.consumer.RequestTransfer(req.TransferFormat, req.AgreementId, req.SinkEndpoint, req.ProviderEndpoint)
	if err != nil {
		middleware.WriteError(w, coreErr.DSPControllerFailed(core.RoleConsumer, `RequestTransfer`, err), http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, transfer.Response{TransferID: trId}, http.StatusOK); err != nil {
		middleware.WriteError(w, external.WriteAckError(`request transfer`,
			err), http.StatusInternalServerError)
	}
}

func (h *Handler) StartTransfer(w http.ResponseWriter, r *http.Request) {
	var req transfer.StartRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		if defaultErr.Is(err, coreErr.TypeUnmarshalError) {
			middleware.WriteError(w, external.InvalidReqBody(`start transfer`,
				err), http.StatusBadRequest)
			return
		}
		middleware.WriteError(w, external.ParseError(`start transfer`, err), http.StatusBadRequest)
		return
	}

	if req.Provider {
		if err := h.provider.StartTransfer(req.TransferId, req.SourceEndpoint); err != nil {
			middleware.WriteError(w, coreErr.DSPControllerFailed(core.RoleProvider,
				`StartTransfer`, err), http.StatusBadRequest)
			return
		}
	} else {
		if err := h.consumer.StartTransfer(req.TransferId); err != nil {
			middleware.WriteError(w, coreErr.DSPControllerFailed(core.RoleConsumer,
				`StartTransfer`, err), http.StatusBadRequest)
			return
		}
	}

	if err := middleware.WriteAck(w, nil, http.StatusOK); err != nil {
		middleware.WriteError(w, external.WriteAckError(`start transfer`,
			err), http.StatusInternalServerError)
	}
}

func (h *Handler) SuspendTransfer(w http.ResponseWriter, r *http.Request) {
	var req transfer.SuspendRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		if defaultErr.Is(err, coreErr.TypeUnmarshalError) {
			middleware.WriteError(w, external.InvalidReqBody(`suspend transfer`,
				err), http.StatusBadRequest)
			return
		}
		middleware.WriteError(w, external.ParseError(`suspend transfer`, err), http.StatusBadRequest)
		return
	}

	if req.Provider {
		if err := h.provider.SuspendTransfer(req.TransferId, req.Code, req.Reasons); err != nil {
			middleware.WriteError(w, coreErr.DSPControllerFailed(core.RoleProvider, `SuspendTransfer`, err), http.StatusBadRequest)
			return
		}
	} else {
		if err := h.consumer.SuspendTransfer(req.TransferId, req.Code, req.Reasons); err != nil {
			middleware.WriteError(w, coreErr.DSPControllerFailed(core.RoleConsumer, `SuspendTransfer`, err), http.StatusBadRequest)
			return
		}
	}

	if err := middleware.WriteAck(w, nil, http.StatusOK); err != nil {
		middleware.WriteError(w, external.WriteAckError(`suspend transfer`,
			err), http.StatusInternalServerError)
	}
}

func (h *Handler) CompleteTransfer(w http.ResponseWriter, r *http.Request) {
	var req transfer.CompleteRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		if defaultErr.Is(err, coreErr.TypeUnmarshalError) {
			middleware.WriteError(w, external.InvalidReqBody(`complete transfer`,
				err), http.StatusBadRequest)
			return
		}
		middleware.WriteError(w, external.ParseError(`complete transfer`, err), http.StatusBadRequest)
		return
	}

	if req.Provider {
		if err := h.provider.CompleteTransfer(req.TransferId); err != nil {
			middleware.WriteError(w, coreErr.DSPControllerFailed(core.RoleProvider, `CompleteTransfer`, err), http.StatusBadRequest)
			return
		}
	} else {
		if err := h.consumer.CompleteTransfer(req.TransferId); err != nil {
			middleware.WriteError(w, coreErr.DSPControllerFailed(core.RoleConsumer, `CompleteTransfer`, err), http.StatusBadRequest)
			return
		}
	}

	if err := middleware.WriteAck(w, nil, http.StatusOK); err != nil {
		middleware.WriteError(w, external.WriteAckError(`complete transfer`,
			err), http.StatusInternalServerError)
	}
}

func (h *Handler) TerminateTransfer(w http.ResponseWriter, r *http.Request) {
	var req transfer.TerminateRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		if defaultErr.Is(err, coreErr.TypeUnmarshalError) {
			middleware.WriteError(w, external.InvalidReqBody(`terminate transfer`,
				err), http.StatusBadRequest)
			return
		}
		middleware.WriteError(w, external.ParseError(`terminate transfer`, err), http.StatusBadRequest)
		return
	}

	if req.Provider {
		if err := h.provider.TerminateTransfer(req.TransferId, req.Code, req.Reasons); err != nil {
			middleware.WriteError(w, coreErr.DSPControllerFailed(core.RoleProvider, `TerminateTransfer`, err), http.StatusBadRequest)
			return
		}
	} else {
		if err := h.consumer.TerminateTransfer(req.TransferId, req.Code, req.Reasons); err != nil {
			middleware.WriteError(w, coreErr.DSPControllerFailed(core.RoleConsumer, `TerminateTransfer`, err), http.StatusBadRequest)
			return
		}
	}

	if err := middleware.WriteAck(w, nil, http.StatusOK); err != nil {
		middleware.WriteError(w, external.WriteAckError(`terminate transfer`,
			err), http.StatusInternalServerError)
	}
}
