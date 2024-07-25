package catalog

import (
	"encoding/json"
	catalog2 "github.com/YasiruR/connector/domain/api/dsp/http/catalog"
	"github.com/YasiruR/connector/domain/core"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/pkg"
)

type Controller struct {
	client pkg.Client
}

func NewController(client pkg.Client) *Controller {
	return &Controller{client: client}
}

func (c *Controller) RequestCatalog(endpoint string) (catalog2.Response, error) {
	req := catalog2.Request{
		Context:      core.Context,
		Type:         catalog2.TypeCatalogRequest,
		DspaceFilter: nil,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return catalog2.Response{}, errors.MarshalError(endpoint, err)
	}

	res, err := c.client.Send(data, endpoint+catalog2.RequestEndpoint)
	if err != nil {
		return catalog2.Response{}, errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var cat catalog2.Response
	if err = json.Unmarshal(res, &cat); err != nil {
		return catalog2.Response{}, errors.UnmarshalError(``, err)
	}

	return cat, nil
}

func (c *Controller) RequestDataset(id, endpoint string) (catalog2.DatasetResponse, error) {
	req := catalog2.DatasetRequest{
		Context:   core.Context,
		Type:      catalog2.TypeDatasetRequest,
		DatasetId: id,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return catalog2.DatasetResponse{}, errors.MarshalError(endpoint, err)
	}

	res, err := c.client.Send(data, endpoint+catalog2.RequestDatasetEndpoint)
	if err != nil {
		return catalog2.DatasetResponse{}, errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var dataset catalog2.DatasetResponse
	if err = json.Unmarshal(res, &dataset); err != nil {
		return catalog2.DatasetResponse{}, errors.UnmarshalError(``, err)
	}

	return dataset, nil
}
