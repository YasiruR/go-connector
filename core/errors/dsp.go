package errors

import "fmt"

const (
	incompatibleValues = "received invalid value for %s (received: %s, expected: %s)"
)

func IncompatibleValues(name, received, expected string) error {
	return fmt.Errorf(incompatibleValues, name, received, expected)
}
