package domain

import (
	"github.com/YasiruR/connector/domain/api"
	"github.com/YasiruR/connector/domain/core"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/stores"
)

type Roles struct {
	core.Provider
	core.Consumer
	core.Owner
}

type Stores struct {
	stores.ProviderCatalog
	stores.ConsumerCatalog
	stores.OfferStore
	stores.ContractNegotiationStore
	stores.AgreementStore
	stores.TransferStore
}

type Servers struct {
	DSP     api.Server
	Gateway api.Server
}

type Plugins struct {
	pkg.Client
	pkg.Database
	pkg.URNService
	pkg.Log
}
