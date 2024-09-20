package middleware

import (
	"fmt"
)

/*
	Parser errors
*/

func unmarshalError(err error) error {
	return fmt.Errorf("unmarshal error - %s", err)
}

func writeAckFailed(err error) error {
	return fmt.Errorf("creating ack response failed - %s", err)
}

func readBodyFailed(err error) error {
	return fmt.Errorf("reading body failed - %s", err)
}
