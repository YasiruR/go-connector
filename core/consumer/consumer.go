package consumer

import (
	"github.com/YasiruR/go-connector/core/consumer/catalog"
	"github.com/YasiruR/go-connector/core/consumer/negotiation"
	"github.com/YasiruR/go-connector/core/consumer/transfer"
	"github.com/YasiruR/go-connector/domain"
	"github.com/YasiruR/go-connector/domain/boot"
	"github.com/YasiruR/go-connector/domain/core/consumer"
)

// validator should verify states before transitioning into next, signatures, authorization

type Consumer struct {
	consumer.CatalogController
	consumer.NegotiationController
	consumer.NegotiationHandler
	consumer.TransferController
	consumer.TransferHandler
}

func New(cfg boot.Config, stores domain.Stores, plugins domain.Plugins) *Consumer {
	return &Consumer{
		CatalogController:     catalog.NewController(stores, plugins.Client, plugins),
		NegotiationController: negotiation.NewController(cfg, stores, plugins),
		NegotiationHandler:    negotiation.NewHandler(stores, plugins),
		TransferController:    transfer.NewController(cfg, stores, plugins),
		TransferHandler:       transfer.NewHandler(stores.TransferStore, plugins.Log),
	}
}
