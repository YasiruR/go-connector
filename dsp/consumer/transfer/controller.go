package transfer

import (
	"encoding/json"
	"fmt"
	"github.com/YasiruR/connector/domain/dsp"
	"github.com/YasiruR/connector/domain/dsp/transfer"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/stores"
)

type Controller struct {
	callbackAddr string
	urn          pkg.URNService
	client       pkg.Client
	trStore      stores.Transfer
	log          pkg.Log
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) RequestTransfer(typ transfer.DataTransferType, agreementId, sinkEndpoint, providerEndpoint string) (tpId string, err error) {
	tpId, err = c.urn.NewURN()
	if err != nil {
		return ``, errors.PkgFailed(pkg.TypeURN, `New`, err)
	}

	var addr transfer.Address
	if typ == transfer.HTTPPush {
		addr = transfer.Address{
			Type:               transfer.TypeDataAddress,
			EndpointType:       transfer.EndpointTypeHTTP,
			Endpoint:           sinkEndpoint, // validate - cannot be null if push
			EndpointProperties: nil,          // e.g. auth tokens
		}
	}

	req := transfer.Request{
		Ctx:          dsp.Context,
		Type:         transfer.TypeTransferRequest,
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

	c.trStore.Set(tpId, transfer.Process(ack)) // validate if received attributes are correct
	c.log.Trace("stored transfer process", ack)
	c.log.Info(fmt.Sprintf("updated transfer process state (id: %s, state: %s)", tpId, transfer.StateRequested))
	return tpId, nil
}
