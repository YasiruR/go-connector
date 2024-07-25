package stores

import (
	"github.com/YasiruR/connector/domain/api/dsp/http/transfer"
)

type Transfer struct{}

func NewTransferStore() *Transfer {
	return &Transfer{}
}

func (t *Transfer) Set(tpId string, val transfer.Process) {}

func (t *Transfer) GetProcess(id string) (transfer.Process, error) {
	return transfer.Process{}, nil
}
