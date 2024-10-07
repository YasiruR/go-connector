package domain

import (
	"github.com/YasiruR/go-connector/domain/api"
	"github.com/YasiruR/go-connector/domain/control-plane"
	"github.com/YasiruR/go-connector/domain/pkg"
	"github.com/YasiruR/go-connector/domain/stores"
)

type Roles struct {
	control_plane.Provider
	control_plane.Consumer
	control_plane.Owner
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
	DSP       api.Server
	Exchanger api.Server
	Gateway   api.Server
}

type Plugins struct {
	pkg.Client
	pkg.Store
	pkg.URNService
	pkg.Log
}
