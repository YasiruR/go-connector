package provider

import (
	"fmt"
	"github.com/YasiruR/connector/core"
	"github.com/YasiruR/connector/pkg"
	"github.com/YasiruR/connector/protocols/negotiation"
)

type Service struct {
	states *stateStore
	urn    core.URN
	client core.HTTPClient
}

func New(client core.HTTPClient) core.Provider {
	return &Service{states: newStateStore(), urn: pkg.NewURN(), client: client}
}

func (s *Service) HandleContractRequest(cr negotiation.ContractRequest) (err error) {
	// return error message if offerId is invalid

	// return error message if callbackAddress is invalid

	// associate with existing contract negotiation if providerPid exists and create a new contract
	// negotiation if otherwise
	provPId := cr.ProvPId
	if provPId != `` {
		state, err := s.states.get(provPId)
		if err != nil {
			return fmt.Errorf("get state failed - %w", err)
		}

		if state != negotiation.StateOffered {
			return fmt.Errorf("incompatible state '%s' (expected '%s')", state, negotiation.StateOffered)
		}
		s.states.set(cr.ProvPId, negotiation.StateRequested)
	} else {
		provPId, err = s.urn.New()
		if err != nil {
			return fmt.Errorf("generate - %w", err)
		}
	}

	s.states.set(provPId, negotiation.StateRequested)
	return nil
}

func (s *Service) CreateAsset() {}

func (s *Service) CreatePolicy() {}

func (s *Service) CreateContractDef() {}

func (s *Service) OfferContract() {}

func (s *Service) AgreeContract() {}

func (s *Service) FinalizeContract() {}
