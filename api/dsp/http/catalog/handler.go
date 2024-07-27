package catalog

import (
	"encoding/json"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/dsp/http/catalog"
	"github.com/YasiruR/connector/domain/core"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/pkg"
	"io"
	"net/http"
)

type Handler struct {
	provider core.Provider
	log      pkg.Log
}

func NewHandler(roles domain.Roles, log pkg.Log) *Handler {
	return &Handler{
		provider: roles.Provider,
		log:      log,
	}
}

func (h *Handler) HandleCatalogRequest(w http.ResponseWriter, r *http.Request) {
	body, err := h.readBody(catalog.RequestEndpoint, w, r)
	if err != nil {
		return
	}

	var req catalog.Request
	if err = json.Unmarshal(body, &req); err != nil {
		h.sendError(w, errors.UnmarshalError(catalog.RequestEndpoint, err), http.StatusBadRequest)
		return
	}

	cat, err := h.provider.HandleCatalogRequest(nil)
	if err != nil {
		h.sendError(w, errors.DSPHandlerFailed(catalog.RequestEndpoint, core.RoleProvider, err), http.StatusBadRequest)
		return
	}

	h.sendAck(w, catalog.RequestEndpoint, cat, http.StatusOK)
}

func (h *Handler) HandleDatasetRequest(w http.ResponseWriter, r *http.Request) {
	body, err := h.readBody(catalog.RequestDatasetEndpoint, w, r)
	if err != nil {
		return
	}

	var req catalog.DatasetRequest
	if err = json.Unmarshal(body, &req); err != nil {
		h.sendError(w, errors.UnmarshalError(catalog.TypeDatasetRequest, err), http.StatusBadRequest)
		return
	}

	ds, err := h.provider.HandleDatasetRequest(req.DatasetId)
	if err != nil {
		h.sendError(w, errors.DSPHandlerFailed(catalog.RequestDatasetEndpoint, core.RoleProvider, err), http.StatusBadRequest)
		return
	}

	h.sendAck(w, catalog.RequestDatasetEndpoint, ds, http.StatusOK)
}

func (h *Handler) readBody(endpoint string, w http.ResponseWriter, r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		err = errors.InvalidRequestBody(endpoint, err)
		w.WriteHeader(http.StatusBadRequest)
		h.log.Error(err)
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
	h.log.Error(errors.APIFailed(`dsp`, err))
}
