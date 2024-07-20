package dsp

import (
	"github.com/YasiruR/connector/core/dsp/negotiation"
	"github.com/YasiruR/connector/core/protocols/odrl"
)

// maybe implement submodules in Provider and Consumer to separate control and data plane

type Provider interface {
	negotiation.Provider
	negotiation.ProviderHandler
}

type Consumer interface {
	negotiation.Consumer
}

type Owner interface {
	CreatePolicy(t odrl.Target, permissions, prohibitions []odrl.Rule) (policyId string, err error)
	CreateDataset(title string, descriptions, keywords, endpoints, policyIds []string) (datasetId string, err error)
}
