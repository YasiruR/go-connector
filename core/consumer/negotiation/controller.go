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
	"github.com/YasiruR/connector/domain/stores"
	"strconv"
)

type Controller struct {
	callbackAddr string
	assigneeId   string
	catalog      stores.ConsumerCatalog
	cnStore      stores.ContractNegotiationStore
	urn          pkg.URNService
	client       pkg.Client
	log          pkg.Log
}

func NewController(port int, stores domain.Stores, plugins domain.Plugins) *Controller {
	return &Controller{
		callbackAddr: `http://localhost:` + strconv.Itoa(port),
		catalog:      stores.ConsumerCatalog,
		cnStore:      stores.ContractNegotiationStore,
		urn:          plugins.URNService,
		client:       plugins.Client,
		log:          plugins.Log,
	}
}

func (c *Controller) RequestContract(consumerPid, providerAddr, offerId string, constraints map[string]string) (cnId string, err error) {
	var providerPid, endpoint string
	if consumerPid != `` {
		cn, err := c.cnStore.Negotiation(consumerPid)
		if err != nil {
			return ``, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
		}

		if cn.State != negotiation.StateOffered {
			return ``, errors.IncompatibleValues(`state`, string(cn.State), string(negotiation.StateOffered))
		}

		// if provider's address is not provided, use the stored one. If provided, it will override
		// the existing address. (alternatively, can fetch from consumer catalog)
		if providerAddr == `` {
			providerAddr, err = c.cnStore.CallbackAddr(consumerPid)
			if err != nil {
				return ``, errors.StoreFailed(stores.TypeContractNegotiation, `CallbackAddr`, err)
			}
		}

		providerPid = cn.ProvPId
		endpoint = api.SetParamProviderPid(negotiation.ContractRequestToOfferEndpoint, providerPid)
		c.log.Trace(fmt.Sprintf("found an existing contract negotiation for the request (id: %s)", consumerPid))
	} else {
		// generate consumerPid
		consumerPid, err = c.urn.NewURN()
		if err != nil {
			return ``, errors.URNFailed(`consumerPid`, `NewURN`, err)
		}
		endpoint = negotiation.ContractRequestEndpoint
	}

	// fetch offer from store
	ofr, err := c.setConstraints(offerId, constraints)
	if err != nil {
		return ``, errors.CustomFuncError(`setConstraints`, err)
	}

	c.cnStore.SetParticipants(consumerPid, providerAddr, ofr.Assigner, ofr.Assignee)
	req := negotiation.ContractRequest{
		Ctx:          core.Context,
		Type:         negotiation.MsgTypeContractRequest,
		ConsPId:      consumerPid,
		ProvPId:      providerPid,
		Offer:        ofr,
		CallbackAddr: c.callbackAddr,
	}

	ack, err := c.send(consumerPid, endpoint, req)
	if err != nil {
		return ``, errors.CustomFuncError(`send`, err)
	}

	if !c.validAck(consumerPid, ack, negotiation.StateRequested) {
		return ``, errors.InvalidAck(`ContractRequest`, ack)
	}

	ack.Type = negotiation.MsgTypeNegotiation
	c.cnStore.AddNegotiation(consumerPid, negotiation.Negotiation(ack))

	c.log.Trace(fmt.Sprintf("consumer stored contract negotiation (id: %s, assigner: %s, assignee: %s, address: %s)",
		consumerPid, ofr.Assigner, ofr.Assignee, providerAddr))
	c.log.Debug(fmt.Sprintf("consumer controller updated negotiation state (id: %s, state: %s)", consumerPid, negotiation.StateRequested))
	return consumerPid, nil
}

func (c *Controller) AcceptOffer(consumerPid string) error {
	cn, err := c.cnStore.Negotiation(consumerPid)
	if err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	if cn.State != negotiation.StateOffered {
		return errors.IncompatibleValues(`state`, string(cn.State), string(negotiation.StateOffered))
	}

	event := negotiation.ContractNegotiationEvent{
		Ctx:       core.Context,
		Type:      negotiation.MsgTypeNegotiationEvent,
		ProvPId:   cn.ProvPId,
		ConsPId:   consumerPid,
		EventType: negotiation.EventAccepted,
	}

	if _, err = c.send(consumerPid, api.SetParamPid(negotiation.EventsEndpoint, cn.ProvPId), event); err != nil {
		return errors.CustomFuncError(`send`, err)
	}

	if err = c.cnStore.UpdateState(consumerPid, negotiation.StateAccepted); err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	c.log.Debug(fmt.Sprintf("consumer controller updated negotiation state (id: %s, state: %s)", consumerPid, negotiation.StateAccepted))
	return nil
}

func (c *Controller) VerifyAgreement(consumerPid string) error {
	cn, err := c.cnStore.Negotiation(consumerPid)
	if err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	if cn.State != negotiation.StateAgreed {
		return errors.IncompatibleValues(`state`, string(cn.State), string(negotiation.StateAgreed))
	}

	req := negotiation.ContractVerification{
		Ctx:     core.Context,
		Type:    negotiation.MsgTypeAgreementVerification,
		ProvPId: cn.ProvPId,
		ConsPId: consumerPid,
	}

	if _, err = c.send(consumerPid, api.SetParamProviderPid(negotiation.AgreementVerificationEndpoint, cn.ProvPId), req); err != nil {
		return errors.CustomFuncError(`send`, err)
	}

	if err = c.cnStore.UpdateState(consumerPid, negotiation.StateVerified); err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	c.log.Debug(fmt.Sprintf("consumer controller updated negotiation state (id: %s, state: %s)", consumerPid, negotiation.StateVerified))
	return nil
}

func (c *Controller) TerminateContract(consumerPid, code string, reasons []string) error {
	cn, err := c.cnStore.Negotiation(consumerPid)
	if err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	var rsnList []negotiation.Reason
	for _, r := range reasons {
		rsnList = append(rsnList, negotiation.Reason{
			Value:    r,
			Language: "en",
		})
	}

	req := negotiation.ContractTermination{
		Ctx:     core.Context,
		Type:    negotiation.MsgTypeTermination,
		ProvPId: cn.ProvPId,
		ConsPId: consumerPid,
		Code:    code,
		Reason:  rsnList,
	}

	if _, err = c.send(consumerPid, api.SetParamPid(negotiation.TerminateEndpoint, cn.ProvPId), req); err != nil {
		return errors.CustomFuncError(`send`, err)
	}

	// clear all store entries for the contract negotiation
	if err = c.cnStore.UpdateState(consumerPid, negotiation.StateTerminated); err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	c.log.Info(fmt.Sprintf("consumer terminated the negotiation flow successfully (id: %s)", consumerPid))
	return nil
}

func (c *Controller) send(consumerPid, endpoint string, req any) (negotiation.Ack, error) {
	providerAddr, err := c.cnStore.CallbackAddr(consumerPid)
	if err != nil {
		return negotiation.Ack{}, errors.StoreFailed(stores.TypeContractNegotiation, `CallBackAddr`, err)
	}

	data, err := json.Marshal(req)
	if err != nil {
		return negotiation.Ack{}, errors.MarshalError(negotiation.TerminateEndpoint, err)
	}

	res, err := c.client.Send(data, providerAddr+endpoint)
	if err != nil {
		return negotiation.Ack{}, errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack negotiation.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return negotiation.Ack{}, errors.UnmarshalError(endpoint, err)
	}

	return ack, nil
}

func (c *Controller) setConstraints(offerId string, vals map[string]string) (odrl.Offer, error) {
	var permList []odrl.Rule
	ofr, err := c.catalog.Offer(offerId)
	if err != nil {
		return odrl.Offer{}, errors.StoreFailed(stores.TypeOffer, `Offer`, err)
	}

	for _, perm := range ofr.Permissions {
		var consList []odrl.Constraint
		for _, cons := range perm.Constraints {
			val, ok := vals[cons.LeftOperand]
			if !ok {
				return odrl.Offer{}, errors.MissingRequiredAttr(cons.LeftOperand, `mandatory constraint`)
			}
			cons.RightOperand = val
			consList = append(consList, cons)
		}

		permList = append(permList, odrl.Rule{
			Action:      perm.Action,
			Constraints: consList,
		})
	}

	ofr.Assignee = odrl.Assignee(c.assigneeId)
	ofr.Permissions = permList
	return ofr, nil
}

func (c *Controller) validAck(pid string, ack negotiation.Ack, state negotiation.State) bool {
	if ack.ConsPId != pid {
		return false
	}

	if ack.State != state {
		return false
	}

	return true
}
