package main

import (
	dspServer "github.com/YasiruR/connector/api/dsp/http"
	managementServer "github.com/YasiruR/connector/api/management/http"
	"github.com/YasiruR/connector/pkg/client/http"
	"github.com/YasiruR/connector/services/consumer"
	"github.com/YasiruR/connector/services/provider"
	"github.com/tryfix/log"
)

const (
	consumerDSPPort        = 8080
	consumerManagementPort = 8081
)

func main() {
	hc := http.NewClient()
	logger := log.Constructor.Log(log.WithColors(true), log.WithLevel("DEBUG"), log.WithFilePath(true))

	cons := consumer.New(consumerDSPPort, hc)
	prov := provider.New(hc)

	ds := dspServer.NewServer(consumerDSPPort, prov, logger)
	ms := managementServer.NewServer(consumerManagementPort, cons, logger)

	go ds.Start()
	ms.Start()
}
