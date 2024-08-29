package boot

import (
	dspSHttp "github.com/YasiruR/connector/api/dsp/http"
	gatewayHttp "github.com/YasiruR/connector/api/gateway/http"
	"github.com/YasiruR/connector/core/consumer"
	"github.com/YasiruR/connector/core/owner"
	"github.com/YasiruR/connector/core/provider"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/pkg/client/http"
	"github.com/YasiruR/connector/pkg/database/memory"
	pkgLog "github.com/YasiruR/connector/pkg/log"
	"github.com/YasiruR/connector/pkg/urn"
	pkgStores "github.com/YasiruR/connector/stores"
)

var log = pkgLog.NewLogger()

var config = loadConfig(log)

var plugins = domain.Plugins{
	Client:     http.NewClient(log),
	Database:   memory.NewStore(log),
	URNService: urn.NewGenerator(),
	Log:        log,
}

var stores = domain.Stores{
	ProviderCatalog:          pkgStores.NewCatalog(config, plugins),
	PolicyStore:              pkgStores.NewPolicyStore(plugins),
	ContractNegotiationStore: pkgStores.NewContractNegotiationStore(plugins),
	AgreementStore:           pkgStores.NewAgreementStore(plugins),
	TransferStore:            pkgStores.NewTransferStore(plugins),
}

var roles = domain.Roles{
	Provider: provider.New(config, stores, plugins),
	Consumer: consumer.New(config.Servers.DSP.HTTP.Port, stores, plugins),
	Owner:    owner.New(config, stores, plugins),
}

var servers = domain.Servers{
	DSP:     dspSHttp.NewServer(config.Servers.DSP.HTTP.Port, roles, plugins.Log),
	Gateway: gatewayHttp.NewServer(config.Servers.Gateway.HTTP.Port, roles, stores, plugins.Log),
}
