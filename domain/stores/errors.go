package stores

import (
	"errors"
	"fmt"
)

var TypeInvalidKey = errors.New("requested key does not exist")

func InvalidKey(val string) error {
	return fmt.Errorf("%w (%s)", TypeInvalidKey, val)
}

func QueryFailed(collection, query string, err error) error {
	return fmt.Errorf("query failed (collection: %s, query: %s) - %s", collection, query, err)
}
