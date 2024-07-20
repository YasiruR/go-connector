package provider

import (
	"github.com/YasiruR/connector/core/dsp"
	"github.com/YasiruR/connector/core/dsp/negotiation"
	"github.com/YasiruR/connector/core/errors"
	"github.com/YasiruR/connector/core/pkg"
	coreStores "github.com/YasiruR/connector/core/stores"
	"github.com/YasiruR/connector/stores"
	"strconv"
)

type Service struct {
	callbackAddr string
	cnStore      *stores.ContractNegotiation
	urn          pkg.URN
	client       pkg.HTTPClient
	log          pkg.Log
}

func New(port int, cnStore *stores.ContractNegotiation, urn pkg.URN, c pkg.HTTPClient, log pkg.Log) dsp.Provider {
	return &Service{
		callbackAddr: `http://localhost:` + strconv.Itoa(port),
		cnStore:      cnStore,
		urn:          urn,
		client:       c,
		log:          log,
	}
}

func (s *Service) CreateAsset() {}

func (s *Service) CreatePolicy() {}

func (s *Service) CreateContractDef() {}

func (s *Service) OfferContract() {}

func (s *Service) AgreeContract() {

}

func (s *Service) FinalizeContract() {}

func (s *Service) HandleNegotiationsRequest(providerPid string) (negotiation.Ack, error) {
	ack, err := s.cnStore.Get(providerPid)
	if err != nil {
		return negotiation.Ack{}, errors.StoreFailed(coreStores.TypeContractNegotiation, `Get`, err)
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
			return negotiation.Ack{}, errors.StoreFailed(coreStores.TypeContractNegotiation, `Get`, err)
		}

		if cn.State != negotiation.StateOffered {
			return negotiation.Ack{}, errors.IncompatibleValues(`state`, string(cn.State), string(negotiation.StateOffered))
		}

		if cn.ConsPId != cr.ConsPId {
			return negotiation.Ack{}, errors.IncompatibleValues(`consumerPid`, cn.ConsPId, cr.ConsPId)
		}

		cn.State = negotiation.StateRequested
		s.log.Info("updated existing contract negotiation", cn)
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
		s.log.Info("stored new contract negotiation", cn)
	}

	s.cnStore.Set(provPId, cn)
	return negotiation.Ack(cn), nil
}
