package postgresql

import data_plane "github.com/YasiruR/go-connector/domain/data-plane"

const (
	PullEndpoint = `/exchanger/` + data_plane.DatabasePostgresql + `/pull`
	PushEndpoint = `/exchanger/` + data_plane.DatabasePostgresql + `/push`
)
