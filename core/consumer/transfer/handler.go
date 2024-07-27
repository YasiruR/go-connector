package transfer

import (
	"fmt"
	"github.com/YasiruR/connector/domain/api/dsp/http/transfer"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/stores"
)

type Handler struct {
	tpStore stores.Transfer
	log     pkg.Log
}

func NewHandler(tpStore stores.Transfer, log pkg.Log) *Handler {
	return &Handler{tpStore: tpStore, log: log}
}

func (h *Handler) HandleTransferStart(sr transfer.StartRequest) (transfer.Ack, error) {
	tp, err := h.tpStore.GetProcess(sr.ConsPId)
	if err != nil {
		return transfer.Ack{}, errors.StoreFailed(stores.TypeTransfer, `GetProcess`, err)
	}

	// validate if received details are compatible with existing TP

	if err = h.tpStore.UpdateState(sr.ConsPId, transfer.StateStarted); err != nil {
		return transfer.Ack{}, errors.StoreFailed(stores.TypeTransfer, `UpdateState`, err)
	}

	tp.State = transfer.StateStarted
	h.log.Info(fmt.Sprintf("updated transfer process (id: %s, state: %s)", sr.ConsPId, transfer.StateStarted))
	return transfer.Ack(tp), nil
}
