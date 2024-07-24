package errors

import "fmt"

// Errors related to Data Space Protocols are defined here

const (
	incompatibleValues = "received invalid value for %s (received: %s, expected: %s)"
)

func IncompatibleValues(name, received, expected string) error {
	return fmt.Errorf(incompatibleValues, name, received, expected)
}
