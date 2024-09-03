package catalog

import (
	defaultErr "errors"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/gateway/http/catalog"
	"github.com/YasiruR/connector/domain/core"
	coreErr "github.com/YasiruR/connector/domain/errors/core"
	"github.com/YasiruR/connector/domain/errors/external"
	"github.com/YasiruR/connector/domain/models/odrl"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/stores"
	"github.com/YasiruR/connector/pkg/middleware"
	"net/http"
)

type Handler struct {
	consumer    core.Consumer
	owner       core.Owner
	consCatalog stores.ConsumerCatalog
	log         pkg.Log
}

func NewHandler(roles domain.Roles, stores domain.Stores, log pkg.Log) *Handler {
	return &Handler{
		consumer:    roles.Consumer,
		owner:       roles.Owner,
		consCatalog: stores.ConsumerCatalog,
		log:         log,
	}
}

func (h *Handler) CreatePolicy(w http.ResponseWriter, r *http.Request) {
	var req catalog.CreatePolicyRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		if defaultErr.Is(err, coreErr.TypeUnmarshalError) {
			middleware.WriteError(w, external.InvalidReqBody(`create policy`,
				err), http.StatusBadRequest)
			return
		}
		middleware.WriteError(w, external.ParseError(`create policy`, err), http.StatusBadRequest)
		return
	}

	var perms []odrl.Rule // handle other policy types
	for _, p := range req.Permissions {
		var cons []odrl.Constraint
		for _, c := range p.Constraints {
			cons = append(cons, odrl.Constraint{
				LeftOperand:  c.LeftOperand,
				Operator:     c.Operator,
				RightOperand: c.RightOperand,
			})
		}
		perms = append(perms, odrl.Rule{Action: odrl.Action(p.Action), Constraints: cons})
	}

	id, err := h.owner.CreatePolicy(req.Target, perms, []odrl.Rule{})
	if err != nil {
		middleware.WriteError(w, coreErr.DSPControllerFailed(core.RoleOwner, `CreatePolicy`, err),
			http.StatusInternalServerError)
		return
	}

	if err = middleware.WriteAck(w, catalog.PolicyResponse{Id: id}, http.StatusOK); err != nil {
		middleware.WriteError(w, external.WriteAckError(`create policy`,
			err), http.StatusInternalServerError)
	}
}

func (h *Handler) CreateDataset(w http.ResponseWriter, r *http.Request) {
	var req catalog.CreateDatasetRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		if defaultErr.Is(err, coreErr.TypeUnmarshalError) {
			middleware.WriteError(w, external.InvalidReqBody(`create dataset`,
				err), http.StatusBadRequest)
			return
		}
		middleware.WriteError(w, external.ParseError(`create dataset`, err), http.StatusBadRequest)
		return
	}

	id, err := h.owner.CreateDataset(req.Title, req.Format, req.Descriptions, req.Keywords, req.Endpoints, req.OfferIds)
	if err != nil {
		middleware.WriteError(w, coreErr.DSPControllerFailed(core.RoleOwner, `CreateDataset`, err),
			http.StatusInternalServerError)
		return
	}

	if err = middleware.WriteAck(w, catalog.DatasetResponse{Id: id}, http.StatusOK); err != nil {
		middleware.WriteError(w, external.WriteAckError(`create dataset`,
			err), http.StatusInternalServerError)
	}
}

func (h *Handler) RequestCatalog(w http.ResponseWriter, r *http.Request) {
	var req catalog.Request
	if err := middleware.ParseRequest(r, &req); err != nil {
		if defaultErr.Is(err, coreErr.TypeUnmarshalError) {
			middleware.WriteError(w, external.InvalidReqBody(`request catalog`,
				err), http.StatusBadRequest)
			return
		}
		middleware.WriteError(w, external.ParseError(`request catalog`, err), http.StatusBadRequest)
		return
	}

	cat, err := h.consumer.RequestCatalog(req.ProviderEndpoint)
	if err != nil {
		middleware.WriteError(w, coreErr.DSPControllerFailed(core.RoleConsumer, `RequestCatalog`, err),
			http.StatusInternalServerError)
		return
	}

	if err = middleware.WriteAck(w, cat, http.StatusOK); err != nil {
		middleware.WriteError(w, external.WriteAckError(`request catalog`,
			err), http.StatusInternalServerError)
	}
}

func (h *Handler) RequestDataset(w http.ResponseWriter, r *http.Request) {
	var req catalog.DatasetRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		if defaultErr.Is(err, coreErr.TypeUnmarshalError) {
			middleware.WriteError(w, external.InvalidReqBody(`request dataset`,
				err), http.StatusBadRequest)
			return
		}
		middleware.WriteError(w, external.ParseError(`request dataset`, err), http.StatusBadRequest)
		return
	}

	ds, err := h.consumer.RequestDataset(req.DatasetId, req.ProviderEndpoint)
	if err != nil {
		middleware.WriteError(w, coreErr.DSPControllerFailed(core.RoleConsumer, `RequestDataset`, err),
			http.StatusInternalServerError)
		return
	}

	if err = middleware.WriteAck(w, ds, http.StatusOK); err != nil {
		middleware.WriteError(w, external.WriteAckError(`request dataset`,
			err), http.StatusInternalServerError)
	}
}

func (h *Handler) GetStoredCatalogs(w http.ResponseWriter, _ *http.Request) {
	cats, err := h.consCatalog.AllCatalogs()
	if err != nil {
		middleware.WriteError(w, coreErr.StoreFailed(stores.TypeConsumerCatalog, `AllCatalogs`, err),
			http.StatusInternalServerError)
		return
	}

	if err = middleware.WriteAck(w, cats, http.StatusOK); err != nil {
		middleware.WriteError(w, external.WriteAckError(`get catalogs`,
			err), http.StatusInternalServerError)
	}
}
