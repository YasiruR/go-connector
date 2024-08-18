package negotiation

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/dsp/http/negotiation"
	"github.com/YasiruR/connector/domain/core"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/models/odrl"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/stores"
	"strconv"
	"strings"
	"time"
)

type Controller struct {
	callbackAddr string
	cnStore      stores.ContractNegotiation
	policyStore  stores.Policy
	agrStore     stores.Agreement
	urn          pkg.URNService
	client       pkg.Client
	log          pkg.Log
}

func NewController(port int, stores domain.Stores, plugins domain.Plugins) *Controller {
	return &Controller{
		callbackAddr: `http://localhost:` + strconv.Itoa(port),
		cnStore:      stores.ContractNegotiation,
		policyStore:  stores.Policy,
		agrStore:     stores.Agreement,
		urn:          plugins.URNService,
		client:       plugins.Client,
		log:          plugins.Log,
	}
}

func (c *Controller) OfferContract(offerId, providerPid, consumerAddr string) (cnId string, err error) {
	// include datasetId (optional) if provider is the initiator

	ofr, err := c.policyStore.GetOffer(offerId)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypePolicy, `GetOffer`, err)
	}

	var endpoint string
	var consumerPid string
	if providerPid != `` {
		cn, err := c.cnStore.GetNegotiation(providerPid)
		if err != nil {
			return ``, errors.StoreFailed(stores.TypeContractNegotiation, `GetNegotiation`, err)
		}

		if cn.State != negotiation.StateRequested {
			return ``, errors.IncompatibleValues(`state`, string(cn.State), string(negotiation.StateRequested))
		}

		assignee, err := c.cnStore.Assignee(providerPid)
		if err != nil {
			return ``, errors.StoreFailed(stores.TypeContractNegotiation, `Assignee`, err)
		}

		consumerAddr, err = c.cnStore.CallbackAddr(providerPid)
		if err != nil {
			return ``, errors.StoreFailed(stores.TypeContractNegotiation, `CallbackAddr`, err)
		}

		consumerPid = cn.ConsPId
		ofr.Assignee = assignee
		endpoint = strings.Replace(negotiation.ContractOfferToRequestEndpoint, `{`+negotiation.ParamConsumerPid+`}`, consumerPid, 1)
	} else {
		if consumerAddr == `` {
			return ``, errors.MissingRequiredAttr(`consumer address`, `must be provided when Provider is the initiator`)
		}

		endpoint = negotiation.ContractOfferEndpoint
		providerPid, err = c.urn.NewURN()
		if err != nil {
			return ``, errors.URNFailed(`providerPid`, `NewURN`, err)
		}
	}

	co := negotiation.ContractOffer{
		Ctx:          core.Context,
		Type:         negotiation.MsgTypeContractOffer,
		ProvPId:      providerPid,
		ConsPId:      consumerPid,
		Offer:        ofr,
		CallbackAddr: c.callbackAddr,
	}

	data, err := json.Marshal(co)
	if err != nil {
		return ``, errors.MarshalError(endpoint, err)
	}

	res, err := c.client.Send(data, consumerAddr+endpoint)
	if err != nil {
		return ``, errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack negotiation.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return ``, errors.UnmarshalError(endpoint, err)
	}

	// validate ack first
	c.cnStore.Set(providerPid, negotiation.Negotiation(ack))
	c.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", providerPid, negotiation.StateOffered))
	return providerPid, nil
}

func (c *Controller) AgreeContract(offerId, providerPid string) (agreementId string, err error) {
	agreementId, err = c.urn.NewURN()
	if err != nil {
		return ``, errors.URNFailed(`providerPid`, `NewURN`, err)
	}

	offer, err := c.policyStore.GetOffer(offerId)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypePolicy, `GetOffer`, err)
	}

	assignee, err := c.cnStore.Assignee(providerPid)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypeContractNegotiation, `Assignee`, err)
	}

	cn, err := c.cnStore.GetNegotiation(providerPid)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypeContractNegotiation, `GetNegotiation`, err)
	}

	ca := negotiation.ContractAgreement{
		Ctx:     core.Context,
		Type:    negotiation.MsgTypeContractAgreement,
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
		CallbackAddr: c.callbackAddr,
	}

	consumerAddr, err := c.cnStore.CallbackAddr(providerPid)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypeContractNegotiation, `CallbackAddr`, err)
	}

	endpoint := strings.Replace(consumerAddr+negotiation.ContractAgreementEndpoint, `{`+negotiation.ParamConsumerPid+`}`, cn.ConsPId, 1)
	data, err := json.Marshal(ca)
	if err != nil {
		return ``, errors.MarshalError(endpoint, err)
	}

	res, err := c.client.Send(data, endpoint)
	if err != nil {
		return ``, errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack negotiation.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return ``, errors.UnmarshalError(negotiation.ContractAgreementEndpoint, err)
	}

	c.agrStore.Set(ca.ConsPId, ca.Agreement)
	c.log.Trace(fmt.Sprintf("stored contract agreement (id: %s) for negotation (id: %s)", ca.Agreement.Id, providerPid))
	if err = c.cnStore.UpdateState(providerPid, negotiation.StateAgreed); err != nil {
		return ``, errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	c.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", providerPid, negotiation.StateAgreed))
	return agreementId, nil
}

func (c *Controller) FinalizeContract(providerPid string) error {
	neg, err := c.cnStore.GetNegotiation(providerPid)
	if err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `GetNegotiation`, err)
	}

	event := negotiation.ContractNegotiationEvent{
		Ctx:       core.Context,
		Type:      negotiation.MsgTypeNegotiationEvent,
		ProvPId:   providerPid,
		ConsPId:   neg.ConsPId,
		EventType: negotiation.EventFinalized,
	}

	consumerAddr, err := c.cnStore.CallbackAddr(providerPid)
	if err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `CallbackAddr`, err)
	}

	data, err := json.Marshal(event)
	if err != nil {
		return errors.MarshalError(``, err)
	}

	endpoint := strings.Replace(consumerAddr+negotiation.EventsEndpoint, `{`+negotiation.ParamContractId+`}`, neg.ConsPId, 1)
	res, err := c.client.Send(data, endpoint)
	if err != nil {
		return errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack negotiation.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return errors.UnmarshalError(negotiation.EventsEndpoint, err)
	}

	if err = c.cnStore.UpdateState(providerPid, negotiation.StateFinalized); err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	c.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", providerPid, negotiation.StateFinalized))
	return nil
}
