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

func (p *Provider) HandleCatalogRequest(_ any) (catalog.Response, error) {
	cat, err := p.catalog.Get()
	if err != nil {
		return catalog.Response{}, errors.StoreFailed(stores.TypeCatalog, `Get`, err)
	}

	return catalog.Response{
		Context:             dsp.Context,
		DspaceParticipantID: p.participantId,
		Catalog:             cat,
	}, nil
}

func (p *Provider) HandleDatasetRequest(id string) (catalog.DatasetResponse, error) {
	ds, err := p.catalog.Dataset(id)
	if err != nil {
		return catalog.DatasetResponse{}, errors.StoreFailed(stores.TypeCatalog, `Dataset`, err)
	}

	return catalog.DatasetResponse{
		Context: dsp.Context,
		Dataset: ds,
	}, nil
}

// Negotiation Protocol (reference: https://docs.internationaldataspaces.org/ids-knowledgebase/v/dataspace-protocol/contract-negotiation/contract.negotiation.protocol)

func (p *Provider) OfferContract() {}

func (p *Provider) AgreeContract(offerId, negotiationId string) (agreementId string, err error) {
	agreementId, err = p.urn.NewURN()
	if err != nil {
		return ``, errors.URNFailed(`providerPid`, `NewURN`, err)
	}

	offer, err := p.policyStore.Offer(offerId)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypePolicy, `Offer`, err)
	}

	assignee, err := p.negStore.Assignee(negotiationId)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypeContractNegotiation, `Assignee`, err)
	}

	cn, err := p.negStore.Negotiation(negotiationId)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
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
		CallbackAddr: p.callbackAddr,
	}

	url, err := p.negStore.CallbackAddr(negotiationId)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypeContractNegotiation, `CallbackAddr`, err)
	}

	data, err := json.Marshal(ca)
	if err != nil {
		return ``, errors.MarshalError(``, err)
	}

	res, err := p.client.Send(data, url+`/negotiations/`+cn.ConsPId+`/agreement`)
	if err != nil {
		return ``, errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack negotiation.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return ``, errors.UnmarshalError(``, err)
	}

	if err = p.negStore.UpdateState(negotiationId, negotiation.StateAgreed); err != nil {
		return ``, errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	p.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", negotiationId, negotiation.StateAgreed))
	return agreementId, nil
}

func (p *Provider) FinalizeContract() {}

func (p *Provider) HandleNegotiationsRequest(providerPid string) (negotiation.Ack, error) {
	ack, err := p.negStore.Negotiation(providerPid)
	if err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	return negotiation.Ack(ack), nil
}

func (p *Provider) HandleContractRequest(cr negotiation.ContractRequest) (ack negotiation.Ack, err error) {
	// return error message if offerId is invalid

	// return error message if callbackAddress is invalid

	// associate with existing contract negotiation if providerPid exists and create a new contract
	// negotiation if otherwise
	var cn negotiation.Negotiation
	provPId := cr.ProvPId
	if provPId != `` {
		cn, err = p.negStore.Negotiation(provPId)
		if err != nil {
			return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
		}

		if cn.State != negotiation.StateOffered {
			return negotiation.Ack{}, errors.IncompatibleValues(`state`, string(cn.State), string(negotiation.StateOffered))
		}

		if cn.ConsPId != cr.ConsPId {
			return negotiation.Ack{}, errors.IncompatibleValues(`consumerPid`, cn.ConsPId, cr.ConsPId)
		}

		cn.State = negotiation.StateRequested
		cn.Type = negotiation.TypeNegotiationAck
		p.log.Trace("a valid contract negotiation exists", cn.ProvPId)
	} else {
		provPId, err = p.urn.NewURN()
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
	}

	p.negStore.Set(provPId, cn)
	p.negStore.SetAssignee(provPId, cr.Offer.Assignee)
	p.negStore.SetCallbackAddr(provPId, cr.CallbackAddr)
	p.log.Trace(fmt.Sprintf("stored contract negotiation (assigner: %s, assignee: %s)", cr.Offer.Assigner, cr.Offer.Assignee), cn)
	p.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", provPId, negotiation.StateRequested))
	return negotiation.Ack(cn), nil
}

func (p *Provider) HandleAgreementVerification(providerPid string) (negotiation.Ack, error) {
	cn, err := p.negStore.Negotiation(providerPid)
	if err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	if err = p.negStore.UpdateState(providerPid, negotiation.StateVerified); err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	cn.State = negotiation.StateVerified
	cn.Type = negotiation.TypeNegotiationAck
	p.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", providerPid, negotiation.StateVerified))
	return negotiation.Ack(cn), nil
}
