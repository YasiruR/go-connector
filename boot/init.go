package boot

import (
	dspSHttp "github.com/YasiruR/connector/api/dsp/http"
	gatewayHttp "github.com/YasiruR/connector/api/gateway/http"
	"github.com/YasiruR/connector/dsp/consumer"
	"github.com/YasiruR/connector/dsp/owner"
	"github.com/YasiruR/connector/dsp/provider"
	"github.com/YasiruR/connector/pkg/client/http"
	"github.com/YasiruR/connector/pkg/database/memory"
	"github.com/YasiruR/connector/pkg/log"
	"github.com/YasiruR/connector/pkg/urn"
	"github.com/YasiruR/connector/stores"
)

func Start() {
	client := http.NewClient()
	memDb := memory.NewStore()
	urnGen := urn.NewGenerator()
	logger := log.NewLogger()

	cnStore := stores.NewContractNegotiationStore(memDb)
	polStore := stores.NewPolicyStore(memDb)
	catalog := stores.NewCatalog(urnGen, memDb)

	cons := consumer.New(dspPort, cnStore, urnGen, client, logger)
	prov := provider.New(dspPort, cnStore, polStore, catalog, urnGen, client, logger)
	ownr := owner.New(polStore, catalog, urnGen, logger)

	dspSvr := dspSHttp.NewServer(dspPort, prov, cons, logger)
	gatewaySvr := gatewayHttp.NewServer(managementPort, prov, cons, ownr, logger)

	go dspSvr.Start()
	gatewaySvr.Start()
}
