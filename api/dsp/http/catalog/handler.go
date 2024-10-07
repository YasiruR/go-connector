package catalog

import (
	"github.com/YasiruR/go-connector/domain"
	"github.com/YasiruR/go-connector/domain/api/dsp/http/catalog"
	"github.com/YasiruR/go-connector/domain/control-plane"
	"github.com/YasiruR/go-connector/domain/errors"
	"github.com/YasiruR/go-connector/domain/pkg"
	"github.com/YasiruR/go-connector/pkg/middleware"
	"net/http"
)

type Handler struct {
	provider control_plane.Provider
	log      pkg.Log
}

func NewHandler(roles domain.Roles, log pkg.Log) *Handler {
	return &Handler{
		provider: roles.Provider,
		log:      log,
	}
}

func (h *Handler) HandleCatalogRequest(w http.ResponseWriter, r *http.Request) {
	var req catalog.Request
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.Catalog(errors.InvalidReqBody(`catalog request`, err)), http.StatusBadRequest)
		return
	}

	cat, err := h.provider.HandleCatalogRequest(nil)
	if err != nil {
		middleware.WriteError(w, errors.DSPHandlerFailed(control_plane.RoleProvider, catalog.RequestEndpoint, err),
			http.StatusInternalServerError)
		return
	}

	if err = middleware.WriteAck(w, cat, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Catalog(errors.WriteAckError(`catalog request`, err)),
			http.StatusInternalServerError)
	}
}

func (h *Handler) HandleDatasetRequest(w http.ResponseWriter, r *http.Request) {
	var req catalog.DatasetRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.Catalog(errors.InvalidReqBody(`dataset request`, err)), http.StatusBadRequest)
		return
	}

	ds, err := h.provider.HandleDatasetRequest(req.DatasetId)
	if err != nil {
		middleware.WriteError(w, errors.DSPHandlerFailed(control_plane.RoleProvider, catalog.RequestDatasetEndpoint, err),
			http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, ds, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Catalog(errors.WriteAckError(`dataset request`, err)),
			http.StatusInternalServerError)
	}
}
