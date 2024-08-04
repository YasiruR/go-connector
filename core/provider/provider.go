package provider

import (
	"github.com/YasiruR/connector/core/provider/catalog"
	"github.com/YasiruR/connector/core/provider/negotiation"
	"github.com/YasiruR/connector/core/provider/transfer"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/core/provider"
)

type Provider struct {
	provider.CatalogHandler
	provider.NegotiationController
	provider.NegotiationHandler
	provider.TransferController
	provider.TransferHandler
}

func New(port int, stores domain.Stores, plugins domain.Plugins) *Provider {
	return &Provider{
		CatalogHandler:        catalog.NewHandler(stores.Catalog, plugins.Log),
		NegotiationController: negotiation.NewController(port, stores, plugins),
		NegotiationHandler:    negotiation.NewHandler(stores.ContractNegotiation, plugins),
		TransferController:    transfer.NewController(stores.Transfer, plugins),
		TransferHandler:       transfer.NewHandler(stores, plugins),
	}
}
