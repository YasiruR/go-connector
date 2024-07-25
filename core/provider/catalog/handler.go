package catalog

import (
	"github.com/YasiruR/connector/domain/api/dsp/http/catalog"
	"github.com/YasiruR/connector/domain/core"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/stores"
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
		Context:             core.Context,
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
		Context: core.Context,
		Dataset: ds,
	}, nil
}
