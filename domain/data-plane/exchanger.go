package data_plane

// several approaches to trigger the transfer
//	- synchronously by the control plane
// 	- asynchronously in a go-routine
//	- manually and separately

type Exchanger interface {
	NewToken(participantId, datasetId string) string
	PushWithCredentials(et ExchangerType, dest Database) error
	PullWithCredentials(src Database) error
}
