package http

import (
	"encoding/json"
	"github.com/YasiruR/connector/core/dsp"
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

const paramProviderPid = `providerPid`

type Server struct {
	port     int
	ownr     dsp.Owner
	router   *mux.Router
	provider dsp.Provider
	log      pkg.Log
}

func NewServer(port int, provider dsp.Provider, log pkg.Log) *Server {
	r := mux.NewRouter()
	s := Server{port: port, router: r, provider: provider, log: log}

	// negotiation protocol related endpoints
	r.HandleFunc(negotiation.NegotiationsEndpoint, s.GetNegotiation).Methods(http.MethodGet)
	r.HandleFunc(negotiation.RequestContractEndpoint, s.HandleContractRequest).Methods(http.MethodPost)

	return &s
}

func (s *Server) Start() {
	if err := http.ListenAndServe(":"+strconv.Itoa(s.port), s.router); err != nil {
		s.log.Fatal(errors.InitFailed(`DSP API`, err))
	}
}

func (s *Server) GetNegotiation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	providerPid, ok := params[paramProviderPid]
	if !ok {
		s.sendError(w, errors.PathParamNotFound(negotiation.NegotiationsEndpoint, paramProviderPid), http.StatusBadRequest)
		return
	}

	neg, err := s.provider.HandleNegotiationsRequest(providerPid)
	if err != nil {
		s.sendError(w, errors.HandlerFailed(negotiation.NegotiationsEndpoint, negotiation.TypeProviderHandler, err), http.StatusBadRequest)
	}

	s.sendAck(w, negotiation.NegotiationsEndpoint, neg, http.StatusOK)
}

func (s *Server) HandleContractRequest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.sendError(w, errors.InvalidRequestBody(negotiation.RequestContractEndpoint, err), http.StatusBadRequest)
		r.Body.Close()
		return
	}
	defer r.Body.Close()

	var req negotiation.ContractRequest
	if err = json.Unmarshal(body, &req); err != nil {
		s.sendError(w, errors.UnmarshalError(negotiation.RequestContractEndpoint, err), http.StatusBadRequest)
		return
	}

	negAck, err := s.provider.HandleContractRequest(req)
	if err != nil {
		s.sendError(w, errors.HandlerFailed(negotiation.RequestContractEndpoint, negotiation.TypeProviderHandler, err), http.StatusBadRequest)
		return
	}

	s.sendAck(w, negotiation.RequestContractEndpoint, negAck, http.StatusCreated)
}

func (s *Server) sendAck(w http.ResponseWriter, endpoint string, data any, code int) {
	body, err := json.Marshal(data)
	if err != nil {
		s.sendError(w, errors.MarshalError(endpoint, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	if _, err = w.Write(body); err != nil {
		s.sendError(w, errors.WriteBodyError(endpoint, err), http.StatusInternalServerError)
	}
}

func (s *Server) sendError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	s.log.Error(err)
}
