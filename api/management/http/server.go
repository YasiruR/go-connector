package http

import (
	"github.com/YasiruR/connector/core"
	"github.com/gorilla/mux"
	"github.com/tryfix/log"
	"net/http"
	"strconv"
)

type Server struct {
	port   int
	ownr   core.DataOwner
	router *mux.Router
	log    log.Logger
}

func NewServer(port int, log log.Logger) *Server {
	s := Server{port: port, log: log}
	r := mux.NewRouter()
	r.HandleFunc(createAssetEndpoint, s.handleCreateAsset).Methods(http.MethodPost)

	return &s
}

func (s *Server) Start() {
	if err := http.ListenAndServe(":"+strconv.Itoa(s.port), s.router); err != nil {
		s.log.Fatal(`http server initialization failed - %v`, err)
	}
}

func (s *Server) handleCreateAsset(w http.ResponseWriter, r *http.Request) {

}
