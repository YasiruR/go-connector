package errors

import (
	"fmt"
)

/*
	High level errors related to main packages of the connector
*/

func ModuleInitFailed(module string, err error) error {
	return fmt.Errorf("initializing failed (module: %s) - %s", module, err)
}

func DSPControllerFailed(role, function string, err error) error {
	return fmt.Errorf("DSP controller failed (role: %s, function: %s) - %w", role, function, err)
}

func DSPHandlerFailed(role, endpoint string, err error) error {
	return fmt.Errorf("DSP handler failed (role: %s, endpoint: %s) - %w", role, endpoint, err)
}

func StoreFailed(name, function string, err error) error {
	return fmt.Errorf("store error (store: %s, function: %s) - %w", name, function, err)
}

func PkgError(pkg, function string, err error, params ...string) error {
	return fmt.Errorf("package error (pkg: %s, function: %s) [%v] - %w", pkg, function, params, err)
}

func TransferFailed(typ, endpoint string, err error) error {
	return fmt.Errorf("transfer error (type: %s, endpoint: %s) - %w", typ, endpoint, err)
}

func CustomFuncError(function string, err error) error {
	return fmt.Errorf("function %s failed - %w", function, err)
}
