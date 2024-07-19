package boot

import (
	dspServer "github.com/YasiruR/connector/api/dsp/http"
	managementServer "github.com/YasiruR/connector/api/gateway/http"
	"github.com/YasiruR/connector/dsp/consumer"
	"github.com/YasiruR/connector/dsp/provider"
	"github.com/YasiruR/connector/pkg/client/http"
	"github.com/YasiruR/connector/pkg/log"
	"github.com/YasiruR/connector/stores"
)

func Start() {
	hc := http.NewClient()
	cnStore := stores.NewContractNegotiationStore() // initialize all stores in single execution
	logger := log.NewLogger()

	cons := consumer.New(dspPort, cnStore, hc, logger)
	prov := provider.New(dspPort, cnStore, hc, logger)

	ds := dspServer.NewServer(dspPort, prov, logger)
	ms := managementServer.NewServer(managementPort, cons, logger)

	go ds.Start()
	ms.Start()
}
