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
	agrStore stores.Agreement
	tpStore  stores.Transfer
	log      pkg.Log
}

func NewHandler(stores domain.Stores, plugins domain.Plugins) *Handler {
	return &Handler{
		agrStore: stores.Agreement,
		tpStore:  stores.Transfer,
		urn:      plugins.URNService,
		log:      plugins.Log,
	}
}

func (h *Handler) HandleTransferRequest(tr transfer.Request) (transfer.Ack, error) {
	// validate agreement
	_, err := h.agrStore.Get(tr.AgreementId)
	if err != nil {
		return transfer.Ack{}, errors.StoreFailed(stores.TypeAgreement, `Get`, err)
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

	h.tpStore.Set(tpId, transfer.Process(ack))
	h.tpStore.SetCallbackAddr(tpId, tr.CallbackAddr)
	h.log.Trace("stored transfer process", ack)
	h.log.Info(fmt.Sprintf("updated transfer process (id: %s, state: %s)", tpId, transfer.StateRequested))
	return ack, nil
}

func (h *Handler) HandleTransferSuspension(sr transfer.SuspendRequest) (transfer.Ack, error) {
	tp, err := h.tpStore.GetProcess(sr.ProvPId)
	if err != nil {
		return transfer.Ack{}, errors.StoreFailed(stores.TypeTransfer, `GetProcess`, err)
	}

	if err := h.tpStore.UpdateState(sr.ProvPId, transfer.StateSuspended); err != nil {
		return transfer.Ack{}, errors.StoreFailed(stores.TypeTransfer, `UpdateState`, err)
	}

	tp.State = transfer.StateSuspended
	h.log.Info(fmt.Sprintf("updated transfer process (id: %s, state: %s)", sr.ProvPId, transfer.StateSuspended))
	return transfer.Ack(tp), nil
}
