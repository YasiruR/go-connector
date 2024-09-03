package dsp

import (
	"fmt"
	"github.com/YasiruR/connector/domain/api/dsp/http/negotiation"
	"github.com/YasiruR/connector/domain/core"
)

type NegotiationError struct {
	err  error
	Body negotiation.Error
}

func (n NegotiationError) Error() string {
	return n.err.Error()
}

func NegotiationStateError(provPid, consPid, operation string, currentState negotiation.State) error {
	return NegotiationError{
		err: fmt.Errorf("incompatible state (current state: %s, operation: %s)", currentState, operation),
		Body: negotiation.Error{
			Ctx:     core.Context,
			Type:    negotiation.MsgTypeError,
			ProvPId: provPid,
			ConsPId: consPid,
			Code:    "3001",
			Reason: []interface{}{
				fmt.Sprintf("current state (%s) is incompatible with the requested operation", currentState),
			},
			Desc: nil,
		},
	}
}

func NegotiationValError(provPid, consPid, key, current, received string) error {
	return NegotiationError{
		err: fmt.Errorf("incompatible values for %s (current: %s, received: %s)", key, current, received),
		Body: negotiation.Error{
			Ctx:     core.Context,
			Type:    negotiation.MsgTypeError,
			ProvPId: provPid,
			ConsPId: consPid,
			Code:    "3002",
			Reason: []interface{}{
				fmt.Sprintf("received value for '%s' does not match with the stored value", key),
			},
			Desc: nil,
		},
	}
}

func NegotiationInvalidKey(provPid, consPid, store, key string, err error) error {
	return NegotiationError{
		err: fmt.Errorf("%s store error - %s for %s", store, err, key),
		Body: negotiation.Error{
			Ctx:     core.Context,
			Type:    negotiation.MsgTypeError,
			ProvPId: provPid,
			ConsPId: consPid,
			Code:    "3003",
			Reason: []interface{}{
				fmt.Sprintf("incorrect value provided for %s", key),
			},
			Desc: nil,
		},
	}
}

func NegotiationInvalidReqBody(provPid, consPid, msgType string, err error) error {
	return NegotiationError{
		err: fmt.Errorf("unmarshal failed for request %s - %w", msgType, err),
		Body: negotiation.Error{
			Ctx:     core.Context,
			Type:    negotiation.MsgTypeError,
			ProvPId: provPid,
			ConsPId: consPid,
			Code:    "3004",
			Reason: []interface{}{
				fmt.Sprintf("invalid request body for %s", msgType),
			},
			Desc: nil,
		},
	}
}

func NegotiationReqParseError(msgType string, err error) error {
	return NegotiationError{
		err: fmt.Errorf("reading request body failed for %s - %s", msgType, err),
		Body: negotiation.Error{
			Ctx:    core.Context,
			Type:   negotiation.MsgTypeError,
			Code:   "3005",
			Reason: []interface{}{"request parser error"},
			Desc:   nil,
		},
	}
}

func NegotiationPathParamError(param string) error {
	return NegotiationError{
		err: fmt.Errorf("path parameter (%s) not found in request", param),
		Body: negotiation.Error{
			Ctx:    core.Context,
			Type:   negotiation.MsgTypeError,
			Code:   "3006",
			Reason: []interface{}{"required path parameter not found"},
			Desc:   nil,
		},
	}
}

func NegotiationWriteAckError(provPid, consPid, msgType string, err error) error {
	return NegotiationError{
		err: fmt.Errorf("%s for %s", err, msgType),
		Body: negotiation.Error{
			Ctx:     core.Context,
			Type:    negotiation.MsgTypeError,
			ProvPId: provPid,
			ConsPId: consPid,
			Code:    "3007",
			Reason:  []interface{}{"internal error"},
			Desc:    nil,
		},
	}
}
