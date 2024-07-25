package transfer

import (
	"encoding/json"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/gateway/http/transfer"
	"github.com/YasiruR/connector/domain/core"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/pkg"
	"io"
	"net/http"
)

type Handler struct {
	consumer core.Consumer
	log      pkg.Log
}

func NewHandler(roles domain.Roles, log pkg.Log) *Handler {
	return &Handler{
		consumer: roles.Consumer,
		log:      log,
	}
}

func (h *Handler) RequestTransfer(w http.ResponseWriter, r *http.Request) {
	body, err := h.readBody(transfer.RequestEndpoint, w, r)
	if err != nil {
		h.sendError(w, errors.InvalidRequestBody(transfer.RequestEndpoint, err), http.StatusBadRequest)
		return
	}

	var req transfer.Request
	if err = json.Unmarshal(body, &req); err != nil {
		h.sendError(w, errors.UnmarshalError(transfer.RequestEndpoint, err), http.StatusBadRequest)
		return
	}

	trId, err := h.consumer.RequestTransfer(req.TransferType, req.AgreementId, req.SinkEndpoint, req.ProviderEndpoint)
	if err != nil {
		h.sendError(w, errors.DSPFailed(core.RoleConsumer, `RequestTransfer`, err), http.StatusBadRequest)
		return
	}

	h.sendAck(w, transfer.RequestEndpoint, transfer.Response{TransferID: trId}, http.StatusOK)
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
		h.sendError(w, errors.MarshalError(receivedEndpoint, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	if _, err = w.Write(body); err != nil {
		h.sendError(w, errors.WriteBodyError(receivedEndpoint, err), http.StatusInternalServerError)
	}
}

func (h *Handler) sendError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	h.log.Error(errors.APIFailed(`gateway`, err))
}
