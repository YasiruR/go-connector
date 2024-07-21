package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/connector/core/dsp"
	"github.com/YasiruR/connector/core/dsp/negotiation"
	"github.com/YasiruR/connector/core/errors"
	"github.com/YasiruR/connector/core/pkg"
	"github.com/YasiruR/connector/core/protocols/odrl"
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

func New(port int, cnStore *stores.ContractNegotiation, urn pkg.URN, hc pkg.HTTPClient, log pkg.Log) dsp.Consumer {
	return &Service{
		callbackAddr: `http://localhost:` + strconv.Itoa(port),
		cnStore:      cnStore,
		urn:          urn,
		client:       hc,
		log:          log,
	}
}

func (s *Service) RequestContract(offerId, providerEndpoint, providerPid, target, assigner, assignee, action string) error {
	// generate consumerPid
	consPId, err := s.urn.New()
	if err != nil {
		return errors.URNFailed(`consumerPid`, `New`, err)
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
		return errors.MarshalError(``, err)
	}

	statusCode, res, err := s.client.Post(providerEndpoint+negotiation.RequestContractEndpoint, data)
	if err != nil {
		return errors.PkgFailed(pkg.TypeHTTPClient, `Post`, err)
	}

	var ack negotiation.Ack
	switch statusCode {
	case 400:
		// read and output error message
		return errors.InvalidStatusCode(400, 200)
	case 201:
		if err = json.Unmarshal(res, &ack); err != nil {
			return fmt.Errorf("unmarshalling ack failed - %w", err)
		}
		s.log.Trace("received ack for the contract request", ack)
		s.cnStore.Set(consPId, negotiation.Negotiation(ack))
		s.cnStore.SetAssigner(consPId, odrl.Assigner(assigner))
		s.cnStore.SetAssignee(consPId, odrl.Assignee(assignee))
		s.log.Info(fmt.Sprintf("stored contract negotiation (id: %s, assigner: %s, assignee: %s)", consPId, assigner, assignee))
	default:
		return errors.InvalidStatusCode(statusCode, 200)
	}

	return nil
}

func (s *Service) AcceptContract() {}

func (s *Service) VerifyAgreement() {}

func (s *Service) TerminateContract() {}
