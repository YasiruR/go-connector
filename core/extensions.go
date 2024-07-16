package core

type IAM interface {
	Register()
	Verify()
}

type Store interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}
