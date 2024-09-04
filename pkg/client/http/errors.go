package http

import (
	"fmt"
)

func invalidStatusCode(code int) error {
	return fmt.Errorf("received invalid status code (%d)", code)
}

func urlStringError(dest any) error {
	return fmt.Errorf("invalid url (received: %v, requires a string)", dest)
}

func sendFailed(endpoint, method string, err error) error {
	return fmt.Errorf("sending message failed (endpoint: %s, method: %s) - %s", endpoint, method, err)
}
