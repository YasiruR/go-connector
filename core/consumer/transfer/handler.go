package transfer

import (
	"fmt"
	"github.com/YasiruR/connector/domain/api/dsp/http/transfer"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/stores"
)

type Handler struct {
	tpStore stores.TransferStore
	log     pkg.Log
}

func NewHandler(tpStore stores.TransferStore, log pkg.Log) *Handler {
	return &Handler{tpStore: tpStore, log: log}
}

func (h *Handler) HandleTransferStart(sr transfer.StartRequest) (transfer.Ack, error) {
	tp, err := h.tpStore.Process(sr.ConsPId)
	if err != nil {
		return transfer.Ack{}, errors.StoreFailed(stores.TypeTransfer, `Process`, err)
	}

	// validate if received details are compatible with existing TP

	if tp.State != transfer.StateRequested && tp.State != transfer.StateSuspended {
		return transfer.Ack{}, errors.IncompatibleValues(`state`, string(tp.State),
			string(transfer.StateRequested)+" or "+string(transfer.StateSuspended))
	}

	if err = h.tpStore.UpdateState(sr.ConsPId, transfer.StateStarted); err != nil {
		return transfer.Ack{}, errors.StoreFailed(stores.TypeTransfer, `UpdateState`, err)
	}

	tp.State = transfer.StateStarted
	h.log.Debug(fmt.Sprintf("consumer handler updated transfer process (id: %s, state: %s)", sr.ConsPId, transfer.StateStarted))
	return transfer.Ack(tp), nil
}

func (h *Handler) HandleTransferSuspension(sr transfer.SuspendRequest) (transfer.Ack, error) {
	tp, err := h.tpStore.Process(sr.ConsPId)
	if err != nil {
		return transfer.Ack{}, errors.StoreFailed(stores.TypeTransfer, `Process`, err)
	}

	if tp.State != transfer.StateStarted {
		return transfer.Ack{}, errors.IncompatibleValues(`state`, string(tp.State), string(transfer.StateStarted))
	}

	if err := h.tpStore.UpdateState(sr.ConsPId, transfer.StateSuspended); err != nil {
		return transfer.Ack{}, errors.StoreFailed(stores.TypeTransfer, `UpdateState`, err)
	}

	tp.State = transfer.StateSuspended
	h.log.Debug(fmt.Sprintf("consumer handler updated transfer process (id: %s, state: %s)", sr.ConsPId, transfer.StateSuspended))
	return transfer.Ack(tp), nil
}

func (h *Handler) HandleTransferCompletion(cr transfer.CompleteRequest) (transfer.Ack, error) {
	tp, err := h.tpStore.Process(cr.ConsPId)
	if err != nil {
		return transfer.Ack{}, errors.StoreFailed(stores.TypeTransfer, `Process`, err)
	}

	if err = h.tpStore.UpdateState(cr.ConsPId, transfer.StateCompleted); err != nil {
		return transfer.Ack{}, errors.StoreFailed(stores.TypeTransfer, `UpdateState`, err)
	}

	tp.State = transfer.StateCompleted
	h.log.Info(fmt.Sprintf("data exchange process completed successfully (id: %s)", cr.ConsPId))
	return transfer.Ack(tp), nil
}

func (h *Handler) HandleTransferTermination(tr transfer.TerminateRequest) (transfer.Ack, error) {
	tp, err := h.tpStore.Process(tr.ProvPId)
	if err != nil {
		return transfer.Ack{}, errors.StoreFailed(stores.TypeTransfer, `Process`, err)
	}

	if tp.State != transfer.StateRequested && tp.State != transfer.StateStarted && tp.State != transfer.StateSuspended {
		return transfer.Ack{}, errors.IncompatibleValues(`state`, string(tp.State),
			string(transfer.StateStarted)+" or "+string(transfer.StateStarted)+" or "+string(transfer.StateSuspended))
	}

	if err = h.tpStore.UpdateState(tr.ProvPId, transfer.StateTerminated); err != nil {
		return transfer.Ack{}, errors.StoreFailed(stores.TypeTransfer, `UpdateState`, err)
	}

	tp.State = transfer.StateTerminated
	h.log.Info(fmt.Sprintf("data exchange process was terminated by provider (id: %s, reasons: %v)",
		tr.ProvPId, tr.Reason))
	return transfer.Ack(tp), nil
}
