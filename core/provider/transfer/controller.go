package transfer

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api"
	"github.com/YasiruR/connector/domain/api/dsp/http/transfer"
	"github.com/YasiruR/connector/domain/core"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/stores"
)

type Controller struct {
	tpStore stores.TransferStore
	client  pkg.Client
	log     pkg.Log
}

func NewController(tpStore stores.TransferStore, plugins domain.Plugins) *Controller {
	return &Controller{
		tpStore: tpStore,
		client:  plugins.Client,
		log:     plugins.Log,
	}
}

func (c *Controller) StartTransfer(tpId, sourceEndpoint string) error {
	tp, err := c.tpStore.Process(tpId)
	if err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `Process`, err)
	}

	if tp.State != transfer.StateRequested && tp.State != transfer.StateSuspended {
		return errors.IncompatibleValues(`state`, string(tp.State),
			string(transfer.StateRequested)+" or "+string(transfer.StateSuspended))
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

	if err = c.send(tpId, api.SetParamPid(transfer.StartEndpoint, tp.ConsPId), req); err != nil {
		return errors.CustomFuncError(`send`, err)
	}

	if err = c.tpStore.UpdateState(tpId, transfer.StateStarted); err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `UpdateState`, err)
	}

	c.log.Debug(fmt.Sprintf("provider controller updated transfer process state (id: %s, state: %s)", tpId, transfer.StateStarted))
	return nil
}

func (c *Controller) SuspendTransfer(tpId, code string, reasons []interface{}) error {
	tp, err := c.tpStore.Process(tpId)
	if err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `Process`, err)
	}

	// validate tp

	if tp.State != transfer.StateStarted {
		return errors.IncompatibleValues(`state`, string(tp.State), string(transfer.StateStarted))
	}

	req := transfer.SuspendRequest{
		Ctx:     core.Context,
		Type:    transfer.MsgTypeSuspend,
		ConsPId: tp.ConsPId,
		ProvPId: tpId,
		Code:    code,
		Reason:  reasons,
	}

	if err = c.send(tpId, api.SetParamPid(transfer.SuspendEndpoint, tp.ConsPId), req); err != nil {
		return errors.CustomFuncError(`send`, err)
	}

	if err = c.tpStore.UpdateState(tpId, transfer.StateSuspended); err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `UpdateState`, err)
	}

	c.log.Debug(fmt.Sprintf("provider controller updated transfer process state (id: %s, state: %s)", tpId, transfer.StateSuspended))
	return nil
}

func (c *Controller) CompleteTransfer(tpId string) error {
	tp, err := c.tpStore.Process(tpId)
	if err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `Process`, err)
	}

	// validate tp

	if tp.State != transfer.StateStarted {
		return errors.IncompatibleValues(`state`, string(tp.State), string(transfer.StateStarted))
	}

	req := transfer.CompleteRequest{
		Ctx:     core.Context,
		Type:    transfer.MsgTypeComplete,
		ConsPId: tp.ConsPId,
		ProvPId: tpId,
	}

	if err = c.send(tpId, api.SetParamPid(transfer.CompleteEndpoint, tp.ConsPId), req); err != nil {
		return errors.CustomFuncError(`send`, err)
	}

	if err = c.tpStore.UpdateState(tpId, transfer.StateCompleted); err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `UpdateState`, err)
	}

	c.log.Info(fmt.Sprintf("data exchange process completed successfully (id: %s)", tpId))
	return nil
}

func (c *Controller) send(tpId, endpoint string, request any) error {
	consumerAddr, err := c.tpStore.CallbackAddr(tpId)
	if err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `CallBackAddr`, err)
	}

	data, err := json.Marshal(request)
	if err != nil {
		return errors.MarshalError(endpoint, err)
	}

	res, err := c.client.Send(data, consumerAddr+endpoint)
	if err != nil {
		return errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack transfer.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return errors.UnmarshalError(endpoint, err)
	}

	// validate ack

	return nil
}
