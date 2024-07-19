package http

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/connector/core/dsp"
	negotiation2 "github.com/YasiruR/connector/core/dsp/negotiation"
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
	ownr     dsp.Owner
	router   *mux.Router
	provider dsp.Provider
	log      pkg.Log
}

func NewServer(port int, provider dsp.Provider, log pkg.Log) *Server {
	r := mux.NewRouter()
	s := Server{port: port, router: r, provider: provider, log: log}

	// negotiation protocol related endpoints
	r.HandleFunc(negotiation2.NegotiationsEndpoint, s.GetNegotiation).Methods(http.MethodGet)
	r.HandleFunc(negotiation2.RequestContractEndpoint, s.HandleContractRequest).Methods(http.MethodPost)

	return &s
}

func (s *Server) Start() {
	if err := http.ListenAndServe(":"+strconv.Itoa(s.port), s.router); err != nil {
		s.log.Fatal(`http server (DSP API) initialization failed - %v`, err)
	}
}

func (s *Server) GetNegotiation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	providerPid, ok := params["providerPid"]
	if !ok {
		s.sendError(w, "no providerPid found in negotiation request", http.StatusBadRequest)
		return
	}

	neg, err := s.provider.HandleNegotiationsRequest(providerPid)
	if err != nil {
		s.sendError(w, fmt.Sprintf("handler failed in negotiation request - %s", err), http.StatusBadRequest)
	}

	s.sendAck(w, neg, http.StatusOK)
}

func (s *Server) HandleContractRequest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.sendError(w, fmt.Sprintf("reading request body failed in handling contract request - %s", err), http.StatusBadRequest)
		r.Body.Close()
		return
	}
	defer r.Body.Close()

	var req negotiation2.ContractRequest
	if err = json.Unmarshal(body, &req); err != nil {
		s.sendError(w, fmt.Sprintf("unmarshalling failed in handling contract request - %s", err), http.StatusBadRequest)
		return
	}

	negAck, err := s.provider.HandleContractRequest(req)
	if err != nil {
		s.sendError(w, fmt.Sprintf("provider failed to handle contract request in handling contract request - %s", err), http.StatusBadRequest)
		return
	}

	s.sendAck(w, negAck, http.StatusCreated)
}

func (s *Server) sendAck(w http.ResponseWriter, data any, code int) {
	body, err := json.Marshal(data)
	if err != nil {
		s.sendError(w, fmt.Sprintf("marshalling failed - %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	if _, err = w.Write(body); err != nil {
		s.sendError(w, fmt.Sprintf("writing failed - %s", err), http.StatusInternalServerError)
	}
}

func (s *Server) sendError(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	s.log.Error(message)
}
