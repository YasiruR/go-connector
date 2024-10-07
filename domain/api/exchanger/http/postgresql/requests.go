package postgresql

type PullRequest struct {
	Database string `json:"database"`
	Endpoint string `json:"endpoint"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type PushRequest struct {
	Database string `json:"database"`
	Endpoint string `json:"endpoint"`
	User     string `json:"user"`
	Password string `json:"password"`
}
