package errors

import "fmt"

// High-level errors related to modules are defined here

const (
	initFailedMsg = "initializing failed (module: %s) - %s"
	dspFailed     = "data space protocol failed (module: %s, function: %s) - %s"
	apiFailed     = "api failed (type: %s) - %s"
	handlerFailed = "handler failed (endpoint: %s, role: %s) - %s"
)

func InitFailed(module string, err error) error {
	return fmt.Errorf(initFailedMsg, module, err)
}

func DSPFailed(role, function string, err error) error {
	return fmt.Errorf(dspFailed, role, function, err)
}

func APIFailed(api string, err error) error {
	return fmt.Errorf(apiFailed, api, err)
}

func HandlerFailed(endpoint, role string, err error) error {
	return fmt.Errorf(handlerFailed, endpoint, role, err)
}
