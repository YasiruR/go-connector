package ror

import "fmt"

func InvalidAck(requestType, reason string, ack any) error {
	return fmt.Errorf("received invalid acknowledgement (request: %v, reason: %s, ack: %v)", requestType, reason, ack)
}
