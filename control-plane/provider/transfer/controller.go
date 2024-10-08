package transfer

import (
	"encoding/json"
	defaultErr "errors"
	"fmt"
	"github.com/YasiruR/go-connector/domain"
	"github.com/YasiruR/go-connector/domain/api"
	"github.com/YasiruR/go-connector/domain/api/dsp/http/transfer"
	"github.com/YasiruR/go-connector/domain/control-plane"
	"github.com/YasiruR/go-connector/domain/errors"
	"github.com/YasiruR/go-connector/domain/models"
	"github.com/YasiruR/go-connector/domain/pkg"
	"github.com/YasiruR/go-connector/domain/stores"
	"strings"
)

type Controller struct {
	tpStore  stores.TransferStore
	agrStore stores.AgreementStore
	client   pkg.Client
	log      pkg.Log
}

func NewController(s domain.Stores, plugins domain.Plugins) *Controller {
	return &Controller{
		tpStore:  s.TransferStore,
		agrStore: s.AgreementStore,
		client:   plugins.Client,
		log:      plugins.Log,
	}
}

func (c *Controller) StartTransfer(tpId string, sourceDb models.Database) error {
	tp, err := c.tpStore.Process(tpId)
	if err != nil {
		if defaultErr.Is(err, stores.TypeInvalidKey) {
			return errors.Client(errors.InvalidKey(stores.TypeTransfer, `transfer process id`, err))
		}
		return errors.StoreFailed(stores.TypeTransfer, `Process`, err)
	}

	if tp.State != transfer.StateRequested && tp.State != transfer.StateSuspended {
		return errors.Client(errors.StateError(`start transfer`, string(tp.State)))
	}

	req := transfer.StartRequest{
		Ctx:     control_plane.Context,
		Type:    transfer.MsgTypeStart,
		ConsPId: tp.ConsPId,
		ProvPId: tpId,
	}

	if strings.HasSuffix(string(tp.Type), `+pull`) {
		if sourceDb.Endpoint == `` {
			return errors.Client(errors.MissingAttrError(`source endpoint`,
				`mandatory for pull transfers`))
		}
		req.Address = transfer.Address{
			Type:         transfer.MsgTypeDataAddress,
			EndpointType: transfer.EndpointTypeHTTP,
			Endpoint:     sourceDb.Endpoint,
			EndpointProperties: []transfer.EndpointProperty{
				{
					Type:  transfer.MsgTypeEndpointProperty,
					Name:  transfer.PropertyUsername,
					Value: sourceDb.User,
				},
				{
					Type:  transfer.MsgTypeEndpointProperty,
					Name:  transfer.PropertyPassword,
					Value: sourceDb.Password,
				},
				{
					Type:  transfer.MsgTypeEndpointProperty,
					Name:  transfer.PropertyDatabase,
					Value: sourceDb.Name,
				},
			}, // e.g. auth tokens
		}
	}

	if err = c.send(tpId, api.SetParamPid(transfer.StartEndpoint, tp.ConsPId), req); err != nil {
		return errors.CustomFuncError(`send`, err)
	}

	if err = c.tpStore.UpdateState(tpId, transfer.StateStarted); err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `UpdateState`, err)
	}

	c.log.Debug(fmt.Sprintf("provider controller updated transfer process state (id: %s, state: %s)",
		tpId, transfer.StateStarted))
	return nil
}

func (c *Controller) SuspendTransfer(tpId, code string, reasons []interface{}) error {
	tp, err := c.tpStore.Process(tpId)
	if err != nil {
		if defaultErr.Is(err, stores.TypeInvalidKey) {
			return errors.Client(errors.InvalidKey(stores.TypeTransfer, `transfer process id`, err))
		}
		return errors.StoreFailed(stores.TypeTransfer, `Process`, err)
	}

	// validate tp

	if tp.State != transfer.StateStarted {
		return errors.Client(errors.StateError(`suspend transfer`, string(tp.State)))
	}

	req := transfer.SuspendRequest{
		Ctx:     control_plane.Context,
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

	c.log.Debug(fmt.Sprintf("provider controller updated transfer process state (id: %s, state: %s)",
		tpId, transfer.StateSuspended))
	return nil
}

func (c *Controller) CompleteTransfer(tpId string) error {
	tp, err := c.tpStore.Process(tpId)
	if err != nil {
		if defaultErr.Is(err, stores.TypeInvalidKey) {
			return errors.Client(errors.InvalidKey(stores.TypeTransfer, `transfer process id`, err))
		}
		return errors.StoreFailed(stores.TypeTransfer, `Process`, err)
	}

	// validate tp

	if tp.State != transfer.StateStarted {
		return errors.Client(errors.StateError(`complete transfer`, string(tp.State)))
	}

	req := transfer.CompleteRequest{
		Ctx:     control_plane.Context,
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

	c.log.Info(fmt.Sprintf("data-plane exchange process completed successfully (id: %s)", tpId))
	return nil
}

func (c *Controller) TerminateTransfer(tpId, code string, reasons []interface{}) error {
	tp, err := c.tpStore.Process(tpId)
	if err != nil {
		if defaultErr.Is(err, stores.TypeInvalidKey) {
			return errors.Client(errors.InvalidKey(stores.TypeTransfer, `transfer process id`, err))
		}
		return errors.StoreFailed(stores.TypeTransfer, `Process`, err)
	}

	// validate tp

	if tp.State != transfer.StateRequested && tp.State != transfer.StateStarted && tp.State != transfer.StateSuspended {
		return errors.Client(errors.StateError(`terminate transfer`, string(tp.State)))
	}

	req := transfer.TerminateRequest{
		Ctx:     control_plane.Context,
		Type:    transfer.MsgTypeTerminate,
		ConsPId: tpId,
		ProvPId: tp.ProvPId,
		Code:    code,
		Reason:  reasons,
	}

	if err = c.send(tpId, api.SetParamPid(transfer.TerminateEndpoint, tp.ConsPId), req); err != nil {
		return errors.CustomFuncError(`send`, err)
	}

	if err = c.tpStore.UpdateState(tpId, transfer.StateTerminated); err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `UpdateState`, err)
	}

	c.log.Info(fmt.Sprintf("terminated data-plane exchange process (id: %s)", tpId))
	return nil
}

func (c *Controller) send(tpId, endpoint string, request any) error {
	consumerAddr, err := c.tpStore.CallbackAddr(tpId)
	if err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `CallBackAddr`, err)
	}

	body, err := json.Marshal(request)
	if err != nil {
		return errors.Client(errors.MarshalError(``, err))
	}

	res, err := c.client.Send(body, consumerAddr+endpoint)
	if err != nil {
		return errors.Client(errors.SendFailed(err))
	}

	var ack transfer.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return errors.Client(errors.UnmarshalError(`transfer ack`, err))
	}

	// validate ack

	return nil
}

//func (c *Controller) transfer(tpId string) error {
//	agrId, err := c.tpStore.AgreementId(tpId)
//	if err != nil {
//		return errors.StoreFailed(stores.TypeTransfer, `AgreementId`, err)
//	}
//
//	agr, err := c.agrStore.Agreement(agrId)
//	if err != nil {
//		return errors.StoreFailed(stores.TypeAgreement, `Agreement`, err)
//	}
//
//	addr, err := c.tpStore.DataSinkAddress(tpId)
//	if err != nil {
//		return errors.StoreFailed(stores.TypeTransfer, `DataSinkAddress`, err)
//	}
//
//	var db string
//	var cr data-plane.Credentials
//	for _, p := range addr.EndpointProperties {
//		switch p.Name {
//		case transfer.PropertyUsername:
//			cr.User = p.Value
//		case transfer.PropertyPassword:
//			cr.Password = p.Value
//		case transfer.PropertyDatabase:
//			db = p.Value
//		}
//	}
//
//	if err = c.exchanger.PushWithCredentials(string(agr.Target),
//		data-plane.Database{Name: db, Endpoint: addr.Endpoint, Credentials: cr}); err != nil {
//		return errors.TransferFailed(addr.Endpoint, err)
//	}
//
//	return nil
//}
