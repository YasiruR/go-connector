package data_plane

import "github.com/YasiruR/go-connector/domain/models"

type ExchangerType string

// supported databases for transfers
const (
	DatabasePostgresql = `postgresql`
)

const (
	PullFilePrefix = "backups/pull/"
	PushFilePrefix = "backups/push/"
)

// several approaches to trigger the transfer
//	- synchronously by the control plane
// 	- asynchronously in a go-routine
//	- manually and separately

type Exchanger interface {
	NewToken(participantId, datasetId string) string
	PushWithCredentials(db string, dest models.Database) error
	PullWithCredentials(db string, src models.Database) error
}
