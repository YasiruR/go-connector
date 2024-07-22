package pkg

const (
	TypeDatabase = `Database`
	TypeURN      = `URNService`
	TypeClient   = `Client`
)

type IAM interface {
	Register()
	Verify()
}

// Database contains one or more Collection to support data storage required
// by the connector
type Database interface {
	NewDataStore() Collection
}

// Collection provides an isolated storage for a single context. For example,
// Collection can be an SQL table, a NoSQL collection or an in-memory map.
type Collection interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}

type Client interface {
	Send(data []byte, destination string) (response []byte, err error)
}

type URNService interface {
	NewURN() (string, error)
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
