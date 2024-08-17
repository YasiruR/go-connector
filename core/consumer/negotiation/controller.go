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

func (c *Controller) RequestContract(consumerPid, providerEndpoint string, ofr odrl.Offer) (cnId string, err error) {
	var providerPid string
	if consumerPid != `` {
		cn, err := c.cnStore.Negotiation(consumerPid)
		if err != nil {
			return ``, errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
		}

		if cn.State != negotiation.StateOffered {
			return ``, errors.IncompatibleValues(`state`, string(cn.State), string(negotiation.StateOffered))
		}

		providerPid = cn.ProvPId
		c.log.Trace("a contract negotiation already exists for the request", consumerPid)
	} else {
		// generate consumerPid
		consumerPid, err = c.urn.NewURN()
		if err != nil {
			return ``, errors.URNFailed(`consumerPid`, `NewURN`, err)
		}
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

	c.cnStore.Set(consumerPid, negotiation.Negotiation(ack)) // check if received state is REQUESTED (and correctness of other attributes)
	c.cnStore.SetAssignee(consumerPid, ofr.Assignee)
	c.log.Trace(fmt.Sprintf("stored contract negotiation (id: %s, assigner: %s, assignee: %s)", consumerPid, ofr.Assigner, ofr.Assignee))
	c.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", consumerPid, negotiation.StateRequested))
	return consumerPid, nil
}

func (c *Controller) AcceptContract() {}

func (c *Controller) VerifyAgreement(consumerPid string) error {
	neg, err := c.cnStore.Negotiation(consumerPid)
	if err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `Negotiation`, err)
	}

	req := negotiation.ContractVerification{
		Ctx:     core.Context,
		Type:    negotiation.MsgTypeAgreementVerification,
		ProvPId: neg.ProvPId,
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

	endpoint := strings.Replace(providerAddr+negotiation.AgreementVerificationEndpoint, `{`+negotiation.ParamProviderId+`}`, neg.ProvPId, 1)
	res, err := c.client.Send(data, endpoint)
	if err != nil {
		return errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack negotiation.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return errors.UnmarshalError(providerAddr+negotiation.AgreementVerificationEndpoint, err)
	}

	if err = c.cnStore.UpdateState(consumerPid, negotiation.StateVerified); err != nil {
		return errors.StoreFailed(stores.TypeContractNegotiation, `UpdateState`, err)
	}

	c.log.Info(fmt.Sprintf("updated negotiation state (id: %s, state: %s)", consumerPid, negotiation.StateVerified))
	return nil
}

func (c *Controller) TerminateContract() {}
