package dsp

import (
	"fmt"
	"github.com/YasiruR/connector/domain/api/dsp/http/transfer"
	"github.com/YasiruR/connector/domain/core"
)

type TransferError struct {
	err  error
	Body transfer.Error
}

func (t TransferError) Error() string {
	return t.err.Error()
}

func TransferStateError(provPid, consPid, operation string, currentState transfer.State) error {
	return TransferError{
		err: fmt.Errorf("incompatible state (current state: %s, operation: %s)", currentState, operation),
		Body: transfer.Error{
			Ctx:     core.Context,
			Type:    transfer.MsgTypeError,
			ProvPId: provPid,
			ConsPId: consPid,
			Code:    "4001",
			Reason: []interface{}{
				fmt.Sprintf("current state (%s) is incompatible with the requested operation", currentState),
			},
		},
	}
}

func TransferInvalidKey(provPid, consPid, store, key string, err error) error {
	return TransferError{
		err: fmt.Errorf("%s store error - %s for %s", store, err, key),
		Body: transfer.Error{
			Ctx:     core.Context,
			Type:    transfer.MsgTypeError,
			ProvPId: provPid,
			ConsPId: consPid,
			Code:    "4002",
			Reason: []interface{}{
				fmt.Sprintf("incorrect value provided for %s", key),
			},
		},
	}
}

func TransferInvalidReqBody(provPid, consPid, msgType string, err error) error {
	return TransferError{
		err: fmt.Errorf("unmarshal failed for request %s - %w", msgType, err),
		Body: transfer.Error{
			Ctx:     core.Context,
			Type:    transfer.MsgTypeError,
			ProvPId: provPid,
			ConsPId: consPid,
			Code:    "4002",
			Reason: []interface{}{
				fmt.Sprintf("invalid request body for %s", msgType),
			},
		},
	}
}

func TransferReqParseError(msgType string, err error) error {
	return TransferError{
		err: fmt.Errorf("reading request body failed for %s - %s", msgType, err),
		Body: transfer.Error{
			Ctx:    core.Context,
			Type:   transfer.MsgTypeError,
			Code:   "4003",
			Reason: []interface{}{"request parser error"},
		},
	}
}

func TransferPathParamError(param string) error {
	return TransferError{
		err: fmt.Errorf("path parameter (%s) not found in request", param),
		Body: transfer.Error{
			Ctx:    core.Context,
			Type:   transfer.MsgTypeError,
			Code:   "4004",
			Reason: []interface{}{"required path parameter not found"},
		},
	}
}

func TransferWriteAckError(provPid, consPid, msgType string, err error) error {
	return TransferError{
		err: fmt.Errorf("%s for %s", err, msgType),
		Body: transfer.Error{
			Ctx:     core.Context,
			Type:    transfer.MsgTypeError,
			ProvPId: provPid,
			ConsPId: consPid,
			Code:    "4005",
			Reason:  []interface{}{fmt.Sprintf("internal error")},
		},
	}
}
