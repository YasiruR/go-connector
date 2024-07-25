package core

import (
	"github.com/YasiruR/connector/domain/core/consumer"
	"github.com/YasiruR/connector/domain/core/provider"
	"github.com/YasiruR/connector/domain/models/odrl"
)

const (
	RoleProvider = `Provider`
	RoleConsumer = `Consumer`
	RoleOwner    = `Owner`
)

type Provider interface {
	provider.CatalogHandler
	provider.NegotiationController
	provider.NegotiationHandler
	provider.TransferHandler
}

type Consumer interface {
	consumer.CatalogController
	consumer.NegotiationController
	consumer.NegotiationHandler
	consumer.TransferController
}

type Owner interface {
	CreatePolicy(target string, permissions, prohibitions []odrl.Rule) (policyId string, err error)
	CreateDataset(title string, descriptions, keywords, endpoints, policyIds []string) (datasetId string, err error)
}
