package http

import (
	httpCatalog "github.com/YasiruR/go-connector/api/dsp/http/catalog"
	httpNegotiation "github.com/YasiruR/go-connector/api/dsp/http/negotiation"
	httpTransfer "github.com/YasiruR/go-connector/api/dsp/http/transfer"
	"github.com/YasiruR/go-connector/domain"
	"github.com/YasiruR/go-connector/domain/api/dsp/http/catalog"
	"github.com/YasiruR/go-connector/domain/api/dsp/http/negotiation"
	"github.com/YasiruR/go-connector/domain/api/dsp/http/transfer"
	"github.com/YasiruR/go-connector/domain/errors"
	"github.com/YasiruR/go-connector/domain/pkg"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"strconv"
)

// dsp.http.Server contains the endpoints defined in data-plane space protocols which will be used
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
	r.HandleFunc(negotiation.ContractRequestToOfferEndpoint, s.nh.HandleContractRequest).Methods(http.MethodPost)
	r.HandleFunc(negotiation.ContractOfferEndpoint, s.nh.HandleContractOffer).Methods(http.MethodPost)
	r.HandleFunc(negotiation.ContractOfferToRequestEndpoint, s.nh.HandleContractOffer).Methods(http.MethodPost)
	r.HandleFunc(negotiation.AgreementVerificationEndpoint, s.nh.HandleAgreementVerification).Methods(http.MethodPost)
	r.HandleFunc(negotiation.ContractAgreementEndpoint, s.nh.HandleContractAgreement).Methods(http.MethodPost)
	r.HandleFunc(negotiation.EventsEndpoint, s.nh.HandleNegotiationEvent).Methods(http.MethodPost)
	r.HandleFunc(negotiation.TerminateEndpoint, s.nh.HandleTermination).Methods(http.MethodPost)

	// transfer protocol related endpoints
	r.HandleFunc(transfer.GetProcessEndpoint, s.th.HandleGetProcess).Methods(http.MethodGet)
	r.HandleFunc(transfer.RequestEndpoint, s.th.HandleTransferRequest).Methods(http.MethodPost)
	r.HandleFunc(transfer.StartEndpoint, s.th.HandleTransferStart).Methods(http.MethodPost)
	r.HandleFunc(transfer.SuspendEndpoint, s.th.HandleTransferSuspension).Methods(http.MethodPost)
	r.HandleFunc(transfer.CompleteEndpoint, s.th.HandleTransferCompletion).Methods(http.MethodPost)
	r.HandleFunc(transfer.TerminateEndpoint, s.th.HandleTransferTermination).Methods(http.MethodPost)

	return &s
}

func (s *Server) Start() {
	s.log.Info("DSP HTTP server is listening on " + strconv.Itoa(s.port))
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost", "http://127.0.0.1"},
		AllowCredentials: true,
	})

	if err := http.ListenAndServe(":"+strconv.Itoa(s.port), c.Handler(s.router)); err != nil {
		s.log.Fatal(errors.ModuleInitFailed(`DSP API`, err))
	}
}
