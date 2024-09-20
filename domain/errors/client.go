package errors

import "fmt"

type ErrorMessage struct {
	Message string                 `json:"message"`
	Params  map[string]interface{} `json:"params"`
	code    string
	err     error
}

func StateError(operation, currentState string) ErrorMessage {
	return ErrorMessage{
		code:    `20001`,
		Message: fmt.Sprintf("current state (%s) is incompatible with the requested operation", currentState),
		err:     fmt.Errorf("incompatible state (current state: %s, operation: %s)", currentState, operation),
	}
}

func InvalidValue(key, current, received string) ErrorMessage {
	return ErrorMessage{
		code:    `20002`,
		Message: fmt.Sprintf("received value for '%s' does not match with the stored value", key),
		err:     fmt.Errorf("incompatible values for %s (current: %s, received: %s)", key, current, received),
	}
}

func InvalidKey(store, key string, err error) ErrorMessage {
	return ErrorMessage{
		code:    `20003`,
		Message: fmt.Sprintf("incorrect value provided for %s", key),
		err:     fmt.Errorf("%s store error - %s for %s", store, err, key),
	}
}

func InvalidReqBody(msgType string, err error) ErrorMessage {
	return ErrorMessage{
		code:    `20004`,
		Message: fmt.Sprintf("invalid request body for '%s' message", msgType),
		err:     fmt.Errorf("reading request body failed for %s - %s", msgType, err),
	}
}

func PathParamNotFound(param string) ErrorMessage {
	return ErrorMessage{
		code:    `20005`,
		Message: "required path parameter not found",
		err:     fmt.Errorf("path parameter (%s) not found in request", param),
	}
}

func WriteAckError(msgType string, err error) ErrorMessage {
	return ErrorMessage{
		code:    `20006`,
		Message: "internal error",
		err:     fmt.Errorf("%s for %s", err, msgType),
	}
}

func SendFailed(err error) ErrorMessage {
	return ErrorMessage{
		code:    `20007`,
		Message: "transport error",
		err:     fmt.Errorf("sending message failed - %s", err),
	}
}

func InvalidAckError(requestType, reason string, ack any) ErrorMessage {
	return ErrorMessage{
		code:    `20008`,
		err:     fmt.Errorf("received invalid acknowledgement (request: %v, reason: %s, ack: %v)", requestType, reason, ack),
		Message: fmt.Sprintf("received an invalid acknowledgement due to %s", reason),
	}
}

func MissingAttrError(attr, reason string) ErrorMessage {
	return ErrorMessage{
		code:    `20009`,
		err:     fmt.Errorf("required attribute was not provided (attribute: %s, reason: %s)", attr, reason),
		Message: fmt.Sprintf("'%s' attribute is required but not provided", attr),
	}
}

func IncorrectReqValues(reason string) ErrorMessage {
	return ErrorMessage{
		code:    `20010`,
		err:     fmt.Errorf("incompatible request body (reason: %s)", reason),
		Message: reason,
	}
}

func UnmarshalError(msgType string, err error) ErrorMessage {
	return ErrorMessage{
		code:    `20011`,
		err:     fmt.Errorf("unmarshal error (message: %s) - %s", msgType, err),
		Message: fmt.Sprintf("received an invalid response body for %s", msgType),
	}
}

func MarshalError(msgType string, err error) ErrorMessage {
	return ErrorMessage{
		code:    `20012`,
		err:     fmt.Errorf("marshal error (message: %s) - %s", msgType, err),
		Message: fmt.Sprintf("internal data parser error"),
	}
}

func ProtocolFailed(typ string, errMsg any, err error) ErrorMessage {
	return ErrorMessage{
		code:    `20013`,
		Message: fmt.Sprintf("%s protocol failed", typ),
		Params:  map[string]interface{}{"response": errMsg},
		err:     fmt.Errorf("%s protocol failed - %s", typ, err),
	}
}
