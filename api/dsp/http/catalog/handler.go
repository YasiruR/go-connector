package catalog

import (
	defaultErr "errors"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/dsp/http/catalog"
	"github.com/YasiruR/connector/domain/core"
	coreErr "github.com/YasiruR/connector/domain/errors/core"
	"github.com/YasiruR/connector/domain/errors/dsp"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/pkg/middleware"
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
	var req catalog.Request
	if err := middleware.ParseRequest(r, &req); err != nil {
		if defaultErr.Is(err, coreErr.TypeUnmarshalError) {
			middleware.WriteError(w, dsp.CatalogInvalidReqBody(`catalog request`,
				err), http.StatusBadRequest)
			return
		}
		middleware.WriteError(w, dsp.CatalogReqParseError(`catalog request`, err), http.StatusBadRequest)
		return
	}

	cat, err := h.provider.HandleCatalogRequest(nil)
	if err != nil {
		middleware.WriteError(w, coreErr.DSPHandlerFailed(core.RoleProvider, catalog.RequestEndpoint, err),
			http.StatusInternalServerError)
		return
	}

	if err = middleware.WriteAck(w, cat, http.StatusOK); err != nil {
		middleware.WriteError(w, dsp.CatalogWriteAckError(`catalog request`, err),
			http.StatusInternalServerError)
	}
}

func (h *Handler) HandleDatasetRequest(w http.ResponseWriter, r *http.Request) {
	var req catalog.DatasetRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		if defaultErr.Is(err, coreErr.TypeUnmarshalError) {
			middleware.WriteError(w, dsp.CatalogInvalidReqBody(`dataset request`,
				err), http.StatusBadRequest)
			return
		}
		middleware.WriteError(w, dsp.CatalogReqParseError(`dataset request`, err), http.StatusBadRequest)
		return
	}

	ds, err := h.provider.HandleDatasetRequest(req.DatasetId)
	if err != nil {
		middleware.WriteError(w, coreErr.DSPHandlerFailed(core.RoleProvider, catalog.RequestDatasetEndpoint, err),
			http.StatusBadRequest)
		return
	}

	if err = middleware.WriteAck(w, ds, http.StatusOK); err != nil {
		middleware.WriteError(w, dsp.CatalogWriteAckError(`dataset request`, err),
			http.StatusInternalServerError)
	}
}
