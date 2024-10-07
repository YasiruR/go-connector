package catalog

import (
	"github.com/YasiruR/go-connector/domain"
	"github.com/YasiruR/go-connector/domain/api/gateway/http/catalog"
	"github.com/YasiruR/go-connector/domain/control-plane"
	"github.com/YasiruR/go-connector/domain/errors"
	"github.com/YasiruR/go-connector/domain/models/odrl"
	"github.com/YasiruR/go-connector/domain/pkg"
	"github.com/YasiruR/go-connector/domain/stores"
	"github.com/YasiruR/go-connector/pkg/middleware"
	"net/http"
)

type Handler struct {
	consumer    control_plane.Consumer
	owner       control_plane.Owner
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
		middleware.WriteError(w, errors.Client(errors.InvalidReqBody(`create policy`,
			err)), http.StatusBadRequest)
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
		middleware.WriteError(w, errors.DSPControllerFailed(control_plane.RoleOwner, `CreatePolicy`, err),
			http.StatusInternalServerError)
		return
	}

	if err = middleware.WriteAck(w, catalog.PolicyResponse{Id: id}, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Client(errors.WriteAckError(`create policy`,
			err)), http.StatusInternalServerError)
	}
}

func (h *Handler) CreateDataset(w http.ResponseWriter, r *http.Request) {
	var req catalog.CreateDatasetRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.Client(errors.InvalidReqBody(`create dataset`,
			err)), http.StatusBadRequest)
		return
	}

	id, err := h.owner.CreateDataset(req.Title, req.Format, req.Descriptions, req.Keywords, req.Endpoints, req.OfferIds)
	if err != nil {
		middleware.WriteError(w, errors.DSPControllerFailed(control_plane.RoleOwner, `CreateDataset`, err),
			http.StatusInternalServerError)
		return
	}

	if err = middleware.WriteAck(w, catalog.DatasetResponse{Id: id}, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Client(errors.WriteAckError(`create dataset`,
			err)), http.StatusInternalServerError)
	}
}

func (h *Handler) RequestCatalog(w http.ResponseWriter, r *http.Request) {
	var req catalog.Request
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.Client(errors.InvalidReqBody(`request catalog`,
			err)), http.StatusBadRequest)
		return
	}

	cat, err := h.consumer.RequestCatalog(req.ProviderEndpoint)
	if err != nil {
		middleware.WriteError(w, errors.DSPControllerFailed(control_plane.RoleConsumer, `RequestCatalog`, err),
			http.StatusInternalServerError)
		return
	}

	if err = middleware.WriteAck(w, cat, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Client(errors.WriteAckError(`request catalog`,
			err)), http.StatusInternalServerError)
	}
}

func (h *Handler) RequestDataset(w http.ResponseWriter, r *http.Request) {
	var req catalog.DatasetRequest
	if err := middleware.ParseRequest(r, &req); err != nil {
		middleware.WriteError(w, errors.Client(errors.InvalidReqBody(`request dataset`,
			err)), http.StatusBadRequest)
		return
	}

	ds, err := h.consumer.RequestDataset(req.DatasetId, req.ProviderEndpoint)
	if err != nil {
		middleware.WriteError(w, errors.DSPControllerFailed(control_plane.RoleConsumer, `RequestDataset`, err),
			http.StatusInternalServerError)
		return
	}

	if err = middleware.WriteAck(w, ds, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Client(errors.WriteAckError(`request dataset`,
			err)), http.StatusInternalServerError)
	}
}

func (h *Handler) GetStoredCatalogs(w http.ResponseWriter, _ *http.Request) {
	cats, err := h.consCatalog.AllCatalogs()
	if err != nil {
		middleware.WriteError(w, errors.StoreFailed(stores.TypeConsumerCatalog, `AllCatalogs`, err),
			http.StatusInternalServerError)
		return
	}

	if err = middleware.WriteAck(w, cats, http.StatusOK); err != nil {
		middleware.WriteError(w, errors.Client(errors.WriteAckError(`get catalogs`,
			err)), http.StatusInternalServerError)
	}
}
