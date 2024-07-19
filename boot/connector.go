package boot

import (
	coreApi "github.com/YasiruR/connector/core/api/dsp"
	"github.com/YasiruR/connector/core/api/gateway"
	"github.com/YasiruR/connector/core/dsp"
	"github.com/YasiruR/connector/core/dsp/catalog"
	"github.com/YasiruR/connector/core/stores"
)

type connector struct {
	provider         dsp.Provider
	consumer         dsp.Consumer
	owner            dsp.Owner
	catalog          catalog.Service
	negotiationStore stores.ContractNegotiation
	dspServer        coreApi.HTTPServer
	gatewayServer    gateway.HTTPServer
}

type plugins struct{}
