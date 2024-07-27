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

	req := transfer.StartRequest{
		Ctx:     core.Context,
		Type:    transfer.TypeTransferStart,
		ConsPId: tp.ConsPId,
		ProvPId: tpId,
	}

	if tp.Type == transfer.HTTPPull {
		if sourceEndpoint == `` {
			return errors.MissingRequiredAttr(`sourceEndpoint`, `mandatory for pull transfers`)
		}

		req.Address = transfer.Address{
			Type:               transfer.TypeDataAddress,
			EndpointType:       transfer.EndpointTypeHTTP,
			Endpoint:           sourceEndpoint,
			EndpointProperties: nil, // e.g. auth tokens
		}
	}

	data, err := json.Marshal(req)
	if err != nil {
		return errors.MarshalError(transfer.StartTransferEndpoint, err)
	}

	consumerAddr, err := c.tpStore.CallbackAddr(tpId)
	if err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `CallBackAddr`, err)
	}

	endpoint := strings.Replace(consumerAddr+transfer.StartTransferEndpoint, `{`+transfer.ParamConsumerPid+`}`, tp.ConsPId, 1)
	res, err := c.client.Send(data, endpoint)
	if err != nil {
		return errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack transfer.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return errors.UnmarshalError(transfer.StartTransferEndpoint, err)
	}

	if err = c.tpStore.UpdateState(tpId, transfer.StateStarted); err != nil {
		return errors.StoreFailed(stores.TypeTransfer, `UpdateState`, err)
	}

	c.log.Info(fmt.Sprintf("updated transfer process state (id: %s, state: %s)", tpId, transfer.StateStarted))
	return nil
}
