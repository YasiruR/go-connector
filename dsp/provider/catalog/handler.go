package catalog

import (
	"github.com/YasiruR/connector/core/dsp"
	"github.com/YasiruR/connector/core/dsp/catalog"
	"github.com/YasiruR/connector/core/errors"
	"github.com/YasiruR/connector/core/pkg"
	"github.com/YasiruR/connector/core/stores"
)

// Catalog Protocol (reference: https://docs.internationaldataspaces.org/ids-knowledgebase/v/dataspace-protocol/catalog/catalog.protocol)

type Handler struct {
	participantId string // data space specific identifier for Provider
	catalog       stores.Catalog
	log           pkg.Log
}

func NewHandler(cnStore stores.Catalog, log pkg.Log) *Handler {
	return &Handler{
		participantId: `participant-id-provider`,
		catalog:       cnStore,
		log:           log,
	}
}

func (h *Handler) HandleCatalogRequest(_ any) (catalog.Response, error) {
	cat, err := h.catalog.Get()
	if err != nil {
		return catalog.Response{}, errors.StoreFailed(stores.TypeCatalog, `Get`, err)
	}

	return catalog.Response{
		Context:             dsp.Context,
		DspaceParticipantID: h.participantId,
		Catalog:             cat,
	}, nil
}

func (h *Handler) HandleDatasetRequest(id string) (catalog.DatasetResponse, error) {
	ds, err := h.catalog.Dataset(id)
	if err != nil {
		return catalog.DatasetResponse{}, errors.StoreFailed(stores.TypeCatalog, `Dataset`, err)
	}

	return catalog.DatasetResponse{
		Context: dsp.Context,
		Dataset: ds,
	}, nil
}
