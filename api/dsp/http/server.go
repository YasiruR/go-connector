package http

import (
	"encoding/json"
	"github.com/YasiruR/connector/core"
	"github.com/YasiruR/connector/protocols/catalog"
	"github.com/gorilla/mux"
	"github.com/tryfix/log"
	"io"
	"net/http"
)

type Server struct {
	port   int
	ownr   core.DataOwner
	router *mux.Router
	log    log.Logger // todo interface

	catSvc catalog.Service
}

func NewServer(port int, log log.Logger) *Server {
	s := Server{port: port, log: log}
	r := mux.NewRouter()

	// catalog protocol related endpoints
	r.HandleFunc(catalog.RequestEndpoint, s.handleCatalogRequest).Methods(http.MethodPost)
	r.HandleFunc(catalog.RequestDatasetEndpoint, s.handleDatasetRequest).Methods(http.MethodGet)

	return &s
}

func (s *Server) handleCatalogRequest(w http.ResponseWriter, r *http.Request) {
	byts, err := io.ReadAll(r.Body)
	if err != nil {

	}

	var req catalog.Request
	if err = json.Unmarshal(byts, &req); err != nil {

	}

	s.catSvc.GetCatalog()
}

func (s *Server) handleDatasetRequest(w http.ResponseWriter, r *http.Request) {

}
