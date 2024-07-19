package provider

import (
	"fmt"
	"github.com/YasiruR/connector/core"
	"github.com/YasiruR/connector/pkg"
	"github.com/YasiruR/connector/protocols/negotiation"
	"github.com/YasiruR/connector/stores"
	"strconv"
)

type Service struct {
	callbackAddr string
	cnStore      *stores.ContractNegotiation
	urn          core.URN
	client       core.HTTPClient
	log          core.Log
}

func New(port int, cnStore *stores.ContractNegotiation, c core.HTTPClient, log core.Log) core.Provider {
	return &Service{
		callbackAddr: `http://localhost:` + strconv.Itoa(port),
		cnStore:      cnStore,
		urn:          pkg.NewURN(),
		client:       c,
		log:          log,
	}
}

func (s *Service) CreateAsset() {}

func (s *Service) CreatePolicy() {}

func (s *Service) CreateContractDef() {}

func (s *Service) OfferContract() {}

func (s *Service) AgreeContract() {}

func (s *Service) FinalizeContract() {}

func (s *Service) HandleNegotiationsRequest(providerPid string) (negotiation.Ack, error) {
	ack, err := s.cnStore.Get(providerPid)
	if err != nil {
		return negotiation.Ack{}, fmt.Errorf("get negotiation failed - %w", err)
	}

	return ack, nil
}

func (s *Service) HandleContractRequest(cr negotiation.ContractRequest) (ack negotiation.Ack, err error) {
	// return error message if offerId is invalid

	// return error message if callbackAddress is invalid

	// associate with existing contract negotiation if providerPid exists and create a new contract
	// negotiation if otherwise
	provPId := cr.ProvPId
	if provPId != `` {
		cn, err := s.cnStore.Get(provPId)
		if err != nil {
			return negotiation.Ack{}, fmt.Errorf("get state failed - %w", err)
		}

		if cn.State != negotiation.StateOffered {
			return negotiation.Ack{}, fmt.Errorf("incompatible state '%s' (expected '%s')", cn.State, negotiation.StateOffered)
		}

		if cn.ConsPId != cr.ConsPId {
			return negotiation.Ack{}, fmt.Errorf("incompatible consumerPid '%s' (expected '%s')", cn.ConsPId, cr.ConsPId)
		}

		cn.State = negotiation.StateRequested
		ack = cn
		s.log.Info("updated existing contract negotiation", ack)
	} else {
		provPId, err = s.urn.New()
		if err != nil {
			return negotiation.Ack{}, fmt.Errorf("generate new URN failed - %w", err)
		}

		ack = negotiation.NewAck()
		ack.ConsPId = cr.ConsPId
		ack.ProvPId = provPId
		ack.State = negotiation.StateRequested
		s.log.Info("stored new contract negotiation", ack)
	}

	s.cnStore.Set(provPId, ack)
	return ack, nil
}
