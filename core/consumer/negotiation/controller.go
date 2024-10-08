package negotiation

import (
	"encoding/json"
	defaultErr "errors"
	"fmt"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api"
	"github.com/YasiruR/connector/domain/api/dsp/http/negotiation"
	"github.com/YasiruR/connector/domain/boot"
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

func NewController(cfg boot.Config, stores domain.Stores, plugins domain.Plugins) *Controller {
	return &Controller{
		callbackAddr: cfg.Servers.IP + `:` + strconv.Itoa(cfg.Servers.DSP.HTTP.Port),
		assigneeId:   cfg.DataSpace.AssigneeId,
		catalog:      stores.ConsumerCatalog,
		cnStore:      stores.ContractNegotiationStore,
		urn:          plugins.URNService,
		client:       plugins.Client,
		log:          plugins.Log,
	}
}

func (c *Controller) RequestContract(consumerPid, providerAddr, offerId string,
	constraints map[string]string) (cnId string, err error) {

	var providerPid, endpoint string
	if consumerPid != `` {
		cn, err := c.cnStore.Negotiation(consumerPid)
		if err != nil {
			if defaultErr.Is(err, stores.TypeInvalidKey) {
				return ``, errors.Client(errors.InvalidKey(stores.TypeContractNegotiation,
					`contract negotiation id`, err))
			}
			return ``, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
		}

		if cn.State != negotiation.StateOffered {
			return ``, errors.Client(errors.StateError(`request contract`, string(cn.State)))
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
			return ``, errors.PkgError(pkg.TypeURN, `NewURN`, err, `contract negotiation id`)
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

	if errMsg := c.validAck(consumerPid, ack, negotiation.StateRequested); errMsg != `` {
		return ``, errors.Client(errors.InvalidAckError(`ContractRequest`, errMsg, ack))
	}

	ack.Type = negotiation.MsgTypeNegotiation
	c.cnStore.AddNegotiation(consumerPid, negotiation.Negotiation(ack))

	c.log.Trace(fmt.Sprintf("consumer stored contract negotiation (id: %s, assigner: %s, "+
		"assignee: %s, address: %s)", consumerPid, ofr.Assigner, ofr.Assignee, providerAddr))
	c.log.Debug(fmt.Sprintf("consumer controller updated negotiation state (id: %s, state: %s)",
		consumerPid, negotiation.StateRequested))
	return consumerPid, nil
}

func (c *Controller) AcceptOffer(consumerPid string) error {
	cn, err := c.cnStore.Negotiation(consumerPid)
	if err != nil {
		if defaultErr.Is(err, stores.TypeInvalidKey) {
			return errors.Client(errors.InvalidKey(stores.TypeContractNegotiation,
				`contract negotiation id`, err))
		}
		return errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	if cn.State != negotiation.StateOffered {
		return errors.Client(errors.StateError(`accept offer`, string(cn.State)))
	}

	event := negotiation.ContractNegotiationEvent{
		Ctx:       core.Context,
		Type:      negotiation.MsgTypeNegotiationEvent,
		ProvPId:   cn.ProvPId,
		ConsPId:   consumerPid,
		EventType: negotiation.EventAccepted,
	}

	ack, err := c.send(consumerPid, api.SetParamPid(negotiation.EventsEndpoint, cn.ProvPId), event)
	if err != nil {
		return errors.CustomFuncError(`send`, err)
	}

	if errMsg := c.validAck(consumerPid, ack, negotiation.StateAccepted); errMsg != `` {
		return errors.Client(errors.InvalidAckError(`AcceptOffer`, errMsg, ack))
	}

	if err = c.cnStore.UpdateState(consumerPid, negotiation.StateAccepted); err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	c.log.Debug(fmt.Sprintf("consumer controller updated negotiation state (id: %s, state: %s)",
		consumerPid, negotiation.StateAccepted))
	return nil
}

func (c *Controller) VerifyAgreement(consumerPid string) error {
	cn, err := c.cnStore.Negotiation(consumerPid)
	if err != nil {
		if defaultErr.Is(err, stores.TypeInvalidKey) {
			return errors.Client(errors.InvalidKey(stores.TypeContractNegotiation,
				`contract negotiation id`, err))
		}
		return errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	if cn.State != negotiation.StateAgreed {
		return errors.Client(errors.StateError(`verify agreement`, string(cn.State)))
	}

	req := negotiation.ContractVerification{
		Ctx:     core.Context,
		Type:    negotiation.MsgTypeAgreementVerification,
		ProvPId: cn.ProvPId,
		ConsPId: consumerPid,
	}

	ack, err := c.send(consumerPid, api.SetParamProviderPid(negotiation.AgreementVerificationEndpoint, cn.ProvPId), req)
	if err != nil {
		return errors.CustomFuncError(`send`, err)
	}

	if errMsg := c.validAck(consumerPid, ack, negotiation.StateVerified); errMsg != `` {
		return errors.Client(errors.InvalidAckError(`VerifyAgreement`, errMsg, ack))
	}

	if err = c.cnStore.UpdateState(consumerPid, negotiation.StateVerified); err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	c.log.Debug(fmt.Sprintf("consumer controller updated negotiation state (id: %s, state: %s)",
		consumerPid, negotiation.StateVerified))
	return nil
}

func (c *Controller) TerminateContract(consumerPid, code string, reasons []string) error {
	cn, err := c.cnStore.Negotiation(consumerPid)
	if err != nil {
		if defaultErr.Is(err, stores.TypeInvalidKey) {
			return errors.Client(errors.InvalidKey(stores.TypeContractNegotiation,
				`contract negotiation id`, err))
		}
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

	ack, err := c.send(consumerPid, api.SetParamPid(negotiation.TerminateEndpoint, cn.ProvPId), req)
	if err != nil {
		return errors.CustomFuncError(`send`, err)
	}

	if errMsg := c.validAck(consumerPid, ack, negotiation.StateTerminated); errMsg != `` {
		return errors.Client(errors.InvalidAckError(`TerminateContract`, errMsg, ack))
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
		return negotiation.Ack{}, errors.Client(errors.MarshalError(``, err))
	}

	res, err := c.client.Send(data, providerAddr+endpoint)
	if err != nil {
		return negotiation.Ack{}, errors.Client(errors.SendFailed(err))
	}

	var ack negotiation.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return negotiation.Ack{}, errors.Client(errors.UnmarshalError(`negotiation ack`, err))
	}

	return ack, nil
}

func (c *Controller) setConstraints(offerId string, vals map[string]string) (odrl.Offer, error) {
	var permList []odrl.Rule
	ofr, err := c.catalog.Offer(offerId)
	if err != nil {
		return odrl.Offer{}, errors.Client(errors.InvalidKey(stores.TypeTransfer, `offer id`, err))
	}

	for _, perm := range ofr.Permissions {
		var consList []odrl.Constraint
		for _, cons := range perm.Constraints {
			val, ok := vals[cons.LeftOperand]
			if !ok {
				return odrl.Offer{}, errors.Client(errors.MissingAttrError(cons.LeftOperand,
					`mandatory constraint`))
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

func (c *Controller) validAck(pid string, ack negotiation.Ack, state negotiation.State) (errMsg string) {
	if ack.Type != negotiation.MsgTypeNegotiationAck {
		return `invalid message type`
	}

	if ack.ConsPId != pid {
		return `incorrect consumer process id`
	}

	if ack.State != state {
		return `invalid state`
	}

	return ``
}
