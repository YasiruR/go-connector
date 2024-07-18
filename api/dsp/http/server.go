package http

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/connector/core"
	"github.com/YasiruR/connector/protocols/catalog"
	"github.com/YasiruR/connector/protocols/negotiation"
	"github.com/gorilla/mux"
	"github.com/tryfix/log"
	"io"
	"net/http"
	"strconv"
)

// dsp.http.Server contains the endpoints defined in data space protocols which will be used
// for the communication between connectors

type Server struct {
	port     int
	ownr     core.Owner
	router   *mux.Router
	provider core.Provider
	log      log.Logger // todo interface
}

func NewServer(port int, provider core.Provider, log log.Logger) *Server {
	r := mux.NewRouter()
	s := Server{port: port, router: r, provider: provider, log: log}

	// catalog protocol related endpoints
	r.HandleFunc(catalog.RequestEndpoint, s.handleCatalogRequest).Methods(http.MethodPost)
	r.HandleFunc(catalog.RequestDatasetEndpoint, s.handleDatasetRequest).Methods(http.MethodGet)

	// negotiation protocol related endpoints
	r.HandleFunc(negotiation.RequestContractEndpoint, s.handleContractRequest).Methods(http.MethodPost)

	return &s
}

func (s *Server) Start() {
	if err := http.ListenAndServe(":"+strconv.Itoa(s.port), s.router); err != nil {
		s.log.Fatal(`http server (DSP API) initialization failed - %v`, err)
	}
}

func (s *Server) handleCatalogRequest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
	}
	defer r.Body.Close()

	var req catalog.Request
	if err = json.Unmarshal(body, &req); err != nil {
	}
}

func (s *Server) handleDatasetRequest(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) handleContractRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("received contract requesttt veeeeee")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.log.Error(fmt.Sprintf("reading request body failed in handling contract request - %s", err))
		w.WriteHeader(http.StatusBadRequest)
		r.Body.Close()
		return
	}
	defer r.Body.Close()

	var req negotiation.ContractRequest
	if err = json.Unmarshal(body, &req); err != nil {
		s.log.Error(fmt.Sprintf("unmarshalling failed in handling contract request - %s", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = s.provider.HandleContractRequest(req); err != nil {
		s.log.Error(fmt.Sprintf("provider failed to handle contract request in handling contract request - %s", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.sendAck(w, negotiation.StateRequested, http.StatusCreated)
}

func (s *Server) sendAck(w http.ResponseWriter, st negotiation.State, code int) {
	res := negotiation.Ack{
		State: st,
	}

	body, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.log.Error(fmt.Errorf("failed to send ack - %w", err))
	}

	w.WriteHeader(code)
	if _, err = w.Write(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.log.Error(fmt.Errorf("failed when writing the ack response body - %w", err))
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) sendError() {

}
