package stores

import (
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/dsp/http/transfer"
	"github.com/YasiruR/connector/domain/pkg"
)

type Transfer struct {
	store pkg.Collection
}

func NewTransferStore(plugins domain.Plugins) *Transfer {
	plugins.Log.Info("initialized transfer process store")
	return &Transfer{store: plugins.Database.NewCollection()}
}

func (t *Transfer) Set(tpId string, val transfer.Process) {
	_ = t.store.Set(tpId, val)
}

func (t *Transfer) GetProcess(id string) (transfer.Process, error) {
	return transfer.Process{}, nil
}
