package transfer

import (
	"encoding/json"
	"fmt"
	transfer2 "github.com/YasiruR/connector/domain/api/dsp/http/transfer"
	"github.com/YasiruR/connector/domain/core"
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

func (c *Controller) RequestTransfer(typ transfer2.DataTransferType, agreementId, sinkEndpoint, providerEndpoint string) (tpId string, err error) {
	tpId, err = c.urn.NewURN()
	if err != nil {
		return ``, errors.PkgFailed(pkg.TypeURN, `New`, err)
	}

	var addr transfer2.Address
	if typ == transfer2.HTTPPush {
		addr = transfer2.Address{
			Type:               transfer2.TypeDataAddress,
			EndpointType:       transfer2.EndpointTypeHTTP,
			Endpoint:           sinkEndpoint, // validate - cannot be null if push
			EndpointProperties: nil,          // e.g. auth tokens
		}
	}

	req := transfer2.Request{
		Ctx:          core.Context,
		Type:         transfer2.TypeTransferRequest,
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

	res, err := c.client.Send(data, providerEndpoint+transfer2.RequestEndpoint)
	if err != nil {
		return ``, errors.PkgFailed(pkg.TypeClient, `Send`, err)
	}

	var ack transfer2.Ack
	if err = json.Unmarshal(res, &ack); err != nil {
		return ``, errors.UnmarshalError(providerEndpoint+transfer2.RequestEndpoint, err)
	}

	c.trStore.Set(tpId, transfer2.Process(ack)) // validate if received attributes are correct
	c.log.Trace("stored transfer process", ack)
	c.log.Info(fmt.Sprintf("updated transfer process state (id: %s, state: %s)", tpId, transfer2.StateRequested))
	return tpId, nil
}
