package catalog

import (
	"encoding/json"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/gateway/http/catalog"
	"github.com/YasiruR/connector/domain/core"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/models/odrl"
	"github.com/YasiruR/connector/domain/pkg"
	"io"
	"net/http"
)

type Handler struct {
	consumer core.Consumer
	owner    core.Owner
	log      pkg.Log
}

func NewHandler(roles domain.Roles, log pkg.Log) *Handler {
	return &Handler{
		consumer: roles.Consumer,
		owner:    roles.Owner,
		log:      log,
	}
}

func (h *Handler) CreatePolicy(w http.ResponseWriter, r *http.Request) {
	body, err := h.readBody(catalog.CreatePolicyEndpoint, w, r)
	if err != nil {
		return
	}

	var req catalog.CreatePolicyRequest
	if err = json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Error(errors.UnmarshalError(catalog.CreatePolicyEndpoint, err))
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

	// todo check if target is required here
	id, err := h.owner.CreatePolicy(`test`, perms, []odrl.Rule{})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Error(errors.DSPControllerFailed(core.RoleOwner, `CreatePolicy`, err))
		return
	}

	h.sendAck(w, catalog.CreatePolicyEndpoint, catalog.PolicyResponse{Id: id}, http.StatusOK)
}

func (h *Handler) CreateDataset(w http.ResponseWriter, r *http.Request) {
	body, err := h.readBody(catalog.CreateDatasetEndpoint, w, r)
	if err != nil {
		return
	}

	var req catalog.CreateDatasetRequest
	if err = json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Error(errors.UnmarshalError(catalog.CreateDatasetEndpoint, err))
		return
	}

	id, err := h.owner.CreateDataset(req.Title, req.Format, req.Descriptions, req.Keywords, req.Endpoints, req.OfferIds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Error(errors.DSPControllerFailed(core.RoleOwner, `CreateDataset`, err))
		return
	}

	h.sendAck(w, catalog.CreateDatasetEndpoint, catalog.DatasetResponse{Id: id}, http.StatusOK)
}

func (h *Handler) RequestCatalog(w http.ResponseWriter, r *http.Request) {
	body, err := h.readBody(catalog.RequestCatalogEndpoint, w, r)
	if err != nil {
		return
	}

	var req catalog.Request
	if err = json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Error(errors.UnmarshalError(catalog.RequestCatalogEndpoint, err))
		return
	}

	cat, err := h.consumer.RequestCatalog(req.ProviderEndpoint)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Error(errors.DSPControllerFailed(core.RoleConsumer, `RequestCatalog`, err))
		return
	}

	h.sendAck(w, catalog.RequestCatalogEndpoint, cat, http.StatusOK)
}

func (h *Handler) RequestDataset(w http.ResponseWriter, r *http.Request) {
	body, err := h.readBody(catalog.RequestDatasetEndpoint, w, r)
	if err != nil {
		return
	}

	var req catalog.DatasetRequest
	if err = json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Error(errors.UnmarshalError(catalog.RequestDatasetEndpoint, err))
		return
	}

	ds, err := h.consumer.RequestDataset(req.DatasetId, req.ProviderEndpoint)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.log.Error(errors.DSPControllerFailed(core.RoleConsumer, `RequestDataset`, err))
		return
	}

	h.sendAck(w, catalog.RequestDatasetEndpoint, ds, http.StatusOK)
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
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(errors.MarshalError(receivedEndpoint, err))
		return
	}

	w.WriteHeader(code)
	if _, err = w.Write(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error(errors.WriteBodyError(receivedEndpoint, err))
	}
}
