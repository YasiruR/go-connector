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
)

type Controller struct {
	callbackAddr string
	cnStore      stores.ContractNegotiation
	urn          pkg.URNService
	client       pkg.Client
	log          pkg.Log
}

func NewController(port int, stores domain.Stores, plugins domain.Plugins) *Controller {
	return &Controller{
		callbackAddr: `http://localhost:` + strconv.Itoa(port),
		cnStore:      stores.ContractNegotiation,
		urn:          plugins.URNService,
		client:       plugins.Client,
		log:          plugins.Log,
	}
}

func (c *Controller) RequestContract(consumerPid, providerAddr string, ofr odrl.Offer) (cnId string, err error) {
	var providerPid string
	var endpoint string
	if consumerPid != `` {
		cn, err := c.cnStore.GetNegotiation(consumerPid)
		if err != nil {
			return ``, errors.StoreFailed(stores.TypeContractNegotiation, `GetNegotiation`, err)
		}

		if cn.State != negotiation.StateOffered {
			return ``, errors.IncompatibleValues(`state`, string(cn.State), string(negotiation.StateOffered))
		}

		providerPid = cn.ProvPId
		endpoint = strings.Replace(negotiation.ContractRequestToOfferEndpoint, `{`+negotiation.ParamProviderId+`}`, cn.ProvPId, 1)
		c.log.Trace("found an existing contract negotiation for the request", "id: "+consumerPid)
	} else {
		// generate consumerPid
		consumerPid, err = c.urn.NewURN()
		if err != nil {
			return ``, errors.URNFailed(`consumerPid`, `NewURN`, err)
		}
		endpoint = negotiation.ContractRequestEndpoint
	}

	// construct payload
	req := negotiation.ContractRequest{
		Ctx:          core.Context,
		Type:         negotiation.MsgTypeContractRequest,
		ConsPId:      consumerPid,
		ProvPId:      providerPid,
		Offer:        ofr,
		CallbackAddr: c.callbackAddr,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return ``, errors.MarshalError(endpoint, err)
	}

	res, err := c.client.Send(data, providerAddr+endpoint)
	if err != nil {
		return ``, errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack negotiation.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return ``, errors.UnmarshalError(providerAddr+negotiation.ContractRequestEndpoint, err)
	}

	if !c.validAck(consumerPid, ack, negotiation.StateRequested) {
		return ``, errors.InvalidAck(`ContractRequest`, ack)
	}

	ack.Type = negotiation.MsgTypeNegotiation
	c.cnStore.Set(consumerPid, negotiation.Negotiation(ack))
	c.cnStore.SetAssignee(consumerPid, ofr.Assignee)
	c.cnStore.SetCallbackAddr(consumerPid, providerAddr)

	c.log.Trace("stored contract negotiation", "id: "+consumerPid, "assigner: "+ofr.Assigner,
		"assignee: "+ofr.Assignee, "address: "+providerAddr)
	c.log.Debug("updated negotiation state", "id: "+consumerPid, "state: "+negotiation.StateRequested)
	return consumerPid, nil
}

func (c *Controller) AcceptOffer(consumerPid string) error {
	cn, err := c.cnStore.GetNegotiation(consumerPid)
	if err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `GetNegotiation`, err)
	}

	if cn.State != negotiation.StateOffered {
		return errors.IncompatibleValues(`state`, string(cn.State), string(negotiation.StateOffered))
	}

	event := negotiation.ContractNegotiationEvent{
		Ctx:       core.Context,
		Type:      negotiation.MsgTypeNegotiationEvent,
		ProvPId:   cn.ProvPId,
		ConsPId:   cn.ConsPId,
		EventType: negotiation.EventAccepted,
	}

	provAddr, err := c.cnStore.CallbackAddr(consumerPid)
	if err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `CallbackAddr`, err)
	}

	data, err := json.Marshal(event)
	if err != nil {
		return errors.MarshalError(negotiation.EventsEndpoint, err)
	}

	res, err := c.client.Send(data, provAddr+negotiation.EventsEndpoint)
	if err != nil {
		return errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack negotiation.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return errors.UnmarshalError(negotiation.EventsEndpoint, err)
	}

	if err = c.cnStore.UpdateState(consumerPid, negotiation.StateAccepted); err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	c.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", consumerPid, negotiation.StateAccepted))
	return nil
}

func (c *Controller) VerifyAgreement(consumerPid string) error {
	cn, err := c.cnStore.GetNegotiation(consumerPid)
	if err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `GetNegotiation`, err)
	}

	req := negotiation.ContractVerification{
		Ctx:     core.Context,
		Type:    negotiation.MsgTypeAgreementVerification,
		ProvPId: cn.ProvPId,
		ConsPId: consumerPid,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return errors.MarshalError(``, err)
	}

	providerAddr, err := c.cnStore.CallbackAddr(consumerPid)
	if err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `CallBackAddr`, err)
	}

	endpoint := strings.Replace(negotiation.AgreementVerificationEndpoint, `{`+negotiation.ParamProviderId+`}`, cn.ProvPId, 1)
	res, err := c.client.Send(data, providerAddr+endpoint)
	if err != nil {
		return errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack negotiation.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return errors.UnmarshalError(negotiation.AgreementVerificationEndpoint, err)
	}

	if err = c.cnStore.UpdateState(consumerPid, negotiation.StateVerified); err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	c.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", consumerPid, negotiation.StateVerified))
	return nil
}

func (c *Controller) TerminateContract(consumerPid, code string, reasons []string) error {
	cn, err := c.cnStore.GetNegotiation(consumerPid)
	if err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `GetNegotiation`, err)
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

	providerAddr, err := c.cnStore.CallbackAddr(consumerPid)
	if err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `CallBackAddr`, err)
	}

	data, err := json.Marshal(req)
	if err != nil {
		return errors.MarshalError(negotiation.TerminateEndpoint, err)
	}

	endpoint := strings.Replace(negotiation.TerminateEndpoint, `{`+negotiation.ParamContractId+`}`, cn.ProvPId, 1)
	res, err := c.client.Send(data, providerAddr+endpoint)
	if err != nil {
		return errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack negotiation.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return errors.UnmarshalError(negotiation.TerminateEndpoint, err)
	}

	// clear all store entries for the contract negotiation
	if err = c.cnStore.UpdateState(consumerPid, negotiation.StateTerminated); err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	c.log.Info("consumer terminated the negotiation flow", consumerPid)
	return nil
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
