package transfer

import (
	"github.com/YasiruR/connector/api/dsp/http/middleware"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api"
	"github.com/YasiruR/connector/domain/api/dsp/http/transfer"
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

func (h *Handler) HandleTransfers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tpId, ok := vars[api.ParamPid]
	if !ok {
		middleware.WriteError(w, errors.PathParamNotFound(transfer.TransfersEndpoint, api.ParamPid), http.StatusBadRequest)
		return
	}

	ack, err := h.provider.HandleTransfers(tpId)
	if err != nil {
		middleware.WriteError(w, errors.DSPHandlerFailed(core.RoleProvider, transfer.TransfersEndpoint, err), http.StatusBadRequest)
		return
	}

	middleware.WriteAck(w, ack, http.StatusOK)
}

func (h *Handler) HandleTransferRequest(w http.ResponseWriter, r *http.Request) {
	var req transfer.Request
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.ParseRequestFailed(transfer.RequestEndpoint, err), http.StatusBadRequest)
		return
	}

	ack, err := h.provider.HandleTransferRequest(req)
	if err != nil {
		middleware.WriteError(w, errors.DSPHandlerFailed(core.RoleProvider, transfer.RequestEndpoint, err), http.StatusBadRequest)
		return
	}

	middleware.WriteAck(w, ack, http.StatusCreated)
}

func (h *Handler) HandleTransferStart(w http.ResponseWriter, r *http.Request) {
	var req transfer.StartRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.ParseRequestFailed(transfer.StartEndpoint, err), http.StatusBadRequest)
		return
	}

	ack, err := h.consumer.HandleTransferStart(req)
	if err != nil {
		middleware.WriteError(w, errors.DSPHandlerFailed(core.RoleProvider, transfer.StartEndpoint, err), http.StatusBadRequest)
		return
	}

	middleware.WriteAck(w, ack, http.StatusOK)
}

func (h *Handler) HandleTransferSuspension(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, ok := vars[api.ParamPid]
	if !ok {
		middleware.WriteError(w, errors.PathParamNotFound(transfer.SuspendEndpoint, api.ParamPid), http.StatusBadRequest)
		return
	}

	var req transfer.SuspendRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.ParseRequestFailed(transfer.SuspendEndpoint, err), http.StatusBadRequest)
		return
	}

	var ack transfer.Ack
	var err error
	switch pid {
	case req.ProvPId:
		ack, err = h.provider.HandleTransferSuspension(req)
		if err != nil {
			middleware.WriteError(w, errors.DSPHandlerFailed(core.RoleProvider, transfer.SuspendEndpoint, err), http.StatusBadRequest)
			return
		}
	case req.ConsPId:
		ack, err = h.consumer.HandleTransferSuspension(req)
		if err != nil {
			middleware.WriteError(w, errors.DSPHandlerFailed(core.RoleConsumer, transfer.SuspendEndpoint, err), http.StatusBadRequest)
			return
		}
	default:
		middleware.WriteError(w, errors.InvalidRequestBody(transfer.SuspendEndpoint, err), http.StatusBadRequest)
		return
	}

	middleware.WriteAck(w, ack, http.StatusOK)
}

func (h *Handler) HandleTransferCompletion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, ok := vars[api.ParamPid]
	if !ok {
		middleware.WriteError(w, errors.PathParamNotFound(transfer.CompleteEndpoint, api.ParamPid), http.StatusBadRequest)
		return
	}

	var req transfer.CompleteRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.ParseRequestFailed(transfer.CompleteEndpoint, err), http.StatusBadRequest)
		return
	}

	var ack transfer.Ack
	var err error
	switch pid {
	case req.ProvPId:
		ack, err = h.provider.HandleTransferCompletion(req)
		if err != nil {
			middleware.WriteError(w, errors.DSPHandlerFailed(core.RoleProvider, transfer.CompleteEndpoint, err), http.StatusBadRequest)
			return
		}
	case req.ConsPId:
		ack, err = h.consumer.HandleTransferCompletion(req)
		if err != nil {
			middleware.WriteError(w, errors.DSPHandlerFailed(core.RoleConsumer, transfer.CompleteEndpoint, err), http.StatusBadRequest)
			return
		}
	default:
		middleware.WriteError(w, errors.InvalidRequestBody(transfer.CompleteEndpoint, err), http.StatusBadRequest)
		return
	}

	middleware.WriteAck(w, ack, http.StatusOK)
}
