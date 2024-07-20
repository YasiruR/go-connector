package errors

import "fmt"

const (
	pathParamNotFound  = "path parameter not found in request (endpoint: %s, param: %s)"
	invalidRequestBody = "invalid request body (endpoint: %s) - %s"
	unmarshalError     = "unmarshal error (endpoint: %s) - %s"
	marshalError       = "marshal error (endpoint: %s) - %s"
	writeBodyError     = "writing response body failed (endpoint: %s) - %s"
	handlerFailed      = "handler failed (endpoint: %s, handler: %s) - %s"
)

func PathParamNotFound(endpoint, param string) error {
	return fmt.Errorf(pathParamNotFound, param, endpoint)
}

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

func HandlerFailed(endpoint, handler string, err error) error {
	return fmt.Errorf(handlerFailed, endpoint, handler, err)
}
