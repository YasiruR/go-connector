package provider

import (
	"github.com/YasiruR/connector/core/provider/catalog"
	"github.com/YasiruR/connector/core/provider/negotiation"
	"github.com/YasiruR/connector/core/provider/transfer"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/boot"
	"github.com/YasiruR/connector/domain/core/provider"
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
		CatalogHandler:        catalog.NewHandler(stores.ProviderCatalog, plugins.Log),
		NegotiationController: negotiation.NewController(cfg.Servers.DSP.HTTP.Port, stores, plugins),
		NegotiationHandler:    negotiation.NewHandler(cfg, stores, plugins),
		TransferController:    transfer.NewController(stores.Transfer, plugins),
		TransferHandler:       transfer.NewHandler(stores, plugins),
	}
}
