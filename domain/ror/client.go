package ror

import (
	"fmt"
)

func IncompatibleState(protocol, received, expected string) error {
	return ClientError{
		err: fmt.Errorf("incompatible state (protocol: %s, received: %s, expected: %s)",
			protocol, received, expected),
		Code: 2001,
		Msg:  fmt.Sprintf("existing state of the %s protocol process does not comply with the requested operation", protocol),
	}
}

func IncompatibleValues(name, received, expected string) error {
	return ClientError{
		err: fmt.Errorf("incompatible value (attribute: %s, received: %s, expected: %s)",
			name, received, expected),
		Code: 2002,
		Msg:  fmt.Sprintf("incorrect value received for '%s'", name),
	}
}

func MissingRequiredAttr(attr, reason string) error {
	return ClientError{
		err:  fmt.Errorf("required attribute was not provided (attribute: %s, reason: %s)", attr, reason),
		Code: 2003,
		Msg:  fmt.Sprintf("'%s' attribute is required but not provided", attr),
	}
}
