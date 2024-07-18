package core

type IAM interface {
	Register()
	Verify()
}

type Store interface {
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
