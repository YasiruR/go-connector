package http

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/connector/core/api/gateway"
	"github.com/YasiruR/connector/core/dsp"
	"github.com/YasiruR/connector/core/models/odrl"
	"github.com/YasiruR/connector/core/pkg"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

// gateway.http.Server contains the endpoints which will be used by a client to initiate
// message flows or manage both control and data planes

type Server struct {
	port     int
	router   *mux.Router
	consumer dsp.Consumer
	log      pkg.Log
}

func NewServer(port int, consumer dsp.Consumer, log pkg.Log) *Server {
	r := mux.NewRouter()
	s := Server{port: port, router: r, consumer: consumer, log: log}

	// endpoints related to asset
	r.HandleFunc(gateway.CreateAssetEndpoint, s.CreateAsset).Methods(http.MethodPost)

	// endpoints related to negotiation
	r.HandleFunc(gateway.ContractRequestEndpoint, s.RequestContract).Methods(http.MethodPost)

	return &s
}

func (s *Server) Start() {
	if err := http.ListenAndServe(":"+strconv.Itoa(s.port), s.router); err != nil {
		s.log.Fatal(`http server (management API) initialization failed - %v`, err)
	}
}

func (s *Server) CreateAsset(w http.ResponseWriter, r *http.Request) {
	// check if authorized to create an asset
}

func (s *Server) RequestContract(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.error(w, fmt.Sprintf("reading request body failed in initializing contract request - %s", err), http.StatusBadRequest)
		r.Body.Close()
		return
	}
	defer r.Body.Close()

	var req gateway.ContractRequest
	if err = json.Unmarshal(body, &req); err != nil {
		s.error(w, fmt.Sprintf("unmarshalling failed in initializing contract request - %s", err), http.StatusBadRequest)
		return
	}

	ot := odrl.Target(req.OdrlTarget)
	a := odrl.Assigner(req.Assigner)
	act := odrl.Action(req.Action)

	if err = s.consumer.RequestContract(req.OfferId, req.ProviderEndpoint, req.ProviderPId, ot, a, act); err != nil {
		s.error(w, fmt.Sprintf("consumer failed to send contract request in initializing contract request - %s", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) AgreeContract(w http.ResponseWriter, r *http.Request) {
	//body, err := io.ReadAll(r.Body)
	//if err != nil {
	//	s.error(w, fmt.Sprintf("reading request body failed in agreeing contract request - %s", err), http.StatusBadRequest)
	//	r.Body.Close()
	//	return
	//}
	//defer r.Body.Close()

}

func (s *Server) error(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	s.log.Error(message)
}
