package errors

import "fmt"

// Errors related to Data Space Protocols are defined here

const (
	incompatibleValues  = "received invalid value for %s (received: %s, expected: %s)"
	missingRequiredAttr = "required attribute was not provided (attribute: %s, reason: %s)"
	dspControllerFailed = "DSP controller failed (role: %s, function: %s) - %s"
	dspHandlerFailed    = "DSP handler failed (endpoint: %s, role: %s) - %s"
)

func IncompatibleValues(name, received, expected string) error {
	return fmt.Errorf(incompatibleValues, name, received, expected)
}

func MissingRequiredAttr(attr, reason string) error {
	return fmt.Errorf(missingRequiredAttr, attr, reason)
}

func DSPControllerFailed(role, function string, err error) error {
	return fmt.Errorf(dspControllerFailed, role, function, err)
}

func DSPHandlerFailed(endpoint, role string, err error) error {
	return fmt.Errorf(dspHandlerFailed, endpoint, role, err)
}
