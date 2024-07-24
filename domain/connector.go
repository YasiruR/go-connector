package domain

import (
	dsp2 "github.com/YasiruR/connector/domain/api/dsp"
	"github.com/YasiruR/connector/domain/api/gateway"
	"github.com/YasiruR/connector/domain/dsp"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/stores"
)

type Roles struct {
	dsp.Provider
	dsp.Consumer
	dsp.Owner
}

type Stores struct {
	stores.Catalog
	stores.Policy
	stores.ContractNegotiation
	stores.Agreement
}

type Servers struct {
	DSP     dsp2.HTTPServer
	Gateway gateway.HTTPServer
}

type Plugins struct {
	pkg.Client
	pkg.Database
	pkg.URNService
	pkg.Log
}
