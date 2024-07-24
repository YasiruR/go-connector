package negotiation

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/connector/core"
	"github.com/YasiruR/connector/core/dsp"
	"github.com/YasiruR/connector/core/dsp/negotiation"
	"github.com/YasiruR/connector/core/errors"
	"github.com/YasiruR/connector/core/pkg"
	"github.com/YasiruR/connector/core/protocols/odrl"
	"github.com/YasiruR/connector/core/stores"
	"strconv"
	"strings"
)

type Controller struct {
	callbackAddr string
	negStore     stores.ContractNegotiation
	urn          pkg.URNService
	client       pkg.Client
	log          pkg.Log
}

func NewController(port int, stores core.Stores, plugins core.Plugins) *Controller {
	return &Controller{
		callbackAddr: `http://localhost:` + strconv.Itoa(port),
		negStore:     stores.ContractNegotiation,
		urn:          plugins.URNService,
		client:       plugins.Client,
		log:          plugins.Log,
	}
}

func (c *Controller) RequestContract(offerId, providerEndpoint, providerPid, target, assigner, assignee, action string) (negotiationId string, err error) {
	// generate consumerPid
	consPId, err := c.urn.NewURN()
	if err != nil {
		return ``, errors.URNFailed(`consumerPid`, `NewURN`, err)
	}

	// construct payload
	req := negotiation.ContractRequest{
		Ctx:     dsp.Context,
		Type:    negotiation.TypeContractRequest,
		ProvPId: providerPid,
		ConsPId: consPId,
		Offer: odrl.Offer{
			Id:          offerId,
			Target:      odrl.Target(target),
			Assigner:    odrl.Assigner(assigner),
			Assignee:    odrl.Assignee(assignee),
			Permissions: []odrl.Rule{{Action: odrl.Action(action)}}, // should handle constraints
		},
		CallbackAddr: c.callbackAddr,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return ``, errors.MarshalError(``, err)
	}

	res, err := c.client.Send(data, providerEndpoint+negotiation.ContractRequestEndpoint)
	if err != nil {
		return ``, errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack negotiation.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return ``, errors.UnmarshalError(providerEndpoint+negotiation.ContractRequestEndpoint, err)
	}

	c.negStore.Set(consPId, negotiation.Negotiation(ack))
	c.negStore.SetAssignee(consPId, odrl.Assignee(assignee))
	c.log.Trace(fmt.Sprintf("stored contract negotiation (id: %s, assigner: %s, assignee: %s)", consPId, assigner, assignee))
	c.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", consPId, negotiation.StateRequested))
	return consPId, nil
}

func (c *Controller) AcceptContract() {}

func (c *Controller) VerifyAgreement(consumerPid string) error {
	neg, err := c.negStore.Negotiation(consumerPid)
	if err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	req := negotiation.ContractVerification{
		Ctx:     dsp.Context,
		Type:    negotiation.TypeAgreementVerification,
		ProvPId: neg.ProvPId,
		ConsPId: consumerPid,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return errors.MarshalError(``, err)
	}

	providerAddr, err := c.negStore.CallbackAddr(consumerPid)
	if err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `CallBackAddr`, err)
	}

	endpoint := strings.Replace(providerAddr+negotiation.AgreementVerificationEndpoint, `{`+negotiation.ParamProviderId+`}`, neg.ProvPId, 1)
	res, err := c.client.Send(data, endpoint)
	if err != nil {
		return errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack negotiation.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return errors.UnmarshalError(providerAddr+negotiation.AgreementVerificationEndpoint, err)
	}

	if err = c.negStore.UpdateState(consumerPid, negotiation.StateVerified); err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	c.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", consumerPid, negotiation.StateVerified))
	return nil
}

func (c *Controller) TerminateContract() {}
