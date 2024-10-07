package pkg

import "context"

const (
	TypeURN    = `URNService`
	TypeClient = `Client`
)

type IAM interface {
	Register()
	Verify()
}

type PolicyEngine interface {
	ValidateOffer()
}

// Store contains one or more Collection to support data-plane storage required
// by the connector
type Store interface {
	NewCollection() Collection
}

// Collection provides an isolated storage for a single context. For example,
// Collection can be an SQL table, a NoSQL collection or an in-memory map.
type Collection interface {
	Get(key string) (interface{}, error)
	GetAll() ([]any, error)
	Set(key string, value interface{}) error
}

type Client interface {
	Send(data []byte, destination any) (response []byte, err error)
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
	FatalContext(ctx context.Context, message interface{}, params ...interface{})
	ErrorContext(ctx context.Context, message interface{}, params ...interface{})
	WarnContext(ctx context.Context, message interface{}, params ...interface{})
	DebugContext(ctx context.Context, message interface{}, params ...interface{})
	InfoContext(ctx context.Context, message interface{}, params ...interface{})
	TraceContext(ctx context.Context, message interface{}, params ...interface{})
}
