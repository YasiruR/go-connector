package provider

import (
	"github.com/YasiruR/go-connector/core/provider/catalog"
	"github.com/YasiruR/go-connector/core/provider/negotiation"
	"github.com/YasiruR/go-connector/core/provider/transfer"
	"github.com/YasiruR/go-connector/domain"
	"github.com/YasiruR/go-connector/domain/boot"
	"github.com/YasiruR/go-connector/domain/core/provider"
)

type Provider struct {
	provider.CatalogHandler
	provider.NegotiationController
	provider.NegotiationHandler
	provider.TransferController
	provider.TransferHandler
}

func New(cfg boot.Config, stores domain.Stores, plugins domain.Plugins) *Provider {
	return &Provider{
		CatalogHandler:        catalog.NewHandler(cfg.DataSpace.ParticipantId, stores.ProviderCatalog, plugins.Log),
		NegotiationController: negotiation.NewController(cfg, stores, plugins),
		NegotiationHandler:    negotiation.NewHandler(cfg, stores, plugins),
		TransferController:    transfer.NewController(stores.TransferStore, plugins),
		TransferHandler:       transfer.NewHandler(stores, plugins),
	}
}
