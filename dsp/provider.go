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
	"time"
)

type Provider struct {
	participantId string // data space specific identifier for Provider
	callbackAddr  string
	negStore      stores.ContractNegotiation
	policyStore   stores.Policy
	catalog       stores.Catalog
	urn           pkg.URNService
	client        pkg.Client
	log           pkg.Log
}

func NewProvider(port int, stores core.Stores, plugins core.Plugins) dsp.Provider {
	plugins.Log.Info("enabled data provider functions")
	return &Provider{
		participantId: `participant-id-provider`,
		callbackAddr:  `http://localhost:` + strconv.Itoa(port),
		negStore:      stores.ContractNegotiation,
		policyStore:   stores.Policy,
		catalog:       stores.Catalog,
		urn:           plugins.URNService,
		client:        plugins.Client,
		log:           plugins.Log,
	}
}

// Catalog Protocol (reference: https://docs.internationaldataspaces.org/ids-knowledgebase/v/dataspace-protocol/catalog/catalog.protocol)

func (s *Provider) HandleCatalogRequest(_ any) (catalog.Response, error) {
	cat, err := s.catalog.Get()
	if err != nil {
		return catalog.Response{}, errors.StoreFailed(stores.TypeCatalog, `Get`, err)
	}

	return catalog.Response{
		Context:             dsp.Context,
		DspaceParticipantID: s.participantId,
		Catalog:             cat,
	}, nil
}

func (s *Provider) HandleDatasetRequest(id string) (catalog.DatasetResponse, error) {
	return catalog.DatasetResponse{}, nil
}

// Negotiation Protocol (reference: https://docs.internationaldataspaces.org/ids-knowledgebase/v/dataspace-protocol/contract-negotiation/contract.negotiation.protocol)

func (s *Provider) OfferContract() {}

func (s *Provider) AgreeContract(offerId, negotiationId string) (agreementId string, err error) {
	agreementId, err = s.urn.NewURN()
	if err != nil {
		return ``, errors.URNFailed(`providerPid`, `NewURN`, err)
	}

	offer, err := s.policyStore.Offer(offerId)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypePolicy, `Offer`, err)
	}

	assignee, err := s.negStore.Assignee(negotiationId)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypeContractNegotiation, `Assignee`, err)
	}

	cn, err := s.negStore.Get(negotiationId)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypeContractNegotiation, `Get`, err)
	}

	ca := negotiation.ContractAgreement{
		Ctx:     dsp.Context,
		Type:    negotiation.TypeContractAgreement,
		ProvPId: cn.ProvPId,
		ConsPId: cn.ConsPId,
		Agreement: odrl.Agreement{
			Id:        agreementId,
			Type:      odrl.TypeAgreement,
			Target:    offer.Target,
			Assigner:  offer.Assigner,
			Assignee:  assignee,
			Timestamp: time.Now().UTC().String(), // change format into XSD
		},
		CallbackAddr: s.callbackAddr,
	}

	url, err := s.negStore.CallbackAddr(negotiationId)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypeContractNegotiation, `CallbackAddr`, err)
	}

	data, err := json.Marshal(ca)
	if err != nil {
		return ``, errors.MarshalError(``, err)
	}

	res, err := s.client.Send(data, url+`/negotiations/`+cn.ConsPId+`/agreement`)
	if err != nil {
		return ``, errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack negotiation.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return ``, errors.UnmarshalError(``, err)
	}

	if err = s.negStore.UpdateState(negotiationId, negotiation.StateAgreed); err != nil {
		return ``, errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	s.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", negotiationId, negotiation.StateAgreed))
	return agreementId, nil
}

func (s *Provider) FinalizeContract() {}

func (s *Provider) HandleNegotiationsRequest(providerPid string) (negotiation.Ack, error) {
	ack, err := s.negStore.Get(providerPid)
	if err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `Get`, err)
	}

	return negotiation.Ack(ack), nil
}

func (s *Provider) HandleContractRequest(cr negotiation.ContractRequest) (ack negotiation.Ack, err error) {
	// return error message if offerId is invalid

	// return error message if callbackAddress is invalid

	// associate with existing contract negotiation if providerPid exists and create a new contract
	// negotiation if otherwise
	var cn negotiation.Negotiation
	provPId := cr.ProvPId
	if provPId != `` {
		cn, err = s.negStore.Get(provPId)
		if err != nil {
			return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `Get`, err)
		}

		if cn.State != negotiation.StateOffered {
			return negotiation.Ack{}, errors.IncompatibleValues(`state`, string(cn.State), string(negotiation.StateOffered))
		}

		if cn.ConsPId != cr.ConsPId {
			return negotiation.Ack{}, errors.IncompatibleValues(`consumerPid`, cn.ConsPId, cr.ConsPId)
		}

		cn.State = negotiation.StateRequested
		s.log.Trace("a valid contract negotiation already exists", cn)
	} else {
		provPId, err = s.urn.NewURN()
		if err != nil {
			return negotiation.Ack{}, errors.URNFailed(`providerPid`, `NewURN`, err)
		}

		cn = negotiation.Negotiation{
			Ctx:     dsp.Context,
			Type:    negotiation.TypeNegotiationAck,
			ConsPId: cr.ConsPId,
			ProvPId: provPId,
			State:   negotiation.StateRequested,
		}
		s.log.Trace("a new contract negotiation was created", cn)
	}

	s.negStore.Set(provPId, cn)
	s.negStore.SetAssignee(provPId, cr.Offer.Assignee)
	s.negStore.SetCallbackAddr(provPId, cr.CallbackAddr)
	s.log.Trace(fmt.Sprintf("stored contract negotiation (id: %s, assigner: %s, assignee: %s)", provPId, cr.Offer.Assigner, cr.Offer.Assignee))
	s.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", provPId, negotiation.StateRequested))
	return negotiation.Ack(cn), nil
}
