package external

import (
	"fmt"
)

/*
	Errors that should be returned to the client through controllers associated with gateway API
*/

type GatewayError struct {
	err  error
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

func (g GatewayError) Error() string {
	return g.err.Error()
}

func StateError(protocol, received, expected string) error {
	return GatewayError{
		err: fmt.Errorf("incompatible state (protocol: %s, received: %s, expected: %s)",
			protocol, received, expected),
		Code: 1001,
		Msg:  fmt.Sprintf("existing state of the %s protocol process does not comply with the requested operation", protocol),
	}
}

func InvalidAckError(requestType, reason string, ack any) error {
	return GatewayError{
		err:  fmt.Errorf("received invalid acknowledgement (request: %v, reason: %s, ack: %v)", requestType, reason, ack),
		Code: 1002,
		Msg:  fmt.Sprintf("received an invalid acknowledgement due to %s", reason),
	}
}

func MissingAttrError(attr, reason string) error {
	return GatewayError{
		err:  fmt.Errorf("required attribute was not provided (attribute: %s, reason: %s)", attr, reason),
		Code: 1003,
		Msg:  fmt.Sprintf("'%s' attribute is required but not provided", attr),
	}
}

func InvalidKeyError(store, key string, err error) error {
	return GatewayError{
		err:  fmt.Errorf("%s store error - %s for %s", store, err, key),
		Code: 1004,
		Msg:  fmt.Sprintf("requested value for '%s' does not exist", key),
	}
}

func InvalidReqBody(msgType string, err error) error {
	return GatewayError{
		err:  fmt.Errorf("unmarshal failed for request %s - %w", msgType, err),
		Code: 1005,
		Msg:  fmt.Sprintf("invalid request body for %s", msgType),
	}
}

func IncompatibleReqBody(reason string) error {
	return GatewayError{
		err:  fmt.Errorf("incompatible request body (reason: %s)", reason),
		Code: 1006,
		Msg:  reason,
	}
}

func UnmarshalError(msgType string, err error) error {
	return GatewayError{
		err:  fmt.Errorf("unmarshal error (message: %s) - %s", msgType, err),
		Code: 1007,
		Msg:  fmt.Sprintf("received an invalid response body for %s", msgType),
	}
}

func MarshalError(msgType string, err error) error {
	return GatewayError{
		err:  fmt.Errorf("marshal error (message: %s) - %s", msgType, err),
		Code: 1008,
		Msg:  fmt.Sprintf("internal data parser error"),
	}
}

func ParseError(msgType string, err error) error {
	return GatewayError{
		err:  fmt.Errorf("reading request body failed for %s - %s", msgType, err),
		Code: 1009,
		Msg:  fmt.Sprintf("internal data parser error"),
	}
}

func PathParamError(param string) error {
	return GatewayError{
		err:  fmt.Errorf("path parameter (%s) not found in request", param),
		Code: 1010,
		Msg:  fmt.Sprintf("required path parameter not found"),
	}
}

func InvalidStatusCode(code int) error {
	return GatewayError{
		err:  fmt.Errorf("received invalid status code (%d)", code),
		Code: 1011,
		Msg:  fmt.Sprintf("received invalid status code (%d)", code),
	}
}

func URLStringError(dest any) error {
	return GatewayError{
		err:  fmt.Errorf("invalid url (received: %v, requires a string)", dest),
		Code: 1012,
		Msg:  fmt.Sprintf("internal error"),
	}
}

func SendFailed(endpoint, method string, err error) error {
	return GatewayError{
		err:  fmt.Errorf("sending message failed (endpoint: %s, method: %s) - %s", endpoint, method, err),
		Code: 1013,
		Msg:  fmt.Sprintf("sending message to %s failed", endpoint),
	}
}

func WriteAckError(msgType string, err error) error {
	return GatewayError{
		err:  fmt.Errorf("%s for %s", err, msgType),
		Code: 1014,
		Msg:  fmt.Sprintf("internal error"),
	}
}
