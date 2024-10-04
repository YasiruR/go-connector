package core

import (
	"github.com/YasiruR/go-connector/domain/core/consumer"
	"github.com/YasiruR/go-connector/domain/core/provider"
	"github.com/YasiruR/go-connector/domain/models/odrl"
)

const (
	RoleProvider = `provider`
	RoleConsumer = `consumer`
	RoleOwner    = `owner`
)

type Provider interface {
	provider.CatalogHandler
	provider.NegotiationController
	provider.NegotiationHandler
	provider.TransferController
	provider.TransferHandler
}

type Consumer interface {
	consumer.CatalogController
	consumer.NegotiationController
	consumer.NegotiationHandler
	consumer.TransferController
	consumer.TransferHandler
}

type Owner interface {
	CreatePolicy(target string, permissions, prohibitions []odrl.Rule) (id string, err error)
	CreateDataset(title, format string, descriptions, keywords, endpoints, offerIds []string) (id string, err error)
}
