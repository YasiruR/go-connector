package transfer

import (
	"github.com/YasiruR/go-connector/domain"
	"github.com/YasiruR/go-connector/domain/api"
	"github.com/YasiruR/go-connector/domain/api/gateway/http/transfer"
	"github.com/YasiruR/go-connector/domain/control-plane"
	"github.com/YasiruR/go-connector/domain/errors"
	"github.com/YasiruR/go-connector/domain/pkg"
	"github.com/YasiruR/go-connector/pkg/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	provider control_plane.Provider
	consumer control_plane.Consumer
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
		middleware.WriteError(w, errors.Client(errors.PathParamNotFound(api.ParamConsumerPid)), http.StatusBadRequest)
		return
	}

	tp, err := h.consumer.GetProviderProcess(tpId)
	if err != nil {
		middleware.WriteError(w, errors.DSPControllerFailed(control_plane.RoleConsumer,
			`GetProviderProcess`, err), http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, tp, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Client(errors.WriteAckError(`get provider process`,
			err)), http.StatusInternalServerError)
	}
}

func (h *Handler) RequestTransfer(w http.ResponseWriter, r *http.Request) {
	var req transfer.Request
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.Client(errors.InvalidReqBody(`request transfer`,
			err)), http.StatusBadRequest)
		return
	}

	trId, err := h.consumer.RequestTransfer(req.TransferFormat, req.AgreementId,
		req.ProviderEndpoint, req.DataSink.Database)
	if err != nil {
		middleware.WriteError(w, errors.DSPControllerFailed(control_plane.RoleConsumer,
			`RequestTransfer`, err), http.StatusInternalServerError)
		return
	}

	if err = middleware.WriteAck(w, transfer.Response{TransferID: trId}, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Client(errors.WriteAckError(`request transfer`,
			err)), http.StatusInternalServerError)
	}
}

func (h *Handler) StartTransfer(w http.ResponseWriter, r *http.Request) {
	var req transfer.StartRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.Client(errors.InvalidReqBody(`start transfer`,
			err)), http.StatusBadRequest)
		return
	}

	if req.Provider {
		if err := h.provider.StartTransfer(req.TransferId, req.DataSource.Database); err != nil {
			middleware.WriteError(w, errors.DSPControllerFailed(control_plane.RoleProvider,
				`StartTransfer`, err), http.StatusInternalServerError)
			return
		}
	} else {
		if err := h.consumer.StartTransfer(req.TransferId); err != nil {
			middleware.WriteError(w, errors.DSPControllerFailed(control_plane.RoleConsumer,
				`StartTransfer`, err), http.StatusInternalServerError)
			return
		}
	}

	if err := middleware.WriteAck(w, nil, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Client(errors.WriteAckError(`start transfer`,
			err)), http.StatusInternalServerError)
	}
}

func (h *Handler) SuspendTransfer(w http.ResponseWriter, r *http.Request) {
	var req transfer.SuspendRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.Client(errors.InvalidReqBody(`suspend transfer`,
			err)), http.StatusBadRequest)
		return
	}

	if req.Provider {
		if err := h.provider.SuspendTransfer(req.TransferId, req.Code, req.Reasons); err != nil {
			middleware.WriteError(w, errors.DSPControllerFailed(control_plane.RoleProvider,
				`SuspendTransfer`, err), http.StatusInternalServerError)
			return
		}
	} else {
		if err := h.consumer.SuspendTransfer(req.TransferId, req.Code, req.Reasons); err != nil {
			middleware.WriteError(w, errors.DSPControllerFailed(control_plane.RoleConsumer,
				`SuspendTransfer`, err), http.StatusInternalServerError)
			return
		}
	}

	if err := middleware.WriteAck(w, nil, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Client(errors.WriteAckError(`suspend transfer`,
			err)), http.StatusInternalServerError)
	}
}

func (h *Handler) CompleteTransfer(w http.ResponseWriter, r *http.Request) {
	var req transfer.CompleteRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.Client(errors.InvalidReqBody(`complete transfer`,
			err)), http.StatusBadRequest)
		return
	}

	if req.Provider {
		if err := h.provider.CompleteTransfer(req.TransferId); err != nil {
			middleware.WriteError(w, errors.DSPControllerFailed(control_plane.RoleProvider,
				`CompleteTransfer`, err), http.StatusInternalServerError)
			return
		}
	} else {
		if err := h.consumer.CompleteTransfer(req.TransferId); err != nil {
			middleware.WriteError(w, errors.DSPControllerFailed(control_plane.RoleConsumer,
				`CompleteTransfer`, err), http.StatusInternalServerError)
			return
		}
	}

	if err := middleware.WriteAck(w, nil, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Client(errors.WriteAckError(`complete transfer`,
			err)), http.StatusInternalServerError)
	}
}

func (h *Handler) TerminateTransfer(w http.ResponseWriter, r *http.Request) {
	var req transfer.TerminateRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.Client(errors.InvalidReqBody(`terminate transfer`,
			err)), http.StatusBadRequest)
		return
	}

	if req.Provider {
		if err := h.provider.TerminateTransfer(req.TransferId, req.Code, req.Reasons); err != nil {
			middleware.WriteError(w, errors.DSPControllerFailed(control_plane.RoleProvider,
				`TerminateTransfer`, err), http.StatusInternalServerError)
			return
		}
	} else {
		if err := h.consumer.TerminateTransfer(req.TransferId, req.Code, req.Reasons); err != nil {
			middleware.WriteError(w, errors.DSPControllerFailed(control_plane.RoleConsumer,
				`TerminateTransfer`, err), http.StatusInternalServerError)
			return
		}
	}

	if err := middleware.WriteAck(w, nil, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Client(errors.WriteAckError(`terminate transfer`,
			err)), http.StatusInternalServerError)
	}
}
