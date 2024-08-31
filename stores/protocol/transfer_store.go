package protocol

import (
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/dsp/http/transfer"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/pkg"
)

const (
	transferCollection = `transfer`
)

type Transfer struct {
	store        pkg.Collection
	callbackAddr pkg.Collection
}

func NewTransferStore(plugins domain.Plugins) *Transfer {
	plugins.Log.Info("initialized transfer process store")
	return &Transfer{store: plugins.Database.NewCollection(), callbackAddr: plugins.Database.NewCollection()}
}

func (t *Transfer) AddProcess(tpId string, val transfer.Process) {
	_ = t.store.Set(tpId, val)
}

func (t *Transfer) Process(id string) (transfer.Process, error) {
	val, err := t.store.Get(id)
	if err != nil {
		return transfer.Process{}, errors.QueryFailed(transferCollection, `Get`, err)
	}
	return val.(transfer.Process), nil
}

func (t *Transfer) SetCallbackAddr(tpId, addr string) {
	_ = t.callbackAddr.Set(tpId, addr)
}

func (t *Transfer) CallbackAddr(tpId string) (string, error) {
	val, err := t.callbackAddr.Get(tpId)
	if err != nil {
		return ``, errors.QueryFailed(callbackAddrCollection, `Get`, err)
	}
	return val.(string), nil
}

func (t *Transfer) UpdateState(tpId string, s transfer.State) error {
	process, err := t.Process(tpId)
	if err != nil {
		return errors.QueryFailed(transferCollection, `Process`, err)
	}

	process.State = s
	t.AddProcess(tpId, process)
	return nil
}
