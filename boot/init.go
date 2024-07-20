package boot

import (
	dspServer "github.com/YasiruR/connector/api/dsp/http"
	managementServer "github.com/YasiruR/connector/api/gateway/http"
	"github.com/YasiruR/connector/dsp/consumer"
	"github.com/YasiruR/connector/dsp/owner"
	"github.com/YasiruR/connector/dsp/provider"
	"github.com/YasiruR/connector/pkg/client/http"
	"github.com/YasiruR/connector/pkg/database/memory"
	"github.com/YasiruR/connector/pkg/log"
	"github.com/YasiruR/connector/stores"
)

func Start() {
	hc := http.NewClient()
	memDb := memory.NewStore()
	logger := log.NewLogger()

	cnStore := stores.NewContractNegotiationStore(memDb)
	pStore := stores.NewPolicyStore(memDb)
	dsStore := stores.NewDatasetStore(memDb)

	cons := consumer.New(dspPort, cnStore, hc, logger)
	prov := provider.New(dspPort, cnStore, hc, logger)
	ownr := owner.New(pStore, dsStore, logger)

	ds := dspServer.NewServer(dspPort, prov, logger)
	ms := managementServer.NewServer(managementPort, cons, ownr, logger)

	go ds.Start()
	ms.Start()
}
