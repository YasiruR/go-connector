package negotiation

import (
	"encoding/json"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/gateway"
	"github.com/YasiruR/connector/domain/dsp"
	"github.com/YasiruR/connector/domain/dsp/negotiation"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/models/odrl"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/stores"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type Handler struct {
	provider dsp.Provider
	consumer dsp.Consumer
	agrStore stores.Agreement
	log      pkg.Log
}

func NewHandler(roles domain.Roles, stores domain.Stores, log pkg.Log) *Handler {
	return &Handler{
		provider: roles.Provider,
		consumer: roles.Consumer,
		agrStore: stores.Agreement,
		log:      log,
	}
}

func (n *Handler) RequestContract(w http.ResponseWriter, r *http.Request) {
	body, err := n.readBody(gateway.RequestContractEndpoint, w, r)
	if err != nil {
		return
	}

	var req gateway.ContractRequest
	if err = json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		n.log.Error(errors.UnmarshalError(gateway.RequestContractEndpoint, err))
		return
	}

	ofr := odrl.Offer{
		Id:          req.OfferId,
		Type:        odrl.TypeOffer,
		Target:      odrl.Target(req.OdrlTarget),
		Assigner:    odrl.Assigner(req.Assigner),
		Assignee:    odrl.Assignee(req.Assignee),
		Permissions: []odrl.Rule{{Action: odrl.Action(req.Action)}}, // should handle constraints
	}

	negId, err := n.consumer.RequestContract(req.ProviderEndpoint, ofr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		n.log.Error(errors.DSPFailed(dsp.RoleConsumer, `RequestContract`, err))
		return
	}

	n.sendAck(w, gateway.RequestContractEndpoint, gateway.ContractRequestResponse{Id: negId}, http.StatusOK)
}

func (n *Handler) AgreeContract(w http.ResponseWriter, r *http.Request) {
	body, err := n.readBody(gateway.AgreeContractEndpoint, w, r)
	if err != nil {
		return
	}

	var req gateway.AgreeContractRequest
	if err = json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		n.log.Error(errors.UnmarshalError(gateway.AgreeContractEndpoint, err))
		return
	}

	agrId, err := n.provider.AgreeContract(req.OfferId, req.NegotiationId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		n.log.Error(errors.DSPFailed(dsp.RoleProvider, `AgreeContract`, err))
		return
	}

	n.sendAck(w, gateway.AgreeContractEndpoint, gateway.ContractAgreementResponse{Id: agrId}, http.StatusOK)
}

func (n *Handler) GetAgreement(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	agreementId, ok := params[gateway.ParamAgreementId]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		n.log.Error(errors.PathParamNotFound(gateway.GetAgreementEndpoint, negotiation.ParamConsumerPid))
		return
	}

	agr, err := n.agrStore.Get(agreementId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		n.log.Error(errors.StoreFailed(stores.TypeAgreement, `Get`, err))
		return
	}

	n.sendAck(w, gateway.GetAgreementEndpoint, agr, http.StatusOK)
}

func (n *Handler) VerifyAgreement(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	consumerPid, ok := params[gateway.ParamConsumerPid]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		n.log.Error(errors.PathParamNotFound(gateway.VerifyAgreementEndpoint, negotiation.ParamConsumerPid))
		return
	}

	if err := n.consumer.VerifyAgreement(consumerPid); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		n.log.Error(errors.DSPFailed(dsp.RoleConsumer, `VerifyAgreement`, err))
		return
	}

	n.sendAck(w, gateway.VerifyAgreementEndpoint, nil, http.StatusOK)
}

func (n *Handler) FinalizeContract(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	providerPid, ok := params[gateway.ParamProviderPid]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		n.log.Error(errors.PathParamNotFound(gateway.FinalizeContractEndpoint, negotiation.ParamProviderId))
		return
	}

	if err := n.provider.FinalizeContract(providerPid); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		n.log.Error(errors.DSPFailed(dsp.RoleProvider, `FinalizeContract`, err))
		return
	}

	n.sendAck(w, gateway.FinalizeContractEndpoint, nil, http.StatusOK)
}

func (n *Handler) readBody(endpoint string, w http.ResponseWriter, r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		err = errors.InvalidRequestBody(endpoint, err)
		w.WriteHeader(http.StatusBadRequest)
		r.Body.Close()
		return nil, err
	}
	defer r.Body.Close()

	return body, nil
}

func (n *Handler) sendAck(w http.ResponseWriter, receivedEndpoint string, data any, code int) {
	body, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		n.log.Error(errors.MarshalError(receivedEndpoint, err))
		return
	}

	w.WriteHeader(code)
	if _, err = w.Write(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		n.log.Error(errors.WriteBodyError(receivedEndpoint, err))
	}
}
