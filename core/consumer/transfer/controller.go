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
)

type Controller struct {
	callbackAddr string
	urn          pkg.URNService
	client       pkg.Client
	trStore      stores.Transfer
	log          pkg.Log
}

func NewController(port int, stores domain.Stores, plugins domain.Plugins) *Controller {
	return &Controller{
		callbackAddr: `http://localhost:` + strconv.Itoa(port),
		trStore:      stores.Transfer,
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
			Type:               transfer.TypeDataAddress,
			EndpointType:       transfer.EndpointTypeHTTP,
			Endpoint:           sinkEndpoint,
			EndpointProperties: nil, // e.g. auth tokens
		}
	}

	req := transfer.Request{
		Ctx:          core.Context,
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
