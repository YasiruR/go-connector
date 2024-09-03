package core

import "fmt"

func NewURNFailed(key string, err error) error {
	return fmt.Errorf("generating new URN for '%s' failed - %s", key, err)
}

func ClientSendError(err error) error {
	return fmt.Errorf("sending message failed - %w", err)
}
