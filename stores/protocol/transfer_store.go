package protocol

import (
	"github.com/YasiruR/go-connector/domain"
	"github.com/YasiruR/go-connector/domain/api/dsp/http/transfer"
	"github.com/YasiruR/go-connector/domain/pkg"
	"github.com/YasiruR/go-connector/domain/stores"
)

const (
	collTransfer   = `transfer`
	collAgreements = `agreements`
	collDataAddr   = `data-sinks`
)

type Transfer struct {
	processes    pkg.Collection
	callbackAddr pkg.Collection
	agreements   pkg.Collection
	dataAddr     pkg.Collection
}

func NewTransferStore(plugins domain.Plugins) *Transfer {
	plugins.Log.Info("initialized transfer process store")
	return &Transfer{
		processes:    plugins.Store.NewCollection(),
		callbackAddr: plugins.Store.NewCollection(),
		agreements:   plugins.Store.NewCollection(),
		dataAddr:     plugins.Store.NewCollection(),
	}
}

func (t *Transfer) AddProcess(tpId string, val transfer.Process) {
	_ = t.processes.Set(tpId, val)
}

func (t *Transfer) Process(id string) (transfer.Process, error) {
	val, err := t.processes.Get(id)
	if err != nil {
		return transfer.Process{}, stores.QueryFailed(collTransfer, `Get`, err)
	}

	if val == nil {
		return transfer.Process{}, stores.InvalidKey(id)
	}

	return val.(transfer.Process), nil
}

func (t *Transfer) SetCallbackAddr(tpId, addr string) {
	_ = t.callbackAddr.Set(tpId, addr)
}

func (t *Transfer) CallbackAddr(tpId string) (string, error) {
	val, err := t.callbackAddr.Get(tpId)
	if err != nil {
		return ``, stores.QueryFailed(collCallbackAddr, `Get`, err)
	}

	if val == nil {
		return ``, stores.InvalidKey(tpId)
	}

	return val.(string), nil
}

func (t *Transfer) SetAgreementId(tpId, agreementId string) {
	_ = t.agreements.Set(tpId, agreementId)
}

func (t *Transfer) AgreementId(tpId string) (string, error) {
	val, err := t.agreements.Get(tpId)
	if err != nil {
		return ``, stores.QueryFailed(collAgreements, `Get`, err)
	}

	if val == nil {
		return ``, stores.InvalidKey(tpId)
	}

	return val.(string), nil
}

func (t *Transfer) SetDataSinkAddress(tpId string, addr transfer.Address) {
	_ = t.dataAddr.Set(tpId, addr)
}

func (t *Transfer) DataSinkAddress(tpId string) (transfer.Address, error) {
	val, err := t.dataAddr.Get(tpId)
	if err != nil {
		return transfer.Address{}, stores.QueryFailed(collDataAddr, `Get`, err)
	}

	if val == nil {
		return transfer.Address{}, stores.InvalidKey(tpId)
	}

	return val.(transfer.Address), nil
}

func (t *Transfer) UpdateState(tpId string, s transfer.State) error {
	process, err := t.Process(tpId)
	if err != nil {
		return stores.QueryFailed(collTransfer, `Process`, err)
	}

	process.State = s
	t.AddProcess(tpId, process)
	return nil
}
