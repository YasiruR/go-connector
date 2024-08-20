package transfer

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/dsp/http/transfer"
	"github.com/YasiruR/connector/domain/core"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/stores"
	"strconv"
	"strings"
)

type Controller struct {
	callbackAddr string
	urn          pkg.URNService
	client       pkg.Client
	tpStore      stores.TransferStore
	log          pkg.Log
}

func NewController(port int, stores domain.Stores, plugins domain.Plugins) *Controller {
	return &Controller{
		callbackAddr: `http://localhost:` + strconv.Itoa(port),
		tpStore:      stores.TransferStore,
		client:       plugins.Client,
		urn:          plugins.URNService,
		log:          plugins.Log,
	}
}

func (c *Controller) RequestTransfer(dataFormat, agreementId, sinkEndpoint, providerEndpoint string) (tpId string, err error) {
	// include validations for format
	typ := transfer.DataTransferType(dataFormat)

	tpId, err = c.urn.NewURN()
	if err != nil {
		return ``, errors.PkgFailed(pkg.TypeURN, `New`, err)
	}

	var addr transfer.Address
	if typ == transfer.HTTPPush {
		if sinkEndpoint == `` {
			return ``, errors.MissingRequiredAttr(`sinkEndpoint`, `mandatory for push transfers`)
		}

		addr = transfer.Address{
			Type:               transfer.MsgTypeDataAddress,
			EndpointType:       transfer.EndpointTypeHTTP,
			Endpoint:           sinkEndpoint,
			EndpointProperties: nil, // e.g. auth tokens
		}
	}

	req := transfer.Request{
		Ctx:          core.Context,
		Type:         transfer.MsgTypeRequest,
		ConsPId:      tpId,
		AgreementId:  agreementId,
		Format:       typ,
		Address:      addr,
		CallbackAddr: c.callbackAddr,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return ``, errors.MarshalError(``, err)
	}

	res, err := c.client.Send(data, providerEndpoint+transfer.RequestEndpoint)
	if err != nil {
		return ``, errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack transfer.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return ``, errors.UnmarshalError(providerEndpoint+transfer.RequestEndpoint, err)
	}

	c.tpStore.AddProcess(tpId, transfer.Process(ack)) // validate if received attributes are correct
	c.tpStore.SetCallbackAddr(tpId, providerEndpoint)
	c.log.Trace("stored transfer process", ack)
	c.log.Info(fmt.Sprintf("updated transfer process state (id: %s, state: %s)", tpId, transfer.StateRequested))
	return tpId, nil
}

func (c *Controller) SuspendTransfer(tpId, code string, reasons []interface{}) error {
	// check if valid tp
	tp, err := c.tpStore.Process(tpId)
	if err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `Process`, err)
	}

	providerAddr, err := c.tpStore.CallbackAddr(tpId)
	if err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `CallBackAddr`, err)
	}

	req := transfer.SuspendRequest{
		Ctx:     core.Context,
		Type:    transfer.MsgTypeSuspend,
		ConsPId: tpId,
		ProvPId: tp.ProvPId,
		Code:    code,
		Reason:  reasons,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return errors.MarshalError(transfer.SuspendEndpoint, err)
	}

	endpoint := strings.Replace(providerAddr+transfer.SuspendEndpoint, `{`+transfer.ParamPid+`}`, tp.ProvPId, 1)
	res, err := c.client.Send(data, endpoint)
	if err != nil {
		return errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack transfer.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return errors.UnmarshalError(transfer.SuspendEndpoint, err)
	}

	if err = c.tpStore.UpdateState(tpId, transfer.StateSuspended); err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `UpdateState`, err)
	}

	c.log.Info(fmt.Sprintf("updated transfer process state (id: %s, state: %s)", tpId, transfer.StateSuspended))
	return nil
}

func (c *Controller) CompleteTransfer(tpId string) error {
	tp, err := c.tpStore.Process(tpId)
	if err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `Process`, err)
	}

	providerAddr, err := c.tpStore.CallbackAddr(tpId)
	if err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `CallBackAddr`, err)
	}

	req := transfer.CompleteRequest{
		Ctx:     core.Context,
		Type:    transfer.MsgTypeComplete,
		ConsPId: tpId,
		ProvPId: tp.ProvPId,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return errors.MarshalError(transfer.CompleteEndpoint, err)
	}

	endpoint := strings.Replace(providerAddr+transfer.CompleteEndpoint, `{`+transfer.ParamPid+`}`, tp.ProvPId, 1)
	res, err := c.client.Send(data, endpoint)
	if err != nil {
		return errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack transfer.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return errors.UnmarshalError(transfer.CompleteEndpoint, err)
	}

	if err = c.tpStore.UpdateState(tpId, transfer.StateCompleted); err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `UpdateState`, err)
	}

	c.log.Info(fmt.Sprintf("updated transfer process state (id: %s, state: %s)", tpId, transfer.StateCompleted))
	return nil
}
