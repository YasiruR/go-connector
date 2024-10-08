package models

type Database struct {
	Name        string `json:"name"`
	Endpoint    string `json:"endpoint"`
	AccessToken string `json:"access_token"`
	Credentials
}

type Credentials struct {
	User     string `json:"user"`
	Password string `json:"password"`
}
