package transfer

import (
	"fmt"
	"github.com/YasiruR/connector/domain"
	"github.com/YasiruR/connector/domain/api/dsp/http/transfer"
	"github.com/YasiruR/connector/domain/core"
	"github.com/YasiruR/connector/domain/errors"
	"github.com/YasiruR/connector/domain/pkg"
	"github.com/YasiruR/connector/domain/stores"
)

type Handler struct {
	urn      pkg.URNService
	agrStore stores.AgreementStore
	tpStore  stores.TransferStore
	log      pkg.Log
}

func NewHandler(stores domain.Stores, plugins domain.Plugins) *Handler {
	return &Handler{
		agrStore: stores.AgreementStore,
		tpStore:  stores.TransferStore,
		urn:      plugins.URNService,
		log:      plugins.Log,
	}
}

func (h *Handler) HandleGetProcess(tpId string) (transfer.Ack, error) {
	tp, err := h.tpStore.Process(tpId)
	if err != nil {
		return transfer.Ack{}, errors.StoreFailed(stores.TypeAgreement, `Process`, err)
	}

	return transfer.Ack(tp), nil
}

func (h *Handler) HandleTransferRequest(tr transfer.Request) (transfer.Ack, error) {
	// validate agreement
	_, err := h.agrStore.Agreement(tr.AgreementId)
	if err != nil {
		return transfer.Ack{}, errors.StoreFailed(stores.TypeAgreement, `Agreement`, err)
	}

	tpId, err := h.urn.NewURN()
	if err != nil {
		return transfer.Ack{}, errors.PkgFailed(pkg.TypeURN, `New`, err)
	}

	ack := transfer.Ack{
		Ctx:     core.Context,
		Type:    transfer.MsgTypeProcess,
		ProvPId: tpId,
		ConsPId: tr.ConsPId,
		State:   transfer.StateRequested,
	}

	h.tpStore.AddProcess(tpId, transfer.Process(ack))
	h.tpStore.SetCallbackAddr(tpId, tr.CallbackAddr)
	h.log.Trace("stored transfer process", ack)
	h.log.Debug(fmt.Sprintf("provider handler updated transfer process (id: %s, state: %s)", tpId, transfer.StateRequested))
	return ack, nil
}

func (h *Handler) HandleTransferSuspension(sr transfer.SuspendRequest) (transfer.Ack, error) {
	tp, err := h.tpStore.Process(sr.ProvPId)
	if err != nil {
		return transfer.Ack{}, errors.StoreFailed(stores.TypeTransfer, `Process`, err)
	}

	// validate tp

	if tp.State != transfer.StateStarted {
		return transfer.Ack{}, errors.IncompatibleValues(`state`, string(tp.State), string(transfer.StateStarted))
	}

	if err = h.tpStore.UpdateState(sr.ProvPId, transfer.StateSuspended); err != nil {
		return transfer.Ack{}, errors.StoreFailed(stores.TypeTransfer, `UpdateState`, err)
	}

	tp.State = transfer.StateSuspended
	h.log.Debug(fmt.Sprintf("provider handler updated transfer process (id: %s, state: %s)", sr.ProvPId, transfer.StateSuspended))
	return transfer.Ack(tp), nil
}

func (h *Handler) HandleTransferStart(sr transfer.StartRequest) (transfer.Ack, error) {
	tp, err := h.tpStore.Process(sr.ProvPId)
	if err != nil {
		return transfer.Ack{}, errors.StoreFailed(stores.TypeTransfer, `Process`, err)
	}

	// validate if received details are compatible with existing TP

	if tp.State != transfer.StateSuspended {
		return transfer.Ack{}, errors.IncompatibleValues(`state`, string(tp.State), string(transfer.StateSuspended))
	}

	if err = h.tpStore.UpdateState(sr.ProvPId, transfer.StateStarted); err != nil {
		return transfer.Ack{}, errors.StoreFailed(stores.TypeTransfer, `UpdateState`, err)
	}

	tp.State = transfer.StateStarted
	h.log.Debug(fmt.Sprintf("provider handler updated transfer process (id: %s, state: %s)", sr.ProvPId, transfer.StateStarted))
	return transfer.Ack(tp), nil
}

func (h *Handler) HandleTransferCompletion(cr transfer.CompleteRequest) (transfer.Ack, error) {
	tp, err := h.tpStore.Process(cr.ProvPId)
	if err != nil {
		return transfer.Ack{}, errors.StoreFailed(stores.TypeTransfer, `Process`, err)
	}

	if tp.State != transfer.StateStarted {
		return transfer.Ack{}, errors.IncompatibleValues(`state`, string(tp.State), string(transfer.StateStarted))
	}

	if err = h.tpStore.UpdateState(cr.ProvPId, transfer.StateCompleted); err != nil {
		return transfer.Ack{}, errors.StoreFailed(stores.TypeTransfer, `UpdateState`, err)
	}

	tp.State = transfer.StateCompleted
	h.log.Info(fmt.Sprintf("data exchange process completed successfully (id: %s)", cr.ProvPId))
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
	h.log.Info(fmt.Sprintf("data exchange process was terminated by consumer (id: %s, reasons: %v)",
		tr.ProvPId, tr.Reason))
	return transfer.Ack(tp), nil
}
