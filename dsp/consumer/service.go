package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/connector/core/dsp"
	"github.com/YasiruR/connector/core/dsp/negotiation"
	"github.com/YasiruR/connector/core/models/odrl"
	"github.com/YasiruR/connector/core/pkg"
	"github.com/YasiruR/connector/pkg/urn"
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

func New(port int, cnStore *stores.ContractNegotiation, hc pkg.HTTPClient, log pkg.Log) dsp.Consumer {
	return &Service{
		callbackAddr: `http://localhost:` + strconv.Itoa(port),
		cnStore:      cnStore,
		urn:          urn.NewGenerator(),
		client:       hc,
		log:          log,
	}
}

func (s *Service) RequestContract(offerId, providerEndpoint, providerPid string, ot odrl.Target, a odrl.Assigner, act odrl.Action) error {
	// generate consumerPid
	consPId, err := s.urn.New()
	if err != nil {
		return fmt.Errorf("generating URN failed - %w", err)
	}

	// construct payload
	req := negotiation.ContractRequest{
		ProvPId: providerPid,
		ConsPId: consPId,
		Offer: odrl.Offer{
			Id:          offerId,
			Target:      ot,
			Assigner:    a,
			Permissions: []odrl.Rule{{Action: act}}, // should handle constraints
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

	var ack negotiation.Ack
	switch statusCode {
	case 400:
		// read and output error message
		return fmt.Errorf("received 400 status code")
	case 201:
		if err = json.Unmarshal(res, &ack); err != nil {
			return fmt.Errorf("unmarshalling ack failed - %w", err)
		}
		s.log.Info("received ack for contract request", ack)
		s.cnStore.Set(consPId, negotiation.Negotiation(ack))
	default:
		return fmt.Errorf("unexpected status code %d (expected 201 or 400)", statusCode)
	}

	return nil
}

func (s *Service) AcceptContract() {}

func (s *Service) VerifyAgreement() {}

func (s *Service) TerminateContract() {}
