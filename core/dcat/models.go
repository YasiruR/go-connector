package dcat

type Distribution struct {
	Id          string
	EndpointURL string
}

type Dataset struct {
	Id          string
	Title       string
	Keyword     string
	Description []struct {
		Value    string
		Language string
	}
}
