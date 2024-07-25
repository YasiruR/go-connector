package http

import (
	"github.com/YasiruR/connector/api/gateway/http/catalog"
	"github.com/YasiruR/connector/api/gateway/http/negotiation"
	"github.com/YasiruR/connector/domain"
	httpCatalog "github.com/YasiruR/connector/domain/api/gateway/http/catalog"
	httpNegotiation "github.com/YasiruR/connector/domain/api/gateway/http/negotiation"
	"github.com/YasiruR/connector/domain/errors"
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
	ch     httpCatalog.Handler
	nh     httpNegotiation.Handler
	log    pkg.Log
}

func NewServer(port int, roles domain.Roles, stores domain.Stores, log pkg.Log) *Server {
	r := mux.NewRouter()
	s := Server{
		port:   port,
		router: r,
		ch:     catalog.NewHandler(roles, log),
		nh:     negotiation.NewHandler(roles, stores, log),
		log:    log,
	}

	// endpoints related to catalog
	r.HandleFunc(httpCatalog.CreatePolicyEndpoint, s.ch.CreatePolicy).Methods(http.MethodPost)
	r.HandleFunc(httpCatalog.CreateDatasetEndpoint, s.ch.CreateDataset).Methods(http.MethodPost)
	r.HandleFunc(httpCatalog.RequestCatalogEndpoint, s.ch.RequestCatalog).Methods(http.MethodPost)
	r.HandleFunc(httpCatalog.RequestDatasetEndpoint, s.ch.RequestDataset).Methods(http.MethodPost)

	// endpoints related to negotiation
	r.HandleFunc(httpNegotiation.RequestContractEndpoint, s.nh.RequestContract).Methods(http.MethodPost)
	r.HandleFunc(httpNegotiation.AgreeContractEndpoint, s.nh.AgreeContract).Methods(http.MethodPost)
	r.HandleFunc(httpNegotiation.GetAgreementEndpoint, s.nh.GetAgreement).Methods(http.MethodGet)
	r.HandleFunc(httpNegotiation.VerifyAgreementEndpoint, s.nh.VerifyAgreement).Methods(http.MethodPost)
	r.HandleFunc(httpNegotiation.FinalizeContractEndpoint, s.nh.FinalizeContract).Methods(http.MethodPost)

	return &s
}

func (s *Server) Start() {
	s.log.Info("gateway HTTP server is listening on " + strconv.Itoa(s.port))
	if err := http.ListenAndServe(":"+strconv.Itoa(s.port), s.router); err != nil {
		s.log.Fatal(errors.InitFailed(`gateway API`, err))
	}
}
