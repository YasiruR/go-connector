package http

import (
	httpCatalog "github.com/YasiruR/connector/api/dsp/http/catalog"
	httpNegotiation "github.com/YasiruR/connector/api/dsp/http/negotiation"
	httpTransfer "github.com/YasiruR/connector/api/dsp/http/transfer"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/dsp/http/catalog"
	"github.com/YasiruR/connector/domain/api/dsp/http/negotiation"
	"github.com/YasiruR/connector/domain/api/dsp/http/transfer"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// dsp.http.Server contains the endpoints defined in data space protocols which will be used
// for the communication between connectors

type Server struct {
	port   int
	ch     catalog.Handler
	nh     negotiation.Handler
	th     transfer.Handler
	router *mux.Router
	log    pkg.Log
}

func NewServer(port int, roles domain.Roles, log pkg.Log) *Server {
	r := mux.NewRouter()
	s := Server{
		port:   port,
		ch:     httpCatalog.NewHandler(roles, log),
		nh:     httpNegotiation.NewHandler(roles, log),
		th:     httpTransfer.NewHandler(roles, log),
		router: r,
		log:    log,
	}

	// catalog protocol related endpoints
	r.HandleFunc(catalog.RequestEndpoint, s.ch.HandleCatalogRequest).Methods(http.MethodPost)
	r.HandleFunc(catalog.RequestDatasetEndpoint, s.ch.HandleDatasetRequest).Methods(http.MethodPost)

	// negotiation protocol related endpoints
	r.HandleFunc(negotiation.RequestEndpoint, s.nh.GetNegotiation).Methods(http.MethodGet)
	r.HandleFunc(negotiation.ContractRequestEndpoint, s.nh.HandleContractRequest).Methods(http.MethodPost)
	r.HandleFunc(negotiation.ContractOfferEndpoint, s.nh.HandleContractOffer).Methods(http.MethodPost)
	r.HandleFunc(negotiation.ContractOfferToRequestEndpoint, s.nh.HandleContractOffer).Methods(http.MethodPost)
	r.HandleFunc(negotiation.ContractAgreementEndpoint, s.nh.HandleContractAgreement).Methods(http.MethodPost)
	r.HandleFunc(negotiation.AgreementVerificationEndpoint, s.nh.HandleAgreementVerification).Methods(http.MethodPost)
	r.HandleFunc(negotiation.EventConsumerEndpoint, s.nh.HandleEventConsumer).Methods(http.MethodPost)

	// transfer process related endpoints
	r.HandleFunc(transfer.RequestEndpoint, s.th.HandleTransferRequest).Methods(http.MethodPost)
	r.HandleFunc(transfer.StartEndpoint, s.th.HandleTransferStart).Methods(http.MethodPost)
	r.HandleFunc(transfer.SuspendEndpoint, s.th.HandleTransferSuspension).Methods(http.MethodPost)
	r.HandleFunc(transfer.CompleteEndpoint, s.th.HandleTransferCompletion).Methods(http.MethodPost)

	return &s
}

func (s *Server) Start() {
	s.log.Info("DSP HTTP server is listening on " + strconv.Itoa(s.port))
	if err := http.ListenAndServe(":"+strconv.Itoa(s.port), s.router); err != nil {
		s.log.Fatal(errors.InitModuleFailed(`DSP API`, err))
	}
}
