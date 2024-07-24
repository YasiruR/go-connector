package catalog

import (
	"encoding/json"
	"github.com/YasiruR/connector/core/dsp"
	"github.com/YasiruR/connector/core/dsp/catalog"
	"github.com/YasiruR/connector/core/errors"
	"github.com/YasiruR/connector/core/pkg"
)

type Controller struct {
	client pkg.Client
}

func NewController(client pkg.Client) *Controller {
	return &Controller{client: client}
}

func (c *Controller) RequestCatalog(endpoint string) (catalog.Response, error) {
	req := catalog.Request{
		Context:      dsp.Context,
		Type:         catalog.TypeCatalogRequest,
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

	return cat, nil
}

func (c *Controller) RequestDataset(id, endpoint string) (catalog.DatasetResponse, error) {
	req := catalog.DatasetRequest{
		Context:   dsp.Context,
		Type:      catalog.TypeDatasetRequest,
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
