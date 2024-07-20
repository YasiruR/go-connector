package errors

import "fmt"

const (
	initFailedMsg = "initializing failed (module: %s) - %s"
	dspFailed     = "data space protocol failed (module: %s, function: %s) - %s"
)

func InitFailed(module string, err error) error {
	return fmt.Errorf(initFailedMsg, module, err)
}

func DSPFailed(role, function string, err error) error {
	return fmt.Errorf(dspFailed, role, function, err)
}
