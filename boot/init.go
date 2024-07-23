package boot

import (
	dspSHttp "github.com/YasiruR/connector/api/dsp/http"
	gatewayHttp "github.com/YasiruR/connector/api/gateway/http"
	"github.com/YasiruR/connector/core"
	"github.com/YasiruR/connector/core/pkg"
	"github.com/YasiruR/connector/dsp"
	"github.com/YasiruR/connector/pkg/client/http"
	"github.com/YasiruR/connector/pkg/database/memory"
	pkgLog "github.com/YasiruR/connector/pkg/log"
	"github.com/YasiruR/connector/pkg/urn"
	"github.com/YasiruR/connector/stores"
)

func Start() {
	log := pkgLog.NewLogger()
	cfg := loadConfig(log)

	plugins := initPlugins(log)
	stors := initStores(plugins)
	roles := initRoles(cfg, stors, plugins)
	servers := initServers(cfg, roles, plugins)

	//stors.Init()
	go servers.DSP.Start()
	servers.Gateway.Start()
}

func initPlugins(log pkg.Log) core.Plugins {
	return core.Plugins{
		Client:     http.NewClient(log),
		Database:   memory.NewStore(log),
		URNService: urn.NewGenerator(),
		Log:        log,
	}
}

func initStores(plugins core.Plugins) core.Stores {
	return core.Stores{
		Catalog:             stores.NewCatalog(plugins),
		Policy:              stores.NewPolicyStore(plugins),
		ContractNegotiation: stores.NewContractNegotiationStore(plugins),
	}
}

func initRoles(cfg core.Config, stores core.Stores, plugins core.Plugins) core.Roles {
	return core.Roles{
		Provider: dsp.NewProvider(cfg.Servers.DSP.HTTP.Port, stores, plugins),
		Consumer: dsp.NewConsumer(cfg.Servers.DSP.HTTP.Port, stores, plugins),
		Owner:    dsp.NewOwner(stores, plugins),
	}
}

func initServers(cfg core.Config, roles core.Roles, plugins core.Plugins) core.Servers {
	return core.Servers{
		DSP:     dspSHttp.NewServer(cfg.Servers.DSP.HTTP.Port, roles, plugins.Log),
		Gateway: gatewayHttp.NewServer(cfg.Servers.Gateway.HTTP.Port, roles, plugins.Log),
	}
}
