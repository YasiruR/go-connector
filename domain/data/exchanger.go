package data

type Exchanger interface {
	NewToken(participantId, datasetId string) string
	PushWithCredentials(datasetId, host, db string, c Credentials) error
	Pull(datasetId, endpoint, token string)
}
