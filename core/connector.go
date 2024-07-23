package core

import (
	dsp2 "github.com/YasiruR/connector/core/api/dsp"
	"github.com/YasiruR/connector/core/api/gateway"
	"github.com/YasiruR/connector/core/dsp"
	"github.com/YasiruR/connector/core/pkg"
	"github.com/YasiruR/connector/core/stores"
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
