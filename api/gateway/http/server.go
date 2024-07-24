package http

import (
	"encoding/json"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/gateway"
	"github.com/YasiruR/connector/domain/dsp"
	"github.com/YasiruR/connector/domain/dsp/negotiation"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/protocols/odrl"
	"github.com/YasiruR/connector/domain/stores"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

// gateway.http.Server contains the endpoints which will be used by a client to initiate
// message flows or manage both control and data planes

// todo check return error codes

type Server struct {
	port     int
	router   *mux.Router
	provider dsp.Provider
	consumer dsp.Consumer
	owner    dsp.Owner
	agrStore stores.Agreement // should this be invoked through a role?
	log      pkg.Log
}

func NewServer(port int, roles domain.Roles, stores domain.Stores, log pkg.Log) *Server {
	r := mux.NewRouter()
	s := Server{
		port:     port,
		router:   r,
		provider: roles.Provider,
		consumer: roles.Consumer,
		owner:    roles.Owner,
		agrStore: stores.Agreement,
		log:      log,
	}

	// endpoints related to catalog
	r.HandleFunc(gateway.CreatePolicyEndpoint, s.CreatePolicy).Methods(http.MethodPost)
	r.HandleFunc(gateway.CreateDatasetEndpoint, s.CreateDataset).Methods(http.MethodPost)
	r.HandleFunc(gateway.RequestCatalogEndpoint, s.RequestCatalog).Methods(http.MethodPost)
	r.HandleFunc(gateway.RequestDatasetEndpoint, s.RequestDataset).Methods(http.MethodPost)

	// endpoints related to negotiation
	r.HandleFunc(gateway.RequestContractEndpoint, s.RequestContract).Methods(http.MethodPost)
	r.HandleFunc(gateway.AgreeContractEndpoint, s.AgreeContract).Methods(http.MethodPost)
	r.HandleFunc(gateway.GetAgreementEndpoint, s.GetAgreement).Methods(http.MethodGet)
	r.HandleFunc(gateway.VerifyAgreementEndpoint, s.VerifyAgreement).Methods(http.MethodPost)

	return &s
}

func (s *Server) Start() {
	s.log.Info("gateway HTTP server is listening on " + strconv.Itoa(s.port))
	if err := http.ListenAndServe(":"+strconv.Itoa(s.port), s.router); err != nil {
		s.log.Fatal(errors.InitFailed(`gateway API`, err))
	}
}

func (s *Server) CreatePolicy(w http.ResponseWriter, r *http.Request) {
	body, err := s.readBody(gateway.CreatePolicyEndpoint, w, r)
	if err != nil {
		return
	}

	var req gateway.CreatePolicyRequest
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
	body, err := s.readBody(gateway.CreateDatasetEndpoint, w, r)
	if err != nil {
		return
	}

	var req gateway.CreateDatasetRequest
	if err = json.Unmarshal(body, &req); err != nil {
		s.sendError(w, errors.UnmarshalError(gateway.CreateDatasetEndpoint, err), http.StatusBadRequest)
	}

	id, err := s.owner.CreateDataset(req.Title, req.Descriptions, req.Keywords, req.Endpoints, req.OfferIds)
	if err != nil {
		s.sendError(w, errors.DSPFailed(dsp.RoleOwner, `CreateDataset`, err), http.StatusBadRequest)
	}

	s.sendAck(w, gateway.CreateDatasetEndpoint, gateway.DatasetResponse{Id: id}, http.StatusOK)
}

func (s *Server) RequestCatalog(w http.ResponseWriter, r *http.Request) {
	body, err := s.readBody(gateway.RequestCatalogEndpoint, w, r)
	if err != nil {
		return
	}

	var req gateway.CatalogRequest
	if err = json.Unmarshal(body, &req); err != nil {
		s.sendError(w, errors.UnmarshalError(gateway.RequestCatalogEndpoint, err), http.StatusBadRequest)
	}

	cat, err := s.consumer.RequestCatalog(req.ProviderEndpoint)
	if err != nil {
		s.sendError(w, errors.DSPFailed(dsp.RoleConsumer, `RequestCatalog`, err), http.StatusBadRequest)
	}

	s.sendAck(w, gateway.RequestCatalogEndpoint, cat, http.StatusOK)
}

func (s *Server) RequestDataset(w http.ResponseWriter, r *http.Request) {
	body, err := s.readBody(gateway.RequestDatasetEndpoint, w, r)
	if err != nil {
		return
	}

	var req gateway.DatasetRequest
	if err = json.Unmarshal(body, &req); err != nil {
		s.sendError(w, errors.UnmarshalError(gateway.RequestDatasetEndpoint, err), http.StatusBadRequest)
	}

	ds, err := s.consumer.RequestDataset(req.DatasetId, req.ProviderEndpoint)
	if err != nil {
		s.sendError(w, errors.DSPFailed(dsp.RoleConsumer, `RequestDataset`, err), http.StatusBadRequest)
	}

	s.sendAck(w, gateway.RequestDatasetEndpoint, ds, http.StatusOK)
}

func (s *Server) RequestContract(w http.ResponseWriter, r *http.Request) {
	body, err := s.readBody(gateway.RequestContractEndpoint, w, r)
	if err != nil {
		return
	}

	var req gateway.ContractRequest
	if err = json.Unmarshal(body, &req); err != nil {
		s.sendError(w, errors.UnmarshalError(gateway.RequestContractEndpoint, err), http.StatusBadRequest)
		return
	}

	// todo check if providerPid should really be requested
	negId, err := s.consumer.RequestContract(req.OfferId, req.ProviderEndpoint, req.ProviderPId, req.OdrlTarget,
		req.Assigner, req.Assignee, req.Action)
	if err != nil {
		s.sendError(w, errors.DSPFailed(dsp.RoleConsumer, `RequestContract`, err), http.StatusBadRequest)
		return
	}

	s.sendAck(w, gateway.RequestContractEndpoint, gateway.ContractRequestResponse{Id: negId}, http.StatusOK)
}

func (s *Server) AgreeContract(w http.ResponseWriter, r *http.Request) {
	body, err := s.readBody(gateway.AgreeContractEndpoint, w, r)
	if err != nil {
		return
	}

	var req gateway.AgreeContractRequest
	if err = json.Unmarshal(body, &req); err != nil {
		s.sendError(w, errors.UnmarshalError(gateway.AgreeContractEndpoint, err), http.StatusBadRequest)
		return
	}

	agrId, err := s.provider.AgreeContract(req.OfferId, req.NegotiationId)
	if err != nil {
		s.sendError(w, errors.DSPFailed(dsp.RoleProvider, `AgreeContract`, err), http.StatusBadRequest)
		return
	}

	s.sendAck(w, gateway.AgreeContractEndpoint, gateway.ContractAgreementResponse{Id: agrId}, http.StatusOK)
}

func (s *Server) GetAgreement(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	agreementId, ok := params[gateway.ParamId]
	if !ok {
		s.sendError(w, errors.PathParamNotFound(gateway.GetAgreementEndpoint, negotiation.ParamConsumerPid), http.StatusBadRequest)
		return
	}

	agr, err := s.agrStore.Get(agreementId)
	if err != nil {
		s.sendError(w, errors.StoreFailed(stores.TypeAgreement, `Get`, err), http.StatusBadRequest)
		return
	}

	s.sendAck(w, gateway.GetAgreementEndpoint, agr, http.StatusOK)
}

func (s *Server) VerifyAgreement(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	consumerPid, ok := params[gateway.ParamConsumerPid]
	if !ok {
		s.sendError(w, errors.PathParamNotFound(gateway.VerifyAgreementEndpoint, negotiation.ParamConsumerPid), http.StatusBadRequest)
		return
	}

	if err := s.consumer.VerifyAgreement(consumerPid); err != nil {
		s.sendError(w, errors.DSPFailed(dsp.RoleConsumer, `VerifyAgreement`, err), http.StatusBadRequest)
		return
	}

	s.sendAck(w, gateway.VerifyAgreementEndpoint, nil, http.StatusOK)
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

// todo remove
func (s *Server) sendError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	s.log.Error(errors.APIFailed(`gateway`, err))
}
