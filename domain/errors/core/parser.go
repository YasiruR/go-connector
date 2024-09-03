package core

import (
	"errors"
	"fmt"
)

var (
	TypeUnmarshalError = errors.New("unmarshal error")
)

func UnmarshalError(err error) error {
	return fmt.Errorf("%w - %s", TypeUnmarshalError, err)
}

func WriteAckFailed(err error) error {
	return fmt.Errorf("creating ack response failed - %s", err)
}

func ReadBodyFailed(err error) error {
	return fmt.Errorf("reading body failed - %s", err)
}
