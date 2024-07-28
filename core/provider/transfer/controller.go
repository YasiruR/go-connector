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
	"strings"
)

type Controller struct {
	tpStore stores.Transfer
	client  pkg.Client
	log     pkg.Log
}

func NewController(tpStore stores.Transfer, plugins domain.Plugins) *Controller {
	return &Controller{
		tpStore: tpStore,
		client:  plugins.Client,
		log:     plugins.Log,
	}
}

func (c *Controller) StartTransfer(tpId, sourceEndpoint string) error {
	tp, err := c.tpStore.GetProcess(tpId)
	if err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `GetProcess`, err)
	}

	consumerAddr, err := c.tpStore.CallbackAddr(tpId)
	if err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `CallBackAddr`, err)
	}

	req := transfer.StartRequest{
		Ctx:     core.Context,
		Type:    transfer.MsgTypeStart,
		ConsPId: tp.ConsPId,
		ProvPId: tpId,
	}

	if tp.Type == transfer.HTTPPull {
		if sourceEndpoint == `` {
			return errors.MissingRequiredAttr(`sourceEndpoint`, `mandatory for pull transfers`)
		}

		req.Address = transfer.Address{
			Type:               transfer.MsgTypeDataAddress,
			EndpointType:       transfer.EndpointTypeHTTP,
			Endpoint:           sourceEndpoint,
			EndpointProperties: nil, // e.g. auth tokens
		}
	}

	data, err := json.Marshal(req)
	if err != nil {
		return errors.MarshalError(transfer.StartEndpoint, err)
	}

	endpoint := strings.Replace(consumerAddr+transfer.StartEndpoint, `{`+transfer.ParamConsumerPid+`}`, tp.ConsPId, 1)
	res, err := c.client.Send(data, endpoint)
	if err != nil {
		return errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack transfer.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return errors.UnmarshalError(transfer.StartEndpoint, err)
	}

	if err = c.tpStore.UpdateState(tpId, transfer.StateStarted); err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `UpdateState`, err)
	}

	c.log.Info(fmt.Sprintf("updated transfer process state (id: %s, state: %s)", tpId, transfer.StateStarted))
	return nil
}

func (c *Controller) SuspendTransfer(tpId, code string, reasons []interface{}) error {
	// check if valid tp
	tp, err := c.tpStore.GetProcess(tpId)
	if err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `GetProcess`, err)
	}

	consumerAddr, err := c.tpStore.CallbackAddr(tpId)
	if err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `CallBackAddr`, err)
	}

	req := transfer.SuspendRequest{
		Ctx:     core.Context,
		Type:    transfer.MsgTypeSuspend,
		ConsPId: tp.ConsPId,
		ProvPId: tpId,
		Code:    code,
		Reason:  reasons,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return errors.MarshalError(transfer.SuspendEndpoint, err)
	}

	endpoint := strings.Replace(consumerAddr+transfer.SuspendEndpoint, `{`+transfer.ParamPid+`}`, tp.ConsPId, 1)
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
	tp, err := c.tpStore.GetProcess(tpId)
	if err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `GetProcess`, err)
	}

	consumerAddr, err := c.tpStore.CallbackAddr(tpId)
	if err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `CallBackAddr`, err)
	}

	req := transfer.CompleteRequest{
		Ctx:     core.Context,
		Type:    transfer.MsgTypeComplete,
		ConsPId: tp.ConsPId,
		ProvPId: tpId,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return errors.MarshalError(transfer.CompleteEndpoint, err)
	}

	endpoint := strings.Replace(consumerAddr+transfer.CompleteEndpoint, `{`+transfer.ParamPid+`}`, tp.ConsPId, 1)
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
