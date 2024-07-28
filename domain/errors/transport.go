package errors

import "fmt"

// Errors related to parsers

const (
	invalidRequestBody = "invalid request body (endpoint: %s) - %s"
	unmarshalError     = "unmarshal error (endpoint: %s) - %s"
	marshalError       = "marshal error (endpoint: %s) - %s"
	writeBodyError     = "writing body failed (endpoint: %s) - %s"
	parseError         = "parsing request failed (endpoint: %s) - %s"
	writeAckFailure    = "ack response failed (endpoint: %s) - %s"
	writeErrFailure    = "error response failed (endpoint: %s) - %s"
)

func InvalidRequestBody(endpoint string, err error) error {
	return fmt.Errorf(invalidRequestBody, endpoint, err)
}

func UnmarshalError(endpoint string, err error) error {
	return fmt.Errorf(unmarshalError, endpoint, err)
}

func MarshalError(endpoint string, err error) error {
	return fmt.Errorf(marshalError, endpoint, err)
}

func WriteBodyError(endpoint string, err error) error {
	return fmt.Errorf(writeBodyError, endpoint, err)
}

func ParseRequestFailed(endpoint string, err error) error {
	return fmt.Errorf(parseError, endpoint, err)
}

func WriteAckFailed(endpoint string, err error) error {
	return fmt.Errorf(writeAckFailure, endpoint, err)
}

func WriteErrorFailed(endpoint string, err error) error {
	return fmt.Errorf(writeErrFailure, endpoint, err)
}

// Errors related to HTTP transport

const (
	pathParamNotFound = "path parameter not found in request (endpoint: %s, param: %s)"
	invalidStatusCode = "received an invalid status code (received: %d)"
	sendError         = "sending message failed (endpoint: %s, method: %s) - %s"
	invalidUrl        = "invalid url (received: %v, requires a string)"
)

func PathParamNotFound(endpoint, param string) error {
	return fmt.Errorf(pathParamNotFound, param, endpoint)
}

func InvalidStatusCode(received int) error {
	return fmt.Errorf(invalidStatusCode, received)
}

func SendFailed(endpoint, method string, err error) error {
	return fmt.Errorf(sendError, endpoint, method, err)
}

func InvalidURL(received any) error {
	return fmt.Errorf(invalidUrl, received)
}
