package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/connector/core"
	"github.com/YasiruR/connector/pkg"
	"github.com/YasiruR/connector/protocols/negotiation"
	"strconv"
)

type Service struct {
	callbackAddr string
	states       *stateStore
	providers    *providerStore
	urn          core.URN
	client       core.HTTPClient
}

func New(port int, client core.HTTPClient) core.Consumer {
	return &Service{
		callbackAddr: `http://localhost:` + strconv.Itoa(port),
		states:       newStateStore(),
		providers:    newProviderStore(),
		urn:          pkg.NewURN(),
		client:       client,
	}
}

func (s *Service) RequestContract(offerId, providerEndpoint, providerPid, odrlTarget, assigner, action string) error {
	// generate consumerPid
	consId, err := s.urn.New()
	if err != nil {
		return fmt.Errorf("generating URN failed - %w", err)
	}

	// construct payload
	req := negotiation.ContractRequest{
		ProvPId: providerPid,
		ConsPId: consId,
		Offer: negotiation.Offer{
			Id:          offerId,
			Target:      odrlTarget,
			Assigner:    assigner,
			Permissions: []negotiation.Permission{{Action: action}}, // should handle constraints
		},
		CallbackAddr: s.callbackAddr,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshalling request failed - %w", err)
	}

	statusCode, res, err := s.client.Post(providerEndpoint+negotiation.RequestContractEndpoint, data)
	if err != nil {
		return fmt.Errorf("posting request failed - %w", err)
	}

	switch statusCode {
	case 400:
		// read and output error message
		s.states.set(consId, negotiation.StateTerminated)
	case 201:
		var ack negotiation.Ack
		if err = json.Unmarshal(res, &ack); err != nil {
			return fmt.Errorf("unmarshalling ack failed - %w", err)
		}
		s.states.set(consId, negotiation.StateRequested)
	default:
		return fmt.Errorf("unexpected status code %d (expected 201 or 400)", statusCode)
	}

	return nil
}

func (s *Service) AcceptContract() {}

func (s *Service) VerifyAgreement() {}

func (s *Service) TerminateContract() {}
