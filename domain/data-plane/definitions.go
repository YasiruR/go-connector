package data_plane

type ExchangerType string

const (
	TypePostgresql ExchangerType = "dspace:postgresql+push"
)

type Database struct {
	Endpoint    string
	Name        string
	AccessToken string
	Credentials
}

type Credentials struct {
	User     string
	Password string
}
