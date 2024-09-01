package catalog

import (
	"encoding/json"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/dsp/http/catalog"
	"github.com/YasiruR/connector/domain/core"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/stores"
)

type Controller struct {
	catalog stores.ConsumerCatalog
	client  pkg.Client
	log     pkg.Log
}

func NewController(s domain.Stores, client pkg.Client, log pkg.Log) *Controller {
	return &Controller{client: client, catalog: s.ConsumerCatalog, log: log}
}

func (c *Controller) RequestCatalog(endpoint string) (catalog.Response, error) {
	req := catalog.Request{
		Context:      core.Context,
		Type:         catalog.MsgTypRequest,
		DspaceFilter: nil,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return catalog.Response{}, errors.MarshalError(endpoint, err)
	}

	res, err := c.client.Send(data, endpoint+catalog.RequestEndpoint)
	if err != nil {
		return catalog.Response{}, errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var cat catalog.Response
	if err = json.Unmarshal(res, &cat); err != nil {
		return catalog.Response{}, errors.UnmarshalError(``, err)
	}

	c.catalog.AddCatalog(cat)
	c.log.Trace("stored the requested catalog", cat.ID)
	return cat, nil
}

func (c *Controller) RequestDataset(id, endpoint string) (catalog.DatasetResponse, error) {
	req := catalog.DatasetRequest{
		Context:   core.Context,
		Type:      catalog.MsgTypDatasetRequest,
		DatasetId: id,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return catalog.DatasetResponse{}, errors.MarshalError(endpoint, err)
	}

	res, err := c.client.Send(data, endpoint+catalog.RequestDatasetEndpoint)
	if err != nil {
		return catalog.DatasetResponse{}, errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var dataset catalog.DatasetResponse
	if err = json.Unmarshal(res, &dataset); err != nil {
		return catalog.DatasetResponse{}, errors.UnmarshalError(``, err)
	}

	return dataset, nil
}
