package dsp

import (
	"github.com/YasiruR/connector/core/dsp/negotiation"
	"github.com/YasiruR/connector/core/models/odrl"
)

type Provider interface {
	negotiation.Provider
	negotiation.ProviderHandler
}

type Consumer interface {
	negotiation.Consumer
}

type Owner interface {
	CreatePolicy(t odrl.Target, permissions, prohibitions []odrl.Rule) (policyId string, err error)
	CreateContractDefinition()
}
