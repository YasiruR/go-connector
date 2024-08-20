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
	ProviderCatalog stores.Catalog
	stores.Policy
	stores.ContractNegotiation
	stores.Agreement
	stores.Transfer
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
