package http

import (
	httpCatalog "github.com/YasiruR/connector/api/gateway/http/catalog"
	httpNegotiation "github.com/YasiruR/connector/api/gateway/http/negotiation"
	httpTransfer "github.com/YasiruR/connector/api/gateway/http/transfer"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/gateway/http/catalog"
	"github.com/YasiruR/connector/domain/api/gateway/http/negotiation"
	"github.com/YasiruR/connector/domain/api/gateway/http/transfer"
	"github.com/YasiruR/connector/domain/errors/core"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// gateway.http.Server contains the endpoints which will be used by a client to initiate
// message flows or manage both control and data planes

// todo check return error codes

type Server struct {
	port   int
	router *mux.Router
	ch     catalog.Handler
	nh     negotiation.Handler
	th     transfer.Handler
	log    pkg.Log
}

func NewServer(port int, roles domain.Roles, stores domain.Stores, log pkg.Log) *Server {
	r := mux.NewRouter()
	s := Server{
		port:   port,
		router: r,
		ch:     httpCatalog.NewHandler(roles, stores, log),
		nh:     httpNegotiation.NewHandler(roles, stores, log),
		th:     httpTransfer.NewHandler(roles, log),
		log:    log,
	}

	// endpoints related to catalog
	r.HandleFunc(catalog.CreatePolicyEndpoint, s.ch.CreatePolicy).Methods(http.MethodPost)
	r.HandleFunc(catalog.CreateDatasetEndpoint, s.ch.CreateDataset).Methods(http.MethodPost)
	r.HandleFunc(catalog.RequestCatalogEndpoint, s.ch.RequestCatalog).Methods(http.MethodPost)
	r.HandleFunc(catalog.RequestDatasetEndpoint, s.ch.RequestDataset).Methods(http.MethodPost)
	r.HandleFunc(catalog.GetStoredCatalogsEndpoint, s.ch.GetStoredCatalogs).Methods(http.MethodGet)

	// endpoints related to negotiation
	r.HandleFunc(negotiation.RequestContractEndpoint, s.nh.RequestContract).Methods(http.MethodPost)
	r.HandleFunc(negotiation.OfferContractEndpoint, s.nh.OfferContract).Methods(http.MethodPost)
	r.HandleFunc(negotiation.AcceptOfferEndpoint, s.nh.AcceptOffer).Methods(http.MethodPost)
	r.HandleFunc(negotiation.AgreeContractEndpoint, s.nh.AgreeContract).Methods(http.MethodPost)
	r.HandleFunc(negotiation.GetAgreementEndpoint, s.nh.GetAgreement).Methods(http.MethodGet)
	r.HandleFunc(negotiation.VerifyAgreementEndpoint, s.nh.VerifyAgreement).Methods(http.MethodPost)
	r.HandleFunc(negotiation.FinalizeContractEndpoint, s.nh.FinalizeContract).Methods(http.MethodPost)
	r.HandleFunc(negotiation.TerminateContractEndpoint, s.nh.TerminateContract).Methods(http.MethodPost)

	// endpoints related to transfer process
	r.HandleFunc(transfer.GetProcessEndpoint, s.th.GetProviderProcess).Methods(http.MethodGet)
	r.HandleFunc(transfer.RequestEndpoint, s.th.RequestTransfer).Methods(http.MethodPost)
	r.HandleFunc(transfer.StartEndpoint, s.th.StartTransfer).Methods(http.MethodPost)
	r.HandleFunc(transfer.SuspendEndpoint, s.th.SuspendTransfer).Methods(http.MethodPost)
	r.HandleFunc(transfer.CompleteEndpoint, s.th.CompleteTransfer).Methods(http.MethodPost)
	r.HandleFunc(transfer.TerminateEndpoint, s.th.TerminateTransfer).Methods(http.MethodPost)

	return &s
}

func (s *Server) Start() {
	s.log.Info("gateway HTTP server is listening on " + strconv.Itoa(s.port))
	if err := http.ListenAndServe(":"+strconv.Itoa(s.port), s.router); err != nil {
		s.log.Fatal(core.ModuleInitFailed(`gateway API`, err))
	}
}
