package errors

import "fmt"

const (
	storeFailed = "store failed (store: %s, function: %s) - %s"
	queryFailed = "get query failed (query: %s) - %s"
)

func StoreFailed(name, function string, err error) error {
	return fmt.Errorf(storeFailed, name, function, err)
}

func QueryFailed(query string, err error) error {
	return fmt.Errorf(queryFailed, query, err)
}
