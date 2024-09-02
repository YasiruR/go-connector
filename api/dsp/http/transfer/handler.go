package transfer

import (
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api"
	"github.com/YasiruR/connector/domain/api/dsp/http/transfer"
	"github.com/YasiruR/connector/domain/core"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/ror"
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

func (h *Handler) HandleGetProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tpId, ok := vars[api.ParamProviderPid]
	if !ok {
		middleware.WriteError(w, errors.PathParamNotFound(transfer.GetProcessEndpoint,
			api.ParamProviderPid), http.StatusBadRequest)
		return
	}

	ack, err := h.provider.HandleGetProcess(tpId)
	if err != nil {
		middleware.WriteError(w, ror.DSPHandlerFailed(core.RoleProvider,
			transfer.GetProcessEndpoint, err), http.StatusBadRequest)
		return
	}

	middleware.WriteAck(w, ack, http.StatusOK)
}

func (h *Handler) HandleTransferRequest(w http.ResponseWriter, r *http.Request) {
	var req transfer.Request
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.ParseRequestFailed(transfer.RequestEndpoint,
			err), http.StatusBadRequest)
		return
	}

	ack, err := h.provider.HandleTransferRequest(req)
	if err != nil {
		middleware.WriteError(w, ror.DSPHandlerFailed(core.RoleProvider,
			transfer.RequestEndpoint, err), http.StatusBadRequest)
		return
	}

	middleware.WriteAck(w, ack, http.StatusCreated)
}

func (h *Handler) HandleTransferStart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, ok := vars[api.ParamPid]
	if !ok {
		middleware.WriteError(w, errors.PathParamNotFound(transfer.StartEndpoint,
			api.ParamPid), http.StatusBadRequest)
		return
	}

	var req transfer.StartRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.ParseRequestFailed(transfer.StartEndpoint, err), http.StatusBadRequest)
		return
	}

	var ack transfer.Ack
	var err error
	switch pid {
	case req.ProvPId:
		ack, err = h.provider.HandleTransferStart(req)
		if err != nil {
			middleware.WriteError(w, ror.DSPHandlerFailed(core.RoleProvider,
				transfer.StartEndpoint, err), http.StatusBadRequest)
			return
		}
	case req.ConsPId:
		ack, err = h.consumer.HandleTransferStart(req)
		if err != nil {
			middleware.WriteError(w, ror.DSPHandlerFailed(core.RoleConsumer,
				transfer.StartEndpoint, err), http.StatusBadRequest)
			return
		}
	default:
		middleware.WriteError(w, errors.InvalidRequestBody(transfer.SuspendEndpoint, err), http.StatusBadRequest)
		return
	}

	middleware.WriteAck(w, ack, http.StatusOK)
}

func (h *Handler) HandleTransferSuspension(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, ok := vars[api.ParamPid]
	if !ok {
		middleware.WriteError(w, errors.PathParamNotFound(transfer.SuspendEndpoint,
			api.ParamPid), http.StatusBadRequest)
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
			middleware.WriteError(w, ror.DSPHandlerFailed(core.RoleProvider,
				transfer.SuspendEndpoint, err), http.StatusBadRequest)
			return
		}
	case req.ConsPId:
		ack, err = h.consumer.HandleTransferSuspension(req)
		if err != nil {
			middleware.WriteError(w, ror.DSPHandlerFailed(core.RoleConsumer,
				transfer.SuspendEndpoint, err), http.StatusBadRequest)
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
		middleware.WriteError(w, errors.PathParamNotFound(transfer.CompleteEndpoint,
			api.ParamPid), http.StatusBadRequest)
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
			middleware.WriteError(w, ror.DSPHandlerFailed(core.RoleProvider,
				transfer.CompleteEndpoint, err), http.StatusBadRequest)
			return
		}
	case req.ConsPId:
		ack, err = h.consumer.HandleTransferCompletion(req)
		if err != nil {
			middleware.WriteError(w, ror.DSPHandlerFailed(core.RoleConsumer,
				transfer.CompleteEndpoint, err), http.StatusBadRequest)
			return
		}
	default:
		middleware.WriteError(w, errors.InvalidRequestBody(transfer.CompleteEndpoint, err), http.StatusBadRequest)
		return
	}

	middleware.WriteAck(w, ack, http.StatusOK)
}

func (h *Handler) HandleTransferTermination(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, ok := vars[api.ParamPid]
	if !ok {
		middleware.WriteError(w, errors.PathParamNotFound(transfer.TerminateEndpoint,
			api.ParamPid), http.StatusBadRequest)
		return
	}

	var req transfer.TerminateRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.ParseRequestFailed(transfer.TerminateEndpoint, err), http.StatusBadRequest)
		return
	}

	var ack transfer.Ack
	var err error
	switch pid {
	case req.ProvPId:
		ack, err = h.provider.HandleTransferTermination(req)
		if err != nil {
			middleware.WriteError(w, ror.DSPHandlerFailed(core.RoleProvider,
				transfer.TerminateEndpoint, err), http.StatusBadRequest)
			return
		}
	case req.ConsPId:
		ack, err = h.consumer.HandleTransferTermination(req)
		if err != nil {
			middleware.WriteError(w, ror.DSPHandlerFailed(core.RoleConsumer,
				transfer.TerminateEndpoint, err), http.StatusBadRequest)
			return
		}
	default:
		middleware.WriteError(w, errors.InvalidRequestBody(transfer.TerminateEndpoint, err), http.StatusBadRequest)
		return
	}

	middleware.WriteAck(w, ack, http.StatusOK)
}
