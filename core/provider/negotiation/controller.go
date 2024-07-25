package negotiation

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/connector/domain"
	negotiation2 "github.com/YasiruR/connector/domain/api/dsp/http/negotiation"
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

	ca := negotiation2.ContractAgreement{
		Ctx:     core.Context,
		Type:    negotiation2.TypeContractAgreement,
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

	endpoint := strings.Replace(consumerAddr+negotiation2.ContractAgreementEndpoint, `{`+negotiation2.ParamConsumerPid+`}`, cn.ConsPId, 1)
	res, err := c.client.Send(data, endpoint)
	if err != nil {
		return ``, errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack negotiation2.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return ``, errors.UnmarshalError(negotiation2.ContractAgreementEndpoint, err)
	}

	if err = c.negStore.UpdateState(providerPid, negotiation2.StateAgreed); err != nil {
		return ``, errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	c.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", providerPid, negotiation2.StateAgreed))
	return agreementId, nil
}

func (c *Controller) FinalizeContract(providerPid string) error {
	neg, err := c.negStore.Negotiation(providerPid)
	if err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	event := negotiation2.ContractNegotiationEvent{
		Ctx:       core.Context,
		Type:      negotiation2.TypeNegotiationEvent,
		ProvPId:   providerPid,
		ConsPId:   neg.ConsPId,
		EventType: negotiation2.EventFinalized,
	}

	consumerAddr, err := c.negStore.CallbackAddr(providerPid)
	if err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `CallbackAddr`, err)
	}

	data, err := json.Marshal(event)
	if err != nil {
		return errors.MarshalError(``, err)
	}

	endpoint := strings.Replace(consumerAddr+negotiation2.EventConsumerEndpoint, `{`+negotiation2.ParamConsumerPid+`}`, neg.ConsPId, 1)
	res, err := c.client.Send(data, endpoint)
	if err != nil {
		return errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack negotiation2.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return errors.UnmarshalError(negotiation2.EventConsumerEndpoint, err)
	}

	if err = c.negStore.UpdateState(providerPid, negotiation2.StateFinalized); err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	c.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", providerPid, negotiation2.StateFinalized))
	return nil
}
