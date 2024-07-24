package errors

import "fmt"

const (
	storeFailed = "store failed (store: %s, function: %s) - %s"
	queryFailed = "query failed (collection: %s, query: %s) - %s"
	invalidKey  = "requested key does not exist (key: %s)"
)

func StoreFailed(name, function string, err error) error {
	return fmt.Errorf(storeFailed, name, function, err)
}

func QueryFailed(collection, query string, err error) error {
	return fmt.Errorf(queryFailed, collection, query, err)
}

func InvalidKey(key string) error {
	return fmt.Errorf(invalidKey, key)
}
