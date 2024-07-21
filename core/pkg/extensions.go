package pkg

const (
	TypeDatabase   = `Database`
	TypeURN        = `URN`
	TypeHTTPClient = `HTTPClient`
)

type IAM interface {
	Register()
	Verify()
}

// Database contains one or more DataStore to support data storage required
// by the connector
type Database interface {
	NewDataStore() DataStore
}

// DataStore provides an isolated storage for a single context. For example,
// DataStore can be a table, a collection or an in-memory map.
type DataStore interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}

type HTTPClient interface {
	Post(url string, data []byte) (statusCode int, resData []byte, err error)
}

type URN interface {
	New() (string, error)
	Validate(urn string) bool
}

type Log interface {
	Fatal(message interface{}, params ...interface{})
	Error(message interface{}, params ...interface{})
	Warn(message interface{}, params ...interface{})
	Debug(message interface{}, params ...interface{})
	Info(message interface{}, params ...interface{})
	Trace(message interface{}, params ...interface{})
}
