package provider

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/connector/core/dsp"
	"github.com/YasiruR/connector/core/dsp/negotiation"
	"github.com/YasiruR/connector/core/errors"
	"github.com/YasiruR/connector/core/pkg"
	"github.com/YasiruR/connector/core/protocols/odrl"
	"github.com/YasiruR/connector/core/stores"
	"strconv"
	"time"
)

type Service struct {
	callbackAddr string
	cnStore      stores.ContractNegotiation
	polStore     stores.Policy
	urn          pkg.URN
	client       pkg.HTTPClient
	log          pkg.Log
}

func New(port int, cnStore stores.ContractNegotiation, polStore stores.Policy, urn pkg.URN, c pkg.HTTPClient, log pkg.Log) dsp.Provider {
	return &Service{
		callbackAddr: `http://localhost:` + strconv.Itoa(port),
		cnStore:      cnStore,
		polStore:     polStore,
		urn:          urn,
		client:       c,
		log:          log,
	}
}

func (s *Service) CreateAsset() {}

func (s *Service) CreatePolicy() {}

func (s *Service) CreateContractDef() {}

func (s *Service) OfferContract() {}

func (s *Service) AgreeContract(offerId, negotiationId string) (agreementId string, err error) {
	agreementId, err = s.urn.New()
	if err != nil {
		return ``, errors.URNFailed(`providerPid`, `New`, err)
	}

	offer, err := s.polStore.Offer(offerId)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypePolicy, `Offer`, err)
	}

	assignee, err := s.cnStore.Assignee(negotiationId)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypeContractNegotiation, `Assignee`, err)
	}

	cn, err := s.cnStore.Get(negotiationId)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypeContractNegotiation, `Get`, err)
	}

	ca := negotiation.ContractAgreement{
		Ctx:     negotiation.Context,
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

	url, err := s.cnStore.CallbackAddr(negotiationId)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypeContractNegotiation, `CallbackAddr`, err)
	}

	data, err := json.Marshal(ca)
	if err != nil {
		return ``, errors.MarshalError(``, err)
	}

	res, statusCode, err := s.client.Post(url+`/negotiations/`+cn.ConsPId+`/agreement`, data)
	if err != nil {
		return ``, errors.PkgFailed(pkg.TypeHTTPClient, `Post`, err)
	}

	var ack negotiation.Ack
	switch statusCode {
	case 200:
		if err = json.Unmarshal(res, &ack); err != nil {
			return ``, errors.UnmarshalError(``, err)
		}

		if err = s.cnStore.UpdateState(negotiationId, negotiation.StateAgreed); err != nil {
			return ``, errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
		}
	default:
		return ``, errors.InvalidStatusCode(statusCode, 200)
	}

	s.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", negotiationId, negotiation.StateAgreed))
	return agreementId, nil
}

func (s *Service) FinalizeContract() {}

func (s *Service) HandleNegotiationsRequest(providerPid string) (negotiation.Ack, error) {
	ack, err := s.cnStore.Get(providerPid)
	if err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `Get`, err)
	}

	return negotiation.Ack(ack), nil
}

func (s *Service) HandleContractRequest(cr negotiation.ContractRequest) (ack negotiation.Ack, err error) {
	// return error message if offerId is invalid

	// return error message if callbackAddress is invalid

	// associate with existing contract negotiation if providerPid exists and create a new contract
	// negotiation if otherwise
	var cn negotiation.Negotiation
	provPId := cr.ProvPId
	if provPId != `` {
		cn, err = s.cnStore.Get(provPId)
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
		provPId, err = s.urn.New()
		if err != nil {
			return negotiation.Ack{}, errors.URNFailed(`providerPid`, `New`, err)
		}

		cn = negotiation.Negotiation{
			Ctx:     negotiation.Context,
			Type:    negotiation.TypeNegotiationAck,
			ConsPId: cr.ConsPId,
			ProvPId: provPId,
			State:   negotiation.StateRequested,
		}
		s.log.Trace("a new contract negotiation was created", cn)
	}

	s.cnStore.Set(provPId, cn)
	s.cnStore.SetAssignee(provPId, cr.Offer.Assignee)
	s.cnStore.SetCallbackAddr(provPId, cr.CallbackAddr)
	s.log.Info(fmt.Sprintf("stored contract negotiation (id: %s, assigner: %s, assignee: %s)", provPId, cr.Offer.Assigner, cr.Offer.Assignee))
	return negotiation.Ack(cn), nil
}
