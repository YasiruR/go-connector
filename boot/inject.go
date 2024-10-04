package boot

import (
	dspSHttp "github.com/YasiruR/go-connector/api/dsp/http"
	gatewayHttp "github.com/YasiruR/go-connector/api/gateway/http"
	"github.com/YasiruR/go-connector/core/consumer"
	"github.com/YasiruR/go-connector/core/owner"
	"github.com/YasiruR/go-connector/core/provider"
	"github.com/YasiruR/go-connector/data"
	"github.com/YasiruR/go-connector/domain"
	"github.com/YasiruR/go-connector/pkg/client/http"
	"github.com/YasiruR/go-connector/pkg/database/memory"
	pkgLog "github.com/YasiruR/go-connector/pkg/log"
	"github.com/YasiruR/go-connector/pkg/urn"
	"github.com/YasiruR/go-connector/stores/catalog"
	"github.com/YasiruR/go-connector/stores/policy"
	"github.com/YasiruR/go-connector/stores/protocol"
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
	ProviderCatalog:          catalog.NewProviderCatalog(config, plugins),
	ConsumerCatalog:          catalog.NewConsumerCatalog(plugins),
	OfferStore:               policy.NewOfferStore(plugins),
	ContractNegotiationStore: protocol.NewContractNegotiationStore(plugins),
	AgreementStore:           policy.NewAgreementStore(plugins),
	TransferStore:            protocol.NewTransferStore(plugins),
}

var exchanger = data.NewExchanger(config, stores, plugins.Log)

var roles = domain.Roles{
	Provider: provider.New(config, stores, plugins),
	Consumer: consumer.New(config, stores, plugins),
	Owner:    owner.New(config, stores, plugins),
}

var servers = domain.Servers{
	DSP:     dspSHttp.NewServer(config.Servers.DSP.HTTP.Port, roles, plugins.Log),
	Gateway: gatewayHttp.NewServer(config.Servers.Gateway.HTTP.Port, roles, stores, plugins.Log),
}
