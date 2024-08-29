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
	"strconv"
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

func (c *Controller) GetProviderProcess(tpId string) (transfer.Process, error) {
	tp, err := c.tpStore.Process(tpId)
	if err != nil {
		return transfer.Process{}, errors.StoreFailed(stores.TypeTransfer, `Process`, err)
	}

	ack, err := c.send(tpId, api.SetParamProviderPid(transfer.GetProcessEndpoint, tp.ProvPId), nil)
	if err != nil {
		return transfer.Process{}, errors.CustomFuncError(`send`, err)
	}

	return transfer.Process(ack), nil
}

func (c *Controller) RequestTransfer(dataFormat, agreementId, sinkEndpoint,
	providerEndpoint string) (tpId string, err error) {
	// include validations for format
	typ := transfer.DataTransferType(dataFormat)

	tpId, err = c.urn.NewURN()
	if err != nil {
		return ``, errors.PkgFailed(pkg.TypeURN, `New`, err)
	}

	req := transfer.Request{
		Ctx:          core.Context,
		Type:         transfer.MsgTypeRequest,
		ConsPId:      tpId,
		AgreementId:  agreementId,
		Format:       typ,
		CallbackAddr: c.callbackAddr,
	}

	if typ == transfer.HTTPPush {
		if sinkEndpoint == `` {
			return ``, errors.MissingRequiredAttr(`sinkEndpoint`, `mandatory for push transfers`)
		}

		req.Address = transfer.Address{
			Type:               transfer.MsgTypeDataAddress,
			EndpointType:       transfer.EndpointTypeHTTP,
			Endpoint:           sinkEndpoint,
			EndpointProperties: nil, // e.g. auth tokens
		}
	}

	c.tpStore.SetCallbackAddr(tpId, providerEndpoint)
	ack, err := c.send(tpId, transfer.RequestEndpoint, req)
	if err != nil {
		return ``, errors.CustomFuncError(`send`, err)
	}

	c.tpStore.AddProcess(tpId, transfer.Process(ack))
	c.log.Trace("stored transfer process", ack)
	c.log.Debug(fmt.Sprintf("consumer controller updated transfer process state (id: %s, state: %s)",
		tpId, transfer.StateRequested))
	return tpId, nil
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
		ConsPId: tpId,
		ProvPId: tp.ProvPId,
		Code:    code,
		Reason:  reasons,
	}

	if _, err = c.send(tpId, api.SetParamPid(transfer.SuspendEndpoint, tp.ProvPId), req); err != nil {
		return errors.CustomFuncError(`send`, err)
	}

	if err = c.tpStore.UpdateState(tpId, transfer.StateSuspended); err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `UpdateState`, err)
	}

	c.log.Debug(fmt.Sprintf("consumer controller updated transfer process state (id: %s, state: %s)",
		tpId, transfer.StateSuspended))
	return nil
}

func (c *Controller) StartTransfer(tpId string) error {
	// todo handler for start in provider
	tp, err := c.tpStore.Process(tpId)
	if err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `Process`, err)
	}

	// validate tp

	if tp.State != transfer.StateSuspended {
		return errors.IncompatibleValues(`state`, string(tp.State), string(transfer.StateSuspended))
	}

	req := transfer.StartRequest{
		Ctx:     core.Context,
		Type:    transfer.MsgTypeStart,
		ConsPId: tpId,
		ProvPId: tp.ProvPId,
	}
	// check if it needs to provide the data address as well

	if _, err = c.send(tpId, api.SetParamPid(transfer.StartEndpoint, tp.ProvPId), req); err != nil {
		return errors.CustomFuncError(`send`, err)
	}

	if err = c.tpStore.UpdateState(tpId, transfer.StateStarted); err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `UpdateState`, err)
	}

	c.log.Debug(fmt.Sprintf("consumer controller updated state of the suspended transfer process (id: %s, state: %s)",
		tpId, transfer.StateStarted))
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
		ConsPId: tpId,
		ProvPId: tp.ProvPId,
	}

	if _, err = c.send(tpId, api.SetParamPid(transfer.CompleteEndpoint, tp.ProvPId), req); err != nil {
		return errors.CustomFuncError(`send`, err)
	}

	if err = c.tpStore.UpdateState(tpId, transfer.StateCompleted); err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `UpdateState`, err)
	}

	c.log.Info(fmt.Sprintf("data exchange process completed successfully (id: %s)", tpId))
	return nil
}

func (c *Controller) TerminateTransfer(tpId string, code string, reasons []interface{}) error {
	tp, err := c.tpStore.Process(tpId)
	if err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `Process`, err)
	}

	if tp.State != transfer.StateRequested && tp.State != transfer.StateStarted && tp.State != transfer.StateSuspended {
		return errors.IncompatibleValues(`state`, string(tp.State),
			string(transfer.StateStarted)+" or "+string(transfer.StateStarted)+" or "+string(transfer.StateSuspended))
	}

	req := transfer.TerminateRequest{
		Ctx:     core.Context,
		Type:    transfer.MsgTypeTerminate,
		ConsPId: tpId,
		ProvPId: tp.ProvPId,
		Code:    code,
		Reason:  reasons,
	}

	if _, err = c.send(tpId, api.SetParamPid(transfer.TerminateEndpoint, tp.ProvPId), req); err != nil {
		return errors.CustomFuncError(`send`, err)
	}

	if err = c.tpStore.UpdateState(tpId, transfer.StateTerminated); err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `UpdateState`, err)
	}

	c.log.Info(fmt.Sprintf("terminated data exchange process (id: %s)", tpId))
	return nil
}

func (c *Controller) send(tpId, endpoint string, req any) (transfer.Ack, error) {
	providerAddr, err := c.tpStore.CallbackAddr(tpId)
	if err != nil {
		return transfer.Ack{}, errors.StoreFailed(stores.TypeTransfer, `CallBackAddr`, err)
	}

	var data []byte
	if req != nil {
		data, err = json.Marshal(req)
		if err != nil {
			return transfer.Ack{}, errors.MarshalError(transfer.CompleteEndpoint, err)
		}
	}

	res, err := c.client.Send(data, providerAddr+endpoint)
	if err != nil {
		return transfer.Ack{}, errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack transfer.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return transfer.Ack{}, errors.UnmarshalError(transfer.CompleteEndpoint, err)
	}

	// validate ack

	return ack, nil
}
