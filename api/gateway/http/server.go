package http

import (
	"encoding/json"
	"github.com/YasiruR/connector/core/api/gateway"
	"github.com/YasiruR/connector/core/dsp"
	"github.com/YasiruR/connector/core/errors"
	"github.com/YasiruR/connector/core/pkg"
	"github.com/YasiruR/connector/core/protocols/odrl"
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
	provider dsp.Provider
	consumer dsp.Consumer
	owner    dsp.Owner
	log      pkg.Log
}

func NewServer(port int, p dsp.Provider, c dsp.Consumer, o dsp.Owner, log pkg.Log) *Server {
	r := mux.NewRouter()
	s := Server{port: port, router: r, provider: p, consumer: c, owner: o, log: log}

	// endpoints related to data assets
	r.HandleFunc(gateway.CreatePolicyEndpoint, s.CreatePolicy).Methods(http.MethodPost)
	r.HandleFunc(gateway.CreateDatasetEndpoint, s.CreateDataset).Methods(http.MethodPost)

	// endpoints related to negotiation
	r.HandleFunc(gateway.ContractRequestEndpoint, s.RequestContract).Methods(http.MethodPost)
	r.HandleFunc(gateway.ContractAgreementEndpoint, s.AgreeContract).Methods(http.MethodPost)

	return &s
}

func (s *Server) Start() {
	if err := http.ListenAndServe(":"+strconv.Itoa(s.port), s.router); err != nil {
		s.log.Fatal(errors.InitFailed(`gateway API`, err))
	}
}

func (s *Server) CreatePolicy(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.sendError(w, errors.InvalidRequestBody(gateway.CreatePolicyEndpoint, err), http.StatusBadRequest)
		r.Body.Close()
		return
	}
	defer r.Body.Close()

	var req gateway.PolicyRequest
	if err = json.Unmarshal(body, &req); err != nil {
		s.sendError(w, errors.UnmarshalError(gateway.CreatePolicyEndpoint, err), http.StatusBadRequest)
		return
	}

	// todo remove odrl bindings from this func
	var perms []odrl.Rule // handle other policy types
	for _, p := range req.Permissions {
		var cons []odrl.Constraint
		for _, c := range p.Constraints {
			cons = append(cons, odrl.Constraint{
				LeftOperand:  c.LeftOperand,
				Operator:     c.Operator,
				RightOperand: c.RightOperand,
			})
		}
		perms = append(perms, odrl.Rule{Action: odrl.Action(p.Action), Constraints: cons})
	}

	// todo check if target is required here
	id, err := s.owner.CreatePolicy(`test`, perms, []odrl.Rule{})
	if err != nil {
		s.sendError(w, errors.DSPFailed(dsp.RoleOwner, `CreatePolicy`, err), http.StatusBadRequest)
		return
	}

	s.sendAck(w, gateway.CreatePolicyEndpoint, gateway.PolicyResponse{Id: id}, http.StatusOK)
}

func (s *Server) CreateDataset(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.sendError(w, errors.InvalidRequestBody(gateway.CreateDatasetEndpoint, err), http.StatusBadRequest)
		r.Body.Close()
		return
	}
	defer r.Body.Close()

	var req gateway.DatasetRequest
	if err = json.Unmarshal(body, &req); err != nil {
		s.sendError(w, errors.UnmarshalError(gateway.CreateDatasetEndpoint, err), http.StatusBadRequest)
	}

	id, err := s.owner.CreateDataset(req.Title, req.Descriptions, req.Keywords, req.Endpoints, req.PolicyIds)
	if err != nil {
		s.sendError(w, errors.DSPFailed(dsp.RoleOwner, `CreateDataset`, err), http.StatusBadRequest)
	}

	s.sendAck(w, gateway.CreateDatasetEndpoint, gateway.DatasetResponse{Id: id}, http.StatusOK)
}

func (s *Server) RequestContract(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.sendError(w, errors.InvalidRequestBody(gateway.ContractRequestEndpoint, err), http.StatusBadRequest)
		r.Body.Close()
		return
	}
	defer r.Body.Close()

	var req gateway.ContractRequest
	if err = json.Unmarshal(body, &req); err != nil {
		s.sendError(w, errors.UnmarshalError(gateway.ContractRequestEndpoint, err), http.StatusBadRequest)
		return
	}

	negId, err := s.consumer.RequestContract(req.OfferId, req.ProviderEndpoint, req.ProviderPId, req.OdrlTarget,
		req.Assigner, req.Assignee, req.Action)
	if err != nil {
		s.sendError(w, errors.DSPFailed(dsp.RoleConsumer, `RequestContract`, err), http.StatusBadRequest)
		return
	}

	s.sendAck(w, gateway.ContractRequestEndpoint, gateway.ContractRequestResponse{Id: negId}, http.StatusOK)
}

func (s *Server) AgreeContract(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.sendError(w, errors.InvalidRequestBody(gateway.ContractAgreementEndpoint, err), http.StatusBadRequest)
		r.Body.Close()
		return
	}
	defer r.Body.Close()

	var req gateway.ContractAgreementRequest
	if err = json.Unmarshal(body, &req); err != nil {
		s.sendError(w, errors.UnmarshalError(gateway.ContractAgreementEndpoint, err), http.StatusBadRequest)
		return
	}

	agrId, err := s.provider.AgreeContract(req.OfferId, req.NegotiationId)
	if err != nil {
		s.sendError(w, errors.DSPFailed(dsp.RoleProvider, `AgreeContract`, err), http.StatusBadRequest)
		return
	}

	s.sendAck(w, gateway.ContractAgreementEndpoint, gateway.ContractAgreementResponse{Id: agrId}, http.StatusOK)
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

// todo remove
func (s *Server) sendError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	s.log.Error(err)
}
