package http

import (
	httpPostgresql "github.com/YasiruR/go-connector/api/exchanger/http/postgresql"
	"github.com/YasiruR/go-connector/domain/api/exchanger/http/postgresql"
	"github.com/YasiruR/go-connector/domain/data-plane"
	"github.com/YasiruR/go-connector/domain/errors"
	"github.com/YasiruR/go-connector/domain/pkg"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"strconv"
)

type Server struct {
	port   int
	router *mux.Router
	psql   postgresql.Handler
	log    pkg.Log
}

func NewServer(port int, e data_plane.Exchanger, log pkg.Log) *Server {
	r := mux.NewRouter()
	s := Server{
		port:   port,
		router: r,
		psql:   httpPostgresql.NewHandler(e),
		log:    log,
	}

	r.HandleFunc(postgresql.PullEndpoint, s.psql.HandlePull).Methods(http.MethodPost)
	r.HandleFunc(postgresql.PushEndpoint, s.psql.HandlePush).Methods(http.MethodPost)

	return &s
}

func (s *Server) Start() {
	s.log.Info("exchanger HTTP server is listening on " + strconv.Itoa(s.port))
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut,
			http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowCredentials: true,
	})

	if err := http.ListenAndServe(":"+strconv.Itoa(s.port), c.Handler(s.router)); err != nil {
		s.log.Fatal(errors.ModuleInitFailed(`exchanger API`, err))
	}
}
