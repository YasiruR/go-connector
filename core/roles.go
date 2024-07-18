package core

import "github.com/YasiruR/connector/protocols/negotiation"

type Provider interface {
	Owner
	negotiation.Provider
	negotiation.ProviderHandler
}

type Consumer interface {
	negotiation.Consumer
}

type Owner interface {
	CreateAsset()
	CreatePolicy()
	CreateContractDef()
}
