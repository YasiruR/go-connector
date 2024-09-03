package core

import (
	"errors"
	"fmt"
	"github.com/YasiruR/connector/domain/errors/dsp"
	"github.com/YasiruR/connector/domain/errors/external"
)

func ModuleInitFailed(module string, err error) error {
	return fmt.Errorf("initializing failed (module: %s) - %s", module, err)
}

func DSPControllerFailed(role, function string, err error) error {
	if errors.As(err, &external.GatewayError{}) || errors.As(err, &dsp.CatalogError{}) {
		return fmt.Errorf("DSP controller failed (role: %s, function: %s) - %w", role, function, err)
	}
	return fmt.Errorf("DSP controller failed (role: %s, function: %s) - %s", role, function, err)
}

func DSPHandlerFailed(role, endpoint string, err error) error {
	if errors.As(err, &external.GatewayError{}) || errors.As(err, &dsp.CatalogError{}) {
		return fmt.Errorf("DSP handler failed (role: %s, endpoint: %s) - %w", role, endpoint, err)
	}
	return fmt.Errorf("DSP handler failed (role: %s, endpoint: %s) - %s", role, endpoint, err)
}

func StoreFailed(name, function string, err error) error {
	if errors.As(err, &external.GatewayError{}) {
		return fmt.Errorf("store error (store: %s, function: %s) - %w", name, function, err)
	}
	return fmt.Errorf("store error (store: %s, function: %s) - %s", name, function, err)
}

func CustomFuncError(function string, err error) error {
	if errors.As(err, &external.GatewayError{}) {
		return fmt.Errorf("function %s failed - %w", function, err)
	}
	return fmt.Errorf("function %s failed - %s", function, err)
}
