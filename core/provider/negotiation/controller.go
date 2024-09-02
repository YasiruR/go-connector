package negotiation

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api"
	"github.com/YasiruR/connector/domain/api/dsp/http/negotiation"
	"github.com/YasiruR/connector/domain/core"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/models/odrl"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/ror"
	"github.com/YasiruR/connector/domain/stores"
	"strconv"
	"time"
)

type Controller struct {
	callbackAddr string
	cnStore      stores.ContractNegotiationStore
	policyStore  stores.OfferStore
	agrStore     stores.AgreementStore
	urn          pkg.URNService
	client       pkg.Client
	log          pkg.Log
}

func NewController(port int, stores domain.Stores, plugins domain.Plugins) *Controller {
	return &Controller{
		callbackAddr: `http://localhost:` + strconv.Itoa(port),
		cnStore:      stores.ContractNegotiationStore,
		policyStore:  stores.OfferStore,
		agrStore:     stores.AgreementStore,
		urn:          plugins.URNService,
		client:       plugins.Client,
		log:          plugins.Log,
	}
}

func (c *Controller) OfferContract(offerId, providerPid, consumerAddr string) (cnId string, err error) {
	// include datasetId (optional) if provider is the initiator

	ofr, err := c.policyStore.Offer(offerId)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypeOffer, `Offer`, err)
	}

	var consumerPid, endpoint string
	if providerPid != `` {
		cn, err := c.cnStore.Negotiation(providerPid)
		if err != nil {
			return ``, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
		}

		if cn.State != negotiation.StateRequested {
			return ``, ror.IncompatibleState(core.NegotiationProtocol, string(cn.State),
				string(negotiation.StateRequested))
		}

		assignee, err := c.cnStore.Assignee(providerPid)
		if err != nil {
			return ``, errors.StoreFailed(stores.TypeContractNegotiation, `Assignee`, err)
		}

		// if consumer's address is not provided, use the stored one. If provided, it will override
		// the existing address.
		if consumerAddr == `` {
			consumerAddr, err = c.cnStore.CallbackAddr(providerPid)
			if err != nil {
				return ``, errors.StoreFailed(stores.TypeContractNegotiation, `CallbackAddr`, err)
			}
		}

		consumerPid = cn.ConsPId
		ofr.Assignee = assignee
		endpoint = api.SetParamConsumerPid(negotiation.ContractOfferToRequestEndpoint, consumerPid)
		c.log.Trace(fmt.Sprintf("found an existing contract negotiation for the request (id: %s)", providerPid))
	} else {
		if consumerAddr == `` {
			return ``, ror.MissingRequiredAttr(`consumer address`,
				`must be provided when Provider is the initiator`)
		}

		providerPid, err = c.urn.NewURN()
		if err != nil {
			return ``, errors.URNFailed(`providerPid`, `NewURN`, err)
		}
		endpoint = negotiation.ContractOfferEndpoint
	}

	c.cnStore.SetParticipants(providerPid, consumerAddr, ofr.Assigner, ofr.Assignee)
	req := negotiation.ContractOffer{
		Ctx:          core.Context,
		Type:         negotiation.MsgTypeContractOffer,
		ProvPId:      providerPid,
		ConsPId:      consumerPid,
		Offer:        ofr,
		CallbackAddr: c.callbackAddr,
	}
	// todo offer must have a target but not in policies

	ack, err := c.send(providerPid, endpoint, req)
	if err != nil {
		return ``, ror.CustomFuncError(`send`, err)
	}

	if !c.validAck(providerPid, ack, negotiation.StateOffered) {
		return ``, ror.InvalidAck(`ContractOffer`, ``, ack)
	}

	ack.Type = negotiation.MsgTypeNegotiation
	c.cnStore.AddNegotiation(providerPid, negotiation.Negotiation(ack))

	c.log.Info(fmt.Sprintf("provider controller updated negotiation state (id: %s, state: %s)", providerPid, negotiation.StateOffered))
	return providerPid, nil
}

func (c *Controller) AgreeContract(offerId, providerPid string) (agreementId string, err error) {
	cn, err := c.cnStore.Negotiation(providerPid)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	if cn.State != negotiation.StateRequested && cn.State != negotiation.StateAccepted {
		return ``, ror.IncompatibleState(core.NegotiationProtocol, string(cn.State),
			string(negotiation.StateRequested)+" or "+string(negotiation.StateAccepted))
	}

	agreementId, err = c.urn.NewURN()
	if err != nil {
		return ``, errors.URNFailed(`providerPid`, `NewURN`, err)
	}

	offer, err := c.policyStore.Offer(offerId)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypeOffer, `Offer`, err)
	}

	assignee, err := c.cnStore.Assignee(providerPid)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypeContractNegotiation, `Assignee`, err)
	}

	req := negotiation.ContractAgreement{
		Ctx:     core.Context,
		Type:    negotiation.MsgTypeContractAgreement,
		ProvPId: cn.ProvPId,
		ConsPId: cn.ConsPId,
		Agreement: odrl.Agreement{
			Id:          agreementId,
			Type:        odrl.TypeAgreement,
			Target:      offer.Target,
			Assigner:    offer.Assigner,
			Assignee:    assignee,
			Timestamp:   time.Now().UTC().String(), // change format into XSD
			Permissions: offer.Permissions,         // should be able to select a subset in future
		},
		CallbackAddr: c.callbackAddr,
	}
	// todo agreement must have a target but not in policies

	if _, err = c.send(providerPid, api.SetParamConsumerPid(negotiation.ContractAgreementEndpoint,
		cn.ConsPId), req); err != nil {
		return ``, ror.CustomFuncError(`send`, err)
	}

	c.agrStore.AddAgreement(req.Agreement.Id, req.Agreement)
	c.log.Trace(fmt.Sprintf("stored contract agreement (id: %s) for negotation (id: %s)",
		req.Agreement.Id, providerPid))

	if err = c.cnStore.UpdateState(providerPid, negotiation.StateAgreed); err != nil {
		return ``, errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	c.log.Debug(fmt.Sprintf("provider controller updated negotiation state (id: %s, state: %s)",
		providerPid, negotiation.StateAgreed))
	return agreementId, nil
}

func (c *Controller) FinalizeContract(providerPid string) error {
	cn, err := c.cnStore.Negotiation(providerPid)
	if err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	if cn.State != negotiation.StateVerified {
		return ror.IncompatibleState(core.NegotiationProtocol, string(cn.State), string(negotiation.StateVerified))
	}

	event := negotiation.ContractNegotiationEvent{
		Ctx:       core.Context,
		Type:      negotiation.MsgTypeNegotiationEvent,
		ProvPId:   providerPid,
		ConsPId:   cn.ConsPId,
		EventType: negotiation.EventFinalized,
	}

	if _, err = c.send(providerPid, api.SetParamPid(negotiation.EventsEndpoint, cn.ConsPId), event); err != nil {
		return ror.CustomFuncError(`send`, err)
	}

	if err = c.cnStore.UpdateState(providerPid, negotiation.StateFinalized); err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	c.log.Info(fmt.Sprintf("provider controller updated negotiation state (id: %s, state: %s)",
		providerPid, negotiation.StateFinalized))
	return nil
}

func (c *Controller) send(providerPid, endpoint string, req any) (negotiation.Ack, error) {
	consumerAddr, err := c.cnStore.CallbackAddr(providerPid)
	if err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `CallbackAddr`, err)
	}

	data, err := json.Marshal(req)
	if err != nil {
		return negotiation.Ack{}, errors.MarshalError(``, err)
	}

	res, err := c.client.Send(data, consumerAddr+endpoint)
	if err != nil {
		return negotiation.Ack{}, errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack negotiation.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return negotiation.Ack{}, errors.UnmarshalError(endpoint, err)
	}

	return ack, nil
}

func (c *Controller) validAck(providerPid string, ack negotiation.Ack, state negotiation.State) bool {
	if ack.ProvPId != providerPid {
		return false
	}

	if ack.State != state {
		return false
	}

	return true
}
