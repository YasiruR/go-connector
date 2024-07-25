package core

import (
	"github.com/YasiruR/connector/domain/core/consumer"
	"github.com/YasiruR/connector/domain/core/provider"
	"github.com/YasiruR/connector/domain/models/odrl"
)

// maybe implement submodules in Provider and Consumer to separate control and data plane

const (
	RoleProvider = `Provider`
	RoleConsumer = `Consumer`
	RoleOwner    = `Owner`
)

type Provider interface {
	provider.CatalogService
	provider.NegotiationController
	provider.NegotiationHandler
}

type Consumer interface {
	consumer.CatalogController
	consumer.NegotiationController
	consumer.NegotiationHandler
}

type Owner interface {
	CreatePolicy(target string, permissions, prohibitions []odrl.Rule) (policyId string, err error)
	CreateDataset(title string, descriptions, keywords, endpoints, policyIds []string) (datasetId string, err error)
}
