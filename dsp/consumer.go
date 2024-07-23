package dsp

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/connector/core"
	"github.com/YasiruR/connector/core/dsp"
	"github.com/YasiruR/connector/core/dsp/catalog"
	"github.com/YasiruR/connector/core/dsp/negotiation"
	"github.com/YasiruR/connector/core/errors"
	"github.com/YasiruR/connector/core/pkg"
	"github.com/YasiruR/connector/core/protocols/odrl"
	"github.com/YasiruR/connector/core/stores"
	"strconv"
)

type Consumer struct {
	callbackAddr string
	negStore     stores.ContractNegotiation
	urn          pkg.URNService
	client       pkg.Client
	log          pkg.Log
}

func NewConsumer(port int, stores core.Stores, plugins core.Plugins) dsp.Consumer {
	plugins.Log.Info("enabled data consumer functions")
	return &Consumer{
		callbackAddr: `http://localhost:` + strconv.Itoa(port),
		negStore:     stores.ContractNegotiation,
		urn:          plugins.URNService,
		client:       plugins.Client,
		log:          plugins.Log,
	}
}

func (s *Consumer) RequestCatalog(endpoint string) (catalog.Response, error) {
	req := catalog.Request{
		Context:      dsp.Context,
		Type:         catalog.TypeCatalogRequest,
		DspaceFilter: nil,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return catalog.Response{}, errors.MarshalError(endpoint, err)
	}

	res, err := s.client.Send(data, endpoint+catalog.RequestEndpoint)
	if err != nil {
		return catalog.Response{}, errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var cat catalog.Response
	if err = json.Unmarshal(res, &cat); err != nil {
		return catalog.Response{}, errors.UnmarshalError(``, err)
	}

	return cat, nil
}

func (s *Consumer) RequestDataset(id, endpoint string) (catalog.DatasetResponse, error) {
	req := catalog.DatasetRequest{
		Context:   dsp.Context,
		Type:      catalog.TypeDatasetRequest,
		DatasetId: id,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return catalog.DatasetResponse{}, errors.MarshalError(endpoint, err)
	}

	res, err := s.client.Send(data, endpoint+catalog.RequestDatasetEndpoint)
	if err != nil {
		return catalog.DatasetResponse{}, errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var dataset catalog.DatasetResponse
	if err = json.Unmarshal(res, &dataset); err != nil {
		return catalog.DatasetResponse{}, errors.UnmarshalError(``, err)
	}

	return dataset, nil
}

func (s *Consumer) RequestContract(offerId, providerEndpoint, providerPid, target, assigner, assignee, action string) (negotiationId string, err error) {
	// generate consumerPid
	consPId, err := s.urn.NewURN()
	if err != nil {
		return ``, errors.URNFailed(`consumerPid`, `NewURN`, err)
	}

	// construct payload
	req := negotiation.ContractRequest{
		ProvPId: providerPid,
		ConsPId: consPId,
		Offer: odrl.Offer{
			Id:          offerId,
			Target:      odrl.Target(target),
			Assigner:    odrl.Assigner(assigner),
			Assignee:    odrl.Assignee(assignee),
			Permissions: []odrl.Rule{{Action: odrl.Action(action)}}, // should handle constraints
		},
		CallbackAddr: s.callbackAddr,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return ``, errors.MarshalError(``, err)
	}

	res, err := s.client.Send(data, providerEndpoint+negotiation.ContractRequestEndpoint)
	if err != nil {
		return ``, errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack negotiation.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return ``, errors.UnmarshalError(``, err)
	}

	s.negStore.Set(consPId, negotiation.Negotiation(ack))
	s.negStore.SetAssignee(consPId, odrl.Assignee(assignee))
	s.log.Trace(fmt.Sprintf("stored contract negotiation (id: %s, assigner: %s, assignee: %s)", consPId, assigner, assignee))
	s.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", consPId, negotiation.StateRequested))
	return consPId, nil
}

func (s *Consumer) AcceptContract() {}

func (s *Consumer) VerifyAgreement() {}

func (s *Consumer) TerminateContract() {}

func (s *Consumer) HandleContractAgreement(consumerPid string, ca negotiation.ContractAgreement) (negotiation.Ack, error) {
	// validate agreement (e.g. consumerPid, target)

	if err := s.negStore.UpdateState(ca.ConsPId, negotiation.StateAgreed); err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	s.negStore.SetCallbackAddr(ca.ConsPId, s.callbackAddr)
	neg, err := s.negStore.Get(ca.ConsPId)
	if err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `Get`, err)
	}

	s.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", ca.ConsPId, negotiation.StateAgreed))
	return negotiation.Ack(neg), nil
}
