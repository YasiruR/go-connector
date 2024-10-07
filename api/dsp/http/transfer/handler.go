package transfer

import (
	"github.com/YasiruR/go-connector/domain"
	"github.com/YasiruR/go-connector/domain/api"
	"github.com/YasiruR/go-connector/domain/api/dsp/http/transfer"
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

func (h *Handler) HandleGetProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tpId, ok := vars[api.ParamProviderPid]
	if !ok {
		middleware.WriteError(w, errors.Transfer(``, ``,
			errors.PathParamNotFound(api.ParamProviderPid)), http.StatusBadRequest)
		return
	}

	ack, err := h.provider.HandleGetProcess(tpId)
	if err != nil {
		middleware.WriteError(w, errors.DSPHandlerFailed(control_plane.RoleProvider,
			transfer.GetProcessEndpoint, err), http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, ack, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Transfer(ack.ProvPId, ack.ConsPId,
			errors.WriteAckError(`get process`, err)), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleTransferRequest(w http.ResponseWriter, r *http.Request) {
	var req transfer.Request
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.Transfer(``, ``,
			errors.InvalidReqBody(`transfer request`, err)), http.StatusBadRequest)
		return
	}

	ack, err := h.provider.HandleTransferRequest(req)
	if err != nil {
		middleware.WriteError(w, errors.DSPHandlerFailed(control_plane.RoleProvider,
			transfer.RequestEndpoint, err), http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, ack, http.StatusCreated); err != nil {
		middleware.WriteError(w, errors.Transfer(ack.ProvPId, ack.ConsPId,
			errors.WriteAckError(`transfer request`, err)), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleTransferStart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, ok := vars[api.ParamPid]
	if !ok {
		middleware.WriteError(w, errors.Transfer(``, ``,
			errors.PathParamNotFound(api.ParamPid)), http.StatusBadRequest)
		return
	}

	var req transfer.StartRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.Transfer(``, ``,
			errors.InvalidReqBody(`transfer start`, err)), http.StatusBadRequest)
		return
	}

	var ack transfer.Ack
	var err error
	switch pid {
	case req.ProvPId:
		ack, err = h.provider.HandleTransferStart(req)
		if err != nil {
			middleware.WriteError(w, errors.DSPHandlerFailed(control_plane.RoleProvider,
				transfer.StartEndpoint, err), http.StatusBadRequest)
			return
		}
	case req.ConsPId:
		ack, err = h.consumer.HandleTransferStart(req)
		if err != nil {
			middleware.WriteError(w, errors.DSPHandlerFailed(control_plane.RoleConsumer,
				transfer.StartEndpoint, err), http.StatusBadRequest)
			return
		}
	default:
		middleware.WriteError(w, errors.Transfer(``, ``, errors.IncorrectReqValues(
			`path parameter should be provider/consumer process ID`)), http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, ack, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Transfer(ack.ProvPId, ack.ConsPId,
			errors.WriteAckError(`transfer start`, err)), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleTransferSuspension(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, ok := vars[api.ParamPid]
	if !ok {
		middleware.WriteError(w, errors.Transfer(``, ``,
			errors.PathParamNotFound(api.ParamPid)), http.StatusBadRequest)
		return
	}

	var req transfer.SuspendRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.Transfer(``, ``,
			errors.InvalidReqBody(`transfer suspend`, err)), http.StatusBadRequest)
		return
	}

	var ack transfer.Ack
	var err error
	switch pid {
	case req.ProvPId:
		ack, err = h.provider.HandleTransferSuspension(req)
		if err != nil {
			middleware.WriteError(w, errors.DSPHandlerFailed(control_plane.RoleProvider,
				transfer.SuspendEndpoint, err), http.StatusBadRequest)
			return
		}
	case req.ConsPId:
		ack, err = h.consumer.HandleTransferSuspension(req)
		if err != nil {
			middleware.WriteError(w, errors.DSPHandlerFailed(control_plane.RoleConsumer,
				transfer.SuspendEndpoint, err), http.StatusBadRequest)
			return
		}
	default:
		middleware.WriteError(w, errors.Transfer(``, ``, errors.IncorrectReqValues(
			`path parameter should be provider/consumer process ID`)), http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, ack, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Transfer(ack.ProvPId, ack.ConsPId,
			errors.WriteAckError(`transfer suspension`, err)), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleTransferCompletion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, ok := vars[api.ParamPid]
	if !ok {
		middleware.WriteError(w, errors.Transfer(``, ``,
			errors.PathParamNotFound(api.ParamPid)), http.StatusBadRequest)
		return
	}

	var req transfer.CompleteRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.Transfer(``, ``,
			errors.InvalidReqBody(`transfer completion`, err)), http.StatusBadRequest)
		return
	}

	var ack transfer.Ack
	var err error
	switch pid {
	case req.ProvPId:
		ack, err = h.provider.HandleTransferCompletion(req)
		if err != nil {
			middleware.WriteError(w, errors.DSPHandlerFailed(control_plane.RoleProvider,
				transfer.CompleteEndpoint, err), http.StatusBadRequest)
			return
		}
	case req.ConsPId:
		ack, err = h.consumer.HandleTransferCompletion(req)
		if err != nil {
			middleware.WriteError(w, errors.DSPHandlerFailed(control_plane.RoleConsumer,
				transfer.CompleteEndpoint, err), http.StatusBadRequest)
			return
		}
	default:
		middleware.WriteError(w, errors.Transfer(``, ``, errors.IncorrectReqValues(
			`path parameter should be provider/consumer process ID`)), http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, ack, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Transfer(ack.ProvPId, ack.ConsPId,
			errors.WriteAckError(`transfer completion`, err)), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleTransferTermination(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, ok := vars[api.ParamPid]
	if !ok {
		middleware.WriteError(w, errors.Transfer(``, ``,
			errors.PathParamNotFound(api.ParamPid)), http.StatusBadRequest)
		return
	}

	var req transfer.TerminateRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.Transfer(``, ``,
			errors.InvalidReqBody(`transfer termination`, err)), http.StatusBadRequest)
		return
	}

	var ack transfer.Ack
	var err error
	switch pid {
	case req.ProvPId:
		ack, err = h.provider.HandleTransferTermination(req)
		if err != nil {
			middleware.WriteError(w, errors.DSPHandlerFailed(control_plane.RoleProvider,
				transfer.TerminateEndpoint, err), http.StatusBadRequest)
			return
		}
	case req.ConsPId:
		ack, err = h.consumer.HandleTransferTermination(req)
		if err != nil {
			middleware.WriteError(w, errors.DSPHandlerFailed(control_plane.RoleConsumer,
				transfer.TerminateEndpoint, err), http.StatusBadRequest)
			return
		}
	default:
		middleware.WriteError(w, errors.Transfer(``, ``, errors.IncorrectReqValues(
			`path parameter should be provider/consumer process ID`)), http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, ack, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Transfer(ack.ProvPId, ack.ConsPId,
			errors.WriteAckError(`transfer termination`, err)), http.StatusInternalServerError)
	}
}
