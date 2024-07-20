package pkg

type IAM interface {
	Register()
	Verify()
}

type Database interface {
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
