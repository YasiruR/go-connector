package catalog

import (
	"github.com/YasiruR/go-connector/domain/api/dsp/http/catalog"
	"github.com/YasiruR/go-connector/domain/control-plane"
	"github.com/YasiruR/go-connector/domain/errors"
	"github.com/YasiruR/go-connector/domain/pkg"
	"github.com/YasiruR/go-connector/domain/stores"
)

// Catalog Protocol (reference: https://docs.internationaldataspaces.org/ids-knowledgebase/v/dataspace-protocol/catalog/catalog.protocol)

type Handler struct {
	participantId string // data-plane space specific identifier for Provider
	catStore      stores.ProviderCatalog
	log           pkg.Log
}

func NewHandler(participantId string, cnStore stores.ProviderCatalog, log pkg.Log) *Handler {
	return &Handler{
		participantId: participantId,
		catStore:      cnStore,
		log:           log,
	}
}

func (h *Handler) HandleCatalogRequest(_ any) (catalog.Response, error) {
	cat, err := h.catStore.Catalog()
	if err != nil {
		return catalog.Response{}, errors.StoreFailed(stores.TypeProviderCatalog, `Get`, err)
	}

	return catalog.Response{
		Context:             control_plane.Context,
		DspaceParticipantID: h.participantId,
		Catalog:             cat,
	}, nil
}

func (h *Handler) HandleDatasetRequest(id string) (catalog.DatasetResponse, error) {
	ds, err := h.catStore.Dataset(id)
	if err != nil {
		return catalog.DatasetResponse{}, errors.Catalog(errors.InvalidKey(
			stores.TypeProviderCatalog, `dataset id`, err))
	}

	return catalog.DatasetResponse{
		Context: control_plane.Context,
		Dataset: ds,
	}, nil
}
