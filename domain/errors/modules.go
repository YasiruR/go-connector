package errors

import "fmt"

// High-level errors related to modules are defined here

const (
	initFailedMsg = "initializing failed (module: %s) - %s"
	apiFailed     = "api failed (type: %s) - %s"
)

func InitModuleFailed(module string, err error) error {
	return fmt.Errorf(initFailedMsg, module, err)
}

func APIFailed(api string, err error) error {
	return fmt.Errorf(apiFailed, api, err)
}
