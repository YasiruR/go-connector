package catalog

import (
	"encoding/json"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/gateway"
	"github.com/YasiruR/connector/domain/dsp"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/models/odrl"
	"github.com/YasiruR/connector/domain/pkg"
	"io"
	"net/http"
)

type Handler struct {
	consumer dsp.Consumer
	owner    dsp.Owner
	log      pkg.Log
}

func NewHandler(roles domain.Roles, log pkg.Log) *Handler {
	return &Handler{
		consumer: roles.Consumer,
		owner:    roles.Owner,
		log:      log,
	}
}

func (c *Handler) CreatePolicy(w http.ResponseWriter, r *http.Request) {
	body, err := c.readBody(gateway.CreatePolicyEndpoint, w, r)
	if err != nil {
		return
	}

	var req gateway.CreatePolicyRequest
	if err = json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.log.Error(errors.UnmarshalError(gateway.CreatePolicyEndpoint, err))
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
	id, err := c.owner.CreatePolicy(`test`, perms, []odrl.Rule{})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.log.Error(errors.DSPFailed(dsp.RoleOwner, `CreatePolicy`, err))
		return
	}

	c.sendAck(w, gateway.CreatePolicyEndpoint, gateway.PolicyResponse{Id: id}, http.StatusOK)
}

func (c *Handler) CreateDataset(w http.ResponseWriter, r *http.Request) {
	body, err := c.readBody(gateway.CreateDatasetEndpoint, w, r)
	if err != nil {
		return
	}

	var req gateway.CreateDatasetRequest
	if err = json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.log.Error(errors.UnmarshalError(gateway.CreateDatasetEndpoint, err))
		return
	}

	id, err := c.owner.CreateDataset(req.Title, req.Descriptions, req.Keywords, req.Endpoints, req.OfferIds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.log.Error(errors.DSPFailed(dsp.RoleOwner, `CreateDataset`, err))
		return
	}

	c.sendAck(w, gateway.CreateDatasetEndpoint, gateway.DatasetResponse{Id: id}, http.StatusOK)
}

func (c *Handler) RequestCatalog(w http.ResponseWriter, r *http.Request) {
	body, err := c.readBody(gateway.RequestCatalogEndpoint, w, r)
	if err != nil {
		return
	}

	var req gateway.CatalogRequest
	if err = json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.log.Error(errors.UnmarshalError(gateway.RequestCatalogEndpoint, err))
		return
	}

	cat, err := c.consumer.RequestCatalog(req.ProviderEndpoint)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.log.Error(errors.DSPFailed(dsp.RoleConsumer, `RequestCatalog`, err))
		return
	}

	c.sendAck(w, gateway.RequestCatalogEndpoint, cat, http.StatusOK)
}

func (c *Handler) RequestDataset(w http.ResponseWriter, r *http.Request) {
	body, err := c.readBody(gateway.RequestDatasetEndpoint, w, r)
	if err != nil {
		return
	}

	var req gateway.DatasetRequest
	if err = json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.log.Error(errors.UnmarshalError(gateway.RequestDatasetEndpoint, err))
		return
	}

	ds, err := c.consumer.RequestDataset(req.DatasetId, req.ProviderEndpoint)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.log.Error(errors.DSPFailed(dsp.RoleConsumer, `RequestDataset`, err))
		return
	}

	c.sendAck(w, gateway.RequestDatasetEndpoint, ds, http.StatusOK)
}

func (c *Handler) readBody(endpoint string, w http.ResponseWriter, r *http.Request) ([]byte, error) {
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

func (c *Handler) sendAck(w http.ResponseWriter, receivedEndpoint string, data any, code int) {
	body, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.log.Error(errors.MarshalError(receivedEndpoint, err))
		return
	}

	w.WriteHeader(code)
	if _, err = w.Write(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.log.Error(errors.WriteBodyError(receivedEndpoint, err))
	}
}
