package data

const (
	PostgreSQLPush = `dspace:postgresql+push`
)

type Database struct {
	Endpoint string
	Name     string
	Credentials
}

type Token string

type Credentials struct {
	User     string
	Password string
}
