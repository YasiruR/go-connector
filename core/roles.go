package core

import "github.com/YasiruR/connector/protocols/negotiation"

type Provider interface {
	DataOwner
	negotiation.Provider
}

type Consumer interface {
	negotiation.Consumer
}

type DataOwner interface {
	CreateAsset()
	CreatePolicy()
	CreateContractDef()
}
