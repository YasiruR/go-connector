package ror

import (
	"errors"
	"fmt"
)

func DSPControllerFailed(role, function string, err error) error {
	if errors.As(err, &ClientError{}) {
		return fmt.Errorf("DSP controller failed (role: %s, function: %s) - %w", role, function, err)
	}
	return fmt.Errorf("DSP controller failed (role: %s, function: %s) - %s", role, function, err)
}

func DSPHandlerFailed(role, endpoint string, err error) error {
	if errors.As(err, &ClientError{}) {
		return fmt.Errorf("DSP handler failed (role: %s, endpoint: %s) - %w", role, endpoint, err)
	}
	return fmt.Errorf("DSP handler failed (role: %s, endpoint: %s) - %s", role, endpoint, err)
}

func CustomFuncError(function string, err error) error {
	if errors.As(err, &ClientError{}) {
		return fmt.Errorf("function %s failed - %w", function, err)
	}
	return fmt.Errorf("function %s failed - %s", function, err)
}
