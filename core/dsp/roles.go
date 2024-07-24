package dsp

import (
	"github.com/YasiruR/connector/core/dsp/catalog"
	"github.com/YasiruR/connector/core/dsp/negotiation"
	"github.com/YasiruR/connector/core/protocols/odrl"
)

// maybe implement submodules in Provider and Consumer to separate control and data plane

const (
	Context      = `https://w3id.org/dspace/2024/1/context.json`
	RoleProvider = `Provider`
	RoleConsumer = `Consumer`
	RoleOwner    = `Owner`
)

type Provider interface {
	catalog.Provider
	negotiation.Provider
}

type Consumer interface {
	catalog.Consumer
	negotiation.Consumer
}

type Owner interface {
	CreatePolicy(target string, permissions, prohibitions []odrl.Rule) (policyId string, err error)
	CreateDataset(title string, descriptions, keywords, endpoints, policyIds []string) (datasetId string, err error)
}
