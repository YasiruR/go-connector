package models

type DCATDistribution struct {
	Id          string
	EndpointURL string
}

type DCATDataset struct {
	Id          string
	Title       string
	Keyword     string
	Description []struct {
		Value    string
		Language string
	}
}
