package http

import (
	"encoding/json"
	"github.com/YasiruR/connector/core"
	"github.com/YasiruR/connector/core/dsp"
	"github.com/YasiruR/connector/core/dsp/catalog"
	"github.com/YasiruR/connector/core/dsp/negotiation"
	"github.com/YasiruR/connector/core/errors"
	"github.com/YasiruR/connector/core/pkg"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

// dsp.http.Server contains the endpoints defined in data space protocols which will be used
// for the communication between connectors

type Server struct {
	port     int
	owner    dsp.Owner
	provider dsp.Provider
	consumer dsp.Consumer
	router   *mux.Router
	log      pkg.Log
}

func NewServer(port int, roles core.Roles, log pkg.Log) *Server {
	r := mux.NewRouter()
	s := Server{port: port, router: r, provider: roles.Provider, consumer: roles.Consumer, log: log}

	// catalog protocol related endpoints
	r.HandleFunc(catalog.RequestEndpoint, s.HandleCatalogRequest).Methods(http.MethodPost)
	r.HandleFunc(catalog.RequestDatasetEndpoint, s.HandleDatasetRequest).Methods(http.MethodPost)

	// negotiation protocol related endpoints
	r.HandleFunc(negotiation.RequestEndpoint, s.GetNegotiation).Methods(http.MethodGet)
	r.HandleFunc(negotiation.ContractRequestEndpoint, s.HandleContractRequest).Methods(http.MethodPost)
	r.HandleFunc(negotiation.ContractAgreementEndpoint, s.HandleContractAgreement).Methods(http.MethodPost)

	return &s
}

func (s *Server) Start() {
	s.log.Info("DSP HTTP server is listening on " + strconv.Itoa(s.port))
	if err := http.ListenAndServe(":"+strconv.Itoa(s.port), s.router); err != nil {
		s.log.Fatal(errors.InitFailed(`DSP API`, err))
	}
}

func (s *Server) HandleCatalogRequest(w http.ResponseWriter, r *http.Request) {
	body, err := s.readBody(catalog.RequestEndpoint, w, r)
	if err != nil {
		return
	}

	var req catalog.Request
	if err = json.Unmarshal(body, &req); err != nil {
		s.sendError(w, errors.UnmarshalError(catalog.RequestEndpoint, err), http.StatusBadRequest)
		return
	}

	cat, err := s.provider.HandleCatalogRequest(nil)
	if err != nil {
		s.sendError(w, errors.HandlerFailed(catalog.RequestEndpoint, dsp.RoleProvider, err), http.StatusBadRequest)
		return
	}

	s.sendAck(w, catalog.RequestEndpoint, cat, http.StatusOK)
}

func (s *Server) HandleDatasetRequest(w http.ResponseWriter, r *http.Request) {
	body, err := s.readBody(catalog.RequestDatasetEndpoint, w, r)
	if err != nil {
		return
	}

	var req catalog.DatasetRequest
	if err = json.Unmarshal(body, &req); err != nil {
		s.sendError(w, errors.UnmarshalError(catalog.TypeDatasetRequest, err), http.StatusBadRequest)
		return
	}

	ds, err := s.provider.HandleDatasetRequest(req.DatasetId)
	if err != nil {
		s.sendError(w, errors.HandlerFailed(catalog.RequestDatasetEndpoint, dsp.RoleProvider, err), http.StatusBadRequest)
		return
	}

	s.sendAck(w, catalog.RequestDatasetEndpoint, ds, http.StatusOK)
}

func (s *Server) GetNegotiation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	providerPid, ok := params[negotiation.ParamProviderId]
	if !ok {
		s.sendError(w, errors.PathParamNotFound(negotiation.RequestEndpoint, negotiation.ParamProviderId), http.StatusBadRequest)
		return
	}

	neg, err := s.provider.HandleNegotiationsRequest(providerPid)
	if err != nil {
		s.sendError(w, errors.HandlerFailed(negotiation.RequestEndpoint, dsp.RoleProvider, err), http.StatusBadRequest)
		return
	}

	s.sendAck(w, negotiation.RequestEndpoint, neg, http.StatusOK)
}

func (s *Server) HandleContractRequest(w http.ResponseWriter, r *http.Request) {
	body, err := s.readBody(catalog.RequestEndpoint, w, r)
	if err != nil {
		return
	}

	var req negotiation.ContractRequest
	if err = json.Unmarshal(body, &req); err != nil {
		s.sendError(w, errors.UnmarshalError(negotiation.ContractRequestEndpoint, err), http.StatusBadRequest)
		return
	}

	ack, err := s.provider.HandleContractRequest(req)
	if err != nil {
		s.sendError(w, errors.HandlerFailed(negotiation.ContractRequestEndpoint, dsp.RoleProvider, err), http.StatusBadRequest)
		return
	}

	s.sendAck(w, negotiation.ContractRequestEndpoint, ack, http.StatusCreated)
}

func (s *Server) HandleContractAgreement(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	consumerPid, ok := params[negotiation.ParamConsumerPid]
	if !ok {
		s.sendError(w, errors.PathParamNotFound(negotiation.ContractAgreementEndpoint, negotiation.ParamConsumerPid), http.StatusBadRequest)
		return
	}

	body, err := s.readBody(catalog.RequestEndpoint, w, r)
	if err != nil {
		return
	}

	var req negotiation.ContractAgreement
	if err = json.Unmarshal(body, &req); err != nil {
		s.sendError(w, errors.UnmarshalError(negotiation.ContractAgreementEndpoint, err), http.StatusBadRequest)
		return
	}

	ack, err := s.consumer.HandleContractAgreement(consumerPid, req)
	if err != nil {
		s.sendError(w, errors.HandlerFailed(negotiation.ContractAgreementEndpoint, dsp.RoleConsumer, err), http.StatusBadRequest)
		return
	}

	s.sendAck(w, negotiation.ContractAgreementEndpoint, ack, http.StatusOK)
}

func (s *Server) readBody(endpoint string, w http.ResponseWriter, r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		err = errors.InvalidRequestBody(endpoint, err)
		w.WriteHeader(http.StatusBadRequest)
		s.log.Error(err)
		r.Body.Close()
		return nil, err
	}
	defer r.Body.Close()

	return body, nil
}

func (s *Server) sendAck(w http.ResponseWriter, receivedEndpoint string, data any, code int) {
	body, err := json.Marshal(data)
	if err != nil {
		s.sendError(w, errors.MarshalError(receivedEndpoint, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	if _, err = w.Write(body); err != nil {
		s.sendError(w, errors.WriteBodyError(receivedEndpoint, err), http.StatusInternalServerError)
	}
}

func (s *Server) sendError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	s.log.Error(err)
}
