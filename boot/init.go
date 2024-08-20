package boot

import (
	dspSHttp "github.com/YasiruR/connector/api/dsp/http"
	gatewayHttp "github.com/YasiruR/connector/api/gateway/http"
	"github.com/YasiruR/connector/core/consumer"
	"github.com/YasiruR/connector/core/owner"
	"github.com/YasiruR/connector/core/provider"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/boot"
	"github.com/YasiruR/connector/domain/pkg"
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
	stors := initStores(cfg, plugins)
	roles := initRoles(cfg, stors, plugins)
	log.Info("enabled consumer, provider and owner roles")
	servers := initServers(cfg, roles, stors, plugins)

	go servers.DSP.Start()
	servers.Gateway.Start()
}

func initPlugins(log pkg.Log) domain.Plugins {
	return domain.Plugins{
		Client:     http.NewClient(log),
		Database:   memory.NewStore(log),
		URNService: urn.NewGenerator(),
		Log:        log,
	}
}

func initStores(cfg boot.Config, plugins domain.Plugins) domain.Stores {
	return domain.Stores{
		ProviderCatalog:          stores.NewCatalog(cfg, plugins),
		PolicyStore:              stores.NewPolicyStore(plugins),
		ContractNegotiationStore: stores.NewContractNegotiationStore(plugins),
		AgreementStore:           stores.NewAgreementStore(plugins),
		TransferStore:            stores.NewTransferStore(plugins),
	}
}

func initRoles(cfg boot.Config, stores domain.Stores, plugins domain.Plugins) domain.Roles {
	return domain.Roles{
		Provider: provider.New(cfg, stores, plugins),
		Consumer: consumer.New(cfg.Servers.DSP.HTTP.Port, stores, plugins),
		Owner:    owner.New(cfg, stores, plugins),
	}
}

func initServers(cfg boot.Config, roles domain.Roles, stores domain.Stores, plugins domain.Plugins) domain.Servers {
	return domain.Servers{
		DSP:     dspSHttp.NewServer(cfg.Servers.DSP.HTTP.Port, roles, plugins.Log),
		Gateway: gatewayHttp.NewServer(cfg.Servers.Gateway.HTTP.Port, roles, stores, plugins.Log),
	}
}
