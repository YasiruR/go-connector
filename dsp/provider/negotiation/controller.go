package negotiation

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/dsp"
	"github.com/YasiruR/connector/domain/dsp/negotiation"
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
	negStore     stores.ContractNegotiation
	policyStore  stores.Policy
	urn          pkg.URNService
	client       pkg.Client
	log          pkg.Log
}

func NewController(port int, stores domain.Stores, plugins domain.Plugins) *Controller {
	return &Controller{
		callbackAddr: `http://localhost:` + strconv.Itoa(port),
		negStore:     stores.ContractNegotiation,
		policyStore:  stores.Policy,
		urn:          plugins.URNService,
		client:       plugins.Client,
		log:          plugins.Log,
	}
}

func (c *Controller) OfferContract() {}

func (c *Controller) AgreeContract(offerId, providerPid string) (agreementId string, err error) {
	agreementId, err = c.urn.NewURN()
	if err != nil {
		return ``, errors.URNFailed(`providerPid`, `NewURN`, err)
	}

	offer, err := c.policyStore.Offer(offerId)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypePolicy, `Offer`, err)
	}

	assignee, err := c.negStore.Assignee(providerPid)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypeContractNegotiation, `Assignee`, err)
	}

	cn, err := c.negStore.Negotiation(providerPid)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	ca := negotiation.ContractAgreement{
		Ctx:     dsp.Context,
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
		CallbackAddr: c.callbackAddr,
	}

	consumerAddr, err := c.negStore.CallbackAddr(providerPid)
	if err != nil {
		return ``, errors.StoreFailed(stores.TypeContractNegotiation, `CallbackAddr`, err)
	}

	data, err := json.Marshal(ca)
	if err != nil {
		return ``, errors.MarshalError(``, err)
	}

	endpoint := strings.Replace(consumerAddr+negotiation.ContractAgreementEndpoint, `{`+negotiation.ParamConsumerPid+`}`, cn.ConsPId, 1)
	res, err := c.client.Send(data, endpoint)
	if err != nil {
		return ``, errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack negotiation.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return ``, errors.UnmarshalError(negotiation.ContractAgreementEndpoint, err)
	}

	if err = c.negStore.UpdateState(providerPid, negotiation.StateAgreed); err != nil {
		return ``, errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	c.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", providerPid, negotiation.StateAgreed))
	return agreementId, nil
}

func (c *Controller) FinalizeContract(providerPid string) error {
	neg, err := c.negStore.Negotiation(providerPid)
	if err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	event := negotiation.ContractNegotiationEvent{
		Ctx:       dsp.Context,
		Type:      negotiation.TypeNegotiationEvent,
		ProvPId:   providerPid,
		ConsPId:   neg.ConsPId,
		EventType: negotiation.EventFinalized,
	}

	consumerAddr, err := c.negStore.CallbackAddr(providerPid)
	if err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `CallbackAddr`, err)
	}

	data, err := json.Marshal(event)
	if err != nil {
		return errors.MarshalError(``, err)
	}

	endpoint := strings.Replace(consumerAddr+negotiation.EventConsumerEndpoint, `{`+negotiation.ParamConsumerPid+`}`, neg.ConsPId, 1)
	res, err := c.client.Send(data, endpoint)
	if err != nil {
		return errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack negotiation.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return errors.UnmarshalError(negotiation.EventConsumerEndpoint, err)
	}

	if err = c.negStore.UpdateState(providerPid, negotiation.StateFinalized); err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	c.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", providerPid, negotiation.StateFinalized))
	return nil
}
