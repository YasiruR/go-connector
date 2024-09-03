package transfer

import (
	defaultErr "errors"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api"
	"github.com/YasiruR/connector/domain/api/dsp/http/transfer"
	"github.com/YasiruR/connector/domain/core"
	coreErr "github.com/YasiruR/connector/domain/errors/core"
	"github.com/YasiruR/connector/domain/errors/dsp"
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

func (h *Handler) HandleGetProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tpId, ok := vars[api.ParamProviderPid]
	if !ok {
		middleware.WriteError(w, dsp.TransferPathParamError(api.ParamProviderPid), http.StatusBadRequest)
		return
	}

	ack, err := h.provider.HandleGetProcess(tpId)
	if err != nil {
		middleware.WriteError(w, coreErr.DSPHandlerFailed(core.RoleProvider,
			transfer.GetProcessEndpoint, err), http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, ack, http.StatusOK); err != nil {
		middleware.WriteError(w, dsp.TransferWriteAckError(ack.ProvPId, ack.ConsPId,
			`get process`, err), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleTransferRequest(w http.ResponseWriter, r *http.Request) {
	var req transfer.Request
	if err := middleware.ParseRequest(r, &req); err != nil {
		if defaultErr.Is(err, coreErr.TypeUnmarshalError) {
			middleware.WriteError(w, dsp.TransferInvalidReqBody(``, ``,
				`transfer request`, err), http.StatusBadRequest)
			return
		}
		middleware.WriteError(w, dsp.TransferReqParseError(`transfer request`, err), http.StatusBadRequest)
		return
	}

	ack, err := h.provider.HandleTransferRequest(req)
	if err != nil {
		middleware.WriteError(w, coreErr.DSPHandlerFailed(core.RoleProvider,
			transfer.RequestEndpoint, err), http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, ack, http.StatusCreated); err != nil {
		middleware.WriteError(w, dsp.TransferWriteAckError(ack.ProvPId, ack.ConsPId,
			`transfer request`, err), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleTransferStart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, ok := vars[api.ParamPid]
	if !ok {
		middleware.WriteError(w, dsp.TransferPathParamError(api.ParamPid), http.StatusBadRequest)
		return
	}

	var req transfer.StartRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		if defaultErr.Is(err, coreErr.TypeUnmarshalError) {
			middleware.WriteError(w, dsp.TransferInvalidReqBody(``, ``,
				`transfer start`, err), http.StatusBadRequest)
			return
		}
		middleware.WriteError(w, dsp.TransferReqParseError(`transfer start`, err), http.StatusBadRequest)
		return
	}

	var ack transfer.Ack
	var err error
	switch pid {
	case req.ProvPId:
		ack, err = h.provider.HandleTransferStart(req)
		if err != nil {
			middleware.WriteError(w, coreErr.DSPHandlerFailed(core.RoleProvider,
				transfer.StartEndpoint, err), http.StatusBadRequest)
			return
		}
	case req.ConsPId:
		ack, err = h.consumer.HandleTransferStart(req)
		if err != nil {
			middleware.WriteError(w, coreErr.DSPHandlerFailed(core.RoleConsumer,
				transfer.StartEndpoint, err), http.StatusBadRequest)
			return
		}
	default:
		middleware.WriteError(w, external.IncompatibleReqBody(
			`path parameter should be provider/consumer process ID`), http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, ack, http.StatusOK); err != nil {
		middleware.WriteError(w, dsp.TransferWriteAckError(ack.ProvPId, ack.ConsPId,
			`transfer start`, err), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleTransferSuspension(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, ok := vars[api.ParamPid]
	if !ok {
		middleware.WriteError(w, dsp.TransferPathParamError(api.ParamPid), http.StatusBadRequest)
		return
	}

	var req transfer.SuspendRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		if defaultErr.Is(err, coreErr.TypeUnmarshalError) {
			middleware.WriteError(w, dsp.TransferInvalidReqBody(``, ``,
				`transfer suspension`, err), http.StatusBadRequest)
			return
		}
		middleware.WriteError(w, dsp.TransferReqParseError(`transfer suspend`, err), http.StatusBadRequest)
		return
	}

	var ack transfer.Ack
	var err error
	switch pid {
	case req.ProvPId:
		ack, err = h.provider.HandleTransferSuspension(req)
		if err != nil {
			middleware.WriteError(w, coreErr.DSPHandlerFailed(core.RoleProvider,
				transfer.SuspendEndpoint, err), http.StatusBadRequest)
			return
		}
	case req.ConsPId:
		ack, err = h.consumer.HandleTransferSuspension(req)
		if err != nil {
			middleware.WriteError(w, coreErr.DSPHandlerFailed(core.RoleConsumer,
				transfer.SuspendEndpoint, err), http.StatusBadRequest)
			return
		}
	default:
		middleware.WriteError(w, external.IncompatibleReqBody(
			`path parameter should be provider/consumer process ID`), http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, ack, http.StatusOK); err != nil {
		middleware.WriteError(w, dsp.TransferWriteAckError(ack.ProvPId, ack.ConsPId,
			`transfer suspension`, err), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleTransferCompletion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, ok := vars[api.ParamPid]
	if !ok {
		middleware.WriteError(w, dsp.TransferPathParamError(api.ParamPid), http.StatusBadRequest)
		return
	}

	var req transfer.CompleteRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		if defaultErr.Is(err, coreErr.TypeUnmarshalError) {
			middleware.WriteError(w, dsp.TransferInvalidReqBody(``, ``,
				`transfer completion`, err), http.StatusBadRequest)
			return
		}
		middleware.WriteError(w, dsp.TransferReqParseError(`transfer completion`, err), http.StatusBadRequest)
		return
	}

	var ack transfer.Ack
	var err error
	switch pid {
	case req.ProvPId:
		ack, err = h.provider.HandleTransferCompletion(req)
		if err != nil {
			middleware.WriteError(w, coreErr.DSPHandlerFailed(core.RoleProvider,
				transfer.CompleteEndpoint, err), http.StatusBadRequest)
			return
		}
	case req.ConsPId:
		ack, err = h.consumer.HandleTransferCompletion(req)
		if err != nil {
			middleware.WriteError(w, coreErr.DSPHandlerFailed(core.RoleConsumer,
				transfer.CompleteEndpoint, err), http.StatusBadRequest)
			return
		}
	default:
		middleware.WriteError(w, external.IncompatibleReqBody(
			`path parameter should be provider/consumer process ID`), http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, ack, http.StatusOK); err != nil {
		middleware.WriteError(w, dsp.TransferWriteAckError(ack.ProvPId, ack.ConsPId,
			`transfer completion`, err), http.StatusInternalServerError)
	}
}

func (h *Handler) HandleTransferTermination(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, ok := vars[api.ParamPid]
	if !ok {
		middleware.WriteError(w, dsp.TransferPathParamError(api.ParamPid), http.StatusBadRequest)
		return
	}

	var req transfer.TerminateRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		if defaultErr.Is(err, coreErr.TypeUnmarshalError) {
			middleware.WriteError(w, dsp.TransferInvalidReqBody(``, ``,
				`transfer termination`, err), http.StatusBadRequest)
			return
		}
		middleware.WriteError(w, dsp.TransferReqParseError(`transfer termination`,
			err), http.StatusBadRequest)
		return
	}

	var ack transfer.Ack
	var err error
	switch pid {
	case req.ProvPId:
		ack, err = h.provider.HandleTransferTermination(req)
		if err != nil {
			middleware.WriteError(w, coreErr.DSPHandlerFailed(core.RoleProvider,
				transfer.TerminateEndpoint, err), http.StatusBadRequest)
			return
		}
	case req.ConsPId:
		ack, err = h.consumer.HandleTransferTermination(req)
		if err != nil {
			middleware.WriteError(w, coreErr.DSPHandlerFailed(core.RoleConsumer,
				transfer.TerminateEndpoint, err), http.StatusBadRequest)
			return
		}
	default:
		middleware.WriteError(w, external.IncompatibleReqBody(
			`path parameter should be provider/consumer process ID`), http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, ack, http.StatusOK); err != nil {
		middleware.WriteError(w, dsp.TransferWriteAckError(ack.ProvPId, ack.ConsPId,
			`transfer termination`, err), http.StatusInternalServerError)
	}
}
