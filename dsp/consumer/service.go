package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/connector/core/dsp"
	negotiation2 "github.com/YasiruR/connector/core/dsp/negotiation"
	pkg2 "github.com/YasiruR/connector/core/pkg"
	"github.com/YasiruR/connector/pkg/urn"
	"github.com/YasiruR/connector/stores"
	"strconv"
)

type Service struct {
	callbackAddr string
	cnStore      *stores.ContractNegotiation
	urn          pkg2.URN
	client       pkg2.HTTPClient
	log          pkg2.Log
}

func New(port int, cnStore *stores.ContractNegotiation, hc pkg2.HTTPClient, log pkg2.Log) dsp.Consumer {
	return &Service{
		callbackAddr: `http://localhost:` + strconv.Itoa(port),
		cnStore:      cnStore,
		urn:          urn.NewGenerator(),
		client:       hc,
		log:          log,
	}
}

func (s *Service) RequestContract(offerId, providerEndpoint, providerPid, odrlTarget, assigner, action string) error {
	// generate consumerPid
	consPId, err := s.urn.New()
	if err != nil {
		return fmt.Errorf("generating URN failed - %w", err)
	}

	// construct payload
	req := negotiation2.ContractRequest{
		ProvPId: providerPid,
		ConsPId: consPId,
		Offer: negotiation2.Offer{
			Id:          offerId,
			Target:      odrlTarget,
			Assigner:    assigner,
			Permissions: []negotiation2.Permission{{Action: action}}, // should handle constraints
		},
		CallbackAddr: s.callbackAddr,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshalling request failed - %w", err)
	}

	statusCode, res, err := s.client.Post(providerEndpoint+negotiation2.RequestContractEndpoint, data)
	if err != nil {
		return fmt.Errorf("posting request failed - %w", err)
	}

	var ack negotiation2.Ack
	switch statusCode {
	case 400:
		// read and output error message
		return fmt.Errorf("received 400 status code")
	case 201:
		if err = json.Unmarshal(res, &ack); err != nil {
			return fmt.Errorf("unmarshalling ack failed - %w", err)
		}
		s.cnStore.Set(consPId, ack)
		s.log.Info("received ack for contract request", ack)
	default:
		return fmt.Errorf("unexpected status code %d (expected 201 or 400)", statusCode)
	}

	return nil
}

func (s *Service) AcceptContract() {}

func (s *Service) VerifyAgreement() {}

func (s *Service) TerminateContract() {}
